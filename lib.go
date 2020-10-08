package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type Snippet struct {
	id      string
	start   int
	end     int
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
	originalLines := strings.Split(content, "\n")
	lines := []string{}
	snippetsToReplace := parseDocument(content)

	for i := 0; i < len(snippetsToReplace.snippets); i++ {
		snippet := snippetsToReplace.snippets[i]

		isFirst := i == 0
		prefix := []string{}

		if isFirst {
			prefix = originalLines[:snippet.start+1]
		} else {
			lastSnippet := snippetsToReplace.snippets[i-1]
			prefix = originalLines[lastSnippet.end+1 : snippet.start+1]
		}

		isLast := !(i < len(snippetsToReplace.snippets)-1)
		postfix := []string{}

		if isLast {
			postfix = originalLines[snippet.end:]
		} else {
			//nextSnippet := snippetsToReplace.snippets[i+1]
			postfix = originalLines[snippet.end : snippet.end+1]
		}

		lines = append(lines, prefix...)

		i := snippets.GetSnippetIndex(snippet.id)
		if i > -1 {
			replacement := snippets.snippets[i]
			lines = append(lines, replacement.content...)
		}
		lines = append(lines, postfix...)
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

	snippetStart := regexp.MustCompile(`snippet:([a-zA-Z0-9_-]*)`)
	snippetEnd := regexp.MustCompile(`/snippet:([a-zA-Z0-9_-]*)`)

	snippets := Snippets{}

	scanner := bufio.NewScanner(strings.NewReader(content))
	index := -1
	var lines []string

	for scanner.Scan() {
		index++

		end := snippetEnd.FindStringSubmatch(scanner.Text())
		if len(end) == 2 {
			snippetIndex := snippets.GetSnippetIndex(end[1])
			if snippetIndex == -1 {
				snippets.snippets = append(snippets.snippets, Snippet{id: end[1], end: index + (index*1), start: -1})
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
				lines = []string{}
			} else {
				snippets.snippets[snippetIndex].start = index
			}
		}
	}

	return snippets
}

func listAllFiles(rootPath string) []string {
	result := []string{}

	err := filepath.Walk(rootPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				result = append(result, path)
			}

			return nil
		})

	if err != nil {
		log.Fatalf("%s", err)
	}

	return result
}

var std = log.New(os.Stderr, "", log.LstdFlags)

func Fatalf(code int, format string, v ...interface{}) {
	std.Output(2, fmt.Sprintf(format, v...))
	os.Exit(code)
}

func fullPath(p string) string {
	if path.IsAbs(p) {
		return p
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return strings.TrimRight(path.Join(dir, p), "/")
}
