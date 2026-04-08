// Package pages contains the individual TUI pages (table, details, edit).
// Each page implements tea.Model so the root model can swap between them.
package pages

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/MMorokk/JobTracker/internal/db"
	"github.com/charmbracelet/lipgloss"
)

// RowsLoadedMsg is delivered to Update when the async DB query in Init completes.
type RowsLoadedMsg struct {
	Rows [][]string
	Err  error
}

// Column describes a single table column.
type Column struct {
	Title string
	Width int // computed by UpdateColumnWidths; initial value is used as a hint only
}

// Table is the job-listing page. It fetches rows asynchronously on Init
// and supports keyboard navigation.
type Table struct {
	db         *db.DB
	rows       [][]string
	cols       []Column
	cursor     int
	height     int // full terminal height
	width      int // full terminal width
	offset     int // first visible row index (for scrolling)
	rowsHeight int // height minus reserved lines (header, separator, footer)
}

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("63"))

	countStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("21"))

	normalStyle = lipgloss.NewStyle()

	helpStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("220"))
)

// NewTable creates a Table. width and height should match the terminal size;
// they will be updated when the first WindowSizeMsg arrives.
func NewTable(db *db.DB, cols []Column, width, height int) *Table {
	return &Table{
		db:         db,
		cols:       cols,
		width:      width,
		height:     height,
		rowsHeight: max(height-6, 0),
	}
}

/*
UpdateColumnWidths recalculates column widths based on current content and
terminal width. Column 0 (ID) gets only its natural content width plus
padding. The remaining space is distributed evenly among the other columns.

Each row inside the box occupies: 2 (cursor "› "/spaces) + sum(col widths)
+ (numCols-1)*2 (the "| " separators between columns). That total must fit
within innerWidth = t.width - 2 (the two │ border characters).
*/
func (t *Table) UpdateColumnWidths() {
	if len(t.cols) == 0 {
		return
	}

	// min width for each column: max of title length and longest cell
	minWidths := make([]int, len(t.cols))
	for i, col := range t.cols {
		minWidths[i] = max(len(col.Title), findLargestRowInColumn(i, t.rows))
	}

	const padding = 2

	// col 0 (ID) stays at its natural min width plus padding
	t.cols[0].Width = minWidths[0] + padding

	// innerWidth = t.width - 2 for the two │ border chars
	// per-row overhead = 2 (cursor prefix) + (numCols-1)*2 (separators)
	innerWidth := t.width - 2
	overhead := 2 + (len(t.cols)-1)*2
	remaining := innerWidth - overhead - t.cols[0].Width

	otherCount := len(t.cols) - 1
	minSum := 0
	for i := 1; i < len(t.cols); i++ {
		minSum += minWidths[i] + padding
	}
	extra := 0
	if remaining > minSum && otherCount > 0 {
		extra = (remaining - minSum) / otherCount
	}
	for i := 1; i < len(t.cols); i++ {
		t.cols[i].Width = minWidths[i] + padding + extra
	}
}

func findLargestRowInColumn(col int, rows [][]string) int {
	maxWidth := 0
	for _, row := range rows {
		if len(row[col]) > maxWidth {
			maxWidth = len(row[col])
		}
	}
	return maxWidth
}

func (t *Table) MoveUp() {
	if t.cursor > 0 {
		t.cursor--
		if t.cursor < t.offset {
			t.offset--
		}
	}
}

func (t *Table) MoveDown() {
	if t.cursor < len(t.rows)-1 {
		t.cursor++
		if t.cursor >= t.offset+t.rowsHeight {
			t.offset++
		}
	}
}

func (t *Table) SetSize(width, height int) {
	t.height = height
	t.width = width
	t.rowsHeight = max(height-6, 0)
}

// SelectedRow returns the data for the currently highlighted row, or nil if
// the table is empty.
func (t *Table) SelectedRow() []string {
	if len(t.rows) == 0 {
		return nil
	}
	return t.rows[t.cursor]
}

func (t *Table) Init() tea.Cmd {
	return func() tea.Msg {
		rows, err := t.db.GetJobs()
		return RowsLoadedMsg{Rows: rows, Err: err}
	}
}

func (t *Table) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case RowsLoadedMsg:
		if msg.Err == nil {
			t.rows = msg.Rows
			t.UpdateColumnWidths()
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "k", "up":
			t.MoveUp()
		case "j", "down":
			t.MoveDown()
		case "q", "ctrl+c":
			return t, tea.Quit
		}
	case tea.WindowSizeMsg:
		t.SetSize(msg.Width, msg.Height)
		t.UpdateColumnWidths()
	}
	return t, cmd
}

func (t *Table) View() tea.View {
	var b strings.Builder

	// ── Title bar ────────────────────────────────────────────────────────────
	title := headerStyle.Render("Job Application Tracker")
	count := countStyle.Render(fmt.Sprintf("[%d applications]", len(t.rows)))
	gap := t.width - lipgloss.Width(title) - lipgloss.Width(count)
	if gap < 1 {
		gap = 1
	}
	b.WriteString(title)
	b.WriteString(strings.Repeat(" ", gap))
	b.WriteString(count)
	b.WriteString("\n\n")

	// ── Box rows ─────────────────────────────────────────────────────────────
	innerWidth := t.width - 2 // subtract the two │ border chars

	end := t.offset + t.rowsHeight
	if end > len(t.rows) {
		end = len(t.rows)
	}
	visibleRows := end - t.offset

	lines := make([]string, 0, t.rowsHeight)
	for i, row := range t.rows[t.offset:end] {
		isSelected := i+t.offset == t.cursor
		cursor := "  "
		if isSelected {
			cursor = "> "
		}

		var line strings.Builder
		line.WriteString(cursor)
		for j, cell := range row {
			style := normalStyle
			if isSelected {
				style = selectedStyle
			}
			line.WriteString(style.Width(t.cols[j].Width).Render(cell))
			if j < len(row)-1 {
				sep := "| "
				if isSelected {
					sep = selectedStyle.Render("| ")
				}
				line.WriteString(sep)
			}
		}
		lines = append(lines, line.String())
	}
	// fill remaining visible area with blank lines
	for range t.rowsHeight - visibleRows {
		lines = append(lines, "")
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(innerWidth)

	b.WriteString(boxStyle.Render(strings.Join(lines, "\n")))
	b.WriteString("\n\n")

	// ── Help bar ─────────────────────────────────────────────────────────────
	helpString := "j/k navigate   enter: details   n: new   q: quit"
	b.WriteString(helpStyle.Render(helpString))

	v := tea.NewView(b.String())
	v.AltScreen = true
	return v
}
