# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run with hot-reload (development)
air

# Build manually
go build -o ./tmp/main .

# Run directly (no hot-reload)
go run .

# Download dependencies
go mod tidy

# Run tests
go test ./...

# Run a single test
go test ./... -run TestFunctionName
```

The server starts on `localhost:8080`. Air watches `.go`, `.tpl`, `.tmpl`, and `.html` files and rebuilds automatically.

## Architecture

This is a flat Go monolith with three files:

- **`main.go`** — all Gin route handlers, CORS middleware, router setup, DB init, and the file-upload/Excel-parsing logic
- **`pkg/structs.go`** — all domain types and DTOs; only `Ingredient` is persisted via GORM
- **`database/database.go`** — exposes a single global `database.DBCon *gorm.DB`

**Stack:** Gin (HTTP router) + GORM (ORM) + SQLite (`techsheets.db` at repo root)

**Database access pattern:** handlers use GORM's generic API — `gorm.G[pkg.Ingredient](database.DBCon).Where(...).First(ctx)` — for typed queries, or `database.DBCon.Find(...)` for bulk fetches.

**File uploads:** Excel files (`.xlsx`) are saved to `./files/` then parsed with `excelize` to bulk-create Ingredients row by row.

**CORS:** A wildcard `CORSMiddleware` is applied globally before all routes.

## Current API Surface

| Method | Path | Description |
|--------|------|-------------|
| POST | `/ingredient` | Create ingredient |
| GET | `/ingredients` | List all ingredients |
| GET | `/ingredient/:id` | Get ingredient by ID |
| PATCH | `/ingredient/:id` | Update ingredient |
| DELETE | `/ingredient/:id` | Delete ingredient |
| POST | `/ingredients/upload` | Bulk import from Excel file (multipart `file` field) |

## Domain Model Notes

`pkg/structs.go` defines many types (`Recipe`, `TechnicalSheet`, `Step`, `Mold`, etc.) that are **not yet wired to the database or any routes** — they represent the intended future domain. Only `Ingredient` (which embeds `gorm.Model`) is currently active. Do not assume other types are persisted or have endpoints.

## API Testing

A Bruno collection lives in `techsheets api bruno/` and can be used to test endpoints manually.
