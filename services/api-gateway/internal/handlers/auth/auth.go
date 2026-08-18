package auth

import (
	authv1 "api-gateway/internal/gen/auth/v1"
	"api-gateway/internal/middleware"
	"encoding/json"
	"fmt"
	"net/http"
)

type RegisterHTTPRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthHTTPResponse struct {
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type LoginHTTPRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RefreshHTTPRequest struct {
	RefreshToken string `json:"refresh_token"`
}
type RefreshHTTPResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type LogoutHTTPRequest struct {
	RefreshToken string `json:"refresh_token"`
}
type Handler struct {
	auth authv1.AuthServiceClient
}

func NewHandler(auth authv1.AuthServiceClient) *Handler {
	return &Handler{auth: auth}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if !isPOSTAndJSON(w, r) {
		return
	}
	var registerRequest RegisterHTTPRequest
	err := json.NewDecoder(r.Body).Decode(&registerRequest)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "error json")
		return
	}
	res, err := h.auth.Register(r.Context(), &authv1.RegisterRequest{Email: registerRequest.Email, Password: registerRequest.Password, Username: registerRequest.Username})
	if err != nil {
		mapErrors(err, w)
		return
	}
	var registerResponse AuthHTTPResponse
	registerResponse = AuthHTTPResponse{UserID: res.UserId, AccessToken: res.AccessToken, RefreshToken: res.RefreshToken, ExpiresIn: res.ExpiresIn}
	writeJSON(w, registerResponse)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if !isPOSTAndJSON(w, r) {
		return
	}
	var loginRequest LoginHTTPRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "error json")
		return
	}

	res, err := h.auth.Login(r.Context(), &authv1.LoginRequest{Login: loginRequest.Login, Password: loginRequest.Password})
	if err != nil {

		mapErrors(err, w)
		return
	}
	var loginResponse AuthHTTPResponse
	loginResponse = AuthHTTPResponse{UserID: res.UserId, AccessToken: res.AccessToken, RefreshToken: res.RefreshToken, ExpiresIn: res.ExpiresIn}
	writeJSON(w, loginResponse)

}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	if !isPOSTAndJSON(w, r) {
		return
	}
	var refreshRequest RefreshHTTPRequest
	err := json.NewDecoder(r.Body).Decode(&refreshRequest)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "error json")
		return
	}
	res, err := h.auth.Refresh(r.Context(), &authv1.RefreshRequest{RefreshToken: refreshRequest.RefreshToken})
	if err != nil {
		mapErrors(err, w)
		return
	}
	var refreshResponse RefreshHTTPResponse
	refreshResponse = RefreshHTTPResponse{AccessToken: res.AccessToken, RefreshToken: res.RefreshToken, ExpiresIn: res.ExpiresIn}
	writeJSON(w, refreshResponse)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if !isPOSTAndJSON(w, r) {
		return
	}

	var logoutHTTPRequest LogoutHTTPRequest
	err := json.NewDecoder(r.Body).Decode(&logoutHTTPRequest)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "error json")
		return
	}
	_, err = h.auth.Logout(r.Context(), &authv1.LogoutRequest{RefreshToken: logoutHTTPRequest.RefreshToken})
	if err != nil {
		mapErrors(err, w)
		return
	}
	writeJSON(w, struct{}{})

}
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	id, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		w.WriteHeader(401)
		return
	}
	writeJSON(w, map[string]string{"user_id": id})
}
