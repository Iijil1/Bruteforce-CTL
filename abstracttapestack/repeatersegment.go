package abstracttapestack

import TM "github.com/Iijil1/Bruteforce-CTL/turingmachines"

type repeater struct {
	list              *stackList
	value             segment
	segmentExtensions map[segment]segment
	orExtensions      map[segment]segment
}

func (r *repeater) pop(tail *stack) ([]TM.Symbol, []AbstractStack) {
	var emptyRepeaterValues []TM.Symbol
	var emptyRepeaterTails []AbstractStack
	if r.list.abstractions.UseStarInsteadOfPlus {
		emptyRepeaterValues, emptyRepeaterTails = tail.Pop()
	} else {
		emptyRepeaterValues, emptyRepeaterTails = r.value.pop(tail)
	}
	fullRepeaterValues, fullRepeaterTails := r.value.pop(tail.pushSegment(r))
	return append(emptyRepeaterValues, fullRepeaterValues...), append(emptyRepeaterTails, fullRepeaterTails...)
}
func (r *repeater) extend(value segment) segment {
	if extension, exists := r.segmentExtensions[value]; exists {
		return extension
	} else {
		newNumber := r.list.newSequence(value, r)
		r.segmentExtensions[value] = newNumber
		return newNumber
	}
}
func (r *repeater) repeat() segment {
	return r
}
func (r *repeater) unrepeat() segment {
	return r.value.unrepeat()
}
func (r *repeater) setupRepeater(currentStack *stack, baseStack *stack) {
	r.value.setupRepeater(currentStack, baseStack)
}
func (r *repeater) registerRepeater(entry segment) {
	r.registerOr(entry, entry)
}
func (r *repeater) numRepeaters() int {
	return r.value.numRepeaters() + 1
}
func (r *repeater) or(value2 segment) segment {
	if extension, exists := r.orExtensions[value2]; exists {
		return extension
	} else {
		return r.list.newOrSegment(r, value2)
	}
}
func (r *repeater) registerOr(value2 segment, extension segment) {
	r.orExtensions[value2] = extension
	r.orExtensions[extension] = extension
}
func (r *repeater) orBranches() int {
	return r.value.orBranches()
}
func (r *repeater) selfString(reverse bool) string {
	sign := "+"
	if r.list.abstractions.UseStarInsteadOfPlus {
		sign = "*"
	}
	return "(" + r.value.selfString(reverse) + ")" + sign
}
func (r *repeater) neatString(reverse bool, tail *stack) string {
	if reverse {
		return r.selfString(reverse) + tail.NeatString(reverse)
	} else {
		return tail.NeatString(reverse) + r.selfString(reverse)
	}
}
func (sl *stackList) newRepeater(value segment) segment {
	newRepeater := &repeater{
		list:              sl,
		value:             value,
		segmentExtensions: map[segment]segment{},
		orExtensions:      map[segment]segment{},
	}
	newRepeater.orExtensions[newRepeater] = newRepeater
	newRepeater.orExtensions[value] = newRepeater
	value.registerRepeater(newRepeater)
	return newRepeater
}
