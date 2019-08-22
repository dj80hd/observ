package app

import "time"

// Job is a simulated amount of work that takes 'secs' seconds to complete
type Job struct {
	secs int
	fail bool
}

// Run returns true if ok
func (j Job) Run() bool {
	time.Sleep(time.Duration(j.secs) * time.Second)
	return !j.fail
}
