package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AnalyzerService struct {
	BaseURL string
}

type AnalyzeRequest struct {
	JobDescription string `json:"jobDescription"`
	ResumeText     string `json:"resumeText"`
}

type AnalyzeResponse struct {
	FitScore  int      `json:"fitScore"`
	Strengths []string `json:"strengths"`
	Gaps      []string `json:"gaps"`
	Summary   string   `json:"summary"`
}

func (s AnalyzerService) Analyze(req AnalyzeRequest) (AnalyzeResponse, error) {
	payload, _ := json.Marshal(req)
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(s.BaseURL+"/analyze", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return AnalyzeResponse{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return AnalyzeResponse{}, fmt.Errorf("python service returned %d", resp.StatusCode)
	}
	var out AnalyzeResponse
	err = json.NewDecoder(resp.Body).Decode(&out)
	return out, err
}
