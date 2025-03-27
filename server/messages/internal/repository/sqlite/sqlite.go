package repository

import (
	"context"
	"errors"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/bd878/gallery/server/logger"
	"github.com/bd878/gallery/server/utils"
	"github.com/bd878/gallery/server/messages/pkg/model"
)

type Repository struct {
	pool         *sql.DB
	insertStmt   *sql.Stmt
	updateStmt   *sql.Stmt
	deleteStmt   *sql.Stmt
	ascStmt *sql.Stmt
	descStmt *sql.Stmt
	ascThreadStmt *sql.Stmt
	descThreadStmt *sql.Stmt
	selectStmt *sql.Stmt
	deleteThreadStmt *sql.Stmt
}

func New(dbFilePath string) *Repository {
	pool, err := sql.Open("sqlite3", "file:" + dbFilePath)
	if err != nil {
		panic(err)
	}

	selectStmt := utils.Must(pool.Prepare(`
SELECT id, user_id, thread_id, create_utc_nano, update_utc_nano, text, file_id, private
FROM messages
WHERE id = :id AND (thread_id ISNULL OR thread_id = 0)
;`,
	))

	insertStmt := utils.Must(pool.Prepare(`
INSERT INTO messages( 
	id,
	thread_id,
	user_id,
	create_utc_nano,
	update_utc_nano,
	text,
	file_id,
	private
) VALUES (:id, :threadId, :userId, :createUtcNano, :updateUtcNano, :text, :fileId, :private)
;`,
	))

	updateStmt := utils.Must(pool.Prepare(`
UPDATE messages SET 
	text = :text,
	update_utc_nano = :updateUtcNano,
	private = :private
WHERE id = :id AND user_id = :userId
;`,
	))

	deleteStmt := utils.Must(pool.Prepare(`
DELETE FROM messages 
WHERE id = :id AND user_id = :userId
;`,
	))

	deleteThreadStmt := utils.Must(pool.Prepare(`
DELETE FROM messages
WHERE thread_id = :threadId
;`,
	))

	ascStmt := utils.Must(pool.Prepare(`
SELECT id, user_id, thread_id, create_utc_nano, update_utc_nano, text, file_id, private
FROM messages
WHERE user_id = :userId AND (thread_id ISNULL OR thread_id = 0)
ORDER BY create_utc_nano ASC
LIMIT :limit OFFSET :offset
;`,
	))

	descStmt := utils.Must(pool.Prepare(`
SELECT id, user_id, thread_id, create_utc_nano, update_utc_nano, text, file_id, private
FROM messages
WHERE user_id = :userId AND (thread_id ISNULL OR thread_id = 0)
ORDER BY create_utc_nano DESC
LIMIT :limit OFFSET :offset
;`,
	))

	ascThreadStmt := utils.Must(pool.Prepare(`
SELECT id, user_id, thread_id, create_utc_nano, update_utc_nano, text, file_id, private
FROM messages
WHERE user_id = :userId AND thread_id = :threadId
ORDER BY create_utc_nano ASC
LIMIT :limit OFFSET :offset
;`,
	))

	descThreadStmt := utils.Must(pool.Prepare(`
SELECT id, user_id, thread_id, create_utc_nano, update_utc_nano, text, file_id, private
FROM messages
WHERE user_id = :userId AND thread_id = :threadId
ORDER BY create_utc_nano DESC
LIMIT :limit OFFSET :offset
;`,
	))

	return &Repository{
		pool: pool,
		insertStmt: insertStmt,
		updateStmt: updateStmt,
		deleteStmt: deleteStmt,
		ascStmt: ascStmt,
		descStmt: descStmt,
		selectStmt: selectStmt,
		ascThreadStmt: ascThreadStmt,
		descThreadStmt: descThreadStmt,
		deleteThreadStmt: deleteThreadStmt,
	}
}

/**
 * Receives message id from params;
 * Does not put message with same id
 * twice
 */
func (r *Repository) Create(ctx context.Context, log *logger.Logger, message *model.Message) error {
	var threadIdCol sql.NullInt32
	if message.ThreadID != 0 {
		threadIdCol.Int32 = int32(message.ThreadID)
		threadIdCol.Valid = true
	}

	var fileIdCol sql.NullInt32
	if message.FileID != 0 {
		fileIdCol.Int32 = int32(message.FileID)
		fileIdCol.Valid = true
	}

	var privateCol int32
	if message.Private {
		privateCol = 1
	}

	_, err := r.insertStmt.ExecContext(ctx,
		sql.Named("id", message.ID),
		sql.Named("threadId", threadIdCol),
		sql.Named("userId", message.UserID),
		sql.Named("createUtcNano", message.CreateUTCNano),
		sql.Named("updateUtcNano", message.UpdateUTCNano),
		sql.Named("text", message.Text),
		sql.Named("fileId", fileIdCol),
		sql.Named("private", privateCol),
	)
	if err != nil {
		log.Errorw("failed to insert new message ", "error", err)
		return errors.New("failed to put message")
	}
	return nil
}

// TODO: utilise ctx timeout
func (r *Repository) Delete(ctx context.Context, log *logger.Logger, params *model.DeleteMessageParams) error {
	var (
		tx *sql.Tx
		err error
	)

	tx, err = r.pool.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		log.Errorln("failed to start transaction")
		return err
	}

	defer func() {
		p := recover()
		switch {
		case p != nil:
			_ = tx.Rollback()
			panic(p)
		case err != nil:
			rErr := tx.Rollback()
			if rErr != nil {
				log.Errorw("failed to rollback delete stmt", "error", rErr)
			}
		default:
			err = tx.Commit()
		}
	}()

	msg, err := r.Read(ctx, log, params.ID)
	if err != nil {
		log.Errorln("cannot delete, no message found")
		return err
	}

	// parent message, no thread id
	if msg.ThreadID == 0 {
		err = r.deleteThreadMessages(ctx, log, tx, params.ID)
		if err != nil {
			log.Errorln("failed to delete thread messages")
			return err
		}
	}

	_, err = tx.StmtContext(ctx, r.deleteStmt).ExecContext(ctx, sql.Named("id", params.ID), sql.Named("userId", params.UserID))
	if err != nil {
		log.Errorln("failed to delete message")
		return err
	}

	return nil
}

func (r *Repository) deleteThreadMessages(ctx context.Context, log *logger.Logger, tx *sql.Tx, threadID int32) error {
	_, err := tx.StmtContext(ctx, r.deleteThreadStmt).ExecContext(ctx, sql.Named("threadId", threadID))
	if err != nil {
		log.Errorw("failed to delete thread messages", "threadId", threadID)
		return err
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, log *logger.Logger, params *model.UpdateMessageParams) error {
	var privateCol int32
	if params.Private {
		privateCol = 1
	}
	_, err := r.updateStmt.ExecContext(ctx,
		sql.Named("text", params.Text),
		sql.Named("updateUtcNano", params.UpdateUTCNano),
		sql.Named("id", params.ID),
		sql.Named("userId", params.UserID),
		sql.Named("private", privateCol))
	if err != nil {
		log.Errorln("failed to update message")
		return err
	}
	return nil
}

func (r *Repository) Truncate(ctx context.Context, log *logger.Logger) error {
	_, err := r.pool.ExecContext(ctx, "DELETE FROM messages")
	if err != nil {
		log.Error("failed to delete messages")
		return err
	}
	return nil
}

func (r *Repository) ReadThreadMessages(ctx context.Context, log *logger.Logger, params *model.ReadThreadMessagesParams) (
	*model.ReadThreadMessagesResult, error,
) {
	var (
		rows *sql.Rows
		err error
	)

	if params.Ascending {
		rows, err = r.ascThreadStmt.QueryContext(ctx,
			sql.Named("userId", params.UserID), sql.Named("threadId", params.ThreadID),
			sql.Named("limit", params.Limit), sql.Named("offset", params.Offset))
	} else {
		rows, err = r.descThreadStmt.QueryContext(ctx,
			sql.Named("userId", params.UserID), sql.Named("threadId", params.ThreadID),
			sql.Named("limit", params.Limit), sql.Named("offset", params.Offset))
	}

	if err != nil {
		log.Errorln("failed to query messages context")
		return nil, err
	}

	defer rows.Close()

	log.Infow("repository read thread messages", "user_id", params.UserID, "thread_id", params.ThreadID)

	selected, err := r.selectMessages(ctx, log, rows, params.UserID, params.Limit, params.Offset)
	if err != nil {
		log.Errorln("cannot select messages")
		return nil, err
	}

	return &model.ReadThreadMessagesResult{
		Messages: selected.List,
		IsLastPage: selected.IsLastPage,
	}, nil
}

type selectedMessages struct {
	List []*model.Message
	IsLastPage bool
}

func (r *Repository) ReadAllMessages(ctx context.Context, log *logger.Logger, params *model.ReadAllMessagesParams) (
	*model.ReadAllMessagesResult, error,
) {
	var (
		rows *sql.Rows
		err error
	)

	if params.Ascending {
		rows, err = r.ascStmt.QueryContext(ctx,
			sql.Named("userId", params.UserID), sql.Named("limit", params.Limit), sql.Named("offset", params.Offset),
		)
	} else {
		rows, err = r.descStmt.QueryContext(ctx,
			sql.Named("userId", params.UserID), sql.Named("limit", params.Limit), sql.Named("offset", params.Offset),
		)
	}

	if err != nil {
		log.Error("failed to query messages context")
		return nil, err
	}

	defer rows.Close()

	selected, err := r.selectMessages(ctx, log, rows, params.UserID, params.Limit, params.Offset)
	if err != nil {
		log.Errorln("cannot select messages")
		return nil, err
	}

	return &model.ReadAllMessagesResult{
		Messages: selected.List,
		IsLastPage: selected.IsLastPage,
	}, nil
}

func (r *Repository) Read(ctx context.Context, log *logger.Logger, id int32) (
	*model.Message, error,
) {
	var (
		_id int32
		userId int32
		threadIdCol sql.NullInt32
		createUtcNano int64
		updateUtcNano int64
		text string
		fileIdCol sql.NullInt32
		privateCol int32
	)

	err := r.selectStmt.QueryRowContext(ctx, sql.Named("id", id)).Scan(
		&_id, &userId, &threadIdCol, &createUtcNano, &updateUtcNano, &text, &fileIdCol, &privateCol)
	if err != nil {
		log.Errorln("failed to select one message")
		return nil, err
	}

	msg := &model.Message{
		ID: id,
		UserID: userId,
		CreateUTCNano: createUtcNano,
		UpdateUTCNano: updateUtcNano,
		Text: text,
		Private: true,
	}

	if threadIdCol.Valid {
		msg.ThreadID = threadIdCol.Int32
	}

	if fileIdCol.Valid {
		msg.FileID = fileIdCol.Int32
	}

	if privateCol == 0 {
		msg.Private = false
	}

	return msg, nil
}

func (r *Repository) selectMessages(ctx context.Context, log *logger.Logger, rows *sql.Rows, userID, limit, offset int32) (
	*selectedMessages, error,
) {
	var isLastPage bool
	var err error

	var res []*model.Message
	for rows.Next() {
		var (
			id int32
			userId int32
			threadIdCol sql.NullInt32
			createUtcNano int64
			updateUtcNano int64
			text string
			fileIdCol sql.NullInt32
			privateCol int32
		)
		if err := rows.Scan(
			&id,
			&userId,
			&threadIdCol,
			&createUtcNano,
			&updateUtcNano,
			&text,
			&fileIdCol,
			&privateCol,
		); err != nil {
			log.Error("failed to scan row")
			return nil, err
		}

		msg := &model.Message{
			ID: id,
			UserID: userId,
			CreateUTCNano: createUtcNano,
			UpdateUTCNano: updateUtcNano,
			Text: text,
			Private: true,
		}

		if threadIdCol.Valid {
			msg.ThreadID = threadIdCol.Int32
		}

		if fileIdCol.Valid {
			msg.FileID = fileIdCol.Int32
		}

		if privateCol == 0 {
			msg.Private = false
		}

		res = append(res, msg)
	}

	if int32(len(res)) < limit {
		isLastPage = true
	} else {
		row := r.pool.QueryRowContext(ctx,
			"SELECT COUNT(*) FROM messages WHERE user_id = :userId",
			sql.Named("userId", userID),
		)
		if err != nil {
			isLastPage = false
		} else {
			var countMessages int32
			if err := row.Scan(&countMessages); err != nil {
				isLastPage = false
			}

			if countMessages <= offset + limit {
				isLastPage = true
			}
		}
	}

	return &selectedMessages{
		List: res,
		IsLastPage: isLastPage,
	}, nil
}