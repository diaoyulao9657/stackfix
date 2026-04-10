package main

import "testing"

func TestDetect(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"python traceback", `Traceback (most recent call last):
  File "app.py", line 12, in <module>
    run()
  File "app.py", line 8, in run
    return data["key"]
KeyError: 'key'`, "python"},

		{"python import error", `ModuleNotFoundError: No module named 'pandas'`, "python"},

		{"node error", `TypeError: Cannot read properties of undefined (reading 'map')
    at Object.<anonymous> (/app/index.js:15:20)
    at Module._compile (node:internal/modules/cjs/loader:1469:14)`, "javascript"},

		{"go panic", `goroutine 1 [running]:
main.handler(0x0?)
	/app/server.go:28 +0x1a4
panic: runtime error: invalid memory address or nil pointer dereference`, "go"},

		{"rust panic", `thread 'main' panicked at 'index out of bounds: the len is 3 but the index is 5', src/main.rs:10:5
note: run with RUST_BACKTRACE=1`, "rust"},

		{"java npe", `Exception in thread "main" java.lang.NullPointerException
	at com.example.App.process(App.java:42)
	at com.example.App.main(App.java:15)`, "java"},

		{"c segfault", `Segmentation fault (core dumped)`, "c/c++"},

		{"unknown", `something broke lol`, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detect(tt.input)
			if got != tt.want {
				t.Errorf("detect(%q) = %q, want %q", tt.name, got, tt.want)
			}
		})
	}
}

func TestContainsAny(t *testing.T) {
	if !containsAny("hello world", "world", "foo") {
		t.Error("expected true")
	}
	if containsAny("hello world", "foo", "bar") {
		t.Error("expected false")
	}
}
