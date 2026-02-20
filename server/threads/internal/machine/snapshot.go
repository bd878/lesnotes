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
	threadsFile *os.File
}

func (f *Machine) Snapshot() (raft.FSMSnapshot, error) {
	var err error

	s := &snapshot{}

	s.threadsFile, err = os.CreateTemp("", "threads_*.bin")
	if err != nil {
		return nil, err
	}

	s.tarFile, err = os.CreateTemp("", "threads_*.tar")
	if err != nil {
		return nil, err
	}

	threadsBuf := bufio.NewWriter(s.threadsFile)
	defer threadsBuf.Flush()

	err = f.threadsRepo.Dump(context.Background(), threadsBuf)
	if err != nil {
		logger.Errorw("failed to dump threads repo", "error", err)
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

		if strings.Contains(hdr.Name, "threads") {
			err = f.threadsRepo.Restore(context.Background(), tr)
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
	err = s.threadsFile.Sync()
	if err != nil {
		return
	}

	_, err = s.threadsFile.Seek(0, os.SEEK_SET)
	if err != nil {
		return
	}

	// write threads to tar
	threadsSize, err := fileSize(s.threadsFile)
	if err != nil {
		return err
	}

	logger.Debugw("persisting", "threads size", threadsSize)

	threadsHdr := &tar.Header{
		Name: s.threadsFile.Name(),
		Mode: 0600,
		Size: threadsSize,
	}

	err = tw.WriteHeader(threadsHdr)
	if err != nil {
		return
	}

	n, err := io.Copy(tw, s.threadsFile)
	if err != nil {
		return
	}

	if err = tw.Flush(); err != nil {
		return
	}

	logger.Debugw("copied threads bytes to tar writer", "bytes", n)

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

	if err := os.Remove(s.threadsFile.Name()); err != nil {
		logger.Errorw("cannot remove threads file", "error", err)
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
