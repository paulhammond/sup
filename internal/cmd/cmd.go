package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	_ "github.com/paulhammond/licensepack/license"
	"github.com/spf13/pflag"

	"github.com/paulhammond/sup/internal/cfg"
	"github.com/paulhammond/sup/internal/filter"
	"github.com/paulhammond/sup/internal/object"
	"github.com/paulhammond/sup/internal/remote"
)

func Run() int {

	ctx := context.Background()

	cmd := pflag.NewFlagSet("sup", pflag.ExitOnError)
	var verbose *bool = cmd.BoolP("verbose", "v", false, "verbose output")
	var help *bool = cmd.BoolP("help", "h", false, "show help")
	var credits *bool = cmd.Bool("credits", false, "show open source credits")

	// this is set by the test scripts to enable bolt-db based fake remotes
	_, allowFakes := (os.LookupEnv("SUP_DEBUG_FAKE_REMOTE"))

	UI := &ui{Verbose: verbose}

	err := cmd.Parse(os.Args[1:])
	if err != nil {
		return UI.Error(err)
	}

	if *credits {
		return printCredits()
	}

	if *help {
		return printHelp()
	}

	args := cmd.Args()
	if len(args) != 2 {
		return printUsage()
	}

	cfg, err := cfg.Parse(args[0])
	if err != nil {
		return UI.Error(err)
	}

	UI.Start("Scanning local files:")
	set, err := object.FS(os.DirFS(cfg.SourceClean()), !cfg.IncludeDotfiles)
	if err != nil {
		return UI.Error(err)
	}

	for _, path := range set.Paths() {
		UI.Debug("· found " + path)
	}
	UI.Done("done")

	UI.Start("Applying filters:")
	err = filter.Filter(set, cfg, func(format string, a ...any) {
		UI.Debug(fmt.Sprintf("· "+format, a...))
	})
	if err != nil {
		return UI.Error(err)
	}
	UI.Done("done")

	UI.Start("Scanning remote files:")

	r, err := remote.Open(ctx, args[1], allowFakes)
	if err != nil {
		return UI.Error(err)
	}
	defer func() {
		err := r.Close()
		if err != nil {
			UI.Error(err)
		}
	}()

	remoteSet, err := r.Set(ctx)
	if err != nil {
		return UI.Error(err)
	}
	for _, path := range remoteSet.Paths() {
		UI.Debug("· found " + path)
	}
	UI.Done("done")

	UI.Start("Comparing:")
	toUpload, toDelete, err := remoteSet.Diff(set)
	if err != nil {
		return UI.Error(err)
	}
	UI.Done("done")

	if len(toUpload) > 0 {
		UI.Output("")
		UI.Output("These files will be uploaded:")
		for _, path := range toUpload.Paths() {
			UI.Output("· " + path)
		}
		y, err := UI.Prompt("Do you want to upload? (y to approve)")
		if err != nil {
			return UI.Error(err)
		}

		if strings.ToLower(strings.TrimSpace(y)) != "y" {
			UI.Output("OK, not uploading")
		} else {
			UI.Start("Uploading:")
			err = r.Upload(ctx, toUpload, func(e remote.Event) {
				UI.Output(fmt.Sprintf("· %s [%s]", e.Path, formatDuration(e.Duration)))
			})
			if err != nil {
				return UI.Error(err)
			}
			UI.Done("done")
		}
	}

	UI.Output("These files should be deleted:")
	for _, path := range toDelete.Paths() {
		UI.Output("· " + path)
	}

	return 0
}
