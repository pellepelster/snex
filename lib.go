package main

import (
	"bufio"
	"regexp"
	"strings"
)

type Snippet struct {
	id    string
	start int
	end   int
}

type Snippets struct {
	snippets []Snippet
}

func (snippets Snippets) GetSnippetIndex(id string) int {
	for index, snippet := range snippets.snippets {
		if snippet.id == id {
			return index
		}
	}

	return -1
}

func parseDocuments(contents []string) Snippets {
	result := Snippets{}

	for _, content := range contents {
		result.snippets = append(result.snippets, parseDocument(content).snippets...)
	}

	return result
}

func parseDocument(content string) Snippets {

	snippetStart := regexp.MustCompile(`snippet:([a-zA-Z0-9]*)`)
	snippetEnd := regexp.MustCompile(`/snippet:([a-zA-Z0-9]*)`)

	snippets := Snippets{}

	scanner := bufio.NewScanner(strings.NewReader(content))
	index := 0
	for scanner.Scan() {
		end := snippetEnd.FindStringSubmatch(scanner.Text())
		if len(end) == 2 {
			snippetIndex := snippets.GetSnippetIndex(end[1])
			if snippetIndex == -1 {
				snippets.snippets = append(snippets.snippets, Snippet{id: end[1], end: index, start: -1})
			} else {
				snippets.snippets[snippetIndex].end = index
			}

			continue
		}

		start := snippetStart.FindStringSubmatch(scanner.Text())
		if len(start) == 2 {
			snippetIndex := snippets.GetSnippetIndex(start[1])
			if snippetIndex == -1 {
				snippets.snippets = append(snippets.snippets, Snippet{id: start[1], start: index, end: -1})
			} else {
				snippets.snippets[snippetIndex].start = index
			}
		}

		index++
	}

	return snippets
}
