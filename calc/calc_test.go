package calc

import (
	"reflect"
	"testing"
)

func TestInfixToRPN(t *testing.T) {
	expr := "3.14*8^2"
	mustbe := []string{"3.14", "8", "2", "^", "*"}
	res, err := InfixToRPN(expr)
	if !reflect.DeepEqual(mustbe, res) || err != nil {
		t.Fatalf(`InfixToRPN test failed, err = %v`, err)
	}
}

func TestCalc(t *testing.T) {
	expr := "3 + 8 ^ 2"
	res, err := Calc(expr)
	if res != 67 || err != nil {
		t.Fatalf(`Calc("3 + 5 * 8") = %f, %v, want 67, nil`, res, err)
	}
}
