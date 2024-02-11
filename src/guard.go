package main

import (
	"log"
	"privacy-guard/src/blocker"
	"privacy-guard/src/tv"
	"time"
)

func Watch(t tv.Tv, b blocker.Blocker, interval int) {
	tvAddress := t.GetAddress()
	initialStatus := t.GetStatus()

	if initialStatus == tv.StandBy {
		b.SetRule(tvAddress)
	} else {
		b.RemoveRule(tvAddress)
	}

	previousStatus := initialStatus
	for {
		currentStatus := t.GetStatus()

		if currentStatus != previousStatus {
			log.Printf("TV status change: %d->%d \n", previousStatus, currentStatus)

			if currentStatus == tv.StandBy {
				b.SetRule(tvAddress)
			} else {
				b.RemoveRule(tvAddress)
			}

			previousStatus = currentStatus
		}

		time.Sleep((time.Duration(interval) * time.Second))
	}
}
