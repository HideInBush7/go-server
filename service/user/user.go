package user

import (
	"context"

	"github.com/HideInBush7/go-server/pb"
)

const server = "/user"

type UserServer struct {
	pb.UnimplementedUserServer
}

func (s *UserServer) Register(context.Context, *pb.RegisterRequest) (*pb.UserInfoReply, error) {
	panic(``)
}

func (s *UserServer) Login(context.Context, *pb.LoginRequest) (*pb.UserInfoReply, error) {
	panic(``)
}

func (s *UserServer) GetUserInfo(context.Context, *pb.UserInfoRequest) (*pb.UserInfoReply, error) {
	panic(``)
}

func (s *UserServer) GetUserList(context.Context, *pb.UserListRequest) (*pb.UserList, error) {
	panic(``)
}

func (s *UserServer) UpdateUser(context.Context, *pb.UpdateUserRequest) (*pb.UserInfoReply, error) {
	panic(``)
}
