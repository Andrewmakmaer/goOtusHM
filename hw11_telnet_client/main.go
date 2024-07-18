package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

var (
	host, port string
	timeout    time.Duration
)

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout for connect")
}

func main() {
	host, port = os.Args[len(os.Args)-2], os.Args[len(os.Args)-1]
	dstAddr := net.JoinHostPort(host, port)

	tcpClient := NewTelnetClient(dstAddr, timeout, os.Stdin, os.Stdout)
	err := tcpClient.Connect()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	go func() {
		for {
			tcpClient.Receive()
		}
	}()

	tcpClient.Send()
	tcpClient.Close()
}
