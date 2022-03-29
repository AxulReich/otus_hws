package main

import (
	"bufio"
	"bytes"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

type scanData struct {
	envName string
	EnvValue
	error
}

const (
	cutSet           = " \t"
	terminalZeroChar = "\x00"
	newLineChar      = "\n"
	forbiddenChar    = "="
)

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
		fileCh = make(chan fs.FileInfo, len(files))
		result = make(Environment)
	)

	resCh := runScan(doneCh, fileCh, absPath)
	defer close(doneCh)

	for j := range files {
		fileCh <- files[j]
	}
	close(fileCh)

	for res := range resCh {
		if res.error == nil && res.envName != "" {
			result[res.envName] = res.EnvValue
		}
	}

	return result, nil
}

func runScan(doneCh chan struct{}, inCh chan fs.FileInfo, absPath string) <-chan scanData {
	workerNum := runtime.NumCPU()

	wg := &sync.WaitGroup{}
	resCh := make(chan scanData)

	go func() {
		wg.Add(workerNum)
		for i := 0; i < workerNum; i++ {
			go func() {
				defer wg.Done()
				for {
					select {
					case <-doneCh:
						return
					default:
					}

					select {
					case <-doneCh:
						return
					case file, ok := <-inCh:
						if !ok {
							return
						}
						resCh <- envDirScan(file, absPath)
					}
				}
			}()
		}
		wg.Wait()
		close(resCh)
	}()
	return resCh
}

func envDirScan(fileInfo fs.FileInfo, dirPath string) scanData {
	var result scanData

	if !fileInfo.IsDir() && fileInfo.Mode().IsRegular() && !strings.Contains(fileInfo.Name(), forbiddenChar) {
		file, err := os.Open(filepath.Join(dirPath, fileInfo.Name()))

		if err == nil {
			result.envName = fileInfo.Name()
			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)

			for scanner.Scan() {
				result.EnvValue.Value = string(bytes.ReplaceAll(
					[]byte(strings.TrimRight(scanner.Text(), cutSet)),
					[]byte(terminalZeroChar),
					[]byte(newLineChar)),
				)
				break
			}

			if err = scanner.Err(); err != nil {
				result.error = err
			}
			if fileInfo.Size() == 0 {
				result.EnvValue.NeedRemove = true
			}

		} else {
			result.error = err
		}
	}
	return result
}
