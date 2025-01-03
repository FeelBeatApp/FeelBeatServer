package api

import (
	"net/http"

	"github.com/feelbeatapp/feelbeatserver/internal/infra/fblog"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/component"
)

func LogApiCall(userId string, req *http.Request) {
	fblog.Info(component.Api, req.Method+" "+req.URL.String(), "user", userId, "ip", req.RemoteAddr)
}

func LogApiError(message string, err error, userId string, req *http.Request) {
	fblog.Error(component.Api, req.Method+" "+req.URL.String()+": "+message, "err", err, "user", userId, "ip", req.RemoteAddr)
}
