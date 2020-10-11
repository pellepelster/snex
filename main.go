package main

import (
	"flag"
	"github.com/jpillora/longestcommon"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func dirExists(filename string) bool {
	info, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}

func main() {
	var snippetsParameter string
	flag.StringVar(&snippetsParameter, "snippets", "", "directory that contains all snippets/sources")

	var sourceParameter string
	flag.StringVar(&sourceParameter, "source", "", "directory that contains all files to populate with snippets")

	var targetParameter string
	flag.StringVar(&targetParameter, "target", "", "target directory for generated files")

	var template string
	flag.StringVar(&template, "template", "{{.Content}}", "template to use for rendering snippet content")

	var templateFile string
	flag.StringVar(&templateFile, "template-file", "", "template file to use for rendering snippet content")

	flag.Parse()

	if sourceParameter == "" {
		Fatalf(1, "no source directory specified")
	}
	sourcePath := fullPath(sourceParameter)

	if snippetsParameter == "" {
		Fatalf(1, "no snippets directory specified")
	}
	snippetsPath := fullPath(snippetsParameter)

	if !fileExists(sourcePath) && targetParameter == "" {
		Fatalf(1, "no target directory specified")
	}
	targetPath := fullPath(targetParameter)

	if !dirExists(sourcePath) && !fileExists(sourcePath) {
		Fatalf(2, "source path '%s' not found\n", snippetsParameter)
	}

	if !dirExists(targetPath) && !fileExists(sourcePath) {
		Fatalf(2, "target path '%s' not found\n", snippetsParameter)
	}

	if !dirExists(snippetsPath) {
		Fatalf(2, "snippets path '%s' not found\n", snippetsParameter)
	}

	if !fileExists(sourcePath) && (strings.HasPrefix(targetPath, snippetsPath) || strings.HasPrefix(snippetsPath, targetPath)) {
		Fatalf(3, "snippets path '%s' and target path '%s' are not distinct\n", snippetsPath, targetPath)
	}

	snippetTemplate := template

	if templateFile != "" {
		if !fileExists(templateFile) {
			Fatalf(3, "template file '%s' not found\n", templateFile)

		}

		content, err := ioutil.ReadFile(templateFile)
		if err != nil {
			Fatalf(99, "unable to read template file '%s' (%s)", templateFile, err)
		}

		snippetTemplate = string(content)
	}

	var parsedDocuments []ParsedDocument

	for _, file := range listAllFiles(snippetsPath) {

		if sourcePath == file {
			continue
		}

		parsedDocument := parseFile(file)
		if len(parsedDocument.snippets) > 0 {
			log.Printf("file '%s' contains %d snippet(s)", file, len(parsedDocument.snippets))
		}
		parsedDocuments = append(parsedDocuments, parsedDocument)
	}

	if dirExists(sourcePath) {
		for _, sourceFile := range listAllFiles(sourcePath) {
			renderFile(sourceFile, targetPath, sourcePath, snippetsPath, snippetTemplate, parsedDocuments)
		}
	} else {
		renderFile(sourcePath, sourcePath, sourcePath, snippetsPath, snippetTemplate, parsedDocuments)
	}
}

func renderFile(sourceFile string, targetPath string, sourcePath string, snippetPath string, snippetTemplate string, parsedDocuments []ParsedDocument) {
	content, err := ioutil.ReadFile(sourceFile)

	if err != nil {
		Fatalf(99, "unable to read file '%s' (%s)", sourceFile, err)
	}

	prefix := longestcommon.Prefix([]string{sourceFile, targetPath})
	prefix = strings.TrimRight(prefix, "/")

	targetFile := path.Join(targetPath, strings.TrimPrefix(sourceFile, sourcePath))
	log.Printf("rendering file '%s' to '%s'", sourceFile, targetFile)

	os.MkdirAll(path.Dir(targetFile), os.ModePerm)
	err = ioutil.WriteFile(targetFile, []byte(replaceSnippets(string(content), snippetPath, snippetTemplate, parsedDocuments)), 0644)
	if err != nil {
		Fatalf(99, "unable to write file '%s' (%s)", targetFile, err)
	}
}
