package db

import (
	"database/sql"
	"os"
	"testing"

	"github.com/itsmetambui/simplebank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		panic("cannot load config: " + err.Error())
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		panic("cannot connect to db: " + err.Error())
	}
	defer testDB.Close()

	testQueries = New(testDB)

	os.Exit(m.Run())
}
