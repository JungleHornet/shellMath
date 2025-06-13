package interpreter

import (
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"
)

func isStopValueRune(r rune) bool {
	_, ok := tokenMap[r]
	if !ok {
		return false
	}
	return true
}

func isValueToken(t token) bool {
	return t.Type == Value || t.Type == CloseParen
}

func tokenize(inptStrRaw string) ([]token, error) {
	inptStr := strings.TrimSpace(inptStrRaw)
	var runes []rune
	for len(inptStr) > 0 {
		r, size := utf8.DecodeRuneInString(inptStr)
		if r == utf8.RuneError {
			if size == 1 {
				return nil, errors.New("error: invalid UTF-8 encoding")
			}
		}
		inptStr = inptStr[size:]
		runes = append(runes, r)
	}
	var tokens []token

	var valueStr string
	for _, r := range runes {
		if valueStr != "" {
			if !isStopValueRune(r) {
				valueStr += string(r)
				continue
			} else {
				value, err := strconv.ParseFloat(valueStr, 64)
				if err != nil {
					return nil, err
				}
				tokens = append(tokens, token{Type: Value, Value: value})
				valueStr = ""
			}
		}
		if v, ok := tokenMap[r]; ok {
			if v == Subtract {
				if len(tokens) == 0 || !isValueToken(tokens[len(tokens)-1]) {
					valueStr = "-"
					continue
				}
			}
			if v == OpenParen {
				if len(tokens) > 0 && isValueToken(tokens[len(tokens)-1]) {
					tokens = append(tokens, token{Type: Multiply, Value: 0})
				}
			}
			tokens = append(tokens, token{Type: v, Value: 0})
		} else {
			valueStr = string(r)
		}
	}

	if valueStr != "" {
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token{Type: Value, Value: value})
		valueStr = ""
	}

	return tokens, nil
}
