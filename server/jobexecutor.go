package server

import (
	"encoding/json"
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
func (jobstore *JobStore) printJobs(w http.ResponseWriter) {
	for _, job := range jobstore.job {
		fmt.Fprintf(w, "Job: %s triggered by %s took %v \n", job.name, job.scm.user, job.duration)
	}
}

//ViewJobs Route to voew the finished Jobs
func (jobstore *JobStore) ViewJobs(w http.ResponseWriter, r *http.Request) {
	// Validate the incoming request
	augment(w, r)
	jobstore.printJobs(w)
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
		//w.WriteHeader(500)
		writeJSON(w, http.StatusInternalServerError, "Url Param 'name' is missing")
		return
	}

	user := r.URL.Query().Get("user")

	if user == "" {
		log.Println("Url param 'user' is missing")
		writeJSON(w, http.StatusInternalServerError, "Url Param 'user' is missing")
		return
	}

	// Same user cannot execute a job under 20 seconds
	if jobstore.controlledUserCheck(user, 20) {
		fmt.Fprintf(w, "Already executing Job for user: %s\n", user)
		return
	}

	// Create Job and push the work to do
	job := Job{name, &Scm{user: user}, time.Second, time.Now().UTC()}
	job.startPipeline(w, r, jobstore)
	return
}

// startPipeline allows to run a particular job and also adds the job into the jobstore.
func (job *Job) startPipeline(w http.ResponseWriter, r *http.Request, jobstore *JobStore) {
	startTime := time.Now()
	fmt.Fprintf(w, "Executing Job: %s\n", job.name)
	time.Sleep(1 * time.Second)
	fmt.Fprintf(w, "Finished executing Job: %s in duration: %s \n", job.name, time.Since(startTime))
	jobstore.addJobs(Job{job.name, job.scm, time.Since(startTime), time.Now().UTC()})
}

// The resolution needed for Authentication and Authorization
func augment(w http.ResponseWriter, r *http.Request) {
	log.Println("Authenticating Job")
	// I would prefer implementing a JWT based solution to authenticate and verify the incoming requests
	// The JWT will contain the needed claims to be matched with then use those information to pass on to headers.
	log.Println("Job Authenticated!")
}

/* Responsible for calculating the difference between two jobs scheduled by the same user.
   It uses the methodology of calculating the current time and allowing the cooling/delay time to pass over before
   starting/executing the new job.
*/
func (jobstore *JobStore) controlledUserCheck(user string, delay int) bool {
	for _, job := range jobstore.job {
		if user == job.scm.user {
			// Taking the current time for reference
			t := time.Now().UTC()
			// Computing the delay time for which the scheduling would be prevented
			delayTime := t.Add(-time.Second * 20)

			// calculate the diff
			diffTime := delayTime.Sub(job.time).Seconds()
			if diffTime < 20 {
				return true
			}
		}
	}
	return false
}

// write the JSON response to the http writer
func writeJSON(w http.ResponseWriter, code int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(value)
}
