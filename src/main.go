package main

import (
	"log"
	"net/url"
	"os"
	"strconv"

	"privacy-guard/src/blocker"
	"privacy-guard/src/guard"
	"privacy-guard/src/tv"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Privacy Guard started")
	tvUrlStr := getEnv("TV_ADDRESS")
	tvUrl, err := url.Parse(tvUrlStr)
	if err != nil {
		log.Fatalln("Invalid TV URL", err)
	}

	adguardUsername, adguardPassword := getEnv("ADGUARD_USERNAME"), getEnv("ADGUARD_PASSWORD")
	adguardUrlStr := getEnv("ADGUARD_ADDRESS")
	adguardUrl, err := url.Parse(adguardUrlStr)
	if err != nil {
		log.Fatalln("Invalid adguard URL", err)
	}

	intervalStr := getEnv("INTERVAL")
	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		log.Fatalln("Interval value must be integer, setting defualt value (2 sec)")
		interval = 2
	}

	sonyTv := tv.NewSony(tvUrl)
	adguard := blocker.NewAdguard(adguardUrl, adguardUsername, adguardPassword)
	sleeper := &guard.DefaultSleeper{
		Duration: interval,
		Break:    false,
	}

	guard.Watch(sonyTv, adguard, sleeper)
}

func getEnv(key string) string {

	filename := ".env" //default env
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}

	err := godotenv.Load(filename)

	if err != nil {
		log.Printf("Error while reading %s key from file, %s. Error: %s", key, filename, err)
	}

	return os.Getenv(key)
}
