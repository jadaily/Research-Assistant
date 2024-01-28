package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func generateQuestions(topic string) []string {

	payload := promptRequestPayload{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "user",
				Content: fmt.Sprintf("Generate 8 research questions about %s.", topic),
			},
		},
		MaxTokens: 275,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error creating JSON payload:", err)
		return nil
	}

	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}

	var chatGPTResponseBody chatGPTResponseBody
	if err := json.Unmarshal(body, &chatGPTResponseBody); err != nil {
		fmt.Println("Error decoding response:", err.Error(), string(body))
		return nil
	}

	var generatedQuestions []string
	for _, choice := range chatGPTResponseBody.Choices {
		generatedQuestions = append(generatedQuestions, choice.Message.Content)
	}

	return generatedQuestions
}

func getArticles(userInput string) []string {
	client := &http.Client{}

	payload := promptRequestPayload{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "user",
				Content: fmt.Sprintf("Give me some scholarly articles that are related to '%s':", userInput),
			},
		},
		MaxTokens: 250,
	}

	payloadBytes, _ := json.Marshal(payload)

	request, _ := http.NewRequest("POST", apiEndpoint, bytes.NewReader(payloadBytes))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+apiKey)

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error making API request:", err)
		return nil
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)

	var gptResponse chatGPTResponseBody
	if err := json.Unmarshal(body, &gptResponse); err != nil {
		fmt.Println("Error decoding response:", err)
		return nil
	}

	var articles []string
	for _, choice := range gptResponse.Choices {
		articles = append(articles, choice.Message.Content)
	}

	return articles

}

func getUserInput() string {
	var input string
	fmt.Scanln(&input)
	return input
}
