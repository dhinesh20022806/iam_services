package cmd

import (
	"fmt"
	"log/slog"
	"net"
	"time"

	genpb "github.com/dhinesh20022806/iam_service/app/proto/gen"
	"github.com/dhinesh20022806/iam_service/app/service"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
)


type gRPCServer struct {
	addr string
}

type CredGrpcHandler struct {
	genpb.UnimplementedLoginServiceServer
	genpb.UnimplementedSignUpServiceServer
}


func NewGrpcServer(addr string) *gRPCServer {
	return &gRPCServer{addr: addr}
}

var secretKey = []byte("secret")

func Login(email string, password string) {
	
}

type User struct {
	email string
	jwt.RegisteredClaims
}

func GenerateToken(email string) (string, error) {


	registed := jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer: "my-go-api",
		}

	user := User{
		email: email,
		RegisteredClaims: registed, 
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, user)
	return token.SignedString(secretKey)
}

func (s *gRPCServer) Run() error {
	lis, err := net.Listen("tcp", s.addr)

	if err != nil {
		return fmt.Errorf("error create tcp port", "err", err)
	}

	grpcServer := grpc.NewServer()



	service.NewGrpcLoginService(grpcServer)

	return grpcServer.Serve(lis)
}


func LoginUser() {


	query := `SELECT email,password FROM users WHERE email=?`;

	

}



func RunServer() int {

	server := NewGrpcServer(":50050")

	err := server.Run()

	if err != nil {
		slog.Error("Error creating server", "err", err)
		return 0
	}



	return 1
}