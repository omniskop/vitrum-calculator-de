package main

import (
	"fmt"
	"strings"

	"github.com/omniskop/vitrum/gui"
)

// Diese Struktur enthält den Zustand des Programms
type Calculator struct {
	history   []fmt.Stringer // Enthält die letzte Aufgabe, welche über dem Ergebnis angezeigt wird
	prevValue Number         // Der zuletzt eingegebene Wert
	value     Number         // Der aktuelle Wert
	result    Number         // Das Ergebnis der Rechnung
	operation Operation      // Die ausgewählte Operation
	window    *gui.Window    // Das Fenster, in dem das Programm läuft
}

// Instanziiert einen neuen Calculator für das Fenster
func NewCalculator(window *gui.Window) *Calculator {
	calc := &Calculator{
		history:   nil,
		prevValue: Number{},
		value:     Number{},
		window:    window,
		operation: NoOperation,
	}

	// Setzt globale Funktionen die von der Oberfläche aufgerufen werden können
	window.SetVariable("clearAll", calc.ClearAll)
	window.SetVariable("clear", calc.Clear)
	window.SetVariable("appendDigit", calc.AppendDigit)
	window.SetVariable("appendPeriod", calc.AppendPeriod)
	window.SetVariable("backspace", calc.Backspace)
	window.SetVariable("operation", calc.Operation)
	window.SetVariable("solve", calc.Solve)
	window.SetVariable("toggleSign", calc.ToggleSign)
	calc.applyState()

	return calc
}

// Leer die aktuelle Zahl
func (c *Calculator) Clear() {
	fmt.Println("Clear")
	c.value = Number{}
	c.result = Number{}
	c.applyState()
}

// Setzt den Zustand zurück
func (c *Calculator) ClearAll() {
	fmt.Println("ClearAll")
	c.prevValue = Number{}
	c.value = Number{}
	c.history = []fmt.Stringer{}
	c.operation = NoOperation
	c.result = Number{}
	c.applyState()
}

// Fügt eine Ziffer zur aktuellen Zahl hinzu
func (c *Calculator) AppendDigit(digit byte) {
	fmt.Println("Append", digit)
	if !c.value.IsSet() && !c.prevValue.IsSet() {
		// only clean up when no numbers are set at all
		c.history = nil
		c.result = Number{}
	}
	c.value.Append(digit)
	c.applyState()
}

// Hängt einen Punkt an die aktuelle Zahl an
func (c *Calculator) AppendPeriod() {
	fmt.Println("Append period")
	c.value.Append(period)
	c.applyState()
}

// Löscht die letzte Ziffer der aktuellen Zahl
func (c *Calculator) Backspace() {
	fmt.Println("Backspace")
	if !c.value.IsSet() && c.result.IsSet() {
		c.value = c.result
		c.result = Number{}
		c.history = nil
	}
	c.value.Pop()
	c.applyState()
}

// Aktiviert einen bestimmten Operator
func (c *Calculator) Operation(opTxt string) {
	fmt.Println("Operation", opTxt)
	op := ParseOperation(opTxt)
	if c.operation != NoOperation {
		// Wenn bereits eine Operation ausgewählt war, wird diese erst einmal ausgeführt
		c.Solve()
		c.prevValue = c.result // Das Ergebnis davon wird für die nächste übernommen
	} else {
		if c.value.IsSet() {
			// Wenn schon eine Zahl eingegeben wurde wird diese genommen ...
			c.prevValue = c.value
		} else {
			// ... wenn nicht, wird das Ergebnis der letzten Aufgabe genommen
			c.prevValue = c.result
			// Die Spuren der letzten Aufgabe werden entfernt
			c.history = nil
			c.result = Number{}
		}
		c.value = Number{} // Die aktuelle Zahl wird auf 0 gesetzt
	}
	c.operation = op // Die operation wird gespeichert
	c.applyState()
}

// Löst die aktuelle Operation aus
func (c *Calculator) Solve() {
	if c.operation == NoOperation {
		return // Wenn noch gar keine Operation ausgewählt wurde wird nichts gemacht
	}
	fmt.Println("Solve", c.prevValue, c.operation, c.value)

	result := c.operation.Perform(c.prevValue, c.value)
	c.history = []fmt.Stringer{c.prevValue, c.operation, c.value}
	c.operation = NoOperation
	c.prevValue = Number{}
	c.value = Number{}
	c.result = result
	c.applyState()
}

// Kehrt das Vorzeichen der aktuellen Zahl um
func (c *Calculator) ToggleSign() {
	fmt.Println("Toggle sign")
	if !c.value.IsSet() {
		if c.result.IsSet() {
			c.value = c.result
		} else {
			c.value = ZeroNumber()
		}
	}
	c.value.ToggleSign()
	c.applyState()
}

// Erstellt den fertigen Text, der ganz oben auf dem Display angezeigt werden soll.
func (c *Calculator) History() string {
	var out []string

	if c.value.IsSet() || c.prevValue.IsSet() {
		// Die aktuelle Aufgabe
		if c.prevValue.IsSet() {
			out = append(out, c.prevValue.String())
		}
		if c.operation != NoOperation {
			out = append(out, c.operation.String())
		}
		if c.value.IsSet() {
			out = append(out, c.value.String())
		}
	} else {
		// Die letzte Aufgabe
		for _, h := range c.history {
			out = append(out, h.String())
		}
		out = append(out, "=", c.result.String())
	}

	return strings.Join(out, " ")
}

// Überträgt den Zustand auf die Oberfläche
func (c *Calculator) applyState() {
	if c.value.IsSet() {
		// Ein Wert wird eingegeben
		c.window.SetVariable("value", c.value.String())
		c.window.SetVariable("binaryValue", c.value.BinaryString())
		c.window.SetVariable("hexValue", c.value.HexString())
	} else {
		// Das letzte Ergebnis wird angezeigt
		c.window.SetVariable("value", c.result.String())
		c.window.SetVariable("binaryValue", c.result.BinaryString())
		c.window.SetVariable("hexValue", c.result.HexString())
	}
	c.window.SetVariable("history", c.History())
}
