package grpc

import (
	"context"
	"messenger/gen"

	"google.golang.org/grpc"
)

type serverAPI struct {
	gen.UnimplementedAuthServer
}

func RegServer(gRPC *grpc.Server) {
	gen.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(ctx context.Context, req *gen.LoginRequest) (*gen.LoginResponse, error) {
	return &gen.LoginResponse{
		Token: req.Email,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *gen.RegisterRequest) (*gen.RegisterResponse, error) {
	panic("implement me")
}

func (s *serverAPI) IsAAdmin(ctx context.Context, req *gen.IsAdminRequest) (*gen.IsAdminResponse, error) {
	panic("implement me")
}

func (s *serverAPI) Logout(ctx context.Context, req *gen.LogoutRequest) (*gen.LogoutResponse, error) {
	panic("implement me")
}
