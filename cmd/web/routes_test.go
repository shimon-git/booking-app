package main

import (
	//"net/http"
	"testing"

	"github.com/go-chi/chi"
	"github.com/shimon-git/booking-app/internal/config"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		// do nothing
	default:
		t.Errorf("type is not *chi.Mux, But is %T", v)
	}
}
