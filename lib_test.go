package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseDocumentsOneSnippet(t *testing.T) {

	const s1 = `first line
<!-- snippet:snippet1 -->
snippet 1 content
<!-- /snippet:snippet1 -->
third line`

	const s2 = `first line
<!-- snippet:snippet2 -->
snippet 2 content
<!-- /snippet:snippet2 -->
third line`

	result := parseDocuments([]string{s1, s2})

	assert.Equal(t, 2, len(result.snippets))
	assert.Equal(t, -1, result.GetSnippetIndex("xxx"))

	snippet1 := result.snippets[result.GetSnippetIndex("snippet1")]
	assert.Equal(t, "snippet1", snippet1.id)
	assert.Equal(t, 1, snippet1.start)
	assert.Equal(t, 3, snippet1.end)

	snippet2 := result.snippets[result.GetSnippetIndex("snippet2")]
	assert.Equal(t, "snippet2", snippet2.id)
	assert.Equal(t, 1, snippet2.start)
	assert.Equal(t, 3, snippet2.end)
}

func TestParseDocumentOneSnippet(t *testing.T) {

	const s = `first line
<!-- snippet:snippet1 -->
second line
<!-- /snippet:snippet1 -->
third line`

	result := parseDocument(s)

	assert.Equal(t, 1, len(result.snippets))
	assert.Equal(t, -1, result.GetSnippetIndex("xxx"))

	snippet := result.snippets[result.GetSnippetIndex("snippet1")]
	assert.Equal(t, "snippet1", snippet.id)
	assert.Equal(t, 1, snippet.start)
	assert.Equal(t, 3, snippet.end)
}

func TestParseDocumentTwoSnippets(t *testing.T) {

	const s = `first line
<!-- snippet:snippet1 -->
snippet1 content
<!-- /snippet:snippet1 -->
third line
fourth line
<!-- snippet:snippet2 -->
snippet2 content
<!-- /snippet:snippet2 -->
last line`

	result := parseDocument(s)

	assert.Equal(t, 2, len(result.snippets))
	assert.Equal(t, -1, result.GetSnippetIndex("xxx"))

	snippet1 := result.snippets[result.GetSnippetIndex("snippet1")]
	assert.Equal(t, "snippet1", snippet1.id)
	assert.Equal(t, 1, snippet1.start)
	assert.Equal(t, 3, snippet1.end)

	snippet2 := result.snippets[result.GetSnippetIndex("snippet2")]
	assert.Equal(t, "snippet2", snippet2.id)
	assert.Equal(t, 5, snippet2.start)
	assert.Equal(t, 7, snippet2.end)
}

func TestParseDocumentSnippetWithoutEnd(t *testing.T) {

	const s = `first line
<!-- snippet:snippet1 -->
second line
third line`

	result := parseDocument(s)

	assert.Equal(t, 1, len(result.snippets))
	assert.Equal(t, -1, result.GetSnippetIndex("xxx"))

	snippet := result.snippets[result.GetSnippetIndex("snippet1")]
	assert.Equal(t, "snippet1", snippet.id)
	assert.Equal(t, 1, snippet.start)
	assert.Equal(t, -1, snippet.end)
}

func TestParseDocumentSnippetWithoutStart(t *testing.T) {

	const s = `first line
second line
<!-- /snippet:snippet1 -->
third line`

	result := parseDocument(s)

	assert.Equal(t, 1, len(result.snippets))
	assert.Equal(t, -1, result.GetSnippetIndex("xxx"))

	snippet := result.snippets[result.GetSnippetIndex("snippet1")]
	assert.Equal(t, "snippet1", snippet.id)
	assert.Equal(t, -1, snippet.start)
	assert.Equal(t, 2, snippet.end)
}
