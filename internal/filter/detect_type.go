package filter

import (
	"io"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/paulhammond/sup/internal/object"
)

func detectType(set object.Set, debug DebugFunc) error {
	for _, path := range set.Paths() {
		o := set[path]
		metadata, err := o.Metadata()

		if err != nil {
			return err
		}
		if metadata.ContentType == nil || *metadata.ContentType == "" || *metadata.ContentType == "application/octet-stream" {
			contentType := ""
			method := ""

			extension := filepath.Ext(path)
			if extension != "" {
				contentType = mime.TypeByExtension(extension)
				method = "extension"
			}

			if contentType == "" {
				buffer := make([]byte, 512)
				err := o.Open(func(r io.Reader) error {
					_, err = r.Read(buffer)
					return err
				})
				if err != nil {
					return err
				}

				contentType = http.DetectContentType(buffer)
				method = "contents"
			}

			metadata.ContentType = &contentType
			debug("detecttype [%s] detected %q via %s", path, contentType, method)
		}
	}
	return nil
}
