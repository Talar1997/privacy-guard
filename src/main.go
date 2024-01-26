package main

import (
	"log"
	"os"
	"time"

	"privacy-guard/src/adguard"
	"privacy-guard/src/sony"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Privacy Guard started")
	tvProtocol, tvAddress, tvPsk := getEnv("TV_PROTOCOL"), getEnv("TV_ADDRESS"), getEnv("TV_PSK")
	adguardProtocol, adguardAddress, adguardUsername, adguardPassword :=
		getEnv("ADGUARD_PROTOCOL"), getEnv("ADGUARD_ADDRESS"), getEnv("ADGUARD_USERNAME"), getEnv("ADGUARD_PASSWORD")
	// interval := getEnv("INTERVAL")

	tv := sony.New(tvProtocol, tvAddress, tvPsk)
	adguard := adguard.New(adguardProtocol, adguardAddress, adguardUsername, adguardPassword)

	initialStatus := tv.GetStatus()
	if initialStatus == "standby" {
		adguard.SetRule(tvAddress)
	} else {
		adguard.RemoveRule(tvAddress)
	}

	previousStatus := initialStatus
	for {
		currentStatus := tv.GetStatus()

		if currentStatus != previousStatus {
			log.Printf("TV status change: %s->%s \n", previousStatus, currentStatus)

			if currentStatus == "standby" {
				adguard.SetRule(tv.Address)
			} else {
				adguard.RemoveRule(tv.Address)
			}

			previousStatus = currentStatus
		}

		time.Sleep(2 * time.Second)
	}
}

func getEnv(key string) string {

	filename := ".env" //default env
	if len(os.Args) > 2 {
		filename = os.Args[1]
	}

	err := godotenv.Load(filename)

	if err != nil {
		log.Fatalf("Error loading %s file", filename)
	}

	return os.Getenv(key)
}
