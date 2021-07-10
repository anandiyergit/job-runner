package server

import "time"

type Scm struct {
	user string
}

type Job struct {
	name     string
	scm      *Scm
	duration time.Duration
	time     time.Time
}

type JobStore struct {
	job []Job
}
