package scrypt

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
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

	fmt.Println("Hashed password parts:", hashedPassword, comparePassword) // Debugging line


    if len(parts) != 7 || parts[0] != "" || parts[1] != "scrypt" {
        return false, errors.New("invalid hash format")
    }

    n, err := strconv.Atoi(parts[2])
    if err != nil {
        return false, errors.New("invalid scrypt N parameter")
    }

    r, err := strconv.Atoi(parts[3])
    if err != nil {
        return false, errors.New("invalid scrypt r parameter")
    }

    p, err := strconv.Atoi(parts[4])
    if err != nil {
        return false, errors.New("invalid scrypt p parameter")
    }

    // Validate parameters before allocating resources.
    if n <= 1 || n&(n-1) != 0 || r <= 0 || p <= 0 {
        return false, errors.New("invalid scrypt parameters")
    }

    salt, err := base64.RawStdEncoding.DecodeString(parts[5])
    if err != nil {
        return false, errors.New("invalid salt encoding")
    }

    expectedHash, err := base64.RawStdEncoding.DecodeString(parts[6])
    if err != nil {
        return false, errors.New("invalid hash encoding")
    }

    if len(salt) == 0 || len(expectedHash) == 0 {
        return false, errors.New("empty salt or hash")
    }

    derivedKey, err := scrypt.Key(
        []byte(comparePassword),
        salt,
        n, r, p,
        len(expectedHash),
    )
    if err != nil {
        return false, err
    }

    return subtle.ConstantTimeCompare(derivedKey, expectedHash) == 1, nil
}