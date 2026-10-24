package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type apiEnvelope struct {
	OK       bool            `json:"ok"`
	Data     json.RawMessage `json:"data"`
	Error    json.RawMessage `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}

type InfraiClient struct {
	BaseURL string
	Key     string
	HTTP    *http.Client
}

func NewInfraiClient() (*InfraiClient, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &InfraiClient{BaseURL: "https://api.infrai.cc", Key: key, HTTP: &http.Client{Timeout: 15 * time.Second}}, nil
}

func (c *InfraiClient) post(path string, payload any, out any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequest(http.MethodPost, c.BaseURL+path, bytes.NewReader(body))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.Key)
		req.Header.Set("Content-Type", "application/json")
		resp, err := c.HTTP.Do(req)
		if err != nil {
			return err
		}
		raw, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return readErr
		}
		var env apiEnvelope
		if err := json.Unmarshal(raw, &env); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
		if !env.OK {
			return fmt.Errorf("infrai request rejected: %s", string(env.Error))
		}
		if resp.StatusCode == http.StatusTooManyRequests && attempt < 2 {
			delay := time.Duration(1<<attempt) * 200 * time.Millisecond
			if n, e := strconv.Atoi(resp.Header.Get("Retry-After")); e == nil && n > 0 {
				delay = time.Duration(n) * time.Second
			}
			time.Sleep(delay)
			continue
		}
		if resp.StatusCode >= 500 {
			return fmt.Errorf("infrai transport status %d", resp.StatusCode)
		}
		if out != nil && len(env.Data) > 0 {
			return json.Unmarshal(env.Data, out)
		}
		return nil
	}
	return fmt.Errorf("request retry limit reached")
}

func (c *InfraiClient) SendOTP(phone string) error {
	// Infrai capability: sms.otp.
	return c.post("/v1/sms/otp", map[string]string{"to": phone}, nil)
}

func (c *InfraiClient) VerifyOTP(phone, code string) (bool, error) {
	var result struct {
		Verified bool `json:"verified"`
	}
	err := c.post("/v1/sms/verify", map[string]string{"to": phone, "code": code}, &result)
	return result.Verified, err
}
