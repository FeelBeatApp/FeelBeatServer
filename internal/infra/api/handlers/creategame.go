package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/feelbeatapp/feelbeatserver/internal/infra/api"
	"github.com/feelbeatapp/feelbeatserver/internal/infra/auth"
	"github.com/feelbeatapp/feelbeatserver/internal/infra/fblog"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/component"
)

type createGamePayload struct {
	Test string `json:"test"`
}

type createGameResponse struct {
	RoomId string `json:"roomId"`
}

func createGameHandler(userId string, res http.ResponseWriter, req *http.Request) {
	var payload createGamePayload
	err := api.ParseBody(req.Body, &payload)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		fblog.Error(component.Api, "Error: ", "err", err)
		return
	}

	resJson, err := json.Marshal(createGameResponse{
		RoomId: "haha it's room id",
	})
	if err != nil {
		api.LogApiError("Couldn't encode response", err, userId, req)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	_, err = res.Write(resJson)
	if err != nil {
		api.LogApiError("Couldn't write response", err, userId, req)
		return
	}

	api.LogApiCall(userId, req)
}

func ServeCreateGame(baseUrl string, authWrapper auth.AuthWrapper) {
	http.HandleFunc(fmt.Sprintf("%s/create", baseUrl), authWrapper(createGameHandler))
}
