package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type ParsedDocument struct {
	snippets []Snippet
}

type Snippet struct {
	id       string
	start    int
	end      int
	content  []string
	filename string
}

func GetSnippetIndex(snippets []Snippet, id string) int {
	for index, snippet := range snippets {
		if snippet.id == id {
			return index
		}
	}

	return -1
}

func replaceSnippets(content string, basePath string, snippets []Snippet) string {
	originalLines := strings.Split(content, "\n")
	var lines []string
	snippetsToReplace := parseDocument(content)

	for i := 0; i < len(snippetsToReplace.snippets); i++ {
		snippet := snippetsToReplace.snippets[i]

		isFirst := i == 0
		var prefix []string

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
			postfix = originalLines[snippet.end : snippet.end+1]
		}

		lines = append(lines, prefix...)

		if len(snippet.filename) > 0 {

			content, err := ioutil.ReadFile(path.Join(basePath, snippet.filename))
			if err != nil {
				log.Fatal(err)
			}

			lines = append(lines, strings.Split(string(content), "\n")...)
		} else {
			i := GetSnippetIndex(snippets, snippet.id)

			if i > -1 {
				replacement := snippets[i]
				lines = append(lines, replacement.content...)
			}
		}

		lines = append(lines, postfix...)
	}

	return strings.Join(lines[:], "\n")
}

var snippetStartExpression = regexp.MustCompile(`snippet:([a-zA-Z0-9_\-]*)`)
var snippetEndExpression = regexp.MustCompile(`/snippet:([a-zA-Z0-9_\-]*)`)

var fileStartExpression = regexp.MustCompile(`file:([a-zA-Z0-9_\-\/\.]*)`)
var fileEndExpression = regexp.MustCompile(`/file:([a-zA-Z0-9_\-\/\.]*)`)

type SnippetMarker struct {
	isSnippet bool
	isFile    bool
	isEnd     bool
	isStart   bool
	id        string
}

func parseSnippetMarker(line string) SnippetMarker {
	snippetEnd := snippetEndExpression.FindStringSubmatch(line)
	if len(snippetEnd) == 2 {
		return SnippetMarker{isSnippet: true, isEnd: true, id: snippetEnd[1]}
	}

	snippetStart := snippetStartExpression.FindStringSubmatch(line)
	if len(snippetStart) == 2 {
		return SnippetMarker{isSnippet: true, isStart: true, id: snippetStart[1]}
	}

	fileEnd := fileEndExpression.FindStringSubmatch(line)
	if len(fileEnd) == 2 {
		return SnippetMarker{isFile: true, isEnd: true, id: fileEnd[1]}
	}

	fileStart := fileStartExpression.FindStringSubmatch(line)
	if len(fileStart) == 2 {
		return SnippetMarker{isFile: true, isStart: true, id: fileStart[1]}
	}

	return SnippetMarker{}
}

func parseDocument(content string) ParsedDocument {

	result := ParsedDocument{}

	scanner := bufio.NewScanner(strings.NewReader(content))
	index := -1
	var lines []string

	for scanner.Scan() {
		index++

		snippetMarker := parseSnippetMarker(scanner.Text())

		if snippetMarker.isEnd {
			snippetIndex := GetSnippetIndex(result.snippets, snippetMarker.id)
			if snippetIndex == -1 {
				if snippetMarker.isFile {
					result.snippets = append(result.snippets, Snippet{id: snippetMarker.id, filename: snippetMarker.id, end: index + (index * 1), start: -1})
				} else {
					result.snippets = append(result.snippets, Snippet{id: snippetMarker.id, end: index + (index * 1), start: -1})
				}
			} else {
				result.snippets[snippetIndex].end = index
				result.snippets[snippetIndex].content = lines
			}
			lines = nil

			continue
		} else {
			if lines != nil {
				lines = append(lines, scanner.Text())
			}
		}

		if snippetMarker.isStart {
			snippetIndex := GetSnippetIndex(result.snippets, snippetMarker.id)
			if snippetIndex == -1 {
				if snippetMarker.isFile {
					result.snippets = append(result.snippets, Snippet{id: snippetMarker.id, filename: snippetMarker.id, start: index, end: -1})
				} else {
					result.snippets = append(result.snippets, Snippet{id: snippetMarker.id, start: index, end: -1})
				}

				lines = []string{}
			} else {
				result.snippets[snippetIndex].start = index
			}

			continue
		}
	}

	return result
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
