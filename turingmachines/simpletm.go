package turingmachines

import "strings"

type undefinedState struct{}

func (undefinedState) IsHalt() bool {
	return true
}
func (undefinedState) IsLoop() bool {
	return false
}
func (undefinedState) IsUndefined() bool {
	return true
}
func (undefinedState) String() string {
	return "-"
}
func (undefinedState) GetInt() int {
	return -3
}

type simpleHaltState struct{}

func (simpleHaltState) IsHalt() bool {
	return true
}
func (simpleHaltState) IsLoop() bool {
	return false
}
func (simpleHaltState) IsUndefined() bool {
	return false
}
func (simpleHaltState) String() string {
	return "Z"
}
func (simpleHaltState) GetInt() int {
	return -1
}

type simpleLoopState struct{}

func (simpleLoopState) IsHalt() bool {
	return false
}
func (simpleLoopState) IsLoop() bool {
	return true
}
func (simpleLoopState) IsUndefined() bool {
	return false
}
func (simpleLoopState) String() string {
	return "@"
}
func (simpleLoopState) GetInt() int {
	return -2
}

type simpleTMState struct {
	value string
}

func (simpleTMState) IsHalt() bool {
	return false
}
func (simpleTMState) IsLoop() bool {
	return false
}
func (simpleTMState) IsUndefined() bool {
	return false
}
func (s simpleTMState) String() string {
	return s.value
}
func (s simpleTMState) GetInt() int {
	return int(s.value[0] - 'A')
}

type simpleTMSymbol struct {
	value string
}

func (s simpleTMSymbol) String() string {
	return s.value
}
func (s simpleTMSymbol) GetInt() int {
	return int(s.value[0] - '0')
}

type transitionOutput struct {
	symbol    Symbol
	direction Direction
	state     State
}

type transitionInput struct {
	state  State
	symbol Symbol
}

type simpleTuringMachine struct {
	transitions map[transitionInput]transitionOutput
	numStates   int
	numSymbols  int
}

func (tm *simpleTuringMachine) NumStates() int {
	return tm.numStates
}
func (tm *simpleTuringMachine) GetState(i int) State {
	return simpleTMState{value: string('A' + byte(i))}
}
func (tm *simpleTuringMachine) StartState() State {
	return simpleTMState{value: string('A')}
}
func (tm *simpleTuringMachine) NumSymbols() int {
	return tm.numSymbols
}
func (tm *simpleTuringMachine) GetSymbol(i int) Symbol {
	return simpleTMSymbol{value: string('0' + byte(i))}
}

func (tm *simpleTuringMachine) Transition(state State, symbol Symbol) (Symbol, Direction, State) {
	input := transitionInput{state, symbol}
	if output, exists := tm.transitions[input]; exists {
		return output.symbol, output.direction, output.state

	}
	return tm.Zero(), LEFT, undefinedState{}
}

func (tm *simpleTuringMachine) Zero() Symbol {
	return simpleTMSymbol{value: "0"}
}

func (tm *simpleTuringMachine) NameString() string {
	return GetStandardTMFormat(tm)
}

func GetTuringMachineFromString(tmString string) TuringMachine {
	stateStrings := strings.Split(tmString, "_")
	numStates := len(stateStrings)
	if !(numStates > 0) {
		return &simpleTuringMachine{transitions: map[transitionInput]transitionOutput{}, numStates: 0, numSymbols: 0}
	}
	numSymbols := len(stateStrings[0]) / 3
	for _, stateString := range stateStrings {
		if len(stateString) != 3*numSymbols {
			return &simpleTuringMachine{transitions: map[transitionInput]transitionOutput{}, numStates: 0, numSymbols: 0}
		}
	}
	tm := &simpleTuringMachine{
		transitions: map[transitionInput]transitionOutput{},
		numStates:   numStates,
		numSymbols:  numSymbols,
	}
	for state := 0; state < numStates; state += 1 {
		for symbol := 0; symbol < numSymbols; symbol += 1 {
			input := transitionInput{simpleTMState{value: string(byte(state) + 'A')}, simpleTMSymbol{value: string(byte(symbol) + '0')}}
			output := transitionOutput{}
			transitionSymbol := stateStrings[state][symbol*3]
			if transitionSymbol >= '0' && transitionSymbol < '0'+byte(numSymbols) {
				output.symbol = simpleTMSymbol{value: string(transitionSymbol)}
			} else {
				continue
			}
			transitionState := stateStrings[state][symbol*3+2]
			if transitionState >= 'A' && transitionState < 'A'+byte(numStates) {
				output.state = simpleTMState{value: string(transitionState)}
			} else {
				continue
			}
			if stateStrings[state][symbol*3+1] == 'R' {
				output.direction = RIGHT
			} else if stateStrings[state][symbol*3+1] == 'L' {
				output.direction = LEFT
			} else {
				continue
			}
			tm.transitions[input] = output
		}
	}
	return tm
}
