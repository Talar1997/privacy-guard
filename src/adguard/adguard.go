package adguard

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
)

type Adguard struct {
	Protocol           string
	Address            string
	Credentials        string
	filteringRulesPath string
	setRulesPath       string
	blockingRule       string
}

func New(protocol string, address string, username string, password string) *Adguard {
	adguardCredentials := fmt.Sprintf("%s:%s", username, password)

	return &Adguard{
		Protocol:           protocol,
		Address:            address,
		Credentials:        base64.StdEncoding.EncodeToString([]byte(adguardCredentials)),
		filteringRulesPath: "/control/filtering/status",
		setRulesPath:       "/control/filtering/set_rules",
		blockingRule:       "||*^$client=", // Block every request from client provided after = sign
	}
}

func (a *Adguard) GetRules() []string {
	adguardRulesUrl := fmt.Sprintf("%s://%s%s", a.Protocol, a.Address, a.filteringRulesPath) // TODO: consider using url.Parse and keep it as struct field
	req, err := http.NewRequest("GET", adguardRulesUrl, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", a.Credentials))

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var response RulesResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalln(err)
	}

	return response.UserRules
}

func (a *Adguard) SetRule(tvAddress string) {
	rule := fmt.Sprintf("%s%s", a.blockingRule, tvAddress)

	userRules := a.GetRules()
	if !slices.Contains(userRules, rule) {
		log.Printf("Setting rule: %s \n", rule)
		userRules = append(userRules, rule)
		a.UpdateRules(userRules)
	}
}

func (a *Adguard) RemoveRule(tvAddress string) {
	existingRule := fmt.Sprintf("%s%s", a.blockingRule, tvAddress)

	userRules := a.GetRules()
	if slices.Contains(userRules, existingRule) {
		log.Printf("Removing rule: %s \n", existingRule)

		var updatedRules []string
		for _, rule := range userRules {
			if rule != existingRule {
				updatedRules = append(updatedRules, rule)
			}
		}

		a.UpdateRules(updatedRules)
	}
}

func (a *Adguard) UpdateRules(rules []string) {
	adguardSetRulesUrl := fmt.Sprintf("%s://%s%s", a.Protocol, a.Address, a.setRulesPath)
	payload := &RulesPayload{
		Rules: rules,
	}
	payloadStr, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", adguardSetRulesUrl, bytes.NewBuffer(payloadStr))
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", a.Credentials))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	//TODO: track status
}

type RulesResponse struct {
	UserRules []string `json:"user_rules"`
}

type RulesPayload struct {
	Rules []string `json:"rules"`
}
