package pkg

import (
	"bytes"
	"sort"
	"strings"
	template2 "text/template"
)

type SnippetTemplateData struct {
	Content  string
	Filename string
}

var TemplateHelp = "\t\t{{.Content}}\t\t snippet content\n" +
	"\t\t{{.Filename}}\t\t the file the snippet content originated from\n"

type SnippetTemplate struct {
	Template   string
	Extensions []string
}

var DefaultSnippetTemplates = []SnippetTemplate{
	{Template: "```\n{{.Content}}\n```\n", Extensions: []string{"md"}},
}

func executeTemplate(template string, snippet []string, file string) ([]string, error) {
	template = strings.ReplaceAll(template, "\\n", "\n")
	tmpl, err := template2.New("snippet").Parse(template)
	if err != nil {
		return nil, err
	}

	templateData := SnippetTemplateData{Content: strings.Join(snippet, "\n"), Filename: file}

	renderedTemplate := new(bytes.Buffer)
	err = tmpl.Execute(renderedTemplate, templateData)
	if err != nil {
		return nil, err
	}

	return strings.Split(renderedTemplate.String(), "\n"), nil
}

func ValidateTemplate(template string) error {
	tmpl, err := template2.New("snippet").Parse(template)
	if err != nil {
		return err
	}

	templateData := SnippetTemplateData{Content: "lorem ipsum"}

	renderedTemplate := new(bytes.Buffer)
	err = tmpl.Execute(renderedTemplate, templateData)
	if err != nil {
		return err
	}

	return nil
}

func executeTemplateWithDefault(lines []string, file string, template string) ([]string, error) {
	if len(template) > 0 {
		return executeTemplate(template, lines, file)
	}

	for _, template := range DefaultSnippetTemplates {
		for _, extension := range template.Extensions {
			if strings.HasSuffix(strings.ToLower(file), extension) {
				return executeTemplate(template.Template, lines, file)
			}
		}
	}

	return lines, nil
}

func longestCommonPrefix(originalLines []string) string {
	var longestPrefix = ""

	lines := make([]string, len(originalLines))
	copy(lines, originalLines)

	if len(lines) > 0 {
		sort.Strings(lines)

		firstLine := lines[0]
		lastLine := lines[len(lines)-1]

		for i := 0; i < len(firstLine); i++ {

			if (string(lastLine[i]) == " " || string(lastLine[i]) == "\t") && string(lastLine[i]) == string(firstLine[i]) {
				longestPrefix += string(lastLine[i])
			} else {
				return longestPrefix
			}
		}
	}

	return longestPrefix
}

func removeIndentation(lines []string) []string {
	prefix := longestCommonPrefix(lines)

	for index, line := range lines {
		lines[index] = strings.TrimPrefix(line, prefix)
	}

	return lines
}
