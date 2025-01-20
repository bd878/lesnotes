package sqlite

const BACKUP_ALL_PAGES = -1

func (s *SqliteSnapshotSink) Write(_ []byte) (int, error) {
  var err error
  s.backup, err = s.targetConn.Backup("main", s.srcConn, "main")
  if err != nil {
    return 0, err
  }

  var done bool = false
  for !done {
    done, err = s.backup.Step(BACKUP_ALL_PAGES)
    if err != nil {
      return 0, err
    }
  }

  return 0, nil
}