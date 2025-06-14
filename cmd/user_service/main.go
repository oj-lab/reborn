package main

import (
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"

	"github.com/oj-lab/reborn/common/app"
	"github.com/oj-lab/reborn/common/gorm_client"
	userpb "github.com/oj-lab/reborn/protobuf/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const configKeyServerPort = "server.port"

var port uint

func init() {
	cwd, _ := os.Getwd()
	app.Init(cwd, "user_service")
	port = app.Config().GetUint(configKeyServerPort)
}

func main() {
	db := gorm_client.GetDB()
	db.AutoMigrate(&UserModel{})

	repo := NewGormUserRepository(db)
	service := NewUserService(repo)

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, service)

	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	slog.Info("user service started", "port", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
