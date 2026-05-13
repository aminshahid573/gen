package ui

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
)

// RenderTable renders any number of headers and rows.
// Each row must have the same number of columns as headers.
func RenderTable(headers []string, rows [][]string) string {
	// find the widest value per column for natural sizing
	colWidths := make([]int, len(headers))
	for i, h := range headers {
		if len(h) > colWidths[i] {
			colWidths[i] = len(h)
		}
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(Purple)).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return HeaderStyle
			}
			var base lipgloss.Style
			if row%2 == 0 {
				base = EvenRowStyle
			} else {
				base = OddRowStyle
			}
			if col < len(colWidths) {
				return base.Width(colWidths[col] + 2) // +2 padding
			}
			return base
		}).
		Headers(headers...)

	for _, row := range rows {
		t = t.Row(row...)
	}

	return t.String()
}
