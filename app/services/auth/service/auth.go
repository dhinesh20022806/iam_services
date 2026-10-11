package service

import (
	"context"
	"fmt"
	"time"

	genpb "github.com/dhinesh20022806/grpc_service/proto/gen"
	"github.com/dhinesh20022806/iam_service/app/pkg/dbx"
	"github.com/dhinesh20022806/iam_service/app/pkg/scrypt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
)

var secretKey = []byte("iloveyou3000")




type User struct {
	email string
	password string
	name string
	jwt.RegisteredClaims
}

func NewUser () *User {
	return &User{}
}


func GenerateToken(email string, name string) (string, error) {
	registed := jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer: "my-go-api",
		}

	user := User{
		email: email,
		name: name,
		RegisteredClaims: registed, 
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, user)
	return token.SignedString(secretKey)
}



func (u *User) IsValidUser(ctx context.Context, user *genpb.IsValidUserRequest) (string, string, error) {

	var userRes User

	query := `SELECT email,password,name FROM users WHERE email=?`;

	err := dbx.Conn.QueryRow(query, user.Email).Scan(
		&userRes.email,
		&userRes.password,
		&userRes.name,
	);

	if err != nil {
		return "", "", fmt.Errorf("invalid user or password!")
	}

	isValid, err := scrypt.ComparePassword(user.Password, userRes.password)

	if err != nil {
		return "", "", fmt.Errorf("invalid user or password!")
	}

	if !isValid {
		return "", "", fmt.Errorf("invalid user or password!")
	}

   token, err :=  GenerateToken(userRes.email, userRes.name)

   if err != nil {
	return "", "", fmt.Errorf("invalid user or password!");
   }

   return userRes.name, token, nil
}

func (u *User) CreateUser(ctx context.Context, user *genpb.CreateUserRequest) (string, error) {

	query := `INSERT INTO users (name, email, password) VALUES (?, ?, ?)`;

	hashedPassword, err := scrypt.HashPassword(user.Password)

	if err != nil {

	}
	_, err = dbx.Conn.Exec(query, user.Name, user.Email, hashedPassword)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {

			// duplicate postgres error code
			if pqErr.Code == "23505" {
				
				return "", fmt.Errorf("User already Created")

			}
		}
	}
	return GenerateToken(user.Email, user.Name)
}






