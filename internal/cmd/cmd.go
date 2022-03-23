package cmd

import (
	"fmt"
	"os"

	"github.com/paulhammond/sup/internal/cfg"
	"github.com/paulhammond/sup/internal/object"
	_ "github.com/rogpeppe/go-internal/testscript"
	"github.com/spf13/pflag"
)

func Run() int {

	cmd := pflag.NewFlagSet("sup", pflag.ExitOnError)
	err := cmd.Parse(os.Args[1:])
	if err != nil {
		return printError(err)
	}

	args := cmd.Args()
	if len(args) != 1 {
		return printUsage()
	}

	cfg, err := cfg.Parse(args[0])
	if err != nil {
		return printError(err)
	}

	set, err := object.FS(os.DirFS(cfg.SourceClean()))
	if err != nil {
		return printError(err)
	}

	for _, path := range set.Paths() {
		fmt.Println(path)
	}

	return 0
}

func printError(err error) int {
	fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
	return 1
}

func printUsage() int {
	fmt.Fprintln(os.Stderr, "usage: sup <config>")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Arguments:")
	fmt.Fprintln(os.Stderr, "<config>    Config File")
	return 2
}
