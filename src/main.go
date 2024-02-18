package main

import (
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	"privacy-guard/src/blocker"
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
	sleeper := &DefaultSleeper{
		Duration: interval,
		Break:    false,
	}

	Watch(sonyTv, adguard, sleeper)
}

func getEnv(key string) string {

	filename := ".env" //default env
	if len(os.Args) > 2 {
		filename = os.Args[1]
	}

	err := godotenv.Load(filename)

	if err != nil {
		log.Printf("Error loading %s file. \n", filename)
	}

	return os.Getenv(key)
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
