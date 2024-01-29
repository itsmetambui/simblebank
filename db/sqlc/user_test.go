package db

import (
	"context"
	"database/sql"
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

func TestUpdateUserOnlyFullName(t *testing.T) {
	oldUser := createRandomUser(t)

	faker := faker.NewFaker()
	newFullName := faker.RandomPersonFullName()

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, updatedUser)
	assert.NotEqual(t, oldUser.FullName, updatedUser.FullName)
	assert.Equal(t, newFullName, updatedUser.FullName)
	assert.Equal(t, oldUser.Username, updatedUser.Username)
	assert.Equal(t, oldUser.Email, updatedUser.Email)
	assert.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := createRandomUser(t)

	faker := faker.NewFaker()
	newEmail := faker.RandomEmail()

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, updatedUser)
	assert.NotEqual(t, oldUser.Email, updatedUser.Email)
	assert.Equal(t, newEmail, updatedUser.Email)
	assert.Equal(t, oldUser.Username, updatedUser.Username)
	assert.Equal(t, oldUser.FullName, updatedUser.FullName)
	assert.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	oldUser := createRandomUser(t)

	faker := faker.NewFaker()
	newPashword := faker.RandomStringWithLength(6)
	newHashedPassword, err := util.HashPassword(newPashword)
	assert.NoError(t, err)

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, updatedUser)
	assert.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	assert.Equal(t, newHashedPassword, updatedUser.HashedPassword)
	assert.Equal(t, oldUser.Username, updatedUser.Username)
	assert.Equal(t, oldUser.FullName, updatedUser.FullName)
	assert.Equal(t, oldUser.Email, updatedUser.Email)
}

func TestUpdateUserAllFields(t *testing.T) {
	oldUser := createRandomUser(t)

	faker := faker.NewFaker()
	newFullName := faker.RandomPersonFullName()
	newEmail := faker.RandomEmail()
	newPashword := faker.RandomStringWithLength(6)
	newHashedPassword, err := util.HashPassword(newPashword)
	assert.NoError(t, err)

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, updatedUser)
	assert.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	assert.NotEqual(t, oldUser.FullName, updatedUser.FullName)
	assert.NotEqual(t, oldUser.Email, updatedUser.Email)
	assert.Equal(t, newHashedPassword, updatedUser.HashedPassword)
	assert.Equal(t, newFullName, updatedUser.FullName)
	assert.Equal(t, newEmail, updatedUser.Email)
}
