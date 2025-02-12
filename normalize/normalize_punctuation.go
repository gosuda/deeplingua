package normalize

import "strings"

var punctuationReplacer = strings.NewReplacer(
	string("\u2010"), "-",
	string("\u2011"), "-",
	string("\u2012"), "-",
	string("\u2013"), "-",
	string("\u2014"), "-",
	string("\u2015"), "-",
	string("\u2212"), "-",
	string("\u2018"), "'",
	string("\u2019"), "'",
	string("\u201A"), "'",
	string("\u201B"), "'",
	string("\u201C"), "\"",
	string("\u201D"), "\"",
	string("\u201E"), "\"",
	string("\u201F"), "\"",
	string("\u2024"), ".",
	string("\u2025"), "..",
	string("\u2026"), "...",
	string("\u203C"), "!!",
	string("\u203E"), "-",
	string("\u2043"), "-",
	string("\u2044"), "/",
	string("\u2047"), "??",
	string("\u2053"), "~",
	string("\u301C"), "~",
	string("\u301D"), "\"",
	string("\u301E"), "\"",
)

func normalizePunctuation(s string) string {
	s = punctuationReplacer.Replace(s)
	return s
}
