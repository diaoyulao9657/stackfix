package main

import "strings"

func detect(input string) string {
	lower := strings.ToLower(input)

	// python
	if strings.Contains(input, "Traceback (most recent call last)") ||
		strings.Contains(input, "File \"") && strings.Contains(input, ", line ") ||
		containsAny(lower, "syntaxerror:", "valueerror:", "keyerror:", "importerror:", "attributeerror:", "nameerror:", "indentationerror:", "modulenotfounderror:") {
		return "python"
	}

	// javascript/node
	if strings.Contains(input, "node_modules/") ||
		strings.Contains(input, "at Object.<anonymous>") ||
		strings.Contains(input, "at Module._compile") ||
		(strings.Contains(input, ".js:") || strings.Contains(input, ".ts:")) && containsAny(lower, "referenceerror:", "typeerror:", "syntaxerror:") {
		return "javascript"
	}

	// go
	if strings.Contains(input, "goroutine ") && strings.Contains(input, "[running]") ||
		strings.Contains(input, "panic:") ||
		strings.Contains(lower, "runtime error:") {
		return "go"
	}

	// rust
	if strings.Contains(input, "thread '") && strings.Contains(input, "panicked at") ||
		strings.Contains(input, "error[E") {
		return "rust"
	}

	// java/kotlin
	if strings.Contains(input, "Exception in thread") ||
		strings.Contains(input, "at java.") || strings.Contains(input, "at com.") || strings.Contains(input, "at org.") ||
		containsAny(lower, "nullpointerexception", "classnotfoundexception", "illegalargumentexception") {
		return "java"
	}

	// c/c++
	if containsAny(lower, "segmentation fault", "undefined reference", "linker command failed") ||
		strings.Contains(input, "error:") && (strings.Contains(input, ".c:") || strings.Contains(input, ".cpp:") || strings.Contains(input, ".h:")) {
		return "c/c++"
	}

	return ""
}

func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
