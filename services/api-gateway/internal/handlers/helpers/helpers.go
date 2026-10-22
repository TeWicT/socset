package helpers

import (
	socialgraphv1 "api-gateway/internal/gen/social-graph/v1"
	"api-gateway/internal/handlers/social-graph/socialdomain"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func IsMethodAndJSON(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		w.WriteHeader(405)
		fmt.Fprint(w, "Неправильный метод")
		return false
	}
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		w.WriteHeader(415)
		fmt.Fprint(w, "Unsupported Media Type")
		return false
	}
	return true
}
func IsMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		w.WriteHeader(405)
		fmt.Fprint(w, "Неправильный метод")
		return false
	}
	return true
}

func MapErrors(err error, w http.ResponseWriter) {
	st, ok := status.FromError(err)
	if ok && st.Code() == codes.InvalidArgument {
		w.WriteHeader(400)
		fmt.Fprint(w, st.Message())
		return
	}

	if ok && st.Code() == codes.Unauthenticated {
		w.WriteHeader(401)
		fmt.Fprint(w, st.Message())
		return
	}
	if ok && st.Code() == codes.PermissionDenied {
		w.WriteHeader(403)
		fmt.Fprint(w, st.Message())
		return
	}
	if ok && st.Code() == codes.NotFound {
		w.WriteHeader(404)
		fmt.Fprint(w, st.Message())
		return
	}

	if ok && st.Code() == codes.AlreadyExists {
		w.WriteHeader(409)
		fmt.Fprint(w, st.Message())
		return
	}
	log.Printf("error: %v", err)
	w.WriteHeader(500)
	fmt.Fprint(w, "internal error")
}

func WriteJSON(w http.ResponseWriter, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RelationStatusMapper(status socialgraphv1.RelationStatus) string {
	switch status {
	case socialgraphv1.RelationStatus_RELATION_STATUS_FRIENDS:
		return "friends"
	case socialgraphv1.RelationStatus_RELATION_STATUS_INCOMING_REQUEST:
		return "incoming_request"
	case socialgraphv1.RelationStatus_RELATION_STATUS_NONE:
		return "none"
	case socialgraphv1.RelationStatus_RELATION_STATUS_OUTGOING_REQUEST:
		return "outgoing_request"
	}
	return ""
}
func RequestStatusMapper(status socialgraphv1.RequestStatus) string {
	switch status {
	case socialgraphv1.RequestStatus_REQUEST_STATUS_ACCEPTED:
		return "accepted"
	case socialgraphv1.RequestStatus_REQUEST_STATUS_CANCELED:
		return "canceled"
	case socialgraphv1.RequestStatus_REQUEST_STATUS_PENDING:
		return "pending"
	case socialgraphv1.RequestStatus_REQUEST_STATUS_REJECTED:
		return "rejected"
	}
	return ""
}

func ToHTTPFriendRequest(request *socialgraphv1.FriendRequest) socialdomain.FriendRequest {
	return socialdomain.FriendRequest{ID: request.Id, FromUserID: request.FromUserId, ToUserID: request.ToUserId, Status: RequestStatusMapper(request.Status), CreatedAtUnix: request.CreatedAtUnix, UpdatedAtUnix: request.UpdatedAtUnix}
}
