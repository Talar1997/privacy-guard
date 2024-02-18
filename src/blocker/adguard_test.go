package blocker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestNewAdguard(t *testing.T) {
	url, _ := url.Parse("http://localhost")
	username, password := "admin", "password"
	a := NewAdguard(url, username, password)

	if !strings.Contains(a.getRulesUrl.String(), getRulesPath) {
		t.Errorf("Expected url (%s) to contains %s", a.getRulesUrl.String(), getRulesPath)
	}
}

func TestSetRule(t *testing.T) {
	initialRules := []string{
		"rule1",
		"rule2",
	}
	newRule := "rule3"

	server := makeAdguardServer(initialRules)
	defer server.Close()
	serverUrl, _ := url.Parse(server.URL)

	adguard := NewAdguard(serverUrl, "admin", "password")
	actualRules, _ := adguard.getRules()
	if !reflect.DeepEqual(initialRules, actualRules) {
		t.Errorf("Expected %s to be %s", initialRules, actualRules)
	}

	adguard.SetRule(newRule)
	actualRules, _ = adguard.getRules()
	wantRules := []string{
		"rule1",
		"rule2",
		fmt.Sprintf("%s%s", blockRule, newRule), //||*^$client=rule3"
	}
	if !reflect.DeepEqual(wantRules, actualRules) {
		t.Errorf("Expected %s to be %s", actualRules, wantRules)
	}

	// Add same rule again to check if it's not duplicated
	adguard.SetRule(newRule)
	actualRules, _ = adguard.getRules()
	if !reflect.DeepEqual(wantRules, actualRules) {
		t.Errorf("Expected %s to be %s", actualRules, wantRules)
	}
}

func TestRemoveRule(t *testing.T) {
	oldRule := "rule3"
	initialRules := []string{
		"rule1",
		"rule2",
		fmt.Sprintf("%s%s", blockRule, oldRule), //||*^$client=rule3"
	}

	server := makeAdguardServer(initialRules)
	defer server.Close()
	serverUrl, _ := url.Parse(server.URL)

	adguard := NewAdguard(serverUrl, "admin", "password")
	adguard.RemoveRule(oldRule)
	actualRules, _ := adguard.getRules()

	wantRules := []string{
		"rule1",
		"rule2",
	}
	if !reflect.DeepEqual(wantRules, actualRules) {
		t.Errorf("Expected %s to be %s", actualRules, wantRules)
	}

	// Remove same rule
	adguard.RemoveRule(oldRule)
	actualRules, _ = adguard.getRules()
	if !reflect.DeepEqual(wantRules, actualRules) {
		t.Errorf("Expected %s to be %s", actualRules, wantRules)
	}
}

func makeAdguardServer(rules []string) *httptest.Server {
	var rulesStorage []string = rules

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case getRulesPath:
			rulesResponse := &RulesResponse{
				rulesStorage,
			}

			response, _ := json.Marshal(rulesResponse)
			w.Header().Set("Content-Type", "application/json")
			w.Write(response)
		case setRulesPath:
			var ru RulesPayload
			json.NewDecoder(r.Body).Decode(&ru)
			rulesStorage = ru.Rules
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not found"))
		}
	}))
}
