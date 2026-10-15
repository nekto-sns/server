package service

import (
	"fmt"
	"time"
	"context"

	"github.com/nekto-sns/server/domain"
)

type scratchAuthRepository interface{
	UserIDByScratchName(context.Context, string) (string, error)
}

type sessionRepository interface{
	Save(context.Context, string, string, time.Duration) error
	ByID(ctx context.Context, sessionID string) (string, time.Duration, error)
	Update(context.Context, string, time.Duration) error
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
		return "", fmt.Errorf("verifying scratch user: %w", err)
	}

	userID, err := s.saRepo.UserIDByScratchName(ctx, scratchName)
	if err != nil {
		return "", fmt.Errorf("verifying scratch user: %w", err)
	}

	sessionID, err := s.sessionMN.Generate()
	if err != nil {
		return "", fmt.Errorf("generating session id: %w: %w", err, domain.ErrInternal)
	}

	err = s.sessionRepo.Save(ctx, sessionID, userID, s.sessionMN.TTL())
	if err != nil {
		return "", fmt.Errorf("saving session: %w", err)
	}

	return sessionID, nil
}
