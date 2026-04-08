// Package tui wires together the top-level bubbletea program.
// mainModel owns the active page and delegates all bubbletea lifecycle
// calls (Init/Update/View) to it, enabling page switching without the
// program needing to know which page is currently active.
package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/MMorokk/JobTracker/internal/db"
	"github.com/MMorokk/JobTracker/internal/tui/pages"
)

// mainModel is the root bubbletea model. It holds the currently active page
// and forwards all messages to it.
type mainModel struct {
	page tea.Model
}

// InitialModel creates the root model with the job table as the starting page.
// No database query is performed here; the table triggers it via Init.
func InitialModel(db *db.DB) *mainModel {
	cols := []pages.Column{
		{Title: "ID", Width: 4},
		{Title: "Role", Width: 40},
		{Title: "Company", Width: 24},
		{Title: "Location", Width: 40},
		{Title: "Status", Width: 10},
	}
	table := pages.NewTable(db, cols, 120, 30)
	return &mainModel{page: table}
}

func (m *mainModel) Init() tea.Cmd {
	return m.page.Init()
}

func (m *mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	page, cmd := m.page.Update(msg)
	m.page = page
	return m, cmd
}

func (m *mainModel) View() tea.View {
	return m.page.View()
}
