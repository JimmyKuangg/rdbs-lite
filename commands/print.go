package commands

import (
	"fmt"
	"rdbslite/data"
	"strings"
)

func Print(db *data.Database, cmd Command) (string, error) {
	var out strings.Builder

	for _, table := range db.Tables {
		out.WriteString(renderTable(table))
		out.WriteString("\n")
	}

	return strings.TrimSpace(out.String()), nil
}

func renderTable(table *data.Table) string {
	headers := buildHeaders(table)
	widths := buildWidths(headers, table.Rows)
	border := buildBorder(widths)

	var out strings.Builder
	out.WriteString("TABLE ")
	out.WriteString(table.Name)
	out.WriteString("\n")

	out.WriteString(border)
	out.WriteString("\n")
	writeRow(&out, headers, widths)
	out.WriteString(border)
	out.WriteString("\n")

	for _, r := range table.Rows {
		cells := make([]string, len(headers))
		for i := 0; i < len(headers); i++ {
			if i < len(r.Values) {
				cells[i] = fmt.Sprint(r.Values[i])
			} else {
				cells[i] = ""
			}
		}
		writeRow(&out, cells, widths)
	}

	out.WriteString(border)
	return out.String()
}

func buildHeaders(table *data.Table) []string {
	headers := make([]string, len(table.Schema))
	for i, col := range table.Schema {
		headers[i] = col.Name
	}
	return headers
}

func buildBorder(widths []int) string {
	var b strings.Builder
	b.WriteString("+")
	for _, w := range widths {
		b.WriteString(strings.Repeat("-", w+2))
		b.WriteString("+")
	}
	return b.String()
}

func buildWidths(headers []string, rows []data.Row) []int {
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}

	for _, r := range rows {
		for i := 0; i < len(headers) && i < len(r.Values); i++ {
			cell := fmt.Sprint(r.Values[i])
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}
	return widths
}

func writeRow(out *strings.Builder, cells []string, widths []int) {
	out.WriteString("|")
	for i, w := range widths {
		cell := ""
		if i < len(cells) {
			cell = cells[i]
		}
		out.WriteString(fmt.Sprintf(" %-*s |", w, cell))
	}
	out.WriteString("\n")
}
