package main

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strings"
)

type Change struct {
	URL        string
	BeforeHash string
	AfterHash  string
	Changed    bool
}

var whitespace = regexp.MustCompile(`\s+`)

func normalizePage(body string) string {
	return strings.TrimSpace(whitespace.ReplaceAllString(body, " "))
}

func pageHash(body string) string {
	sum := sha256.Sum256([]byte(normalizePage(body)))
	return hex.EncodeToString(sum[:])
}

func detectChange(url, before, after string) Change {
	b, a := pageHash(before), pageHash(after)
	return Change{URL: url, BeforeHash: b, AfterHash: a, Changed: b != a}
}
