package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
)

func fullPath(p string) string {
	if path.IsAbs(p) {
		return p
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return path.Join(dir, p)
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
		fmt.Println("no source directory specified")
		os.Exit(1)
	}

	if snippetsParameter == "" {
		fmt.Println("no snippets directory specified")
		os.Exit(1)
	}

	if targetParameter == "" {
		fmt.Println("no target directory specified")
		os.Exit(1)
	}

	sourcePath := fullPath(sourceParameter)
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		fmt.Printf("source path '%s' not found\n", snippetsParameter)
		os.Exit(2)
	}

	targetPath := fullPath(targetParameter)
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		log.Printf("target path '%s' not found\n", snippetsParameter)
		os.Exit(2)
	}

	snippetsPath := fullPath(snippetsParameter)
	if _, err := os.Stat(snippetsPath); os.IsNotExist(err) {
		log.Printf("snippets path '%s' not found\n", snippetsParameter)
		os.Exit(2)
	}

}
