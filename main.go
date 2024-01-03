package main

import (
	"database/sql"

	"github.com/itsmetambui/simplebank/api"
	db "github.com/itsmetambui/simplebank/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbProvider    = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "localhost:8080"
)

func main() {
	conn, err := sql.Open(dbProvider, dbSource)
	if err != nil {
		panic("cannot connect to db: " + err.Error())
	}
	defer conn.Close()

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		panic("cannot start server: " + err.Error())
	}
}
