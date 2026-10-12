package service

import (
	"fmt"
	"context"
	"errors"

	"github.com/nekto-sns/server/model"
)

type sessionService struct{
	repo sessionRepository
	mn   sessionManager
}

func NewSession(s sessionRepository, m sessionManager) *sessionService {
	return &sessionService{
		repo: s,
		mn:   m,
	}
}

func (s *sessionService) Update(ctx context.Context, sessionID string) (string, bool, error) {
	userID, ttl, err := s.repo.ByID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return "", false, fmt.Errorf("Session not found: %w", err)
		}
		return "", false, fmt.Errorf("Failed to get session: %w", err)
	}

	if ttl < (s.mn.TTL() / 2) {
		err = s.repo.Update(ctx, sessionID, s.mn.TTL())
		if err != nil {
			return "", false, fmt.Errorf("Failed to update session: %w", err)
		}
		return userID, true, nil
	}

	return userID, false, nil
}
