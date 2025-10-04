package interpreter

import (
	"errors"
	"slices"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

/*
func isStopValueRune(r rune) bool {
	_, ok := tokenMap[r]
	if !ok {
		return false
	}
	return true
}
*/

func DecodeTokenTypeInString(inptStr string) (TokenType, error, int) {
	tokenStr := ""
	totalSize := 0
	digits := false
	for len(inptStr) > 0 {
		r, size := utf8.DecodeRuneInString(inptStr)
		if r == utf8.RuneError {
			if size == 1 {
				return 0, errors.New("error: invalid UTF-8 encoding"), 0
			}
		}

		if slices.Contains(digitRunes, r) {
			if !digits && totalSize == 0 {
				digits = true
			} else if !digits {
				if tokenStr == "-" {
					digits = true
				} else {
					break
				}
			}
		} else {
			if digits {
				break
			} else if totalSize > 0 && r == '-' {
				break
			} else if tokenStr == ")" {
				break
			}
		}

		totalSize += size
		tokenStr += string(r)
		inptStr = inptStr[size:]
	}

	op, ok := tokenMap[tokenStr]
	if !ok && !digits {
		return 0, errors.New("error: invalid operator"), totalSize
	} else {
		return op, nil, totalSize
	}
}

func isValueToken(t token) bool {
	return t.Type == Value || t.Type == CloseParen
}

func tokenize(inptStr string) ([]token, error) {
	inptStr = strings.Map(func(r rune) rune {
		if !unicode.IsSpace(r) {
			return r
		}
		return rune(-1)
	}, inptStr)

	var tokens []token
	for len(inptStr) > 0 {
		thisType, err, size := DecodeTokenTypeInString(inptStr)
		if err != nil {
			return nil, err
		}
		/*
			if (len(tokens) == 0 || slices.Contains(operators, tokens[len(tokens)-1].Type)) && TokenType(thisType) == Subtract {
				newType, err, newSize := DecodeTokenTypeInString(inptStr[size:])
				if err != nil {
					return nil, err
				}
				if newType != 0 {
					return nil, errors.New("error: invalid token sequence")
				}
				size += newSize
				thisType = 0
			}
		*/
		if thisType == 0 {
			var lastToken token
			if len(tokens) > 0 {
				lastToken = tokens[len(tokens)-1]
			} else {
				lastToken = token{Type: None}
			}

			if lastToken.Type == CloseParen {
				tokens = append(tokens, token{Type: Multiply})
			}
			valStr := inptStr[:size]

			value, err := strconv.ParseFloat(valStr, 64)
			if err != nil {
				return nil, err
			}
			if lastToken.Type == Value {
				if value < 0 {
					value *= -1
					tokens = append(tokens, token{Type: Subtract})
				} else {
					return nil, errors.New("error: non-negative Value token following Value token")
				}
			}
			tokens = append(tokens, token{Type: Value, Value: value})
		} else {
			if len(tokens) > 0 /* && tokens[len(tokens)-1].Type == Value */ && thisType == OpenParen {
				tokens = append(tokens, token{Type: Multiply})
			}
			tokens = append(tokens, token{Type: thisType})
		}
		inptStr = inptStr[size:]
	}

	return tokens, nil
}
