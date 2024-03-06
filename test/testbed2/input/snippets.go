package input

import "testing"

func TestcaseForCodeSnippet1(t *testing.T) {
	// snippet[snippet1]
	var lines = []string{"unit", "tested", "code"}
	for line := range lines {
		println(line)
	}
	// /snippet
}
