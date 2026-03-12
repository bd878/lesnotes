package postgres

import (
	"fmt"
	"time"
	"sync"
	"errors"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/bd878/gallery/server/api"
	"github.com/bd878/gallery/server/internal/logger"
)

type Dumper struct {
	pool                   *pgxpool.Pool
	messagesTableName      string
	filesTableName         string
	translationsTableName  string
	ctx                    context.Context
	cancel                 context.CancelCauseFunc
	ch                     chan *api.MessagesSnapshot
	wg                     sync.WaitGroup
}

func NewDumper(pool *pgxpool.Pool, messagesTableName, filesTableName, translationsTableName string) *Dumper {
	return &Dumper{
		pool:                    pool,
		messagesTableName:       messagesTableName,
		filesTableName:          filesTableName,
		translationsTableName:   translationsTableName,
	}
}

func (r *Dumper) Open(ctx context.Context) (ch chan *api.MessagesSnapshot, err error) {
	r.ctx, r.cancel = context.WithCancelCause(ctx)
	ch = make(chan *api.MessagesSnapshot, 100)
	r.ch = ch

	r.wg.Add(1)
	go r.runMessages()
	r.wg.Add(1)
	go r.runFiles()
	r.wg.Add(1)
	go r.runTranslations()

	go func() {
		r.wg.Wait()
		close(r.ch)
	}()

	return
}

func (r *Dumper) runMessages() {
	query := "SELECT id, text, private, name, user_id, title, created_at, updated_at FROM %s"

	defer r.wg.Done()
	defer logger.Debugln("messages dump finished")

	rows, err := r.pool.Query(r.ctx, r.messagesTable(query))
	if err != nil {
		logger.Errorln(err)
		r.cancel(err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		message := &api.MessageSnapshotItem{}

		var createdAt, updatedAt *time.Time
		err = rows.Scan(&message.Id, &message.Text, &message.Private, &message.Name,
			&message.UserId, &message.Title, &createdAt, &updatedAt)
		if err != nil {
			logger.Errorln(err)
			r.cancel(err)
			return
		}

		message.CreatedAt = createdAt.Format(time.RFC3339)
		message.UpdatedAt = updatedAt.Format(time.RFC3339)

		select {
		case <-r.ctx.Done():
			return
		default:
		}

		r.ch <- &api.MessagesSnapshot{
			Item: &api.MessagesSnapshot_Message{
				Message: message,
			},
		}
	}

	if err := rows.Err(); err != nil {
		logger.Errorln(err)
		r.cancel(err)
		return
	}
}

func (r *Dumper) runFiles() {
	query := "SELECT file_id, message_id, user_id FROM %s"

	defer r.wg.Done()
	defer logger.Debugln("files dump finished")

	rows, err := r.pool.Query(r.ctx, r.filesTable(query))
	if err != nil {
		logger.Errorln(err)
		r.cancel(err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		file := &api.FileSnapshotItem{}

		err = rows.Scan(&file.FileId, &file.MessageId, &file.UserId)
		if err != nil {
			logger.Errorln(err)
			r.cancel(err)
			return
		}

		select {
		case <-r.ctx.Done():
			return
		default:
		}

		r.ch <- &api.MessagesSnapshot{
			Item: &api.MessagesSnapshot_File{
				File: file,
			},
		}
	}

	if err := rows.Err(); err != nil {
		logger.Errorln(err)
		r.cancel(err)
		return
	}
}

func (r *Dumper) runTranslations() {
	query := "SELECT message_id, lang, text, title, created_at, updated_at FROM %s"

	defer r.wg.Done()
	defer logger.Debugln("translations dump finished")

	rows, err := r.pool.Query(r.ctx, r.translationsTable(query))
	if err != nil {
		logger.Errorln(err)
		r.cancel(err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		translation := &api.TranslationSnapshotItem{}

		var createdAt, updatedAt *time.Time
		err = rows.Scan(&translation.MessageId, &translation.Lang, &translation.Text,
			&translation.Title, &createdAt, &updatedAt)
		if err != nil {
			logger.Errorln(err)
			r.cancel(err)
			return
		}

		translation.CreatedAt = createdAt.Format(time.RFC3339)
		translation.UpdatedAt = updatedAt.Format(time.RFC3339)

		select {
		case <-r.ctx.Done():
			return
		default:
		}

		r.ch <- &api.MessagesSnapshot{
			Item: &api.MessagesSnapshot_Translation{
				Translation: translation,
			},
		}
	}

	if err := rows.Err(); err != nil {
		logger.Errorln(err)
		r.cancel(err)
		return
	}
}

func (r *Dumper) Close() (err error) {
	logger.Debugln("close dumper")
	r.cancel(nil)
	r.wg.Wait()
	return nil
}

func (r *Dumper) Restore(ctx context.Context, snapshot *api.MessagesSnapshot) (err error) {
	switch v := snapshot.Item.(type) {
	case *api.MessagesSnapshot_Message:

		query := "INSERT INTO %s(id, text, private, name, user_id, title, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)"

		_, err = r.pool.Exec(ctx, r.messagesTable(query), v.Message.Id, v.Message.Text, v.Message.Private,
			v.Message.Name, v.Message.UserId, v.Message.Title, v.Message.CreatedAt, v.Message.UpdatedAt)

		return

	case *api.MessagesSnapshot_File:

		query := "INSERT INTO %s(file_id, message_id, user_id) VALUES ($1,$2,$3)"

		_, err = r.pool.Exec(ctx, r.filesTable(query), v.File.FileId, v.File.MessageId, v.File.UserId)

		return

	case *api.MessagesSnapshot_Translation:

		query := "INSERT INTO %s(message_id, lang, text, title, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6)"

		_, err = r.pool.Exec(ctx, r.translationsTable(query), v.Translation.MessageId, v.Translation.Lang, v.Translation.Text,
			v.Translation.Title, v.Translation.CreatedAt, v.Translation.UpdatedAt)

		return

	default:
		return errors.New("unknown snapshot item")
	}
}

func (r Dumper) messagesTable(query string) string {
	return fmt.Sprintf(query, r.messagesTableName)
}

func (r Dumper) filesTable(query string) string {
	return fmt.Sprintf(query, r.filesTableName)
}

func (r Dumper) translationsTable(query string) string {
	return fmt.Sprintf(query, r.translationsTableName)
}