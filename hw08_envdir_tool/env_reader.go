package main

import (
	"bufio"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type Environment map[string]EnvValue

const name =

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

type scanData struct {
	isSuccess bool
	envName   string
	EnvValue
	error
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
		workerNum = runtime.NumCPU()
		doneCh    = make(chan struct{})
		fileCh    = make(chan fs.FileInfo)
		resCh     = make(chan scanData)

		mu     sync.WaitGroup
		result = make(Environment)

		wg sync.WaitGroup
	)

	wg.Add(workerNum)
	for i := 0; i < workerNum; i++ {
		go func() {
			defer wg.Done()
			for file := range fileCh {
				select {
				case <-doneCh:
					return
				default:
				}

				select {
				case <-doneCh:
					return
				default:
					if !file.IsDir() && file.Mode().IsRegular() {

					}
				}

			}
		}()
	}

	for j := range files {
		fileCh <- files[j]
	}

	return nil, nil
}

func envScan(fileInfo fs.FileInfo, dirPath string) scanData {
	var result scanData

	if !fileInfo.IsDir() && fileInfo.Mode().IsRegular() && !strings.Contains(fileInfo.Name(), "=") {
		file, err := os.Open(filepath.Join(dirPath, fileInfo.Name()))
		if err != nil {
			log.Fatalf("failed opening fileInfo: %s", err)
		}

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		var (
			line          string
			hasSecondLine bool
			scanCounter   int64
		)

		for scanner.Scan() {
			line = scanner.Text()

			if line == "" {

			}
		}

	}
}
