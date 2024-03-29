package pkg

import (
	"github.com/alecthomas/assert/v2"
	"testing"
)

func TestExecuteTemplate(t *testing.T) {
	snippets, err := executeTemplate("begin\n{{.Content}}\nend", []string{"line1", "line2"}, "file1")
	assert.NoError(t, err)
	assert.Equal(t, []string{"begin", "line1", "line2", "end"}, snippets)
}

func TestExecuteTemplateTrailingNewline(t *testing.T) {
	snippets, err := executeTemplate("begin\n{{.Content}}\nend\n", []string{"line1", "line2"}, "file1")
	assert.NoError(t, err)
	assert.Equal(t, []string{"begin", "line1", "line2", "end", ""}, snippets)
}

func TestExecuteTemplateMarkdown(t *testing.T) {
	snippets, err := executeTemplateWithDefault([]string{"line1", "line2"}, "test.md", "")
	assert.NoError(t, err)
	assert.Equal(t, []string{"```", "line1", "line2", "```", ""}, snippets)
}

func TestExecuteTemplateMarkdownUppercase(t *testing.T) {
	snippets, err := executeTemplateWithDefault([]string{"line1", "line2"}, "test.MD", "")
	assert.NoError(t, err)
	assert.Equal(t, []string{"```", "line1", "line2", "```", ""}, snippets)
}

func TestExecuteTemplateUnknownExtension(t *testing.T) {
	snippets, err := executeTemplateWithDefault([]string{"line1", "line2"}, "test.yolo", "")
	assert.NoError(t, err)
	assert.Equal(t, []string{"line1", "line2"}, snippets)
}

func TestValidateTemplate(t *testing.T) {
	err := ValidateTemplate("")
	assert.NoError(t, err)

	err = ValidateTemplate("lorem ipsum")
	assert.NoError(t, err)

	err = ValidateTemplate("{{.NotExistingVariable}}")
	assert.Error(t, err)
}

func TestRemoveIndentationWhitespaceBlank(t *testing.T) {
	assert.Equal(t, []string{"one space", " two spaces", "one space", "  three spaces"}, removeIndentation([]string{" one space", "  two spaces", " one space", "   three spaces"}))
}

func TestRemoveIndentationWhitespaceTab(t *testing.T) {
	assert.Equal(t, []string{"one tab", "\ttwo tabs", "one tab"}, removeIndentation([]string{"\tone tab", "\t\ttwo tabs", "\tone tab"}))
}

func TestRemoveIndentationNonWhitespace(t *testing.T) {
	assert.Equal(t, []string{"aabb", "aabb"}, removeIndentation([]string{"aabb", "aabb"}))
}

func TestRemoveIndentationSingleLine(t *testing.T) {
	assert.Equal(t, []string{"one space"}, removeIndentation([]string{" one space"}))
}
