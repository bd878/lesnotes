package distributed

import (
	"io"
	"os"
	"bufio"
	"strings"
	"context"
	"archive/tar"

	"github.com/hashicorp/raft"
	"github.com/bd878/gallery/server/logger"
)

/**
 * Merge two dumps into tar archive
 */

type snapshot struct {
	tarFile      *os.File
	messagesFile *os.File
	filesFile    *os.File
	translationsFile *os.File
}

func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	var err error

	s := &snapshot{}

	s.filesFile, err = os.CreateTemp("", "files_*.bin")
	if err != nil {
		return nil, err
	}

	s.translationsFile, err = os.CreateTemp("", "translations_*.bin")
	if err != nil {
		return nil, err
	}

	s.messagesFile, err = os.CreateTemp("", "messages_*.bin")
	if err != nil {
		return nil, err
	}

	s.tarFile, err = os.CreateTemp("", "messages_*.tar")
	if err != nil {
		return nil, err
	}

	messagesBuf := bufio.NewWriter(s.messagesFile)
	defer messagesBuf.Flush()

	err = f.messagesRepo.Dump(context.Background(), messagesBuf)
	if err != nil {
		logger.Errorw("failed to dump messages repo", "error", err)
		return nil, err
	}

	filesBuf := bufio.NewWriter(s.filesFile)
	defer filesBuf.Flush()

	err = f.filesRepo.Dump(context.Background(), filesBuf)
	if err != nil {
		logger.Errorw("failed to dump files repo", "error", err)
		return nil, err
	}

	translationsBuf := bufio.NewWriter(s.translationsFile)
	defer translationsBuf.Flush()

	err = f.translationsRepo.Dump(context.Background(), translationsBuf)
	if err != nil {
		logger.Errorw("failed to dump translations repo", "error", err)
		return nil, err
	}

	return s, nil
}

func (f *fsm) Restore(reader io.ReadCloser) (err error) {
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

		if strings.Contains(hdr.Name, "files") {
			err = f.filesRepo.Restore(context.Background(), tr)
			if err != nil {
				return err
			}
		} else if strings.Contains(hdr.Name, "messages") {
			err = f.messagesRepo.Restore(context.Background(), tr)
			if err != nil {
				return err
			}
		} else if strings.Contains(hdr.Name, "translations") {
			err = f.translationsRepo.Restore(context.Background(), tr)
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
	err = s.messagesFile.Sync()
	if err != nil {
		return
	}

	_, err = s.messagesFile.Seek(0, os.SEEK_SET)
	if err != nil {
		return
	}

	err = s.filesFile.Sync()
	if err != nil {
		return
	}

	_, err = s.filesFile.Seek(0, os.SEEK_SET)
	if err != nil {
		return
	}

	err = s.translationsFile.Sync()
	if err != nil {
		return
	}

	_, err = s.translationsFile.Seek(0, os.SEEK_SET)
	if err != nil {
		return
	}

	// write messages to tar
	messagesSize, err := fileSize(s.messagesFile)
	if err != nil {
		return err
	}

	logger.Debugw("persisting", "messages size", messagesSize)

	messagesHdr := &tar.Header{
		Name: s.messagesFile.Name(),
		Mode: 0600,
		Size: messagesSize,
	}

	err = tw.WriteHeader(messagesHdr)
	if err != nil {
		return
	}

	n, err := io.Copy(tw, s.messagesFile)
	if err != nil {
		return
	}

	if err = tw.Flush(); err != nil {
		return
	}

	logger.Debugw("copied messages bytes to tar writer", "bytes", n)

	// write files to tar
	filesSize, err := fileSize(s.filesFile)
	if err != nil {
		return err
	}

	logger.Debugw("persisting", "files size", filesSize)

	filesHdr := &tar.Header{
		Name: s.filesFile.Name(),
		Mode: 0600,
		Size: filesSize,
	}

	err = tw.WriteHeader(filesHdr)
	if err != nil {
		return
	}

	n, err = io.Copy(tw, s.filesFile)
	if err != nil {
		return
	}

	if err = tw.Flush(); err != nil {
		return
	}

	logger.Debugw("copied files bytes to tar writer", "bytes", n)

	// write translations to tar
	translationsSize, err := fileSize(s.translationsFile)
	if err != nil {
		return err
	}

	logger.Debugw("persisting", "translations size", translationsSize)

	translationsHdr := &tar.Header{
		Name: s.translationsFile.Name(),
		Mode: 0600,
		Size: translationsSize,
	}

	err = tw.WriteHeader(translationsHdr)
	if err != nil {
		return
	}

	n, err = io.Copy(tw, s.translationsFile)
	if err != nil {
		return
	}

	if err = tw.Flush(); err != nil {
		return
	}

	logger.Debugw("copied translations bytes to tar writer", "bytes", n)

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

	if err := os.Remove(s.messagesFile.Name()); err != nil {
		logger.Errorw("cannot remove messages file", "error", err)
	}

	if err := os.Remove(s.filesFile.Name()); err != nil {
		logger.Errorw("cannot remove files file", "error", err)
	}

	if err := os.Remove(s.translationsFile.Name()); err != nil {
		logger.Errorw("cannot remove translations file", "error", err)
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
