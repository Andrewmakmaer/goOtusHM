package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type client struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (c *client) Connect() error {
	dialer := &net.Dialer{Timeout: c.timeout}
	conn, err := dialer.Dial("tcp", c.address)
	if err != nil {
		return fmt.Errorf("connection error: %w", err)
	}
	c.conn = conn
	return nil
}

func (c *client) Send() error {
	scanner := bufio.NewScanner(c.in)
	for scanner.Scan() {
		_, err := c.conn.Write(append(scanner.Bytes(), '\n'))
		if err != nil {
			return fmt.Errorf("send error: %w", err)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return fmt.Errorf("input error: %w", err)
	}
	return nil
}

func (c *client) Receive() error {
	_, err := io.Copy(c.out, c.conn)
	if err != nil {
		return fmt.Errorf("receive error: %w", err)
	}
	return nil
}

func (c *client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}
