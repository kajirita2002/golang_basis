package entwrapper

import (
	"context"

	"github.com/kajirita2002/golang_basis/domain/model/db"
	"github.com/kajirita2002/golang_basis/external/ent"
)

type UserImpl struct {
	client *ent.Client
}

func NewUser(client *ent.Client) *UserImpl {
	return &UserImpl{
		client: client,
	}
}

func (m *UserImpl) Get(ctx context.Context, userID int64) (*db.User, error) {
	user, err := m.client.User.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	dbUser := &db.User{
		ID: user.ID,
	}
	return dbUser, nil
}
