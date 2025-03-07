package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	stdhttp "net/http"
	"time"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/http"
)

// Configuration
const (
	defaultPort = 8080
	sseEndpoint = "/sse"
)

// TiDB AI API response structure
type AskResponse struct {
	Content string `json:"content"`
}

// askQuestion sends a question to the TiDB AI API and returns the response
func askQuestion(question string) (string, error) {
	url := "https://tidb.ai/api/v1/chats"

	// Construct request body
	requestBody, err := json.Marshal(map[string]interface{}{
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": question,
			},
		},
		"chat_engine": "default",
		"stream":      false,
	})
	if err != nil {
		return "", fmt.Errorf("error marshaling request body: %v", err)
	}

	// Create HTTP request
	req, err := stdhttp.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// Set request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")

	// Send request
	client := &stdhttp.Client{
		Timeout: 60 * time.Second, // Set a timeout to prevent hanging
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	// Parse response
	var askResp AskResponse
	err = json.Unmarshal(body, &askResp)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling response: %v", err)
	}

	return askResp.Content, nil
}

// QuestionArgs defines the arguments for the ask tool
type QuestionArgs struct {
	Question string `json:"question" jsonschema:"required,description=The question to ask TiDB AI"`
}

// registerTools registers the tools for the server
func registerTools(server *mcp_golang.Server) error {
	// Register the ask tool
	err := server.RegisterTool("ask", "Ask a question to TiDB AI", func(args QuestionArgs) (*mcp_golang.ToolResponse, error) {
		log.Printf("Processing question: %s", args.Question)

		answer, err := askQuestion(args.Question)
		if err != nil {
			return nil, fmt.Errorf("error asking question: %v", err)
		}

		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(answer)), nil
	})
	if err != nil {
		return fmt.Errorf("error registering tool: %v", err)
	}

	return nil
}

func main() {
	// Define command-line flags
	port := flag.Int("port", defaultPort, "Port to listen on for SSE transport")
	flag.Parse()

	// Create an HTTP transport for SSE
	addr := fmt.Sprintf(":%d", *port)
	transport := http.NewHTTPTransport(sseEndpoint).WithAddr(addr)

	// Create a new server with the transport
	server := mcp_golang.NewServer(transport,
		mcp_golang.WithName("TiDB AI"),
		mcp_golang.WithVersion("1.0.0"),
	)

	// Register tools
	if err := registerTools(server); err != nil {
		log.Fatalf("Error registering tools: %v", err)
	}

	// Start the server
	log.Printf("Starting SSE MCP server on %s with endpoint %s", addr, sseEndpoint)
	if err := server.Serve(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
