package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"testing"
	// this package is used to make test suite
	"github.com/stretchr/testify/suite"
)

type StoreSuite struct {
	suite.Suite
	// suite is defined as a struct with store and db as attribuites
	// any variables that are to be shared between tests in a suite should be
	// stored as attribuites of the suite instance.
	store *dbStore
	db    *sql.DB
}

func (s *StoreSuite) SetupSuite() {
	// open database connection and stored as instance variable.
	connString := "dbname=bird_encyclopedia sslmode=disable"
	db, err := sql.Open("postgres", connString)
	if err != nil {
		s.T().Fatal(err)
	}
	s.db = db
	s.store = &dbStore{db: db}
}

func (s *StoreSuite) SetupTest() {
	// delete all entries from the table before each test run.
	_, err := s.db.Query("DELETE FROM birds")
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *StoreSuite) TearDownSuite() {
	// close connection after tests finish running
	s.db.Close()
}

// actual test
func TestStoreSuite(t *testing.T) {
	s := new(StoreSuite)
	suite.Run(t, s)
}

func (s *StoreSuite) TestCreateBird() {
	// create a bird through the store `CreateBird` method
	s.store.CreateBird(&Bird{
		Description: "test description",
		Species:     "test species",
	})

	// query the database for the entry created
	res, err := s.db.Query("SELECT COUNT(*) FROM birds WHERE description='test description' AND species='test species'")
	if err != nil {
		s.T().Fatal(err)
	}

	// get count result
	var count int
	for res.Next() {
		err := res.Scan(&count)
		if err != nil {
			s.T().Error(err)
		}
	}

	// assert there must be one entry with the properties of the bird that was
	// inserted
	if count != 1 {
		s.T().Errorf("incorrect count, wanted 1, got %d", count)
	}
}

func (s *StoreSuite) TestGetBird() {
	// insert a sample bird into the birds table
	_, err := s.db.Query(`INSERT INTO birds (species, description) VALUES('bird', 'description')`)
	if err != nil {
		s.T().Fatal(err)
	}

	// get the list of birds through GetBirds method
	birds, err := s.store.GetBirds()
	if err != nil {
		s.T().Fatal(err)
	}

	// assert the count of birds must be 1
	nBirds := len(birds)
	if nBirds != 1 {
		s.T().Errorf("incorrect count, wanted 1, got %d", nBirds)
	}
}
