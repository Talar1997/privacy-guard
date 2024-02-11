package main

import (
	"log"
	"os"
	"strconv"

	"privacy-guard/src/blocker"
	"privacy-guard/src/tv"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Privacy Guard started")
	tvProtocol, tvAddress, tvPsk := getEnv("TV_PROTOCOL"), getEnv("TV_ADDRESS"), getEnv("TV_PSK")
	adguardProtocol, adguardAddress := getEnv("ADGUARD_PROTOCOL"), getEnv("ADGUARD_ADDRESS")
	adguardUsername, adguardPassword := getEnv("ADGUARD_USERNAME"), getEnv("ADGUARD_PASSWORD")
	interval, err := strconv.Atoi(getEnv("INTERVAL"))

	if err != nil {
		log.Fatalln("Interval value must be integer, setting defualt value (2 sec)")
		interval = 2
	}

	sonyTv := tv.New(tvProtocol, tvAddress, tvPsk)
	adguard := blocker.New(adguardProtocol, adguardAddress, adguardUsername, adguardPassword)

	Watch(sonyTv, adguard, interval)
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
