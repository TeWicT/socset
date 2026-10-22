package socialgraph

import (
	socialgraphv1 "api-gateway/internal/gen/social-graph/v1"
	"api-gateway/internal/handlers/helpers"
	"api-gateway/internal/handlers/social-graph/socialdomain"
	"api-gateway/internal/middleware"
	"encoding/json"
	"fmt"
	"net/http"
)

type Handler struct {
	socialgraphv1.SocialGraphServiceClient
}

func NewHandler(socialgraph socialgraphv1.SocialGraphServiceClient) *Handler {
	return &Handler{SocialGraphServiceClient: socialgraph}
}

func (h *Handler) CreateFriendRequest(w http.ResponseWriter, r *http.Request) {

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		w.WriteHeader(401)
		fmt.Fprint(w, "unauthenticated")
		return
	}
	if !helpers.IsMethodAndJSON(w, r, "POST") {
		return
	}
	var req socialdomain.CreateRequestHTTPRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "error json")
		return
	}
	toUserID := req.ToUserID
	if toUserID == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "empty to_user_id")
		return
	}
	res, err := h.CreateRequest(r.Context(), &socialgraphv1.CreateRequestRequest{FromUserId: userID, ToUserId: toUserID})
	if err != nil {
		helpers.MapErrors(err, w)
		return
	}

	helpers.WriteJSON(w, socialdomain.CreateRequestHTTPResponse{FriendRequest: helpers.ToHTTPFriendRequest(res.Request)})
}

func (h *Handler) AcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		w.WriteHeader(401)
		fmt.Fprint(w, "unauthenticated")
		return
	}
	if !helpers.IsMethod(w, r, "POST") {
		return
	}

	reqID := r.PathValue("id")
	if reqID == "" {
		w.WriteHeader(400)
		fmt.Fprint(w, "empty id")
		return
	}
	res, err := h.AcceptRequest(r.Context(), &socialgraphv1.AcceptRequestRequest{RequestId: reqID, UserId: userID})
	if err != nil {
		helpers.MapErrors(err, w)
		return
	}
	helpers.WriteJSON(w, socialdomain.AcceptRequestHTTPResponse{FriendRequest: helpers.ToHTTPFriendRequest(res.Request)})
}

func (h *Handler) RejectFriendRequest(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		w.WriteHeader(401)
		fmt.Fprint(w, "unauthenticated")
		return
	}
	if !helpers.IsMethod(w, r, "POST") {
		return
	}
	reqID := r.PathValue("id")
	if reqID == "" {
		w.WriteHeader(400)
		fmt.Fprint(w, "empty id")
		return
	}
	res, err := h.RejectRequest(r.Context(), &socialgraphv1.RejectRequestRequest{RequestId: reqID, UserId: userID})
	if err != nil {
		helpers.MapErrors(err, w)
		return
	}
	helpers.WriteJSON(w, socialdomain.RejectRequestHTTPResponse{FriendRequest: helpers.ToHTTPFriendRequest(res.Request)})
}
func (h *Handler) CancelFriendRequest(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		w.WriteHeader(401)
		fmt.Fprint(w, "unauthenticated")
		return
	}
	if !helpers.IsMethod(w, r, "DELETE") {
		return
	}
	reqID := r.PathValue("id")
	if reqID == "" {
		w.WriteHeader(400)
		fmt.Fprint(w, "empty id")
		return
	}
	res, err := h.CancelRequest(r.Context(), &socialgraphv1.CancelRequestRequest{RequestId: reqID, UserId: userID})
	if err != nil {
		helpers.MapErrors(err, w)
		return
	}
	helpers.WriteJSON(w, socialdomain.CancelRequestHTTPResponse{FriendRequest: helpers.ToHTTPFriendRequest(res.Request)})
}
