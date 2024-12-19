package main

import (
	"errors"
	"io/fs"
	"strings"

	"github.com/feelbeatapp/feelbeatserver/internal/lib/component"
	"github.com/feelbeatapp/feelbeatserver/internal/infra/fblog"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func setupConfig(config *koanf.Koanf) error {
	err := config.Load(file.Provider(tomlPath), toml.Parser())

	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}

		fblog.Info(component.Config, "config file not detected")
	} else {
		fblog.Info(component.Config, "config file loaded")
	}

	envProvider := env.Provider(envPrefix, ".", func(key string) string {
		envvar := strings.ToLower(strings.TrimPrefix(key, envPrefix))

		return strings.ReplaceAll(envvar, "_", ".")
	})

	return config.Load(envProvider, nil)
}
