package interpreter

import (
	"strconv"
	"testing"
)

var tokenizerTestMap = map[string][]token{
	/* 1 */ "-1+-1": {{Value, -1}, {Add, 0}, {Value, -1}},
	/* 2 */ "1+-1": {{Value, 1}, {Add, 0}, {Value, -1}},
	/* 3 */ "-1--1": {{Value, -1}, {Subtract, 0}, {Value, -1}},
	/* 4 */ "-12*123": {{Value, -12}, {Multiply, 0}, {Value, 123}},
	/* 5 */ "-1(13/3)-14+-3": {{Value, -1}, {Multiply, 0}, {OpenParen, 0}, {Value, 13}, {Divide, 0}, {Value, 3}, {CloseParen, 0}, {Subtract, 0}, {Value, 14}, {Add, 0}, {Value, -3}},
	/* 6 */ "15(-1+-1)-12*4": {{Value, 15}, {Multiply, 0}, {OpenParen, 0}, {Value, -1}, {Add, 0}, {Value, -1}, {CloseParen, 0}, {Subtract, 0}, {Value, 12}, {Multiply, 0}, {Value, 4}},
	/* 7 */ "-7(2-4)(16+2)": {{Value, -7}, {Multiply, 0}, {OpenParen, 0}, {Value, 2}, {Subtract, 0}, {Value, 4}, {CloseParen, 0}, {Multiply, 0}, {OpenParen, 0}, {Value, 16}, {Add, 0}, {Value, 2}, {CloseParen, 0}},
}

var tokenizerTests = []string{
	"-1+-1",
	"1+-1",
	"-1--1",
	"-12*123",
	"-1(13/3)-14+-3",
	"15(-1+-1)-12*4",
	"-7(2-4)(16+2)",
}

func tokenSlicesEqual(s1, s2 []token) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, t1 := range s1 {
		t2 := s2[i]
		if t1 != t2 {
			return false
		}
	}
	return true
}

func TestTokenize(t *testing.T) {
	for i, testStr := range tokenizerTests {
		result := tokenizerTestMap[testStr]
		if output, err := tokenize(testStr); err != nil || !tokenSlicesEqual(output, result) {
			t.Log("Failure on test #" + strconv.Itoa(i+1))
			t.Log("\t\tExpected:", result)
			t.Log("\t\tGot:     ", output)
			t.Log("\t\terr:     ", err)
			t.Fail()
		}
	}
}
