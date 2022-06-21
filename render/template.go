package render

import (
	"bytes"
	"os"
	"strings"
	"text/template"
)

// FileToString will accept a template filename and return the template
// as a string. This is just a convenience function.
func FileToString(tplFilename string) (string, error) {
	readBytes, err := os.ReadFile(tplFilename)
	if err != nil {
		return "", err
	}
	return string(readBytes), nil
}

func Commands(data map[string]interface{}, tplString string) []string {
	tpl := template.Must(template.New("").Parse(tplString))
	var tplBuffer bytes.Buffer
	tpl.Execute(&tplBuffer, data)
	return strings.Split(tplBuffer.String(), "\n")
}
