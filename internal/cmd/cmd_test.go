package cmd_test

import (
	"os"
	"testing"

	"github.com/paulhammond/sup/internal/cmd"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestIntegration2(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "../../tests",
	})
}

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"sup": cmd.Run,
	}))
}
