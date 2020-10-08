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

	var allSnippets = Snippets{}

	for _, file := range listAllFiles(snippetsPath) {

		if sourcePath == file {
			continue
		}

		content, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}

		snippets := parseDocument(string(content))

		if len(snippets.snippets) > 0 {
			log.Printf("file '%s' contains %d snippet(s)", file, len(snippets.snippets))

			/*
			for _, t := range snippets.snippets {
				fmt.Printf("%s: %d -> %d\n", t.id, t.start, t.end)
			}
			 */

			allSnippets.snippets = append(allSnippets.snippets, snippets.snippets...)
		}
	}

	if dirExists(sourcePath) {
		for _, sourceFile := range listAllFiles(sourcePath) {
			renderFile(sourceFile, targetPath, sourcePath, allSnippets)
		}
	} else {
		renderFile(sourcePath, sourcePath, sourcePath, allSnippets)
	}
}

func renderFile(sourceFile string, targetPath string, sourcePath string, allSnippets Snippets) {
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
