package balus

import (
	"fmt"
	"io"
	"net"
	"os"
	"syscall"

	"github.com/lestrrat/go-server-starter/listener"
	"github.com/uber-go/zap"
)

func (balus *balus) RunServer() error {
	var (
		path string
		l    net.Listener
	)

	listeners, err := listener.ListenAll()
	if err != nil {
		return err
	}

	if len(listeners) > 0 {
		l = listeners[0]
		if v, ok := l.(*net.UnixListener); ok {
			path = v.Addr().String()
		}
	}

	if l == nil {
		var err error
		path = "/tmp/laputa.sock"
		l, err = net.Listen("unix", path)
		if err != nil {
			return err
		}
	}

	balus.logger.Info(
		"Graceful start balus...",
		zap.String("Path", path),
		zap.String("Akatsuki", balus.env.Akatsuki),
	)

	return balus.serve(l)
}

func (balus *balus) serve(l net.Listener) error {
	// connect to nfc reader client
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		balus.logger.Info("unix socket is listning...")
		go balus.listenSocket(conn)
	}
}

func (balus *balus) listenSocket(conn net.Conn) {
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

func (balus *balus) listenSignal(conn net.Listener, sigchan <-chan os.Signal) {
	for sig := range sigchan {
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
