package utils

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
)

var funcMap = template.FuncMap{
	"join": strings.Join,
	"sum":  func(x, y int) int { return x + y },
}

// ExecTemplateToFile generates content from a template and writes it to a file.
func ExecTemplateToFile(text string, data interface{}, fileName string) error {
	// Parse the template with custom functions
	tmpl, err := template.New("config").Funcs(funcMap).Parse(text)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute the template and capture the output
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Write the generated content to the specified file
	if err := os.WriteFile(fileName, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
