package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

var version = "0.1.0"

func main() {
	loadEnv()

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-v", "--version":
			fmt.Println("stackfix", version)
			return
		case "-h", "--help":
			usage()
			return
		}
	}

	input := readInput()
	if input == "" {
		usage()
		os.Exit(1)
	}

	checkAPIKey()

	lang := detect(input)
	explain(input, lang)
}

func readInput() string {
	var parts []string

	stat, _ := os.Stdin.Stat()
	piped := (stat.Mode() & os.ModeCharDevice) == 0

	if piped {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			die("failed to read stdin: %v", err)
		}
		parts = append(parts, string(data))
	}

	for _, a := range os.Args[1:] {
		if strings.HasPrefix(a, "-") {
			continue
		}
		parts = append(parts, a)
	}

	return strings.TrimSpace(strings.Join(parts, "\n"))
}

func checkAPIKey() {
	if os.Getenv("API_KEY") != "" {
		return
	}
	fmt.Fprintln(os.Stderr, "Error: API key not configured.")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "To get started:")
	fmt.Fprintln(os.Stderr, "  1. Get a free API key at https://tokenmix.ai ($1 free credit)")
	fmt.Fprintln(os.Stderr, "     Or use any OpenAI-compatible API provider")
	fmt.Fprintln(os.Stderr, "  2. Create ~/.config/stackfix/.env (or .env in current dir)")
	fmt.Fprintln(os.Stderr, "  3. Set your API_KEY in the .env file")
	os.Exit(1)
}

func loadEnv() {
	// check several locations — first match wins per-key
	home, _ := os.UserHomeDir()
	paths := []string{
		".env",
		home + "/.config/stackfix/.env",
	}

	for _, p := range paths {
		if p == "" {
			continue
		}
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || line[0] == '#' {
				continue
			}
			k, v, ok := strings.Cut(line, "=")
			if !ok {
				continue
			}
			k = strings.TrimSpace(k)
			v = strings.Trim(strings.TrimSpace(v), `"'`)
			if os.Getenv(k) == "" {
				os.Setenv(k, v)
			}
		}
	}
}

func usage() {
	fmt.Print(`stackfix - explain errors with AI

Usage:
  stackfix "error message"
  command 2>&1 | stackfix
  stackfix < error.log

Environment:
  API_KEY    API key (required)
  BASE_URL   API endpoint (default: https://api.tokenmix.ai/v1)
  MODEL      Model name (default: gpt-4o-mini)
`)
}

func die(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "stackfix: "+format+"\n", args...)
	os.Exit(1)
}
