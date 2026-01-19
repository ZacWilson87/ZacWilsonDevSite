# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Task Board

**Always consult `Todo.md` before starting work.** This file is the project's working memory and task board. Update it when:

- Starting a task (move to "In Progress")
- Completing a task (move to "Completed")
- Discovering new work (add to "Up Next")

## Commands

```bash
# Run the server (default port 8080)
go run ./cmd/web

# Run with custom port
PORT=3000 go run ./cmd/web

# Build binary
go build -o devsite ./cmd/web
```

## Architecture

This is a server-rendered personal developer site built with Go. No database—all content is loaded from markdown files at startup and cached in memory.

### Request Flow

1. `cmd/web/main.go` - Entry point. Creates content loader, loads all markdown, starts HTTP server
2. `internal/content/Loader` - Reads markdown files from `content/` directories at startup, caches parsed items in memory (thread-safe with RWMutex)
3. `internal/server/Server` - Chi router handles requests, renders Go html/templates with content data

### Content Model

Two content types stored as markdown in `content/`:

- `content/projects/*.md` → `/projects` and `/projects/{slug}`
- `content/case-studies/*.md` → `/projects` and `/case-study/{slug}`

Filename (without .md) becomes the URL slug. Frontmatter supports `title`, `date`, `status`, and `tags` fields.

### Templates

Templates use Go's html/template with layout inheritance:

- `templates/layout.html` - Base layout with `{{block}}` for title/content
- `templates/*.html` - Page templates that extend layout
- `templates/partials/*.html` - Shared components (header, nav, footer)

Templates are parsed once at server startup from `templates/` directory.

### Static Assets

Served from `/static/*` → `static/` directory. Vanilla CSS and minimal JS only.

## Design Constraints

From `developer_site_plan.md`:

- Server-rendered first, minimal JavaScript
- No SPA frameworks, animation libraries, or heavy client-side state
- One accent color, neutral palette, generous whitespace
- Content-driven with explicit trade-offs

## Typography

From `typography_plan.md`:

- **Headings/UI:** Inter (400–600)
- **Body text:** Source Serif 4 (400)
- **Scale:** Base ~1.05rem, line-height 1.6–1.7 (body), 1.2–1.3 (headings)

No display fonts, no trendy typefaces. Typography carries hierarchy through spacing, not animation.

## Theme

From `theme.md`. Automatic light/dark mode via `prefers-color-scheme`.

**Light Mode (Refined Paper):** Background #F9FAFB, Text #111827, Accent #2563EB, Muted #6B7280

**Dark Mode (Deep Console):** Background #0F172A, Text #F1F5F9, Accent #38BDF8, Muted #94A3B8
