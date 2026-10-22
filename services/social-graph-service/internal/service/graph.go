package service

import (
	"context"
	"social-graph-service/internal/domain"
	"social-graph-service/internal/events"
	"social-graph-service/internal/repository"

	"github.com/google/uuid"
)

type GraphService struct {
	graphRepo repository.GraphRepo
}

func NewGraphService(graphRepo repository.GraphRepo) *GraphService {
	return &GraphService{graphRepo: graphRepo}
}

func (svc *GraphService) CreateRequest(ctx context.Context, fromUserID, toUserID uuid.UUID) (*domain.FriendRequest, error) {
	if fromUserID == toUserID {
		return nil, domain.ErrSelfRequest
	}
	req := domain.FriendRequest{ID: uuid.New(), FromUserID: fromUserID, ToUserID: toUserID, Status: domain.RequestStatusPending}
	event, err := events.BuildFriendRequested(req.ID, req.FromUserID, req.ToUserID)
	if err != nil {
		return nil, err
	}
	res, err := svc.graphRepo.CreateRequestWithOutbox(ctx, req, event)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (svc *GraphService) AcceptRequest(ctx context.Context, userID, reqID uuid.UUID) (*domain.FriendRequest, error) {
	req, err := svc.graphRepo.AcceptRequestWithOutbox(ctx, userID, reqID, events.BuildFriendAccepted)
	if err != nil {
		return nil, err
	}
	return req, nil
}
func (svc *GraphService) RejectRequest(ctx context.Context, userID, reqID uuid.UUID) (*domain.FriendRequest, error) {
	req, err := svc.graphRepo.RejectRequestWithOutbox(ctx, userID, reqID, events.BuildFriendRejected)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (svc *GraphService) CancelRequest(ctx context.Context, userID, reqID uuid.UUID) (*domain.FriendRequest, error) {
	req, err := svc.graphRepo.CancelRequestWithOutbox(ctx, userID, reqID, events.BuildFriendCanceled)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (svc *GraphService) ListFriends(ctx context.Context, userID uuid.UUID) ([]domain.Friend, error) {
	friends, err := svc.graphRepo.ListFriends(ctx, userID)
	if err != nil {
		return nil, err
	}
	return friends, nil
}
func (svc *GraphService) ListIncoming(ctx context.Context, userID uuid.UUID) ([]domain.FriendRequest, error) {
	incoming, err := svc.graphRepo.ListIncoming(ctx, userID)
	if err != nil {
		return nil, err
	}
	return incoming, nil
}

func (svc *GraphService) ListOutgoing(ctx context.Context, userID uuid.UUID) ([]domain.FriendRequest, error) {
	outgoing, err := svc.graphRepo.ListOutgoing(ctx, userID)
	if err != nil {
		return nil, err
	}
	return outgoing, nil
}

func (svc *GraphService) GetStatus(ctx context.Context, userID, toUserID uuid.UUID) (*domain.Relation, error) {
	relation, err := svc.graphRepo.GetStatus(ctx, userID, toUserID)
	if err != nil {
		return nil, err
	}
	return relation, nil
}
func (svc *GraphService) RemoveFriendship(ctx context.Context, userIDA, userIDB uuid.UUID) error {
	userA, userB, err := domain.CanonicalPair(userIDA, userIDB)
	if err != nil {
		return err
	}
	event, err := events.BuildFriendRemoved(userIDA, userIDB)
	if err != nil {
		return err
	}
	err = svc.graphRepo.RemoveFriendshipWithOutbox(ctx, userA, userB, event)
	if err != nil {
		return err
	}
	return nil
}
