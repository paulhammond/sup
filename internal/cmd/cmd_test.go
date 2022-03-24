package cmd_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/paulhammond/sup/internal/cmd"
	"github.com/paulhammond/sup/internal/remote"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestIntegration2(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "../../tests",
	})
}

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"sup":  cmd.Run,
		"prep": prep,
	}))
}

func prep() int {
	err := remote.CreateFake(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return 1
	}
	return 0
}
