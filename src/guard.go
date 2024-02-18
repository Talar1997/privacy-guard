package main

import (
	"log"
	"privacy-guard/src/blocker"
	"privacy-guard/src/tv"
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
