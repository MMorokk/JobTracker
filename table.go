package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Column struct {
	Title string
	Width int
}

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("240"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57"))

	normalStyle = lipgloss.NewStyle()
)

type Table struct {
	cols   []Column
	rows   [][]string
	cursor int
	height int
	width  int
	offset int
}

func NewTable(cols []Column, rows [][]string, width, height int) Table {
	return Table{
		cols:   cols,
		rows:   rows,
		width:  width,
		height: height,
	}
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
		if t.cursor >= t.offset+t.height {
			t.offset++
		}
	}
}

func (t *Table) SetSize(width, height int) {
	t.height = height
	t.width = width
}

func (t *Table) SelectedRow() []string {
	if len(t.rows) == 0 {
		return nil
	}
	return t.rows[t.cursor]
}

func (t *Table) View() string {
	var b strings.Builder

	// header
	for _, col := range t.cols {
		b.WriteString(headerStyle.Width(col.Width).Render(col.Title))
	}
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", t.width) + "\n")

	// visible rows (clamped to avoid panic when rows < height)
	end := t.offset + t.height
	if end > len(t.rows) {
		end = len(t.rows)
	}
	for i, row := range t.rows[t.offset:end] {
		style := normalStyle
		if i+t.offset == t.cursor {
			style = selectedStyle
		}
		for j, cell := range row {
			b.WriteString(style.Width(t.cols[j].Width).Render(cell))
		}
		b.WriteString("\n")
	}

	return b.String()
}
