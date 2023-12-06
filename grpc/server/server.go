package server

import (
	"context"
	"log"
	"net"

	user_service "github.com/BerryTracer/user-service/grpc/proto"
	"github.com/BerryTracer/user-service/service"
	"google.golang.org/grpc"
)

type UserGRPCServer struct {
	UserService service.UserService
	user_service.UnimplementedUserServiceServer
}

func NewUserGRPCServer(userService service.UserService) *UserGRPCServer {
	return &UserGRPCServer{
		UserService: userService,
	}
}

func (s *UserGRPCServer) Run(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
		return err
	}
	server := grpc.NewServer()
	user_service.RegisterUserServiceServer(server, s)
	log.Printf("gRPC server listening on port %s\n", port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
		return err
	}
	return nil
}

func (s *UserGRPCServer) GetUserById(ctx context.Context, req *user_service.GetUserRequest) (*user_service.User, error) {
	user, err := s.UserService.GetUserById(ctx, req.GetId())
	if err != nil {
		log.Fatalf("failed to get user by id: %v\n", err)
		return nil, err
	}

	return user.ConvertToProto(), nil
}

func (s *UserGRPCServer) GetUserByEmail(ctx context.Context, req *user_service.GetUserByEmailRequest) (*user_service.User, error) {
	user, err := s.UserService.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		log.Fatalf("failed to get user by email: %v\n", err)
		return nil, err
	}

	return user.ConvertToProto(), nil
}

func (s *UserGRPCServer) GetUserByUsername(ctx context.Context, req *user_service.GetUserByUsernameRequest) (*user_service.User, error) {
	user, err := s.UserService.GetUserByUsername(ctx, req.GetUsername())
	if err != nil {
		log.Fatalf("failed to get user by username: %v\n", err)
		return nil, err
	}

	return user.ConvertToProto(), nil
}
