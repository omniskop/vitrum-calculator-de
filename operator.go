package main

type Operation byte

const (
	NoOperation Operation = iota
	AddOperation
	SubtractOperation
	MultiplyOperation
	DivideOperation
)

func ParseOperation(op string) Operation {
	switch op {
	case "+":
		return AddOperation
	case "-":
		return SubtractOperation
	case "*":
		return MultiplyOperation
	case "/":
		return DivideOperation
	}
	return NoOperation
}

func (o Operation) String() string {
	switch o {
	case AddOperation:
		return "+"
	case SubtractOperation:
		return "-"
	case MultiplyOperation:
		return "*"
	case DivideOperation:
		return "/"
	}
	return ""
}

func (o Operation) Perform(a, b Number) Number {
	switch o {
	case AddOperation:
		return Add(a, b)
	case SubtractOperation:
		return Subtract(a, b)
	case MultiplyOperation:
		return Multiply(a, b)
	case DivideOperation:
		return Div(a, b)
	default:
		return Number{}
	}
}
