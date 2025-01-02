package main

import (
	"github.com/feelbeatapp/feelbeatserver/internal/infra/api/handlers"
	"github.com/feelbeatapp/feelbeatserver/internal/infra/auth"
	"github.com/feelbeatapp/feelbeatserver/internal/infra/fblog"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/component"
)

type Authorizer interface {
	authorize()
}

func setupAPI(authWrapper auth.AuthWrapper) {
	handlers := []func(auth.AuthWrapper){handlers.ServeCreateGame}

	fblog.Info(component.Api, "Setting up REST API", "handlers count", len(handlers))

	for _, f := range handlers {
		f(authWrapper)
	}
}
