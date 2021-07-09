package main

import (
	_ "expvar"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/anandiyergit/job-runner/server"
)

func main() {
	// handler for jobStore
	handlers := server.JobStore{}

	// Defining the routes

	// For simplicity to use GET /executeJob. In ideal world this should be a post request but just for the sake of ease of test.
	http.HandleFunc("/executeJob", handlers.ExecuteJob)

	// Route for Viewing the Finished Jobs
	http.HandleFunc("/viewJobs", handlers.ViewJobs)

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
