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
	nameWidth := len("Name")
	typeWidth := len("Type")

	for _, col := range table.Schema {
		if len(col.Name) > nameWidth {
			nameWidth = len(col.Name)
		}
		typeText := string(col.Type)
		if len(typeText) > typeWidth {
			typeWidth = len(typeText)
		}
	}

	border := "+" + strings.Repeat("-", nameWidth+2) + "+" + strings.Repeat("-", typeWidth+2) + "+"

	var out strings.Builder
	out.WriteString("TABLE ")
	out.WriteString(table.Name)
	out.WriteString("\n")
	out.WriteString(border)
	out.WriteString("\n")
	out.WriteString(fmt.Sprintf("| %-*s | %-*s |\n", nameWidth, "Name", typeWidth, "Type"))
	out.WriteString(border)
	out.WriteString("\n")

	for _, col := range table.Schema {
		out.WriteString(fmt.Sprintf("| %-*s | %-*s |\n", nameWidth, col.Name, typeWidth, col.Type))
	}

	out.WriteString(border)
	return out.String()
}
