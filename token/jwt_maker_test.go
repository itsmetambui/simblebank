package token

import (
	"testing"
	"time"

	"github.com/ddosify/go-faker/faker"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJWTMaker(t *testing.T) {
	faker := faker.NewFaker()

	maker, err := NewJWTMaker(faker.RandomStringWithLength(32))
	assert.NoError(t, err)

	username := faker.RandomUsername()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(username, duration)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotEmpty(t, payload)
	assert.Equal(t, username, payload.Username)

	payload, err = maker.VerifyToken(token)
	assert.NoError(t, err)
	assert.NotEmpty(t, payload)

	assert.Equal(t, username, payload.Username)
	assert.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	assert.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJWT(t *testing.T) {
	faker := faker.NewFaker()

	maker, err := NewJWTMaker(faker.RandomStringWithLength(32))
	assert.NoError(t, err)

	username := faker.RandomUsername()

	token, payload, err := maker.CreateToken(username, -time.Minute)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotEmpty(t, payload)
	assert.Equal(t, username, payload.Username)

	payload, err = maker.VerifyToken(token)
	assert.Error(t, err, ErrExpiredToken.Error())
	assert.Nil(t, payload)
}

func TestInvalidJWTAlgNone(t *testing.T) {
	faker := faker.NewFaker()

	payload, err := NewPayload(faker.RandomUsername(), time.Minute)
	assert.NoError(t, err)

	token := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	signedToken, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	assert.NoError(t, err)

	maker, err := NewJWTMaker(faker.RandomStringWithLength(32))
	assert.NoError(t, err)

	payload, err = maker.VerifyToken(signedToken)
	assert.Error(t, err, ErrInvalidToken.Error())
	assert.Nil(t, payload)
}
