package util

import (
	"testing"

	"github.com/ddosify/go-faker/faker"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	faker := faker.NewFaker()

	password := faker.RandomPassword()
	hashedPassword1, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword1)

	err = CheckPassword(password, hashedPassword1)
	assert.NoError(t, err)

	wrongPassword := faker.RandomPassword()
	err = CheckPassword(wrongPassword, hashedPassword1)
	assert.EqualError(t, bcrypt.ErrMismatchedHashAndPassword, err.Error())

	hashedPassword2, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEqual(t, hashedPassword1, hashedPassword2)
}
