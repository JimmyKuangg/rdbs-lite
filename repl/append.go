package repl

import (
	"os"

	"rdbslite/data"
)

func AppendAOF(command string) (err error) {
	aofPath := data.AOFPath()

	file, err := os.OpenFile(
		aofPath,
		os.O_APPEND|os.O_WRONLY,
		0o644,
	)
	if err != nil {
		return err
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	_, err = file.WriteString(command + "\n")
	return err
}
