package server

import "time"

type Scm struct {
	user string
}

type Job struct {
	name     string
	scm      *Scm
	duration time.Duration
}

type JobStore struct {
	job []Job
}
