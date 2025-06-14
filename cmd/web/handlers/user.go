package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	userpb "github.com/oj-lab/reborn/protobuf/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var userServiceClient userpb.UserServiceClient

func init() {
	var err error
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("failed to connect to user service: " + err.Error())
	}
	userServiceClient = userpb.NewUserServiceClient(conn)
}

func CreateUser(ctx echo.Context) error {
	var req userpb.CreateUserRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	_, err := userServiceClient.CreateUser(ctx.Request().Context(), &req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, map[string]string{"status": "success"})
}
