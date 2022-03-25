package cmd

import (
	"fmt"
	"os"

	"github.com/paulhammond/sup/internal/cfg"
	"github.com/paulhammond/sup/internal/object"
	"github.com/paulhammond/sup/internal/remote"
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
	if len(args) != 2 {
		return printUsage()
	}

	cfg, err := cfg.Parse(args[0])
	if err != nil {
		return printError(err)
	}

	r, err := remote.Open(args[1])
	if err != nil {
		return printError(err)
	}
	defer func() {
		err := r.Close()
		if err != nil {
			printError(err)
		}
	}()

	set, err := object.FS(os.DirFS(cfg.SourceClean()))
	if err != nil {
		return printError(err)
	}

	fmt.Println("local files:")
	for _, path := range set.Paths() {
		fmt.Println(path)
	}

	remoteSet, err := r.Set()
	if err != nil {
		return printError(err)
	}
	fmt.Println("remote files:")
	for _, path := range remoteSet.Paths() {
		fmt.Println(path)
	}

	toUpload, toDelete, err := remoteSet.Diff(set)
	if err != nil {
		return printError(err)
	}
	fmt.Println("upload:")
	for _, path := range toUpload.Paths() {
		fmt.Println(path)
	}
	fmt.Println("delete:")
	for _, path := range toDelete.Paths() {
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
