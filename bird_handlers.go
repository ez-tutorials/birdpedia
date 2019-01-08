package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Bird struct {
	Species     string `json:"species"`
	Description string `json:"description"`
}

var birds []Bird

func getBirdHandler(w http.ResponseWriter, r *http.Request) {
	// convert "birds" to json
	birdListBytes, err := json.Marshal(birds)

	// if there is an error, print it to the console, and return a server
	// error response to the user
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// no errors, write the JSON list of birds to the response
	w.Write(birdListBytes)
}

func createBirdHandler(w http.ResponseWriter, r *http.Request) {
	// create a new instance of bird
	bird := Bird{}

	// send all data as HTML form data
	// `PostForm` method of the request, prses the form values
	err := r.ParseForm()
	// in case of error, respond with error to user
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get info about the bird from the form info
	bird.Species = r.Form.Get("species")
	bird.Description = r.Form.Get("description")

	// append the existing list of birds with a new entry
	birds = append(birds, bird)

	// redirect the user to the original HTML page
	http.Redirect(w, r, "/assets/", http.StatusFound)
}
