package auth

import (
	"net/http"

	"github.com/feelbeatapp/feelbeatserver/internal/lib"
)

type User struct {
	Profile lib.UserProfile
	Token   string
}

type AuthWrapper func(func(User, http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request)
