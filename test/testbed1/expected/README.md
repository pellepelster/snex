# Example 1

## Include snippet1

<!-- insertSnippet[snippet1] -->
```
	var lines = []string{"unit", "tested", "code"}
	for line := range lines {
		println(line)
	}
```

<!-- /insertSnippet -->

## Include full file

<!-- insertFile[file1.go] -->
```
package input

func includeFullFile() {
	println("file1")
}
```

<!-- /insertFile -->
