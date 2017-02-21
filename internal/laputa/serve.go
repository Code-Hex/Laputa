package laputa

import (
	"context"
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

	s := laputa.Echo.Server

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
