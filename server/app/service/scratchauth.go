package service

import (
	"fmt"
	"time"
	"errors"
	"context"

	"github.com/nekto-sns/nekto/server/app/model"
)

type scratchAuthRepository interface{
	UserIDByScratchName(context.Context, string) (string, error)
}

type sessionRepository interface{
	Save(context.Context, string, string, time.Duration) error
}

type scratchAuthClient interface{
	Verify(context.Context, string) (string, error)
}

type sessionManager interface{
	Generate() (string, error)
	TTL() (time.Duration)
}

type scratchAuthService struct{
	saRepo      scratchAuthRepository
	sessionRepo sessionRepository
	saClient    scratchAuthClient
	sessionMN   sessionManager
}

func NewScratchAuth(sessionRepo sessionRepository, saRepo scratchAuthRepository, saClient scratchAuthClient, sessionMN sessionManager) *scratchAuthService {
	return &scratchAuthService{
		saRepo:      saRepo,
		sessionRepo: sessionRepo,
		saClient:    saClient,
		sessionMN:   sessionMN,
	}
}

func (s *scratchAuthService) LoginCallback(ctx context.Context, privateCode string) (string, error) {
	scratchName, err := s.saClient.Verify(ctx, privateCode)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return "", fmt.Errorf("Scratch account not found: %w", err)
		}
		return "", fmt.Errorf("Authentication failed: %w", err)
	}

	userID, err := s.saRepo.UserIDByScratchName(ctx, scratchName)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return "", fmt.Errorf("User not found: %w", err)
		}
		return "", fmt.Errorf("Failed to get user: %w", err)
	}

	sessionID, err := s.sessionMN.Generate()
	if err != nil {
		return "", fmt.Errorf("Failed to generate session ID(%v): %w", err, model.ErrInternal)
	}

	err = s.sessionRepo.Save(ctx, sessionID, userID, s.sessionMN.TTL())
	if err != nil {
		return "", fmt.Errorf("Failed to save session ID: %w", err)
	}

	return sessionID, nil
}
