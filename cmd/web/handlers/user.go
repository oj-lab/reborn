package handlers

import (
	"github.com/gin-gonic/gin"
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

// CreateUser
//
//	@Summary		Create a new user
//	@Description	Create a new user with the provided details
//	@Tags			User
//	@Router			/user [post]
//	@Accept			json
//	@Produce		json
//	@Param			user	body	userpb.CreateUserRequest	true	"User details"
//	@Success		200
func CreateUser(ginCtx *gin.Context) {
	var req userpb.CreateUserRequest
	if err := ginCtx.Bind(&req); err != nil {
		return
	}
	_, err := userServiceClient.CreateUser(ginCtx, &req)
	if err != nil {
		_ = ginCtx.AbortWithError(500, err)
		return
	}
	ginCtx.Status(200)
}
