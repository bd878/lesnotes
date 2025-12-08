package distributed

import (
	"io"
	"os"
	"bufio"
	"context"
	"archive/tar"
	"path/filepath"

	"github.com/hashicorp/raft"
	"github.com/bd878/gallery/server/logger"
)

/**
 * Merge two dumps into tar archive
 */

var (
	invoicesFileName = "invoices.bin"
	paymentsFileName = "payments.bin"
	tarFileName      = "billing_snapshot.tar"
)

type snapshot struct {
	tarFile      *os.File
	paymentsFile *os.File
	invoicesFile *os.File
}

func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	// remove file after Persist failure
	err := cleanup()
	if err != nil {
		return nil, err
	}

	s := &snapshot{}

	s.invoicesFile, err = os.CreateTemp("", invoicesFileName)
	if err != nil {
		return nil, err
	}

	s.paymentsFile, err = os.CreateTemp("", paymentsFileName)
	if err != nil {
		return nil, err
	}

	paymentsBuf := bufio.NewWriter(s.paymentsFile)
	defer paymentsBuf.Flush()

	err = f.paymentsRepo.Dump(context.Background(), paymentsBuf)
	if err != nil {
		logger.Errorw("failed to dump payments repo", "error", err)
		return nil, err
	}

	invoicesBuf := bufio.NewWriter(s.invoicesFile)
	defer invoicesBuf.Flush()

	err = f.invoicesRepo.Dump(context.Background(), invoicesBuf)
	if err != nil {
		logger.Errorw("failed to dump invoices repo", "error", err)
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

		if hdr.Name == invoicesFileName {
			err = f.invoicesRepo.Restore(context.Background(), tr)
			if err != nil {
				return err
			}
		} else if hdr.Name == paymentsFileName {
			err = f.paymentsRepo.Restore(context.Background(), tr)
			if err != nil {
				return err
			}
		}
	}

	defer reader.Close()

	return
}

func (s *snapshot) Persist(sink raft.SnapshotSink) (err error) {
	// create tar file
	s.tarFile, err = os.CreateTemp("", tarFileName)
	if err != nil {
		return err
	}

	tarBuf := bufio.NewWriter(s.tarFile)

	tw := tar.NewWriter(tarBuf)

	// seek files
	_, err = s.paymentsFile.Seek(0, os.SEEK_SET)
	if err != nil {
		return
	}

	_, err = s.invoicesFile.Seek(0, os.SEEK_SET)
	if err != nil {
		return
	}

	// write payments to tar
	paymentsSize, err := fileSize(s.paymentsFile)
	if err != nil {
		return err
	}

	paymentsHdr := &tar.Header{
		Name: s.paymentsFile.Name(),
		Mode: 0600,
		Size: paymentsSize,
	}

	err = tw.WriteHeader(paymentsHdr)
	if err != nil {
		return
	}

	_, err = io.Copy(tw, s.paymentsFile)
	if err != nil {
		return
	}

	// write invoices to tar
	invoicesSize, err := fileSize(s.invoicesFile)
	if err != nil {
		return err
	}

	invoicesHdr := &tar.Header{
		Name: s.invoicesFile.Name(),
		Mode: 0600,
		Size: invoicesSize,
	}

	err = tw.WriteHeader(invoicesHdr)
	if err != nil {
		return
	}

	_, err = io.Copy(tw, s.invoicesFile)
	if err != nil {
		return
	}

	// dump tar to sink
	if err = tarBuf.Flush(); err != nil {
		return
	}

	if err = tw.Close(); err != nil {
		return
	}

	_, err = s.tarFile.Seek(0, os.SEEK_SET)
	if err != nil {
		return
	}

	tr := tar.NewReader(s.tarFile)

	_, err = io.Copy(sink, tr)
	defer sink.Cancel()
	if err != nil {
		return
	}

	return sink.Close()
}

func (s *snapshot) Release() {
	if err := s.tarFile.Close(); err != nil {
		logger.Errorw("failed to close tar file", "error", err)
	}

	if err := s.paymentsFile.Close(); err != nil {
		logger.Errorw("failed to close payments file", "error", err)
	}

	if err := s.invoicesFile.Close(); err != nil {
		logger.Errorw("failed to close invoices file", "error", err)
	}

	if err := os.Remove(filepath.Join(os.TempDir(), s.tarFile.Name())); err != nil {
		logger.Errorw("cannot remove tar file", "error", err)
	}

	if err := os.Remove(filepath.Join(os.TempDir(), s.paymentsFile.Name())); err != nil {
		logger.Errorw("cannot remove payments file", "error", err)
	}

	if err := os.Remove(filepath.Join(os.TempDir(), s.invoicesFile.Name())); err != nil {
		logger.Errorw("cannot remove invoices file", "error", err)
	}
}

/**
 * Utils
 */
func cleanup() (err error) {
	_, err = os.Stat(tarFileName)
	if err == nil || err != os.ErrNotExist {
		err = os.Remove(filepath.Join(os.TempDir(), tarFileName))
		if err != nil {
			return
		}
	}

	_, err = os.Stat(invoicesFileName)
	if err == nil || err != os.ErrNotExist {
		err = os.Remove(filepath.Join(os.TempDir(), invoicesFileName))
		if err != nil {
			return
		}
	}

	_, err = os.Stat(paymentsFileName)
	if err == nil || err != os.ErrNotExist {
		err = os.Remove(filepath.Join(os.TempDir(), paymentsFileName))
		if err != nil {
			return
		}
	}

	return
}

func fileSize(f *os.File) (size int64, err error) {
	info, err := f.Stat()
	if err != nil {
		return 0, err
	}

	size = info.Size()

	return
}
