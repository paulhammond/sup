package main

import (
	"os"

	"github.com/paulhammond/sup/internal/cmd"
)

func main() {
	os.Exit(cmd.Run())
}
