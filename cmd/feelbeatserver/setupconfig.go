package main

import (
	"errors"
	"io/fs"
	"strings"

	"github.com/feelbeatapp/feelbeatserver/internal/component"
	"github.com/feelbeatapp/feelbeatserver/internal/fblog"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func setupConfig(config *koanf.Koanf) error {
	err := config.Load(file.Provider(TOML_PATH), toml.Parser())

	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}

		fblog.Info(component.Config, "config file not detected")
	} else {
		fblog.Info(component.Config, "config file loaded")
	}

	envProvider := env.Provider(ENV_PREFIX, ".", func(key string) string {
		envvar := strings.ToLower(strings.TrimPrefix(key, ENV_PREFIX))

		return strings.ReplaceAll(envvar, "_", ".")
	})

	return config.Load(envProvider, nil)
}
