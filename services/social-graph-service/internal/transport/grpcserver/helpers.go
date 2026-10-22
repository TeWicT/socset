package grpcserver

import (
	"social-graph-service/internal/domain"
	socialgraphv1 "social-graph-service/internal/gen/social-graph/v1"
)

func buildGRPCFriendRequest(req *domain.FriendRequest) *socialgraphv1.FriendRequest {
	return &socialgraphv1.FriendRequest{Id: req.ID.String(), FromUserId: req.FromUserID.String(), ToUserId: req.ToUserID.String(), Status: statusMapper(req.Status), CreatedAtUnix: req.CreatedAt.Unix(), UpdatedAtUnix: req.UpdatedAt.Unix()}
}
func statusMapper(status string) socialgraphv1.RequestStatus {
	switch status {
	case domain.RequestStatusAccepted:
		return socialgraphv1.RequestStatus_REQUEST_STATUS_ACCEPTED
	case domain.RequestStatusRejected:
		return socialgraphv1.RequestStatus_REQUEST_STATUS_REJECTED
	case domain.RequestStatusPending:
		return socialgraphv1.RequestStatus_REQUEST_STATUS_PENDING
	case domain.RequestStatusCanceled:
		return socialgraphv1.RequestStatus_REQUEST_STATUS_CANCELED
	default:
		return socialgraphv1.RequestStatus_REQUEST_STATUS_UNSPECIFIED
	}

}

func relationStatusMapper(status string) socialgraphv1.RelationStatus {
	switch status {
	case domain.RelationStatusFriends:
		return socialgraphv1.RelationStatus_RELATION_STATUS_FRIENDS
	case domain.RelationStatusIncomingRequest:
		return socialgraphv1.RelationStatus_RELATION_STATUS_INCOMING_REQUEST
	case domain.RelationStatusOutgoingRequest:
		return socialgraphv1.RelationStatus_RELATION_STATUS_OUTGOING_REQUEST
	case domain.RelationStatusNone:
		return socialgraphv1.RelationStatus_RELATION_STATUS_NONE
	default:
		return socialgraphv1.RelationStatus_RELATION_STATUS_UNSPECIFIED
	}
}
