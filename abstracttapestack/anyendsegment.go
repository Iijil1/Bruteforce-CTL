package abstracttapestack

import TM "github.com/Iijil1/Bruteforce-CTL/turingmachines"

type infAnys struct {
	list              *stackList
	segmentExtensions map[segment]segment
	orExtensions      map[segment]segment
}

func (ae *infAnys) pop(tail *stack) ([]TM.Symbol, []AbstractStack) {
	numSymbols := ae.list.symbolsInUse.NumSymbols()
	symbols := make([]TM.Symbol, numSymbols)
	tails := make([]AbstractStack, numSymbols)
	for i := 0; i < numSymbols; i += 1 {
		symbols[i] = ae.list.symbolsInUse.GetSymbol(i)
		tails[i] = tail
	}
	return symbols, tails
}
func (ae *infAnys) repeat() segment {
	return ae
}
func (ae *infAnys) unrepeat() segment {
	return ae
}
func (ae *infAnys) setupRepeater(*stack, *stack) {}
func (ae *infAnys) registerRepeater(entry segment) {
	ae.registerOr(entry, entry)
}
func (*infAnys) numRepeaters() int {
	return 1
}
func (ae *infAnys) extend(value segment) segment {
	if extension, exists := ae.segmentExtensions[value]; exists {
		return extension
	} else {
		newNumber := ae.list.newSequence(value, ae)
		ae.segmentExtensions[value] = newNumber
		return newNumber
	}
}
func (ae *infAnys) or(value2 segment) segment {
	if extension, exists := ae.orExtensions[value2]; exists {
		return extension
	} else {
		return ae.list.newOrSegment(ae, value2)
	}
}
func (ae *infAnys) registerOr(value2 segment, extension segment) {
	ae.orExtensions[value2] = extension
	ae.orExtensions[extension] = extension
}
func (*infAnys) orBranches() int {
	return 1
}
func (*infAnys) selfString(bool) string {
	return "(.)@"
}
func (ae *infAnys) neatString(reverse bool, tail *stack) string {
	return ae.selfString(reverse)
}
