package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/uber-go/zap"
)

type Kicker struct {
	logger zap.Logger
}

func main() {
	kicker := New(os.Stdout)
	err := kicker.UnixSocket()
	if err != nil {
		kicker.logger.Error(err.Error())
	}
}

func New(Out zap.WriteSyncer) *Kicker {
	return &Kicker{
		logger: zap.New(
			zap.NewTextEncoder(zap.TextTimeFormat(time.ANSIC)),
			zap.AddCaller(), // Add line number option
			zap.Output(Out),
		),
	}
}

func (kicker *Kicker) IsRegisteredEdy(number int) bool {
	kicker.logger.Info(fmt.Sprintf("Checking edy number: %d", number))
	req, err := http.NewRequest("GET", os.Getenv("LAPUTA_TARGET"), nil)
	if err != nil {
		kicker.logger.Error(err.Error())
		return false
	}

	if err != nil {
		kicker.logger.Error(err.Error())
		return false
	}

	secret, err := GetDeviceSecret()
	if err != nil {
		kicker.logger.Error(err.Error())
		return false
	}

	req.Header.Set("X-Secret", secret)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		kicker.logger.Error(err.Error())
		return false
	}
	defer res.Body.Close()

	ldap := new(LDAP)
	json.NewDecoder(res.Body).Decode(ldap)

	return ldap != nil
}

func (kicker *Kicker) UnixSocket() error {

	sigchan := make(chan os.Signal, 1)

	signal.Notify(sigchan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	file := "/tmp/laputa.sock"
	defer os.Remove(file)
	ln, err := net.Listen("unix", file)
	if err != nil {
		return err
	}

	go kicker.listenSignal(ln, sigchan)

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		kicker.logger.Info("unix socket is listning...")
		go kicker.listenSocket(conn)
	}

	close(sigchan)

	return nil
}

func (kicker *Kicker) listenSocket(conn net.Conn) {
	for {
		// create buffer for felica id(max digit is 16)
		buf := make([]byte, 16)
		nr, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				kicker.logger.Info("shutdown unix socket...")
				return
			}
			kicker.logger.Error(err.Error())
		}
		buf = buf[:nr]
		kicker.logger.Info(fmt.Sprintf("receive: %s", string(buf)))

		conn.Write([]byte("Success"))
	}
}

func (kicker *Kicker) listenSignal(conn net.Listener, sigchan <-chan os.Signal) {
	for {
		sig := <-sigchan
		switch sig {
		case syscall.SIGHUP:
			fallthrough
		case syscall.SIGINT:
			fallthrough
		case Tsyscall.SIGQUIT:
			fallthrough
		case syscall.SIGABR:
			fallthrough
		case syscall.SIGKILL:
			fallthrough
		case syscall.SIGTERM:
			kicker.logger.Info(fmt.Sprintf("Caught signal %s: shutting down.", sig))
			conn.Close()
			return
		}
	}
}

func GetDeviceSecret() (string, error) {
	db, err := leveldb.OpenFile("secret", nil)
	if err != nil {
		return "", err
	}
	defer db.Close()

	secret, err := db.Get([]byte("secret"), nil)
	if err != nil {
		return "", err
	}

	return string(secret), nil
}
