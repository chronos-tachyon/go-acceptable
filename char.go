package acceptable

import (
	"bytes"
	"strings"
	"sync"
)

var gPool = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(make([]byte, 0, 64))
	},
}

func appendToken(out []byte, token string) []byte {
	if token == "*" {
		return append(out, '*')
	}

	if stringMatches(token, isToken) {
		return append(out, token...)
	}

	out = append(out, '"')
	n := uint(len(token))
	for i := uint(0); i < n; i++ {
		ch := token[i]
		if ch == '"' || ch == '\\' {
			out = append(out, '\\')
		}
		out = append(out, ch)
	}
	out = append(out, '"')
	return out
}

func stringMatches(str string, pred func(byte) bool) bool {
	n := uint(len(str))
	ok := true
	for i := uint(0); i < n; i++ {
		ch := str[i]
		if !pred(ch) {
			ok = false
			break
		}
	}
	return ok
}

func consumeSpace(input string) string {
	for input != "" {
		if !isLWS(input[0]) {
			break
		}
		input = input[1:]
	}
	return input
}

func consumeToken(input string) (token string, rest string, ok bool) {
	if input == "" {
		return
	}

	n := uint(len(input))
	i := uint(0)
	for i < n {
		ch := input[i]
		if !isToken(ch) {
			break
		}
		i++
	}

	token = input[:i]
	rest = input[i:]
	ok = (i != 0)
	return
}

func consumeQuoted(input string) (quoted string, rest string, ok bool) {
	if input == "" {
		return
	}

	if input[0] != '"' {
		return consumeToken(input)
	}

	buf := gPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		gPool.Put(buf)
	}()

	n := uint(len(input))
	i := uint(1)
	inEscape := false
	for i < n {
		ch := input[i]
		i++

		switch {
		case isControl(ch):
			return

		case inEscape:
			buf.WriteByte(ch)
			inEscape = false

		case ch == '\\':
			inEscape = true

		case ch == '"':
			quoted = buf.String()
			rest = input[i:]
			ok = true
			return

		default:
			buf.WriteByte(ch)
		}
	}
	return
}

func isLWS(ch byte) bool    { return ch == ' ' || ch == '\t' }
func isQuote(ch byte) bool  { return ch == '"' }
func isComma(ch byte) bool  { return ch == ',' }
func isPeriod(ch byte) bool { return ch == '.' }
func isDigit(ch byte) bool  { return ch >= '0' && ch <= '9' }
func isColon(ch byte) bool  { return ch == ':' }
func isSemi(ch byte) bool   { return ch == ';' }
func isEqual(ch byte) bool  { return ch == '=' }
func isUpper(ch byte) bool  { return ch >= 'A' && ch <= 'Z' }
func isLower(ch byte) bool  { return ch >= 'a' && ch <= 'z' }

func isLetter(ch byte) bool { return isUpper(ch) || isLower(ch) }
func isAlnum(ch byte) bool  { return isDigit(ch) || isLetter(ch) }
func isToken(ch byte) bool  { return isAlnum(ch) || isTokenMisc(ch) }

func isTokenMisc(ch byte) bool {
	const SET = "!#$%&'*+-.^_`|~"
	return strings.IndexByte(SET, ch) >= 0
}

func isControl(ch byte) bool {
	if isLWS(ch) {
		return false
	}
	return ch < 0x20 || ch == 0x7f
}
