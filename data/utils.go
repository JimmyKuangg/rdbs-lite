package data

import (
	"fmt"
	"strings"
)

type ColumnType string

const (
	IntType  ColumnType = "INT"
	TextType ColumnType = "TEXT"
	BoolType ColumnType = "BOOL"
)

var reservedIdentifiers = map[string]struct{}{
	"CREATE": {},
	"TABLE":  {},
}

func normalizeIdent(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}

func ParseColumnType(input string) (ColumnType, error) {
	inputType := normalizeIdent(input)

	switch inputType {
	case "INT":
		return IntType, nil
	case "TEXT":
		return TextType, nil
	case "BOOL":
		return BoolType, nil
	default:
		return "", fmt.Errorf("invalid column type: %q", input)
	}
}

func IsReservedIdentifier(name string) bool {
	n := normalizeIdent(name)

	// Check reserved keywords
	if _, exists := reservedIdentifiers[n]; exists {
		return true
	}

	// Type names are also reserved for identifiers
	_, err := ParseColumnType(n)
	return err == nil
}
