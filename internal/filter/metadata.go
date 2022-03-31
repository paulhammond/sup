package filter

import (
	"github.com/bmatcuk/doublestar/v2"
	"github.com/paulhammond/sup/internal/cfg"
	"github.com/paulhammond/sup/internal/object"
)

func addMetadata(cfg cfg.Config, path string, o object.Object, debug DebugFunc) error {
	metadata, err := o.Metadata()
	if err != nil {
		return err
	}
	for _, rule := range cfg.Metadata {
		match, err := doublestar.Match(rule.Pattern, path)
		if err != nil {
			return err
		}
		if match {
			if rule.ContentType != nil {
				metadata.ContentType = rule.ContentType
				debug("metadata [%s] matches %q set ContentType %q", path, rule.Pattern, *rule.ContentType)
			}
		}
	}
	return nil
}
