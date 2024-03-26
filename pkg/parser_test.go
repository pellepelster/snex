package pkg

import (
	"github.com/alecthomas/assert/v2"
	"testing"
)

func TestParseDocumentEmptyDocument(t *testing.T) {
	document, err := ParseDocument(Document{File: "file1", Content: ""})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(document.Lines))
}

func TestParseDocumentFullSnippet(t *testing.T) {

	content := `some preface
snippet[id1]
some Content
/snippet
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
snippet[id1]
some Content
/snippet
more text`

	document, err := ParseDocument(Document{File: "file1", Content: content})
	assert.NoError(t, err)

	errors := ValidateDocuments([]ParsedDocument{document})
	assert.Equal(t, 0, len(errors))
}

func TestValidateDocumentsDuplicateSnippetStart(t *testing.T) {

	content := `some preface
snippet[id1]
something else
snippet[id1]
some Content
/snippet
more text`

	document, err := ParseDocument(Document{File: "file1", Content: content})
	assert.NoError(t, err)

	errors := ValidateDocuments([]ParsedDocument{document})
	assert.Equal(t, 1, len(errors))
	assert.Equal(t, "start marker for snippet 'id1' found more than once (file1:2, file1:4)", errors[0].Error())
}

func TestValidateDocumentsDuplicateSnippetEnd(t *testing.T) {

	content := `some preface
snippet[id1]
something else
/snippet
some Content
/snippet
more text`

	document, err := ParseDocument(Document{File: "file1", Content: content})
	assert.NoError(t, err)

	errors := ValidateDocuments([]ParsedDocument{document})
	assert.Equal(t, 1, len(errors))
	assert.Equal(t, "too many end markers found in 'file1'", errors[0].Error())
}

func TestValidateDocumentsNoSnippetEnd(t *testing.T) {

	content := `some preface
snippet[id1]
more text`

	document, err := ParseDocument(Document{File: "file1", Content: content})
	assert.NoError(t, err)

	errors := ValidateDocuments([]ParsedDocument{document})
	assert.Equal(t, 1, len(errors))
	assert.Equal(t, "not all start markers are closed in 'file1'", errors[0].Error())
}

func TestValidateDocumentsNoSnippetStart(t *testing.T) {

	content := `some preface
/snippet
more text`

	document, err := ParseDocument(Document{File: "file1", Content: content})
	assert.NoError(t, err)

	errors := ValidateDocuments([]ParsedDocument{document})
	assert.Equal(t, 1, len(errors))
	assert.Equal(t, "too many end markers found in 'file1'", errors[0].Error())
}

func TestCountMarkers(t *testing.T) {

	content := `
# Example 1

## Include snippet1

<!-- insertSnippet[snippet1] -->
<!-- /insertSnippet -->

## Include full file

<!-- insertFile[file1.go] -->
<!-- /insertFile -->
`
	document, err := ParseDocument(Document{File: "file1", Content: content})
	assert.NoError(t, err)
	assert.Equal(t, 2, countStartMarkers(document.Lines))
	assert.Equal(t, 2, countEndMarkers(document.Lines))
}

func TestValidateDocumentsInsertFileSelfReference(t *testing.T) {

	content := `lorem 
insertFile[file1]
content
/insertFile
ipsum`

	document, err := ParseDocument(Document{File: "file1", Content: content})
	assert.NoError(t, err)

	errors := ValidateDocuments([]ParsedDocument{document})
	assert.Equal(t, 1, len(errors))
	assert.Equal(t, "insert file snippet 'file1' references itself", errors[0].Error())
}

func TestValidateDocumentsSnippetMissing(t *testing.T) {

	content := `lorem 
insertSnippet[snippet1] 
content
/insertSnippet
ipsum`

	document, err := ParseDocument(Document{File: "file1", Content: content})
	assert.NoError(t, err)

	errors := ValidateDocuments([]ParsedDocument{document})
	assert.Equal(t, 1, len(errors))
	assert.Equal(t, "referenced snippet 'snippet1' not found", errors[0].Error())
}

func TestValidateDocumentsStartEndMultipleDocuments(t *testing.T) {

	content1 := `some preface
snippet[id1]
some Content
more text`

	content2 := `some preface
/snippet
more text`

	document1, err := ParseDocument(Document{File: "file1", Content: content1})
	assert.NoError(t, err)

	document2, err := ParseDocument(Document{File: "file2", Content: content2})
	assert.NoError(t, err)

	errors := ValidateDocuments([]ParsedDocument{document1, document2})
	assert.Equal(t, 2, len(errors))
	assert.Equal(t, "not all start markers are closed in 'file1'", errors[0].Error())
	assert.Equal(t, "too many end markers found in 'file2'", errors[1].Error())
}

func TestValidateDocumentsStartEndMultipleDocumentsInsert(t *testing.T) {

	content1 := `some preface
snippet[id1]
some Content
/snippet
more text`

	content2 := `some preface
insertSnippet[id1]
lorem ipsum
/insertSnippet
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
snippet[id1]
some new Content
/snippet
more text`

	target := `lorem
insertSnippet[id1]
some old Content
/insertSnippet
ipsum`

	targetReplaced := `lorem
insertSnippet[id1]
some new Content
/insertSnippet
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

func TestGetSnippetLines(t *testing.T) {

	source := `some preface
snippet[id1]
some new content
/snippet
more text`

	document, err := ParseDocument(Document{File: "source", Content: source})
	assert.NoError(t, err)

	lines := getSnippetLines([]ParsedDocument{document}, "id1")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(lines))
	assert.Equal(t, "some new content", lines[0])
}

func TestReplaceFiles(t *testing.T) {

	source := `yolo1
yolo2
yolo3`

	target := `lorem
insertFile[source]
some old Content
/insertFile
ipsum`

	targetReplaced := `lorem
insertFile[source]
yolo1
yolo2
yolo3
/insertFile
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
snippet[id1]
some new Content
/snippet
more text
`

	target := `lorem
insertSnippet[id1]
some old Content
/insertSnippet
ipsum
`

	targetReplaced := `lorem
insertSnippet[id1]
some new Content
/insertSnippet
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
snippet[id1]
some new Content
/snippet
more text`

	target := `
lorem
insertSnippet[id1]
some old Content
/insertSnippet
ipsum`

	targetReplaced := `
lorem
insertSnippet[id1]
some new Content
/insertSnippet
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

func TestReplaceSnippetsRemoveIndent(t *testing.T) {

	source := `some preface
snippet[id1]
	snippet line 1
		snippet line 2
	snippet line 3
/snippet
more text
`

	target := `lorem
insertSnippet[id1]
some old Content
/insertSnippet
ipsum
`

	targetReplaced := `lorem
insertSnippet[id1]
snippet line 1
	snippet line 2
snippet line 3
/insertSnippet
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
