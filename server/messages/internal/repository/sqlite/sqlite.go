package repository

import (
  "context"
  "database/sql"

  _ "github.com/mattn/go-sqlite3"
  "github.com/bd878/gallery/server/logger"
  "github.com/bd878/gallery/server/messages/pkg/model"
)

const NullMessageID int32 = 0

type Repository struct {
  db           *sql.DB
  insertStmt   *sql.Stmt
}

func New(dbFilePath string) (*Repository, error) {
  db, err := sql.Open("sqlite3", "file:" + dbFilePath)
  if err != nil {
    logger.Error("failed to establish connection with sql by given path", dbFilePath)
    return nil, err
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
    logger.Error("failed to prepare insert statement")
    return nil, err
  }

  return &Repository{
    db: db,
    insertStmt: insertStmt,
  }, nil
}

func (r *Repository) Put(ctx context.Context, log *logger.Logger, params *model.PutParams) (
  int32, error,
) {
  _, err := r.insertStmt.ExecContext(ctx,
    params.Message.ID,
    params.Message.UserID,
    params.Message.CreateUTCNano,
    params.Message.UpdateUTCNano,
    params.Message.Text,
    params.Message.FileID,
  )
  if err != nil {
    log.Error("failed to insert new message")
    return NullMessageID, err
  }
  return params.Message.ID, nil
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
ORDER BY id ASC
LIMIT ? OFFSET ?
`

const descendingStmt = `
SELECT id, user_id, create_utc_nano, update_utc_nano, text, file_id
FROM messages
WHERE user_id = ?
ORDER BY id DESC
LIMIT ? OFFSET ?
`

func (r *Repository) Get(ctx context.Context, log *logger.Logger, params *model.GetParams) (
  *model.MessagesList, error,
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

  return &model.MessagesList{
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