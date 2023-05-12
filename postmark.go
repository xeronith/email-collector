package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var client = &http.Client{}

func SendEmail(token, from, to, templateAlias string) error {
	emailBody := map[string]interface{}{
		"From":          from,
		"To":            to,
		"TemplateAlias": templateAlias,
		"TemplateModel": map[string]interface{}{},
	}

	requestBody, _ := json.Marshal(emailBody)
	buffer := bytes.NewBuffer(requestBody)

	r, err := http.NewRequest("POST", "https://api.postmarkapp.com/email/withTemplate", buffer)
	if err != nil {
		return err
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Accept", "application/json")
	r.Header.Add("X-Postmark-Server-Token", token)

	resp, err := client.Do(r)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error sending email to %s", to)
	}

	return nil
}
