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

func (s *Server) CancelRequest(ctx context.Context, req *socialgraphv1.CancelRequestRequest) (*socialgraphv1.CancelRequestResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid userID")
	}
	reqID, err := uuid.Parse(req.RequestId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid reqID")
	}
	res, err := s.GraphService.CancelRequest(ctx, userID, reqID)
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
	return &socialgraphv1.CancelRequestResponse{Request: buildGRPCFriendRequest(res)}, nil
}

func (s *Server) RejectRequest(ctx context.Context, req *socialgraphv1.RejectRequestRequest) (*socialgraphv1.RejectRequestResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid userID")
	}
	reqID, err := uuid.Parse(req.RequestId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid reqID")
	}
	res, err := s.GraphService.RejectRequest(ctx, userID, reqID)
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
	return &socialgraphv1.RejectRequestResponse{Request: buildGRPCFriendRequest(res)}, nil
}

func (s *Server) RemoveFriend(ctx context.Context, req *socialgraphv1.RemoveFriendRequest) (*socialgraphv1.RemoveFriendResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid userID")
	}
	friendID, err := uuid.Parse(req.FriendUserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid friendID")
	}
	err = s.GraphService.RemoveFriendship(ctx, userID, friendID)
	if err != nil {
		if errors.Is(err, domain.ErrNoFriends) {
			return nil, status.Error(codes.NotFound, "no friend")
		}
		if errors.Is(err, domain.ErrSelfRequest) {
			return nil, status.Error(codes.InvalidArgument, "self request")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &socialgraphv1.RemoveFriendResponse{}, nil
}

func (s *Server) ListFriends(ctx context.Context, req *socialgraphv1.ListFriendsRequest) (*socialgraphv1.ListFriendsResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid userID")
	}
	friends, err := s.GraphService.ListFriends(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	protofriends := make([]*socialgraphv1.Friendship, len(friends))
	for idx, friend := range friends {
		protofriends[idx] = &socialgraphv1.Friendship{UserId: friend.UserID.String(), CreatedAtUnix: friend.CreatedAt.Unix()}
	}
	return &socialgraphv1.ListFriendsResponse{Friends: protofriends}, nil
}

func (s *Server) ListIncoming(ctx context.Context, req *socialgraphv1.ListIncomingRequest) (*socialgraphv1.ListIncomingResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid userID")
	}
	incoming, err := s.GraphService.ListIncoming(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	protoincoming := make([]*socialgraphv1.FriendRequest, len(incoming))
	for idx, people := range incoming {
		protoincoming[idx] = &socialgraphv1.FriendRequest{Id: people.ID.String(), FromUserId: people.FromUserID.String(), ToUserId: people.ToUserID.String(), Status: statusMapper(people.Status), CreatedAtUnix: people.CreatedAt.Unix(), UpdatedAtUnix: people.UpdatedAt.Unix()}
	}
	return &socialgraphv1.ListIncomingResponse{Requests: protoincoming}, nil
}

func (s *Server) ListOutgoing(ctx context.Context, req *socialgraphv1.ListOutgoingRequest) (*socialgraphv1.ListOutgoingResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid userID")
	}
	outgoing, err := s.GraphService.ListOutgoing(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	protooutgoing := make([]*socialgraphv1.FriendRequest, len(outgoing))
	for idx, people := range outgoing {
		protooutgoing[idx] = &socialgraphv1.FriendRequest{Id: people.ID.String(), FromUserId: people.FromUserID.String(), ToUserId: people.ToUserID.String(), Status: statusMapper(people.Status), CreatedAtUnix: people.CreatedAt.Unix(), UpdatedAtUnix: people.UpdatedAt.Unix()}
	}
	return &socialgraphv1.ListOutgoingResponse{Requests: protooutgoing}, nil
}

func (s *Server) GetStatus(ctx context.Context, req *socialgraphv1.GetStatusRequest) (*socialgraphv1.GetStatusResponse, error) {
	userID, err := uuid.Parse(req.ViewerUserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid ViewerUserID")
	}
	toUserID, err := uuid.Parse(req.TargetUserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid TargetUserID")
	}
	relation, err := s.GraphService.GetStatus(ctx, userID, toUserID)
	if err != nil {
		if errors.Is(err, domain.ErrSelfRequest) {
			return nil, status.Error(codes.InvalidArgument, "self request")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	switch relation.Status {
	case domain.RelationStatusFriends, domain.RelationStatusNone:
		return &socialgraphv1.GetStatusResponse{Status: relationStatusMapper(relation.Status)}, nil
	case domain.RelationStatusIncomingRequest, domain.RelationStatusOutgoingRequest:
		return &socialgraphv1.GetStatusResponse{Status: relationStatusMapper(relation.Status), RequestId: relation.ID.String()}, nil

	}
	return &socialgraphv1.GetStatusResponse{Status: socialgraphv1.RelationStatus_RELATION_STATUS_UNSPECIFIED}, nil
}
