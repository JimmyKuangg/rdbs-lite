package commands

import (
	"fmt"
	"strings"

	"rdbslite/data"
)

func Print(db *data.Database, cmd Command) (string, error) {
	var out strings.Builder

	for _, table := range db.Tables {
		out.WriteString(renderTable(table, []string{"*"}))
		out.WriteString("\n")
	}

	return strings.TrimSpace(out.String()), nil
}

func renderTable(table *data.Table, cols []string) string {
	idxs := selectedIndexes(table, cols)
	headers := buildHeaders(table, idxs)
	widths := buildWidths(table, idxs, headers)
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
		cells := make([]string, len(idxs))
		for i, idx := range idxs {
			if idx < len(r.Values) {
				cells[i] = fmt.Sprint(r.Values[idx])
			} else {
				cells[i] = ""
			}
		}
		writeRow(&out, cells, widths)
	}

	out.WriteString(border)
	return out.String()
}

func selectedIndexes(table *data.Table, cols []string) []int {
	// SELECT * (or Print path) -> all schema columns in order
	if len(cols) == 1 && cols[0] == "*" {
		idxs := make([]int, len(table.Schema))
		for i := range table.Schema {
			idxs[i] = i
		}
		return idxs
	}

	// Specific projection order, e.g. name,id
	idxs := make([]int, len(cols))
	for i, c := range cols {
		idxs[i] = table.ColumnIndex[strings.ToLower(c)]
	}
	return idxs
}

func buildHeaders(table *data.Table, idxs []int) []string {
	headers := make([]string, len(idxs))
	for i, idx := range idxs {
		headers[i] = table.Schema[idx].Name
	}
	return headers
}

func buildWidths(table *data.Table, idxs []int, headers []string) []int {
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}

	for _, r := range table.Rows {
		for i, idx := range idxs {
			if idx >= len(r.Values) {
				continue
			}
			cell := fmt.Sprint(r.Values[idx])
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}
	return widths
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

func writeRow(out *strings.Builder, cells []string, widths []int) {
	out.WriteString("|")

	for i, w := range widths {
		cell := ""
		if i < len(cells) {
			cell = cells[i]
		}

		fmt.Fprintf(out, " %-*s |", w, cell)
	}

	out.WriteString("\n")
}
