package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	tcpClient := NewTelnetClient(dstAddr, timeout, os.Stdin, os.Stdout)
	err := tcpClient.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tcpClient.Close()

	go func() {
		select {
		case <-signalChan:
			cancel()
		case <-ctx.Done():
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := tcpClient.Receive(); err != nil {
					cancel()
				}
			}
		}
	}()

	go func() {
		if err := tcpClient.Send(); err != nil {
			cancel()
		} else {
			cancel()
			return
		}
	}()

	<-ctx.Done()
}
