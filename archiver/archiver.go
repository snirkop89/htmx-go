// Package archiver mimics the usage of archiving data in a long-process.
// It can update its status and run asynchrouncly.
package archiver

import (
	"math/rand"
	"time"
)

type Archiver struct {
	status   string
	progress float64
}

func New() *Archiver {
	return &Archiver{
		status:   "Waiting",
		progress: 0,
	}
}

// Starts a new job if its status is "Waiting"
func (a *Archiver) Run() {
	completed := make(chan struct{})
	prog := make(chan float64)

	a.status = "Running"
	go func() {
		for i := 0; i < 10; i++ {
			r := rand.Intn(3)
			time.Sleep(time.Second * time.Duration(r))
			if a.status != "Running" {
				continue
			}
			prog <- (1 + float64(i)) / 10
		}
		time.Sleep(time.Second)
		if a.status != "Running" {
			return
		}
		completed <- struct{}{}
	}()

	go func() {
		for {
			select {
			case <-completed:
				a.status = "Complete"
				return
			case val := <-prog:
				a.progress = val
			}
		}
	}()
}

// Status returns the currect state of the download:
// Waiting, Running, Complete
func (a *Archiver) Status() string {
	return a.status
}

// Progress is a number between 0 and 1 incidating how much
// progress the archve job has made.
func (a *Archiver) Progress() float64 {
	return a.progress
}

// File returns the path of the file that has been created
func (a *Archiver) File() string {
	return "contacts.json"
}

func (a *Archiver) Reset() {
	a.status = "Waiting"
	a.progress = 0
}
