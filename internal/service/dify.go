package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type DifySyncer struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

type Dataset struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Document struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DatasetsResponse struct {
	Data []Dataset `json:"data"`
}

type DocumentsResponse struct {
	Data []Document `json:"data"`
}

func NewDifySyncer(baseURL, apiKey string) *DifySyncer {
	return &DifySyncer{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  &http.Client{},
	}
}

func (d *DifySyncer) CreateDataset(name string) (string, error) {
	payload := map[string]string{
		"name":       name,
		"permission": "only_me",
	}

	resp, err := d.makeRequest("POST", "/v1/datasets", payload)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if id, ok := result["id"].(string); ok {
		return id, nil
	}
	return "", fmt.Errorf("failed to get dataset ID from response")
}

func (d *DifySyncer) CreateDocument(datasetID, name, text string) (string, error) {
	payload := map[string]interface{}{
		"name":               name,
		"text":              text,
		"indexing_technique": "high_quality",
		"process_rule": map[string]string{
			"mode": "automatic",
		},
	}

	resp, err := d.makeRequest("POST", fmt.Sprintf("/v1/datasets/%s/document/create-by-text", datasetID), payload)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if id, ok := result["id"].(string); ok {
		return id, nil
	}
	return "", fmt.Errorf("failed to get document ID from response")
}

func (d *DifySyncer) UpdateDocument(datasetID, documentID, name, text string) error {
	payload := map[string]string{
		"name": name,
		"text": text,
	}

	resp, err := d.makeRequest("POST", fmt.Sprintf("/v1/datasets/%s/documents/%s/update-by-text", datasetID, documentID), payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (d *DifySyncer) GetDatasets() ([]Dataset, error) {
	resp, err := d.makeRequest("GET", "/v1/datasets?page=1&limit=20", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result DatasetsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (d *DifySyncer) GetDocuments(datasetID string) ([]Document, error) {
	resp, err := d.makeRequest("GET", fmt.Sprintf("/v1/datasets/%s/documents", datasetID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result DocumentsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (d *DifySyncer) makeRequest(method, path string, payload interface{}) (*http.Response, error) {
	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, d.baseURL+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+d.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return resp, nil
} 