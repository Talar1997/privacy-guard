package sony

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type TvStatusPayload struct {
	Method  string   `json:"method"`
	Id      int      `json:"id"`
	Params  []string `json:"params"`
	Version string   `json:"version"`
}

type TvStatusResponse struct {
	Id     int        `json:"id"`
	Result []TvResult `json:"result"`
}

type TvResult struct {
	Status string `json:"status"`
}

type Tv struct {
	Protocol      string
	Address       string
	Psk           string
	SystemPath    string
	getStatusUrl  string
	StatusPayload TvStatusPayload
}

func New(protocol string, address string, psk string) *Tv {
	systemPath := "/sony/system"

	return &Tv{
		Protocol:     protocol,
		Address:      address,
		Psk:          psk,
		SystemPath:   systemPath,
		getStatusUrl: fmt.Sprintf("%s://%s%s", protocol, address, systemPath), // TODO: url.Parse
		StatusPayload: TvStatusPayload{
			Method:  "getPowerStatus",
			Id:      50,
			Params:  []string{},
			Version: "1.0",
		},
	}
}

func (t *Tv) GetStatus() string {
	jsonData, err := json.Marshal(t.StatusPayload)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(t.getStatusUrl, "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var response TvStatusResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalln(err)
	}

	return response.Result[0].Status
}
