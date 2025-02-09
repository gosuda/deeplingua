package normalize

import "strings"

var punctuationReplacer = strings.NewReplacer(
	string("\x20\x10"), "-",
	string("\x20\x11"), "-",
	string("\x20\x12"), "-",
	string("\x20\x13"), "-",
	string("\x20\x14"), "-",
	string("\x20\x15"), "-",

	string("\x20\x18"), "'",
	string("\x20\x19"), "'",
	string("\x20\x1A"), "'",
	string("\x20\x1B"), "'",
	string("\x20\x1C"), "\"",
	string("\x20\x1D"), "\"",
	string("\x20\x1E"), "\"",
	string("\x20\x1F"), "\"",

	string("\x20\x24"), ".",
	string("\x20\x25"), "..",
	string("\x20\x26"), "...",

	string("\x20\x3C"), "!!",

	string("\x20\x3E"), "-",
	string("\x20\x43"), "-",
	string("\x20\x44"), "/",
	string("\x20\x47"), "??",

	string("\x20\x53"), "~",

	string("\x30\x1C"), "~",
	string("\x30\x1D"), "\"",
	string("\x30\x1E"), "\"",
)

func normalizePunctuation(s string) string {
	s = punctuationReplacer.Replace(s)
	return s
}
