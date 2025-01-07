package auth

import (
	"net/http"
	"strings"

	"github.com/feelbeatapp/feelbeatserver/internal/infra/fblog"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/component"
	"github.com/feelbeatapp/feelbeatserver/internal/lib/feelbeaterror"
	"github.com/feelbeatapp/feelbeatserver/internal/thirdparty/spotify"
)

func AuthorizeThroughSpotify(handler func(User, http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")

		splits := strings.Split(authHeader, "Bearer ")
		if len(splits) != 2 {
			http.Error(res, "Incorrect authorization format", http.StatusBadRequest)
			fblog.Error(component.Auth, "Incorrect authorization format", "url", req.URL, "addr", req.RemoteAddr)
			return
		}
		token := splits[1]

		user, err := spotify.GetUserProfile(token)

		if err != nil {
			http.Error(res, feelbeaterror.AuthFailed, http.StatusForbidden)
			fblog.Error(component.Auth, "Access denied", "reason", err)
			return
		}

		handler(User{Profile: user, Token: token}, res, req)
	}
}
