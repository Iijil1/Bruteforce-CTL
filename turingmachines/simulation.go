package turingmachines

func tmStep(tm TuringMachine, tape []Symbol, state State, position int) (State, int) {
	nextSymbol, nextDirection, nextState := tm.Transition(state, tape[position])
	if nextState.IsUndefined() {
		return nextState, position
	}
	tape[position] = nextSymbol
	shift := -1
	if nextDirection == RIGHT {
		shift = 1
	}
	return nextState, position + shift
}

type finiteTMConfig struct {
	TapeIndex int
	State     State
	Position  int
}

func indexForTape(tape []Symbol, base int) int {
	index := 0
	for i := 0; i < len(tape); i += 1 {
		index = index * base
		index += tape[i].GetInt()
	}
	return index
}

func tmRunOnFiniteTape(tm TuringMachine, tape []Symbol, state State, position int) (State, int) {
	seen := map[finiteTMConfig]bool{{TapeIndex: indexForTape(tape, tm.NumSymbols()), State: state, Position: position}: true}
	for i := 0; position >= 0 && position < len(tape) && i < 100; i += 1 {
		state, position = tmStep(tm, tape, state, position)
		if state.IsHalt() || state.IsLoop() {
			return state, position
		}
		newConfig := finiteTMConfig{TapeIndex: indexForTape(tape, tm.NumSymbols()), State: state, Position: position}
		if _, exists := seen[newConfig]; exists {
			return simpleLoopState{}, position
		}
		seen[newConfig] = true
	}
	return state, position
}
