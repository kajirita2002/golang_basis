package handler

import (
	"context"

	"github.com/kajirita2002/golang_basis/server/api/grpc/pb"
	"github.com/kajirita2002/golang_basis/usecase"
)

// User service handler
type User struct {
	user usecase.IFUser
	pb.UnimplementedUserServiceServer
}

// NewUser returns new User.
func NewUser(user usecase.IFUser) *User {
	return &User{
		user: user,
	}
}

func (u *User) Get(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := u.user.Get(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserResponse{
		User: &pb.User{
			Id: user.ID,
		},
	}, nil
}
