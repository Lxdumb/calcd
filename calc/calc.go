package calc

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"unicode"
)

var priorities = map[rune]int{
	'+': 1,
	'-': 1,
	'*': 2,
	'/': 2,
	'^': 3,
}

func InfixToRPN(expression string) ([]string, error) {
	if len(expression) == 0 {
		return nil, errors.New("expression is not valid")
	}
	outputQueue := []string{}
	operatorStack := []rune{}
	var i int
	for i < len(expression) {
		char := rune(expression[i])
		if unicode.IsDigit(char) || char == '.' {
			num := parseNumber(expression, &i)
			outputQueue = append(outputQueue, num)
			continue
		}
		if isOperator(char) {
			for len(operatorStack) > 0 && operatorStack[len(operatorStack)-1] != '(' && priorities[operatorStack[len(operatorStack)-1]] >= priorities[char] {
				outputQueue = append(outputQueue, string(operatorStack[len(operatorStack)-1]))
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
			operatorStack = append(operatorStack, char)
		} else if char == '(' {
			operatorStack = append(operatorStack, char)
		} else if char == ')' {
			for len(operatorStack) > 0 && operatorStack[len(operatorStack)-1] != '(' {
				outputQueue = append(outputQueue, string(operatorStack[len(operatorStack)-1]))
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
			if len(operatorStack) == 0 {
				return nil, errors.New("expression is not valid")
			}
			operatorStack = operatorStack[:len(operatorStack)-1]
		} else {
			return nil, errors.New("expression is not valid")
		}
		i++
	}
	for len(operatorStack) > 0 {
		if operatorStack[len(operatorStack)-1] == '(' {
			return nil, errors.New("expression is not valid")
		}
		outputQueue = append(outputQueue, string(operatorStack[len(operatorStack)-1]))
		operatorStack = operatorStack[:len(operatorStack)-1]
	}

	return outputQueue, nil
}

func isOperator(c rune) bool {
	_, isExisting := priorities[c]
	return isExisting
}

func isOpStr(s string) bool {
	return s == "+" || s == "-" || s == "*" || s == "/"
}

func parseNumber(expr string, index *int) string {
	start := *index
	for *index < len(expr) && (unicode.IsDigit(rune(expr[*index])) || rune(expr[*index]) == '.') {
		*index++
	}
	return expr[start:*index]
}

func remove(slice []float64, s int) []float64 {
	return append(slice[:s], slice[s+1:]...)
}

func evalRPN(expression []string) (float64, error) {
	var stack []float64
	for _, el := range expression {
		if !(strings.Contains("+-/*()^", el)) {
			num, err := strconv.ParseFloat(el, 64)
			if err != nil {
				return 0, errors.New("expression is not valid")
			}
			stack = append(stack, num)
		} else {
			if len(stack) < 2 {
				return 0, errors.New("expression is not valid")
			}
			right := stack[len(stack)-1]
			stack = remove(stack, len(stack)-1)
			left := stack[len(stack)-1]
			stack = remove(stack, len(stack)-1)
			if el == "+" {
				stack = append(stack, left+right)
			} else if el == "-" {
				stack = append(stack, left-right)
			} else if el == "*" {
				stack = append(stack, left*right)
			} else if el == "/" {
				if right == 0 {
					return 0, errors.New("expression is not valid")
				}
				stack = append(stack, left/right)
			} else if el == "^" {
				stack = append(stack, math.Pow(left, right))
			} else {
				return 0, errors.New("expression is not valid")
			}
		}
	}
	return stack[len(stack)-1], nil
}

func Calc(expression string) (float64, error) {
	expression = strings.ReplaceAll(expression, " ", "")
	expression = strings.ReplaceAll(expression, ",", ".")
	if len(expression) == 0 {
		return 0, errors.New("internal server error")
	}
	if strings.Count(expression, "(") != strings.Count(expression, ")") {
		return 0, errors.New("internal server error")
	}
	if isOpStr(expression[len(expression)-1:]) {
		return 0, errors.New("internal server error")
	}
	if strings.Count(expression, "+)") != 0 {
		return 0, errors.New("internal server error")
	}
	if strings.Count(expression, "-)") != 0 {
		return 0, errors.New("internal server error")
	}
	if strings.Count(expression, "*)") != 0 {
		return 0, errors.New("internal server error")
	}
	if strings.Count(expression, "/)") != 0 {
		return 0, errors.New("internal server error")
	}
	if strings.Count(expression, "^)") != 0 {
		return 0, errors.New("internal server error")
	}
	rpn_expr, err := InfixToRPN(expression)
	if err != nil {
		return 0, err
	}
	return evalRPN(rpn_expr)
}
