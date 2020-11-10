package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func mustLogin(email, password string) (int, []*http.Cookie) {
	dataJSON := fmt.Sprintf(`{"email":"%s","password":"%s","remember":true}`, email, password)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "https://www.endomondo.com/rest/session", strings.NewReader(dataJSON))
	if err != nil {
		log.Fatalf("Error create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error get session: %v", err)
	}
	defer resp.Body.Close()

	bb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read body: %v", err)
	}

	var sessionResp struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(bb, &sessionResp); err != nil {
		log.Println(string(bb))
		log.Fatalf("Error decode session: %v", err)
	}
	userID := sessionResp.ID

	return userID, resp.Cookies()
}
