package main

import (
	"testing"

	"github.com/GekixD/Bookings/internal/config"
	"github.com/go-chi/chi"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig
	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		// do nothing, test passed
	default:
		t.Errorf("Type mismatch, expected '*chi.Mux' but instead got: %T", v)

	}
}
