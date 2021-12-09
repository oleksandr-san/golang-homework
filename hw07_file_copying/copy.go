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

func min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromStat, err := os.Stat(fromPath)
	if err != nil {
		return err
	}

	if fromStat.Size() == 0 {
		return ErrUnsupportedFile
	}

	if offset > fromStat.Size() {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 {
		limit = fromStat.Size() - offset
	} else {
		limit = min(fromStat.Size()-offset, limit)
	}

	fromFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}

	_, err = fromFile.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}

	bar := pb.Full.Start64(limit)
	defer bar.Finish()

	_, err = io.CopyN(toFile, bar.NewProxyReader(fromFile), limit)
	return err
}
