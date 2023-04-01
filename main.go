package main

import (
	"context"
	"fmt"
	"os"
	"bufio"
	openai "github.com/sashabaranov/go-openai"
)

func chatgpt(message string) string {
	client := openai.NewClient("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")
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: message,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return ""
	}

	return resp.Choices[0].Message.Content
}

func main() {
	fmt.Println("commands: exit, help")
	var exit bool = false
	for !exit {
		
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("~> ")
		scanner.Scan()
		message := scanner.Text()

		if message == "exit" {
			exit = true
		} else if message == "help" {
			fmt.Println("help")
		} else {
			fmt.Println(chatgpt(message))
		}
	}
}