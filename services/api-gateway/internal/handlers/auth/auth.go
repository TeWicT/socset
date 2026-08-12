package auth

import (
	authv1 "api-gateway/internal/gen/auth/v1"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RegisterHTTPRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterHTTPResponse struct {
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type LoginHTTPRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginHTTPResponse struct {
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}
type Handler struct {
	auth authv1.AuthServiceClient
}

func NewHandler(auth authv1.AuthServiceClient) *Handler {
	return &Handler{auth: auth}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		fmt.Fprint(w, "Неправильный метод")
		return
	}
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		w.WriteHeader(415)
		fmt.Fprint(w, "Unsupported Media Type")
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
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.InvalidArgument {
			w.WriteHeader(400)
			fmt.Fprint(w, "Invalid argument")
			return
		}
		if ok && st.Code() == codes.AlreadyExists {
			w.WriteHeader(409)
			fmt.Fprint(w, "Пользователь с таким email или username уже существует")
			return
		}

		log.Printf("error: %v", err)
		w.WriteHeader(500)
		fmt.Fprint(w, "internal error")
		return
	}
	var registerResponse RegisterHTTPResponse
	registerResponse = RegisterHTTPResponse{UserID: res.UserId, AccessToken: res.AccessToken, RefreshToken: res.RefreshToken, ExpiresIn: res.ExpiresIn}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(registerResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

}
