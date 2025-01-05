package repository

import (
  "context"
  "errors"
  "database/sql"

  _ "github.com/mattn/go-sqlite3"
  "github.com/bd878/gallery/server/logger"
  "github.com/bd878/gallery/server/messages/pkg/model"
  "github.com/bd878/gallery/server/messages/internal/repository"
)

type Repository struct {
  db *sql.DB

  insertSt *sql.Stmt
}

func New(dbfilepath string) (*Repository, error) {
  db, err := sql.Open("sqlite3", "file:" + dbfilepath)
  if err != nil {
    return nil, err
  }

  insertSt, err := db.Prepare(
    "INSERT INTO messages(" +
      "user_id, " +
      "createtime, " +
      "message, " +
      "file, " +
      "file_id, " +
      "log_index, " +
      "log_term" +
    ") VALUES (?,?,?,?,?,?,?)",
  )
  if err != nil {
    return nil, err
  }

  return &Repository{
    db: db,

    insertSt: insertSt,
  }, nil
}

func (r *Repository) Put(ctx context.Context, _ *logger.Logger, params *model.PutParams) (
  model.MessageId, error,
) {
  res, err := r.insertSt.ExecContext(ctx,
    params.Message.UserId,
    params.Message.CreateTime,
    params.Message.Value,
    params.Message.FileName,
    params.Message.FileId,
    params.Message.LogIndex,
    params.Message.LogTerm,
  )
  if err != nil {
    return model.NullMsgId, err
  }
  id, _ := res.LastInsertId()
  return model.MessageId(id), nil
}

func (r *Repository) Truncate(ctx context.Context, _ *logger.Logger) error {
  _, err := r.db.ExecContext(ctx,
    "DELETE FROM messages",
  )
  if err != nil {
    return err
  }
  return nil
}

func (r *Repository) FindByIndexTerm(ctx context.Context, _ *logger.Logger, params *model.FindByIndexParams) (
  *model.Message, error,
) {
  row := r.db.QueryRowContext(ctx,
    "SELECT id, user_id, createtime, message, file, file_id, log_index, log_term " +
    "FROM messages WHERE log_index = ? AND log_term = ?",
    params.LogIndex, params.LogTerm,
  )

  var msg model.Message
  var fileCol sql.NullString
  var fileIdCol sql.NullString
  var logIndexCol sql.NullInt64
  var logTermCol sql.NullInt64
  if err := row.Scan(
    &msg.Id,
    &msg.UserId,
    &msg.CreateTime,
    &msg.Value,
    &fileCol,
    &fileIdCol,
    &logIndexCol,
    &logTermCol,
  ); err != nil {
    if errors.Is(err, sql.ErrNoRows) {
      return nil, repository.ErrNotFound
    }
    return nil, err
  }
  if fileCol.Valid {
    msg.FileName = fileCol.String
  }
  if fileIdCol.Valid {
    msg.FileId = model.FileId(fileIdCol.String)
  }
  if logIndexCol.Valid {
    msg.LogIndex = uint64(logIndexCol.Int64)
  }
  if logTermCol.Valid {
    msg.LogTerm = uint64(logTermCol.Int64)
  }
  return &msg, nil
}

const ascStmt = `
SELECT id, user_id, createtime, message, file, file_id, log_index, log_term 
FROM messages
WHERE user_id = ?
ORDER BY id ASC
LIMIT ? OFFSET ?
`

const descStmt = `
SELECT id, user_id, createtime, message, file, file_id, log_index, log_term 
FROM messages
WHERE user_id = ?
ORDER BY id DESC
LIMIT ? OFFSET ?
`

func (r *Repository) Get(ctx context.Context, _ *logger.Logger, params *model.GetParams) (
  *model.MessagesList, error,
) {
  var rows *sql.Rows
  var err error
  var isLastPage bool

  if params.Ascending {
    rows, err = r.db.QueryContext(ctx, ascStmt, int(params.UserId), params.Limit, params.Offset)
  } else {
    rows, err = r.db.QueryContext(ctx, descStmt, int(params.UserId), params.Limit, params.Offset)
  }

  if err != nil {
    return nil, err
  }

  defer rows.Close()

  var res []*model.Message
  for rows.Next() {
    var id int
    var userId int
    var createtime string
    var value string
    var fileCol sql.NullString
    var fileIdCol sql.NullString
    var logIndexCol sql.NullInt64
    var logTermCol sql.NullInt64
    if err := rows.Scan(
      &id,
      &userId,
      &createtime,
      &value,
      &fileCol,
      &fileIdCol,
      &logIndexCol,
      &logTermCol,
    ); err != nil {
      return nil, err
    }

    var fileName string
    if fileCol.Valid {
      fileName = fileCol.String
    }
    var fileId string
    if fileIdCol.Valid {
      fileId = fileIdCol.String
    }
    var logIndex uint64
    if logIndexCol.Valid {
      logIndex = uint64(logIndexCol.Int64)
    }
    var logTerm uint64
    if logTermCol.Valid {
      logTerm = uint64(logTermCol.Int64)
    }
    res = append(res, &model.Message{
      Id: model.MessageId(id),
      UserId: userId,
      CreateTime: createtime,
      Value: value,
      FileName: fileName,
      FileId: model.FileId(fileId),
      LogIndex: logIndex,
      LogTerm: logTerm,
    })
  }

  if int32(len(res)) < params.Limit {
    isLastPage = true
  } else {
    row := r.db.QueryRowContext(ctx,
      "SELECT COUNT(*) FROM messages WHERE user_id = ?",
      int(params.UserId),
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

func (r *Repository) GetOne(ctx context.Context, _ *logger.Logger, params *model.GetOneParams) (
  *model.Message, error,
) {
  row := r.db.QueryRowContext(ctx,
    "SELECT id, user_id, createtime, message, file, file_id, log_index, log_term " +
    "FROM messages WHERE user_id = ? AND id = ?",
    int(params.UserId), int(params.MessageId),
  )

  var message model.Message
  var fileCol sql.NullString
  var fileIdCol sql.NullString
  var logIndexCol sql.NullInt64
  var logTermCol sql.NullInt64
  if err := row.Scan(
    &message.Id,
    &message.UserId,
    &message.CreateTime,
    &message.Value,
    &fileCol,
    &fileIdCol,
    &logIndexCol,
    &logTermCol,
  ); err != nil {
    if errors.Is(err, sql.ErrNoRows) {
      return &message, repository.ErrNotFound
    }
    return &message, err
  }
  if fileCol.Valid {
    message.FileName = fileCol.String
  }
  if fileIdCol.Valid {
    message.FileId = model.FileId(fileIdCol.String)
  }
  if logIndexCol.Valid {
    message.LogIndex = uint64(logIndexCol.Int64)
  }
  if logTermCol.Valid {
    message.LogTerm = uint64(logTermCol.Int64)
  }
  return &message, nil
}

func (r *Repository) PutBatch(ctx context.Context, _ *logger.Logger, params *model.PutBatchParams) error {
  for _, message := range params.MessagesList {
    _, err := r.insertSt.ExecContext(ctx,
      message.UserId,
      message.CreateTime,
      message.Value,
      message.FileName,
      message.FileId,
      message.LogIndex,
      message.LogTerm,
    )
    if err != nil {
      return err
    }
  }
  return nil
}

func (r *Repository) GetBatch(ctx context.Context, _ *logger.Logger) ([]*model.Message, error) {
  rows, err := r.db.QueryContext(ctx,
    "SELECT id, user_id, createtime, message, file, file_id, log_index, log_term " +
    "FROM messages",
  )
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  var res []*model.Message
  for rows.Next() {
    var id int
    var userId int
    var value string
    var createtime string
    var fileCol sql.NullString
    var fileIdCol sql.NullString
    var logIndexCol sql.NullInt64
    var logTermCol sql.NullInt64
    if err := rows.Scan(
      &id,
      &userId,
      &createtime,
      &value,
      &fileCol,
      &fileIdCol,
      &logIndexCol,
      &logTermCol,
    ); err != nil {
      return nil, err
    }

    var fileName string
    if fileCol.Valid {
      fileName = fileCol.String
    }
    var fileId string
    if fileIdCol.Valid {
      fileId = fileIdCol.String
    }
    var logIndex uint64
    if logIndexCol.Valid {
      logIndex = uint64(logIndexCol.Int64)
    }
    var logTerm uint64
    if logTermCol.Valid {
      logTerm = uint64(logTermCol.Int64)
    }
    res = append(res, &model.Message{
      Id: model.MessageId(id),
      UserId: userId,
      CreateTime: createtime,
      Value: value,
      FileName: fileName,
      FileId: model.FileId(fileId),
      LogIndex: logIndex,
      LogTerm: logTerm,
    })
  }

  return res, nil
}