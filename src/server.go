package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"unicode"
)

var priorities = map[rune]int{
	'+': 1,
	'-': 1,
	'*': 2,
	'/': 2,
}

func InfixToRPN(expression string) ([]string, error) {
	if len(expression) == 0 {
		return nil, errors.New("Expression is not valid")
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
				return nil, errors.New("Expression is not valid")
			}
			operatorStack = operatorStack[:len(operatorStack)-1]
		} else {
			return nil, errors.New("Expression is not valid")
		}
		i++
	}
	for len(operatorStack) > 0 {
		if operatorStack[len(operatorStack)-1] == '(' {
			return nil, errors.New("Expression is not valid")
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
		if !(strings.Contains("+-/*()", el)) {
			num, err := strconv.ParseFloat(el, 64)
			if err != nil {
				return 0, errors.New("Expression is not valid")
			}
			stack = append(stack, num)
		} else {
			if len(stack) < 2 {
				return 0, errors.New("Expression is not valid")
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
					return 0, errors.New("Expression is not valid")
				}
				stack = append(stack, left/right)
			} else {
				return 0, errors.New("Expression is not valid")
			}
		}
	}
	return stack[len(stack)-1], nil
}

func Calc(expression string) (float64, error) {
	expression = strings.ReplaceAll(expression, " ", "")
	expression = strings.ReplaceAll(expression, ",", ".")
	if len(expression) == 0 {
		return 0, errors.New("Internal server error")
	}
	if strings.Count(expression, "(") != strings.Count(expression, ")") {
		return 0, errors.New("Internal server error")
	}
	if isOpStr(expression[len(expression)-1:]) {
		return 0, errors.New("Internal server error")
	}
	if strings.Count(expression, "+)") != 0 {
		return 0, errors.New("Internal server error")
	}
	if strings.Count(expression, "-)") != 0 {
		return 0, errors.New("Internal server error")
	}
	if strings.Count(expression, "*)") != 0 {
		return 0, errors.New("Internal server error")
	}
	if strings.Count(expression, "/)") != 0 {
		return 0, errors.New("Internal server error")
	}
	rpn_expr, err := InfixToRPN(expression)
	if err != nil {
		return 0, err
	}
	return evalRPN(rpn_expr)
}

func calcHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer w.Write([]byte("\n"))
	if r.Method != http.MethodPost {
		w.WriteHeader(500)
		w.Write([]byte("{\n    \"error\": \"Internal server error\"\n}"))
		return
	}
	var req map[string]string
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("{\n    \"error\": \"Internal server error\"\n}"))
		return
	}
	res, calcerr := Calc(req["expression"])
	if calcerr != nil {
		if calcerr.Error() == "Expression is not valid" {
			w.WriteHeader(422)
			w.Write([]byte("{\n    \"error\": \"Expression is not valid\"\n}"))
			return
		}
		w.WriteHeader(500)
		w.Write([]byte("{\n    \"error\": \"Internal server error\"\n}"))
		return
	}
	w.WriteHeader(200)
	resmap := make(map[string]string)
	resb := strconv.FormatFloat(res, 'f', -1, 64)
	resmap["result"] = resb
	resjson, err2 := json.MarshalIndent(resmap, "", "    ")
	if err2 != nil {
		w.WriteHeader(500)
		w.Write([]byte("{\n    \"error\": \"Internal server error\"\n}"))
		return
	}
	w.Write(resjson)
}

func main() {
	http.HandleFunc("/api/v1/calculate", calcHandler)
	http.ListenAndServe(":8080", nil)
}
