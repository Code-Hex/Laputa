package laputa

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/uber-go/zap"
)

func (l laputa) SetMiddleware() {
	l.Echo.Use(middleware.Recover())
	l.Echo.Use(l.LogHandler())
}

func (l *laputa) LogHandler() echo.MiddlewareFunc {
	return func(before echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := before(c)
			stop := time.Now()

			w, r := c.Response(), c.Request()
			l.logger.Info(
				"Detected access",
				zap.String("status", fmt.Sprintf("%d: %s", w.Status, http.StatusText(w.Status))),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("useragent", r.UserAgent()),
				zap.String("remote_ip", r.RemoteAddr),
				zap.Int64("latency", stop.Sub(start).Nanoseconds()/int64(time.Microsecond)),
			)
			return err
		}
	}
}

func JSTFormatter(key string) zap.TimeFormatter {
	return zap.TimeFormatter(func(t time.Time) zap.Field {
		jst := time.FixedZone("Asia/Tokyo", 9*3600)
		return zap.String(key, t.In(jst).Format(time.ANSIC))
	})
}
