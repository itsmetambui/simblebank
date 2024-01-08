package db

import (
	"context"
	"testing"
	"time"

	"github.com/ddosify/go-faker/faker"
	"github.com/itsmetambui/simplebank/util"
	"github.com/stretchr/testify/assert"
)

func createRandomUser(t *testing.T) User {
	faker := faker.NewFaker()
	password := faker.RandomPassword()
	hashedPassword, err := util.HashPassword(password)
	assert.NoError(t, err)

	arg := CreateUserParams{
		Username:       faker.RandomUsername(),
		FullName:       faker.RandomPersonFullName(),
		Email:          faker.RandomEmail(),
		HashedPassword: hashedPassword,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	assert.NoError(t, err)

	assert.Equal(t, arg.Username, user.Username)
	assert.Equal(t, arg.FullName, user.FullName)
	assert.Equal(t, arg.Email, user.Email)
	assert.Equal(t, arg.HashedPassword, user.HashedPassword)
	assert.True(t, user.PasswordChangedAt.IsZero())
	assert.NotZero(t, user.CreatedAt)

	return user
}
func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	assert.NoError(t, err)
	assert.NotEmpty(t, user2)

	assert.Equal(t, user1.Username, user2.Username)
	assert.Equal(t, user1.FullName, user2.FullName)
	assert.Equal(t, user1.Email, user2.Email)
	assert.Equal(t, user1.HashedPassword, user2.HashedPassword)
	assert.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	assert.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}
