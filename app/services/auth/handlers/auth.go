package handler

import (
	"context"

	genpb "github.com/dhinesh20022806/iam_service/app/proto/gen"
	"github.com/dhinesh20022806/iam_service/app/services/auth/types"
	"google.golang.org/grpc"
)

type AuthServiceHandler struct {
	authService types.AuthService
	genpb.UnimplementedAuthServiceServer
}

func NewGrpcAuthService(grpc *grpc.Server, user types.AuthService) {
	cred := &AuthServiceHandler{
		authService: user,
	}

	genpb.RegisterAuthServiceServer(grpc, cred)
}



func (h *AuthServiceHandler) CreateUser(ctx context.Context, req *genpb.CreateUserRequest) (*genpb.CreateUserResponse, error) {
	token , err := h.authService.CreateUser(ctx, req)

	res := &genpb.CreateUserResponse{
		Email: req.Email,
		Name: req.Name,
		Jwt: token,
		Valid: true,
	}

	if err != nil {
		res.Valid = false
		return res, err
	}

	return res, nil
}



func (h *AuthServiceHandler) IsValidUser(ctx context.Context, req *genpb.IsValidUserRequest) (*genpb.IsValidUserResponse, error) {

	
	name, token, err := h.authService.IsValidUser(ctx, req)

	res := &genpb.IsValidUserResponse{
		Email: req.Email,
		Name: name,
		Jwt: token,
		Valid: true,
	}

	if err != nil {
		res.Valid = false

		return res, err
	}

	return res, nil
} 
