package laputa

import (
	"crypto/sha256"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/Code-Hex/Laputa/internal/context"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	envs := []context.Env{
		{
			Floor:    "F321",
			Certfile: "/hello/cert",
			Keyfile:  "/test/key",
			Akatsuki: "localhost:8080",
		},
		{
			Floor:    "F322",
			Certfile: "cert",
			Keyfile:  "key",
			Akatsuki: "localhost",
		},
	}

	data := []struct {
		Mode string
		Envs []context.Env
	}{
		{
			Mode: "develop",
			Envs: envs,
		},
		{
			Mode: "staging",
			Envs: envs,
		},
	}

	for _, d := range data {
		for _, env := range d.Envs {
			os.Setenv("LAPUTA_FLOOR", env.Floor)
			os.Setenv("LAPUTA_CERTFILE", env.Certfile)
			os.Setenv("LAPUTA_KEYFILE", env.Keyfile)
			os.Setenv("LAPUTA_AKATSUKI", env.Akatsuki)
			os.Setenv("LOG_DIR", "_testdata")

			l := New(d.Mode)
			if !reflect.DeepEqual(l.env, env) {
				t.Errorf("Unexpected environment struct")
			}

			assert.Equal(t, l.GetDeviceHash(), makeDeviceHash(env.Floor))

			if _, err := os.Stat("_testdata/laputa.log"); err == nil {
				os.Remove("_testdata/laputa.log")
			}
		}
	}

}

func makeDeviceHash(floor string) string {
	v := fmt.Sprintf("felica_device_%s", floor)
	return fmt.Sprintf("%x", sha256.Sum256([]byte(v)))
}
