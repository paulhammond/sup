package cmd

import (
	"fmt"
	"os"

	"github.com/paulhammond/sup/internal/cfg"
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

	fmt.Println(cfg.SourceClean())
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
