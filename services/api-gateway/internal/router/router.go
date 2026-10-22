package router

import (
	"api-gateway/internal/denylist"
	authv1 "api-gateway/internal/gen/auth/v1"
	profilev1 "api-gateway/internal/gen/profile/v1"
	socialgraphv1 "api-gateway/internal/gen/social-graph/v1"
	"api-gateway/internal/handlers/auth"
	"api-gateway/internal/handlers/profile"
	socialgraph "api-gateway/internal/handlers/social-graph"
	"api-gateway/internal/middleware"
	"fmt"
	"net/http"
)

func NewRouter(authclient authv1.AuthServiceClient, profileclient profilev1.ProfileServiceClient, socialGraphClient socialgraphv1.SocialGraphServiceClient, jwtSecret string, deny *denylist.DenyList) *http.ServeMux {
	mux := http.NewServeMux()
	authHandler := auth.NewHandler(authclient, deny)
	profileHandler := profile.NewHandler(profileclient)
	socialGraphHandler := socialgraph.NewHandler(socialGraphClient)
	mux.HandleFunc("/{$}", homeHandler)
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/v1/auth/refresh", authHandler.Refresh)
	mux.Handle("GET /api/v1/profile/{user_id}", middleware.JWTOptional(jwtSecret, deny)(http.HandlerFunc(profileHandler.GetProfile)))
	jwtMw := middleware.JWT(jwtSecret, deny)
	mux.Handle("GET /api/v1/auth/me", jwtMw(http.HandlerFunc(authHandler.Me)))
	mux.Handle("POST /api/v1/auth/logout", jwtMw(http.HandlerFunc(authHandler.Logout)))
	mux.Handle("GET /api/v1/profile/me", jwtMw(http.HandlerFunc(profileHandler.GetProfileMe)))
	mux.Handle("PATCH /api/v1/profile/me", jwtMw(http.HandlerFunc(profileHandler.UpdateProfile)))
	mux.Handle("POST /api/v1/friends/requests", jwtMw(http.HandlerFunc(socialGraphHandler.CreateFriendRequest)))
	mux.Handle("POST /api/v1/friends/requests/{id}/accept", jwtMw(http.HandlerFunc(socialGraphHandler.AcceptFriendRequest)))
	mux.Handle("POST /api/v1/friends/requests/{id}/reject", jwtMw(http.HandlerFunc(socialGraphHandler.RejectFriendRequest)))
	mux.Handle("DELETE /api/v1/friends/requests/{id}", jwtMw(http.HandlerFunc(socialGraphHandler.CancelFriendRequest)))
	mux.Handle("GET /api/v1/friends", jwtMw(http.HandlerFunc(socialGraphHandler.ListFriendsHandler)))
	mux.Handle("GET /api/v1/friends/requests/incoming", jwtMw(http.HandlerFunc(socialGraphHandler.ListIncomingHandler)))
	mux.Handle("GET /api/v1/friends/requests/outgoing", jwtMw(http.HandlerFunc(socialGraphHandler.ListOutgoingHandler)))
	mux.Handle("GET /api/v1/friends/status/{user_id}", jwtMw(http.HandlerFunc(socialGraphHandler.GetStatusHandler)))
	mux.Handle("DELETE /api/v1/friends/{user_id}", jwtMw(http.HandlerFunc(socialGraphHandler.RemoveFriendHandler)))
	return mux
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Главная")
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "200")
}
