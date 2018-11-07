package sche

import (
	"reflect"
	"time"
)

type job struct {
	id      string
	handler reflect.Value
	in      []reflect.Value

	sche       *sche
	mode       int // mode: period or trigger
	frequency  time.Duration
	unit       time.Duration
	times      int          // mark if using DoTimes(), unspecified default
	expiration bool         // mark if using DoUntil(), false default
	until      time.Time    // expiration time
	at         bool         // mark if using At()
	h, m, s    int          // At() parameters
	point      time.Time    // next time
	week       time.Weekday // week flag, unspecified default
}

// unit defined here
const (
	period      = 0 // default
	trigger     = 1
	unspecified = -1
)

var now = func() time.Time {
	return time.Now()
}

func newJob() *job {
	return &job{
		unit:  1,
		mode:  unspecified,
		week:  unspecified,
		times: unspecified,
	}
}

// mode : unspecified --> period
func (j *job) Every(f ...time.Duration) *job {
	if j.mode != unspecified {
		panic("mode is set")
	}
	l := len(f)
	j.mode = period
	switch l {
	case 0:
		j.frequency = 1
	case 1:
		if f[0] <= 0 {
			panic("need a non-zero integer")
		}
		j.frequency = f[0]
	default:
		panic("need a non-zero integer")
	}
	return j
}

// mode : unspecified --> period
func (j *job) EverySingle() *job {
	return j.Every()
}

// mode : unspecified --> trigger
func (j *job) After(d time.Duration) *job {
	if j.mode != unspecified {
		panic("mode is set")
	}
	j.mode = trigger
	j.frequency = d
	return j
}

// mode : unspecified --> trigger or period --> period
func (j *job) At(h, m, s int) *job {
	if j.mode == trigger {
		panic("error context")
	}
	if j.mode == unspecified {
		j.mode = trigger
	}
	if j.mode == period && j.unit < 24*time.Hour {
		panic("error context")
	}
	j.at = true
	if h < 0 || h > 24 || m < 0 || m > 60 || s < 0 || s > 60 {
		panic("error context")
	}
	j.h, j.m, j.s = h, m, s
	return j
}

// mode : not unspecified and keep previous status
func (j *job) Do(f interface{}, a ...interface{}) *job {
	if j.mode == unspecified {
		panic("error context")
	}
	j.check(f, a)
	j.sche.jobs = append(j.sche.jobs, j)
	return j
}

func (j *job) ForTimes(c int) *job {
	if j.mode != period {
		panic("error context")
	}
	j.times = c
	return j
}

func (j *job) Until(e time.Time) *job {
	if j.mode != period {
		panic("error context")
	}
	j.until = e
	return j
}

func (j *job) Second() *job {
	if j.mode != period {
		panic("error context")
	}
	j.unit = time.Second
	return j
}
func (j *job) Minute() *job {
	if j.mode != period {
		panic("error context")
	}
	j.unit = time.Minute
	return j
}
func (j *job) Hour() *job {
	if j.mode != period {
		panic("error context")
	}
	j.unit = time.Hour
	return j
}
func (j *job) Day() *job {
	if j.mode != period {
		panic("error context")
	}
	j.unit = time.Hour * 24
	return j
}
func (j *job) Week() *job {
	if j.mode != period {
		panic("error context")
	}
	j.unit = 7 * 24 * time.Hour
	return j
}
func (j *job) Sunday() *job {
	if j.mode != period {
		panic("error context")
	}
	j.unit = 7 * 24 * time.Hour
	j.week = time.Sunday
	return j
}
func (j *job) Monday() *job {
	if j.mode != period {
		panic("error context")
	}
	j.unit = 7 * 24 * time.Hour
	j.week = time.Monday
	return j
}
func (j *job) Tuesday() *job {
	if j.mode != period {
		panic("error context")
	}
	j.unit = 7 * 24 * time.Hour
	j.week = time.Tuesday
	return j
}
func (j *job) Wednesday() *job {
	if j.mode != period {
		panic("error context")
	}
	j.unit = 7 * 24 * time.Hour
	j.week = time.Wednesday
	return j
}
func (j *job) Thursday() *job {
	if j.mode != period {
		panic("error context")
	}
	j.unit = 7 * 24 * time.Hour
	j.week = time.Thursday
	return j
}
func (j *job) Friday() *job {
	if j.mode != period {
		panic("error context")
	}
	j.unit = 7 * 24 * time.Hour
	j.week = time.Friday
	return j
}
func (j *job) Saturday() *job {
	if j.mode != period {
		panic("error context")
	}
	j.unit = 7 * 24 * time.Hour
	j.week = time.Saturday
	return j
}

func (j *job) check(f interface{}, args []interface{}) {
	if f == nil {
		panic("handler type error")
	}
	j.handler = reflect.ValueOf(f)
	if j.handler.Kind() != reflect.Func {
		panic("handler type error")
	}
	ht := j.handler.Type()
	if ht.NumIn() != len(args) || ht.NumOut() != 0 {
		panic("arguments count not match")
	}
	j.in = make([]reflect.Value, 0)
	ft := reflect.TypeOf(f)
	for i, n := 0, ht.NumIn(); i < n; i++ {
		if ft.In(i) != reflect.TypeOf(args[i]) {
			panic("arguments type not match")
		}
		j.in = append(j.in, reflect.ValueOf(args[i]))
	}
}

func (j *job) next() {
	j.point = now().Add(j.unit * j.frequency)
}

func (j *job) due() bool {
	if j.point == (time.Time{}) {
		j.first()
		return false
	}
	if j.mode == period {
		j.next()
		if j.expiration == true && j.until.After(j.point) {
			return false
		}
		if j.times == unspecified {
			return false
		}
		if j.times > 0 {
			j.times--
			return false
		}
		return true
	}
	return true
}

func (j *job) first() {
	if j.mode == trigger {
		if j.at {
			j.point = nextAtPoint(j.h, j.m, j.s)
		} else {
			j.point = now().Add(j.frequency)
		}
	} else {
		if j.at {
			if j.unit < time.Hour*24 {
				panic("error context")
			}
			if j.week == unspecified {
				j.point = nextAtPoint(j.h, j.m, j.s)
			} else {
				j.point = nextAtPointWithWeek(j.h, j.m, j.s, j.week)
			}
		} else {
			j.point = now().Add(j.unit * j.frequency)
		}

		if j.times > 0 {
			j.times--
		}
	}
}

func nextAtPoint(h, m, s int) time.Time {
	ts := now()
	point := time.Date(ts.Year(), ts.Month(), ts.Day(),
		h, m, s, 0, time.Local)
	if ts.After(point) {
		point.Add(24 * time.Hour)
	}
	return point
}

func nextAtPointWithWeek(h, m, s int, w time.Weekday) time.Time {
	point := nextAtPoint(h, m, s)
	wp := point.Weekday()
	for wp != w {
		wp = point.Add(24 * time.Hour).Weekday()
	}
	return point
}
