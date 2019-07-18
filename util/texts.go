package util

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/fatih/camelcase"
	"github.com/jinzhu/inflection"
)

const (
	// Id is camelCased Id
	Id = "Id"
	// ID is golang ID
	ID = "ID"
	// Rel if suffix for Relation
	Rel = "Rel"
)

// Singular makes singular of plural english word
func Singular(s string) string {
	return inflection.Singular(s)
}

// IsUpper check rune for upper case
func IsUpper(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

// IsLower check rune for lower case
func IsLower(c byte) bool {
	return c >= 'a' && c <= 'z'
}

// ToUpper converts rune to upper
func ToUpper(c byte) byte {
	return c - 32
}

// ToLower converts rune to lower
func ToLower(c byte) byte {
	return c + 32
}

// CamelCased converts string to camelCase
// from github.com/go-pg/pg/internal
func CamelCased(s string) string {
	r := make([]byte, 0, len(s))
	upperNext := true
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' {
			upperNext = true
			continue
		}
		if upperNext {
			if IsLower(c) {
				c = ToUpper(c)
			}
			upperNext = false
		}
		r = append(r, c)
	}
	return string(r)
}

func separated(s string, delim byte) string {
	r := make([]byte, 0, len(s)+5)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if IsUpper(c) {
			if i > 0 && i+1 < len(s) && (IsLower(s[i-1]) || IsLower(s[i+1])) {
				r = append(r, delim, ToLower(c))
			} else {
				r = append(r, ToLower(c))
			}
		} else {
			r = append(r, c)
		}
	}
	return string(r)
}

// Underscore converts string to under_scored
// from github.com/go-pg/pg/internal
func Underscore(s string) string {
	return separated(s, '_')
}

// Dash converts string to dash-separated
func Dash(s string) string {
	return strings.ReplaceAll(separated(s, '-'), "_", "-")
}

// Sanitize makes string suitable for golang var, const, field, type name
func Sanitize(s string) string {
	rgxp := regexp.MustCompile(`[^a-zA-Z\d\-_]`)
	sanitized := strings.Replace(rgxp.ReplaceAllString(s, ""), "-", "_", -1)

	if len(sanitized) != 0 && ((sanitized[0] >= '0' && sanitized[0] <= '9') || sanitized[0] == '_') {
		sanitized = "T" + sanitized
	}

	return sanitized
}

// PackageName gets string usable as package name
func PackageName(s string) string {
	return strings.ToLower(Sanitize(s))
}

// EntityName gets string usable as struct name
func EntityName(s string) string {
	splitted := camelcase.Split(CamelCased(Sanitize(s)))

	ln := len(splitted) - 1
	for i := ln; i >= 0; i-- {
		split := splitted[i]
		singular := Singular(split)
		if strings.ToLower(singular) != strings.ToLower(split) {
			splitted[i] = strings.Title(singular)
			break
		}
	}

	return strings.Join(splitted, "")
}

// ColumnName gets string usable as struct field name
func ColumnName(s string) string {
	camelCased := ReplaceSuffix(CamelCased(Sanitize(s)), Id, ID)

	return strings.Title(camelCased)
}

// HasUpper checks if string contains upper case
func HasUpper(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if IsUpper(c) {
			return true
		}
	}
	return false
}

// ReplaceSuffix replaces substring on the end of string
func ReplaceSuffix(in, suffix, replace string) string {
	if strings.HasSuffix(in, suffix) {
		in = in[:len(in)-len(suffix)] + replace
	}
	return in
}

// LowerFirst lowers the first letter
func LowerFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[n:]
}
