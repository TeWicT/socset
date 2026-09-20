package router

import (
	"api-gateway/internal/denylist"
	authv1 "api-gateway/internal/gen/auth/v1"
	profilev1 "api-gateway/internal/gen/profile/v1"
	"api-gateway/internal/handlers/auth"
	"api-gateway/internal/handlers/profile"
	"api-gateway/internal/middleware"
	"fmt"
	"net/http"
)

func NewRouter(authclient authv1.AuthServiceClient, profileclient profilev1.ProfileServiceClient, jwtSecret string, deny *denylist.DenyList) *http.ServeMux {
	mux := http.NewServeMux()
	authHandler := auth.NewHandler(authclient, deny)
	profileHandler := profile.NewHandler(profileclient)

	mux.HandleFunc("/{$}", homeHandler)
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/v1/auth/refresh", authHandler.Refresh)
	mux.Handle("GET /api/v1/profiles/{user_id}", middleware.JWTOptional(jwtSecret, deny)(http.HandlerFunc(profileHandler.GetProfile)))
	jwtMw := middleware.JWT(jwtSecret, deny)
	mux.Handle("GET /api/v1/auth/me", jwtMw(http.HandlerFunc(authHandler.Me)))
	mux.Handle("POST /api/v1/auth/logout", jwtMw(http.HandlerFunc(authHandler.Logout)))
	mux.Handle("GET /api/v1/profiles/me", jwtMw(http.HandlerFunc(profileHandler.GetProfileMe)))
	mux.Handle("PATCH /api/v1/profiles/me", jwtMw(http.HandlerFunc(profileHandler.UpdateProfile)))
	return mux
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Главная")
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "200")
}
