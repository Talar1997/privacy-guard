package guard

import (
	"log"
	"privacy-guard/src/blocker"
	"privacy-guard/src/tv"
	"time"
)

func Watch(t tv.Tv, b blocker.Blocker, s Sleeper) {
	initialStatus := t.GetStatus()
	Init(t, b, initialStatus)
	WatchInLoop(t, b, s, initialStatus)
}

func Init(t tv.Tv, b blocker.Blocker, initialStatus tv.Status) {
	tvAddress := t.GetAddress()

	if initialStatus == tv.StandBy {
		b.SetRule(tvAddress)
	} else {
		b.RemoveRule(tvAddress)
	}
}

func WatchInLoop(t tv.Tv, b blocker.Blocker, s Sleeper, initialStatus tv.Status) {
	tvAddress := t.GetAddress()

	previousStatus := initialStatus
	for {
		currentStatus := t.GetStatus()

		if currentStatus != previousStatus {
			log.Printf("TV status change: %d->%d", previousStatus, currentStatus)

			if currentStatus == tv.StandBy {
				b.SetRule(tvAddress)
			} else {
				b.RemoveRule(tvAddress)
			}

			previousStatus = currentStatus
		}

		if s.Stop() {
			break
		} else {
			s.Sleep()
		}
	}
}

type Sleeper interface {
	Sleep()
	Stop() bool
}

type DefaultSleeper struct {
	Duration int
	Break    bool
}

func (d *DefaultSleeper) Sleep() {
	time.Sleep(time.Duration(d.Duration) * time.Second)
}

func (d *DefaultSleeper) Stop() bool {
	return d.Break
}
