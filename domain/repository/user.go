package repository

import (
	"context"

	"github.com/kajirita2002/golang_basis/domain/model/db"
)

type IFUser interface {
	Get(ctx context.Context, userID int64) (*db.User, error)
}
