package utils

import (
	"HarvestBox/config"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

func GenerateToken(userId int, email string, role string) (string, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Failed to load configuration:", err)
		return "", err
	}
	now := time.Now()

	claims := CustomClaims{
		Email: email,
		Role:  role,
		StandardClaims: jwt.StandardClaims{
			Id:        fmt.Sprint(userId),
			ExpiresAt: now.Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JwtSceret))
	if err != nil {
		fmt.Println("Failed JWT setup")
		return "", err
	}

	return tokenString, nil
}

func RenewToken(email string, role string) (string, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Failed to load configuration:", err)
		return "", err
	}
	now := time.Now()

	claims := CustomClaims{
		Email: email,
		Role:  role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(time.Hour * 24).Unix(), // Token valid for 24 hours
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JwtSceret))
	if err != nil {
		fmt.Println("Failed JWT setup")
		return "", err
	}

	fmt.Println("JWT token successfully renewed")
	return tokenString, nil
}

func ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			cfg, err := config.LoadConfig()
			if err != nil {
				return nil, fmt.Errorf("error loading configuration: %w", err)
			}
			return []byte(cfg.JwtSceret), nil
		},
	)

	if err != nil || !token.Valid {
		fmt.Println("Not a valid token")
		return nil, errors.New("not a valid token")
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		fmt.Println("Can't parse the claims")
		return nil, errors.New("can't parse the claims")
	}

	return claims, nil
}
