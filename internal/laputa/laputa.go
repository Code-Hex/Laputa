package laputa

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Code-Hex/Laputa/internal/context"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/uber-go/zap"
)

type laputa struct {
	env    context.Env
	art    []byte
	Echo   *echo.Echo
	logger zap.Logger
	Secret string `json:"device_secret"`
}

type Response struct {
	Status string
}

func New(mode string) *laputa {
	laputa := &laputa{Echo: echo.New(), art: []byte(LAPUTA)}
	err := envconfig.Process("laputa", &laputa.env)
	if err != nil {
		log.Fatal(err.Error())
	}
	return laputa.setup(mode)
}

func (l *laputa) Run() int {
	l.SetMiddleware()
	l.RegisterRoute()
	// See serve.go
	if err := l.RunServer(); err != nil {
		l.logger.Error("Failed to run server", zap.String("reason", err.Error()))
		return 1
	}
	return 0
}

func (l *laputa) setup(mode string) *laputa {
	switch mode {
	case "develop":
		l.logger = context.Setlogger(os.Stderr)
	case "staging":
		logdir := os.Getenv("LOG_DIR")
		if logdir == "" {
			log.Fatal("LOG_DIR env was not set")
		}

		f, err := rotatelogs.New(
			filepath.Join(logdir, "laputa_log.%Y%m%d%H%M"),
			rotatelogs.WithLinkName(filepath.Join(logdir, "laputa_log")),
			rotatelogs.WithMaxAge(24*time.Hour),
			rotatelogs.WithRotationTime(time.Hour),
		)
		if err != nil {
			log.Fatalf("failed to create rotatelogs: %s", err)
		}
		defer f.Close()
		l.logger = context.Setlogger(zap.AddSync(f))
	default:
		log.Fatal("main.mode was not set")
	}

	return l
}

func (l *laputa) RegisterRoute() {
	l.Echo.Any("/", l.HealthCheck)
	l.Echo.HEAD("/information", l.GetInfo)
	l.Echo.POST("/register", l.Register)
}

func (laputa *laputa) HealthCheck(c echo.Context) error {
	art, err := laputa.Art()
	if err != nil {
		return c.String(http.StatusOK, "Good")
	}
	return c.String(http.StatusOK, art)
}

func (laputa *laputa) GetInfo(c echo.Context) error {
	laputa.logger.Info("Generate response header to get secret")
	c.Response().Header().Set("X-Device", laputa.GetDeviceHash())
	return c.String(http.StatusOK, "")
}

func (laputa *laputa) Register(c echo.Context) error {
	laputa.logger.Info("Registration processing...")
	header := c.Request().Header

	if subtle.ConstantTimeCompare([]byte(header.Get("X-Device")), []byte(laputa.GetDeviceHash())) != 1 {
		return c.JSON(http.StatusBadRequest, Response{Status: "Bad request"})
	}
	json.NewDecoder(c.Request().Body).Decode(laputa)

	return laputa.Store(c)
}

func (laputa *laputa) Store(c echo.Context) error {
	db, err := leveldb.OpenFile("secret", nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Status: err.Error()})
	}
	defer db.Close()
	if err := db.Put([]byte("secret"), []byte(laputa.Secret), nil); err != nil {
		return c.JSON(http.StatusInternalServerError, Response{Status: err.Error()})
	}

	return c.JSON(http.StatusOK, Response{})
}

func (laputa *laputa) GetDeviceHash() string {
	v := fmt.Sprintf("felica_device_%s", laputa.env.Floor)
	return fmt.Sprintf("%x", sha256.Sum256([]byte(v)))
}
