package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func (balus *Balus) UnixSocket() error {

	sigchan := make(chan os.Signal, 1)

	signal.Notify(sigchan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGABRT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)

	file := "/tmp/laputa.sock"
	defer os.Remove(file)
	ln, err := net.Listen("unix", file)
	if err != nil {
		return err
	}

	// listen to graceful shutdown
	go balus.listenSignal(ln, sigchan)

	// connect to nfc reader client
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		balus.logger.Info("unix socket is listning...")
		go balus.listenSocket(conn)
	}

	close(sigchan)

	return nil
}

func (balus *Balus) listenSocket(conn net.Conn) {
	for {
		// create buffer for felica id(max digit is 16)
		buf := make([]byte, 16)
		nr, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				balus.logger.Info("shutdown unix socket...")
				return
			}
			balus.logger.Error(err.Error())
		}
		buf = buf[:nr]
		balus.logger.Info(fmt.Sprintf("receive: %s", string(buf)))
		balus.Balus(conn, string(buf))
	}
}

func (balus *Balus) listenSignal(conn net.Listener, sigchan <-chan os.Signal) {
	for {
		sig := <-sigchan
		switch sig {
		case syscall.SIGHUP:
			fallthrough
		case syscall.SIGINT:
			fallthrough
		case syscall.SIGQUIT:
			fallthrough
		case syscall.SIGABRT:
			fallthrough
		case syscall.SIGKILL:
			fallthrough
		case syscall.SIGTERM:
			balus.logger.Info(fmt.Sprintf("Caught signal %s: shutting down.", sig))
			conn.Close()
			return
		}
	}
}
