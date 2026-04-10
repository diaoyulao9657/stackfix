package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type chatReq struct {
	Model    string `json:"model"`
	Messages []msg  `json:"messages"`
	Stream   bool   `json:"stream"`
}

type msg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type streamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

func explain(input, lang string) {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.tokenmix.ai/v1"
	}
	baseURL = strings.TrimRight(baseURL, "/")

	model := os.Getenv("MODEL")
	if model == "" {
		model = "gpt-4o-mini"
	}

	// don't blow up the context window
	if len(input) > 8000 {
		input = input[:8000] + "\n... (truncated)"
	}

	prompt := buildPrompt(lang)

	body, _ := json.Marshal(chatReq{
		Model: model,
		Messages: []msg{
			{Role: "system", Content: prompt},
			{Role: "user", Content: input},
		},
		Stream: true,
	})

	req, err := http.NewRequest("POST", baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		die("request setup failed: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("API_KEY"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		die("couldn't reach API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		handleErr(resp)
		return
	}

	scanner := bufio.NewScanner(resp.Body)
	// bump scanner buffer for big responses
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}
		var chunk streamChunk
		if json.Unmarshal([]byte(data), &chunk) != nil {
			continue
		}
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			fmt.Print(chunk.Choices[0].Delta.Content)
		}
	}
	fmt.Println()
}

func buildPrompt(lang string) string {
	p := `You are stackfix, a terse debugging assistant. A developer piped an error into you.

Do this:
1. What went wrong — one sentence, plain english
2. Why — brief root cause
3. How to fix — code snippet or concrete steps

No preamble. No "I see that..." or "It looks like...". Just the answer. Write like a senior dev helping a teammate.`

	if lang != "" {
		p += fmt.Sprintf("\n\nDetected language: %s", lang)
	}
	return p
}

func handleErr(resp *http.Response) {
	body, _ := io.ReadAll(resp.Body)

	switch resp.StatusCode {
	case 401, 403:
		fmt.Fprintln(os.Stderr, "Error: Authentication failed — your API key may be invalid or expired.")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "To fix this:")
		fmt.Fprintln(os.Stderr, "  - Check that your API_KEY in .env is correct")
		fmt.Fprintln(os.Stderr, "  - Get a new API key at https://tokenmix.ai (new accounts get $1 free credit)")
		fmt.Fprintln(os.Stderr, "  - Or use any OpenAI-compatible API by updating BASE_URL in .env")
	case 429:
		fmt.Fprintln(os.Stderr, "Error: Rate limited. Wait a moment and try again.")
	default:
		fmt.Fprintf(os.Stderr, "Error: API returned HTTP %d\n", resp.StatusCode)
		var errBody struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		if json.Unmarshal(body, &errBody) == nil && errBody.Error.Message != "" {
			fmt.Fprintf(os.Stderr, "  %s\n", errBody.Error.Message)
		}
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "If this persists, check:")
		fmt.Fprintln(os.Stderr, "  - Your API key is valid")
		fmt.Fprintln(os.Stderr, "  - Your BASE_URL is correct")
		fmt.Fprintln(os.Stderr, "  - The API provider is not experiencing downtime")
	}
	os.Exit(1)
}
