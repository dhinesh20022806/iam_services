package scrypt

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/scrypt"
)


const (
	scryptN = 32768
	scryptR = 8
	scryptP = 1
	saltBytes = 16
	keyHashBytes = 32
)



func HashPassword(password string) (string, error) {

	salt := make([]byte, saltBytes)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	derivedKey, err := scrypt.Key([]byte(password), salt, scryptN, scryptR, scryptP, keyHashBytes)
	if err != nil {
		return "", err
	}
	saltEncoded := base64.RawStdEncoding.EncodeToString(salt)
	hashEncoded := base64.RawStdEncoding.EncodeToString(derivedKey)

	formattedHash := fmt.Sprintf("$scrypt$%d$%d$%d$%s$%s", scryptN, scryptR, scryptP, saltEncoded, hashEncoded)

	return formattedHash, nil
}

func ComparePassword(comparePassword, hashedPassword string) (bool, error) {
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 7|| parts[1] != "script" {
		return false, errors.New("invalid hash format")
	}
	
	var n, r, p int

	_, err := fmt.Sscanf(parts[2]+"-"+parts[3]+"-"+parts[4], "%d-%d-%d", &n, &r, &p)

	if err != nil {
		return false, errors.New("failed to parse scrypt parameter")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[6])

	if err != nil {
		return false, err
	}

	derivedKey, err := scrypt.Key([]byte(comparePassword), salt, n, r, p, len(expectedHash))

	if err != nil {
		return false, err
	}

	if subtle.ConstantTimeCompare(derivedKey, expectedHash) == 1 {
		return true, nil
	}

	return false, nil
}
