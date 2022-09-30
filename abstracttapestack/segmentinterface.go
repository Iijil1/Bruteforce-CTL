package abstracttapestack

import TM "github.com/Iijil1/Bruteforce-CTL/turingmachines"

type segment interface {
	pop(tail *stack) ([]TM.Symbol, []AbstractStack)
	selfString(reverse bool) string
	neatString(reverse bool, tail *stack) string
	extend(value segment) segment
	repeat() segment
	unrepeat() segment
	setupRepeater(currentStack *stack, baseStack *stack)
	registerRepeater(segment)
	or(segment) segment
	registerOr(segment, segment)
	numRepeaters() int
	orBranches() int
}
