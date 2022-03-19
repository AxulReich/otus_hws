package main

import (
	"errors"
	"os"
	"path/filepath"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromPath, err := filepath.Abs(fromPath)
	if err != nil {
		return err
	}
	toPath, err = filepath.Abs(toPath)
	if err != nil {
		return err
	}

	fileFrom, err := os.Open(fromPath)
	if err != nil {
		return err
	}

	if _, err = os.Stat(toPath); !os.IsNotExist(err) {
		return err
	}

	return nil
}
