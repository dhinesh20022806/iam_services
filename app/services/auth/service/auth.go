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
	id string
	email string
	password string
	name string
	jwt.RegisteredClaims
}

func NewUser () *User {
	return &User{}
}


func GenerateToken(email string, id string) (string, error) {
	registed := jwt.RegisteredClaims{
		Subject: email, ID: id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer: "my-go-api",
		}

	user := User{
		email: email,
		id: id,
		RegisteredClaims: registed, 
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, user)
	return token.SignedString(secretKey)
}



func (u *User) IsValidUser(ctx context.Context, user *genpb.IsValidUserRequest) (string, string, error) {

	var userRes User

	query := `SELECT id,email,password,name FROM users WHERE email=$1`;

	err := dbx.Conn.QueryRow(query, user.Email).Scan(
		&userRes.id,
		&userRes.email,
		&userRes.password,
		&userRes.name,
	);

	if err != nil {
		fmt.Printf("Error retrieving user from database: %v\n", err)
		return "", "", fmt.Errorf("invalid user or password!")
	}

	isValid, err := scrypt.ComparePassword(user.Password, userRes.password)

	if err != nil {
		fmt.Printf("Error comparing password: %v\n", err)
		return "", "", fmt.Errorf("invalid user or password!")
	}

	if !isValid {
		return "", "", fmt.Errorf("invalid user or password!")
	}

   token, err :=  GenerateToken(userRes.email, userRes.id)



   if err != nil {
	fmt.Printf("Error generating token: %v\n", err)
	return "", "", fmt.Errorf("invalid user or password!");
   }

   return userRes.name, token, nil
}

func (u *User) CreateUser(ctx context.Context, user *genpb.CreateUserRequest) (string, error) {

	fmt.Printf("Creating user with email: %s and name: %s\n", user.Email, user.Name)
	query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3)`;

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

		return "",fmt.Errorf("Failed to create user: %v", err)
	}
	var newUser User
	query = `SELECT id,email FROM users WHERE email=?`;
	
	err = dbx.Conn.QueryRow(query, user.Email).Scan(
		&newUser.id,
		&newUser.email,
	);

	if err != nil {
		return "", fmt.Errorf("Failed to retrieve user after creation: %v", err)
	}

	return GenerateToken(newUser.email, newUser.id)
}






