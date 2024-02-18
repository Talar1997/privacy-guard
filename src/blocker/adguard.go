package blocker

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"slices"
)

type Adguard struct {
	Url         *url.URL
	getRulesUrl *url.URL
	setRulesUrl *url.URL
	Credentials string
}

type UpdateStatus int

const (
	Success UpdateStatus = iota
	Fail
)

type RulesResponse struct {
	UserRules []string `json:"user_rules"`
}

type RulesPayload struct {
	Rules []string `json:"rules"`
}

const getRulesPath = "/control/filtering/status"
const setRulesPath = "/control/filtering/set_rules"
const blockRule = "||*^$client="

// https://github.com/AdguardTeam/AdGuardHome/tree/master/openapi
func NewAdguard(u *url.URL, username string, password string) *Adguard {
	adguardCredentials := fmt.Sprintf("%s:%s", username, password)

	return &Adguard{
		Url:         u,
		getRulesUrl: u.JoinPath(getRulesPath),
		setRulesUrl: u.JoinPath(setRulesPath),
		Credentials: base64.StdEncoding.EncodeToString([]byte(adguardCredentials)),
	}
}

func (a *Adguard) SetRule(tvAddress string) {
	rule := fmt.Sprintf("%s%s", blockRule, tvAddress)

	userRules, err := a.getRules()
	if err != nil {
		log.Println("Couldn't set rule", err)
		return
	}

	if !slices.Contains(userRules, rule) {
		log.Printf("Setting rule: %s \n", rule)
		userRules = append(userRules, rule)
		a.updateRules(userRules)
	}
}

func (a *Adguard) RemoveRule(tvAddress string) {
	existingRule := fmt.Sprintf("%s%s", blockRule, tvAddress)

	userRules, err := a.getRules()
	if err != nil {
		log.Println("Couldn't remove rule", err)
		return
	}

	if slices.Contains(userRules, existingRule) {
		log.Printf("Removing rule: %s \n", existingRule)

		var updatedRules []string
		for _, rule := range userRules {
			if rule != existingRule {
				updatedRules = append(updatedRules, rule)
			}
		}

		a.updateRules(updatedRules)
	}
}

func (a *Adguard) getRules() ([]string, error) {
	req, err := http.NewRequest("GET", a.getRulesUrl.String(), nil)
	if err != nil {
		return []string{}, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", a.Credentials))

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []string{}, err
	}

	var response RulesResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return []string{}, err
	}

	return response.UserRules, nil
}

func (a *Adguard) updateRules(rules []string) (UpdateStatus, error) {
	payload := &RulesPayload{
		Rules: rules,
	}
	payloadStr, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", a.setRulesUrl.String(), bytes.NewBuffer(payloadStr))
	if err != nil {
		log.Println("Update failed", err)
		return Fail, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", a.Credentials))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println("Update failed", err)
		return Fail, err
	}
	resp.Body.Close()

	if resp.StatusCode == 200 {
		return Success, nil
	} else {
		return Fail, nil
	}
}
