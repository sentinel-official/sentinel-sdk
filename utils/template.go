package utils

import (
	"bytes"
	"os"
	"strings"
	"text/template"
)

var funcMap = template.FuncMap{
	"join": strings.Join,
}

// ExecTemplateToFile generates content from a template and writes it to a file.
func ExecTemplateToFile(text string, data interface{}, fileName string) error {
	// Parse the template with custom functions
	tmpl, err := template.New("config").Funcs(funcMap).Parse(text)
	if err != nil {
		return err
	}

	// Execute the template and capture the output
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	// Write the generated content to the specified file
	return os.WriteFile(fileName, buf.Bytes(), 0644)
}
