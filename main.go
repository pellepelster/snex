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

	var verboseParameter bool
	flag.BoolVar(&verboseParameter, "verbose", false, "log detailed information about the snippet extraction process")

	var snippetsParameter string
	flag.StringVar(&snippetsParameter, "snippets", "", "directory that contains all snippets/sources")

	var sourceParameter string
	flag.StringVar(&sourceParameter, "source", "", "directory that contains all files to populate with snippets")

	var targetParameter string
	flag.StringVar(&targetParameter, "target", "", "target directory for generated files")

	var templateHelp = "available variables: \n" +
		"\t{{.Content}}: the snippet content\n" +
		"\t{{.Filename}}: the file content originated from\n" +
		"\t{{.IsFullFile}}: true if the snippet was a file include" +
		"\t{{.Start}}: include start line\n" +
		"\t{{.End}}: include end line\n\n"

	var template string
	flag.StringVar(&template, "template", "{{.Content}}", "template to use for rendering snippet content\n" + templateHelp)

	var templateFile string
	flag.StringVar(&templateFile, "template-file", "", "template file to use for rendering snippet content\n" + templateHelp)

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
		Fatalf(2, "source path '%s' not found\n", sourcePath)
	}

	if !dirExists(targetPath) && !fileExists(sourcePath) {
		Fatalf(2, "target path '%s' not found\n", targetParameter)
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
		    if (verboseParameter) {
    			log.Printf("ignoring file '%s'", file)
		    }
			continue
		}

		parsedDocument := parseFile(file, verboseParameter)
		if len(parsedDocument.snippets) > 0 {
			log.Printf("file '%s' contains %d snippet(s)", file, len(parsedDocument.snippets))
		}

		parsedDocuments = append(parsedDocuments, parsedDocument)
	}

	for _, parsedDocument := range parsedDocuments {
		for _, snippet := range parsedDocument.snippets {
			if snippet.filename != "" && !fileExists(path.Join(snippet.filename)) && snippet.end > -1 && snippet.start > -1 {
				Fatalf(4, "file include '%s' not found", snippet.filename)
			}
		}
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
