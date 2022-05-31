package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/spf13/pflag"
)

const (
	maxPortValue   = 65535
	minPortValue   = 1
	defaultTimeout = 10 * time.Second
)

type telnetFunc = func() error

func main() {
	address, timeout, err := parseParams()
	if err != nil {
		log.Fatal(err)
	}

	telnet := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)
	if err = telnet.Connect(); err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := telnet.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	ctx, cancelFunc := context.WithCancel(context.Background())

	go run(cancelFunc, telnet.Send)
	go run(cancelFunc, telnet.Receive)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	select {
	case <-sigCh:
	case <-ctx.Done():
		close(sigCh)
	}
}

func run(cancelFunc context.CancelFunc, tlFunc telnetFunc) {
	defer cancelFunc()

	err := tlFunc()
	if err != nil {
		_, err = fmt.Fprintln(os.Stderr, err)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func parseParams() (address string, timeout *time.Duration, err error) {
	timeout = pflag.Duration("timeout", defaultTimeout, "specify connection timeout, default: 10s")
	pflag.Parse()

	args := pflag.CommandLine.Args()
	if len(args) != 2 {
		log.Fatal("you must pass host & port")
	}

	host := args[0]
	port := args[1]

	if i, err := strconv.Atoi(port); err != nil || i < minPortValue || i > maxPortValue {
		log.Fatal("invalid port")
	}
	return net.JoinHostPort(host, port), timeout, nil
}
