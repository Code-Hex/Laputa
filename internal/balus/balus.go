package balus

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Code-Hex/Laputa/internal/context"
	"github.com/kelseyhightower/envconfig"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/uber-go/zap"
)

/*
 * Balus is a spell to collapse Laputa.
 * However, Since it is not possible to merge even if you create a destructive function,
 * here it means to confirm that felica's id is registered.
 */

type balus struct {
	env    context.Env
	logger zap.Logger
}

func New(mode string) *balus {
	balus := new(balus)
	err := envconfig.Process("laputa", &balus.env)
	if err != nil {
		log.Fatal(err.Error())
	}
	return balus.setup(mode)
}

func (b *balus) setup(mode string) *balus {
	switch mode {
	case "develop":
		b.logger = context.Setlogger(os.Stderr)
	case "staging":
		logdir := os.Getenv("LOG_DIR")
		if logdir == "" {
			log.Fatal("LOG_DIR env was not set")
		}

		f, err := rotatelogs.New(
			filepath.Join(logdir, "balus_log.%Y%m%d%H%M"),
			rotatelogs.WithLinkName(filepath.Join(logdir, "balus_log")),
			rotatelogs.WithMaxAge(24*time.Hour),
			rotatelogs.WithRotationTime(time.Hour),
		)
		if err != nil {
			log.Fatalf("failed to create rotatelogs: %s", err)
		}
		defer f.Close()
		b.logger = context.Setlogger(zap.AddSync(f))
	default:
		log.Fatal("main.mode was not set")
	}

	return b
}

func (b *balus) Run() int {
	if err := b.RunServer(); err != nil {
		b.logger.Error("Failed to run server", zap.String("reason", err.Error()))
		return 1
	}
	return 0
}

// What is your answer...
func (balus *balus) Balus(conn net.Conn, id string) {
	// Muska said,
	if balus.isRegisteredFelica(id) {
		conn.Write([]byte("I can't see!! I can't see!!!!"))
	} else {
		conn.Write([]byte("Get down on your knee. Beg your life."))
	}
}

func (balus *balus) isRegisteredFelica(id string) bool {
	balus.logger.Info("Checking felica", zap.String("id", id))

	url := os.Getenv("LAPUTA_AKATSUKI")
	if url == "" {
		balus.logger.Error("LAPUTA_AKATSUKI env is empty")
		return false
	}

	req, err := http.NewRequest("GET", url+id, nil)
	if err != nil {
		balus.logger.Error(err.Error())
		return false
	}

	secret, err := getDeviceSecret()
	if err != nil {
		balus.logger.Error(err.Error())
		return false
	}

	req.Header.Set("X-Secret", secret)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		balus.logger.Error(err.Error())
		return false
	}
	defer res.Body.Close()

	ldap := new(LDAP)
	json.NewDecoder(res.Body).Decode(ldap)

	return ldap.UID != ""
}

func getDeviceSecret() (string, error) {
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
