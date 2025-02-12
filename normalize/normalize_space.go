package normalize

import "strings"

var spaceReplacer = strings.NewReplacer(
	string("\u00A0"), " ",
	string("\u1680"), "_",
	string("\u2000"), " ",
	string("\u2001"), " ",
	string("\u2002"), " ",
	string("\u2003"), " ",
	string("\u2004"), " ",
	string("\u2005"), " ",
	string("\u2006"), " ",
	string("\u2007"), " ",
	string("\u2008"), " ",
	string("\u2009"), " ",
	string("\u200A"), " ",
	string("\u202F"), " ",

	string("\u2061"), "",
	string("\u2062"), "",
	string("\u2063"), "",
	string("\u2064"), "",

	string("\u3000"), " ",
)

func normalizeSpace(s string) string {
	s = spaceReplacer.Replace(s)
	return s
}
