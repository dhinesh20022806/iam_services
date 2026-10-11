package service

import (
	genpb "github.com/dhinesh20022806/iam_service/app/proto/gen"
	"google.golang.org/grpc"
)


type AuthCred struct {
	genpb.UnimplementedLoginServiceServer

}

func NewGrpcLoginService(grpc *grpc.Server) {
	cred := AuthCred{} 


	genpb.RegisterLoginServiceServer(grpc, cred)
}