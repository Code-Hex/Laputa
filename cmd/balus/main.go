package main

import (
	"os"

	"github.com/Code-Hex/Laputa/internal/balus"
)

var mode string

func main() {
	os.Exit(balus.New(mode).Run())
}
