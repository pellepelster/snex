package pkg

import (
	"github.com/alecthomas/assert/v2"
	"testing"
)

func TestParseMarkersEmptyLine(t *testing.T) {
	assert.Zero(t, parseMarker(""))
}

func TestParseMarkersSnippetStart(t *testing.T) {

	lines := []string{"snippet:id1", "snippet :id1", "snippet: id1", "snippet : id1", "prefix snippet: id1", "snippet: id1 postfix", "prefix snippet: id1 postfix"}

	for _, line := range lines {
		snippet := parseMarker(line)
		assert.NotZero(t, snippet, line)
		assert.True(t, snippet.IsStart)
		assert.True(t, snippet.IsSnippet)
		assert.False(t, snippet.IsInsertFile)
		assert.False(t, snippet.IsEnd)
		assert.Equal(t, "id1", snippet.Id)
	}
}

func TestParseDocumentEmptyDocument(t *testing.T) {
	document, err := ParseDocument(Document{File: "file1", Content: ""})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(document.Lines))
}

func TestParseDocumentFullSnippet(t *testing.T) {

	content := `some preface
snippet: id1
some Content
/snippet: id1
more text`

	document, err := ParseDocument(Document{File: "file1", Content: content})
	assert.NoError(t, err)
	assert.Equal(t, 5, len(document.Lines))
	assert.Zero(t, document.Lines[0].Snippet)
	assert.NotZero(t, document.Lines[1].Snippet)
	assert.Zero(t, document.Lines[2].Snippet)
	assert.NotZero(t, document.Lines[3].Snippet)
	assert.Zero(t, document.Lines[4].Snippet)
}

func TestParseDocumentTrailingNewline(t *testing.T) {

	content := `line1
line2
`

	document, err := ParseDocument(Document{File: "file1", Content: content})
	assert.NoError(t, err)
	assert.Equal(t, 3, len(document.Lines))
}

func TestValidateDocumentsFullSnippet(t *testing.T) {

	content := `some preface
snippet: id1
some Content
/snippet: id1
more text`

	document, err := ParseDocument(Document{File: "file1", Content: content})
	assert.NoError(t, err)

	errors := ValidateDocuments([]ParsedDocument{document})
	assert.Equal(t, 0, len(errors))
}

func TestValidateDocumentsDuplicateSnippetStart(t *testing.T) {

	content := `some preface
snippet: id1
something else
snippet: id1
some Content
/snippet: id1
more text`

	document, err := ParseDocument(Document{File: "file1", Content: content})
	assert.NoError(t, err)

	errors := ValidateDocuments([]ParsedDocument{document})
	assert.Equal(t, 1, len(errors))
	assert.Equal(t, "start marker for snippet 'id1' found more than once", errors[0].Error())
}

func TestValidateDocumentsDuplicateSnippetEnd(t *testing.T) {

	content := `some preface
snippet: id1
something else
/snippet: id1
some Content
/snippet: id1
more text`

	document, err := ParseDocument(Document{File: "file1", Content: content})
	assert.NoError(t, err)

	errors := ValidateDocuments([]ParsedDocument{document})
	assert.Equal(t, 1, len(errors))
	assert.Equal(t, "end marker for snippet 'id1' found more than once", errors[0].Error())
}

func TestValidateDocumentsNoSnippetEnd(t *testing.T) {

	content := `some preface
snippet: id1
more text`

	document, err := ParseDocument(Document{File: "file1", Content: content})
	assert.NoError(t, err)

	errors := ValidateDocuments([]ParsedDocument{document})
	assert.Equal(t, 1, len(errors))
	assert.Equal(t, "snippet 'id1' has no end marker", errors[0].Error())
}

func TestValidateDocumentsNoSnippetStart(t *testing.T) {

	content := `some preface
/snippet: id1
more text`

	document, err := ParseDocument(Document{File: "file1", Content: content})
	assert.NoError(t, err)

	errors := ValidateDocuments([]ParsedDocument{document})
	assert.Equal(t, 1, len(errors))
	assert.Equal(t, "snippet 'id1' has no start marker", errors[0].Error())
}

func TestValidateDocumentsInsertFileSelfReference(t *testing.T) {

	content := `lorem 
insertFile: file1
content
/insertFile: file1
ipsum`

	document, err := ParseDocument(Document{File: "file1", Content: content})
	assert.NoError(t, err)

	errors := ValidateDocuments([]ParsedDocument{document})
	assert.Equal(t, 1, len(errors))
	assert.Equal(t, "insert file snippet 'file1' references itself", errors[0].Error())
}

func TestValidateDocumentsSnippetMissing(t *testing.T) {

	content := `lorem 
insertSnippet: snippet1
content
/insertSnippet: snippet1
ipsum`

	document, err := ParseDocument(Document{File: "file1", Content: content})
	assert.NoError(t, err)

	errors := ValidateDocuments([]ParsedDocument{document})
	assert.Equal(t, 1, len(errors))
	assert.Equal(t, "referenced snippet 'snippet1' not found", errors[0].Error())
}

func TestValidateDocumentsStartEndMultipleDocuments(t *testing.T) {

	content1 := `some preface
snippet: id1
some Content
more text`

	content2 := `some preface
/snippet: id1
more text`

	document1, err := ParseDocument(Document{File: "file1", Content: content1})
	assert.NoError(t, err)

	document2, err := ParseDocument(Document{File: "file2", Content: content2})
	assert.NoError(t, err)

	errors := ValidateDocuments([]ParsedDocument{document1, document2})
	assert.Equal(t, 1, len(errors))
	assert.Equal(t, "snippet 'id1' found in more than one document (file1, file2)", errors[0].Error())
}

func TestValidateDocumentsStartEndMultipleDocumentsInsert(t *testing.T) {

	content1 := `some preface
snippet: id1
some Content
/snippet: id1
more text`

	content2 := `some preface
insertSnippet: id1
lorem ipsum
/insertSnippet: id1
more text`

	document1, err := ParseDocument(Document{File: "file1", Content: content1})
	assert.NoError(t, err)

	document2, err := ParseDocument(Document{File: "file2", Content: content2})
	assert.NoError(t, err)

	errors := ValidateDocuments([]ParsedDocument{document1, document2})
	assert.Equal(t, 0, len(errors))
}

func TestReplaceSnippets(t *testing.T) {

	source := `some preface
snippet: id1
some new Content
/snippet: id1
more text`

	target := `lorem
insertSnippet: id1
some old Content
/insertSnippet: id1
ipsum`

	targetReplaced := `lorem
insertSnippet: id1
some new Content
/insertSnippet: id1
ipsum`

	document1, err := ParseDocument(Document{File: "source", Content: source})
	assert.NoError(t, err)

	document2, err := ParseDocument(Document{File: "target", Content: target})
	assert.NoError(t, err)

	documents, err := ReplaceSnippets([]ParsedDocument{document1, document2}, "")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(documents))
	assert.Equal(t, "source", documents[0].File)
	assert.Equal(t, source, documents[0].Content)
	assert.Equal(t, "target", documents[1].File)
	assert.Equal(t, targetReplaced, documents[1].Content)
}

func TestReplaceFiles(t *testing.T) {

	source := `yolo1
yolo2
yolo3`

	target := `lorem
insertFile: source
some old Content
/insertFile: source
ipsum`

	targetReplaced := `lorem
insertFile: source
yolo1
yolo2
yolo3
/insertFile: source
ipsum`

	document1, err := ParseDocument(Document{File: "source", Content: source})
	assert.NoError(t, err)

	document2, err := ParseDocument(Document{File: "target", Content: target})
	assert.NoError(t, err)

	documents, err := ReplaceSnippets([]ParsedDocument{document1, document2}, "")
	assert.NoError(t, err)

	assert.Equal(t, 2, len(documents))
	assert.Equal(t, "source", documents[0].File)
	assert.Equal(t, source, documents[0].Content)
	assert.Equal(t, "target", documents[1].File)
	assert.Equal(t, targetReplaced, documents[1].Content)
}

func TestReplaceSnippetsTrailingNewline(t *testing.T) {

	source := `some preface
snippet: id1
some new Content
/snippet: id1
more text
`

	target := `lorem
insertSnippet: id1
some old Content
/insertSnippet: id1
ipsum
`

	targetReplaced := `lorem
insertSnippet: id1
some new Content
/insertSnippet: id1
ipsum
`

	document1, err := ParseDocument(Document{File: "source", Content: source})
	assert.NoError(t, err)

	document2, err := ParseDocument(Document{File: "target", Content: target})
	assert.NoError(t, err)

	documents, err := ReplaceSnippets([]ParsedDocument{document1, document2}, "")
	assert.NoError(t, err)

	assert.Equal(t, 2, len(documents))
	assert.Equal(t, "source", documents[0].File)
	assert.Equal(t, source, documents[0].Content)
	assert.Equal(t, "target", documents[1].File)
	assert.Equal(t, targetReplaced, documents[1].Content)
}

func TestReplaceSnippetsLeadingNewline(t *testing.T) {

	source := `
some preface
snippet: id1
some new Content
/snippet: id1
more text`

	target := `
lorem
insertSnippet: id1
some old Content
/insertSnippet: id1
ipsum`

	targetReplaced := `
lorem
insertSnippet: id1
some new Content
/insertSnippet: id1
ipsum`

	document1, err := ParseDocument(Document{File: "source", Content: source})
	assert.NoError(t, err)

	document2, err := ParseDocument(Document{File: "target", Content: target})
	assert.NoError(t, err)

	documents, err := ReplaceSnippets([]ParsedDocument{document1, document2}, "")
	assert.NoError(t, err)

	assert.Equal(t, 2, len(documents))
	assert.Equal(t, "source", documents[0].File)
	assert.Equal(t, source, documents[0].Content)
	assert.Equal(t, "target", documents[1].File)
	assert.Equal(t, targetReplaced, documents[1].Content)
}
