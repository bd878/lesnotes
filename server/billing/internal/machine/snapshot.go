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
	paymentsFile *os.File
	invoicesFile *os.File
}

func (f *Machine) Snapshot() (raft.FSMSnapshot, error) {
	var err error

	s := &snapshot{}

	s.invoicesFile, err = os.CreateTemp("", "invoices_*.bin")
	if err != nil {
		return nil, err
	}

	s.paymentsFile, err = os.CreateTemp("", "payments_*.bin")
	if err != nil {
		return nil, err
	}

	s.tarFile, err = os.CreateTemp("", "billing_*.tar")
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

		if strings.Contains(hdr.Name, "invoices") {
			err = f.invoicesRepo.Restore(context.Background(), tr)
			if err != nil {
				return err
			}
		} else if strings.Contains(hdr.Name, "payments") {
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
	tarBuf := bufio.NewWriter(s.tarFile)
	tw := tar.NewWriter(tarBuf)

	// seek files
	err = s.paymentsFile.Sync()
	if err != nil {
		return
	}

	_, err = s.paymentsFile.Seek(0, os.SEEK_SET)
	if err != nil {
		return
	}

	err = s.invoicesFile.Sync()
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

	logger.Debugw("persisting", "payments size", paymentsSize)

	paymentsHdr := &tar.Header{
		Name: s.paymentsFile.Name(),
		Mode: 0600,
		Size: paymentsSize,
	}

	err = tw.WriteHeader(paymentsHdr)
	if err != nil {
		return
	}

	n, err := io.Copy(tw, s.paymentsFile)
	if err != nil {
		return
	}

	if err = tw.Flush(); err != nil {
		return
	}

	logger.Debugw("copied payments bytes to tar writer", "bytes", n)

	// write invoices to tar
	invoicesSize, err := fileSize(s.invoicesFile)
	if err != nil {
		return err
	}

	logger.Debugw("persisting", "invoices size", invoicesSize)

	invoicesHdr := &tar.Header{
		Name: s.invoicesFile.Name(),
		Mode: 0600,
		Size: invoicesSize,
	}

	err = tw.WriteHeader(invoicesHdr)
	if err != nil {
		return
	}

	n, err = io.Copy(tw, s.invoicesFile)
	if err != nil {
		return
	}

	if err = tw.Flush(); err != nil {
		return
	}

	logger.Debugw("copied invoices bytes to tar writer", "bytes", n)

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

	if err := os.Remove(s.paymentsFile.Name()); err != nil {
		logger.Errorw("cannot remove payments file", "error", err)
	}

	if err := os.Remove(s.invoicesFile.Name()); err != nil {
		logger.Errorw("cannot remove invoices file", "error", err)
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
