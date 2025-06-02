package main

import (
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"path"

	"github.com/oj-lab/reborn/common/app"
	"github.com/oj-lab/reborn/common/gorm_client"
	userpb "github.com/oj-lab/reborn/protobuf/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var port uint

func init() {
	cwd, _ := os.Getwd()
	app.InitConfig(path.Join(cwd, "configs", "user_service"))
	port = app.Config().GetUint("server.port")
}

func main() {
	db := gorm_client.GetDB()
	db.AutoMigrate(&UserModel{})

	repo := NewGormUserRepository(db)
	service := NewUserService(repo)

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, service)

	// 开启 gRPC ServerReflection
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
