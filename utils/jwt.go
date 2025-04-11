package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretkey = "secret"

func GenerateToken(email string, adminId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":   email,
		"adminId": adminId,
		"exp":     time.Now().Add(time.Hour).Unix(),
	})

	return token.SignedString([]byte(secretkey))
}

func VerifyToken(token string) (int64, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretkey), nil
	})

	if err != nil {
		return 0, fmt.Errorf("could not parse token: %w", err)
	}

	tokenIsVaild := parsedToken.Valid

	if !tokenIsVaild {
		return 0, errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return 0, errors.New("invalid token claims")
	}

	adminIdFloat, ok := claims["adminId"].(float64)

	if !ok {
		return 0, errors.New("adminId not found or invalid type")
	}
	adminId := int64(adminIdFloat)

	return adminId, nil
}

func GenerateBusToken(busId int64, name string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"busId": busId,
		"name":  name,
	})

	return token.SignedString([]byte(secretkey))
}

func VerifyBusToken(token string) (int64, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretkey), nil
	})

	if err != nil {
		return 0, fmt.Errorf("could not parse token: %w", err)
	}

	tokenIsVaild := parsedToken.Valid

	if !tokenIsVaild {
		return 0, errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return 0, errors.New("invalid token claims")
	}

	busIdFloat, ok := claims["busId"].(float64)

	if !ok {
		return 0, errors.New("busId not found or invalid type")
	}
	busId := int64(busIdFloat)

	return busId, nil
}
