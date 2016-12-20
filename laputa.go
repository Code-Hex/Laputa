package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/k0kubun/pp"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/uber-go/zap"
)

type Env struct {
	Mode     string
	Debug    bool
	Port     int
	Floor    string
	Akatsuki string
	Certfile string
	Keyfile  string
}

type Laputa struct {
	env    Env
	art    []byte
	balus  *Balus
	Secret string
}

type Response struct {
	Status string
}

func New() *Laputa {
	l := &Laputa{
		art:   []byte(LAPUTA),
		balus: BalusNew(os.Stdout),
	}

	err := envconfig.Process("laputa", &l.env)
	pp.Print(l.env)
	if err != nil {
		log.Fatal(err.Error())
	}

	go func() {
		err := l.balus.UnixSocket()
		if err != nil {
			l.Logger().Error(err.Error())
		}
	}()

	return l
}

func main() {
	laputa := New()
	e := echo.New()
	laputa.SetMiddleware(e)
	laputa.RegisterRoute(e)
	err := e.StartTLS(
		laputa.Port(),
		laputa.env.Certfile,
		laputa.env.Keyfile,
	)
	if err != nil {
		laputa.Logger().Error(err.Error())
	}
}

func (laputa Laputa) Logger() zap.Logger {
	return laputa.balus.logger
}

func (laputa Laputa) SetMiddleware(e *echo.Echo) {
	e.Use(middleware.Recover())
}

func (laputa *Laputa) RegisterRoute(e *echo.Echo) {
	e.GET("/", laputa.HealthCheck)
	e.GET("/information", laputa.HealthCheck) // fake
	e.HEAD("/information", laputa.GetInfo)
	e.POST("/register", laputa.Register)
}

func (laputa *Laputa) HealthCheck(c echo.Context) error {
	art, err := laputa.Art()
	if err != nil {
		return c.String(http.StatusOK, "Good")
	}
	return c.String(http.StatusOK, art)
}

func (laputa *Laputa) GetInfo(c echo.Context) error {
	laputa.Logger().Info("Generate response header to get secret")
	c.Response().Header().Set("X-Device", laputa.GetDeviceHash())
	return c.String(http.StatusOK, "")
}

func (laputa *Laputa) Register(c echo.Context) error {
	laputa.Logger().Info("Registration processing...")
	header := c.Request().Header
	if subtle.ConstantTimeCompare([]byte(header.Get("X-Device")), []byte(laputa.GetDeviceHash())) != 1 {
		return c.JSON(http.StatusBadRequest, Response{Status: "Bad request"})
	}
	decoder := json.NewDecoder(c.Request().Body)
	decoder.Decode(laputa)

	return laputa.Store(c)
}

func (laputa *Laputa) Store(c echo.Context) error {
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

func (laputa *Laputa) GetDeviceHash() string {
	v := fmt.Sprintf("edy_device_%s", laputa.env.Floor)
	return fmt.Sprintf("%x", sha256.Sum256([]byte(v)))
}

func (laputa *Laputa) Port() string {
	return fmt.Sprintf(":%d", laputa.env.Port)
}
