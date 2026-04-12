package main

import (
	"fmt"
	"os"

	"sse/cmd/keygen"
	"sse/cmd/server"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {
	case "serve":
		server.Run()
	case "keygen":
		keygen.Run()
	case "help", "--help", "-h":
		printHelp()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", os.Args[1])
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("  sse-hub serve   Start the SSE hub server")
	fmt.Println("  sse-hub keygen  Generate an API key for X-API-Key authentication")
}
