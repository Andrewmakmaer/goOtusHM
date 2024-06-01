package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == "/dev/urandom" || fromPath == "/dev/random" {
		return ErrUnsupportedFile
	}

	fileStat, err := os.Stat(fromPath)
	if err != nil {
		return err
	}

	N := fileStat.Size() - offset
	if N < 0 {
		return ErrOffsetExceedsFileSize
	} else if limit > 0 {
		N = limit
	}

	file, err := os.OpenFile(fromPath, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Seek(offset, io.SeekStart)

	fileReader := io.LimitReader(file, N)
	bar := pb.Full.Start64(N)
	barReader := bar.NewProxyReader(fileReader)

	dstFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	io.Copy(dstFile, barReader)
	bar.Finish()

	return nil
}
