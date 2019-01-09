/*
	This is the name of the package
	Everything with this package anme can see everything
	else inside the same package, regardless of the file they are in.
*/
package main

// These are the libraries we are going to use
// Both "fmt" and "net" are part of the Go standard library
import (
	"database/sql"
	"fmt" // formatted I/O operations
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http" // implement HTTP client and server
)

func newRouter() *mux.Router {
	connString := "dbname=bird_encyclopedia sslmode=disable"
	db, err := sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	InitStore(&dbStore{db: db})
	r := mux.NewRouter()
	r.HandleFunc("/hello", handler).Methods("GET")
	// declare static file directory
	staticFileDirectory := http.Dir("./assets/")

	/*
		Declare the handler, that routes requests to their respective filename.
		The fileserver is wrapped in the 'stripPrefix' method, because we want
		to remove the "/assets/" prefix when looking for files.
		if we type "/assets/index.html", the file serve will look for "index.html"
	*/
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
	// the "PathPrefix" matches all routes starting with "/assets/"
	r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")

	r.HandleFunc("/bird", getBirdHandler).Methods("GET")
	r.HandleFunc("/bird", createBirdHandler).Methods("POST")

	return r
}

func main() {

	r := newRouter()

	http.ListenAndServe(":8080", r)

}

func handler(w http.ResponseWriter, r *http.Request) {
	// Fprintf takes a "writer" as its first arguments.
	// The 2nd argument is the data that is piped into this writer.
	// The output appears according to where the writer moves it. In this case
	// the ResponseWriter w writes the output as the response to the users
	// request.
	fmt.Fprintf(w, "Hello World!")
}
