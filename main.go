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

func main() {
	var snippetsParameter string
	flag.StringVar(&snippetsParameter, "snippets", "", "directory that contains all snippets/sources")

	var sourceParameter string
	flag.StringVar(&sourceParameter, "source", "", "directory that contains all files to populate with snippets")

	var targetParameter string
	flag.StringVar(&targetParameter, "target", "", "target directory for generated files")

	flag.Parse()

	if sourceParameter == "" {
		Fatalf(1, "no source directory specified")
	}

	if snippetsParameter == "" {
		Fatalf(1, "no snippets directory specified")
	}

	if targetParameter == "" {
		Fatalf(1, "no target directory specified")
	}

	sourcePath := fullPath(sourceParameter)
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		Fatalf(2, "source path '%s' not found\n", snippetsParameter)
	}

	targetPath := fullPath(targetParameter)
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		Fatalf(2, "target path '%s' not found\n", snippetsParameter)
	}

	snippetsPath := fullPath(snippetsParameter)
	if _, err := os.Stat(snippetsPath); os.IsNotExist(err) {
		Fatalf(2, "snippets path '%s' not found\n", snippetsParameter)
	}

	if strings.HasPrefix(targetPath, snippetsPath) || strings.HasPrefix(snippetsPath, targetPath) {
		Fatalf(3, "snippets path '%s' and target path '%s' are not distinct\n", snippetsPath, targetPath)
	}

	var allSnippets = Snippets{}

	for _, file := range listAllFiles(snippetsPath) {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}

		snippets := parseDocument(string(content))
		log.Printf("file '%s' contains %d snippet(s)", file, len(snippets.snippets))
		allSnippets.snippets = append(allSnippets.snippets, snippets.snippets...)
	}

	for _, sourceFile := range listAllFiles(sourcePath) {
		content, err := ioutil.ReadFile(sourceFile)

		if err != nil {
			log.Fatal(err)
		}

		prefix := longestcommon.Prefix([]string{sourceFile, targetPath})
		prefix = strings.TrimRight(prefix, "/")

		targetFile := path.Join(targetPath, strings.TrimPrefix(sourceFile, sourcePath))
		log.Printf("rendering file '%s' to '%s'", sourceFile, targetFile)

		os.MkdirAll(path.Dir(targetFile), os.ModePerm)
		error := ioutil.WriteFile(targetFile, []byte(replaceSnippets(string(content), allSnippets)), 0644)
		if error != nil {
			log.Fatalf("%s", error)
		}
	}
}
