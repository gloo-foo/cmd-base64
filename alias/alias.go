// Package alias provides unprefixed names for base64 command flags.
//
//	import b64 "github.com/gloo-foo/cmd-base64/alias"
//	b64.Base64(b64.Decode)
package alias

import command "github.com/gloo-foo/cmd-base64"

// Base64 re-exports the constructor.
var Base64 = command.Base64

// Decode is the -d / --decode flag.
const Decode = command.Base64Decode

// NoDecode is the default: encode.
const NoDecode = command.Base64NoDecode

// IgnoreGarbage is the -i / --ignore-garbage flag.
const IgnoreGarbage = command.Base64IgnoreGarbage

// NoIgnoreGarbage is the default: reject non-alphabet bytes when decoding.
const NoIgnoreGarbage = command.Base64NoIgnoreGarbage

// Wrap re-exports the -w / --wrap column flag type. Wrap(0) disables wrapping.
type Wrap = command.Base64Wrap
