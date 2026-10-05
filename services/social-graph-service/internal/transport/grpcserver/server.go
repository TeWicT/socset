package grpcserver

import (
	"context"
	"errors"
	"social-graph-service/internal/domain"
	socialgraphv1 "social-graph-service/internal/gen/social-graph/v1"
	"social-graph-service/internal/service"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	socialgraphv1.UnimplementedSocialGraphServiceServer
	GraphService *service.GraphService
}

func (s *Server) CreateRequest(ctx context.Context, req *socialgraphv1.CreateRequestRequest) (*socialgraphv1.CreateRequestResponse, error) {
	fromUserID, err := uuid.Parse(req.FromUserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid fromUserID")
	}
	toUserID, err := uuid.Parse(req.ToUserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid toUserID")
	}

	resp, err := s.GraphService.CreateRequest(ctx, fromUserID, toUserID)
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyFriends) {
			return nil, status.Error(codes.AlreadyExists, "already friends")
		}
		if errors.Is(err, domain.ErrPendingExist) {
			return nil, status.Error(codes.AlreadyExists, "pending exist")
		}
		if errors.Is(err, domain.ErrSelfRequest) {
			return nil, status.Error(codes.InvalidArgument, "self request")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &socialgraphv1.CreateRequestResponse{Request: buildGRPCFriendRequest(resp)}, nil
}

func (s *Server) AcceptRequest(ctx context.Context, req *socialgraphv1.AcceptRequestRequest) (*socialgraphv1.AcceptRequestResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid userID")
	}
	reqID, err := uuid.Parse(req.RequestId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid reqID")
	}
	res, err := s.GraphService.AcceptRequest(ctx, userID, reqID)
	if err != nil {
		if errors.Is(err, domain.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, "forbidden")
		}
		if errors.Is(err, domain.ErrNoPending) {
			return nil, status.Error(codes.NotFound, "not found pending")
		}
		if errors.Is(err, domain.ErrNoRequest) {
			return nil, status.Error(codes.NotFound, "not found request")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &socialgraphv1.AcceptRequestResponse{Request: buildGRPCFriendRequest(res)}, nil
}
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
