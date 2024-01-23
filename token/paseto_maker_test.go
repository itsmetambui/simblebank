package token

import (
	"testing"
	"time"

	"github.com/ddosify/go-faker/faker"
	"github.com/stretchr/testify/assert"
)

func TestPasetoMaker(t *testing.T) {
	faker := faker.NewFaker()

	maker, err := NewPasetoMaker(faker.RandomStringWithLength(32))
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

func TestExpiredPasetoToken(t *testing.T) {
	faker := faker.NewFaker()

	maker, err := NewPasetoMaker(faker.RandomStringWithLength(32))
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
