package router

import (
	authv1 "api-gateway/internal/gen/auth/v1"
	"api-gateway/internal/handlers/auth"
	"fmt"
	"net/http"
)

func NewRouter(authclient authv1.AuthServiceClient) *http.ServeMux {
	mux := http.NewServeMux()
	authHandler := auth.NewHandler(authclient)

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("/api/v1/auth/login", authHandler.Login)
	return mux
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Главная")
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "200")
}
