package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Client struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (c *Client) Connect() error {
	var err error
	c.conn, err = net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return fmt.Errorf("connect err: %w", err)
	}
	log.Printf("...Connected to %s", c.address)
	return nil
}

func (c *Client) Close() error {
	if c.conn == nil {
		return fmt.Errorf("connection is nil")
	}
	if c.in == nil {
		return fmt.Errorf("telnet in is nil")
	}

	if err := c.in.Close(); err != nil {
		return fmt.Errorf("in io.ReadCloser close err: %w", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("connection close err: %w", err)
	}
	return nil
}

func (c *Client) Send() error {
	if c.conn == nil {
		return fmt.Errorf("connection is nil")
	}
	if c.in == nil {
		return fmt.Errorf("telnet in is nil")
	}

	if _, err := io.Copy(c.conn, c.in); err != nil {
		return fmt.Errorf("send err: %w", err)
	}
	return nil
}

func (c *Client) Receive() error {
	if c.conn == nil {
		return fmt.Errorf("connection is nil")
	}
	if c.out == nil {
		return fmt.Errorf("telnet out is nil")
	}

	if _, err := io.Copy(c.out, c.conn); err != nil {
		return fmt.Errorf("receive err: %w", err)
	}
	return nil
}
