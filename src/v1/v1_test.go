package v1

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"gopkg.in/gin-gonic/gin.v1"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
)

func createRoutes() *gin.Engine {
	router := gin.New()

	Bind(router)

	return router
}

func TestGetToken200(t *testing.T) {
	userPassPair := fmt.Sprintf("%s:%s", os.Getenv("ADMIN_NAME"), os.Getenv("ADMIN_PASS"))

	authZ := base64.StdEncoding.EncodeToString([]byte(userPassPair))

	type token struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	gofight.New().
		GET("/v1/token").
		SetHeader(gofight.H{
			"Authorization": fmt.Sprintf("Basic %s", authZ),
		}).
		Run(createRoutes(), func(res gofight.HTTPResponse, req gofight.HTTPRequest) {
			assert.Equal(t, 200, res.Code)

			var tok token

			if err := json.Unmarshal(res.Body.Bytes(), &tok); err != nil {
				t.Fatalf(`Failed to parse response for /v1/token: %s`, err)
				return
			}

			assert.Nil(t, verifyAccessToken([]byte(os.Getenv("TOKEN_SECRET")), tok.AccessToken))
			assert.Nil(t, verifyRefreshToken([]byte(os.Getenv("TOKEN_SECRET")), tok.RefreshToken))
		})
}

func TestGetToken401(t *testing.T) {
	userPassPair := "not_admin:invalid_password"

	authZ := base64.StdEncoding.EncodeToString([]byte(userPassPair))

	gofight.New().
		GET("/v1/token").
		SetHeader(gofight.H{
			"Authorization": fmt.Sprintf("Basic %s", authZ),
		}).
		Run(createRoutes(), func(res gofight.HTTPResponse, req gofight.HTTPRequest) {
			assert.Equal(t, 401, res.Code)
		})
}

func TestGetAccessToken200(t *testing.T) {
	secret := []byte(os.Getenv("TOKEN_SECRET"))

	refreshToken, err := generateRefreshToken(secret)

	if err != nil {
		t.Fatalf(`Failed to generate refresh token: %s`, err)
		return
	}

	gofight.New().
		GET("/v1/token/access_token").
		SetHeader(gofight.H{
			"Authorization": fmt.Sprintf("Bearer %s", refreshToken),
		}).
		Run(createRoutes(), func(res gofight.HTTPResponse, req gofight.HTTPRequest) {
			assert.Equal(t, 200, res.Code)

			assert.Nil(t, verifyAccessToken(secret, res.Body.String()))
		})
}
