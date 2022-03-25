package filter

import (
	"io"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/paulhammond/sup/internal/object"
)

func detectType(path string, o object.Object) error {
	metadata, err := o.Metadata()
	if err != nil {
		return err
	}
	if metadata.ContentType == nil || *metadata.ContentType == "" || *metadata.ContentType == "application/octet-stream" {
		contentType := ""

		extension := filepath.Ext(path)
		if extension != "" {
			contentType = mime.TypeByExtension(extension)
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
		}

		metadata.ContentType = &contentType
	}
	return nil
}
