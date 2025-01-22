package repository

import (
  "sync/atomic"
  "context"
  "errors"
  "database/sql"

  _ "github.com/mattn/go-sqlite3"
  "github.com/bd878/gallery/server/logger"
  "github.com/bd878/gallery/server/messages/pkg/model"
)

type Repository struct {
  dbFilePath    string
  pool          atomic.Pointer[sql.DB]
  insertStmt   *sql.Stmt
  updateStmt   *sql.Stmt
  deleteStmt   *sql.Stmt
}

func New(dbFilePath string) *Repository {
  pool, err := sql.Open("sqlite3", "file:" + dbFilePath)
  if err != nil {
    panic(err)
  }

  insertStmt, err := pool.Prepare(`
INSERT INTO messages( 
  id,
  user_id,
  create_utc_nano,
  update_utc_nano,
  text,
  file_id
) VALUES (:id, :userId, :createUtcNano, :updateUtcNano, :text, :fileId)
;`,
  )
  if err != nil {
    panic(err)
  }

  updateStmt, err := pool.Prepare(`
UPDATE messages SET 
  text = :text,
  update_utc_nano = :updateUtcNano 
WHERE id = :id AND user_id = :userId
;`,
  )
  if err != nil {
    panic(err)
  }

  deleteStmt, err := pool.Prepare(`
DELETE FROM messages 
WHERE id = :id AND user_id = :userId
;`,
  )
  if err != nil {
    panic(err)
  }

  repo := &Repository{
    dbFilePath: dbFilePath,
    insertStmt: insertStmt,
    updateStmt: updateStmt,
    deleteStmt: deleteStmt,
  }

  repo.pool.Store(pool)
  return repo
}

func (r *Repository) Restore(nextDBFilePath string) error {
  // TODO: mv snapshot replica to main destination
  nextPool, err := sql.Open("sqlite3", "file:" + nextDBFilePath)
  if err != nil {
    logger.Error("failed to open next pool")
    return err
  }

  r.dbFilePath = nextDBFilePath

  go func(pool *sql.DB) {
    if err := pool.Close(); err != nil {
      logger.Error("failed to close current pool", "error=", err)
    }
  }(r.pool.Swap(nextPool))

  return nil
}

/**
 * Receives message id from params;
 * Does not put message with same id
 * twice
 */
func (r *Repository) Create(ctx context.Context, log *logger.Logger, params *model.SaveMessageParams) error {
  _, err := r.insertStmt.ExecContext(ctx,
    sql.Named("id", params.Message.ID),
    sql.Named("userId", params.Message.UserID),
    sql.Named("createUtcNano", params.Message.CreateUTCNano),
    sql.Named("updateUtcNano", params.Message.UpdateUTCNano),
    sql.Named("text", params.Message.Text),
    sql.Named("fileId", params.Message.FileID),
  )
  if err != nil {
    log.Error("failed to insert new message", err)
    return errors.New("failed to put message")
  }
  return nil
}

func (r *Repository) Delete(ctx context.Context, log *logger.Logger, params *model.DeleteMessageParams) error {
  _, err := r.deleteStmt.ExecContext(ctx, sql.Named("id", params.ID), sql.Named("userId", params.UserID))
  if err != nil {
    log.Error("failed to delete message", err)
    return errors.New("failed to delete message")
  }
  return nil
}

func (r *Repository) Update(ctx context.Context, log *logger.Logger, params *model.UpdateMessageParams) error {
  _, err := r.updateStmt.ExecContext(ctx,
    sql.Named("text", params.Text),
    sql.Named("updateUtcNano", params.UpdateUTCNano),
    sql.Named("id", params.ID),
    sql.Named("userId", params.UserID),
  )
  if err != nil {
    log.Error("failed to update message", err)
    return errors.New("failed to update message")
  }
  return nil
}

func (r *Repository) Truncate(ctx context.Context, log *logger.Logger) error {
  _, err := r.pool.Load().ExecContext(ctx, "DELETE FROM messages")
  if err != nil {
    log.Error("failed to delete messages")
    return err
  }
  return nil
}

const ascendingStmt = `
SELECT id, user_id, create_utc_nano, update_utc_nano, text, file_id
FROM messages
WHERE user_id = :userId
ORDER BY update_utc_nano ASC
LIMIT :limit OFFSET :offset
`

const descendingStmt = `
SELECT id, user_id, create_utc_nano, update_utc_nano, text, file_id
FROM messages
WHERE user_id = :userId
ORDER BY update_utc_nano DESC
LIMIT :limit OFFSET :offset
`

func (r *Repository) ReadUserMessages(ctx context.Context, log *logger.Logger, params *model.ReadUserMessagesParams) (
  *model.ReadUserMessagesResult, error,
) {
  var rows *sql.Rows
  var err error
  var isLastPage bool

  if params.Ascending {
    rows, err = r.pool.Load().QueryContext(ctx, ascendingStmt,
      sql.Named("userId", params.UserID), sql.Named("limit", params.Limit), sql.Named("offset", params.Offset),
    )
  } else {
    rows, err = r.pool.Load().QueryContext(ctx, descendingStmt,
      sql.Named("userId", params.UserID), sql.Named("limit", params.Limit), sql.Named("offset", params.Offset),
    )
  }

  if err != nil {
    log.Error("failed to query messages context")
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
      FileID: fileId,
    })
  }

  if int32(len(res)) < params.Limit {
    isLastPage = true
  } else {
    row := r.pool.Load().QueryRowContext(ctx,
      "SELECT COUNT(*) FROM messages WHERE user_id = :userId",
      sql.Named("userId", params.UserID),
    )
    if err != nil {
      isLastPage = false
    } else {
      var countMessages int32
      if err := row.Scan(&countMessages); err != nil {
        isLastPage = false
      }

      if countMessages <= params.Offset + params.Limit {
        isLastPage = true
      }
    }
  }

  return &model.ReadUserMessagesResult{
    Messages: res,
    IsLastPage: isLastPage,
  }, nil
}

func (r *Repository) GetBatch(ctx context.Context, log *logger.Logger) ([]*model.Message, error) {
  rows, err := r.pool.Load().QueryContext(ctx,
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
      FileID: fileId,
    })
  }

  return res, nil
}