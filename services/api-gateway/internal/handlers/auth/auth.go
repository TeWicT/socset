package auth

import (
	authv1 "api-gateway/internal/gen/auth/v1"
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
