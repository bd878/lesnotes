package repository

import (
  "context"
  "errors"
  "database/sql"

  _ "github.com/mattn/go-sqlite3"
  "github.com/bd878/gallery/server/messages/pkg/model"
  "github.com/bd878/gallery/server/messages/internal/repository"
  usermodel "github.com/bd878/gallery/server/users/pkg/model"
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

func (r *Repository) Put(ctx context.Context, msg *model.Message) (model.MessageId, error) {
  res, err := r.insertSt.ExecContext(ctx,
    msg.UserId,
    msg.CreateTime,
    msg.Value,
    msg.FileName,
    msg.FileId,
    msg.LogIndex,
    msg.LogTerm,
  )
  if err != nil {
    return model.NullMsgId, err
  }
  id, _ := res.LastInsertId()
  return model.MessageId(id), nil
}

func (r *Repository) Truncate(ctx context.Context) error {
  _, err := r.db.ExecContext(ctx,
    "DELETE FROM messages",
  )
  if err != nil {
    return err
  }
  return nil
}

func (r *Repository) FindByIndexTerm(ctx context.Context, logIndex, logTerm uint64) (*model.Message, error) {
  row := r.db.QueryRowContext(ctx,
    "SELECT id, user_id, createtime, message, file, file_id, log_index, log_term " +
    "FROM messages WHERE log_index = ? AND log_term = ?",
    logIndex, logTerm,
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

func (r *Repository) Get(ctx context.Context, userId usermodel.UserId, limit, offset int32) ([]*model.Message, error) {
  rows, err := r.db.QueryContext(ctx,
    "SELECT id, user_id, createtime, message, file, file_id, log_index, log_term " +
    "FROM messages WHERE user_id = ? LIMIT ? OFFSET ?",
    int(userId), limit, offset,
  )
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
  return res, nil
}

func (r *Repository) GetOne(ctx context.Context, userId usermodel.UserId, id model.MessageId) (*model.Message, error) {
  row := r.db.QueryRowContext(ctx,
    "SELECT id, user_id, createtime, message, file, file_id, log_index, log_term " +
    "FROM messages WHERE user_id = ? AND id = ?",
    int(userId), int(id),
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
      return &msg, repository.ErrNotFound
    }
    return &msg, err
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

func (r *Repository) PutBatch(ctx context.Context, batch []*model.Message) error {
  for _, msg := range batch {
    _, err := r.insertSt.ExecContext(ctx,
      msg.UserId,
      msg.CreateTime,
      msg.Value,
      msg.FileName,
      msg.FileId,
      msg.LogIndex,
      msg.LogTerm,
    )
    if err != nil {
      return err
    }
  }
  return nil
}

func (r *Repository) GetBatch(ctx context.Context) ([]*model.Message, error) {
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