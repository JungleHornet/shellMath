package interpreter

import (
	"errors"
	"slices"
)

func countOccurrences[T comparable](slice []T, target T) int {
	count := 0
	for _, value := range slice {
		if value == target {
			count++
		}
	}
	return count
}

func buildParentheses(tokens []token) ([]token, []AST, error) {
	openParentheses := countOccurrences(tokens, token{OpenParen, 0})
	closeParentheses := countOccurrences(tokens, token{CloseParen, 0})
	if openParentheses != closeParentheses {
		return tokens, nil, errors.New("parentheses parsing error: unmatched parentheses")
	}

	var parenPairs [][]int

	lvl := 0
	for i, t := range tokens {
		if t.Type == CloseParen {
			lvl--
		}

		if lvl == 0 {
			if t.Type == OpenParen {
				parenPairs = append(parenPairs, []int{i})
			} else if t.Type == CloseParen {
				last := len(parenPairs) - 1
				parenPairs[last] = append(parenPairs[last], i)
			}
		}

		if t.Type == OpenParen {
			lvl++
		}
	}

	var trees []AST

	for _, p := range parenPairs {
		tree, err := buildAST(tokens[p[0]+1 : p[1]])
		if err != nil {
			return tokens, nil, err
		}

		tokens = append(append(tokens[:p[0]], token{Type: Tree, Value: float64(len(trees))}), tokens[p[1]+1:]...)

		trees = append(trees, tree)
	}

	return tokens, trees, nil
}

func buildAST(tokens []token) (AST, error) {
	var tree AST

	tokens, trees, err := buildParentheses(tokens)
	if err != nil {
		return nil, err
	}

	var operatorIndices []int
	for i, t := range tokens {
		if IsOperatorToken(t.Type) {
			operatorIndices = append(operatorIndices, i)
		}
	}

	if len(tokens) == 0 {
		return ASTValue(0), nil
	}

	if len(operatorIndices) == 0 {
		return ASTValue(tokens[0].Value), nil
	}

	for x := 0; x < len(orderOfOperations); x++ {
		current := orderOfOperations[x]
		for i := 0; i < len(operatorIndices); i++ {
			index := operatorIndices[i]
			this := tokens[index]
			if slices.Contains(current, this.Type) {
				if i == len(operatorIndices)-1 {
					operatorIndices = append(operatorIndices[:i])
				} else {
					operatorIndices = append(operatorIndices[:i], operatorIndices[i+1:]...)
				}

				i--
				var t1, t2 SubAST

				t := tokens[index-1]
				if t.Type == Tree {
					t1 = trees[int(t.Value)]
				} else if t.Type == Value {
					t1 = ASTValue(t.Value)
				} else {
					return nil, errors.New("operator parsing error: operator not acting on valid value")
				}
				t = tokens[index+1]

				t = tokens[index+1]
				if t.Type == Tree {
					t2 = trees[int(t.Value)]
				} else if t.Type == Value {
					t2 = ASTValue(t.Value)
				} else {
					return nil, errors.New("operator parsing error: operator not acting on valid value")
				}
				t = tokens[index+1]

				subtree := ASTBranch{function: operatorFuncs[this.Type], t1: t1, t2: t2}
				tokens = append(append(tokens[:index-1], token{Type: Tree, Value: float64(len(trees))}), tokens[index+2:]...)
				trees = append(trees, subtree)

				for i, v := range operatorIndices {
					if v >= index {
						operatorIndices[i] = v - 2
					}
				}
			}

			if len(tokens) == 1 {
				break
			}
		}
	}

	tree = trees[len(trees)-1]
	return tree, nil
}
