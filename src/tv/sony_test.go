package tv

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestNewSony(t *testing.T) {
	tvUrl, _ := url.Parse("http://localhost")
	tvSony := NewSony(tvUrl)

	if !strings.Contains(tvSony.statusUrl.String(), systemPath) {
		t.Errorf("Expected url (%s) to contains %s", tvSony.statusUrl.String(), systemPath)
	}
}

func TestGetStatus(t *testing.T) {
	standByServer := makeTvServer("standby")
	tvUrl, _ := url.Parse(standByServer.URL)

	tvSony := NewSony(tvUrl)
	status := tvSony.GetStatus()
	if status != StandBy {
		t.Errorf("Expected %d to be %d", status, StandBy)
	}
	defer standByServer.Close()

	activeServer := makeTvServer("active")
	tvUrl, _ = url.Parse(activeServer.URL)
	tvSony = NewSony(tvUrl)
	status = tvSony.GetStatus()
	if status != Active {
		t.Errorf("Expected %d to be %d", status, StandBy)
	}
	defer activeServer.Close()

	tvUrl, _ = url.Parse("http://localhost:2137")
	tvSony = NewSony(tvUrl)
	status = tvSony.GetStatus()
	if status != Off {
		t.Errorf("Expected %d to be %d", status, Off)
	}
}

func makeTvServer(status string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := &SonyStatusResponse{
			Id: 0,
			Result: []SonyResult{
				{
					Status: status,
				},
			},
		}

		responseData, _ := json.Marshal(response)

		w.Header().Set("Content-Type", "application/json")
		w.Write(responseData)
	}))
}
