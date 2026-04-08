// JobTracker is a terminal UI for tracking job applications.
// It stores job postings in a local SQLite database and provides
// a keyboard-driven interface for browsing and editing them.
package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/MMorokk/JobTracker/internal/db"
	"github.com/MMorokk/JobTracker/internal/tui"
)

func main() {
	database, err := db.NewDB("jobs.db")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open database:", err)
		os.Exit(1)
	}

	p := tea.NewProgram(tui.InitialModel(database))
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error running program:", err)
		os.Exit(1)
	}
}