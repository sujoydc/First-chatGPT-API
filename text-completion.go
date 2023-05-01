package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type ChatGPTRequest struct {
	Prompt      string  `json:"prompt"`
	Model       string  `json:"model"`
	MaxTokens   int32   `json:"max_tokens"`
	Temperature float32 `json:"temparature"`
}

type ChatGPTResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

func postToChatGPT(payload ChatGPTRequest) (string, error) {
	// Replace with the actual ChatGPT API URL
	url := "https://api.openai.com/v1/completions"

	// Create the request payload
	//payload := &ChatGPTRequest{Prompt: prompt}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	// Set required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer MY_TOKEN")

	// Send the HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Check the response status code
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to make POST request to ChatGPT API: %s", response.Status)
	}

	// Read and parse the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var chatGPTResponse ChatGPTResponse
	err = json.Unmarshal(body, &chatGPTResponse)
	if err != nil {
		return "", err
	}

	// Return the first choice's text (assuming there is at least one choice)
	if len(chatGPTResponse.Choices) > 0 {
		return chatGPTResponse.Choices[0].Text, nil
	}

	return "", fmt.Errorf("no choices found in the ChatGPT API response")
}

func main() {
	// Replace with your desired prompt
	prompt := "Come up with a boy's name for my friend Sneha's second child"
	model := "text-davinci-003" //"text-ada-001" // "text-davinci-003"
	maxTokens := int32(7)
	temperature := float32(1.0)

	chatGPTRequest := ChatGPTRequest{
		Prompt:      prompt,
		Model:       model,
		MaxTokens:   maxTokens,
		Temperature: temperature,
	}
	responseText, err := postToChatGPT(chatGPTRequest)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Response: %s\n", responseText)
}
