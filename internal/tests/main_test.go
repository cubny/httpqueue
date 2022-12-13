package tests_test

import (
	"flag"
	"github.com/cubny/cart/internal/infra/http/api"
	"github.com/cubny/cart/internal/tests/testdb"
	"os"
	"testing"

	"github.com/cubny/cart/internal/service"
	"github.com/cubny/cart/internal/storage/sqlite3"

	log "github.com/sirupsen/logrus"
)

var (
	// "a" is a test global handler. it is used in integration tests so that each test does not
	// setup it's own database and other requirements for the handler
	a *api.Router

	// testDB is a helper to manipulate the database such as migrations, seeding, etc.
	testDB *testdb.TestDB
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	if !testing.Short() { // if it's an integration test, then setup everything
		log.Info("setting up test db")
		db, err := setupDatabase()
		if err != nil {
			log.WithError(err).Info("cannot setup test database")
			return 1
		}
		defer db.Close()

		service, err := service.New(db)
		if err != nil {
			log.WithError(err).Info("cannot instantiate cart service")
			return 1
		}

		a, err = api.New(service)
		if err != nil {
			log.WithError(err).Infof("cannot instantiate handler, %s", err)
			return 1
		}

		testDB = testdb.New(db, service)
		if err := testDB.Refresh(); err != nil {
			log.WithError(err).Infof("cannot refresh db, %s", err)
			return 1
		}
	}
	return m.Run()
}

func setupDatabase() (*sqlite3.Sqlite3, error) {
	dbfile, err := os.Create("test.db")
	if err != nil {
		return nil, err
	}
	db, err := sqlite3.New(dbfile)
	if err != nil {
		return nil, err
	}

	return db, nil
}
