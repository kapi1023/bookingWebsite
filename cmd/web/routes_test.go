package main

import (
	"fmt"
	"testing"

	"github.com/go-chi/chi"
	"github.com/kapi1023/bookingWebsite/internal/config"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig
	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		//do nothink
	default:
		t.Error(fmt.Sprintf("type is not chi.mux, type is %T", v))
	}
}
