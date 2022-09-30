package turingmachines

import (
	"fmt"

	"github.com/rgeoghegan/tabulate"
)

type Direction bool

const (
	LEFT  Direction = false
	RIGHT Direction = true
)

func (d Direction) String() string {
	if d == RIGHT {
		return "R"
	}
	return "L"
}

type StateCollection interface {
	NumStates() int
	GetState(int) State
	StartState() State
}

type State interface {
	IsHalt() bool
	IsLoop() bool
	IsUndefined() bool
	String() string
	GetInt() int
}
type SymbolCollection interface {
	NumSymbols() int
	GetSymbol(int) Symbol
	Zero() Symbol
}
type Symbol interface {
	String() string
	GetInt() int
}

type TuringMachine interface {
	StateCollection
	SymbolCollection
	Transition(State, Symbol) (Symbol, Direction, State)
	NameString() string
}

//----------- functions for printing ------------------

func GetTMTable(tm TuringMachine) string {
	var table [][]string
	for i := 0; i < tm.NumStates(); i += 1 {
		newRow := []string{tm.GetState(i).String()}
		for j := 0; j < tm.NumSymbols(); j += 1 {
			symbol, direction, state := tm.Transition(tm.GetState(i), tm.GetSymbol(j))

			if state.IsUndefined() {
				newRow = append(newRow, "---")
			} else {
				newRow = append(newRow, fmt.Sprintf("%v%v%v", symbol, direction, state))
			}
		}
		table = append(table, newRow)

	}
	header := []string{""}
	for j := 0; j < tm.NumSymbols(); j += 1 {
		header = append(header, tm.GetSymbol(j).String())
	}
	layout := &tabulate.Layout{Headers: header, Format: tabulate.SimpleFormat}
	asText, _ := tabulate.Tabulate(
		table, layout,
	)
	return asText
}

func GetStandardTMFormat(tm TuringMachine) string {
	tmString := ""
	for i := 0; i < tm.NumStates(); i += 1 {
		for j := 0; j < tm.NumSymbols(); j += 1 {
			symbol, direction, state := tm.Transition(tm.GetState(i), tm.GetSymbol(j))

			if state.IsUndefined() {
				tmString += "---"
			} else {
				tmString += fmt.Sprintf("%v%v%v", symbol, direction, state)
			}
		}
		tmString += "_"
	}
	return tmString[:len(tmString)-1]
}
