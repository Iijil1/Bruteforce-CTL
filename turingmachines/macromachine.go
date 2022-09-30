package turingmachines

import (
	"fmt"
	"math"
)

type macroMachine struct {
	internalTM  TuringMachine
	blockSize   int
	backSymbols int
	offset      int
	transitions map[transitionInput]transitionOutput
}

type macroSymbol struct {
	tm       *macroMachine
	encoding int
}

func (ms macroSymbol) String() string {
	blockString := ""
	if ms.tm.blockSize > 1 {
		blockString = "]"
	}
	restEncoding := ms.encoding
	for i := 0; i < ms.tm.blockSize; i += 1 {
		blockString = ms.tm.internalTM.GetSymbol(restEncoding%ms.tm.internalTM.NumSymbols()).String() + blockString
		restEncoding = restEncoding / ms.tm.internalTM.NumSymbols()
	}
	if ms.tm.blockSize > 1 {
		blockString = "[" + blockString
	}
	return blockString
}
func (ms macroSymbol) GetInt() int {
	return ms.encoding
}

type macroSpecial struct {
	tm                  *macroMachine
	backSymbolsEncoding int
	internalState       State
}

func (ms macroSpecial) IsHalt() bool {
	return ms.internalState.IsHalt()
}
func (ms macroSpecial) IsLoop() bool {
	return ms.internalState.IsLoop()
}
func (ms macroSpecial) IsUndefined() bool {
	return ms.internalState.IsUndefined()
}
func (ms macroSpecial) GetInt() int {
	return ms.internalState.GetInt()
}

func (msb macroSpecial) String() string {
	blockString := ""
	if msb.tm.backSymbols > 0 {
		blockString = "]"
	}
	restEncoding := msb.backSymbolsEncoding
	for i := 0; i < msb.tm.backSymbols; i += 1 {
		blockString = msb.tm.internalTM.GetSymbol(restEncoding%msb.tm.internalTM.NumSymbols()).String() + blockString
		restEncoding = restEncoding / msb.tm.internalTM.NumSymbols()
	}
	if msb.tm.backSymbols > 0 {
		blockString = "[" + blockString
	}
	stateString := msb.internalState.String()
	if restEncoding == 0 {
		return blockString + stateString + ">"
	} else {
		return "<" + stateString + blockString
	}
}

type macroStart struct {
}

func (macroStart) IsHalt() bool {
	return false
}
func (macroStart) IsLoop() bool {
	return false
}
func (macroStart) IsUndefined() bool {
	return false
}
func (macroStart) GetInt() int {
	return -3
}

func (macroStart) String() string {
	return "START"
}

type macroState struct {
	tm       *macroMachine
	encoding int
}

func (macroState) IsHalt() bool {
	return false
}
func (macroState) IsLoop() bool {
	return false
}
func (macroState) IsUndefined() bool {
	return false
}
func (ms macroState) GetInt() int {
	return ms.encoding
}

func (msb macroState) String() string {
	blockString := ""
	if msb.tm.backSymbols > 0 {
		blockString = "]"
	}
	restEncoding := msb.encoding
	for i := 0; i < msb.tm.backSymbols; i += 1 {
		blockString = msb.tm.internalTM.GetSymbol(restEncoding%msb.tm.internalTM.NumSymbols()).String() + blockString
		restEncoding = restEncoding / msb.tm.internalTM.NumSymbols()
	}
	if msb.tm.backSymbols > 0 {
		blockString = "[" + blockString
	}
	stateString := msb.tm.internalTM.GetState(restEncoding / 2).String()
	if restEncoding%2 == 0 {
		return blockString + stateString + ">"
	} else {
		return "<" + stateString + blockString
	}
}

func (mm *macroMachine) NumStates() int {
	return 2 * mm.internalTM.NumStates() * int(math.Pow(float64(mm.internalTM.NumSymbols()), float64(mm.backSymbols)))
}

func (mm *macroMachine) NumSymbols() int {
	return int(math.Pow(float64(mm.internalTM.NumSymbols()), float64(mm.blockSize)))
}
func (mm *macroMachine) Transition(state State, symbol Symbol) (Symbol, Direction, State) {
	input := transitionInput{state, symbol}
	if transition, exists := mm.transitions[input]; exists {
		return transition.symbol, transition.direction, transition.state
	}
	output := transitionOutput{mm.Zero(), LEFT, undefinedState{}}
	if macroState, isMacroState := state.(macroState); !isMacroState || macroState.tm != mm {
		return output.symbol, output.direction, output.state
	} else if macroSymbol, isMacroSymbol := symbol.(macroSymbol); !isMacroSymbol || macroSymbol.tm != mm {
		return output.symbol, output.direction, output.state
	} else {
		blockTape := make([]Symbol, mm.blockSize)
		restSymbol := macroSymbol.encoding
		for i := mm.blockSize - 1; i >= 0; i -= 1 {
			blockTape[i] = mm.internalTM.GetSymbol(restSymbol % mm.internalTM.NumSymbols())
			restSymbol = restSymbol / mm.internalTM.NumSymbols()
		}
		backTape := make([]Symbol, mm.backSymbols)
		restSymbol = macroState.encoding
		for i := mm.backSymbols - 1; i >= 0; i -= 1 {
			backTape[i] = mm.internalTM.GetSymbol(restSymbol % mm.internalTM.NumSymbols())
			restSymbol = restSymbol / mm.internalTM.NumSymbols()
		}
		internalState := mm.internalTM.GetState(restSymbol / 2)

		var tape []Symbol
		var position int
		if restSymbol%2 == 0 {
			tape = append(backTape, blockTape...)
			position = mm.backSymbols
		} else {
			tape = append(blockTape, backTape...)
			position = mm.blockSize - 1
		}
		newInternalState, newPosition := tmRunOnFiniteTape(mm.internalTM, tape, internalState, position)
		output = mm.getOutputFromInternalTape(tape, newInternalState, newPosition)
		mm.transitions[input] = output
		return output.symbol, output.direction, output.state
	}
}

func (mm *macroMachine) getOutputFromInternalTape(tape []Symbol, internalState State, position int) transitionOutput {
	output := transitionOutput{}
	newDirection := RIGHT
	if position < 0 {
		newDirection = LEFT
	}

	newStateEncoding := internalState.GetInt()
	if internalState.IsLoop() || internalState.IsHalt() {
		newStateEncoding = 0
	}
	newStateEncoding = newStateEncoding * 2
	if newDirection == LEFT {
		newStateEncoding += 1
	}
	offset := 0
	if newDirection == RIGHT {
		offset = mm.blockSize
	}
	for i := 0; i < mm.backSymbols; i += 1 {
		newStateEncoding = newStateEncoding * mm.internalTM.NumSymbols()
		newStateEncoding += tape[offset+i].GetInt()
	}

	newSymbolEncoding := 0
	offset = 0
	if newDirection == LEFT {
		offset = mm.backSymbols
	}
	for i := 0; i < mm.blockSize; i += 1 {
		newSymbolEncoding = newSymbolEncoding * mm.internalTM.NumSymbols()
		newSymbolEncoding += tape[offset+i].GetInt()
	}

	output.symbol = mm.GetSymbol(newSymbolEncoding)
	output.direction = newDirection
	if internalState.IsHalt() || internalState.IsLoop() {
		output.state = macroSpecial{
			tm:                  mm,
			backSymbolsEncoding: newStateEncoding,
			internalState:       internalState,
		}
	} else {
		output.state = mm.GetState(newStateEncoding)
	}
	return output
}

func (mm *macroMachine) Zero() Symbol {
	return macroSymbol{
		tm:       mm,
		encoding: 0,
	}
}

func (mm *macroMachine) GetSymbol(i int) Symbol {
	return macroSymbol{
		tm:       mm,
		encoding: i,
	}
}

func (mm *macroMachine) GetState(i int) State {
	return macroState{
		tm:       mm,
		encoding: i,
	}
}
func (mm *macroMachine) StartState() State {
	return macroStart{}
}

func (mm *macroMachine) NameString() string {
	return fmt.Sprintf("MM-%v-%v-%v-", mm.blockSize, mm.backSymbols, mm.offset) + mm.internalTM.NameString()
}

func CreateMacroMachine(tm TuringMachine, blockSize int, backSymbols int, offset int) *macroMachine {
	mm := &macroMachine{
		internalTM:  tm,
		blockSize:   blockSize,
		backSymbols: backSymbols,
		transitions: map[transitionInput]transitionOutput{},
	}
	tape := make([]Symbol, blockSize+backSymbols)
	for i := 0; i < blockSize+backSymbols; i += 1 {
		tape[i] = tm.Zero()
	}
	newInternalState, newPosition := tmRunOnFiniteTape(tm, tape, tm.StartState(), backSymbols+offset)

	output := mm.getOutputFromInternalTape(tape, newInternalState, newPosition)
	mm.transitions[transitionInput{state: mm.StartState(), symbol: mm.Zero()}] = output
	return mm
}
