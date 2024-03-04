# Example 1

## Include snippet1

<!--- insertSnippet: snippet1 -->
start
	println("snippet1")
end
<!--- /insertSnippet: snippet1 -->

## Include full file

<!--- insertFile: file1.go -->
start
package input

func includeFullFile() {
	println("file1")
}

end
<!--- /insertFile: file1.go -->