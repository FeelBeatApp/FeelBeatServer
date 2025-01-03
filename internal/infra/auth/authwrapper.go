package auth

import "net/http"

type AuthWrapper func(func(string, http.ResponseWriter, *http.Request)) func (http.ResponseWriter, *http.Request)
