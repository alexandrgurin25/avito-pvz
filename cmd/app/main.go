package main

import (
	"avito-pvz/internal/transport/http/handlers/auth"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	handler := auth.NewHandler()

	r := chi.NewRouter()
	r.Post("/dummyLogin", handler.DummyLogin)
	http.ListenAndServe(":8080", r)
}
