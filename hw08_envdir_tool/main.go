package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Printf("must pass path & at least bin, you pass: %v\n", args)
	}

	path, err := filepath.Abs(args[2])
	if err != nil {
		log.Fatalf("can't get absolute path err: %v, passed path: %v", err, args[2])
	}
	environment, err := ReadDir(path)
	if err != nil {
		log.Fatalf("can't get environment variables err: %v, passed path: %v", err, args[2])
	}

	os.Exit(RunCmd(args[1:], environment))
}
