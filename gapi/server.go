package gapi

import (
	"fmt"

	db "github.com/itsmetambui/simplebank/db/sqlc"
	"github.com/itsmetambui/simplebank/pb"
	"github.com/itsmetambui/simplebank/token"
	"github.com/itsmetambui/simplebank/util"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
