package service

import (
	"fmt"
	"context"
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
		return "", false, fmt.Errorf("updating session: %w", err)
	}

	if ttl < (s.mn.TTL() / 2) {
		err = s.repo.Update(ctx, sessionID, s.mn.TTL())
		if err != nil {
			return "", false, fmt.Errorf("updating session: %w", err)
		}
		return userID, true, nil
	}

	return userID, false, nil
}
