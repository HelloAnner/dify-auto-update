package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// DifySyncer 用于与Dify API进行交互的结构体
type DifySyncer struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

// Dataset 数据集结构体
type Dataset struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Document 文档结构体
type Document struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// DatasetsResponse 数据集响应结构体
type DatasetsResponse struct {
	Data  []Dataset `json:"data"`
	Total int       `json:"total"`
	Page  int       `json:"page"`
	Limit int       `json:"limit"`
}

// DocumentsResponse 文档响应结构体
type DocumentsResponse struct {
	Data []Document `json:"data"`
}

// NewDifySyncer 创建一个新的DifySyncer实例
func NewDifySyncer(baseURL, apiKey string) *DifySyncer {
	return &DifySyncer{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  &http.Client{},
	}
}

// CreateDataset 创建一个新的数据集
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
	return "", fmt.Errorf("获取数据集ID失败")
}

// CreateDocument 在指定数据集中创建新文档
func (d *DifySyncer) CreateDocument(datasetID, name, text string) (string, error) {
	payload := map[string]interface{}{
		"name":               name,
		"text":               text,
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
	return "", nil
}

// UpdateDocument 更新指定数据集中的文档
func (d *DifySyncer) UpdateDocument(datasetID, documentID, name, text string) error {
	payload := map[string]interface{}{
		"name": name,
		"text": text,
		"process_rule": map[string]string{
			"mode": "automatic",
		},
	}

	resp, err := d.makeRequest("POST", fmt.Sprintf("/v1/datasets/%s/documents/%s/update-by-text", datasetID, documentID), payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// GetDatasets 获取所有数据集列表
func (d *DifySyncer) GetDatasets() ([]Dataset, error) {
	var allDatasets []Dataset
	page := 1
	limit := 20

	for {
		resp, err := d.makeRequest("GET", fmt.Sprintf("/v1/datasets?page=%d&limit=%d", page, limit), nil)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var result DatasetsResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, err
		}

		allDatasets = append(allDatasets, result.Data...)

		if len(result.Data) < limit || result.Total <= len(allDatasets) {
			break
		}
		page++
	}

	return allDatasets, nil
}

// GetDocuments 获取指定数据集中的所有文档
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

// DeleteDocument 删除指定数据集中的文档
func (d *DifySyncer) DeleteDocument(datasetID, documentID string) error {
	resp, err := d.makeRequest("DELETE", fmt.Sprintf("/v1/datasets/%s/documents/%s", datasetID, documentID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// makeRequest 发送HTTP请求到Dify API
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
		return nil, fmt.Errorf("API请求失败， url: %s, 状态码：%d，错误信息：%s", d.baseURL+path, resp.StatusCode, string(bodyBytes))
	}

	return resp, nil
}
