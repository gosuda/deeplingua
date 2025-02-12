package normalize

func Normalize(s string) string {
	s = normalizeSpace(s)
	s = normalizePunctuation(s)
	return s
}
