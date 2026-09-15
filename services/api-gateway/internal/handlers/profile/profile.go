package profile

import (
	profilev1 "api-gateway/internal/gen/profile/v1"
	"api-gateway/internal/handlers/helpers"
	"api-gateway/internal/middleware"
	"encoding/json"
	"fmt"
	"net/http"
)

type Handler struct {
	profile profilev1.ProfileServiceClient
}

func NewHandler(profile profilev1.ProfileServiceClient) *Handler {
	return &Handler{profile: profile}
}

type GetProfileHTTPResponse struct {
	UserID      string `json:"user_id"`
	DisplayName string `json:"display_name"`
	Bio         string `json:"bio"`
	AvatarKey   string `json:"avatar_key"`
	IsPrivate   bool   `json:"is_private"`
}

type UpdateProfileHTTPRequest struct {
	DisplayName *string `json:"display_name"`
	Bio         *string `json:"bio"`
	AvatarKey   *string `json:"avatar_key"`
	IsPrivate   *bool   `json:"is_private"`
}
type UpdateProfileHTTPResponse struct {
	UserID      string `json:"user_id"`
	DisplayName string `json:"display_name"`
	Bio         string `json:"bio"`
	AvatarKey   string `json:"avatar_key"`
	IsPrivate   bool   `json:"is_private"`
}

func (h *Handler) GetProfileMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		w.WriteHeader(401)
		fmt.Fprint(w, "unauthenticated")
		return
	}
	res, err := h.profile.GetProfile(r.Context(), &profilev1.GetProfileRequest{UserId: userID})
	if err != nil {
		helpers.MapErrors(err, w)
		return
	}
	var getProfileResponse GetProfileHTTPResponse
	getProfileResponse = GetProfileHTTPResponse{UserID: res.UserId, DisplayName: res.DisplayName, Bio: res.Bio, AvatarKey: res.AvatarKey, IsPrivate: res.IsPrivate}
	helpers.WriteJSON(w, getProfileResponse)
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	clientUserID, _ := middleware.UserIDFromContext(r.Context())

	userID := r.PathValue("user_id")
	res, err := h.profile.GetProfile(r.Context(), &profilev1.GetProfileRequest{UserId: userID})
	if err != nil {
		helpers.MapErrors(err, w)
		return
	}
	var getProfileResponse GetProfileHTTPResponse
	getProfileResponse = GetProfileHTTPResponse{UserID: res.UserId, DisplayName: res.DisplayName, Bio: res.Bio, AvatarKey: res.AvatarKey, IsPrivate: res.IsPrivate}
	if getProfileResponse.IsPrivate == true && getProfileResponse.UserID != clientUserID {
		getProfileResponse.Bio = ""
	}
	helpers.WriteJSON(w, getProfileResponse)
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		w.WriteHeader(401)
		fmt.Fprint(w, "unauthenticated")
		return
	}

	if !helpers.IsMethodAndJSON(w, r, "PATCH") {
		return
	}
	var req UpdateProfileHTTPRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "error json")
		return
	}

	res, err := h.profile.UpdateProfile(r.Context(), &profilev1.UpdateProfileRequest{UserId: userID, DisplayName: req.DisplayName, Bio: req.Bio, AvatarKey: req.AvatarKey, IsPrivate: req.IsPrivate})
	if err != nil {
		helpers.MapErrors(err, w)
		return
	}
	var UpdateProfileResponse UpdateProfileHTTPResponse
	UpdateProfileResponse = UpdateProfileHTTPResponse{UserID: res.UserId, DisplayName: res.DisplayName, Bio: res.Bio, IsPrivate: res.IsPrivate, AvatarKey: res.AvatarKey}
	helpers.WriteJSON(w, UpdateProfileResponse)
}
