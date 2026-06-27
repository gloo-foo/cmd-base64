package command

import (
	"bytes"
	"encoding/base64"

	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"
)

// defaultWrap is the column at which GNU base64 wraps encoded output by default.
const defaultWrap wrapColumn = 76

// Base64 returns a command that base64-encodes or -decodes its whole input,
// matching GNU coreutils base64.
//
// Encode (default): the entire input is treated as one byte stream and encoded;
// output is wrapped at column 76, or at the column given by Base64Wrap(COLS).
// Base64Wrap(0) disables wrapping.
//
// Decode (Base64Decode / -d): the entire input is concatenated and decoded.
// Without Base64IgnoreGarbage, a byte outside the base64 alphabet is an error;
// with it (-i / --ignore-garbage), such bytes are discarded before decoding.
func Base64(opts ...any) gloo.Command[[]byte, []byte] {
	f := gloo.NewParameters[gloo.File, flags](opts...).Flags
	return patterns.Accumulate(transform(f))
}

// transform selects the encode or decode accumulator for the resolved flags.
func transform(f flags) func([][]byte) ([][]byte, error) {
	if bool(f.decode) {
		return decode(bool(f.ignoreGarbage))
	}
	return encode(wrapAt(f))
}

// wrapAt resolves the effective wrap column: an explicit -w value when set,
// otherwise the GNU default of 76.
func wrapAt(f flags) wrapColumn {
	if f.wrapSet {
		return f.wrap
	}
	return defaultWrap
}

// encode returns an accumulator that joins all input lines into one byte stream,
// base64-encodes it, and wraps the result at the given column.
func encode(cols wrapColumn) func([][]byte) ([][]byte, error) {
	return func(lines [][]byte) ([][]byte, error) {
		if len(lines) == 0 {
			return [][]byte{}, nil
		}
		encoded := base64.StdEncoding.EncodeToString(bytes.Join(lines, newline))
		return wrap(encoded, cols), nil
	}
}

// decode returns an accumulator that concatenates all input lines and decodes
// the result. When ignoreGarbage is set, non-alphabet bytes are stripped first.
func decode(ignoreGarbage bool) func([][]byte) ([][]byte, error) {
	return func(lines [][]byte) ([][]byte, error) {
		if len(lines) == 0 {
			return [][]byte{}, nil
		}
		joined := bytes.Join(lines, nil)
		if ignoreGarbage {
			joined = stripGarbage(joined)
		}
		decoded, err := base64.StdEncoding.DecodeString(string(joined))
		if err != nil {
			return nil, err
		}
		return [][]byte{decoded}, nil
	}
}

// newline is the separator used to rejoin the line-split input into the original
// byte stream before encoding.
var newline = []byte{'\n'}

// wrap splits encoded into lines of at most cols characters. A cols of zero
// (GNU's -w 0) emits the encoding as a single unwrapped line.
func wrap(encoded string, cols wrapColumn) [][]byte {
	if cols == 0 {
		return [][]byte{[]byte(encoded)}
	}
	width := int(cols)
	lines := make([][]byte, 0, len(encoded)/width+1)
	for len(encoded) > width {
		lines = append(lines, []byte(encoded[:width]))
		encoded = encoded[width:]
	}
	return append(lines, []byte(encoded))
}

// stripGarbage removes every byte that is not part of the standard base64
// alphabet (A-Z, a-z, 0-9, +, /, =), implementing -i / --ignore-garbage.
func stripGarbage(in []byte) []byte {
	out := make([]byte, 0, len(in))
	for _, b := range in {
		if isBase64Byte(b) {
			out = append(out, b)
		}
	}
	return out
}

// isBase64Byte reports whether b belongs to the standard base64 alphabet,
// including the '=' padding character.
func isBase64Byte(b byte) bool {
	switch {
	case b >= 'A' && b <= 'Z', b >= 'a' && b <= 'z', b >= '0' && b <= '9':
		return true
	default:
		return b == '+' || b == '/' || b == '='
	}
}
