---
title: "SeedGenius: AI-Powered Test Data Generation"
date: 2025-08-12
status: Open Source
tags: [go, ai, cli, postgres, developer-tools]
---

Every developer knows the pain: you need test data, but creating it manually is tedious, and copying production data is a security risk. SeedGenius solves this by generating realistic, constraint-aware mock data using AI.

## The Problem

Setting up test databases is a recurring friction point in development workflows. The options are typically:

1. **Manual creation** — Time-consuming and produces unrealistic patterns
2. **Production snapshots** — Security and privacy concerns, plus stale data
3. **Static fixtures** — Brittle, don't adapt to schema changes
4. **Generic fakers** — Ignore foreign keys and unique constraints, causing insert failures

I needed something that understood database relationships and could generate data that *actually works* when inserted.

## The Approach

SeedGenius takes a different angle: let AI understand your schema and generate contextually appropriate data.

The core workflow:

```
seedgenius schema          # Inspect tables interactively
seedgenius generate --tables=users,orders --rows=50
```

The tool reads your PostgreSQL schema, identifies relationships and constraints, then uses GPT-4 to generate data that respects:

- **Foreign key relationships** — Orders reference real user IDs
- **Unique constraints** — No duplicate emails or usernames
- **Data types** — Proper formats for dates, JSON fields, enums
- **Semantic context** — User names look like names, not random strings

## Technical Decisions

**Go for the CLI.** Fast startup, single binary distribution, excellent PostgreSQL driver support. No runtime dependencies for end users.

**Bubble Tea for the interface.** The schema browser needed to be interactive but stay in the terminal. Bubble Tea's Elm-inspired architecture made building a responsive TUI straightforward.

**Batch operations with dependency ordering.** When generating for multiple tables, SeedGenius topologically sorts based on foreign keys. Parent tables populate first, so child records can reference valid IDs.

**Local connection storage.** Database credentials stay on your machine in a SQLite file. No cloud sync, no config servers. Run `seedgenius connections add` once, then just reference by name.

## Export Flexibility

Beyond direct database insertion, SeedGenius exports to JSON and CSV:

```
seedgenius export --tables=users,orders --format=json --output=./fixtures
```

This covers CI pipelines, frontend development with mock APIs, and sharing test datasets across teams.

## What I Learned

**AI prompting for structured output is tricky.** Early versions produced data that *looked* right but violated constraints. The solution was feeding the AI more schema context—column types, nullable flags, and explicit constraint descriptions—rather than just table names.

**Terminal UIs have surprising depth.** Bubble Tea forced me to think about state management in a way web frameworks abstract away. The result is snappier than I expected.

**Developer tools need to be fast.** If generating 10 rows takes 30 seconds, nobody will use it. Batching OpenAI requests and parallelizing inserts where safe cut generation time significantly.

## Try It

SeedGenius is open source and available on GitHub:

```bash
git clone https://github.com/ZacWilson87/SeedGeniusCLI
cd SeedGeniusCLI
go mod tidy
go build -o seedgenius ./cmd/seedgenius
```

Point it at a PostgreSQL database, add your OpenAI key, and generate some data. Contributions welcome.
