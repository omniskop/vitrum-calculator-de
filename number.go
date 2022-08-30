package main

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/ericlagergren/decimal"
)

const period = byte(255)

// Unlimitierte Präzision
var numberContext = decimal.ContextUnlimited

// Eine Nummer repräsentiert eine Zahl im Dezimalformat. Sie enthält jede Ziffer einzeln.
type Number struct {
	digits   []byte
	negative bool
}

func ZeroNumber() Number {
	return Number{
		digits:   []byte{0},
		negative: false,
	}
}

func BigNumber(f *decimal.Big) Number {
	str := fmt.Sprintf("%f", f)
	var n = Number{
		negative: f.Sign() == -1,
	}
	for _, r := range str {
		if r >= '0' && r <= '9' {
			n.digits = append(n.digits, byte(r-'0'))
		}
		if r == '.' {
			n.digits = append(n.digits, period)
		}
	}
	if n.hasPeriod() {
		// remove trailing zeros on floating point numbers
		n.digits = bytes.TrimRight(n.digits, string([]byte{0}))
		if n.digits[len(n.digits)-1] == period {
			n.digits = n.digits[:len(n.digits)-1] // if there is nothing left after the period, remove it as well
		}
	}
	return n
}

func Add(a, b Number) Number {
	return BigNumber(decimal.WithContext(numberContext).Add(a.Big(), b.Big()))
}

func Subtract(a, b Number) Number {
	return BigNumber(decimal.WithContext(numberContext).Sub(a.Big(), b.Big()))
}

func Multiply(a, b Number) Number {
	return BigNumber(decimal.WithContext(numberContext).Mul(a.Big(), b.Big()))
}

func Div(a, b Number) Number {
	// we won't perform divisions with unlimited precision as repeating decimals will result in NaN
	return BigNumber(decimal.WithContext(decimal.Context128).Quo(a.Big(), b.Big()))
}

func (n *Number) Append(digit byte) {
	if digit == period && n.hasPeriod() {
		return
	}
	n.digits = append(n.digits, digit)
}

func (n *Number) Pop() {
	if len(n.digits) == 0 {
		return
	}
	n.digits = n.digits[:len(n.digits)-1]
}

func (n *Number) ToggleSign() {
	n.negative = !n.negative
}

func (n *Number) IsZero() bool {
	if len(n.digits) == 0 {
		return true
	}
	for _, d := range n.digits {
		if d != 0 && d != period {
			return false
		}
	}
	return true
}

func (n *Number) hasPeriod() bool {
	for _, d := range n.digits {
		if d == period {
			return true
		}
	}
	return false
}

func (n *Number) DecimalPlaces() int {
	if !n.hasPeriod() {
		return 0
	}
	return len(n.digits) - bytes.Index(n.digits, []byte{period}) - 1
}

func (n *Number) IsSet() bool {
	return len(n.digits) > 0
}

func (n *Number) Big() *decimal.Big {
	f, _ := decimal.WithContext(numberContext).SetString(n.String())
	return f
}

func (n *Number) Float() float64 {
	var out float64
	var shiftCounter int
	for i, d := range n.digits {
		if d == period {
			shiftCounter = len(n.digits) - i - 1
			continue
		}
		out = out*10 + float64(d)
	}
	out *= math.Pow10(-shiftCounter)
	if n.negative {
		return -out
	}
	return out
}

func (n Number) String() string {
	var out strings.Builder
	for _, d := range n.digits {
		if d == period {
			if out.Len() == 0 {
				out.WriteRune('0') // add a zero if it starts with a period
			}
			out.WriteRune('.')
			continue
		}
		if d == 0 && out.Len() == 0 {
			continue // skip leading zeros
		}
		out.WriteRune('0' + rune(d))
	}
	if n.negative {
		return "-" + out.String()
	}
	if out.Len() == 0 {
		return "0"
	}
	return out.String()
}

func (n *Number) BinaryString() string {
	if n.hasPeriod() || n.negative {
		return ""
	}
	return strconv.FormatInt(int64(n.Float()), 2)
}

func (n *Number) HexString() string {
	if n.hasPeriod() || n.negative {
		return ""
	}
	return "0x" + strconv.FormatInt(int64(n.Float()), 16)
}
