package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	table  Table
	width  int
	height int
}

var style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4"))

func initialModel(db *sql.DB) model {
	rows, err := tableViewQuery(db)
	if err != nil {
		log.Fatal("Error querying database:", err)
	}
	cols := []Column{
		{Title: "ID", Width: 4},
		{Title: "Role", Width: 40},
		{Title: "Company", Width: 24},
		{Title: "Location", Width: 40},
		{Title: "Status", Width: 10},
	}

	t := NewTable(cols, rows, 80, 20)
	return model{table: t, width: 80, height: 20}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "j", "down":
			m.table.MoveDown()
		case "k", "up":
			m.table.MoveUp()
			//case "enter":
			//	row := m.table.SelectedRow()
			// do something with row...
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.table.SetSize(msg.Width, msg.Height-3) // -1 footer, -2 header+separator
	}

	return m, cmd
}

func (m model) View() tea.View {
	helpString := "press q to quit | j/k to navigate | enter to select row"
	footer := style.Width(m.width).Render(helpString)

	tableView := m.table.View()
	usedLines := strings.Count(tableView, "\n") + 1 // +1 for footer line
	padding := m.height - usedLines
	if padding < 0 {
		padding = 0
	}

	v := tea.NewView(tableView + strings.Repeat("\n", padding) + footer)
	v.AltScreen = true
	return v
}

func main() {

	dataBase, err := NewDatabase("./jobtracker.db")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println("Error closing database:", err)
		}
	}(dataBase.db)

	m := initialModel(dataBase.db)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		log.Fatal("Error running program:", err)
	}
}
