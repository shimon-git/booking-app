package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

// Nosurf adds CSRF protection to all post requests
func Nosurf(next http.Handler) http.Handler {
	csfrHandler := nosurf.New(next)

	csfrHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csfrHandler
}

// SessionLoad loads and saves the sessions on every request.
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
