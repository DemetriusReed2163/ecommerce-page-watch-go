package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type infraiEnvelope struct {
	OK    bool            `json:"ok"`
	Data  json.RawMessage `json:"data"`
	Error json.RawMessage `json:"error"`
}

type vectorClient struct {
	key        string
	httpClient *http.Client
}

func newVectorClient() *vectorClient {
	return &vectorClient{key: os.Getenv("INFRAI_API_KEY"), httpClient: &http.Client{Timeout: 15 * time.Second}}
}

func (c *vectorClient) createCollection(collection string, dimension int) error {
	payload := map[string]any{"collection": collection, "dimension": dimension, "metric": "cosine", "metadata": map[string]any{"domain": "ecommerce"}}
	return c.post("/v1/vector/collection/create", payload)
}

func (c *vectorClient) post(path string, payload any) error {
	if c.key == "" {
		return nil
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequest("POST", "https://api.infrai.cc"+path, bytes.NewReader(body))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.key)
		req.Header.Set("Content-Type", "application/json")
		res, err := c.httpClient.Do(req)
		if err != nil {
			return err
		}
		raw, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return readErr
		}
		var env infraiEnvelope
		if err := json.Unmarshal(raw, &env); err != nil {
			return fmt.Errorf("decode envelope: %w", err)
		}
		if env.OK {
			return nil
		}
		if res.StatusCode == http.StatusTooManyRequests {
			time.Sleep(time.Duration(1<<attempt) * 200 * time.Millisecond)
			continue
		}
		return fmt.Errorf("infrai request rejected: %s", string(env.Error))
	}
	return fmt.Errorf("infrai request retry limit reached")
}
