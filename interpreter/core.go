package interpreter

import (
	"fmt"
	"math"
	"slices"
)

type TokenType uint64

type token struct {
	Type  TokenType
	Value float64
}

// order matters! lower values will be earlier in order of operations.

const None TokenType = math.MaxUint64 // special TokenType representing a nil or non-existent token
const (
	Value TokenType = iota
	Tree            // special TokenType representing an already constructed subtree during AST construction
	OpenParen
	CloseParen
	Exponent
	Multiply
	Divide
	Modulo
	IntegerDivision
	Add
	Subtract
)

var tokenMap = map[string]TokenType{
	"(":  OpenParen,
	")":  CloseParen,
	"^":  Exponent,
	"*":  Multiply,
	"/":  Divide,
	"%":  Modulo,
	"//": IntegerDivision,
	"+":  Add,
	"-":  Subtract,
}

var digitRunes = []rune{
	'0',
	'1',
	'2',
	'3',
	'4',
	'5',
	'6',
	'7',
	'8',
	'9',
	'.',
}

var operators = []TokenType{
	Exponent,
	Multiply,
	Divide,
	Modulo,
	IntegerDivision,
	Add,
	Subtract,
}

type Operator func(num1, num2 float64) float64

var operatorFuncs = map[TokenType]Operator{
	Exponent: func(num1, num2 float64) float64 {
		return math.Pow(num1, num2)
	},
	Multiply: func(num1, num2 float64) float64 {
		return num1 * num2
	},
	Divide: func(num1, num2 float64) float64 {
		return num1 / num2
	},
	Modulo: func(num1, num2 float64) float64 {
		return math.Mod(num1, num2)
	},
	IntegerDivision: func(num1, num2 float64) float64 {
		return (num1 - math.Mod(num1, num2)) / num2
	},
	Add: func(num1, num2 float64) float64 {
		return num1 + num2
	},
	Subtract: func(num1, num2 float64) float64 {
		return num1 - num2
	},
}

var orderOfOperations = [][]TokenType{
	{Exponent},
	{Multiply, Divide, Modulo, IntegerDivision},
	{Add, Subtract},
}

func IsOperatorToken(t TokenType) bool {
	return slices.Contains(operators, t)
}

func IsOperatorString(s string) bool {
	tType, ok := tokenMap[s]
	if !ok {
		return false
	}
	return slices.Contains(operators, tType)
}

type AST interface {
	evaluate() float64
}

type ASTBranch struct {
	function Operator
	t1, t2   SubAST
}

type SubAST interface {
	evaluate() float64
}

type ASTValue float64

func (val ASTValue) evaluate() float64 {
	return float64(val)
}

func (ast ASTBranch) evaluate() float64 {
	return ast.function(ast.t1.evaluate(), ast.t2.evaluate())
}

var debug = false

func Evaluate(equation string) (float64, error) {
	if equation == "debug=true" {
		debug = true
		return 0, nil
	} else if equation == "debug=false" {
		debug = false
		return 0, nil
	}
	tokens, err := tokenize(equation)
	if debug {
		fmt.Println(tokens)
	}
	if err != nil {
		return 0, err
	}

	res, err := buildAST(tokens)
	if err != nil {
		return 0, err
	}

	return res.evaluate(), nil
}
