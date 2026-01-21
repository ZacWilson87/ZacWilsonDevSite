# Zac Wilson Dev Site

A server-rendered personal developer portfolio site built with Go. Fast, simple, and content-driven—no database required.

## Overview

This is a minimal, performance-focused personal site that prioritizes server-side rendering and content over complexity. All content is loaded from markdown files at startup and cached in memory for instant response times.

## Features

- **Server-rendered** - Pure Go templates with layout inheritance
- **Markdown content** - Write projects and case studies in markdown with YAML frontmatter
- **In-memory caching** - Content loaded once at startup for zero-latency responses
- **Automatic theming** - Light/dark mode via `prefers-color-scheme`
- **Zero dependencies** - Single binary deployment, no database needed
- **Responsive design** - Clean, minimal interface optimized for readability

## Tech Stack

- **Go** - Server, routing, templating
- **Chi** - HTTP router
- **Goldmark** - Markdown parsing
- **HTML/CSS** - Semantic markup with vanilla CSS
- **Inter + Source Serif 4** - Typography stack

## Getting Started

### Prerequisites

- Go 1.21 or higher

### Installation

```bash
# Clone the repository
git clone https://github.com/ZacWilson87/ZacWilsonDevSite.git
cd ZacWilsonDevSite

# Install dependencies
go mod download
```

### Running the Server

```bash
# Run with default port (8080)
go run ./cmd/web

# Run with custom port
PORT=3000 go run ./cmd/web
```

Visit `http://localhost:8080` in your browser.

### Building for Production

```bash
# Build binary
go build -o devsite ./cmd/web

# Run the binary
./devsite

# Or with custom port
PORT=3000 ./devsite
```

## Project Structure

```
.
├── cmd/
│   └── web/           # Application entry point
├── internal/
│   ├── content/       # Content loader and models
│   └── server/        # HTTP server and handlers
├── templates/         # Go html/template files
│   ├── partials/      # Reusable components (header, nav, footer)
│   └── *.html         # Page templates
├── static/            # Static assets (CSS, images, JS)
├── content/           # Markdown content
│   ├── projects/      # Project pages
│   └── case-studies/  # Detailed case study pages
├── CLAUDE.md          # Project guidance for Claude Code
└── Todo.md            # Task board and working memory
```

## Content Management

### Adding Projects

Create a markdown file in `content/projects/`:

```markdown
---
title: Your Project Name
date: 2025-01-15
status: Open Source
tags: [go, typescript, api]
---

Your project description and content here...
```

The filename (without `.md`) becomes the URL slug: `content/projects/my-project.md` → `/projects/my-project`

### Adding Case Studies

Create a markdown file in `content/case-studies/`:

```markdown
---
title: "Project Name: Detailed Case Study"
date: 2025-01-15
status: Production
tags: [backend, infrastructure]
---

Detailed case study content...
```

Accessed via: `/case-study/{slug}`

### Frontmatter Fields

- `title` (required) - Display title
- `date` (optional) - Publication date (YYYY-MM-DD)
- `status` (optional) - Project status (e.g., "Open Source", "Production")
- `tags` (optional) - Array of technology tags

## Development Workflow

This project uses `Todo.md` as the working task board. Before starting work:

1. Check `Todo.md` for current tasks
2. Move items to "In Progress" when starting
3. Add new discoveries to "Up Next"
4. Move completed items to "Completed"

See `CLAUDE.md` for detailed development guidelines.

## Architecture

### Request Flow

1. **Startup** - `cmd/web/main.go` creates content loader, parses all markdown, starts server
2. **Content Loading** - `internal/content/Loader` reads markdown files, caches parsed items with thread-safe RWMutex
3. **Request Handling** - `internal/server/Server` uses Chi router to handle requests, renders templates with cached content

### Template System

Templates use Go's `html/template` with layout inheritance:

- `templates/layout.html` - Base layout with `{{block}}` definitions
- Page templates extend the layout with `{{define "title"}}` and `{{define "content"}}`
- Partials (header, nav, footer) are reusable components

Templates are parsed once at startup for optimal performance.

## Design Philosophy

From the project's design constraints:

- **Server-rendered first** - Minimal JavaScript, no SPA frameworks
- **Content-driven** - Typography and spacing over animation
- **Explicit trade-offs** - Simple and fast over complex and flexible
- **No over-engineering** - Build what's needed, nothing more

## License

MIT License - See LICENSE file for details

## Contact

Zac Wilson - [https://zacwilson.dev](https://zacwilson.dev)
