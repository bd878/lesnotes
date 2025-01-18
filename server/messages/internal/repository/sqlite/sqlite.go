package repository

import (
  "context"
  "errors"
  "database/sql"

  _ "github.com/mattn/go-sqlite3"
  "github.com/bd878/gallery/server/logger"
  "github.com/bd878/gallery/server/messages/pkg/model"
)

type Repository struct {
  db           *sql.DB
  insertStmt   *sql.Stmt
  updateStmt   *sql.Stmt
  deleteStmt   *sql.Stmt
}

func New(dbFilePath string) *Repository {
  db, err := sql.Open("sqlite3", "file:" + dbFilePath)
  if err != nil {
    panic(err)
  }

  insertStmt, err := db.Prepare(
    "INSERT INTO messages(" +
      "id, " +
      "user_id, " +
      "create_utc_nano, " +
      "update_utc_nano, " +
      "text, " +
      "file_id " +
    ") VALUES (?,?,?,?,?,?)",
  )
  if err != nil {
    panic(err)
  }

  updateStmt, err := db.Prepare(
    "UPDATE messages SET " +
      "text = ?, " +
      "update_utc_nano = ? " +
      "WHERE id = ? AND user_id = ?",
  )
  if err != nil {
    panic(err)
  }

  deleteStmt, err := db.Prepare(
    "DELETE FROM messages " +
    "WHERE id = ? AND user_id = ?",
  )
  if err != nil {
    panic(err)
  }

  return &Repository{
    db: db,
    insertStmt: insertStmt,
    updateStmt: updateStmt,
    deleteStmt: deleteStmt,
  }
}

/**
 * Receives message id from params;
 * Does not put message with same id
 * twice
 */
func (r *Repository) Create(ctx context.Context, log *logger.Logger, params *model.SaveMessageParams) error {
  _, err := r.insertStmt.ExecContext(ctx,
    params.Message.ID,
    params.Message.UserID,
    params.Message.CreateUTCNano,
    params.Message.UpdateUTCNano,
    params.Message.Text,
    params.Message.FileID,
  )
  if err != nil {
    log.Error("failed to insert new message", err)
    return errors.New("failed to put message")
  }
  return nil
}

func (r *Repository) Delete(ctx context.Context, log *logger.Logger, params *model.DeleteMessageParams) error {
  _, err := r.deleteStmt.ExecContext(ctx, params.ID, params.UserID)
  if err != nil {
    log.Error("failed to delete message", err)
    return errors.New("failed to delete message")
  }
  return nil
}

func (r *Repository) Update(ctx context.Context, log *logger.Logger, params *model.UpdateMessageParams) error {
  _, err := r.updateStmt.ExecContext(ctx,
    params.Text,
    params.UpdateUTCNano,
    params.ID,
    params.UserID,
  )
  if err != nil {
    log.Error("failed to update message", err)
    return errors.New("failed to update message")
  }
  return nil
}

func (r *Repository) Truncate(ctx context.Context, log *logger.Logger) error {
  _, err := r.db.ExecContext(ctx, "DELETE FROM messages")
  if err != nil {
    log.Error("failed to delete messages")
    return err
  }
  return nil
}

const ascendingStmt = `
SELECT id, user_id, create_utc_nano, update_utc_nano, text, file_id
FROM messages
WHERE user_id = ?
ORDER BY update_utc_nano ASC
LIMIT ? OFFSET ?
`

const descendingStmt = `
SELECT id, user_id, create_utc_nano, update_utc_nano, text, file_id
FROM messages
WHERE user_id = ?
ORDER BY update_utc_nano DESC
LIMIT ? OFFSET ?
`

func (r *Repository) ReadUserMessages(ctx context.Context, log *logger.Logger, params *model.ReadUserMessagesParams) (
  *model.ReadUserMessagesResult, error,
) {
  var rows *sql.Rows
  var err error
  var isLastPage bool

  if params.Ascending {
    rows, err = r.db.QueryContext(ctx, ascendingStmt, params.UserID, params.Limit, params.Offset)
  } else {
    rows, err = r.db.QueryContext(ctx, descendingStmt, params.UserID, params.Limit, params.Offset)
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
    row := r.db.QueryRowContext(ctx,
      "SELECT COUNT(*) FROM messages WHERE user_id = ?",
      params.UserID,
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
  rows, err := r.db.QueryContext(ctx,
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