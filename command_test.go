package command_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gloo-foo/testable"

	command "github.com/gloo-foo/cmd-base64"
)

// These tests pin GNU coreutils base64 behavior with fixed vectors, not
// round-trips: encode treats the whole input as one byte stream (lines rejoined
// with '\n'), wraps at column 76 by default, and decode concatenates all input.

func run(t *testing.T, input string, opts ...any) []string {
	t.Helper()
	lines, err := testable.TestLines(command.Base64(opts...), input)
	if err != nil {
		t.Fatalf("Base64 failed: %v", err)
	}
	return lines
}

func TestEncode_RFC4648Vectors(t *testing.T) {
	// RFC 4648 section 10 vectors. Single-line inputs carry no trailing newline
	// through the line-stream model, so they encode the literal bytes.
	cases := map[string]string{
		"":       "",
		"f":      "Zg==",
		"fo":     "Zm8=",
		"foo":    "Zm9v",
		"foob":   "Zm9vYg==",
		"fooba":  "Zm9vYmE=",
		"foobar": "Zm9vYmFy",
	}
	for in, want := range cases {
		lines := run(t, in)
		got := strings.Join(lines, "\n")
		if got != want {
			t.Errorf("encode(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestEncode_MultiLineIsOneStream(t *testing.T) {
	// GNU base64 encodes the whole stream, not each line: the two lines are
	// rejoined with '\n' before encoding. base64("hello\nworld") == this vector.
	lines := run(t, "hello\nworld\n")
	got := strings.Join(lines, "\n")
	want := "aGVsbG8Kd29ybGQ="
	if got != want {
		t.Errorf("encode multiline = %q, want %q", got, want)
	}
}

func TestEncode_WrapsAt76ByDefault(t *testing.T) {
	in := strings.Repeat("a", 120) // 120 bytes -> 160 base64 chars -> 76 + 76 + 8
	lines := run(t, in)
	if len(lines) != 3 {
		t.Fatalf("expected 3 wrapped lines, got %d: %v", len(lines), lines)
	}
	for i, l := range lines[:2] {
		if len(l) != 76 {
			t.Errorf("line %d length = %d, want 76", i, len(l))
		}
	}
	if len(lines[2]) != 8 {
		t.Errorf("final line length = %d, want 8", len(lines[2]))
	}
}

func TestEncode_CustomWrapColumn(t *testing.T) {
	lines := run(t, "The quick brown fox jumps over the lazy dog.", command.Base64Wrap(10))
	for i, l := range lines[:len(lines)-1] {
		if len(l) != 10 {
			t.Errorf("line %d length = %d, want 10", i, len(l))
		}
	}
	if got := strings.Join(lines, ""); got != "VGhlIHF1aWNrIGJyb3duIGZveCBqdW1wcyBvdmVyIHRoZSBsYXp5IGRvZy4=" {
		t.Errorf("rejoined encoding = %q", got)
	}
}

func TestEncode_WrapZeroDisablesWrapping(t *testing.T) {
	in := strings.Repeat("a", 120)
	lines := run(t, in, command.Base64Wrap(0))
	if len(lines) != 1 {
		t.Fatalf("expected 1 unwrapped line, got %d", len(lines))
	}
	if len(lines[0]) != 160 {
		t.Errorf("line length = %d, want 160", len(lines[0]))
	}
}

func TestDecode_Vector(t *testing.T) {
	lines := run(t, "Zm9vYmFy\n", command.Base64Decode)
	if len(lines) != 1 || lines[0] != "foobar" {
		t.Fatalf("decode = %q, want [foobar]", lines)
	}
}

func TestDecode_ConcatenatesWrappedInput(t *testing.T) {
	// A decoder must rejoin wrapped lines into one stream before decoding:
	// "Zm9v" + "YmFy" decodes to "foobar", not two independent fragments.
	lines := run(t, "Zm9v\nYmFy", command.Base64Decode)
	if len(lines) != 1 || lines[0] != "foobar" {
		t.Fatalf("decode wrapped = %q, want [foobar]", lines)
	}
}

func TestDecode_RejectsGarbageByDefault(t *testing.T) {
	_, err := testable.TestLines(command.Base64(command.Base64Decode), "Zm9v!!YmFy\n")
	if err == nil {
		t.Fatal("expected error decoding garbage, got nil")
	}
}

func TestDecode_IgnoreGarbage(t *testing.T) {
	lines := run(t, "Zm9v!! Ym\tFy\n", command.Base64Decode, command.Base64IgnoreGarbage)
	if len(lines) != 1 || lines[0] != "foobar" {
		t.Fatalf("decode ignore-garbage = %q, want [foobar]", lines)
	}
}

func TestDecode_EmptyInput(t *testing.T) {
	lines := run(t, "", command.Base64Decode)
	if len(lines) != 0 {
		t.Fatalf("decode empty = %q, want no lines", lines)
	}
}

func TestEncode_EmptyInput(t *testing.T) {
	lines := run(t, "")
	if len(lines) != 0 {
		t.Fatalf("encode empty = %q, want no lines", lines)
	}
}

func ExampleBase64() {
	lines, _ := testable.TestLines(command.Base64(), "Hello World\n")
	for _, line := range lines {
		fmt.Println(line)
	}
	// Output:
	// SGVsbG8gV29ybGQ=
}

func ExampleBase64_decode() {
	lines, _ := testable.TestLines(command.Base64(command.Base64Decode), "SGVsbG8gV29ybGQ=\n")
	for _, line := range lines {
		fmt.Println(line)
	}
	// Output:
	// Hello World
}

func ExampleBase64_wrap() {
	lines, _ := testable.TestLines(command.Base64(command.Base64Wrap(10)), "Hello World\n")
	for _, line := range lines {
		fmt.Println(line)
	}
	// Output:
	// SGVsbG8gV2
	// 9ybGQ=
}
