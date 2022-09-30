package turingmachines

import "os"

const DBFILEPATH = "all_5_states_undecided_machines_with_global_header"
const DBENTRIES = 88664064

func GetTMFromDBIndex(dbFilePath string, index uint32) TuringMachine {
	if index < 0 || index >= DBENTRIES {
		return nil
	}
	db, err := os.Open(dbFilePath)
	defer db.Close()
	if err != nil {
		return nil
	}
	bytes := make([]byte, 30)
	db.Seek(int64((index+1)*30), 0)
	l, err := db.Read(bytes)
	if l != 30 || err != nil {
		return nil
	}
	return get5x2MachineFrom30Bytes(bytes)
}

func get5x2MachineFrom30Bytes(bytes []byte) TuringMachine {
	if len(bytes) != 30 {
		return nil
	}
	tm := &simpleTuringMachine{
		transitions: map[transitionInput]transitionOutput{},
		numStates:   5,
		numSymbols:  2,
	}
	for state := 0; state < 5; state += 1 {
		for symbol := 0; symbol < 2; symbol += 1 {
			input := transitionInput{simpleTMState{value: string(byte(state) + 'A')}, simpleTMSymbol{value: string(byte(symbol) + '0')}}
			output := transitionOutput{}
			transitionSymbol := bytes[state*6+symbol*3]
			if transitionSymbol >= 0 && transitionSymbol <= 1 {
				output.symbol = simpleTMSymbol{value: string(transitionSymbol + '0')}
			} else {
				continue
			}
			transitionState := bytes[state*6+symbol*3+2]
			if transitionState >= 1 && transitionState <= 5 {
				output.state = simpleTMState{value: string(transitionState - 1 + 'A')}
			} else {
				continue
			}
			if bytes[state*6+symbol*3+1] == 0 {
				output.direction = RIGHT
			} else if bytes[state*6+symbol*3+1] == 1 {
				output.direction = LEFT
			} else {
				continue
			}
			tm.transitions[input] = output
		}
	}
	return tm
}
