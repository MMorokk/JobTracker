package main

import (
	"fmt"
	"os"
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

func initialModel() model {
	rows := [][]string{
		{"1", "Acme Corp", "Backend Dev", "Applied"},
		{"2", "Globex", "SRE", "Interview"},
		{"3", "Initech", "DevOps Engineer", "Applied"},
		{"4", "Umbrella Ltd", "Platform Engineer", "Rejected"},
		{"5", "Soylent Co", "Software Engineer", "Offer"},
		{"6", "Hooli", "Staff Engineer", "Applied"},
		{"7", "Pied Piper", "Backend Engineer", "Interview"},
		{"8", "Dunder Mifflin", "Full Stack Dev", "Applied"},
		{"9", "Vandelay Tech", "Cloud Architect", "Rejected"},
		{"10", "Gekko & Co", "Site Reliability Engineer", "Interview"},
		{"11", "Initrode", "Go Developer", "Applied"},
		{"12", "Oscorp", "Systems Engineer", "Offer"},
		{"13", "Massive Dynamic", "Infrastructure Engineer", "Applied"},
		{"14", "Weyland Corp", "Backend Developer", "Ghosted"},
		{"15", "Cyberdyne Sys", "Platform Dev", "Rejected"},
		{"16", "Tyrell Corp", "Software Architect", "Interview"},
		{"17", "Stark Industries", "API Engineer", "Applied"},
		{"18", "Wayne Enterprises", "DevOps Lead", "Offer"},
		{"19", "Pendant Publishing", "Backend Dev", "Rejected"},
		{"20", "Veridian Dynamics", "SRE", "Applied"},
		{"21", "Rekall Inc", "Cloud Engineer", "Ghosted"},
		{"22", "Momcorp", "Go Engineer", "Interview"},
		{"23", "Aperture Science", "Software Engineer", "Applied"},
		{"24", "Black Mesa", "Systems Dev", "Rejected"},
		{"25", "Bluth Company", "Backend Engineer", "Interview"},
		{"26", "Sterling Archer", "Platform Engineer", "Applied"},
		{"27", "Soylent Systems", "SRE Lead", "Offer"},
		{"28", "Lacroix Tech", "Infrastructure Dev", "Applied"},
		{"29", "Delos Inc", "Backend Dev", "Ghosted"},
		{"30", "Tesslamax", "Cloud Architect", "Interview"},
		{"31", "Initech EMEA", "DevOps Engineer", "Applied"},
		{"32", "Goliath National", "Backend Developer", "Rejected"},
		{"33", "Altair Nanotech", "Systems Engineer", "Applied"},
		{"34", "Gringotts Digital", "Staff Backend Dev", "Interview"},
		{"35", "Abstergo Ltd", "Platform Dev", "Offer"},
		{"36", "Netlink Corp", "Cloud Engineer", "Applied"},
		{"37", "Prestige Global", "API Engineer", "Rejected"},
		{"38", "Monarch Solutions", "Backend Engineer", "Applied"},
		{"39", "Krusty Krab Tech", "SRE", "Ghosted"},
		{"40", "Nakatomi Corp", "DevOps Engineer", "Interview"},
		{"41", "Zorg Industries", "Software Architect", "Applied"},
		{"42", "Planet Express", "Go Developer", "Offer"},
		{"43", "Multivac Systems", "Backend Dev", "Applied"},
		{"44", "Sombra Corp", "Infrastructure Engineer", "Rejected"},
		{"45", "Primatech Paper", "Platform Engineer", "Interview"},
		{"46", "Virtucon Ltd", "SRE", "Applied"},
		{"47", "Consolidated", "Backend Engineer", "Ghosted"},
		{"48", "Roxxon Energy", "Cloud Dev", "Applied"},
		{"49", "Frobozz Magic", "Software Engineer", "Interview"},
		{"50", "Aperture Labs", "Staff SRE", "Offer"},
	}

	cols := []Column{
		{Title: "ID", Width: 4},
		{Title: "Company", Width: 24},
		{Title: "Role", Width: 40},
		{Title: "Status", Width: 12},
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
		case "1", "2", "3", "4": // sort by column number
			col := int(msg.String()[0] - '1')
			m.table.SortByCol(col)
			//case "enter":
			//row := m.table.SelectedRow()
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
	helpstring := "press q to quit | j/k to navigate | 1-4 to sort by column | enter to select row"
	footer := style.Width(m.width).Render(helpstring)

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
	m := initialModel()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
