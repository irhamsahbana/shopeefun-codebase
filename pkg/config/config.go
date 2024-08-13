package config

import (
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Opts struct {
		Config    any
		Paths     []string
		Filenames []string
	}
)

func Load(opts Opts) error {
	for _, p := range opts.Paths {
		fp := filepath.Join(p, ".env")
		// load env from file
		if _, fileErr := os.Stat(fp); fileErr == nil {
			// Set ENV for development
			_ = cleanenv.ReadConfig(fp, opts.Config)
		}
	}
	var err error
	for _, f := range opts.Filenames {
		for _, p := range opts.Paths {
			fp := filepath.Join(p, f)
			if _, fileErr := os.Stat(fp); fileErr != nil {
				return fileErr
			}
			if err = cleanenv.ReadConfig(fp, opts.Config); err != nil {
				return err
			}
		}
	}

	return err
}
