package abstracttapestack

import TM "github.com/Iijil1/Bruteforce-CTL/turingmachines"

type infZeroes struct {
	list              *stackList
	segmentExtensions map[segment]segment
	orExtensions      map[segment]segment
}

func (iz *infZeroes) pop(tail *stack) ([]TM.Symbol, []AbstractStack) {
	return []TM.Symbol{iz.list.symbolsInUse.Zero()}, []AbstractStack{tail}
}
func (iz *infZeroes) extend(value segment) segment {
	if extension, exists := iz.segmentExtensions[value]; exists {
		return extension
	} else {
		newNumber := segment(iz.list.newSequence(value, iz))
		iz.segmentExtensions[value] = newNumber
		return newNumber
	}
}
func (iz *infZeroes) repeat() segment {
	return iz
}
func (iz *infZeroes) unrepeat() segment {
	return iz
}
func (iz *infZeroes) setupRepeater(*stack, *stack) {}
func (iz *infZeroes) selfString(bool) string {
	return "(" + iz.list.symbolsInUse.Zero().String() + ")@"
}
func (iz *infZeroes) registerRepeater(entry segment) {
	iz.registerOr(entry, entry)
}
func (*infZeroes) numRepeaters() int {
	return 0
}
func (iz *infZeroes) or(value2 segment) segment {
	if extension, exists := iz.orExtensions[value2]; exists {
		return extension
	} else {
		return iz.list.newOrSegment(iz, value2)
	}
}
func (iz *infZeroes) registerOr(value2 segment, extension segment) {
	iz.orExtensions[value2] = extension
	iz.orExtensions[extension] = extension
}
func (*infZeroes) orBranches() int {
	return 1
}
func (iz *infZeroes) neatString(reverse bool, _ *stack) string {
	return iz.selfString(reverse)
}
