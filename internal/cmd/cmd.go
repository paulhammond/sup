package cmd

import (
	"fmt"
	"os"

	"github.com/paulhammond/sup/internal/cfg"
	"github.com/paulhammond/sup/internal/filter"
	"github.com/paulhammond/sup/internal/object"
	"github.com/paulhammond/sup/internal/remote"
	_ "github.com/rogpeppe/go-internal/testscript"
	"github.com/spf13/pflag"
)

func Run() int {

	cmd := pflag.NewFlagSet("sup", pflag.ExitOnError)
	var verbose *bool = cmd.BoolP("verbose", "v", false, "verbose output")

	UI := &ui{Verbose: verbose}

	err := cmd.Parse(os.Args[1:])
	if err != nil {
		return UI.Error(err)
	}

	args := cmd.Args()
	if len(args) != 2 {
		return printUsage()
	}

	cfg, err := cfg.Parse(args[0])
	if err != nil {
		return UI.Error(err)
	}

	r, err := remote.Open(args[1])
	if err != nil {
		return UI.Error(err)
	}
	defer func() {
		err := r.Close()
		if err != nil {
			UI.Error(err)
		}
	}()

	set, err := object.FS(os.DirFS(cfg.SourceClean()))
	if err != nil {
		return UI.Error(err)
	}

	UI.Debug("local files:")
	for _, path := range set.Paths() {
		UI.Debug(path)
	}

	UI.Start("applying filters:")
	err = filter.Filter(&set, UI.DebugF)
	if err != nil {
		return UI.Error(err)
	}
	UI.Done("done")

	remoteSet, err := r.Set()
	if err != nil {
		return UI.Error(err)
	}
	UI.Debug("remote files:")
	for _, path := range remoteSet.Paths() {
		UI.Debug(path)
	}

	toUpload, toDelete, err := remoteSet.Diff(set)
	if err != nil {
		return UI.Error(err)
	}
	UI.Output("upload:")
	for _, path := range toUpload.Paths() {
		UI.Output(path)
	}
	UI.Output("delete:")
	for _, path := range toDelete.Paths() {
		UI.Output(path)
	}

	if len(toUpload) > 0 {
		UI.Start("uploading:")
		err = r.Upload(toUpload, func(e remote.Event) {
			UI.Output(fmt.Sprintf("%s [%s]", e.Path, formatDuration(e.Duration)))
		})
		if err != nil {
			return UI.Error(err)
		}
		UI.Done("done")
	}

	return 0
}
