package cmd

import (
	"fmt"
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

func printUsage() int {
	fmt.Fprintln(os.Stderr, "usage: sup <config>")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Arguments:")
	fmt.Fprintln(os.Stderr, "<config>    Config File")
	return 2
}

func formatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return "~0ms"
	}
	return d.String()
}
