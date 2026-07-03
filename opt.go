package command

// base64DecodeFlag selects decode mode (-d / --decode). The zero value encodes.
type base64DecodeFlag bool

const (
	Base64Decode   base64DecodeFlag = true
	Base64NoDecode base64DecodeFlag = false
)

// base64IgnoreGarbageFlag selects ignore-garbage mode (-i / --ignore-garbage):
// when decoding, characters outside the base64 alphabet are discarded rather
// than treated as an error. The zero value rejects garbage.
type base64IgnoreGarbageFlag bool

const (
	Base64IgnoreGarbage   base64IgnoreGarbageFlag = true
	Base64NoIgnoreGarbage base64IgnoreGarbageFlag = false
)

// Base64Wrap sets the column at which encoded output is wrapped (-w / --wrap).
// Base64Wrap(0) disables wrapping. When no Base64Wrap is supplied the GNU
// default of 76 applies (see defaultWrap).
type Base64Wrap uint

// wrapColumn is the resolved wrap width carried in flags. Zero means no wrap.
type wrapColumn uint

type flags struct {
	wrap                 wrapColumn
	decodeEnabled        base64DecodeFlag
	ignoreGarbageEnabled base64IgnoreGarbageFlag
	isWrapSet            bool
}

// fold partitions opts: base64's own option values are folded into the flag
// set, and every other argument is passed through unchanged for the
// framework's positional classifier.
func fold(opts []any) (flags, []any) {
	var f flags
	rest := make([]any, 0, len(opts))
	for _, o := range opts {
		switch v := o.(type) {
		case base64DecodeFlag:
			f.decodeEnabled = v
		case base64IgnoreGarbageFlag:
			f.ignoreGarbageEnabled = v
		case Base64Wrap:
			f.wrap = wrapColumn(v)
			f.isWrapSet = true
		default:
			rest = append(rest, o)
		}
	}
	return f, rest
}
