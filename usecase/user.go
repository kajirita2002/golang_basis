package usecase

import (
	"context"

	"github.com/kajirita2002/golang_basis/domain/model"
	"github.com/kajirita2002/golang_basis/domain/service"
)

type IFUser interface {
	Get(ctx context.Context, userID int64) (*model.User, error)
}

type UserImpl struct {
	user service.IFUser
}

// NewUser returns new UserImpl
func NewUser(user service.IFUser) *UserImpl {
	return &UserImpl{
		user: user,
	}
}

// Get get user.
func (u *UserImpl) Get(ctx context.Context, userID int64) (*model.User, error) {
	user, err := u.user.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
