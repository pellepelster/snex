package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReplaceSnippetsWithLongLeadIn(t *testing.T) {

	const snippets = `first line
<!-- snippet:snippet1 -->
xxx1
<!-- /snippet:snippet1 -->
<!-- snippet:snippet2 -->
xxx2
<!-- /snippet:snippet2 -->
<!-- snippet:snippet3 -->
xxx3
<!-- /snippet:snippet3 -->
<!-- snippet:snippet4 -->
xxx4
<!-- /snippet:snippet4 -->
snippet third line`

	const input = `first line
<!-- snippet:snippet1 -->
yyy
<!-- /snippet:snippet1 -->
second line
<!-- snippet:snippet2 -->
yyy
<!-- /snippet:snippet2 -->
third line
<!-- snippet:snippet3 -->
yyy
<!-- /snippet:snippet3 -->
fourth line
<!-- snippet:snippet4 -->
yyy
<!-- /snippet:snippet4 -->
fifth line`

	const expected = `first line
<!-- snippet:snippet1 -->
xxx1
<!-- /snippet:snippet1 -->
second line
<!-- snippet:snippet2 -->
xxx2
<!-- /snippet:snippet2 -->
third line
<!-- snippet:snippet3 -->
xxx3
<!-- /snippet:snippet3 -->
fourth line
<!-- snippet:snippet4 -->
xxx4
<!-- /snippet:snippet4 -->
fifth line`

	assert.Equal(t, expected, replaceSnippets(input, "", "{{.Content}}", []ParsedDocument{{snippets: parseSnippets(snippets)}}))
}

func TestReplaceSnippetsOnÖastLine(t *testing.T) {

	const snippets = `first line
<!-- snippet:snippet1 -->
xxx
<!-- /snippet:snippet1 -->`

	const input = `first line
<!-- snippet:snippet1 -->
zzz
<!-- /snippet:snippet1 -->
third line`

	const expected = `first line
<!-- snippet:snippet1 -->
xxx
<!-- /snippet:snippet1 -->
third line`

	assert.Equal(t, expected, replaceSnippets(input, "", "{{.Content}}", []ParsedDocument{{snippets: parseSnippets(snippets)}}))
}

func TestReplaceSnippetsWithSingleLineContent(t *testing.T) {

	const snippets = `snippet first line
<!-- snippet:snippet1 -->
xxx
<!-- /snippet:snippet1 -->
snippet third line`

	const input = `first line
<!-- snippet:snippet1 -->
zzz
<!-- /snippet:snippet1 -->
third line`

	const expected = `first line
<!-- snippet:snippet1 -->
xxx
<!-- /snippet:snippet1 -->
third line`

	assert.Equal(t, expected, replaceSnippets(input, "", "{{.Content}}", []ParsedDocument{ ParsedDocument{ snippets: parseSnippets(snippets)}}))
}

func TestReplaceSnippetsWithMultipleLineContent(t *testing.T) {

	const snippets = `snippet first line
<!-- snippet:snippet1 -->
xxx
<!-- /snippet:snippet1 -->
snippet third line`

	const input = `first line
<!-- snippet:snippet1 -->
zzz1
zzz2
zzz3
<!-- /snippet:snippet1 -->
third line`

	const expected = `first line
<!-- snippet:snippet1 -->
xxx
<!-- /snippet:snippet1 -->
third line`

	assert.Equal(t, expected, replaceSnippets(input, "", "{{.Content}}", []ParsedDocument{{snippets: parseSnippets(snippets)}}))
}

func TestReplaceSnippetsWithMixedLineContent(t *testing.T) {

	const snippets = `snippet first line
<!-- snippet:snippet1 -->
snippet 1 content
<!-- /snippet:snippet1 -->
line in between
<!-- snippet:snippet2 -->
snippet 2 content
<!-- /snippet:snippet2 -->
snippet third line`

	const input = `first line
<!-- snippet:snippet1 -->
zzz1
zzz2
zzz3
<!-- /snippet:snippet1 -->
line in between
<!-- snippet:snippet2 -->
<!-- /snippet:snippet2 -->
third line`

	const expected = `first line
<!-- snippet:snippet1 -->
snippet 1 content
<!-- /snippet:snippet1 -->
line in between
<!-- snippet:snippet2 -->
snippet 2 content
<!-- /snippet:snippet2 -->
third line`

	assert.Equal(t, expected, replaceSnippets(input, "", "{{.Content}}", []ParsedDocument{{snippets: parseSnippets(snippets)}}))
}

func TestReplaceSnippetsWithMixedLineContentNoSpaceInBetween(t *testing.T) {

	const snippets = `snippet first line
<!-- snippet:snippet1 -->
snippet 1 content
<!-- /snippet:snippet1 -->
<!-- snippet:snippet2 -->
snippet 2 content
<!-- /snippet:snippet2 -->
snippet third line`

	const input = `first line
<!-- snippet:snippet1 -->
zzz1
zzz2
zzz3
<!-- /snippet:snippet1 -->
<!-- snippet:snippet2 -->
<!-- /snippet:snippet2 -->
third line`

	const expected = `first line
<!-- snippet:snippet1 -->
snippet 1 content
<!-- /snippet:snippet1 -->
<!-- snippet:snippet2 -->
snippet 2 content
<!-- /snippet:snippet2 -->
third line`

	assert.Equal(t, expected, replaceSnippets(input, "", "{{.Content}}", []ParsedDocument{{snippets: parseSnippets(snippets)}}))
}

func TestReplaceSnippetsWithoutContent(t *testing.T) {

	const snippets = `snippet first line
<!-- snippet:snippet1 -->
xxx
<!-- /snippet:snippet1 -->
snippet third line`

	const input = `first line
<!-- snippet:snippet1 -->
<!-- /snippet:snippet1 -->
third line`

	const expected = `first line
<!-- snippet:snippet1 -->
xxx
<!-- /snippet:snippet1 -->
third line`

	assert.Equal(t, expected, replaceSnippets(input, "", "{{.Content}}", []ParsedDocument{{snippets: parseSnippets(snippets)}}))
}

/*
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
*/

func TestParseDocumentOneSnippet(t *testing.T) {

	const s = `first line
<!-- snippet:snippet1 -->
snippet1 content
<!-- /snippet:snippet1 -->
third line`

	snippets := parseSnippets(s)

	assert.Equal(t, 1, len(snippets))
	assert.Equal(t, -1, GetSnippetIndex(snippets, "xxx"))

	snippet := snippets[GetSnippetIndex(snippets, "snippet1")]
	assert.Equal(t, "snippet1", snippet.id)
	assert.Equal(t, []string{"snippet1 content"}, snippet.content)
	assert.Equal(t, 1, snippet.start)
	assert.Equal(t, 3, snippet.end)
}

func TestParseDocumentFileInclude(t *testing.T) {

	const s = `first line
<!-- file:folder/file1.txt -->
<!-- /file:folder/file1.txt -->
third line`

	snippets := parseSnippets(s)

	assert.Equal(t, 1, len(snippets))
	assert.Equal(t, "folder/file1.txt", snippets[0].id)
	assert.Equal(t, "folder/file1.txt", snippets[0].filename)
	assert.Equal(t, 1, snippets[0].start)
	assert.Equal(t, 2, snippets[0].end)
}

func TestParseDocumentTwoSnippets(t *testing.T) {

	const s = `first line
<!-- snippet:snippet1 -->
snippet1 content 1
snippet1 content 2
<!-- /snippet:snippet1 -->
third line
fourth line
<!-- snippet:snippet2 -->
snippet2 content
<!-- /snippet:snippet2 -->
last line`

	snippets := parseSnippets(s)

	assert.Equal(t, 2, len(snippets))
	assert.Equal(t, -1, GetSnippetIndex(snippets, "xxx"))

	snippet1 := snippets[GetSnippetIndex(snippets, "snippet1")]
	assert.Equal(t, "snippet1", snippet1.id)
	assert.Equal(t, []string{"snippet1 content 1", "snippet1 content 2"}, snippet1.content)
	assert.Equal(t, 1, snippet1.start)
	assert.Equal(t, 4, snippet1.end)

	snippet2 := snippets[GetSnippetIndex(snippets, "snippet2")]
	assert.Equal(t, "snippet2", snippet2.id)
	assert.Equal(t, []string{"snippet2 content"}, snippet2.content)
	assert.Equal(t, 7, snippet2.start)
	assert.Equal(t, 9, snippet2.end)
}

func TestParseDocumentSnippetWithoutEnd(t *testing.T) {

	const s = `first line
<!-- snippet:snippet1 -->
second line
third line`

	snippets := parseSnippets(s)

	assert.Equal(t, 1, len(snippets))
	assert.Equal(t, -1, GetSnippetIndex(snippets, "xxx"))

	snippet := snippets[GetSnippetIndex(snippets, "snippet1")]
	assert.Equal(t, "snippet1", snippet.id)
	assert.Equal(t, 1, snippet.start)
	assert.Equal(t, -1, snippet.end)
}

func TestParseDocumentSnippetWithoutStart(t *testing.T) {

	const s = `first line
second line
<!-- /snippet:snippet1 -->
third line`

	snippets := parseSnippets(s)

	assert.Equal(t, 1, len(snippets))
	assert.Equal(t, -1, GetSnippetIndex(snippets, "xxx"))

	snippet := snippets[GetSnippetIndex(snippets, "snippet1")]
	assert.Equal(t, "snippet1", snippet.id)
	assert.Equal(t, -1, snippet.start)
	assert.Equal(t, 4, snippet.end)
}
