package laputa

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lestrrat/go-server-starter/listener"
	"github.com/uber-go/zap"
)

func (laputa *laputa) RunServer() error {
	if laputa.hasNotSecurityFiles() {
		return errors.New("invalid tls configuration")
	}

	var l net.Listener

	port := os.Getenv("SERVER_STARTER_PORT")
	if port != "" {
		listeners, err := listener.ListenAll()
		if err != nil {
			return err
		}

		if len(listeners) > 0 {
			l = listeners[0]
		}
	}

	if l == nil {
		var err error
		port = ":8080"
		l, err = net.Listen("tcp", port)
		if err != nil {
			return err
		}
	}

	laputa.logger.Info(
		"Start laputa...",
		zap.String("Port", port),
		zap.String("Akatsuki", laputa.env.Akatsuki),
	)

	s := laputa.Echo.TLSServer
	s.TLSConfig = new(tls.Config)
	s.TLSConfig.Certificates = make([]tls.Certificate, 1)

	var err error
	s.TLSConfig.Certificates[0], err = tls.LoadX509KeyPair(laputa.pairKeyFiles())
	if err != nil {
		return err
	}
	if !laputa.Echo.DisableHTTP2 {
		s.TLSConfig.NextProtos = append(s.TLSConfig.NextProtos, "h2")
	}

	go func() {
		if err := serve(s, l); err != nil {
			laputa.logger.Error("serve error", zap.String("reason", err.Error()))
		}
	}()

	// Graceful shutdown (signal by TERM).
	termCh := make(chan os.Signal, 1)
	signal.Notify(termCh, syscall.SIGTERM)

	// Block until a term signal is coming
	<-termCh
	laputa.logger.Info(
		"Shutdown laputa...",
		zap.String("Port", port),
		zap.String("Akatsuki", laputa.env.Akatsuki),
	)

	timeout := time.Duration(10 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return shutdown(ctx, s)
}

func serve(server *http.Server, l net.Listener) error {
	return server.Serve(l)
}

func shutdown(ctx context.Context, server *http.Server) error {
	return server.Shutdown(ctx)
}

func (l *laputa) hasNotSecurityFiles() bool {
	return l.env.Certfile == "" || l.env.Keyfile == ""
}

func (l *laputa) pairKeyFiles() (string, string) {
	return l.env.Certfile, l.env.Keyfile
}
