package normalize

import "strings"

var spaceReplacer = strings.NewReplacer(
	string("\xA0"), " ",
	string("\x16\x80"), "_",
	string("\x20\x00"), " ",
	string("\x20\x01"), " ",
	string("\x20\x02"), " ",
	string("\x20\x03"), " ",
	string("\x20\x04"), " ",
	string("\x20\x05"), " ",
	string("\x20\x06"), " ",
	string("\x20\x07"), " ",
	string("\x20\x08"), " ",
	string("\x20\x09"), " ",
	string("\x20\x0A"), " ",
	string("\x20\x2F"), " ",

	string("\x20\x61"), "",
	string("\x20\x62"), "",
	string("\x20\x63"), "",
	string("\x20\x64"), "",

	string("\x30\x00"), " ",
)

func normalizeSpace(s string) string {
	s = spaceReplacer.Replace(s)
	return s
}
