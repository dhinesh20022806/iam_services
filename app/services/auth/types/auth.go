package types

import (
	"context"

	genpb "github.com/dhinesh20022806/grpc_service/proto/gen"
)




type AuthService interface {
	IsValidUser(context.Context, *genpb.IsValidUserRequest) (string, string, error)
	CreateUser(context.Context, *genpb.CreateUserRequest) (string, error)
}


