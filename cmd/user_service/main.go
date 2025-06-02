package main

import (
	"log"
	"net"

	"github.com/oj-lab/reborn/common/gorm_client"
	userpb "github.com/oj-lab/reborn/protobuf/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	db := gorm_client.GetDB()
	db.AutoMigrate(&UserModel{})

	repo := NewGormUserRepository(db)
	service := NewUserService(repo)

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, service)

	// 开启 gRPC ServerReflection
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("UserService gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
