# AGENTS.md — hfpaper

You have access to `hfpaper`, a CLI for searching, reading, and exploring AI research papers from Hugging Face and arXiv.

## Install

```bash
go install github.com/zakelfassi/hfpaper@latest
```

Or download a prebuilt binary from [Releases](https://github.com/zakelfassi/hfpaper/releases).

## Commands

```bash
hfpaper search <query> [--limit N]         # Semantic + full-text search
hfpaper get <paper_id>                      # Structured metadata (JSON)
hfpaper read <paper_id>                     # Full paper as markdown
hfpaper daily [--date YYYY-MM-DD] [--trending] [--limit N]  # Daily papers
hfpaper models <paper_id>                   # HF models linked to paper
hfpaper datasets <paper_id>                 # HF datasets linked to paper
hfpaper spaces <paper_id>                   # HF spaces linked to paper
hfpaper index <paper_id>                    # Index a paper (needs HF_TOKEN)
```

## Paper ID formats

All commands accept any of these as `<paper_id>`:
- `2602.08025` (arXiv ID)
- `2602.08025v1` (versioned)
- `https://huggingface.co/papers/2602.08025`
- `https://arxiv.org/abs/2602.08025`
- `https://arxiv.org/pdf/2602.08025`

## Output

- `--json` — raw JSON (default when piped / non-TTY)
- `--table` — human-readable table
- `--markdown` — markdown output

When running non-interactively (piped or redirected), output defaults to JSON for easy parsing.

## When to use

- User asks about a specific paper → `hfpaper get <id>` then `hfpaper read <id>`
- User asks "what's new in AI" → `hfpaper daily --trending --limit 10`
- User asks to find papers on a topic → `hfpaper search "<topic>" --limit 10`
- User asks what models implement a paper → `hfpaper models <id>`
- You need to summarize or analyze a paper → `hfpaper read <id> | head -500`

## Examples

```bash
# Find papers about multimodal reasoning
hfpaper search "multimodal reasoning" --limit 5 --json

# Read a paper and pipe to yourself for analysis
hfpaper read 2602.08025

# What's trending today
hfpaper daily --trending --limit 10 --json

# Check if a paper has linked model weights
hfpaper models 2602.08025 --json

# Get metadata (authors, GitHub, upvotes, abstract)
hfpaper get 2602.08025 --json
```

## Tips

- `read` output can be long (full paper). Use `head -200` or `--limit` if you just need the abstract and intro.
- `search` uses HF's hybrid semantic + full-text search — natural language queries work well.
- `daily` without `--date` returns today's papers. Add `--trending` to sort by upvotes.
- Chain commands: `hfpaper search "RL alignment" --json | jq '.[0].paper.id'` → get the top result's ID → `hfpaper read <id>`
