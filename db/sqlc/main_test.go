package db

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

const (
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbSource)
	if err != nil {
		panic("cannot connect to db: " + err.Error())
	}
	defer conn.Close(ctx)

	testQueries = New(conn)

	os.Exit(m.Run())
}
