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

var (
	port          uint
	publicMethods map[string]bool
)

func init() {
	cwd, _ := os.Getwd()
	app.Init(cwd, "user_service")
	port = app.Config().GetUint(configKeyServerPort)
	jwtSecret = []byte(app.Config().GetString("jwt.secret"))

	// Load public methods for auth interceptor
	methods := app.Config().GetStringSlice("auth.public_methods")
	publicMethods = make(map[string]bool, len(methods))
	for _, method := range methods {
		publicMethods[method] = true
	}
}

func main() {
	db := gorm_client.GetDB()
	db.AutoMigrate(&UserModel{})

	repo := NewGormUserRepository(db)
	service := NewUserService(repo)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(AuthInterceptor),
	)
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
