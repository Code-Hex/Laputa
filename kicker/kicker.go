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
	req, err := http.NewRequest("GET", "http://127.0.0.1:3000/", nil)
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
	file := "/tmp/laputa.sock"
	defer os.Remove(file)
	listener, err := net.Listen("unix", file)
	if err != nil {
		return err
	}
	conn, err := listener.Accept()
	if err != nil {
		return err
	}

	sigchan := make(chan os.Signal)
	sigquit := make(chan bool)

	signal.Notify(sigchan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go listenSignal(sigchan, sigquit)
	go kicker.listenSocket(conn, sigquit)
	<-sigquit // wait for goroutines

	close(sigquit)
	close(sigchan)

	return nil
}

func (kicker *Kicker) listenSocket(conn net.Conn, quit chan bool) {
	defer conn.Close()
	for {
		select {
		case <-quit:
			kicker.logger.Info("killed with signal")
			return
		default:
			buf := make([]byte, 16)
			nr, err := conn.Read(buf)
			if err != nil {
				if err == io.EOF {
					quit <- true
					return
				}
				kicker.logger.Error(err.Error())
			}
			buf = buf[:nr]
			kicker.logger.Info(fmt.Sprintf("receive: %s", string(buf)))
			conn.Write([]byte("Success"))
		}
	}
}

func listenSignal(sigchan <-chan os.Signal, quit chan bool) {
	for {
		select {
		case <-quit:
			return
		default:
			switch <-sigchan {
			case syscall.SIGHUP:
				fallthrough
			case syscall.SIGINT:
				fallthrough
			case syscall.SIGTERM:
				fallthrough
			case syscall.SIGQUIT:
				quit <- true
			}
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
