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

func (svc *GraphService) CreateRequest(ctx context.Context, fromUserID, toUserID uuid.UUID) (domain.FriendRequest, error) {
	if fromUserID == toUserID {
		return domain.FriendRequest{}, domain.ErrSelfRequest
	}
	req := domain.FriendRequest{ID: uuid.New(), FromUserID: fromUserID, ToUserID: toUserID, Status: domain.RequestStatusPending}
	event, err := events.BuildFriendRequested(req.ID, req.FromUserID, req.ToUserID)
	if err != nil {
		return domain.FriendRequest{}, err
	}
	res, err := svc.graphRepo.CreateRequestWithOutbox(ctx, req, event)
	if err != nil {
		return domain.FriendRequest{}, err
	}
	return res, nil
}
