package api

import (
	"os"
	"testing"
	"time"

	"github.com/ddosify/go-faker/faker"
	"github.com/gin-gonic/gin"
	db "github.com/itsmetambui/simplebank/db/sqlc"
	"github.com/itsmetambui/simplebank/util"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func newTestServer(t *testing.T, store db.Store) *Server {
	faker := faker.NewFaker()
	config := util.Config{
		TokenSymmetricKey:   faker.RandomStringWithLength(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	assert.NoError(t, err)

	return server
}
