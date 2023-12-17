package main

import (
	"log"
	"net"

	"github.com/BerryTracer/common-service/adapter/database/mongodb"
	"github.com/BerryTracer/common-service/config"
	"github.com/BerryTracer/common-service/crypto"
	"github.com/BerryTracer/user-service/database"
	user_service "github.com/BerryTracer/user-service/grpc/proto"
	"github.com/BerryTracer/user-service/grpc/server"
	"github.com/BerryTracer/user-service/repository"
	"github.com/BerryTracer/user-service/service"
	"google.golang.org/grpc"
)

func main() {
	// Load environment configurations
	mongodbURI := getEnvOrPanic("MONGODB_URI")
	grpcPort := getEnvWithDefaultOrPanic("GRPC_PORT", "50051")

	// Initialize the database
	db := initDatabase(mongodbURI)
	defer func(db *database.UserMongoDatabase) {
		err := db.Disconnect()
		if err != nil {
			log.Fatalf("failed to listen: %v\n", err)
		}
	}(db)

	// Set up the gRPC server and start listening
	grpcServer := setupGRPCServer(db)
	startGRPCServer(grpcServer, grpcPort)
}

func getEnvOrPanic(key string) string {
	value, err := config.LoadEnv(config.NewRealEnvLoader(), key)
	if err != nil {
		panic(err)
	}
	return value
}

func getEnvWithDefaultOrPanic(key, defaultValue string) string {
	value, err := config.LoadEnvWithDefault(config.NewRealEnvLoader(), key, defaultValue)
	if err != nil {
		panic(err)
	}
	return value
}

func initDatabase(mongodbURI string) *database.UserMongoDatabase {
	db, err := database.NewUserMongoDatabaseConnection(mongodbURI, "user", "user")
	if err != nil {
		panic(err)
	}
	return db
}

func setupGRPCServer(db *database.UserMongoDatabase) *grpc.Server {
	mongoDBAdapter := mongodb.NewMongoAdapter(db.Collection)
	userRepository := repository.NewUserMongoRepository(mongoDBAdapter)
	passwordHasher := crypto.NewBcryptHasher()
	userService := service.NewUserService(userRepository, passwordHasher)

	gGRPCServer := server.NewUserGRPCServer(userService)

	grpcServer := grpc.NewServer()
	user_service.RegisterUserServiceServer(grpcServer, gGRPCServer)

	return grpcServer
}

func startGRPCServer(grpcServer *grpc.Server, grpcPort string) {
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	log.Println("gRPC server listening on port " + grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC server: %v\n", err)
	}
}
