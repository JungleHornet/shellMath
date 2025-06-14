package interpreter

import (
	"math"
	"slices"
)

type tokenType uint32

type token struct {
	Type  tokenType
	Value float64
}

// order matters! lower values will be earlier in order of operations.
const (
	Value tokenType = iota
	Tree            // special tokentype representing an already constructed subtree during AST construction
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

var tokenMap = map[rune]tokenType{
	'(': OpenParen,
	')': CloseParen,
	'^': Exponent,
	'*': Multiply,
	'/': Divide,
	'%': Modulo,
	'÷': IntegerDivision,
	'+': Add,
	'-': Subtract,
}

var operators = []tokenType{
	Exponent,
	Multiply,
	Divide,
	Modulo,
	IntegerDivision,
	Add,
	Subtract,
}

type Operator func(num1, num2 float64) float64

var operatorFuncs = map[tokenType]Operator{
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

var orderOfOperations = [][]tokenType{
	{Exponent},
	{Multiply, Divide, Modulo, IntegerDivision},
	{Add, Subtract},
}

func IsOperatorToken(t tokenType) bool {
	return slices.Contains(operators, t)
}

func IsOperatorRune(r rune) bool {
	tType, ok := tokenMap[r]
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

func Evaluate(equation string) (float64, error) {
	tokens, err := tokenize(equation)
	if err != nil {
		return 0, err
	}

	res, err := buildAST(tokens)
	if err != nil {
		return 0, err
	}

	return res.evaluate(), nil
}
