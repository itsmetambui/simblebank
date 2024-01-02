package db

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbProvider = "postgres"
	dbSource   = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	testDB, err = sql.Open(dbProvider, dbSource)
	if err != nil {
		panic("cannot connect to db: " + err.Error())
	}
	defer testDB.Close()

	testQueries = New(testDB)

	os.Exit(m.Run())
}
