package services

import (
	"fmt"
	"net"

	handler "github.com/dhinesh20022806/iam_service/app/services/auth/handlers"
	authService "github.com/dhinesh20022806/iam_service/app/services/auth/service"
	"google.golang.org/grpc"
)



type gRPCServer struct {
	addr string
}

func NewGrpcServer(addr string) *gRPCServer {
	return &gRPCServer{addr: addr}
}





func (s *gRPCServer) Run() error {
	lis, err := net.Listen("tcp", s.addr)

	if err != nil {
		return fmt.Errorf("error create tcp port", "err", err)
	}

	user := authService.NewUser();
	grpcServer := grpc.NewServer()
	
	fmt.Printf("grpc service running on %s\n", s.addr)
	handler.NewGrpcAuthService(grpcServer, user)


	return grpcServer.Serve(lis)
}
