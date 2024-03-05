![snex](https://github.com/pellepelster/snex/workflows/snex/badge.svg)

# SNippet EXtractor (snex)

`snex` keeps the code snippets inside your documentation in sync with real code from your sources. The issue this solves is that source examples in documentation tend to get outdated very quick. By pulling the snippets directly from a working project you can make sure the source examples used in you docs are always up-to-date and functional.

## Downloads

* Linux (amd64) [snex_linux_amd64](https://github.com/pellepelster/snex/releases/latest/download/snex_linux_amd64)
* Linux (arm64) [snex_linux_arm64](https://github.com/pellepelster/snex/releases/latest/download/snex_linux_arm64)
* Linux (386) [snex_linux_386](https://github.com/pellepelster/snex/releases/latest/download/snex_linux_386)
* Windows (amd64) [snex_windows_amd64](https://github.com/pellepelster/snex/releases/latest/download/snex_windows_amd64)
* Windows (arm64) [snex_windows_arm64](https://github.com/pellepelster/snex/releases/latest/download/snex_windows_arm64)
* Windows (386) [snex_windows_386](https://github.com/pellepelster/snex/releases/latest/download/snex_windows_386)
* Darwin (amd64) [snex_darwin_amd64](https://github.com/pellepelster/snex/releases/latest/download/snex_darwin_amd64)
* Darwin (arm64) [snex_darwin_arm64](https://github.com/pellepelster/snex/releases/latest/download/snex_darwin_arm64)
* Freebsd (amd64) [snex_freebsd_amd64](https://github.com/pellepelster/snex/releases/latest/download/snex_freebsd_amd64)

## How it works

`snex` works line based on text files and iterates through all files inside of one or more folders, where it searches for snippet start- and end-markers and replaces all content between those markers with the referenced snippet.

Those markers are language agnostic, so you can embed them in a way that the source file is not corrupted by the markers, typically you want to choose comments for that, 

e.g. for Java

```java
// snippet[snippet1]
```

or inside of HTML

```html
<!-- snippet[snippet1] -->
```

There are three types of markers available that must be opened and closed like HTML tags

* `snippet[${id}]` and `/snippet` define the beginning and end of a snippet that can be inserted somewhere else

* `insertSnippet[${id}]` and `/insertSnippet` define the bounds where the snipped with the id `${id}` will be inserted

* `insertFile[${file}]` and `/insertFile` define the bounds where the whole file `${file}` will be inserted

### Example 1

Given the following files (see also example folder `examples/example1`)

**examples/example1/README.md**
```markdown
# Example 1

## Include snippet1

<!-- insertSnippet[snippet1] -->
<!-- /insertSnippet -->

## Include full file

<!-- insertFile[file1.go] -->
<!-- /insertFile -->
```

**examples/example1/src/snippets.go**
```go
package input

func snippet1() {
	// snippet[snippet1]
	println("snippet1")
	//  /snippet
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

<!-- insertSnippet[snippet1] -->
```
println("snippet1")
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
````

`snex` will keep the original markers to ensure it can be re-run anytime on the documentation sources. 

## Features

### Templates

`snex` has default replacement templates for different well-known files extensions. E.g. replacements inside a `.md` will automatically be surrounded by markdown code block markers.

You can override the used template with

```shell
snex --template 'begin\n{{.Content}}\nend' ./
```

To show the list of default templates run

```shell
snex show-templates
```
