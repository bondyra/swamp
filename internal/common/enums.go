package common

import "fmt"

type Operator int

// for language parser
func (op *Operator) Capture(values []string) error {
	res, err := NewOperator(values[0])
	if err != nil {
		return fmt.Errorf("Operator.Capture: %w", err)
	}
	*op = res
	return nil
}

const (
	EqualsTo Operator = iota
	NotEqualsTo
)

func NewOperator(op string) (Operator, error) {
	switch op {
	case "eq":
		return EqualsTo, nil
	case "neq":
		return NotEqualsTo, nil
	default:
		return EqualsTo, fmt.Errorf("unknown operator: %s", op)
	}
}

func (o Operator) String() string {
	switch o {
	case EqualsTo:
		return "eq"
	case NotEqualsTo:
		return "neq"
	default:
		panic("unknown operator")
	}
}
