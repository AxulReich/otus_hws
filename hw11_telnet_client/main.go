package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"time"

	"github.com/spf13/pflag"
)

var reg *regexp.Regexp

type telnetFunc = func() error

func init() {
	reg = regexp.MustCompile(hostRegexp)
}

const (
	maxPortValue   = 65535
	minPortValue   = 1
	defaultTimeout = 10 * time.Second
	hostRegexp     = `\blocalhost\b|\b((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4}\b`
)

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
			log.Println(err)
		}
	}()

	ctx, cancelFunc := context.WithCancel(context.Background())

	go run(cancelFunc, telnet.Receive)
	go run(cancelFunc, telnet.Send)

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
			log.Println(err)
		}
	}
}

func parseParams() (address string, timeout *time.Duration, err error) {
	timeout = pflag.Duration("timeout", defaultTimeout, "specify connection timeout, default: 10s")
	pflag.Parse()

	args := pflag.CommandLine.Args()
	if len(args) != 2 {
		return "", nil, fmt.Errorf("you must pass host & port")
	}

	host := args[0]
	port := args[1]

	if ok := reg.MatchString(host); !ok {
		return "", nil, fmt.Errorf("invalid host")
	}

	if i, err := strconv.Atoi(port); err != nil || i < minPortValue || i > maxPortValue {
		return "", nil, fmt.Errorf("invalid port")
	}
	return net.JoinHostPort(host, port), timeout, nil
}
