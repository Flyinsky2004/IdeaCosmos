package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

/**
 * @author Flyinsky
 * @email w2084151024@gmail.com
 * @date 2024/12/29 17:25
 */
// Dalle3Request represents the request body for DALL-E 3 API
// Dalle3Request represents the request body for DALL-E 3 API
type Dalle3Request struct {
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
	Model  string `json:"model"`
}

// Dalle3Response represents the response from DALL-E 3 API
type Dalle3Response struct {
	Data []struct {
		URL            string `json:"url"`
		Revised_Prompt string `json:"revised_prompt"`
	} `json:"data"`
}

// GenerateImage generates an image using DALL-E 3 API
func GenerateImage(prompt, baseURL, apiKey string) (string, error) {
	// Construct the request body
	requestBody := Dalle3Request{
		Prompt: prompt,
		N:      1,           // Number of images to generate
		Size:   "1024x1024", // Image size
		Model:  "dall-e-3",
	}
	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Create the HTTP request
	url := fmt.Sprintf("%s/v1/images/generations", baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var dalleResponse Dalle3Response
	err = json.Unmarshal(body, &dalleResponse)
	if err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	// Check if there is at least one URL in the response
	if len(dalleResponse.Data) == 0 {
		return "", fmt.Errorf("no image URL returned in response")
	}

	// Return the first image URL
	return dalleResponse.Data[0].URL, nil
}

// type GenerateResultData struct {
// 	Prompt string `json:"revised_prompt"`
// 	Url    string `json:"url"`
// }
// type GenerateResult struct {
// 	Created int                  `json:"created"`
// 	Data    []GenerateResultData `json:"data"`
// }

// func GenerateImageVAPI(prompt string) (string, error) {
// 	url := "https://api.gpt.ge/v1/images/generations"
// 	method := "POST"

// 	payload := strings.NewReader(`{
// 		"model": "dall-e-3",
// 		"prompt": ` + prompt + `,
// 		"n": 1,
// 		"size": "1024x1024"
// 	}`)

// 	client := &http.Client{}
// 	req, err := http.NewRequest(method, url, payload)

// 	if err != nil {
// 		return "", err
// 	}
// 	req.Header.Add("Content-Type", "application/json")
// 	req.Header.Add("Authorization", "Bearer sk-hySadfvZfjMxfWx12b302e8c832c4aEeBf7e44C5138bE860")
// 	res, err := client.Do(req)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer res.Body.Close()

// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		return "", err
// 	}

// 	var result GenerateResult
// 	err = json.Unmarshal(string(body), &result)
// 	if err != nil {
// 		return "", err
// 	}

// 	if len(result.Data) == 0 {
// 		return "", fmt.Errorf("no image generated")
// 	}

// 	return result.Data[0].Url, nil
// }

func DownloadImage(imageURL string) (string, error) {
	// Create the /uploads directory if it doesn't exist
	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}

	// Generate a unique file name
	randomSuffix := make([]byte, 3) // 3 bytes = 6 hex characters
	if _, err := rand.Read(randomSuffix); err != nil {
		return "", fmt.Errorf("failed to generate random suffix: %v", err)
	}
	fileName := fmt.Sprintf("%s_%06x.webp", time.Now().Format("20060102"), randomSuffix)
	filePath := filepath.Join(uploadDir, fileName)

	// Download the image
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected HTTP status: %d", resp.StatusCode)
	}

	// Create the output file
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Write the image data to the file
	if _, err := io.Copy(file, resp.Body); err != nil {
		return "", fmt.Errorf("failed to save image: %v", err)
	}

	return fileName, nil
}
