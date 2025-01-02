package handlers

import (
	"net/http"

	"github.com/feelbeatapp/feelbeatserver/internal/infra/api"
	"github.com/feelbeatapp/feelbeatserver/internal/infra/auth"
	"github.com/feelbeatapp/feelbeatserver/internal/infra/fblog"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/component"
)

type createGamePayload struct {
	Test string `json:"test"`
}

func createGameHandler(userId string, res http.ResponseWriter, req *http.Request) {
	var payload createGamePayload
	err := api.ParseBody(req.Body, &payload)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		fblog.Error(component.Api, "Error: ", "err", err)
		return
	}

	_, err = res.Write([]byte("Endpoint hit!, Auth success, You are: " + userId))
	if err != nil {
		api.LogApiError("Couldn't write response", err, userId, req)
		return
	}

	api.LogApiCall(userId, req)
}

func ServeCreateGame(authWrapper auth.AuthWrapper) {
	http.HandleFunc("/api/v1/create", authWrapper(createGameHandler))
}
