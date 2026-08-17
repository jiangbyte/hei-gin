package main

import (
	"io"
	"net/http"
	"strings"
	"time"
)

func doRaw(method, url, token, body string) (status int, raw []byte, ar apiResp, err error) {
	var reqBody io.Reader
	if body != "" {
		reqBody = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return 0, nil, ar, err
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", token)
	}
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, ar, err
	}
	defer resp.Body.Close()
	raw, _ = io.ReadAll(resp.Body)
	status = resp.StatusCode
	ar, _, _ = parseLoose(raw)
	return status, raw, ar, nil
}
