package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start MCP server (stdio transport)",
	Long:  "Start a Model Context Protocol server over stdin/stdout for use with Claude, Cursor, Codex, and other MCP clients.",
	Run: func(cmd *cobra.Command, args []string) {
		runMCPServer()
	},
}

func runMCPServer() {
	s := server.NewMCPServer(
		"hfpaper",
		"0.2.0",
		server.WithToolCapabilities(true),
	)

	// search_papers
	s.AddTool(mcp.Tool{
		Name:        "search_papers",
		Description: "Search AI research papers on Hugging Face using semantic and full-text search. Returns titles, abstracts, authors, upvotes, and paper IDs.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "Search query (natural language or keywords)",
				},
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "Number of results to return (1-120, default 10)",
					"default":     10,
				},
			},
			Required: []string{"query"},
		},
	}, handleSearchPapers)

	// get_paper
	s.AddTool(mcp.Tool{
		Name:        "get_paper",
		Description: "Get structured metadata for an AI research paper including authors, abstract, AI summary, GitHub repo, project page, upvotes, and linked resources. Accepts arXiv IDs or Hugging Face paper URLs.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"paper_id": map[string]interface{}{
					"type":        "string",
					"description": "arXiv paper ID (e.g. 2602.08025) or Hugging Face/arXiv URL",
				},
			},
			Required: []string{"paper_id"},
		},
	}, handleGetPaper)

	// read_paper
	s.AddTool(mcp.Tool{
		Name:        "read_paper",
		Description: "Read the full content of an AI research paper as markdown. Ideal for summarization, analysis, or citation. Accepts arXiv IDs or Hugging Face paper URLs.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"paper_id": map[string]interface{}{
					"type":        "string",
					"description": "arXiv paper ID (e.g. 2602.08025) or Hugging Face/arXiv URL",
				},
			},
			Required: []string{"paper_id"},
		},
	}, handleReadPaper)

	// daily_papers
	s.AddTool(mcp.Tool{
		Name:        "daily_papers",
		Description: "Get the Hugging Face Daily Papers feed. Returns today's submitted and trending AI research papers with titles, abstracts, and upvote counts.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"date": map[string]interface{}{
					"type":        "string",
					"description": "Date in YYYY-MM-DD format (default: today)",
				},
				"trending": map[string]interface{}{
					"type":        "boolean",
					"description": "Sort by trending/upvotes instead of publish date",
					"default":     false,
				},
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "Number of results (1-100, default 20)",
					"default":     20,
				},
			},
		},
	}, handleDailyPapers)

	// paper_models
	s.AddTool(mcp.Tool{
		Name:        "paper_models",
		Description: "Find Hugging Face model checkpoints linked to a research paper. Returns model IDs, names, and download counts.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"paper_id": map[string]interface{}{
					"type":        "string",
					"description": "arXiv paper ID or URL",
				},
			},
			Required: []string{"paper_id"},
		},
	}, handlePaperModels)

	// paper_datasets
	s.AddTool(mcp.Tool{
		Name:        "paper_datasets",
		Description: "Find Hugging Face datasets linked to a research paper.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"paper_id": map[string]interface{}{
					"type":        "string",
					"description": "arXiv paper ID or URL",
				},
			},
			Required: []string{"paper_id"},
		},
	}, handlePaperDatasets)

	// paper_spaces
	s.AddTool(mcp.Tool{
		Name:        "paper_spaces",
		Description: "Find Hugging Face Spaces (demos/apps) linked to a research paper.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"paper_id": map[string]interface{}{
					"type":        "string",
					"description": "arXiv paper ID or URL",
				},
			},
			Required: []string{"paper_id"},
		},
	}, handlePaperSpaces)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "MCP server error: %v\n", err)
		os.Exit(1)
	}
}

// Tool handlers

func handleSearchPapers(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := getArgs(request.Params.Arguments)
	query := getStringArg(args, "query")
	limit := getIntArg(args, "limit", 10)

	path := fmt.Sprintf("/api/papers/search?q=%s&limit=%d", url.QueryEscape(query), limit)
	body, err := mcpFetch("GET", path, nil)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(body)), nil
}

func handleGetPaper(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := getArgs(request.Params.Arguments)
	paperID := parsePaperID(getStringArg(args, "paper_id"))
	path := fmt.Sprintf("/api/papers/%s", paperID)
	body, err := mcpFetch("GET", path, nil)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// Parse and return a cleaner summary for the agent
	var paper map[string]interface{}
	if json.Unmarshal(body, &paper) == nil {
		clean := map[string]interface{}{
			"id":          paper["id"],
			"title":       paper["title"],
			"summary":     paper["summary"],
			"ai_summary":  paper["ai_summary"],
			"upvotes":     paper["upvotes"],
			"authors":     paper["authors"],
			"publishedAt": paper["publishedAt"],
			"githubRepo":  paper["githubRepo"],
			"projectPage": paper["projectPage"],
		}
		if b, err := json.MarshalIndent(clean, "", "  "); err == nil {
			return mcp.NewToolResultText(string(b)), nil
		}
	}
	return mcp.NewToolResultText(string(body)), nil
}

func handleReadPaper(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := getArgs(request.Params.Arguments)
	paperID := parsePaperID(getStringArg(args, "paper_id"))
	path := fmt.Sprintf("/papers/%s.md", paperID)
	body, err := mcpFetch("GET", path, map[string]string{"Accept": "text/markdown"})
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(body)), nil
}

func handleDailyPapers(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := getArgs(request.Params.Arguments)
	params := url.Values{}
	if date := getStringArg(args, "date"); date != "" {
		params.Set("date", date)
	}
	if getBoolArg(args, "trending") {
		params.Set("sort", "trending")
	} else {
		params.Set("sort", "publishedAt")
	}
	limit := getIntArg(args, "limit", 20)
	params.Set("limit", fmt.Sprintf("%d", limit))

	path := "/api/daily_papers?" + params.Encode()
	body, err := mcpFetch("GET", path, nil)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(body)), nil
}

func handlePaperModels(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := getArgs(request.Params.Arguments)
	paperID := parsePaperID(getStringArg(args, "paper_id"))
	path := fmt.Sprintf("/api/models?filter=arxiv:%s", paperID)
	body, err := mcpFetch("GET", path, nil)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(body)), nil
}

func handlePaperDatasets(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := getArgs(request.Params.Arguments)
	paperID := parsePaperID(getStringArg(args, "paper_id"))
	path := fmt.Sprintf("/api/datasets?filter=arxiv:%s", paperID)
	body, err := mcpFetch("GET", path, nil)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(body)), nil
}

func handlePaperSpaces(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := getArgs(request.Params.Arguments)
	paperID := parsePaperID(getStringArg(args, "paper_id"))
	path := fmt.Sprintf("/api/spaces?filter=arxiv:%s", paperID)
	body, err := mcpFetch("GET", path, nil)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(body)), nil
}

// mcpFetch is a lightweight HTTP helper for MCP handlers
func mcpFetch(method, path string, headers map[string]string) ([]byte, error) {
	resp, err := makeRequest(method, path, nil, headers)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}
