package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// Constants
const (
	BaseURL = "https://huggingface.co"
	UserAgent = "hfpaper-cli/1.0"
)

// Global Flags
var (
	jsonFlag     bool
	tableFlag    bool
	markdownFlag bool
	limitFlag    int
	dateFlag     string
	trendingFlag bool
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "hfpaper",
		Short: "Hugging Face Papers CLI",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Auto-detect JSON mode if not a TTY
			if !jsonFlag && !tableFlag && !markdownFlag {
				if !isTTY() {
					jsonFlag = true
				}
			}
		},
	}

	rootCmd.PersistentFlags().BoolVar(&jsonFlag, "json", false, "Output as JSON")
	rootCmd.PersistentFlags().BoolVar(&tableFlag, "table", false, "Output as table")
	rootCmd.PersistentFlags().BoolVar(&markdownFlag, "markdown", false, "Output as markdown")

	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(readCmd)
	rootCmd.AddCommand(dailyCmd)
	rootCmd.AddCommand(modelsCmd)
	rootCmd.AddCommand(datasetsCmd)
	rootCmd.AddCommand(spacesCmd)
	rootCmd.AddCommand(indexCmd)
	rootCmd.AddCommand(citeCmd)
	rootCmd.AddCommand(summaryCmd)
	rootCmd.AddCommand(openCmd)
	rootCmd.AddCommand(mcpCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// Utils
func isTTY() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

func parsePaperID(input string) string {
	// regex for 2602.08025 or 2602.08025v1
	re := regexp.MustCompile(`(\d{4,5}\.\d{4,5}(v\d+)?)`)
	matches := re.FindStringSubmatch(input)
	if len(matches) > 0 {
		return matches[1]
	}
	// Fallback to simple split if it's just the ID
	return strings.TrimSuffix(input, ".md")
}

func makeRequest(method, path string, body io.Reader, headers map[string]string) (*http.Response, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(method, BaseURL+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", UserAgent)
	if token := os.Getenv("HF_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return client.Do(req)
}

func handleResponse(resp *http.Response, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		fmt.Fprintf(os.Stderr, "Error: Not Found\n")
		os.Exit(2)
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(os.Stderr, "Error: %s (Status: %d)\n", string(body), resp.StatusCode)
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading response: %v\n", err)
		os.Exit(1)
	}

	if jsonFlag {
		fmt.Println(string(body))
	} else {
		// Simple text/human output for now, could be improved with table writers
		fmt.Println(string(body))
	}
}

// Commands Implementation

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for papers",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := url.QueryEscape(args[0])
		path := fmt.Sprintf("/api/papers/search?q=%s&limit=%d", query, limitFlag)
		resp, err := makeRequest("GET", path, nil, nil)
		handleResponse(resp, err)
	},
}

var getCmd = &cobra.Command{
	Use:   "get <paper_id>",
	Short: "Get paper metadata",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		paperID := parsePaperID(args[0])
		path := fmt.Sprintf("/api/papers/%s", paperID)
		resp, err := makeRequest("GET", path, nil, nil)
		handleResponse(resp, err)
	},
}

var readCmd = &cobra.Command{
	Use:   "read <paper_id>",
	Short: "Read paper content as markdown",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		paperID := parsePaperID(args[0])
		path := fmt.Sprintf("/papers/%s.md", paperID)
		resp, err := makeRequest("GET", path, nil, map[string]string{"Accept": "text/markdown"})
		handleResponse(resp, err)
	},
}

var dailyCmd = &cobra.Command{
	Use:   "daily",
	Short: "Get daily papers",
	Run: func(cmd *cobra.Command, args []string) {
		params := url.Values{}
		if dateFlag != "" {
			params.Set("date", dateFlag)
		}
		if trendingFlag {
			params.Set("sort", "trending")
		} else {
			params.Set("sort", "publishedAt")
		}
		params.Set("limit", fmt.Sprintf("%d", limitFlag))

		path := "/api/daily_papers?" + params.Encode()
		resp, err := makeRequest("GET", path, nil, nil)
		handleResponse(resp, err)
	},
}

var modelsCmd = &cobra.Command{
	Use:   "models <paper_id>",
	Short: "List linked models",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		paperID := parsePaperID(args[0])
		path := fmt.Sprintf("/api/models?filter=arxiv:%s", paperID)
		resp, err := makeRequest("GET", path, nil, nil)
		handleResponse(resp, err)
	},
}

var datasetsCmd = &cobra.Command{
	Use:   "datasets <paper_id>",
	Short: "List linked datasets",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		paperID := parsePaperID(args[0])
		path := fmt.Sprintf("/api/datasets?filter=arxiv:%s", paperID)
		resp, err := makeRequest("GET", path, nil, nil)
		handleResponse(resp, err)
	},
}

var spacesCmd = &cobra.Command{
	Use:   "spaces <paper_id>",
	Short: "List linked spaces",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		paperID := parsePaperID(args[0])
		path := fmt.Sprintf("/api/spaces?filter=arxiv:%s", paperID)
		resp, err := makeRequest("GET", path, nil, nil)
		handleResponse(resp, err)
	},
}

var indexCmd = &cobra.Command{
	Use:   "index <paper_id>",
	Short: "Index a paper",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		paperID := parsePaperID(args[0])
		payload := map[string]string{"arxivId": paperID}
		body, _ := json.Marshal(payload)
		resp, err := makeRequest("POST", "/api/papers/index", strings.NewReader(string(body)), map[string]string{"Content-Type": "application/json"})
		handleResponse(resp, err)
	},
}

func init() {
	searchCmd.Flags().IntVar(&limitFlag, "limit", 20, "Number of results")
	dailyCmd.Flags().IntVar(&limitFlag, "limit", 20, "Number of results")
	dailyCmd.Flags().StringVar(&dateFlag, "date", "", "Date in YYYY-MM-DD format")
	dailyCmd.Flags().BoolVar(&trendingFlag, "trending", false, "Sort by trending")
}
