# Example 1

## Include snippet1

<!-- insertSnippet[snippet1] -->
start
	var lines = []string{"unit", "tested", "code"}
	for line := range lines {
		println(line)
	}
end
<!-- /insertSnippet -->

## Include full file

<!-- insertFile[file1.go] -->
start
package input

func includeFullFile() {
	println("file1")
}
end
<!-- /insertFile -->