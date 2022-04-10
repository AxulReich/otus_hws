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

	path, err := filepath.Abs(args[1])
	if err != nil {
		log.Fatalf("can't get absolute path err: %v, passed path: %v", err, args[1])
	}
	environment, err := ReadDir(path)
	if err != nil {
		log.Fatalf("can't get environment variables err: %v, passed path: %v", err, path)
	}
	// fmt.Printf("AXE env: %#v\n", environment)
	// fmt.Printf("AXE args: %#v\n", args[2:])
	res, errors := RunCmd(os.Stdin, os.Stdout, args[2:], environment)

	for _, err = range errors {
		fmt.Println(err)
	}

	os.Exit(res)
}
