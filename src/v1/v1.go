package v1

import (
	"fmt"
	"os"

	"gopkg.in/gin-gonic/gin.v1"
)

// Bind v1 routes to main router
func Bind(router *gin.Engine) error {
	v1 := router.Group("/v1")

	accounts, err := buildAdminAccounts()

	if err != nil {
		return fmt.Errorf(`Failed to bind v1 routes: %s`, err)
	}

	v1.GET("/token", getAdminToken, gin.BasicAuth(accounts))

	v1.GET("/token/access_token", refreshAccessToken)

	return nil
}

// Build admin account from enviromental variables
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
