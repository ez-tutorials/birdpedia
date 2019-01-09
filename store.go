package main

// sql lib to interact with database
import (
	"database/sql"
)

// create 2 method: get existing birds and add a new bird
type Store interface {
	CreateBird(bird *Bird) error
	GetBirds() ([]*Bird, error)
}

// dbStore implements the "Store" interface
// it takes sql DB connection object that represents the database connection
type dbStore struct {
	db *sql.DB
}

func (store *dbStore) CreateBird(bird *Bird) error {
	// 'Bird' is a simple struct which has "species" and "description"
	// attributes
	// The first underscore means that we don't care about what's returned
	// from this insert query. We just want to know if it was inserted
	// correctly, and the error will be populated if it wasn't
	_, err := store.db.Query("INSERT INTO birds(species, description) VALUES ($1, $2)", bird.Species, bird.Description)
	return err
}

func (store *dbStore) GetBirds() ([]*Bird, error) {
	// qery database for all birds, and return the result to the `row` object
	rows, err := store.db.Query("SELECT species, description from birds")
	// return incase of error, and defer the closing of the row structure
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// create the data structure that is returned from the function.
	// default to empty array of birds
	birds := []*Bird{}
	for rows.Next() {
		// for each row returned by the table, create a pointer to a bird
		bird := &Bird{}
		// populate the `Species` and `Description` attribuites
		if err := rows.Scan(&bird.Species, &bird.Description); err != nil {
			return nil, err
		}
		// append the result to the returned array and repeat for the next row
		birds = append(birds, bird)
	}
	return birds, nil
}

// the store variable is a package level variable that will be available for
// use throughout the application
var store Store

/*
We will need to call the InitStore method to initialize the store. This will
typically be done at the beginning of our application (in this case, when the server starts up)
This can also be used to set up the store as a mock, which we will be observing
later on
*/
func InitStore(s Store) {
	store = s
}
