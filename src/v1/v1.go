package v1

import (
	"fmt"
	"os"
	"regexp"

	"gopkg.in/gin-gonic/gin.v1"
)

// Bind v1 routes to main router
func Bind(router *gin.Engine) error {
	v1 := router.Group("/v1")

	accounts, err := buildAdminAccounts()

	if err != nil {
		return fmt.Errorf(`Failed to bind v1 routes: %s`, err)
	}

	v1.GET("/token", gin.BasicAuth(accounts), getAdminToken)

	v1.GET("/token/access_token", refreshAccessToken)

	return nil
}

// Build admin account from environmental variables
// ADMIN_NAME ... Name of admin user (required).
// ADMIN_PASS ... Password for admin user (required).
func buildAdminAccounts() (map[string]string, error) {
	adminName := os.Getenv("ADMIN_NAME")
	adminPass := os.Getenv("ADMIN_PASS")

	if !(adminName != "" && adminPass != "") {
		return nil, fmt.Errorf(`You MUST specify ADMIN_USER and ADMIN_PASS`)
	}

	accounts := make(map[string]string)

	accounts[adminName] = adminPass

	return accounts, nil
}

// Returns access token and refresh token for admin user
func getAdminToken(c *gin.Context) {
	secret := []byte(os.Getenv("TOKEN_SECRET"))

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

// Get new access token using refresh token
func refreshAccessToken(c *gin.Context) {
	secret := []byte(os.Getenv("TOKEN_SECRET"))

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
