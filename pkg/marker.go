package pkg

import "regexp"

var snippetStartExpression = regexp.MustCompile(`[^|\s]*snippet\[\s*([a-zA-Z0-9_\-]*)\s*\][\s|$]*`)
var snippetEndExpression = regexp.MustCompile(`[^|\s]*/snippet[\s|$]*`)

var insertSnippetStartExpression = regexp.MustCompile(`[^|\s]*insertSnippet\[\s*([a-zA-Z0-9_\-]*)\s*\][\s|$]*`)
var insertSnippetEndExpression = regexp.MustCompile(`[^|\s]*/insertSnippet[\s|$]*`)

var insertFileStartExpression = regexp.MustCompile(`[^|\s]*insertFile\[\s*([a-zA-Z0-9_\-\\.]*)\s*\][\s|$]*`)
var insertFileEndExpression = regexp.MustCompile(`[^|\s]*/insertFile[\s|$]*`)

func ParseMarker(line string) *SnippetMarker {
	snippetStart := snippetStartExpression.FindStringSubmatch(line)
	if len(snippetStart) == 2 {
		return &SnippetMarker{IsSnippet: true, IsStart: true, Id: snippetStart[1]}
	}

	if snippetEndExpression.MatchString(line) {
		return &SnippetMarker{IsSnippet: true, IsEnd: true}
	}

	if insertSnippetEndExpression.MatchString(line) {
		return &SnippetMarker{IsInsertSnippet: true, IsEnd: true}
	}

	insertSnippetStart := insertSnippetStartExpression.FindStringSubmatch(line)
	if len(insertSnippetStart) == 2 {
		return &SnippetMarker{IsInsertSnippet: true, IsStart: true, Id: insertSnippetStart[1]}
	}

	fileStart := insertFileStartExpression.FindStringSubmatch(line)
	if len(fileStart) == 2 {
		return &SnippetMarker{IsInsertFile: true, IsStart: true, Id: fileStart[1]}
	}

	if insertFileEndExpression.MatchString(line) {
		return &SnippetMarker{IsInsertFile: true, IsEnd: true}
	}

	return nil
}
