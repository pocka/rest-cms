package v1

import (
	"fmt"
	"regexp"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"gopkg.in/gin-gonic/gin.v1"
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
		IssuedAt:  time.Now().Local().Unix(),
		Issuer:    jwtIssuer,
		Subject:   jwtRefreshTokenSubject,
	}).SignedString(secret)
}

// Returns access token and refresh token for admin user
func getAdminToken(c *gin.Context) {
	secret := []byte("admin")

	accessToken, err := generateAccessToken(secret)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := generateRefreshToken(secret)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	c.JSON(200, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
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

// Get new access token using refresh token
func refreshAccessToken(c *gin.Context) {
	secret := []byte("admin")

	authHeader := c.Request.Header.Get("Authorization")

	refreshToken := regexp.MustCompile(`^Bearer\s`).ReplaceAllString(authHeader, "")

	err := verifyRefreshToken(secret, refreshToken)

	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid refresh token"})
		return
	}

	accessToken, err := generateAccessToken(secret)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate access token"})
		return
	}

	c.String(200, accessToken)
}
