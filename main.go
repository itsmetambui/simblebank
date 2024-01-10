package main

import (
	"database/sql"

	"github.com/itsmetambui/simplebank/api"
	db "github.com/itsmetambui/simplebank/db/sqlc"
	"github.com/itsmetambui/simplebank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		panic("cannot load config: " + err.Error())
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		panic("cannot connect to db: " + err.Error())
	}
	defer conn.Close()

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		panic("cannot create server: " + err.Error())
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		panic("cannot start server: " + err.Error())
	}
}
