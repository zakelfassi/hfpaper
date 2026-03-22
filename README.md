# hfpaper

Hugging Face Papers CLI - Fast, agent-consumable access to AI research papers.

## Features

- **Search**: Semantic and full-text search across Hugging Face papers.
- **Read**: Fetch any paper directly as Markdown.
- **Metadata**: Structured JSON for paper details, authors, upvotes, and linked assets.
- **Daily**: Get the daily trending papers feed.
- **Linked Assets**: Find models, datasets, and spaces linked to a paper.
- **Index**: Trigger indexing of new papers on Hugging Face (requires `HF_TOKEN`).

## Installation

### From Source
```bash
git clone https://github.com/zakelfassi/hfpaper
cd hfpaper
make build
# Or install to $GOPATH/bin
make install
```

## Usage

### Examples

```bash
# Search for vision-language papers
hfpaper search "vision language" --limit 5

# Get paper metadata
hfpaper get 2602.08025

# Read paper as markdown (perfect for LLMs)
hfpaper read 2602.08025 > paper.md

# Get daily papers feed
hfpaper daily --trending

# Find models linked to a paper
hfpaper models 2602.08025
```

## Output Formats

`hfpaper` automatically detects if it's being run in a TTY.
- **Human mode** (TTY): Readable text/markdown.
- **Agent mode** (Redirected/Pipe): Raw JSON output.
- Explicit flags: `--json`, `--markdown`, `--table`.

## Authentication

Set the `HF_TOKEN` environment variable for authenticated requests and write operations (like `index`).

```bash
export HF_TOKEN=your_token_here
```

## MCP Integration

This tool is designed to be easily wrapped by Model Context Protocol (MCP) servers. You can call its subcommands from any agentic workflow to fetch research context on the fly.

## Development

- `make build`: Build local binary.
- `make test`: Run tests and vet.
- `make release`: Cross-compile for all targets.
