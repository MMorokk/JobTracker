```
██████╗ ██████╗  ██████╗ ██╗  ██╗███████╗███╗   ██╗
██╔══██╗██╔══██╗██╔═══██╗██║ ██╔╝██╔════╝████╗  ██║
██████╔╝██████╔╝██║   ██║█████╔╝ █████╗  ██╔██╗ ██║
██╔══██╗██╔══██╗██║   ██║██╔═██╗ ██╔══╝  ██║╚██╗██║
██████╔╝██║  ██║╚██████╔╝██║  ██╗███████╗██║ ╚████║
╚═════╝ ╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝╚══════╝╚═╝  ╚═══╝

 █████╗ ████████╗    ████████╗██╗  ██╗██╗███████╗
██╔══██╗╚══██╔══╝    ╚══██╔══╝██║  ██║██║██╔════╝
███████║   ██║          ██║   ███████║██║███████╗
██╔══██║   ██║          ██║   ██╔══██║██║╚════██║
██║  ██║   ██║          ██║   ██║  ██║██║███████║
╚═╝  ╚═╝   ╚═╝          ╚═╝   ╚═╝  ╚═╝╚═╝╚══════╝

███╗   ███╗ ██████╗ ███╗   ███╗███████╗███╗   ██╗████████╗
████╗ ████║██╔═══██╗████╗ ████║██╔════╝████╗  ██║╚══██╔══╝
██╔████╔██║██║   ██║██╔████╔██║█████╗  ██╔██╗ ██║   ██║
██║╚██╔╝██║██║   ██║██║╚██╔╝██║██╔══╝  ██║╚██╗██║   ██║
██║ ╚═╝ ██║╚██████╔╝██║ ╚═╝ ██║███████╗██║ ╚████║   ██║
╚═╝     ╚═╝ ╚═════╝ ╚═╝     ╚═╝╚══════╝╚═╝  ╚═══╝   ╚═╝
```

> **This project is currently broken and has not yet reached a first fully working version.**
> The code compiles but is incomplete — the database is not wired up, features are stubbed out,
> and the app uses hardcoded dummy data instead of real persistence.

---

# JobTracker

A terminal UI (TUI) app for tracking job applications, written in Go. Paste a job posting URL and let it automatically fill in the details for you.

## Planned Features

- **TUI table view** — browse all your applications in a scrollable, sortable table right in your terminal
- **URL autofill** — paste a job posting URL and have the app scrape and parse it automatically using a headless browser (Chromium via `chromedp`)
- **Local LLM extraction** — job page text is sent to a locally running [Ollama](https://ollama.com) model, which extracts structured fields (title, company, location, type, salary, requirements, etc.) without any data leaving your machine
- **SQLite persistence** — all applications stored locally in a single SQLite file via `modernc.org/sqlite` (no CGO required)
- **Application statuses** — track where each application stands: Applied, Interview, Offer, Rejected, Ghosted

## Tech Stack

| Concern | Library |
|---|---|
| TUI framework | [bubbletea v2](https://charm.land/bubbletea) |
| Styling | [lipgloss](https://github.com/charmbracelet/lipgloss) |
| JS-rendered scraping | [chromedp](https://github.com/chromedp/chromedp) |
| LLM inference | [ollama-go](https://github.com/ollama/ollama) (local models) |
| Database | [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) |

## Keybindings (current)

| Key | Action |
|---|---|
| `j` / `↓` | Move cursor down |
| `k` / `↑` | Move cursor up |
| `1`–`4` | Sort by column (ID, Company, Role, Status) |
| `q` / `Ctrl+C` | Quit |

## Prerequisites

- Go 1.21+
- [Ollama](https://ollama.com) running locally with at least one model pulled (e.g. `ollama pull llama3`)
- Chromium or Google Chrome installed (used by `chromedp` for headless scraping)

## Running

```sh
go run .
```