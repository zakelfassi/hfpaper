#!/usr/bin/env python3
"""Generate beautiful SVG terminal screenshots for hfpaper README."""

import html
import os

# Terminal theme (Warm Scholar)
BG = "#2D2A26"
FG = "#E8DFD3"
ACCENT = "#D4622B"
YELLOW = "#E8A87C"
DIM = "#8A7E6E"
GREEN = "#A8C97F"
BLUE = "#7CAFC2"
WHITE = "#F5F0EB"
BOLD_WHITE = "#FFFFFF"

FONT = "JetBrains Mono, Menlo, Monaco, Consolas, monospace"
FONT_SIZE = 13
LINE_HEIGHT = 22
PADDING_X = 24
PADDING_Y = 20
TITLE_BAR_H = 36
DOT_R = 6
DOT_Y = TITLE_BAR_H // 2
DOT_COLORS = ["#FF5F56", "#FFBD2E", "#27C93F"]

def svg_header(width, height, title="hfpaper"):
    total_h = TITLE_BAR_H + PADDING_Y * 2 + height
    dots = ""
    for i, c in enumerate(DOT_COLORS):
        dots += f'<circle cx="{16 + i * 22}" cy="{DOT_Y}" r="{DOT_R}" fill="{c}"/>'
    return f'''<svg xmlns="http://www.w3.org/2000/svg" width="{width}" height="{total_h}" viewBox="0 0 {width} {total_h}">
  <defs>
    <style>
      @import url('https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;700&amp;display=swap');
    </style>
  </defs>
  <rect width="{width}" height="{total_h}" rx="10" fill="{BG}"/>
  <rect width="{width}" height="{TITLE_BAR_H}" rx="10" fill="rgba(0,0,0,0.2)"/>
  <rect y="{TITLE_BAR_H - 10}" width="{width}" height="10" fill="rgba(0,0,0,0.2)"/>
  {dots}
  <text x="{width//2}" y="{DOT_Y + 4}" text-anchor="middle" fill="{DIM}" font-family="{FONT}" font-size="12">{title}</text>
  <g transform="translate({PADDING_X}, {TITLE_BAR_H + PADDING_Y})">'''

def svg_footer():
    return "  </g>\n</svg>"

def text(y, content, color=FG, bold=False, size=None):
    weight = 'font-weight="700"' if bold else ''
    sz = f'font-size="{size}"' if size else f'font-size="{FONT_SIZE}"'
    return f'    <text y="{y}" fill="{color}" font-family="{FONT}" {sz} {weight}>{html.escape(content)}</text>'

def span(y, parts):
    """Render a line with multiple colored spans."""
    x = 0
    elements = []
    for content, color, bold in parts:
        weight = 'font-weight="700"' if bold else ''
        elements.append(f'<tspan x="{x}" fill="{color}" {weight}>{html.escape(content)}</tspan>')
        x += len(content) * 7.8  # approximate char width
    return f'    <text y="{y}" font-family="{FONT}" font-size="{FONT_SIZE}">{"".join(elements)}</text>'

def gen_search():
    """Search command screenshot."""
    lines = []
    y = 0
    
    lines.append(span(y, [("$ ", ACCENT, True), ("hfpaper search ", FG, False), ('"multimodal reasoning"', YELLOW, False), (" --limit 3", FG, False)]))
    y += LINE_HEIGHT * 1.5
    
    lines.append(text(y, "🔍 3 results", WHITE, bold=True))
    y += LINE_HEIGHT * 1.3
    
    # Result 1
    lines.append(text(y, "1. Why Reasoning Matters? A Survey of Advancements in", WHITE, bold=True))
    y += LINE_HEIGHT
    lines.append(text(y, "   Multimodal Reasoning", WHITE, bold=True))
    y += LINE_HEIGHT
    lines.append(span(y, [("   ", FG, False), ("2504.03151", YELLOW, False), (" · 15 ⬆", DIM, False)]))
    y += LINE_HEIGHT
    lines.append(text(y, "   Overview of reasoning techniques in LLMs handling", DIM))
    y += LINE_HEIGHT
    lines.append(text(y, "   both textual and multimodal inputs.", DIM))
    y += LINE_HEIGHT * 1.5
    
    # Result 2
    lines.append(text(y, "2. Multimodal Reasoning for Science: 1st Place Solution", WHITE, bold=True))
    y += LINE_HEIGHT
    lines.append(text(y, "   to the ICML 2025 SeePhys Challenge", WHITE, bold=True))
    y += LINE_HEIGHT
    lines.append(span(y, [("   ", FG, False), ("2509.06079", YELLOW, False), (" · 6 ⬆", DIM, False)]))
    y += LINE_HEIGHT
    lines.append(text(y, "   Caption-assisted reasoning bridges visual and", DIM))
    y += LINE_HEIGHT
    lines.append(text(y, "   textual modalities for top performance.", DIM))
    y += LINE_HEIGHT * 1.5
    
    # Result 3
    lines.append(text(y, "3. MDK12-Bench: Multi-Discipline Benchmark for MLLMs", WHITE, bold=True))
    y += LINE_HEIGHT
    lines.append(span(y, [("   ", FG, False), ("2504.05782", YELLOW, False), (" · 3 ⬆", DIM, False)]))
    y += LINE_HEIGHT
    lines.append(text(y, "   Evaluates multimodal reasoning with diverse", DIM))
    y += LINE_HEIGHT
    lines.append(text(y, "   real-world educational tests.", DIM))
    
    width = 620
    height = y + LINE_HEIGHT
    return svg_header(width, height, "hfpaper — search") + "\n" + "\n".join(lines) + "\n" + svg_footer()

def gen_get():
    """Get command screenshot."""
    lines = []
    y = 0
    
    lines.append(span(y, [("$ ", ACCENT, True), ("hfpaper get ", FG, False), ("2602.08025", YELLOW, False)]))
    y += LINE_HEIGHT * 1.5
    
    lines.append(text(y, "📄 MIND: Benchmarking Memory Consistency and", WHITE, bold=True))
    y += LINE_HEIGHT
    lines.append(text(y, "   Action Control in World Models", WHITE, bold=True))
    y += LINE_HEIGHT * 1.3
    
    fields = [
        ("ID:       ", "2602.08025", YELLOW),
        ("Published:", " Feb 8, 2026", FG),
        ("Authors:  ", " Yixuan Ye, Xuanyu Lu, +8 more", FG),
        ("Upvotes:  ", " 12", FG),
        ("GitHub:   ", " github.com/CSU-JPG/MIND", BLUE),
    ]
    for label, val, color in fields:
        lines.append(span(y, [("  ", FG, False), (label, DIM, False), (val, color, False)]))
        y += LINE_HEIGHT
    
    y += LINE_HEIGHT * 0.5
    lines.append(text(y, "  Summary:", WHITE, bold=True))
    y += LINE_HEIGHT
    lines.append(text(y, "  First open-domain closed-loop benchmark for", DIM))
    y += LINE_HEIGHT
    lines.append(text(y, "  evaluating memory consistency and action control", DIM))
    y += LINE_HEIGHT
    lines.append(text(y, "  in world models.", DIM))
    
    width = 560
    height = y + LINE_HEIGHT
    return svg_header(width, height, "hfpaper — get") + "\n" + "\n".join(lines) + "\n" + svg_footer()

def gen_daily():
    """Daily trending screenshot."""
    lines = []
    y = 0
    
    lines.append(span(y, [("$ ", ACCENT, True), ("hfpaper daily ", FG, False), ("--trending", YELLOW, False), (" --limit 5", FG, False)]))
    y += LINE_HEIGHT * 1.5
    
    lines.append(text(y, "📰 Daily Papers — 5 results", WHITE, bold=True))
    y += LINE_HEIGHT * 1.3
    
    papers = [
        ("1.", "Attention Residuals", "2603.15031", "137", "Kimi Team"),
        ("2.", "TradingAgents: Multi-Agent LLM Trading", "2412.20138", "26", "Yijia Xiao et al."),
        ("3.", "AutoDev: Automated AI-Driven Development", "2403.08299", "12", "Michele Tufano et al."),
        ("4.", "DreamPartGen: Part-Level 3D Generation", "2603.19216", "8", "Tianjiao Yu et al."),
        ("5.", "Temporal Reasoning in LLMs: Tokenisation", "2603.19017", "4", "Gagan Bhatia et al."),
    ]
    
    for num, title, pid, ups, author in papers:
        lines.append(span(y, [(f"  {num} ", DIM, False), (title, WHITE, True)]))
        y += LINE_HEIGHT
        lines.append(span(y, [("     ", FG, False), (pid, YELLOW, False), (f" · {ups} ⬆", DIM, False), (f" · by {author}", DIM, False)]))
        y += LINE_HEIGHT * 1.4
    
    width = 600
    height = y + LINE_HEIGHT * 0.5
    return svg_header(width, height, "hfpaper — daily") + "\n" + "\n".join(lines) + "\n" + svg_footer()

def gen_cite():
    """Cite command screenshot."""
    lines = []
    y = 0
    
    lines.append(span(y, [("$ ", ACCENT, True), ("hfpaper cite ", FG, False), ("2602.08025", YELLOW, False)]))
    y += LINE_HEIGHT * 1.3
    
    bibtex = [
        "@article{ye2026260208025,",
        "  title={MIND: Benchmarking Memory Consistency",
        "         and Action Control in World Models},",
        "  author={Yixuan Ye and Xuanyu Lu and",
        "          Yuxin Jiang and ...},",
        "  journal={arXiv preprint arXiv:2602.08025},",
        "  year={2026}",
        "}",
    ]
    for line in bibtex:
        lines.append(text(y, line, GREEN))
        y += LINE_HEIGHT
    
    y += LINE_HEIGHT * 0.8
    lines.append(span(y, [("$ ", ACCENT, True), ("hfpaper cite ", FG, False), ("2602.08025", YELLOW, False), (" --format apa", FG, False)]))
    y += LINE_HEIGHT * 1.3
    lines.append(text(y, "Ye, Y., Lu, X., Jiang, Y., ... Wang, A. J.", GREEN))
    y += LINE_HEIGHT
    lines.append(text(y, "(2026). MIND: Benchmarking Memory Consistency", GREEN))
    y += LINE_HEIGHT
    lines.append(text(y, "and Action Control. arXiv:2602.08025.", GREEN))
    
    width = 560
    height = y + LINE_HEIGHT
    return svg_header(width, height, "hfpaper — cite") + "\n" + "\n".join(lines) + "\n" + svg_footer()

def gen_mcp():
    """MCP config screenshot."""
    lines = []
    y = 0
    
    lines.append(span(y, [("$ ", ACCENT, True), ("hfpaper mcp", FG, False)]))
    y += LINE_HEIGHT
    lines.append(text(y, "MCP server running on stdio...", DIM))
    y += LINE_HEIGHT * 1.8
    
    lines.append(text(y, "# Claude Desktop / Cursor config:", DIM))
    y += LINE_HEIGHT * 1.3
    
    json_lines = [
        ('{', FG),
        ('  "mcpServers": {', FG),
        ('    "hfpaper": {', FG),
        ('      "command": ', FG),
        ('"npx"', GREEN),
        ('      "args": [', FG),
        ('"-y", "hfpaper", "mcp"', GREEN),
        (']', FG),
        ('    }', FG),
        ('  }', FG),
        ('}', FG),
    ]
    
    config = [
        '{',
        '  "mcpServers": {',
        '    "hfpaper": {',
        '      "command": "npx",',
        '      "args": ["-y", "hfpaper", "mcp"]',
        '    }',
        '  }',
        '}',
    ]
    
    for line in config:
        # Color the string values
        colored = line
        lines.append(text(y, colored, GREEN if '"' in line else FG))
        y += LINE_HEIGHT
    
    y += LINE_HEIGHT * 0.8
    lines.append(text(y, "7 tools · zero install · works everywhere", DIM))
    
    width = 500
    height = y + LINE_HEIGHT
    return svg_header(width, height, "hfpaper — mcp") + "\n" + "\n".join(lines) + "\n" + svg_footer()

# Generate all screenshots
os.makedirs("assets", exist_ok=True)

screenshots = {
    "search": gen_search(),
    "get": gen_get(),
    "daily": gen_daily(),
    "cite": gen_cite(),
    "mcp": gen_mcp(),
}

for name, svg in screenshots.items():
    path = f"assets/{name}.svg"
    with open(path, "w") as f:
        f.write(svg)
    print(f"✅ {path}")
