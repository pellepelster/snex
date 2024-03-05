package main

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/pellepelster/snex/pkg"
	"os"
	"path"
	"path/filepath"
	"unicode/utf8"
)

func listAllFiles(rootPath string) []string {
	var result []string

	err := filepath.Walk(rootPath, func(file string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if !info.IsDir() {
			result = append(result, file)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("%s", err)
	}

	return result
}

// borrowed from "golang.org/x/tools/godoc/util"
// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// IsText reports whether a significant prefix of s looks like correct UTF-8;
// that is, if it is likely that s is human-readable text.
func isText(s []byte) bool {

	const MAX = 1024 // at least utf8.UTFMax

	if len(s) > MAX {
		s = s[0:MAX]
	}

	for i, c := range string(s) {
		if i+utf8.UTFMax > len(s) {
			// last char may be incomplete - ignore
			break
		}
		if c == 0xFFFD || c < ' ' && c != '\n' && c != '\t' && c != '\f' {
			// decoding error or control character - not a text file
			return false
		}
	}

	return true
}

func fileReadHeadBytes(file string, n int) []byte {
	xfile, err := os.Open(file)

	if err != nil {
		panic(err)
	}

	defer xfile.Close()

	headBytes := make([]byte, n)
	m, err := xfile.Read(headBytes)
	if err != nil {
		panic(err)
	}

	return headBytes[:m]
}

func processFiles(folderOrFiles []string, template string) error {
	var files []string
	for _, folderOrFile := range folderOrFiles {
		log.Infof("collecting files from '%s'", folderOrFile)

		for _, file := range append(files, listAllFiles(folderOrFile)...) {
			headBytes := fileReadHeadBytes(path.Join(file), 1024)

			if isText(headBytes) {
				log.Infof("found text file '%s'", file)
				files = append(files, file)
			} else {
				log.Infof("ignoring non-text file '%s'", file)
			}
		}
	}

	var documents []pkg.ParsedDocument

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		document, err := pkg.ParseDocument(pkg.Document{File: file, Content: string(content)})
		if err != nil {
			return err
		}

		documents = append(documents, document)
	}

	errors := pkg.ValidateDocuments(documents)

	if len(errors) > 0 {
		for _, err := range errors {
			log.Error(err)
		}
		return fmt.Errorf("validating snippets failed")
	} else {
		log.Info("snippets successfully validated")
	}

	for _, document := range documents {
		snippetCount := pkg.CountSnippets(document)

		if snippetCount == 0 {
			log.Infof("no snippets found in '%s'", document.File)
		} else {
			log.Infof("found %d snippets in '%s'", snippetCount, document.File)

			for _, line := range document.Lines {
				if line.Snippet != nil && line.Snippet.IsSnippet && line.Snippet.IsStart {
					log.Infof("\t%s", line.Snippet.Id)
				}
			}
		}
	}

	replacedDocuments, err := pkg.ReplaceSnippets(documents, template)
	if err != nil {
		return err
	}

	for _, document := range replacedDocuments {
		file, err := os.Create(document.File)
		if err != nil {
			return err
		}

		_, err = file.Write([]byte(document.Content))
		if err != nil {
			return err
		}

		err = file.Close()
		if err != nil {
			return err
		}
	}

	log.Info("snippets successfully replaced")

	return nil
}

func fileOrDirExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
