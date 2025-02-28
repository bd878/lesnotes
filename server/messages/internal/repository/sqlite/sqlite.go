package repository

import (
  "context"
  "errors"
  "database/sql"

  _ "github.com/mattn/go-sqlite3"
  "github.com/bd878/gallery/server/logger"
  "github.com/bd878/gallery/server/utils"
  "github.com/bd878/gallery/server/messages/pkg/model"
  filesmodel "github.com/bd878/gallery/server/files/pkg/model"
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
}

func New(dbFilePath string) *Repository {
  pool, err := sql.Open("sqlite3", "file:" + dbFilePath)
  if err != nil {
    panic(err)
  }

  insertStmt := utils.Must(pool.Prepare(`
INSERT INTO messages( 
  id,
  thread_id,
  user_id,
  create_utc_nano,
  update_utc_nano,
  text,
  file_id
) VALUES (:id, :threadId, :userId, :createUtcNano, :updateUtcNano, :text, :fileId)
;`,
  ))

  updateStmt := utils.Must(pool.Prepare(`
UPDATE messages SET 
  text = :text,
  update_utc_nano = :updateUtcNano ,
  thread_id = :threadId
WHERE id = :id AND user_id = :userId
;`,
  ))

  deleteStmt := utils.Must(pool.Prepare(`
DELETE FROM messages 
WHERE id = :id AND user_id = :userId
;`,
  ))

  ascStmt := utils.Must(pool.Prepare(`
SELECT id, user_id, create_utc_nano, update_utc_nano, text, file_id
FROM messages
WHERE user_id = :userId
ORDER BY create_utc_nano ASC
LIMIT :limit OFFSET :offset
;`,
  ))

  descStmt := utils.Must(pool.Prepare(`
SELECT id, user_id, create_utc_nano, update_utc_nano, text, file_id
FROM messages
WHERE user_id = :userId
ORDER BY create_utc_nano DESC
LIMIT :limit OFFSET :offset
;`,
  ))

  ascThreadStmt := utils.Must(pool.Prepare(`
SELECT id, user_id, create_utc_nano, update_utc_nano, text, file_id
FROM messages
WHERE user_id = :userId AND thread_id = :threadId
ORDER BY create_utc_nano ASC
LIMIT :limit OFFSET :offset
;`,
  ))

  descThreadStmt := utils.Must(pool.Prepare(`
SELECT id, user_id, create_utc_nano, update_utc_nano, text, file_id
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
    ascThreadStmt: ascThreadStmt,
    descThreadStmt: descThreadStmt,
  }
}

/**
 * Receives message id from params;
 * Does not put message with same id
 * twice
 */
func (r *Repository) Create(ctx context.Context, log *logger.Logger, message *model.Message) error {
  _, err := r.insertStmt.ExecContext(ctx,
    sql.Named("id", message.ID),
    sql.Named("threadId", message.ThreadID),
    sql.Named("userId", message.UserID),
    sql.Named("createUtcNano", message.CreateUTCNano),
    sql.Named("updateUtcNano", message.UpdateUTCNano),
    sql.Named("text", message.Text),
    sql.Named("fileId", message.File.ID),
  )
  if err != nil {
    log.Errorw("failed to insert new message ", "error", err)
    return errors.New("failed to put message")
  }
  return nil
}

func (r *Repository) Delete(ctx context.Context, log *logger.Logger, params *model.DeleteMessageParams) error {
  _, err := r.deleteStmt.ExecContext(ctx, sql.Named("id", params.ID), sql.Named("userId", params.UserID))
  if err != nil {
    log.Errorw("failed to delete message", "error", err)
    return errors.New("failed to delete message")
  }
  return nil
}

func (r *Repository) Update(ctx context.Context, log *logger.Logger, params *model.UpdateMessageParams) error {
  _, err := r.updateStmt.ExecContext(ctx,
    sql.Named("text", params.Text),
    sql.Named("threadId", params.ThreadID),
    sql.Named("updateUtcNano", params.UpdateUTCNano),
    sql.Named("id", params.ID),
    sql.Named("userId", params.UserID),
  )
  if err != nil {
    log.Errorw("failed to update message", "error", err)
    return errors.New("failed to update message")
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
  var rows *sql.Rows
  var err error

  if params.Ascending {
    rows, err = r.ascThreadStmt.QueryContext(ctx,
      sql.Named("userId", params.UserID), sql.Named("threadId", params.ThreadID),
      sql.Named("limit", params.Limit), sql.Named("offset", params.Offset),
    )
  } else {
    rows, err = r.descStmt.QueryContext(ctx,
      sql.Named("userId", params.UserID), sql.Named("threadId", params.ThreadID),
      sql.Named("limit", params.Limit), sql.Named("offset", params.Offset),
    )
  }

  if err != nil {
    log.Errorln("failed to query messages context")
    return nil, err
  }

  defer rows.Close()

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
  var rows *sql.Rows
  var err error

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

func (r *Repository) selectMessages(ctx context.Context, log *logger.Logger, rows *sql.Rows, userID, limit, offset int32) (
  *selectedMessages, error,
) {
  var isLastPage bool
  var err error

  var res []*model.Message
  for rows.Next() {
    var id int32
    var userId int32
    var createUtcNano int64
    var updateUtcNano int64
    var text string
    var fileIdCol sql.NullInt32
    if err := rows.Scan(
      &id,
      &userId,
      &createUtcNano,
      &updateUtcNano,
      &text,
      &fileIdCol,
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
      File: &filesmodel.File{},
    }

    if fileIdCol.Valid {
      msg.File.ID = fileIdCol.Int32
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

func (r *Repository) GetBatch(ctx context.Context, log *logger.Logger) ([]*model.Message, error) {
  rows, err := r.pool.QueryContext(ctx,
    "SELECT id, user_id, create_utc_nano, update_utc_nano, text, file_id " +
    "FROM messages",
  )
  if err != nil {
    log.Error("failed to query get batch context")
    return nil, err
  }
  defer rows.Close()

  var res []*model.Message
  for rows.Next() {
    var id int32
    var userId int32
    var createUtcNano int64
    var updateUtcNano int64
    var text string
    var fileIdCol sql.NullInt32
    if err := rows.Scan(
      &id,
      &userId,
      &createUtcNano,
      &updateUtcNano,
      &text,
      &fileIdCol,
    ); err != nil {
      log.Error("failed to scan row")
      return nil, err
    }
    var fileId int32
    if fileIdCol.Valid {
      fileId = fileIdCol.Int32
    }
    res = append(res, &model.Message{
      ID: id,
      UserID: userId,
      CreateUTCNano: createUtcNano,
      UpdateUTCNano: updateUtcNano,
      Text: text,
      File: &filesmodel.File{
        ID: fileId,
      },
    })
  }

  return res, nil
}