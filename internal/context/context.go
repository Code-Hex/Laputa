package context

import (
	"time"

	"github.com/uber-go/zap"
)

type Env struct {
	Floor    string
	Akatsuki string
	Certfile string
	Keyfile  string
}

func JSTFormatter(key string) zap.TimeFormatter {
	return zap.TimeFormatter(func(t time.Time) zap.Field {
		jst := time.FixedZone("Asia/Tokyo", 9*3600)
		return zap.String(key, t.In(jst).Format(time.ANSIC))
	})
}

func Setlogger(Out zap.WriteSyncer) zap.Logger {
	return zap.New(
		zap.NewJSONEncoder(JSTFormatter("time")),
		zap.AddCaller(), // Add line number option
		zap.Output(Out),
	)
}
