package render

import (
	"bytes"
	"os"
	"strings"
	"text/template"
)

// FileToString will accept a render file and return the render as a string.
func FileToString(tplFilename string) (string, error) {
	readBytes, err := os.ReadFile(tplFilename)
	if err != nil {
		return "", err
	}
	return string(readBytes), nil
}

func RenderCommands(data map[string]interface{}, tplString string) []string {
	tpl := template.Must(template.New("").Parse(tplString))
	var tplBuffer bytes.Buffer
	tpl.Execute(&tplBuffer, data)
	return strings.Split(tplBuffer.String(), "\n")
}
