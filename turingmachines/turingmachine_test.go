package turingmachines

import (
	"testing"
)

func TestTuringMachineFromString(t *testing.T) {
	tm := GetTuringMachineFromString("1RB1LC_1RC1RB_1RD0LE_1LA1LD_1RZ0LA")
	t.Run("FirstTransition", func(t *testing.T) {
		if symbol, direction, state := tm.Transition(tm.GetState(0), tm.GetSymbol(0)); symbol != tm.GetSymbol(1) || direction != RIGHT || state != tm.GetState(1) {
			t.Fail()
		}
	})
	t.Run("LastTransition", func(t *testing.T) {
		if symbol, direction, state := tm.Transition(tm.GetState(4), tm.GetSymbol(1)); symbol != tm.GetSymbol(0) || direction != LEFT || state != tm.GetState(0) {
			t.Fail()
		}
	})
	t.Run("CentralTransition", func(t *testing.T) {
		if symbol, direction, state := tm.Transition(tm.GetState(2), tm.GetSymbol(1)); symbol != tm.GetSymbol(0) || direction != LEFT || state != tm.GetState(4) {
			t.Fail()
		}
	})
	t.Run("HaltingTranstion", func(t *testing.T) {
		if _, _, state := tm.Transition(tm.GetState(4), tm.GetSymbol(0)); !state.IsHalt() {
			t.Fail()
		}
	})
}

//TODO: Tests for simulation, macromachines, database, ...
