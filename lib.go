package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

type ParsedDocument struct {
	snippets []Snippet
	file     string
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

func GetSnippet(id string, parsedDocuments []ParsedDocument) (string, Snippet) {
	for _, parsedDocument := range parsedDocuments {
		for _, snippet := range parsedDocument.snippets {
			if snippet.id == id {
				return parsedDocument.file, snippet
			}
		}
	}

	return "", Snippet{}
}

func parseFile(file string) ParsedDocument {

	content, err := ioutil.ReadFile(file)
	if err != nil {
		Fatalf(99, "unable to read file '%s' (%s)", file, err)
	}

	snippets, err := parseSnippets(string(content))
	if err != nil {
		Fatalf(5, "%s", err)
	}

	return ParsedDocument{snippets: snippets, file: file}
}

func replaceSnippets(content string, basePath string, snippetTemplate string, parsedDocuments []ParsedDocument) string {
	originalLines := strings.Split(content, "\n")
	var lines []string
	snippetsToReplace, err := parseSnippets(content)
	if err != nil {
		Fatalf(5, "%s", err)
	}


	for i := 0; i < len(snippetsToReplace); i++ {
		snippetToReplace := snippetsToReplace[i]

		var prefix []string
		var postfix []string

		isFirst := i == 0

		if isFirst {
			prefix = originalLines[:snippetToReplace.start+1]
		} else {
			lastSnippet := snippetsToReplace[i-1]
			prefix = originalLines[lastSnippet.end+1 : snippetToReplace.start+1]
		}

		isLast := !(i < len(snippetsToReplace)-1)

		if isLast {
			postfix = originalLines[snippetToReplace.end:]
		} else {
			postfix = originalLines[snippetToReplace.end : snippetToReplace.end+1]
		}

		lines = append(lines, prefix...)

		type TemplateData struct {
			Content    string
			Filename   string
			IsFullFile bool
			Start      int
			End        int
		}

		filename, snippetToInsert := getSnippet(snippetToReplace, basePath, parsedDocuments)

		templateData := TemplateData{Content: strings.Join(snippetToInsert.content, "\n")}

		templateData.Filename = strings.TrimPrefix(filename, basePath)
		templateData.Filename = strings.TrimPrefix(templateData.Filename, "/")


		if snippetToInsert.filename != "" {
			templateData.Start = -1
			templateData.End = -1
			templateData.IsFullFile = true
		} else {
			templateData.Start = snippetToInsert.start + 2
			templateData.End = snippetToInsert.end
			templateData.IsFullFile = false
		}

		tmpl, err := template.New("test").Parse(snippetTemplate)

		if err != nil {
			Fatalf(99, "could not parse template: %s", err)
		}

		var renderedTemplate bytes.Buffer
		err = tmpl.Execute(&renderedTemplate, templateData)
		if err != nil {
			Fatalf(99, "could not execute template: %s", err)
		}

		renderedLines := strings.Split(renderedTemplate.String(), "\n")

		lines = append(lines, renderedLines...)
		lines = append(lines, postfix...)
	}

	return strings.Join(lines[:], "\n")
}

func getSnippet(snippet Snippet, basePath string, parsedDocuments []ParsedDocument) (string, Snippet) {

	if len(snippet.filename) > 0 {

		content, err := ioutil.ReadFile(path.Join(basePath, snippet.filename))
		if err != nil {
			log.Fatal(err)
		}

		return snippet.filename, Snippet{filename: snippet.filename, content: strings.Split(string(content), "\n")}
	} else {
		return GetSnippet(snippet.id, parsedDocuments)
	}

	return "", Snippet{}
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

func parseSnippets(content string) ([]Snippet, error) {

	var snippets []Snippet

	scanner := bufio.NewScanner(strings.NewReader(content))
	index := -1
	var lines []string

	for scanner.Scan() {
		index++

		snippetMarker := parseSnippetMarker(scanner.Text())

		if snippetMarker.isEnd {
			snippetIndex := GetSnippetIndex(snippets, snippetMarker.id)
			if snippetIndex == -1 {
				if snippetMarker.isFile {
					snippets = append(snippets, Snippet{id: snippetMarker.id, filename: snippetMarker.id, end: index + (index * 1), start: -1})
				} else {
					snippets = append(snippets, Snippet{id: snippetMarker.id, end: index + (index * 1), start: -1})
				}
			} else {
				snippets[snippetIndex].end = index
				snippets[snippetIndex].content = lines
			}
			lines = nil

			continue
		} else {
			if lines != nil {
				lines = append(lines, scanner.Text())
			}
		}

		if snippetMarker.isStart {
			snippetIndex := GetSnippetIndex(snippets, snippetMarker.id)
			if snippetIndex == -1 {
				if snippetMarker.isFile {
					snippets = append(snippets, Snippet{id: snippetMarker.id, filename: snippetMarker.id, start: index, end: -1})
				} else {
					snippets = append(snippets, Snippet{id: snippetMarker.id, start: index, end: -1})
				}

				lines = []string{}
			} else {
				snippets[snippetIndex].start = index
			}

			continue
		}
	}

	for _, snippet := range snippets {
		if snippet.end == -1 || snippet.start == -1 {
			return []Snippet{}, errors.New("unbalanced snippet markers")
		}
	}

	return snippets, nil
}

func listAllFiles(rootPath string) []string {
	var result []string

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
