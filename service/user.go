package service

import (
	"fmt"
	"context"

	"github.com/nekto-sns/server/model"
)

type userRepository interface{
	ByName(context.Context, string) (*model.User, error)
	ByID(context.Context, string) (*model.User, error)
	Create(context.Context, string) (*model.User, error)
	IDByScratchName(context.Context, string) (string, error)
}

type userService struct{
	repo userRepository
}

func NewUser(repo userRepository) (*userService) {
	return &userService{
		repo: repo,
	}
}

func (u *userService) ByName(ctx context.Context, name string) (*model.User, error) {
	user, err := u.repo.ByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("getting user %s: %w", name, err)
	}

	return user, nil
}

func (u *userService) ByID(ctx context.Context, id string) (*model.User, error) {
	user, err := u.repo.ByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting user %s: %w", id, err)
	}

	return user, nil
}
