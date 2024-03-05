![snex](https://github.com/pellepelster/snex/workflows/snex/badge.svg)

# SNippet EXtractor (snex)

snex keep the code snippets inside your documentation in sync with real code from your sources. The issues this solves is that source examples in documentation tend to get outdated very quick. By pulling the snippets directly from a working project you can make sure the source examples used in you docs are always up to date.

## Downloads
* Linux (AMD64) [snex_linux_amd64](https://github.com/pellepelster/snex/releases/latest/download/snex_linux_amd64)
* Windows (AMD64) [snex_windows_amd64](https://github.com/pellepelster/snex/releases/latest/download/snex_windows_amd64)
* Darwin (AMD64) [snex_darwin_amd64](https://github.com/pellepelster/snex/releases/latest/download/snex_darwin_amd64)

## How it works
`snex` works line based on text files and iterates through all files inside one or more folders. It searches for snippet start- and end-markers and replaces all content between those markers with the referenced snippet.

<!-- To keep things simple and language agnostic it does not care for comment markers (which differ between languages) and just looks for snippet start- and end-markers. -->


### Example 1

Given the following files (see example folder `examples/example1`)

**examples/example1/README.md**
```markdown
# Example 1

## Include snippet1

<!--- insertSnippet: snippet1 -->
<!--- /insertSnippet: snippet1 -->

## Include full file

<!--- insertFile: file1.go -->
<!--- /insertFile: file1.go -->
```

**examples/example1/src/snippets.go**
```go
package input

func snippet1() {
	// snippet: snippet1
	println("snippet1")
	//  /snippet: snippet1
}
```

**examples/example1/src/file1.go**
```go
package input

func includeFullFile() {
	println("file1")
}
```

after running `snex` inside the `example1` folder via

```shell
snex ./examples/example1
```

The `README.md` looks like this

**examples/example1/README.md**
````markdown
# Example 1

## Include snippet1

<!--- insertSnippet: snippet1 -->
```
println("snippet1")
```
<!--- /insertSnippet: snippet1 -->

## Include full file

<!--- insertFile: file1.go -->
```
package input

func includeFullFile() {
println("file1")
}

```
<!--- /insertFile: file1.go -->
````

`snex` will keep the original markers to ensure it can be re-run anytime on the documentation sources. 

## Features

### Templates

`snex` has default replacement templates for different well-known files extensions. E.g. replacements inside a `.md` will automatically be surrounded by the markdown code block markers.

You can override the used template with

```shell
snex --template 'begin\n{{.Content}}\nend' ./
```

To show the list of default templates run

```shell
snex show-templates
```


