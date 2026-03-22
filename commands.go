package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var formatFlag string

var citeCmd = &cobra.Command{
	Use:   "cite <paper_id>",
	Short: "Generate a citation for a paper",
	Long:  "Generate a citation in BibTeX, APA, or MLA format",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		paperID := parsePaperID(args[0])
		path := fmt.Sprintf("/api/papers/%s", paperID)
		resp, err := makeRequest("GET", path, nil, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		var paper map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&paper); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing response: %v\n", err)
			os.Exit(1)
		}

		title, _ := paper["title"].(string)
		publishedAt, _ := paper["publishedAt"].(string)

		// Parse authors
		var authorNames []string
		if authors, ok := paper["authors"].([]interface{}); ok {
			for _, a := range authors {
				if author, ok := a.(map[string]interface{}); ok {
					if name, ok := author["name"].(string); ok {
						authorNames = append(authorNames, name)
					}
				}
			}
		}

		// Parse year
		year := "2026"
		if t, err := time.Parse(time.RFC3339, publishedAt); err == nil {
			year = fmt.Sprintf("%d", t.Year())
		}

		switch formatFlag {
		case "apa":
			// APA format
			var authorStr string
			if len(authorNames) > 0 {
				parts := make([]string, 0, len(authorNames))
				for _, name := range authorNames {
					nameParts := strings.Fields(name)
					if len(nameParts) >= 2 {
						last := nameParts[len(nameParts)-1]
						initials := ""
						for _, p := range nameParts[:len(nameParts)-1] {
							initials += string(p[0]) + ". "
						}
						parts = append(parts, fmt.Sprintf("%s, %s", last, strings.TrimSpace(initials)))
					} else {
						parts = append(parts, name)
					}
				}
				if len(parts) > 7 {
					authorStr = strings.Join(parts[:6], ", ") + ", ... " + parts[len(parts)-1]
				} else if len(parts) > 1 {
					authorStr = strings.Join(parts[:len(parts)-1], ", ") + ", & " + parts[len(parts)-1]
				} else {
					authorStr = parts[0]
				}
			}
			fmt.Printf("%s (%s). %s. arXiv preprint arXiv:%s.\n", authorStr, year, title, paperID)

		case "mla":
			var authorStr string
			if len(authorNames) > 0 {
				if len(authorNames) == 1 {
					authorStr = authorNames[0]
				} else if len(authorNames) == 2 {
					authorStr = authorNames[0] + ", and " + authorNames[1]
				} else {
					authorStr = authorNames[0] + ", et al"
				}
			}
			fmt.Printf("%s. \"%s.\" arXiv preprint arXiv:%s (%s).\n", authorStr, title, paperID, year)

		default: // bibtex
			firstAuthor := "unknown"
			if len(authorNames) > 0 {
				parts := strings.Fields(authorNames[0])
				if len(parts) > 0 {
					firstAuthor = strings.ToLower(parts[len(parts)-1])
				}
			}
			citeKey := fmt.Sprintf("%s%s%s", firstAuthor, year, strings.ReplaceAll(paperID, ".", ""))
			authorBibtex := strings.Join(authorNames, " and ")
			fmt.Printf("@article{%s,\n  title={%s},\n  author={%s},\n  journal={arXiv preprint arXiv:%s},\n  year={%s}\n}\n",
				citeKey, title, authorBibtex, paperID, year)
		}
	},
}

var summaryCmd = &cobra.Command{
	Use:   "summary <paper_id>",
	Short: "Get AI-generated summary of a paper",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		paperID := parsePaperID(args[0])
		path := fmt.Sprintf("/api/papers/%s", paperID)
		resp, err := makeRequest("GET", path, nil, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		var paper map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&paper); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing response: %v\n", err)
			os.Exit(1)
		}

		if jsonFlag {
			out := map[string]interface{}{
				"id":         paperID,
				"title":      paper["title"],
				"ai_summary": paper["ai_summary"],
			}
			b, _ := json.MarshalIndent(out, "", "  ")
			fmt.Println(string(b))
		} else {
			title, _ := paper["title"].(string)
			summary, _ := paper["ai_summary"].(string)
			if summary == "" {
				summary, _ = paper["summary"].(string)
			}
			fmt.Printf("%s\n\n%s\n", title, summary)
		}
	},
}

var openCmd = &cobra.Command{
	Use:   "open <paper_id>",
	Short: "Open paper in browser",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		paperID := parsePaperID(args[0])
		url := fmt.Sprintf("https://huggingface.co/papers/%s", paperID)

		var openErr error
		switch runtime.GOOS {
		case "darwin":
			openErr = exec.Command("open", url).Start()
		case "windows":
			openErr = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
		default:
			openErr = exec.Command("xdg-open", url).Start()
		}

		if openErr != nil {
			fmt.Fprintf(os.Stderr, "Error opening browser: %v\nURL: %s\n", openErr, url)
			os.Exit(1)
		}
		fmt.Printf("Opening %s\n", url)
	},
}

func init() {
	citeCmd.Flags().StringVar(&formatFlag, "format", "bibtex", "Citation format: bibtex, apa, mla")
}
