package v1

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	jwtIssuer = "admin.rest-cms"

	jwtAccessTokenSubject  = "acc"
	jwtRefreshTokenSubject = "ref"
)

// Generates access token
func generateAccessToken(secret []byte) (string, error) {
	if len(secret) == 0 {
		return "", fmt.Errorf(`Empty secret is not supported`)
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		IssuedAt:  time.Now().Local().Unix(),
		ExpiresAt: time.Now().Local().Add(time.Minute * 10).Unix(),
		Issuer:    jwtIssuer,
		Subject:   jwtAccessTokenSubject,
	}).SignedString(secret)
}

// Generates refresh token
func generateRefreshToken(secret []byte) (string, error) {
	if len(secret) == 0 {
		return "", fmt.Errorf(`Empty secret is not supported`)
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		IssuedAt: time.Now().Local().Unix(),
		Issuer:   jwtIssuer,
		Subject:  jwtRefreshTokenSubject,
	}).SignedString(secret)
}

func verifyAccessToken(secret []byte, tokenString string) error {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(`Unexpected sigining method: %v`, token.Header["alg"])
		}

		return secret, nil
	})

	claims, ok := token.Claims.(*jwt.StandardClaims)

	if !(ok && token.Valid && err == nil) {
		return fmt.Errorf(`Invalid token`)
	}

	now := time.Now().Unix()

	if !claims.VerifyExpiresAt(now, true) {
		delta := time.Unix(now, 0).Sub(time.Unix(claims.ExpiresAt, 0))

		return fmt.Errorf(`Token is expired by %v`, delta)
	}

	if !claims.VerifyIssuedAt(now, true) {
		return fmt.Errorf(`Token used before issued`)
	}

	if !claims.VerifyIssuer(jwtIssuer, true) {
		return fmt.Errorf(`Invalid token issuer was passed`)
	}

	if claims.Subject != jwtAccessTokenSubject {
		return fmt.Errorf(`Passed token is not access token`)
	}

	return nil
}

func verifyRefreshToken(secret []byte, tokenString string) error {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(`Unexpected sigining method: %v`, token.Header["alg"])
		}

		return secret, nil
	})

	claims, ok := token.Claims.(*jwt.StandardClaims)

	if !(ok && token.Valid && err == nil) {
		return fmt.Errorf(`Invalid token`)
	}

	if !claims.VerifyIssuedAt(time.Now().Unix(), true) {
		return fmt.Errorf(`Token used before issued`)
	}

	if !claims.VerifyIssuer(jwtIssuer, true) {
		return fmt.Errorf(`Invalid token issuer was passed`)
	}

	if claims.Subject != jwtRefreshTokenSubject {
		return fmt.Errorf(`Passed token is not access token`)
	}

	return nil
}
