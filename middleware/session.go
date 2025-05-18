package middleware

import (
	"net/http"
)

var MockGetUserRole func(*http.Request) (string, bool)

func GetUserRole(r *http.Request) (string, bool) {
	if MockGetUserRole != nil {
		return MockGetUserRole(r)
	}
	cookie, err := r.Cookie("session")
	if err != nil {
		return "", false
	}
	return cookie.Value, true
}

func RequireLogin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := GetUserRole(r)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role, ok := GetUserRole(r)
		if !ok || role != "admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}
