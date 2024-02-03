package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func printCandinates(cs []*genai.Candidate) {
	for _, c := range cs {
		for _, p := range c.Content.Parts {
			fmt.Println(p)
		}
	}
}

func run(ctx context.Context) error {
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("API_KEY")))
	if err != nil {
		return err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-pro-vision")

	tableImage, err := os.ReadFile("images/table.jpeg")
	if err != nil {
		return err
	}

	handImage, err := os.ReadFile("images/hand.jpeg")
	if err != nil {
		return err
	}

	prompt := []genai.Part{
		genai.ImageData("jpeg", tableImage),
		genai.ImageData("jpeg", handImage),
		genai.Text("あなたはポーカープレイヤーです。1枚目の画像はテーブルの画像。2枚目はハンドの画像です。この場合の役を答えてください。"),
	}

	resp, err := model.GenerateContent(ctx, prompt...)
	if err != nil {
		return err
	}

	printCandinates(resp.Candidates)

	return nil
}
