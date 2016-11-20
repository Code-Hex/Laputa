package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	fmt.Println()
}

func (laputa *Laputa) IsRegisteredEdy(number int) bool {

	c.Logger().Infof("Checking edy number: %d", number)
	req, err := http.NewRequest("GET", "http://127.0.0.1:3000/", nil)
	if err != nil {
		c.Logger().Errorf(err.Error())
		return false
	}

	if err != nil {
		c.Logger().Errorf(err.Error())
		return false
	}

	secret, err := laputa.DeviceSecret()
	if err != nil {
		c.Logger().Errorf(err.Error())
		return false
	}

	req.Header.Set("X-Secret", secret)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		c.Logger().Errorf(err.Error())
		return false
	}
	defer res.Body.Close()

	ldap := new(LDAP)
	json.NewDecoder(res.Body).Decode(ldap)

	return ldap != nil
}

func (laputa *Laputa) DeviceSecret(c echo.Context) (string, error) {
	db, err := leveldb.OpenFile("secret", nil)
	if err != nil {
		return "", err
	}
	defer db.Close()

	secret, err := db.Get([]byte("secret", nil))
	if err != nil {
		return "", err
	}

	return string(secret), nil
}
