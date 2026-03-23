package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// handleFormattedResponse replaces the generic handleResponse with format-aware output
func handleFormattedResponse(resp *http.Response, err error, kind string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		fmt.Fprintf(os.Stderr, "Error: not found\n")
		os.Exit(2)
	}
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(os.Stderr, "Error: %s (HTTP %d)\n", string(body), resp.StatusCode)
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading response: %v\n", err)
		os.Exit(1)
	}

	// JSON or markdown mode: pass through
	if jsonFlag {
		fmt.Println(string(body))
		return
	}
	if markdownFlag || kind == "read" {
		fmt.Println(string(body))
		return
	}

	// Human-readable formatting
	switch kind {
	case "search":
		formatSearch(body)
	case "daily":
		formatDaily(body)
	case "get":
		formatGet(body)
	case "models", "datasets", "spaces":
		formatResources(body, kind)
	default:
		fmt.Println(string(body))
	}
}

func formatSearch(body []byte) {
	var results []map[string]interface{}
	if err := json.Unmarshal(body, &results); err != nil {
		fmt.Println(string(body))
		return
	}

	if len(results) == 0 {
		fmt.Println("  No results found.")
		return
	}

	fmt.Printf("\n  \033[1m🔍 %d results\033[0m\n\n", len(results))

	for i, item := range results {
		paper, _ := item["paper"].(map[string]interface{})
		if paper == nil {
			continue
		}
		title := cleanStr(paper["title"])
		id := cleanStr(paper["id"])
		upvotes := int(numVal(paper["upvotes"]))
		summary := cleanStr(paper["ai_summary"])
		if summary == "" {
			summary = truncate(cleanStr(paper["summary"]), 120)
		}

		fmt.Printf("  \033[1m%d. %s\033[0m\n", i+1, title)
		fmt.Printf("     \033[33m%s\033[0m · %d ⬆\n", id, upvotes)
		if summary != "" {
			fmt.Printf("     \033[2m%s\033[0m\n", summary)
		}
		fmt.Println()
	}
}

func formatDaily(body []byte) {
	var results []map[string]interface{}
	if err := json.Unmarshal(body, &results); err != nil {
		fmt.Println(string(body))
		return
	}

	if len(results) == 0 {
		fmt.Println("  No papers today.")
		return
	}

	fmt.Printf("\n  \033[1m📰 Daily Papers — %d results\033[0m\n\n", len(results))

	for i, item := range results {
		paper, _ := item["paper"].(map[string]interface{})
		if paper == nil {
			continue
		}
		title := cleanStr(paper["title"])
		id := cleanStr(paper["id"])
		upvotes := int(numVal(paper["upvotes"]))

		// Authors
		authors := extractAuthors(paper, 2)

		fmt.Printf("  \033[1m%d. %s\033[0m\n", i+1, title)
		fmt.Printf("     \033[33m%s\033[0m · %d ⬆", id, upvotes)
		if authors != "" {
			fmt.Printf(" · by %s", authors)
		}
		fmt.Println()
		fmt.Println()
	}
}

func formatGet(body []byte) {
	var paper map[string]interface{}
	if err := json.Unmarshal(body, &paper); err != nil {
		fmt.Println(string(body))
		return
	}

	title := cleanStr(paper["title"])
	id := cleanStr(paper["id"])
	upvotes := int(numVal(paper["upvotes"]))
	summary := cleanStr(paper["summary"])
	aiSummary := cleanStr(paper["ai_summary"])
	github := cleanStr(paper["githubRepo"])
	project := cleanStr(paper["projectPage"])
	published := cleanStr(paper["publishedAt"])
	authors := extractAuthors(paper, 5)

	// Parse date
	dateStr := published
	if t, err := time.Parse(time.RFC3339, published); err == nil {
		dateStr = t.Format("Jan 2, 2006")
	}

	fmt.Printf("\n  \033[1m📄 %s\033[0m\n\n", title)
	fmt.Printf("  ID:        \033[33m%s\033[0m\n", id)
	fmt.Printf("  Published: %s\n", dateStr)
	if authors != "" {
		fmt.Printf("  Authors:   %s\n", authors)
	}
	fmt.Printf("  Upvotes:   %d\n", upvotes)
	if github != "" {
		fmt.Printf("  GitHub:    %s\n", github)
	}
	if project != "" {
		fmt.Printf("  Project:   %s\n", project)
	}

	if aiSummary != "" {
		fmt.Printf("\n  \033[1mSummary:\033[0m\n  %s\n", aiSummary)
	}
	if summary != "" {
		fmt.Printf("\n  \033[1mAbstract:\033[0m\n  %s\n", truncate(summary, 500))
	}
	fmt.Println()
}

func formatResources(body []byte, kind string) {
	var items []map[string]interface{}
	if err := json.Unmarshal(body, &items); err != nil {
		fmt.Println(string(body))
		return
	}

	if len(items) == 0 {
		fmt.Printf("  No linked %s found.\n", kind)
		return
	}

	emoji := map[string]string{"models": "🤖", "datasets": "📊", "spaces": "🚀"}
	fmt.Printf("\n  \033[1m%s %d linked %s\033[0m\n\n", emoji[kind], len(items), kind)

	for i, item := range items {
		id := cleanStr(item["id"])
		if id == "" {
			id = cleanStr(item["modelId"])
		}
		if id == "" {
			id = cleanStr(item["_id"])
		}
		downloads := int(numVal(item["downloads"]))
		likes := int(numVal(item["likes"]))

		fmt.Printf("  %d. \033[1m%s\033[0m", i+1, id)
		if downloads > 0 {
			fmt.Printf(" · %d downloads", downloads)
		}
		if likes > 0 {
			fmt.Printf(" · %d ❤", likes)
		}
		fmt.Println()
	}
	fmt.Println()
}

// Helpers

func cleanStr(v interface{}) string {
	s, _ := v.(string)
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "  ", " ")
	return s
}

func numVal(v interface{}) float64 {
	f, _ := v.(float64)
	return f
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func extractAuthors(paper map[string]interface{}, max int) string {
	authorsRaw, ok := paper["authors"].([]interface{})
	if !ok || len(authorsRaw) == 0 {
		return ""
	}
	var names []string
	for i, a := range authorsRaw {
		if i >= max {
			names = append(names, fmt.Sprintf("+%d more", len(authorsRaw)-max))
			break
		}
		if author, ok := a.(map[string]interface{}); ok {
			if name, ok := author["name"].(string); ok {
				names = append(names, name)
			}
		}
	}
	return strings.Join(names, ", ")
}
