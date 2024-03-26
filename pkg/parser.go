package pkg

import (
	"bufio"
	"fmt"
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

func ParseDocument(document Document) (ParsedDocument, error) {
	var lines []DocumentLine
	scanner := bufio.NewScanner(strings.NewReader(document.Content))

	lineNumber := 0
	for scanner.Scan() {
		lines = append(lines, DocumentLine{line: scanner.Text(), number: lineNumber, Snippet: ParseMarker(scanner.Text())})
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

	errors = append(errors, validateSnippetMarkerDuplicates(documents)...)
	if len(errors) > 0 {
		return errors
	}

	errors = append(errors, validateNoInsertFileSelfReference(documents)...)
	errors = append(errors, validateMarkerStartEnd(documents)...)
	errors = append(errors, validateSnippetsMissing(documents)...)

	return errors
}

func getSnippetLines(documents []ParsedDocument, id string) []string {
	var foundSnippet = false

	for _, document := range documents {
		var lines []string
		for _, line := range document.Lines {

			if line.Snippet != nil && !foundSnippet && line.Snippet.IsSnippet && line.Snippet.Id == id && line.Snippet.IsStart {
				foundSnippet = true
			}

			if line.Snippet != nil && foundSnippet && line.Snippet.IsSnippet && line.Snippet.IsEnd {
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
					snippetLines := getSnippetLines(documents, snippet.Id)
					snippetLines = removeIndentation(snippetLines)

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

	errors = append(errors, validateDuplicates(documents, "start marker for snippet '%s' found more than once (%s)", func(marker *SnippetMarker) bool {
		return marker.IsSnippet && marker.IsStart
	})...)

	return errors
}

func validateMarkerStartEnd(documents []ParsedDocument) []error {
	var errors []error

	for _, document := range documents {

		if countStartMarkers(document.Lines) > countEndMarkers(document.Lines) {
			errors = append(errors, fmt.Errorf("not all start markers are closed in '%s'", document.File))
		}

		if countStartMarkers(document.Lines) < countEndMarkers(document.Lines) {
			errors = append(errors, fmt.Errorf("too many end markers found in '%s'", document.File))
		}

	}

	return errors
}

func countStartMarkers(lines []DocumentLine) int {
	return countLines(lines, func(line DocumentLine) bool {
		return line.Snippet != nil && line.Snippet.IsStart
	})
}

func countEndMarkers(lines []DocumentLine) int {
	return countLines(lines, func(line DocumentLine) bool {
		return line.Snippet != nil && line.Snippet.IsEnd
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

type DocumentSnippet struct {
	line DocumentLine
	file string
}

func collectSnippets(documents []ParsedDocument, predicate SnippetMarkerPredicate) map[string][]DocumentSnippet {
	snippets := make(map[string][]DocumentSnippet)

	for _, document := range documents {
		for _, line := range document.Lines {
			if line.Snippet != nil {
				id := line.Snippet.Id
				if predicate(line.Snippet) {

					val, hasSnippet := snippets[id]
					if hasSnippet {
						snippets[id] = append(val, DocumentSnippet{line, document.File})
					} else {
						snippets[id] = []DocumentSnippet{{line, document.File}}
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

	collectedSnippets := collectSnippets(documents, predicate)

	for id, documentSnippets := range collectedSnippets {
		if len(documentSnippets) > 1 {

			var files []string
			for _, documentSnippet := range documentSnippets {
				files = append(files, fmt.Sprintf("%s:%d", documentSnippet.file, documentSnippet.line.number+1))
			}
			return []error{fmt.Errorf(message, id, strings.Join(files, ", "))}
		}
	}

	return []error{}
}

func CountSnippets(document ParsedDocument) int {
	count := 0
	for _, line := range document.Lines {
		if line.Snippet != nil && line.Snippet.IsSnippet && line.Snippet.IsStart {
			count++
		}
	}

	return count
}
