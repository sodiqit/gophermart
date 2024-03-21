package auth

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type AuthController struct {
}

func (c *AuthController) Route() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/register", c.handleRegister)

	return r
}

func (c *AuthController) handleRegister(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	if contentType != "application/json" {
		http.Error(w, "need provide Content-Type: application/json", http.StatusBadRequest)
		return
	}

	w.Write([]byte(""))
}

func NewController() *AuthController {
	return &AuthController{}
}
