package command

// base64DecodeFlag selects decode mode (-d / --decode). The zero value encodes.
type base64DecodeFlag bool

const (
	Base64Decode   base64DecodeFlag = true
	Base64NoDecode base64DecodeFlag = false
)

func (d base64DecodeFlag) Configure(flags *flags) { flags.decode = d }

// base64IgnoreGarbageFlag selects ignore-garbage mode (-i / --ignore-garbage):
// when decoding, characters outside the base64 alphabet are discarded rather
// than treated as an error. The zero value rejects garbage.
type base64IgnoreGarbageFlag bool

const (
	Base64IgnoreGarbage   base64IgnoreGarbageFlag = true
	Base64NoIgnoreGarbage base64IgnoreGarbageFlag = false
)

func (g base64IgnoreGarbageFlag) Configure(flags *flags) { flags.ignoreGarbage = g }

// Base64Wrap sets the column at which encoded output is wrapped (-w / --wrap).
// Base64Wrap(0) disables wrapping. When no Base64Wrap is supplied the GNU
// default of 76 applies (see defaultWrap).
type Base64Wrap uint

func (w Base64Wrap) Configure(flags *flags) {
	flags.wrap = wrapColumn(w)
	flags.wrapSet = true
}

// wrapColumn is the resolved wrap width carried in flags. Zero means no wrap.
type wrapColumn uint

type flags struct {
	wrap          wrapColumn
	decode        base64DecodeFlag
	ignoreGarbage base64IgnoreGarbageFlag
	wrapSet       bool
}
