package pkg

import (
	"bufio"
	"fmt"
	"regexp"
	"slices"
	"strings"
)

type Document struct {
	File    string
	Content string
}

type ParsedDocument struct {
	File  string
	Lines []DocumentLine
}

type DocumentLine struct {
	line    string
	number  int
	Snippet *SnippetMarker
}

type SnippetMarker struct {
	Id              string
	IsSnippet       bool
	IsInsertSnippet bool
	IsInsertFile    bool
	IsStart         bool
	IsEnd           bool
}

type SnippetMarkerPredicate func(marker *SnippetMarker) bool

var snippetStartExpression = regexp.MustCompile(`[^\s]*snippet\s*:\s*([a-zA-Z0-9_\-]*)\s*`)
var snippetEndExpression = regexp.MustCompile(`[^\s]*/snippet\s*:\s*([a-zA-Z0-9_\-]*)\s*`)

var insertSnippetStartExpression = regexp.MustCompile(`[^\s]*insertSnippet\s*:\s*([a-zA-Z0-9_\-]*)\s*`)
var insertSnippetEndExpression = regexp.MustCompile(`[^\s]*/insertSnippet\s*:\s*([a-zA-Z0-9_\-]*)\s*`)

var insertFileStartExpression = regexp.MustCompile(`[^\s]*insertFile\s*:\s*([^\s]*)\s*`)
var insertFileEndExpression = regexp.MustCompile(`[^\s]*/insertFile\s*:\s*([^\s]*)\s*`)

func parseMarker(line string) *SnippetMarker {
	snippetEnd := snippetEndExpression.FindStringSubmatch(line)
	if len(snippetEnd) == 2 {
		return &SnippetMarker{IsSnippet: true, IsEnd: true, Id: snippetEnd[1]}
	}

	snippetStart := snippetStartExpression.FindStringSubmatch(line)
	if len(snippetStart) == 2 {
		return &SnippetMarker{IsSnippet: true, IsStart: true, Id: snippetStart[1]}
	}

	insertSnippetEnd := insertSnippetEndExpression.FindStringSubmatch(line)
	if len(insertSnippetEnd) == 2 {
		return &SnippetMarker{IsInsertSnippet: true, IsEnd: true, Id: insertSnippetEnd[1]}
	}

	insertSnippetStart := insertSnippetStartExpression.FindStringSubmatch(line)
	if len(insertSnippetStart) == 2 {
		return &SnippetMarker{IsInsertSnippet: true, IsStart: true, Id: insertSnippetStart[1]}
	}

	fileEnd := insertFileEndExpression.FindStringSubmatch(line)
	if len(fileEnd) == 2 {
		return &SnippetMarker{IsInsertFile: true, IsEnd: true, Id: fileEnd[1]}
	}

	fileStart := insertFileStartExpression.FindStringSubmatch(line)
	if len(fileStart) == 2 {
		return &SnippetMarker{IsInsertFile: true, IsStart: true, Id: fileStart[1]}
	}

	return nil
}

func ParseDocument(document Document) (ParsedDocument, error) {
	var lines []DocumentLine
	scanner := bufio.NewScanner(strings.NewReader(document.Content))

	lineNumber := 0
	for scanner.Scan() {
		lines = append(lines, DocumentLine{line: scanner.Text(), number: lineNumber, Snippet: parseMarker(scanner.Text())})
		lineNumber++
	}
	lineNumber++

	if strings.HasSuffix(document.Content, "\n") {
		lines = append(lines, DocumentLine{line: "", number: lineNumber})
	}

	return ParsedDocument{Lines: lines, File: document.File}, nil
}

func ValidateDocuments(documents []ParsedDocument) []error {
	var errors []error

	errors = append(errors, validateNoInsertFileSelfReference(documents)...)
	errors = append(errors, validateSnippetsMultipleDocuments(documents)...)
	errors = append(errors, validateSnippetMarkerDuplicates(documents)...)
	errors = append(errors, validateMarkerStartEnd(documents)...)
	errors = append(errors, validateSnippetsMissing(documents)...)

	return errors
}

func getContentForSnippet(documents []ParsedDocument, id string) []string {
	var foundSnippet = false

	for _, document := range documents {
		var lines []string
		for _, line := range document.Lines {

			if line.Snippet != nil && !foundSnippet && line.Snippet.IsSnippet && line.Snippet.Id == id && line.Snippet.IsStart {
				foundSnippet = true
			}

			if line.Snippet != nil && foundSnippet && line.Snippet.IsSnippet && line.Snippet.Id == id && line.Snippet.IsEnd {
				return lines
			}

			if foundSnippet && line.Snippet == nil {
				lines = append(lines, line.line)
			}
		}
	}

	return []string{}
}

func getContentForFile(documents []ParsedDocument, file string) []string {
	for _, document := range documents {
		if strings.HasSuffix(document.File, file) {
			var lines []string
			for _, line := range document.Lines {
				lines = append(lines, line.line)
			}

			return lines
		}
	}

	return []string{}
}

func hasSnippet(documents []ParsedDocument, id string) bool {
	for _, document := range documents {
		for _, line := range document.Lines {
			if line.Snippet != nil && line.Snippet.IsSnippet && line.Snippet.Id == id && line.Snippet.IsStart {
				return true
			}
		}
	}

	return false
}

func ReplaceSnippets(documents []ParsedDocument, template string) ([]Document, error) {
	var replacedDocuments []Document

	for _, document := range documents {
		var lines []string
		isSnippet := false
		for _, line := range document.Lines {
			snippet := line.Snippet

			if isSnippet {
				if snippet != nil && (snippet.IsInsertSnippet || snippet.IsInsertFile) && snippet.IsEnd {
					lines = append(lines, line.line)
					isSnippet = false
					continue
				} else {
					continue
				}
			}

			if snippet != nil && !isSnippet && snippet.IsStart {

				if snippet.IsInsertSnippet {
					snippetLines := getContentForSnippet(documents, snippet.Id)
					renderedLines, err := executeTemplateWithDefault(snippetLines, document.File, template)
					if err != nil {
						return nil, err
					}

					isSnippet = true
					lines = append(lines, line.line)
					lines = append(lines, renderedLines...)
					continue
				}

				if snippet.IsInsertFile {
					snippetLines := getContentForFile(documents, snippet.Id)
					renderedLines, err := executeTemplateWithDefault(snippetLines, document.File, template)
					if err != nil {
						return nil, err
					}

					isSnippet = true
					lines = append(lines, line.line)
					lines = append(lines, renderedLines...)
					continue
				}
			}

			lines = append(lines, line.line)
		}

		replacedDocuments = append(replacedDocuments, Document{document.File, strings.Join(lines, "\n")})
	}

	return replacedDocuments, nil
}

func validateSnippetMarkerDuplicates(documents []ParsedDocument) []error {
	var errors []error

	errors = append(errors, validateDuplicates(documents, "start marker for snippet '%s' found more than once", func(marker *SnippetMarker) bool {
		return marker.IsSnippet && marker.IsStart
	})...)

	errors = append(errors, validateDuplicates(documents, "end marker for snippet '%s' found more than once", func(marker *SnippetMarker) bool {
		return marker.IsSnippet && marker.IsEnd
	})...)

	return errors
}

func validateMarkerStartEnd(documents []ParsedDocument) []error {

	var errors []error

	snippets := collectSnippets(documents, func(marker *SnippetMarker) bool {
		return marker.IsSnippet
	})

	for id, lines := range snippets {

		if countStartMarkers(lines) == 1 && countEndMarkers(lines) == 0 {
			errors = append(errors, fmt.Errorf("snippet '%s' has no end marker", id))
		}

		if countStartMarkers(lines) == 0 && countEndMarkers(lines) == 1 {
			errors = append(errors, fmt.Errorf("snippet '%s' has no start marker", id))
		}
	}

	return errors
}

func countStartMarkers(lines []DocumentLine) int {
	return countLines(lines, func(line DocumentLine) bool {
		return line.Snippet.IsStart
	})
}

func countEndMarkers(lines []DocumentLine) int {
	return countLines(lines, func(line DocumentLine) bool {
		return line.Snippet.IsEnd
	})
}

func countLines(lines []DocumentLine, predicate func(line DocumentLine) bool) int {
	count := 0
	for _, line := range lines {
		if predicate(line) {
			count++
		}
	}
	return count
}

func collectSnippets(documents []ParsedDocument, predicate SnippetMarkerPredicate) map[string][]DocumentLine {
	snippets := make(map[string][]DocumentLine)

	for _, document := range documents {
		for _, line := range document.Lines {

			if line.Snippet != nil {
				id := line.Snippet.Id
				if predicate(line.Snippet) {

					val, hasSnippet := snippets[id]
					if hasSnippet {
						snippets[id] = append(val, line)
					} else {
						snippets[id] = []DocumentLine{line}
					}
				}
			}
		}
	}

	return snippets
}

func validateNoInsertFileSelfReference(documents []ParsedDocument) []error {
	var errors []error

	for _, document := range documents {
		for _, line := range document.Lines {
			snippet := line.Snippet
			if snippet != nil && snippet.IsInsertFile && snippet.IsStart && snippet.Id == document.File {
				errors = append(errors, fmt.Errorf("insert file snippet '%s' references itself", document.File))
			}
		}
	}

	return errors
}

func validateSnippetsMultipleDocuments(documents []ParsedDocument) []error {
	snippets := make(map[string][]string)

	for _, document := range documents {
		for _, line := range document.Lines {
			if line.Snippet != nil && line.Snippet.IsSnippet {
				id := line.Snippet.Id
				val, hasSnippet := snippets[id]
				if hasSnippet {
					if !slices.Contains(val, document.File) {
						snippets[id] = append(val, document.File)
					}
				} else {
					snippets[id] = []string{document.File}
				}
			}
		}
	}

	var errors []error
	for id, documents := range snippets {
		if len(documents) > 1 {
			errors = append(errors, fmt.Errorf("snippet '%s' found in more than one document", id))
		}
	}

	return errors
}

func validateSnippetsMissing(documents []ParsedDocument) []error {
	var errors []error

	for _, document := range documents {
		for _, line := range document.Lines {
			snippet := line.Snippet
			if snippet != nil && snippet.IsInsertSnippet && snippet.IsStart {

				if !hasSnippet(documents, snippet.Id) {
					errors = append(errors, fmt.Errorf("referenced snippet '%s' not found", snippet.Id))
				}
			}
		}
	}

	return errors
}

func validateDuplicates(documents []ParsedDocument, message string, predicate SnippetMarkerPredicate) []error {

	snippets := collectSnippets(documents, predicate)

	for id, lines := range snippets {
		if len(lines) > 1 {
			return []error{fmt.Errorf(message, id)}
		}
	}

	return []error{}
}
