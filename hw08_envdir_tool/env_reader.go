package main

import (
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"sync"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	absPath, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	files, err := ioutil.ReadDir(absPath)
	if err != nil {
		return nil, err
	}

	var (
		doneCh = make(chan struct{})
		fileCh = make(chan fs.FileInfo)

		mu      sync.WaitGroup
		result  = make(Environment)
		wg      sync.WaitGroup
		envScan = func() {
			defer wg.Done()
			for file := range fileCh {
				select {
				case <-doneCh:
					return
				default:
					if !file.IsDir() && file.Mode().IsRegular() {

					}
				}

			}
		}
	)

	for _, file := range files {
		if !file.IsDir() && file.Mode().IsRegular() {

		}

	}

	return nil, nil
}

func envScan(done chan struct{}, data chan fs.FileInfo) {
}
