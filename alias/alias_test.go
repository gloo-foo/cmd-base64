package alias_test

import (
	"slices"
	"testing"

	"github.com/gloo-foo/testable"

	b64 "github.com/gloo-foo/cmd-base64/alias"
)

// The alias package re-exports the constructor and flags under unprefixed names.
// A mis-wired re-export (Decode bound to the disabled constant, Wrap bound to
// the wrong type) compiles cleanly, so only behavior proves the wiring. Each
// test exercises one re-export and asserts the GNU base64 output it produces.

func assertLines(t *testing.T, got, want []string) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func lines(t *testing.T, input string, opts ...any) []string {
	t.Helper()
	out, err := testable.TestLines(b64.Base64(opts...), input)
	if err != nil {
		t.Fatalf("Base64 failed: %v", err)
	}
	return out
}

func TestAlias_DefaultEncodes(t *testing.T) {
	assertLines(t, lines(t, "foobar"), []string{"Zm9vYmFy"})
}

func TestAlias_DecodeDecodes(t *testing.T) {
	// Decode must be the enabled -d constant: decoding the vector yields foobar.
	assertLines(t, lines(t, "Zm9vYmFy\n", b64.Decode), []string{"foobar"})
}

func TestAlias_NoDecodeEncodes(t *testing.T) {
	// NoDecode is the disabled form: it must behave like the default (encode).
	assertLines(t, lines(t, "foobar", b64.NoDecode), []string{"Zm9vYmFy"})
}

func TestAlias_IgnoreGarbageStripsNonAlphabet(t *testing.T) {
	assertLines(t, lines(t, "Zm9v!! YmFy\n", b64.Decode, b64.IgnoreGarbage), []string{"foobar"})
}

func TestAlias_NoIgnoreGarbageRejectsGarbage(t *testing.T) {
	// NoIgnoreGarbage is the default decode behavior: garbage is an error.
	if _, err := testable.TestLines(b64.Base64(b64.Decode, b64.NoIgnoreGarbage), "Zm9v!!YmFy\n"); err == nil {
		t.Fatal("expected error decoding garbage, got nil")
	}
}

func TestAlias_WrapColumnsOutput(t *testing.T) {
	// Wrap(10) must reach the -w column flag: a 16-char encoding splits 10 + 6.
	assertLines(t, lines(t, "Hello World", b64.Wrap(10)), []string{"SGVsbG8gV2", "9ybGQ="})
}
