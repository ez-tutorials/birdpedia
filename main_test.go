package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	/*
		form a new HTTP request that's going to be pased to the handler.
		1st argument is the method
		2nd argument is the route, leave it blank for now
		3rd argument is the request body
	*/
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	// use Go's httptest lib to create an http recorder, it acts as the target
	// of the http request.
	recorder := httptest.NewRecorder()

	// create a http handler for the hander function.
	// "handler" is the handler function definer in main.go to be tested.
	hf := http.HandlerFunc(handler)

	// serve the http request to the recorder.
	// this line executes the handler func that needs to be tested.
	hf.ServeHTTP(recorder, req)

	// check the status code
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v instead of %v", status, http.StatusOK)
	}

	// check the response body
	expected := "Hello World!"
	actual := recorder.Body.String()
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v instead of %v", actual, expected)
	}
}

func TestRouter(t *testing.T) {
	// instantiate the router using the constructer defined in main.go
	r := newRouter()

	// create a new server using the httptest ib
	mockServer := httptest.NewServer(r)

	// the mock server runs a server and exposes its location in the URL
	// attribute.
	resp, err := http.Get(mockServer.URL + "/hello")

	// handle any unexpected error
	if err != nil {
		t.Fatal(err)
	}

	// check status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status should be ok, got %d", resp.StatusCode)
	}

	// read response body and convert to string
	defer resp.Body.Close()
	// read the body into bytes
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	// convert byte to string
	respString := string(b)
	expected := "Hello World!"

	// check the response
	if respString != expected {
		t.Errorf("Response should be %s, got %s", expected, respString)
	}
}

func TestRouterForNonExistentRoute(t *testing.T) {
	r := newRouter()
	mockServer := httptest.NewServer(r)
	// make a request to a route that wasn't defined
	resp, err := http.Post(mockServer.URL+"/hello", "", nil)

	if err != nil {
		t.Fatal(err)
	}

	// check status code is 405
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Status should be 405, got %d", resp.StatusCode)
	}

	// check the body is empty
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	respString := string(b)
	expected := ""

	if respString != expected {
		t.Errorf("Response should be %s, got %s", expected, respString)
	}

}

func TestStaticFileServer(t *testing.T) {
	r := newRouter()
	mockServer := httptest.NewServer(r)

	// we want to hit the "GET /assets/" route to get the index.html
	resp, err := http.Get(mockServer.URL + "/assets/")
	if err != nil {
		t.Fatal(err)
	}

	// status code needs to be 200
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status should be 200, got %d", resp.StatusCode)
	}

	// test that the content-type header is "text/html; charset=utf-8"
	contentType := resp.Header.Get("Content-Type")
	expectedContentType := "text/html; charset=utf-8"

	if expectedContentType != contentType {
		t.Errorf("Wrong content type, expected %s, got %s", expectedContentType, contentType)
	}
}
