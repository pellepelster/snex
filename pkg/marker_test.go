package pkg

import (
	"github.com/alecthomas/assert/v2"
	"testing"
)

func TestParseMarkerEmptyLine(t *testing.T) {
	assert.Zero(t, ParseMarker(""))
}

func TestParseMarkerSnippetStart(t *testing.T) {

	lines := []string{"snippet[id1]", " snippet[id1]", "snippet[id1] ", " snippet[id1] ", "snippet[ id1]", "snippet[id1 ]", "snippet[ id1 ]", "prefix snippet[id1]", "snippet[id1] postfix", "prefix snippet[id1] postfix"}

	for _, line := range lines {
		marker := ParseMarker(line)
		assert.NotZero(t, marker, line)
		assert.True(t, marker.IsStart)
		assert.True(t, marker.IsSnippet)
		assert.False(t, marker.IsInsertFile)
		assert.False(t, marker.IsEnd)
		assert.Equal(t, "id1", marker.Id)
	}
}

func TestParseMarkerSnippetEnd(t *testing.T) {

	lines := []string{"/snippet", " /snippet", "/snippet ", " /snippet ", "prefix /snippet", "/snippet postfix", "prefix /snippet postfix"}

	for _, line := range lines {
		marker := ParseMarker(line)
		assert.NotZero(t, marker, line)
		assert.False(t, marker.IsStart)
		assert.True(t, marker.IsSnippet)
		assert.False(t, marker.IsInsertSnippet, line)
		assert.False(t, marker.IsInsertFile)
		assert.True(t, marker.IsEnd)
		assert.Zero(t, marker.Id)
	}
}

func TestParseMarkerInsertSnippetStart(t *testing.T) {

	lines := []string{"insertSnippet[id1]", " insertSnippet[id1]", "insertSnippet[id1] ", " insertSnippet[id1] ", "insertSnippet[ id1]", "insertSnippet[id1 ]", "insertSnippet[ id1 ]", "prefix insertSnippet[id1]", "insertSnippet[id1] postfix", "prefix insertSnippet[id1] postfix"}

	for _, line := range lines {
		marker := ParseMarker(line)
		assert.NotZero(t, marker, line)
		assert.True(t, marker.IsStart, line)
		assert.False(t, marker.IsSnippet, line)
		assert.False(t, marker.IsInsertFile, line)
		assert.True(t, marker.IsInsertSnippet, line)
		assert.False(t, marker.IsEnd, line)
		assert.Equal(t, "id1", marker.Id)
	}
}

func TestParseMarkerInsertSnippetEnd(t *testing.T) {

	lines := []string{"/insertSnippet", " /insertSnippet", "/insertSnippet ", " /insertSnippet ", "prefix /insertSnippet", "/insertSnippet postfix", "prefix /insertSnippet postfix"}

	for _, line := range lines {
		marker := ParseMarker(line)
		assert.NotZero(t, marker, line)
		assert.False(t, marker.IsStart, line)
		assert.False(t, marker.IsSnippet, line)
		assert.True(t, marker.IsInsertSnippet, line)
		assert.False(t, marker.IsInsertFile, line)
		assert.True(t, marker.IsEnd, line)
		assert.Zero(t, marker.Id, line)
	}
}

func TestParseMarkerInsertFileStart(t *testing.T) {

	lines := []string{"insertFile[file1.txt]", " insertFile[file1.txt]", "insertFile[file1.txt] ", " insertFile[file1.txt] ", "insertFile[ file1.txt]", "insertFile[file1.txt ]", "insertFile[ file1.txt ]", "prefix insertFile[file1.txt]", "insertFile[file1.txt] postfix", "prefix insertFile[file1.txt] postfix"}

	for _, line := range lines {
		marker := ParseMarker(line)
		assert.NotZero(t, marker, line)
		assert.True(t, marker.IsStart, line)
		assert.False(t, marker.IsSnippet, line)
		assert.True(t, marker.IsInsertFile, line)
		assert.False(t, marker.IsInsertSnippet, line)
		assert.False(t, marker.IsEnd, line)
		assert.Equal(t, "file1.txt", marker.Id, line)
	}
}

func TestParseMarkerInsertFileEnd(t *testing.T) {

	lines := []string{"/insertFile", " /insertFile", "/insertFile ", " /insertFile ", "prefix /insertFile", "/insertFile postfix", "prefix /insertFile postfix"}

	for _, line := range lines {
		marker := ParseMarker(line)
		assert.NotZero(t, marker, line)
		assert.False(t, marker.IsStart, line)
		assert.False(t, marker.IsSnippet, line)
		assert.False(t, marker.IsInsertSnippet, line)
		assert.True(t, marker.IsInsertFile, line)
		assert.True(t, marker.IsEnd, line)
		assert.Zero(t, marker.Id, line)
	}
}
