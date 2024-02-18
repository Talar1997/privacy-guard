package tv

import (
	"bytes"
	"encoding/json"
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
	Url           *url.URL
	statusUrl     *url.URL
	StatusPayload SonyStatusPayload
}

const systemPath = "/sony/system"
const method = "getPowerStatus"
const id = 50
const version = "1.0"

// https://pro-bravia.sony.net/develop/index.html
func NewSony(u *url.URL) *SonyTv {
	statusUrl := u.JoinPath(systemPath)

	sonyTv := &SonyTv{
		Url:       u,
		statusUrl: statusUrl,
		StatusPayload: SonyStatusPayload{
			Method:  method,
			Id:      id,
			Params:  []string{},
			Version: version,
		},
	}

	return sonyTv
}

func (t *SonyTv) GetStatus() Status {
	jsonData, err := json.Marshal(t.StatusPayload)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(t.statusUrl.String(), "application/json", bytes.NewBuffer(jsonData))

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

	switch status := response.Result[0].Status; status {
	case "standby":
		return StandBy
	case "active":
		return Active
	default:
		return Off
	}
}

func (t *SonyTv) GetAddress() string {
	return t.Url.Host
}
