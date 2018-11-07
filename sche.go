package sche

import (
	"github.com/google/uuid"
	"runtime"
	"time"
)

/*
usage:
Sche.Push().At(9,10,0).Do(...) 							// do... at the next nearest 9:10
Sche.Push().After(10*time.Second).Do(...) 				// do... at 10 seconds later
Sche.Push().Every(30).Second().Do(...)					// do... each 30 second
Sche.Push().Every(30).Minute().Do(...)					// do... each 30 minute
Sche.Push().Every(3).Hour().Do(...)						// do... each 3 hour
Sche.Push().Every().Day().Do(...)						// do... everyday on this moment
Sche.Push().Every().Week().Do(...)						// do... each week
Sche.Push().Every(2).Monday().Do(...) 					// do... every 2 Monday
Sche.Push().Every().Day().At(10,10,0).Do(...)			// do... at 10:10 everyday
Sche.Push().Every(3).Day().At(10,10,0).Do(...)			// do... at 10:10 every 3 day
Sche.Push().Every(3).Week().At(10,10,0).Do(...)			// do... at 10:10 every 3 week
Sche.Push().Every(3).Monday().At(10,10,0).Do(...)		// do... at 10:10 on every 3 Monday
Sche.Push().Every(3).Monday().DoTimes(3, ...)			// do... ... 3 times
Sche.Push().Every(3).Monday().DoUntil(endpoint, ...)	// do... ... until endpoint

todos:
support cron expression
*/

type sche struct {
	id   string // if many schedulers ...
	jobs []*job
}

var Sche = sche{"", make([]*job, 0)} // single instance

func (s *sche) Push() *job {
	j := newJob()
	j.sche = s
	j.id = uuid.New().String()
	return j
}

func (s *sche) Run() {
	for _, j := range s.jobs {
		go func(j *job) {
			for !j.due() {
				d := time.Until(j.point)
				select {
				case <-time.After(d):
					j.handler.Call(j.in)
				}
			}
		}(j)
	}
	runtime.Goexit()
}
