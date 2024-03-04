package pkg

import (
	"bytes"
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

func executeTemplateWithDefault(snippet []string, file string, template string) ([]string, error) {

	if len(template) > 0 {
		return executeTemplate(template, snippet, file)
	}

	for _, template := range DefaultSnippetTemplates {
		for _, extension := range template.Extensions {
			if strings.HasSuffix(strings.ToLower(file), extension) {
				return executeTemplate(template.Template, snippet, file)
			}
		}
	}

	return snippet, nil
}
