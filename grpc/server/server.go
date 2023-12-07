package server

import (
	"context"
	"log"
	"net"

	proto "github.com/BerryTracer/user-service/grpc/proto"
	"github.com/BerryTracer/user-service/service"
	"google.golang.org/grpc"
)

type UserGRPCServer struct {
	UserService service.UserService
	proto.UnimplementedUserServiceServer
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
	proto.RegisterUserServiceServer(server, s)
	log.Printf("gRPC server listening on port %s\n", port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
		return err
	}
	return nil
}

func (s *UserGRPCServer) GetUserById(ctx context.Context, req *proto.GetUserByIdRequest) (*proto.User, error) {
	user, err := s.UserService.GetUserById(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return user.ConvertToProto(), nil
}

func (s *UserGRPCServer) GetUserByEmail(ctx context.Context, req *proto.GetUserByEmailRequest) (*proto.User, error) {
	user, err := s.UserService.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, err
	}

	return user.ConvertToProto(), nil
}

func (s *UserGRPCServer) GetUserByUsername(ctx context.Context, req *proto.GetUserByUsernameRequest) (*proto.User, error) {
	user, err := s.UserService.GetUserByUsername(ctx, req.GetUsername())
	if err != nil {
		return nil, err
	}

	return user.ConvertToProto(), nil
}

func (s *UserGRPCServer) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.User, error) {
	user, err := s.UserService.CreateUser(ctx, req.GetUsername(), req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return user.ConvertToProto(), nil
}
