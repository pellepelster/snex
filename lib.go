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
	content []string
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

func replaceSnippets(content string, snippets Snippets) string {
	lines := strings.Split(content,"\n")
	snippetsToReplace := parseDocument(content)

	for _, snippet := range snippetsToReplace.snippets {

		i := snippets.GetSnippetIndex(snippet.id)
		if i > -1 {
			temp := append(lines[:snippet.start+1], snippets.snippets[i].content...)
			lines = append(temp, lines[snippet.end:]...)
		}

	}

	return strings.Join(lines[:], "\n")
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
	var lines []string

	for scanner.Scan() {

		end := snippetEnd.FindStringSubmatch(scanner.Text())
		if len(end) == 2 {
			snippetIndex := snippets.GetSnippetIndex(end[1])
			if snippetIndex == -1 {
				snippets.snippets = append(snippets.snippets, Snippet{id: end[1], end: index, start: -1})
			} else {
				snippets.snippets[snippetIndex].end = index
				snippets.snippets[snippetIndex].content = lines
			}
			lines = nil

			continue
		} else {
			if lines != nil {
				lines = append(lines, scanner.Text())
			}
		}

		start := snippetStart.FindStringSubmatch(scanner.Text())
		if len(start) == 2 {
			snippetIndex := snippets.GetSnippetIndex(start[1])
			if snippetIndex == -1 {
				snippets.snippets = append(snippets.snippets, Snippet{id: start[1], start: index, end: -1})
				lines = []string {}
			} else {
				snippets.snippets[snippetIndex].start = index
			}
		}

		index++
	}

	return snippets
}