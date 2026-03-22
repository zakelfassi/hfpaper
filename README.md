# hfpaper

AI research papers from your terminal. Search, read, and explore the full Hugging Face Papers ecosystem from the command line.

**[zakelfassi.github.io/hfpaper](https://zakelfassi.github.io/hfpaper)** · **[GitHub](https://github.com/zakelfassi/hfpaper)**

## Install

```bash
# From source
go install github.com/zakelfassi/hfpaper@latest

# Or clone and build
git clone https://github.com/zakelfassi/hfpaper
cd hfpaper && make build
```

Prebuilt binaries for Linux, macOS, and Windows are available on the [Releases](https://github.com/zakelfassi/hfpaper/releases) page.

## Usage

```bash
# Search for papers
hfpaper search "vision language models" --limit 5

# Get structured metadata (authors, abstract, GitHub, upvotes)
hfpaper get 2602.08025

# Read the full paper as markdown
hfpaper read 2602.08025

# Today's trending papers
hfpaper daily --trending

# Find linked models, datasets, or spaces
hfpaper models 2602.08025
hfpaper datasets 2602.08025
hfpaper spaces 2602.08025

# Index a new paper on HF (requires HF_TOKEN)
hfpaper index 2503.12345
```

## Paper ID Formats

hfpaper accepts any of these and extracts the arXiv ID automatically:

| Input | Parsed ID |
|-------|-----------|
| `2602.08025` | `2602.08025` |
| `2602.08025v1` | `2602.08025v1` |
| `https://huggingface.co/papers/2602.08025` | `2602.08025` |
| `https://huggingface.co/papers/2602.08025.md` | `2602.08025` |
| `https://arxiv.org/abs/2602.08025` | `2602.08025` |
| `https://arxiv.org/pdf/2602.08025` | `2602.08025` |

## Output Formats

hfpaper auto-detects context:

- **TTY** (interactive terminal): human-readable text/markdown
- **Piped/redirected** (non-interactive): raw JSON for parsing

Override with flags: `--json`, `--table`, `--markdown`

```bash
# Pipe JSON to jq
hfpaper search "RLHF" --json | jq '.[0].paper.title'

# Save a paper as markdown
hfpaper read 2602.08025 > paper.md

# Feed a paper to an LLM
hfpaper read 2602.08025 | llm "summarize this paper in 3 bullets"
```

## Use with AI Agents

hfpaper is designed to be consumed by AI coding agents (Claude, Codex, Cursor, etc). Drop [`AGENTS.md`](./AGENTS.md) into your project or reference it in your agent's skill/tool config.

### Claude Code / Codex

Add to your project's `AGENTS.md` or system prompt:

```
You have access to `hfpaper`, a CLI for AI research papers.
- Search: `hfpaper search "<query>" --json`
- Read: `hfpaper read <arxiv_id>`
- Trending: `hfpaper daily --trending --json`
- Metadata: `hfpaper get <arxiv_id> --json`
- Models: `hfpaper models <arxiv_id> --json`
```

### OpenClaw / Custom Agents

Copy `AGENTS.md` to your agent's workspace. It contains the full command reference, paper ID format docs, and usage examples optimized for agent consumption.

### MCP Server

`hfpaper mcp` starts a Model Context Protocol server over stdio, exposing all commands as tools. Works with Claude Desktop, Cursor, Windsurf, and any MCP-compatible client.

**Claude Desktop** — add to `claude_desktop_config.json`:
```json
{
  "mcpServers": {
    "hfpaper": {
      "command": "hfpaper",
      "args": ["mcp"]
    }
  }
}
```

**Cursor** — add to `.cursor/mcp.json`:
```json
{
  "mcpServers": {
    "hfpaper": {
      "command": "hfpaper",
      "args": ["mcp"]
    }
  }
}
```

**Available MCP tools:** `search_papers`, `get_paper`, `read_paper`, `daily_papers`, `paper_models`, `paper_datasets`, `paper_spaces`

## Authentication

Most commands work without authentication. Set `HF_TOKEN` for write operations:

```bash
export HF_TOKEN=hf_xxxxx
hfpaper index 2503.12345
```

## Development

```bash
make build      # Build local binary
make install    # Install to $GOPATH/bin
make test       # Run tests + vet
make release    # Cross-compile (linux/mac/windows, amd64/arm64)
```

## License

MIT — see [LICENSE](./LICENSE)

---

Built by [Zak El Fassi](https://github.com/zakelfassi)
