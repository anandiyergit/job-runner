package server

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Append the completed job
func (jobstore *JobStore) addJobs(job Job) {
	jobstore.job = append(jobstore.job, job)
}

// Print the jobs in writer console
func (jobStore *JobStore) printJobs(w http.ResponseWriter) {
	for _, job := range jobStore.job {
		fmt.Fprintf(w, "Job: %s took %v \n", job.name, job.duration)
	}
}

//ViewJobs Route to voew the finished Jobs
func (jobStore *JobStore) ViewJobs(w http.ResponseWriter, r *http.Request) {
	// Validate the incoming request
	augment(w, r)
	jobStore.printJobs(w)
	return
}

//ExecuteJob Responsible for executing the Builds and publishes
func (jobstore *JobStore) ExecuteJob(w http.ResponseWriter, r *http.Request) {
	// Validate the incoming request
	augment(w, r)
	// Job name is passed by the SCM Trigger
	name := r.URL.Query().Get("name")

	if name == "" {
		log.Println("Url Param 'name' is missing")
		return
	}

	user := r.URL.Query().Get("user")

	if user == "" {
		log.Println("Url param 'user' is missing")
	}
	// Create Job and push the work to do
	job := Job{name, &Scm{user: user}, time.Second}
	job.startPipeline(w, r, jobstore)
	return
}

func (job *Job) startPipeline(w http.ResponseWriter, r *http.Request, jobstore *JobStore) {
	startTime := time.Now()
	fmt.Fprintf(w, "Executing Job: %s\n", job.name)
	time.Sleep(1 * time.Second)
	fmt.Fprintf(w, "Finished executing Job: %s in duration: %s \n", job.name, time.Since(startTime))
	jobstore.addJobs(Job{job.name, job.scm, time.Since(startTime)})
}

func augment(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Authenticating Job \n")
	// I would prefer implementing a JWT based solution to authenticate and verify the incoming requests
	// The JWT will contain the needed claims to be matched with then use those information to pass on to headers.
	fmt.Fprintf(w, "Job Authenticated! \n")
}
