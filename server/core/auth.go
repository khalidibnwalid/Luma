package core

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const jwtTimeToLive = 10 * 24 * time.Hour

func GenerateJwtToken(secret, userId string) (string, error) {
	headerAndPayload := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(jwtTimeToLive).Unix(),
	})

	jwtToken, err := headerAndPayload.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return jwtToken, nil
}

func ValidateJwtToken(secret, tokenString string) (*jwt.Token, jwt.MapClaims, bool) {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if !token.Valid {
		return nil, nil, false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, false
	}

	return token, claims, true
}
