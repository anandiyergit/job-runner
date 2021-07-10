package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestExecuteJob(testing *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/executeJob?name=frontJob&user=bob", nil)
	if err != nil {
		testing.Fatal(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handlers := JobStore{}
	handler := http.HandlerFunc(handlers.ExecuteJob)
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		testing.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Check the response body is what we expect.
	testing.Log("code", rr.Code)
	//testing.Log("status",status)
	expected := "frontJob"
	testing.Logf(rr.Body.String())
	if !strings.Contains(rr.Body.String(), expected) {
		testing.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// Similar Test cases to be designed for Failure scenarios for User and name parameters with 500
