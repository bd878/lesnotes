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
	threadsTableName       string
	ctx                    context.Context
	cancel                 context.CancelCauseFunc
	ch                     chan *api.SearchSnapshot
	wg                     sync.WaitGroup
}

func NewDumper(pool *pgxpool.Pool, messagesTableName, filesTableName, threadsTableName, translationsTableName string) *Dumper {
	return &Dumper{
		pool:                    pool,
		messagesTableName:       messagesTableName,
		filesTableName:          filesTableName,
		threadsTableName:        threadsTableName,
		translationsTableName:   translationsTableName,
	}
}

func (r *Dumper) Open(ctx context.Context) (ch chan *api.SearchSnapshot, err error) {
	r.ctx, r.cancel = context.WithCancelCause(ctx)
	ch = make(chan *api.SearchSnapshot, 100)
	r.ch = ch

	r.wg.Add(1)
	go r.runMessages()
	r.wg.Add(1)
	go r.runFiles()
	r.wg.Add(1)
	go r.runTranslations()
	r.wg.Add(1)
	go r.runThreads()

	go func() {
		r.wg.Wait()
		close(r.ch)
	}()

	return
}

func (r *Dumper) runMessages() {
	query := "SELECT id, user_id, name, text, title, private, created_at, updated_at FROM %s"

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
		message := &api.SearchMessageSnapshotItem{}

		var createdAt, updatedAt *time.Time
		err = rows.Scan(&message.Id, &message.UserId, &message.Name, &message.Text, &message.Title,
			&message.Private, &createdAt, &updatedAt)
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

		r.ch <- &api.SearchSnapshot{
			Item: &api.SearchSnapshot_Message{
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
	query := "SELECT id, owner_id, name, mime, created_at, updated_at, size, private, description FROM %s"

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
		file := &api.SearchFileSnapshotItem{}

		var createdAt, updatedAt *time.Time
		err = rows.Scan(&file.Id, &file.UserId, &file.Name, &file.Mime, &createdAt, &updatedAt,
			&file.Size, &file.Private, &file.Description)
		if err != nil {
			logger.Errorln(err)
			r.cancel(err)
			return
		}

		file.CreatedAt = createdAt.Format(time.RFC3339)
		file.UpdatedAt = updatedAt.Format(time.RFC3339)

		select {
		case <-r.ctx.Done():
			return
		default:
		}

		r.ch <- &api.SearchSnapshot{
			Item: &api.SearchSnapshot_File{
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
	query := "SELECT message_id, user_id, lang, text, title, created_at, updated_at FROM %s"

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
		translation := &api.SearchTranslationSnapshotItem{}

		var createdAt, updatedAt *time.Time
		err = rows.Scan(&translation.MessageId, &translation.UserId, &translation.Lang, &translation.Text,
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

		r.ch <- &api.SearchSnapshot{
			Item: &api.SearchSnapshot_Translation{
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

func (r *Dumper) runThreads() {
	query := "SELECT id, user_id, parent_id, name, description, private, created_at, updated_at FROM %s"

	defer r.wg.Done()
	defer logger.Debugln("threads dump finished")

	rows, err := r.pool.Query(r.ctx, r.threadsTable(query))
	if err != nil {
		logger.Errorln(err)
		r.cancel(err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		thread := &api.SearchThreadSnapshotItem{}

		var createdAt, updatedAt *time.Time
		err = rows.Scan(&thread.Id, &thread.UserId, &thread.ParentId, &thread.Name,
			&thread.Description, &thread.Private, &createdAt, &updatedAt)
		if err != nil {
			logger.Errorln(err)
			r.cancel(err)
			return
		}

		thread.CreatedAt = createdAt.Format(time.RFC3339)
		thread.UpdatedAt = updatedAt.Format(time.RFC3339)

		select {
		case <-r.ctx.Done():
			return
		default:
		}

		r.ch <- &api.SearchSnapshot{
			Item: &api.SearchSnapshot_Thread{
				Thread: thread,
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

func (r *Dumper) Restore(ctx context.Context, snapshot *api.SearchSnapshot) (err error) {
	switch v := snapshot.Item.(type) {
	case *api.SearchSnapshot_Message:

		query := "INSERT INTO %s(id, user_id, name, text, title, private, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)"

		_, err = r.pool.Exec(ctx, r.messagesTable(query), v.Message.Id, v.Message.UserId, v.Message.Name,
			v.Message.Text, v.Message.Title, v.Message.Private, v.Message.CreatedAt, v.Message.UpdatedAt)

		return

	case *api.SearchSnapshot_File:

		query := "INSERT INTO %s(id, owner_id, name, mime, created_at, updated_at, size, private, description) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)"

		_, err = r.pool.Exec(ctx, r.filesTable(query), v.File.Id, v.File.UserId, v.File.Name, v.File.Mime, v.File.CreatedAt, v.File.UpdatedAt,
			v.File.Size, v.File.Private, v.File.Description)

		return

	case *api.SearchSnapshot_Translation:

		query := "INSERT INTO %s(message_id, user_id, lang, text, title, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7)"

		_, err = r.pool.Exec(ctx, r.translationsTable(query), v.Translation.MessageId, v.Translation.UserId, v.Translation.Lang,
			v.Translation.Text, v.Translation.Title, v.Translation.CreatedAt, v.Translation.UpdatedAt)

		return

	case *api.SearchSnapshot_Thread:

		query := "INSERT INTO %s(id, user_id, parent_id, name, description, private, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)"

		_, err = r.pool.Exec(ctx, r.threadsTable(query), v.Thread.Id, v.Thread.UserId, v.Thread.ParentId,
			v.Thread.Name, v.Thread.Description, v.Thread.Private, v.Thread.CreatedAt, v.Thread.UpdatedAt)

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

func (r Dumper) threadsTable(query string) string {
	return fmt.Sprintf(query, r.threadsTableName)
}