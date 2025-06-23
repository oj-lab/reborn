package main

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oj-lab/reborn/common/app"
	userpb "github.com/oj-lab/reborn/protobuf/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var jwtSecret []byte

func init() {
	cwd, _ := os.Getwd()
	app.Init(cwd, "user_service")
	jwtSecret = []byte(app.Config().GetString("jwt.secret"))
}

type Claims struct {
	UserID uint64          `json:"user_id"`
	Role   userpb.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(user *userpb.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: user.Id,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "user_login",
			Issuer:    app.Config().GetString("jwt.issuer"),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type contextKey string

const claimsContextKey = contextKey("claims")

func ValidateJWT(ctx context.Context) (*Claims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	authHeader := values[0]
	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
		return nil, status.Errorf(codes.Unauthenticated, "invalid authorization token format")
	}
	tokenString := bearerToken[1]

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	}, jwt.WithIssuer(app.Config().GetString("jwt.issuer")))

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, status.Errorf(codes.Unauthenticated, "invalid signature")
		}
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	if !token.Valid {
		return nil, status.Errorf(codes.Unauthenticated, "token is not valid")
	}

	return claims, nil
}

func AuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if publicMethods[info.FullMethod] {
		return handler(ctx, req)
	}

	claims, err := ValidateJWT(ctx)
	if err != nil {
		return nil, err
	}

	// Store claims in context for use in RPC methods
	newCtx := context.WithValue(ctx, claimsContextKey, claims)

	return handler(newCtx, req)
}
