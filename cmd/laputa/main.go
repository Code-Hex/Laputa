package main

import (
	"os"

	"github.com/Code-Hex/Laputa/internal/laputa"
)

var mode string

func main() {
	os.Exit(laputa.New(mode).Run())
}
