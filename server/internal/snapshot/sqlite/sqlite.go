package sqlite

import (
  "github.com/bd878/gallery/server/logger"
)

const BACKUP_ALL_PAGES = -1

func (s *SqliteSnapshotSink) Write(_ []byte) (int, error) {
  var err error
  logger.Info("writing snapshot ", "from=", s.srcConn.GetFilename("main"), "to=", s.targetConn.GetFilename("main"))
  s.backup, err = s.targetConn.Backup("main", s.srcConn, "main")
  if err != nil {
    return 0, err
  }

  var done bool = false
  for !done {
    done, err = s.backup.Step(BACKUP_ALL_PAGES)
    if err != nil {
      logger.Error("step error=", err)
      return 0, err
    }
    logger.Info("done=", done)
  }

  logger.Info("done")

  return 0, nil
}