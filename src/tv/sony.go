package tv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type SonyStatusPayload struct {
	Method  string   `json:"method"`
	Id      int      `json:"id"`
	Params  []string `json:"params"`
	Version string   `json:"version"`
}

type SonyStatusResponse struct {
	Id     int          `json:"id"`
	Result []SonyResult `json:"result"`
}

type SonyResult struct {
	Status string `json:"status"`
}

type SonyTv struct {
	Protocol      string
	Address       string
	Psk           string
	SystemPath    string
	getStatusUrl  string
	StatusPayload SonyStatusPayload
}

func New(protocol string, address string, psk string) *SonyTv {
	systemPath := "/sony/system"

	_, err := url.Parse(fmt.Sprintf("%s://%s", protocol, address))
	if err != nil {
		log.Fatalln("Invalid SonyTV URL")
	}

	return &SonyTv{
		Protocol:     protocol,
		Address:      address,
		Psk:          psk,
		SystemPath:   systemPath,
		getStatusUrl: fmt.Sprintf("%s://%s%s", protocol, address, systemPath),
		StatusPayload: SonyStatusPayload{
			Method:  "getPowerStatus",
			Id:      50,
			Params:  []string{},
			Version: "1.0",
		},
	}
}

func (t *SonyTv) GetStatus() Status {
	jsonData, err := json.Marshal(t.StatusPayload)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(t.getStatusUrl, "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		log.Println(err)
		return Off
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return Off
	}

	var response SonyStatusResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println(err)
		return Off
	}

	status := response.Result[0].Status

	if status == "standby" {
		return StandBy
	} else if status == "active" {
		return Active
	} else {
		return Off
	}
}

func (t *SonyTv) GetAddress() string {
	return t.Address
}
