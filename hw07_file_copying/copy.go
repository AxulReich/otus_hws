package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
)

const readOnlyPermission = 0444

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrOffsetNegativeValue   = errors.New("invalid negative offset value")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromPath, err := filepath.Abs(fromPath)
	if err != nil {
		return fmt.Errorf("can't get abs path: %v, err: %w", from, err)
	}

	fileFrom, err := os.OpenFile(fromPath, os.O_RDONLY, readOnlyPermission)
	if err != nil {
		return fmt.Errorf("open fromPath err: %w", err)
	}
	defer fileCloser(fileFrom)

	fileFromInfo, err := fileFrom.Stat()
	if err != nil {
		return fmt.Errorf("can't get file info, path: %v err: %w", from, err)
	}

	if !fileFromInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	if offset < 0 {
		return ErrOffsetNegativeValue
	}

	if offset > fileFromInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	_, err = fileFrom.Seek(offset, 0)
	if err != nil {
		return fmt.Errorf("can't get offset from file: %v err: %w", from, err)
	}

	toPath, err = filepath.Abs(toPath)
	if err != nil {
		return fmt.Errorf("can't get abs path: %v, err: %w", toPath, err)
	}
	fileTo, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("can't create or open to file: %v err: %w", to, err)
	}
	defer fileCloser(fileTo)

	if limit == 0 || fileFromInfo.Size()-offset < limit {
		limit = fileFromInfo.Size() - offset
	}

	reader := io.LimitReader(fileFrom, limit)
	bar := pb.Full.Start64(limit)
	barReader := bar.NewProxyReader(reader)
	defer bar.Finish()

	_, err = io.Copy(fileTo, barReader)
	if err != nil {
		return fmt.Errorf("can't copy: %w", err)
	}
	return nil
}

func fileCloser(f *os.File) {
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
