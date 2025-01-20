package sqlite

import (
  "time"
  "fmt"
  "io"
  "bufio"
  "errors"
  "strings"
  "os"
  "sort"
  "path/filepath"
  "encoding/json"
  "database/sql"

  sqlite3 "github.com/mattn/go-sqlite3"

  "github.com/hashicorp/raft"
  "github.com/bd878/gallery/server/logger"
)

const (
  snapPath      = "snapshots"
  metaFilePath  = "meta.json"
  stateFilePath = "state.bin"
  tmpSuffix     = ".tmp"
)

type SqliteSnapshotStore struct {
  log         *logger.Logger
  path         string
  dbFilePath   string
}

var _ raft.SnapshotStore = (*SqliteSnapshotStore)(nil)

type snapMetaSlice []*sqliteSnapshotMeta

type sqliteSnapshotMeta struct {
  raft.SnapshotMeta
}

type SqliteSnapshotSink struct {
  log         *logger.Logger
  dir          string
  parentDir    string
  srcConn     *sqlite3.SQLiteConn
  targetConn  *sqlite3.SQLiteConn
  backup      *sqlite3.SQLiteBackup
  store       *SqliteSnapshotStore
  meta         sqliteSnapshotMeta
  statePath    string

  closed       bool
}

var _ raft.SnapshotSink = (*SqliteSnapshotSink)(nil)

func New(base, dbFilePath string, log *logger.Logger) *SqliteSnapshotStore {
  path := filepath.Join(base, snapPath)
  if err := os.MkdirAll(path, 0o755); err != nil && !os.IsExist(err) {
    panic(err)
  }

  store := &SqliteSnapshotStore{
    path:       path,
    dbFilePath: dbFilePath,
    log:        log,
  }

  return store
}

func (s *SqliteSnapshotStore) Create(version raft.SnapshotVersion, index, term uint64, configuration raft.Configuration,
  configurationIndex uint64, trans raft.Transport) (raft.SnapshotSink, error) {
  if version != 1 {
    return nil, fmt.Errorf("unsupported snapshot version %d", version)
  }

  name := snapshotName(term, index)

  path := filepath.Join(s.path, name+tmpSuffix)
  s.log.Info("creating new snapshot", "path", path)

  sink := &SqliteSnapshotSink{
    log:        s.log,
    store:      s,
    dir:        path,
    parentDir:  s.path,
    meta:       sqliteSnapshotMeta{
      SnapshotMeta: raft.SnapshotMeta{
        Version:              version,
        ID:                   name,
        Index:                index,
        Term:                 term,
        Configuration:        configuration,
        ConfigurationIndex:   configurationIndex,
      },
    },
  }

  if err := sink.writeMeta(); err != nil {
    s.log.Error("failed to write metadata", "error", err)
    return nil, err
  }

  statePath := filepath.Join(path, stateFilePath)

  srcDb, err := sql.Open("sqlite3", "file:" + s.dbFilePath)
  if err != nil {
    return nil, err
  }

  srcConn, err := srcDb.Driver().Open("sqlite3")
  defer srcConn.Close()
  if err != nil {
    return nil, err
  }

  sqliteSrcConn, ok := srcConn.(*sqlite3.SQLiteConn)
  if !ok {
    return nil, errors.New("driver conn is not SQLiteConn")
  }

  targetDb, err := sql.Open("sqlite3", "file:" + statePath)
  if err != nil {
    return nil, err
  }

  targetConn, err := targetDb.Driver().Open("sqlite3")
  defer targetConn.Close()
  if err != nil {
    return nil, err
  }

  sqliteTargetConn, ok := targetConn.(*sqlite3.SQLiteConn)
  if !ok {
    return nil, errors.New("driver conn is not SQLiteConn")
  }

  sink.statePath = statePath
  sink.srcConn = sqliteSrcConn
  sink.targetConn = sqliteTargetConn

  return sink, nil
}

func (sink *SqliteSnapshotSink) writeMeta() error {
  metaPath := filepath.Join(sink.dir, metaFilePath)
  fh, err := os.Create(metaPath)
  if err != nil {
    return err
  }
  defer fh.Close()

  buffered := bufio.NewWriter(fh)

  enc := json.NewEncoder(buffered)
  if err := enc.Encode(&sink.meta); err != nil {
    return err
  }

  if err := buffered.Flush(); err != nil {
    return err
  }

  if err := fh.Sync(); err != nil {
    return err
  }

  return nil
}

func snapshotName(term, index uint64) string {
  now := time.Now()
  msec := now.UnixNano() / int64(time.Millisecond)
  return fmt.Sprintf("%d-%d-%d", term, index, msec)
}

func (s *SqliteSnapshotStore) List() ([]*raft.SnapshotMeta, error) {
  snapshots, err := s.getSnapshots()
  if err != nil {
    s.log.Error("failed to get snapshots", "error", err)
    return nil, err
  }

  var snapMeta []*raft.SnapshotMeta
  for _, meta := range snapshots {
    snapMeta = append(snapMeta, &meta.SnapshotMeta)
  }
  return snapMeta, nil
}

func (s *SqliteSnapshotStore) getSnapshots() ([]*sqliteSnapshotMeta, error) {
  snapshots, err := os.ReadDir(s.path)
  if err != nil {
    s.log.Error("failed to scan snapshots directory", "error", err)
    return nil, err
  }

  var snapMeta []*sqliteSnapshotMeta
  for _, snap := range snapshots {
    if !snap.IsDir() {
      continue
    }

    dirName := snap.Name()
    if strings.HasSuffix(dirName, tmpSuffix) {
      s.log.Warn("found temporary snapshot", "name", dirName)
      continue
    }

    meta, err := s.readMeta(dirName)
    if err != nil {
      s.log.Warn("failed to read metadata", "name", dirName, "error", err)
      continue
    }

    if meta.Version < raft.SnapshotVersionMin || meta.Version > raft.SnapshotVersionMax {
      s.log.Warn("snapshot version not supported", "name", dirName, "version", meta.Version)
      continue
    }

    snapMeta = append(snapMeta, meta)
  }

  sort.Sort(sort.Reverse(snapMetaSlice(snapMeta)))

  return snapMeta, nil
}

func (s *SqliteSnapshotStore) readMeta(name string) (*sqliteSnapshotMeta, error) {
  metaPath := filepath.Join(s.path, name, metaFilePath)
  fh, err := os.Open(metaPath)
  if err != nil {
    return nil, err
  }
  defer fh.Close()

  buffered := bufio.NewReader(fh)

  meta := &sqliteSnapshotMeta{}
  dec := json.NewDecoder(buffered)
  if err := dec.Decode(meta); err != nil {
    return nil, err
  }
  return meta, err
}

func (s *SqliteSnapshotStore) ReapSnapshots() error {
  snapshots, err := s.getSnapshots()
  if err != nil {
    s.log.Error("failed to get snapshots", "error", err)
    return err
  }

  path := filepath.Join(s.path, snapshots[0].ID)
  s.log.Info("reaping snapshot", "path", path)
  if err := os.RemoveAll(path); err != nil {
    s.log.Error("failed to reap snapshot", "path", path, "error", err)
    return err
  }

  return nil
}

type bufferedFile struct {
  bh *bufio.Reader
  fh *os.File
}

func (b *bufferedFile) Read(p []byte) (n int, err error) {
  return b.bh.Read(p)
}

func (b *bufferedFile) Close() error {
  return b.fh.Close()
}

func (s *SqliteSnapshotStore) Open(id string) (*raft.SnapshotMeta, io.ReadCloser, error) {
  meta, err := s.readMeta(id)
  if err != nil {
    s.log.Error("failed to get meta data to open snapshot", "error", err)
    return nil, nil, err
  }

  statePath := filepath.Join(s.path, id, stateFilePath)
  fh, err := os.Open(statePath)
  if err != nil {
    s.log.Error("failed to open state file", "error", err)
    return nil, nil, err
  }

  if _, err := fh.Seek(0, 0); err != nil {
    s.log.Error("state file seek failed", "error", err)
    fh.Close()
    return nil, nil, err
  }

  buffered := &bufferedFile{bufio.NewReader(fh), fh}

  return &meta.SnapshotMeta, buffered, nil
}

func (s *SqliteSnapshotSink) Close() error {
  if s.closed {
    return nil
  }
  s.closed = true

  if err := s.finalize(); err != nil {
    s.log.Error("failed to finalize snapshot", "error", err)
    if delErr := os.RemoveAll(s.dir); delErr != nil {
      s.log.Error("failed to delete temp snapshot dir", "path", s.dir, "error", "delErr")
      return delErr
    }
    return err
  }

  if err := s.writeMeta(); err != nil {
    s.log.Error("failed to write metadata", "error", err)
    return err
  }

  newPath := strings.TrimSuffix(s.dir, tmpSuffix)
  if err := os.Rename(s.dir, newPath); err != nil {
    s.log.Error("failed to move snapshot into place", "error", err)
    return err
  }

  parentFh, err := os.Open(s.parentDir)
  if err != nil {
    s.log.Error("faield to open snapshot parent directory", "path", s.parentDir, "error", err)
    return err
  }
  defer parentFh.Close()

  if err := parentFh.Sync(); err != nil {
    s.log.Error("failed syncing parent dir", "path", s.parentDir, "error", err)
    return err
  }

  // Reap any old snapshots
  if err := s.store.ReapSnapshots(); err != nil {
    return err
  }

  return nil
}

func (s *SqliteSnapshotSink) ID() string {
  return s.meta.ID  
}

func (s *SqliteSnapshotSink) Cancel() error {
  if s.closed {
    return nil
  }

  s.closed = true
  if err := s.finalize(); err != nil {
    s.log.Error("failed to finalize snapshot", "error", err)
    return err
  }

  return os.RemoveAll(s.dir)
}

func (s *SqliteSnapshotSink) finalize() error {
  if err := s.backup.Finish(); err != nil {
    return err
  }

  fh, err := os.Open(s.statePath)
  if err != nil {
    return err
  }

  stat, stateErr := fh.Stat()
  if stateErr != nil {
    return stateErr
  }

  s.meta.Size = stat.Size()

  return nil
}

func (s snapMetaSlice) Len() int {
  return len(s)
}

func (s snapMetaSlice) Less(i, j int) bool {
  if s[i].Term != s[j].Term {
    return s[i].Term < s[j].Term
  }
  if s[i].Index != s[j].Index {
    return s[i].Index < s[j].Index
  }
  return s[i].ID < s[j].ID
}

func (s snapMetaSlice) Swap(i, j int) {
  s[i], s[j] = s[j], s[i]
}
