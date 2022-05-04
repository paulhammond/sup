package cmd

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"os"
	"time"
)

type ui struct {
	Verbose *bool
	started bool
}

func (u *ui) Output(s string) {
	if u.started {
		u.started = false
		fmt.Fprint(os.Stderr, "\n")
	}
	fmt.Fprint(os.Stderr, s+"\n")
}

func (u *ui) Debug(s string) {
	if !*u.Verbose {
		return
	}
	u.Output(s)
}

func (u *ui) Error(err error) int {
	if u.started {
		fmt.Fprint(os.Stderr, "\n")
	}
	fmt.Fprintf(os.Stderr, "error: %s\n", err)
	return 1
}

func (u *ui) Start(s string) {
	if u.started {
		fmt.Fprint(os.Stderr, "\n")
	}
	u.started = true
	fmt.Fprint(os.Stderr, s)
}

func (u *ui) Done(s string) {
	prefix := ""
	if u.started {
		prefix = " "
		u.started = false
	}
	fmt.Fprint(os.Stderr, prefix+s+"\n")
}

func (u *ui) Prompt(s string) (string, error) {
	u.Output(s)

	response, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err == io.EOF {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return response, nil
}

func printUsage() int {
	fmt.Fprintln(os.Stderr, "Usage: sup [options] <configfile> <remote>")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "For documentation run 'sup --help'")
	return 2
}

//go:embed help.txt
var helpTxt string

func printHelp() int {
	fmt.Println(helpTxt)
	return 0
}

func formatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return "~0ms"
	}
	return d.String()
}
