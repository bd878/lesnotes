package machine

import (
	"io"
	"os"
	"bufio"
	"strings"
	"context"
	"archive/tar"

	"github.com/hashicorp/raft"
	"github.com/bd878/gallery/server/internal/logger"
)

/**
 * Merge two dumps into tar archive
 */

type snapshot struct {
	tarFile      *os.File
	sessionsFile *os.File
}

func (f *Machine) Snapshot() (raft.FSMSnapshot, error) {
	var err error

	s := &snapshot{}

	s.sessionsFile, err = os.CreateTemp("", "sessions_*.bin")
	if err != nil {
		return nil, err
	}

	s.tarFile, err = os.CreateTemp("", "sessions_*.tar")
	if err != nil {
		return nil, err
	}

	sessionsBuf := bufio.NewWriter(s.sessionsFile)
	defer sessionsBuf.Flush()

	err = f.sessionsRepo.Dump(context.Background(), sessionsBuf)
	if err != nil {
		logger.Errorw("failed to dump sessions repo", "error", err)
		return nil, err
	}

	return s, nil
}

func (f *Machine) Restore(reader io.ReadCloser) (err error) {
	logger.Debugln("restoring fsm from snapshot")

	tr := tar.NewReader(reader)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if strings.Contains(hdr.Name, "sessions") {
			err = f.sessionsRepo.Restore(context.Background(), tr)
			if err != nil {
				return err
			}
		}
	}

	defer reader.Close()

	return
}

func (s *snapshot) Persist(sink raft.SnapshotSink) (err error) {
	tarBuf := bufio.NewWriter(s.tarFile)
	tw := tar.NewWriter(tarBuf)

	// seek files
	err = s.sessionsFile.Sync()
	if err != nil {
		return
	}

	_, err = s.sessionsFile.Seek(0, os.SEEK_SET)
	if err != nil {
		return
	}

	// write sessions to tar
	sessionsSize, err := fileSize(s.sessionsFile)
	if err != nil {
		return err
	}

	logger.Debugw("persisting", "sessions size", sessionsSize)

	sessionsHdr := &tar.Header{
		Name: s.sessionsFile.Name(),
		Mode: 0600,
		Size: sessionsSize,
	}

	err = tw.WriteHeader(sessionsHdr)
	if err != nil {
		return
	}

	n, err := io.Copy(tw, s.sessionsFile)
	if err != nil {
		return
	}

	if err = tw.Flush(); err != nil {
		return
	}

	logger.Debugw("copied sessions bytes to tar writer", "bytes", n)

	// dump tar to sink
	if err = tw.Close(); err != nil {
		return
	}

	if err = tarBuf.Flush(); err != nil {
		return
	}

	err = s.tarFile.Sync()
	if err != nil {
		return
	}

	_, err = s.tarFile.Seek(0, os.SEEK_SET)
	if err != nil {
		return
	}

	tarSize, err := fileSize(s.tarFile)
	if err != nil {
		return err
	}

	logger.Debugw("tar size", "bytes", tarSize)

	n, err = io.Copy(sink, s.tarFile)
	defer sink.Cancel()
	if err != nil {
		return err
	}

	logger.Debugw("persisting", "copied bytes", n)

	return sink.Close()
}

func (s *snapshot) Release() {
	if err := os.Remove(s.tarFile.Name()); err != nil {
		logger.Errorw("cannot remove tar file", "error", err)
	}

	if err := os.Remove(s.sessionsFile.Name()); err != nil {
		logger.Errorw("cannot remove sessions file", "error", err)
	}
}

/**
 * Utils
 */

func fileSize(f *os.File) (size int64, err error) {
	info, err := f.Stat()
	if err != nil {
		return 0, err
	}

	size = info.Size()

	return
}
