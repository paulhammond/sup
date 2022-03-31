package cfg

import (
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

type Config struct {
	path       string
	Source     string     `hcl:"source,optional"`
	Ignore     []string   `hcl:"ignore,optional"`
	Metadata   []Metadata `hcl:"metadata,block"`
	Redirects  bool       `hcl:"redirects,optional"`
	TrimSuffix []string   `hcl:"trim_suffix,optional"`
}

type Metadata struct {
	Pattern     string  `hcl:"type,label"`
	ContentType *string `hcl:"content_type"`
}

func (c Config) SourceClean() string {
	src := c.Source
	if src == "" {
		src = "."
	}

	cfgDir := filepath.Dir(c.path)
	return filepath.Clean(filepath.Join(cfgDir, src))
}

func Parse(path string) (Config, error) {
	var hclCfg Config
	err := hclsimple.DecodeFile(path, nil, &hclCfg)
	if err != nil {
		return hclCfg, err
	}
	hclCfg.path = path
	return hclCfg, nil
}
