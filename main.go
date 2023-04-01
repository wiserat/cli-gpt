package main

import (
	"context"
	"fmt"
	"os"
	"bufio"
	"io"
	"errors"
	"strings"
	"os/exec"
	"strconv"
	openai "github.com/sashabaranov/go-openai"
)



func chatgpt(message string, tokens int) string {
	c := openai.NewClient("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")
	ctx := context.Background()

	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: tokens,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: message,
			},
		},
		Stream: true,
	}
	stream, err := c.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return ""
	}
	defer stream.Close()

	var gptOutput string = ""
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return gptOutput
		}

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			return ""
		}

		fmt.Printf(response.Choices[0].Delta.Content)
		gptOutput += response.Choices[0].Delta.Content
	}

}


func create_file() {
	file, err := os.Create("history.txt")
	if err != nil {
		fmt.Println("Error creating file")
		return
	}
	defer file.Close()
}


func write_file(message, gptOutput string) {

}


func main() {
	fmt.Println("help for help")
	var exit bool = false
	for !exit {
		
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("~> ")
		scanner.Scan()
		message := scanner.Text()

		if message == "exit" {
			exit = true
		} else if message == "help" {
			fmt.Println("commands:\n \thelp \n\texit \n\tclear \n\thistory\nswitches:\n \t-{number of tokens} (lenght of the answer)\n\t-c (like I am 10 years old)\n\t-e (like i am expert)")
		} else if message == "clear" {
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()
		} else if message == "" {
			continue
		} else {
			parts := strings.Split(message, " -")
			var tokens int
			message = ""
			for i, s:= range parts {
				if i == 0 {
					message = message + s
					i++
					continue
				}

				switch s {
					case "e":
						s = " like I am expert"
						message = message + s
						continue
					case "c":
						s = " like I am child"
						message = message + s
						continue
				}
				partsT := strings.Split(s, "t")
				tPart, err := strconv.Atoi(partsT[0])
				if err != nil {
					fmt.Println("error: switch does not exist")
					return
				}
				tokens = tPart
			}
			
			gptOutput := chatgpt(message, tokens)
			create_file()
			write_file(message, gptOutput)
		}
	}
}