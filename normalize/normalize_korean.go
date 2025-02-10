package normalize

import "strings"

var koreanReplacer = strings.NewReplacer(
	string("\x11\x5F"), "",
	string("\x11\x60"), "",
	string("\x11\x00"), "ㄱ",
	string("\x11\x01"), "ㄲ",
	string("\x11\x02"), "ㄴ",
	string("\x11\x03"), "ㄷ",
	string("\x11\x04"), "ㄸ",
	string("\x11\x05"), "ㄹ",
	string("\x11\x06"), "ㅁ",
	string("\x11\x07"), "ㅂ",
	string("\x11\x08"), "ㅃ",
	string("\x11\x09"), "ㅅ",
	string("\x11\x0A"), "ㅆ",
	string("\x11\x0B"), "ㅇ",
	string("\x11\x0C"), "ㅈ",
	string("\x11\x0D"), "ㅉ",
	string("\x11\x0E"), "ㅊ",
	string("\x11\x0F"), "ㅋ",
	string("\x11\x10"), "ㅌ",
	string("\x11\x11"), "ㅍ",
	string("\x11\x12"), "ㅎ",
	string("\x11\x61"), "ㅏ",
	string("\x11\x62"), "ㅐ",
	string("\x11\x63"), "ㅑ",
	string("\x11\x64"), "ㅒ",
	string("\x11\x65"), "ㅓ",
	string("\x11\x66"), "ㅔ",
	string("\x11\x67"), "ㅕ",
	string("\x11\x68"), "ㅖ",
	string("\x11\x69"), "ㅗ",
	string("\x11\x6A"), "ㅚ",
	string("\x11\x6B"), "ㅙ",
	string("\x11\x6C"), "ㅚ",
	string("\x11\x6D"), "ㅛ",
	string("\x11\x6E"), "ㅜ",
	string("\x11\x6F"), "ㅝ",
	string("\x11\x70"), "ㅞ",
	string("\x11\x71"), "ㅟ",
	string("\x11\x72"), "ㅠ",
	string("\x11\x73"), "ㅡ",
	string("\x11\x74"), "ㅢ",
	string("\x11\x75"), "ㅣ",
)

func normalizeKorean(s string) string {
	s = koreanReplacer.Replace(s)

	return s
}
