package filter

import (
	"bufio"
	"io"
	"strings"

	"github.com/paulhammond/sup/internal/cfg"
	"github.com/paulhammond/sup/internal/object"
)

var textPlainType = "text/plain; charset=utf-8"

func processRedirect(config cfg.Config, set object.Set, debug DebugFunc) error {
	if !config.Redirects {
		return nil
	}

	for _, path := range set.Paths() {
		if strings.HasSuffix(path, ".redirect") {
			o := set[path]
			err := o.Open(func(r io.Reader) error {
				br := bufio.NewReader(r)
				redirect, err := br.ReadString('\n')
				if err != nil && err != io.EOF {
					return err
				}

				redirect = strings.TrimSuffix(redirect, "\n")

				delete(set, path)

				blob := object.Empty(object.Metadata{
					ContentType:             &textPlainType,
					WebsiteRedirectLocation: &redirect,
				})
				newPath := strings.TrimSuffix(path, ".redirect")
				set[newPath] = blob
				debug("redirect [%s] created redirect to %q", newPath, redirect)
				return nil
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
