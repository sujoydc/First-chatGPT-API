package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type ChatGPTRequest struct {
	Prompt      string  `json:"prompt"`
	Model       string  `json:"model"`
	MaxTokens   int32   `json:"max_tokens"`
	Temperature float32 `json:"temperature"`
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

func main() {

	// read the key=value pair file first
	keyMapVals, err := readPropsFromFile(".sujix-prop")
	if err != nil {
		fmt.Println("Mandatory file .sujix-prop not found. Please create this file with the following:")
		fmt.Println("API-KEY=[Your Key]")
		os.Exit(9)
	}

	logPropsForDebug(keyMapVals)

	fmt.Println("Welcome to my chatGPT prompt.")
	fmt.Println("For Text Completion, enter 1")
	fmt.Println("For Chat Completion, enter 2")
	fmt.Println("Anytime to exit, please enter 0")

	var option int16
	fmt.Print("Please enter your option? - ")
	fmt.Scanln(&option)

	for option != 0 {

		if option == 1 {

			prompt := getPromptInput()
			fmt.Printf("Your question is: %s\n", prompt)

			model := "text-davinci-003" //"text-ada-001"
			maxTokens := int32(7)
			temperature := float32(1.0)

			chatGPTRequest := ChatGPTRequest{
				Prompt:      prompt,
				Model:       model,
				MaxTokens:   maxTokens,
				Temperature: temperature,
			}

			// Replace with the actual ChatGPT API URL
			url := "https://api.openai.com/v1/completions"

			responseText, err := postToChatGPT(chatGPTRequest, url, keyMapVals)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Response: %s\n", responseText)

		} else if option == 2 {
			fmt.Println("Work in progress! ")

		} else {
			fmt.Println("Invalid option! Bye bye!")
			os.Exit(0)
		}

		fmt.Print("Please enter your option? - ")
		fmt.Scanln(&option)
	}

}

func postToChatGPT(payload ChatGPTRequest, url string, kv map[string]string) (string, error) {

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
	req.Header.Set("Authorization", "Bearer "+kv["API-KEY"])

	// Send the HTTP request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Check the response status code
	if response.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(response.Body)
		return "", fmt.Errorf("failed to make POST request to ChatGPT API: %s\nResponse body: %s", response.Status, string(body))
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

func readPropsFromFile(filename string) (map[string]string, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2) // Split on the first "="
		if len(parts) != 2 {
			fmt.Println("Ignoring line:", line)
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		result[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil

}

func getPromptInput() string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("What is your question? - ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return ""
	}

	// Remove the newline character at the end of the input
	return input[:len(input)-1]
}

func logPropsForDebug(kv map[string]string) {
	for key, value := range kv {
		fmt.Printf("Key: %s, Value: %s\n", key, value)
	}
}
