package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	openai "github.com/sashabaranov/go-openai"
)

func chatgptStream(message string, tokens int) string {
	c := openai.NewClient("your api key")
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
		color.Set(color.FgHiRed, color.Bold)
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		color.Unset()
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
			color.Set(color.FgHiRed, color.Bold)
			fmt.Printf("\nStream error: %v\n", err)
			color.Unset()
			return ""
		}

		color.Set(color.FgHiCyan)
		fmt.Printf(response.Choices[0].Delta.Content)
		color.Unset()
		gptOutput += response.Choices[0].Delta.Content
	}

}

func createFile() {
	if _, err := os.Stat("history.txt"); os.IsNotExist(err) {
		file, err := os.Create("history.txt")
		if err != nil {
			color.Set(color.FgHiRed, color.Bold)
			fmt.Println("Error creating file")
			color.Unset()
			return
		}
		defer file.Close()
	}
}

func writeFile(message, gptOutput string) {
	file, err := os.OpenFile("history.txt",
		os.O_APPEND|os.O_WRONLY,
		fs.ModeAppend)
	if err != nil {
		color.Set(color.FgHiRed, color.Bold)
		fmt.Println(err.Error())
		color.Unset()
		return
	}
	defer file.Close()

	_, err = fmt.Fprintln(file, message, "\n", gptOutput)
	if err != nil {
		color.Set(color.FgHiRed, color.Bold)
		fmt.Println(err.Error())
		color.Unset()
		return
	}
}

func readFile() {
	file, err := os.ReadFile("history.txt")
	if err != nil {
		color.Set(color.FgHiRed, color.Bold)
		fmt.Println("There's no history yet")
		color.Unset()
		return
	}
	parts := strings.Split(string(file), "\n")
	for i, s := range parts {
		if i%2 == 0 {
			color.Set(color.FgHiYellow, color.Bold)
			fmt.Println(s)
			color.Unset()
		} else {
			color.Set(color.FgHiCyan)
			fmt.Println(s)
			color.Unset()
		}
	}
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		color.Set(color.FgHiRed, color.Bold)
		fmt.Println("\nExiting...")
		color.Unset()
		os.Exit(1)
	}()

	color.Set(color.FgHiGreen, color.Bold)
	fmt.Println("type \"help\" for help")
	color.Unset()
	var exit bool = false
	for !exit {
		scanner := bufio.NewScanner(os.Stdin)
		color.Set(color.FgHiMagenta)
		fmt.Print("~> ")
		color.Unset()
		scanner.Scan()
		message := scanner.Text()

		switch message {
		case "exit":
			exit = true
			continue
		case "help":
			color.Set(color.FgHiGreen)
			fmt.Println("commands:\n \thelp \n\texit \n\thistory clear \n\tclear \n\thistory\nswitches:\n \t-{number of tokens} (lenght of the answer)\n\t-c (like I am 10 years old)\n\t-e (like i am expert)")
			color.Unset()
			continue
		case "clear":
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()
			continue
		case "":
			continue
		case "history":
			readFile()
			continue
		case "history clear":
			err := os.Remove("history.txt")
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			continue
		}

		parts := strings.Split(message, " -")
		var tokens int
		message = ""
		for i, s := range parts {
			if i == 0 {
				message = message + s
				i++
				continue
			}

			switch s {
			case "e":
				s = " like I'm an expert"
				message = message + s
				continue
			case "c":
				s = " like I'm a child"
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

		color.Set(color.FgHiYellow, color.Bold)
		gptOutput := chatgptStream(message, tokens)
		color.Unset()
		createFile()
		writeFile(message, gptOutput)
	}
}
