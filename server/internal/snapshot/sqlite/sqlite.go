package sqlite

import (
  "errors"
  "database/sql"
  "google.golang.org/protobuf/proto"
)

const BACKUP_ALL_PAGES = -1

func (s *SqliteSnapshotSink) Write(p []byte) (int, error) {
  var params BackupParams
  proto.Unmarshal(p, &params)

  makeBackup(params.DbFilePath)

  return 0, nil
}

func makeBackup(dbFilePath, backupPath string) error {
  srcDb, err := sql.Open("sqlite3", "file:" + dbFilePath)
  defer srcConn.Close()
  if err != nil {
    return err
  }

  srcConn, err := srcDb.Driver.Open("sqlite3")
  defer srcConn.Close()
  if err != nil {
    return err
  }

  sqliteSrcConn, ok := srcConn.(*sqlite3.SQLiteConn)
  if !ok {
    return errors.New("driver conn is not SQLiteConn")
  }

  targetDb, err := sql.Open("sqlite3", "file:" + backupPath)
  if err != nil {
    return err
  }

  targetConn, err := targetDb.Driver.Open("sqlite3")
  defer targetConn.Close()
  if err != nil {
    return err
  }

  sqliteTargetConn, ok := targetConn.(*sqlite3.SQLiteConn)
  if !ok {
    return errors.New("driver conn is not SQLiteConn")
  }

  backup, err := sqliteTargetConn.Backup("main", sqliteSrcConn, "main")
  if err != nil {
    backup.Close()
    return err
  }

  var done bool = false
  for !done {
    var err error
    done, err = backup.Step(BACKUP_ALL_PAGES)
    if err != nil {
      backup.Close()
      return err
    }
  }

  if err := backup.Finish(); err != nil {
    return err
  }
}