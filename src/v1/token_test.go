package v1

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAccessTokenOK(t *testing.T) {
	secret := []byte("foo")

	token, err := generateAccessToken(secret)

	assert.Nil(t, err)
	assert.NotEqual(t, "", token)
}

func TestGenerateAccessTokenNG(t *testing.T) {
	secret := []byte("")

	token, err := generateAccessToken(secret)

	assert.NotNil(t, err)
	assert.Equal(t, "", token)
}

func TestVerifyAccessTokenOK(t *testing.T) {
	secret := []byte("foo")

	token, _ := generateAccessToken(secret)

	err := verifyAccessToken(secret, token)

	assert.Nil(t, err)
}

func TestVerifyAccessTokenNGSecretNotMatch(t *testing.T) {
	token, _ := generateAccessToken([]byte("foo"))

	err := verifyAccessToken([]byte("bar"), token)

	assert.NotNil(t, err)
}

func TestVerifyAccessTokenNGModifiedAlg(t *testing.T) {
	secret := []byte("foo")

	token, _ := generateAccessToken(secret)

	bodies := strings.Split(token, ".")

	headerBytes, _ := base64.StdEncoding.DecodeString(bodies[0])

	var header struct {
		Alg string `json:"alg"`
		Typ string `json:"typ"`
	}

	err := json.Unmarshal(headerBytes, &header)

	if err != nil {
		t.Fatalf(`generateAccessToken returned jwt with no json serializable header`)
	}

	header.Alg = "none"

	b, _ := json.Marshal(header)

	bodies[0] = base64.StdEncoding.EncodeToString(b)

	modifiedToken := strings.Join(bodies, ".")

	err = verifyAccessToken(secret, modifiedToken)

	assert.NotNil(t, err)
}

func TestVerifyAccessTokenNGInvalidIssuedAt(t *testing.T) {
	secret := []byte("foo")

	issuedAt := time.Now().Local().Add(time.Minute * 5).Unix()
	expiresIn := time.Now().Local().Add(time.Minute * 10).Unix()

	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		IssuedAt:  issuedAt,
		ExpiresAt: expiresIn,
		Issuer:    jwtIssuer,
		Subject:   jwtAccessTokenSubject,
	}).SignedString(secret)

	err := verifyAccessToken(secret, token)

	assert.NotNil(t, err)
}

func TestVerifyAccessTokenNGInvalidIssuer(t *testing.T) {
	secret := []byte("foo")

	issuedAt := time.Now().Local().Add(time.Minute * -10).Unix()
	expiresIn := time.Now().Local().Add(time.Minute * 10).Unix()

	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		IssuedAt:  issuedAt,
		ExpiresAt: expiresIn,
		Issuer:    "rest-cms",
		Subject:   jwtAccessTokenSubject,
	}).SignedString(secret)

	err := verifyAccessToken(secret, token)

	assert.NotNil(t, err)
}

func TestVerifyAccessTokenNGPassedRefreshToken(t *testing.T) {
	secret := []byte("foo")

	issuedAt := time.Now().Local().Add(time.Minute * -10).Unix()
	expiresIn := time.Now().Local().Add(time.Minute * 10).Unix()

	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		IssuedAt:  issuedAt,
		ExpiresAt: expiresIn,
		Issuer:    jwtIssuer,
		Subject:   jwtRefreshTokenSubject,
	}).SignedString(secret)

	err := verifyAccessToken(secret, token)

	assert.NotNil(t, err)
}

func TestVerifyAccessTokenNGExpired(t *testing.T) {
	secret := []byte("foo")

	issuedAt := time.Now().Local().Add(time.Minute * -10).Unix()
	expiresIn := time.Now().Local().Add(time.Minute * -5).Unix()

	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		IssuedAt:  issuedAt,
		ExpiresAt: expiresIn,
		Issuer:    jwtIssuer,
		Subject:   jwtAccessTokenSubject,
	}).SignedString(secret)

	err := verifyAccessToken(secret, token)

	assert.NotNil(t, err)
}

func TestGenerateRefreshTokenOK(t *testing.T) {
	secret := []byte("foo")

	token, err := generateRefreshToken(secret)

	assert.Nil(t, err)
	assert.NotEqual(t, "", token)
}

func TestGenerateRefreshTokenNG(t *testing.T) {
	secret := []byte("")

	token, err := generateRefreshToken(secret)

	assert.NotNil(t, err)
	assert.Equal(t, "", token)
}

func TestVerifyRefreshTokenOK(t *testing.T) {
	secret := []byte("foo")

	token, _ := generateRefreshToken(secret)

	err := verifyRefreshToken(secret, token)

	assert.Nil(t, err)
}

func TestVerifyRefreshTokenNGSecretNotMatch(t *testing.T) {
	token, _ := generateRefreshToken([]byte("foo"))

	err := verifyRefreshToken([]byte("bar"), token)

	assert.NotNil(t, err)
}

func TestVerifyRefreshTokenNGModifiedAlg(t *testing.T) {
	secret := []byte("foo")

	token, _ := generateRefreshToken(secret)

	bodies := strings.Split(token, ".")

	headerBytes, _ := base64.StdEncoding.DecodeString(bodies[0])

	var header struct {
		Alg string `json:"alg"`
		Typ string `json:"typ"`
	}

	err := json.Unmarshal(headerBytes, &header)

	if err != nil {
		t.Fatalf(`generateRefreshToken returned jwt with no json serializable header`)
	}

	header.Alg = "none"

	b, _ := json.Marshal(header)

	bodies[0] = base64.StdEncoding.EncodeToString(b)

	modifiedToken := strings.Join(bodies, ".")

	err = verifyRefreshToken(secret, modifiedToken)

	assert.NotNil(t, err)
}

func TestVerifyRefreshTokenNGInvalidIssuedAt(t *testing.T) {
	secret := []byte("foo")

	issuedAt := time.Now().Local().Add(time.Minute * 5).Unix()

	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		IssuedAt: issuedAt,
		Issuer:   jwtIssuer,
		Subject:  jwtRefreshTokenSubject,
	}).SignedString(secret)

	err := verifyRefreshToken(secret, token)

	assert.NotNil(t, err)
}

func TestVerifyRefreshTokenNGInvalidIssuer(t *testing.T) {
	secret := []byte("foo")

	issuedAt := time.Now().Local().Add(time.Minute * -10).Unix()

	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		IssuedAt: issuedAt,
		Issuer:   "rest-cms",
		Subject:  jwtRefreshTokenSubject,
	}).SignedString(secret)

	err := verifyRefreshToken(secret, token)

	assert.NotNil(t, err)
}

func TestVerifyRefreshTokenNGPassedAccessToken(t *testing.T) {
	secret := []byte("foo")

	issuedAt := time.Now().Local().Add(time.Minute * -10).Unix()
	expiresIn := time.Now().Local().Add(time.Minute * 10).Unix()

	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		IssuedAt:  issuedAt,
		ExpiresAt: expiresIn,
		Issuer:    jwtIssuer,
		Subject:   jwtAccessTokenSubject,
	}).SignedString(secret)

	err := verifyRefreshToken(secret, token)

	assert.NotNil(t, err)
}
