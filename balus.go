package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/uber-go/zap"
)

/*
 * Balus is a spell where Laputa collapses.
 * However, Since it is not possible to merge even if you create a destructive function,
 * here it means to confirm that felica's id is registered.
 */

type Balus struct {
	logger zap.Logger
}

func BalusNew(Out zap.WriteSyncer) *Balus {
	return &Balus{
		logger: zap.New(
			zap.NewTextEncoder(zap.TextTimeFormat(time.ANSIC)),
			zap.AddCaller(), // Add line number option
			zap.Output(Out),
		),
	}
}

// What is your answer...
func (balus *Balus) Balus(conn net.Conn, id string) {
	// Muska said,
	if balus.isRegisteredFelica(id) {
		conn.Write([]byte("I can't see!! I can't see!!!!"))
	} else {
		conn.Write([]byte("Get down on your knee. Beg your life."))
	}
}

func (balus *Balus) isRegisteredFelica(id string) bool {
	balus.logger.Info(fmt.Sprintf("Checking felica id: %s", id))

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
