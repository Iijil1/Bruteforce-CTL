package abstracttapestack

import TM "github.com/Iijil1/Bruteforce-CTL/turingmachines"

type orSegment struct {
	list              *stackList
	value1            segment
	value2            segment
	segmentExtensions map[segment]segment
	repeater          segment
	orExtensions      map[segment]segment
}

func (os *orSegment) pop(tail *stack) ([]TM.Symbol, []AbstractStack) {
	values1, tails1 := os.value1.pop(tail)
	values2, tails2 := os.value2.pop(tail)
	return append(values1, values2...), append(tails1, tails2...)
}
func (os *orSegment) extend(value segment) segment {
	if extension, exists := os.segmentExtensions[value]; exists {
		return extension
	} else {
		newNumber := os.list.newSequence(value, os)
		os.segmentExtensions[value] = newNumber
		return newNumber
	}
}
func (os *orSegment) repeat() segment {
	if os.repeater != nil {
		return os.repeater
	} else {
		var newSegment segment
		if os.unrepeat() == os {
			newSegment = os.list.newRepeater(os)
		} else {
			newSegment = os.unrepeat().repeat()
			os.registerOr(newSegment, newSegment)
		}
		os.repeater = newSegment
		return newSegment
	}
}
func (os *orSegment) unrepeat() segment {
	segment1 := os.value1.unrepeat()
	segment2 := os.value2.unrepeat()
	return segment1.or(segment2)
}

func (os *orSegment) setupRepeater(currentStack *stack, baseStack *stack) {
	os.value1.setupRepeater(currentStack, baseStack)
	os.value2.setupRepeater(currentStack, baseStack)
}
func (s *orSegment) registerRepeater(entry segment) {
	s.repeater = entry
	s.registerOr(entry, entry)
}
func (os *orSegment) numRepeaters() int {
	return os.value1.numRepeaters() + os.value2.numRepeaters()
}
func (os *orSegment) or(value3 segment) segment {
	if extension, exists := os.orExtensions[value3]; exists {
		return extension
	} else {
		newNumber := os.value1.or(os.value2.or(value3))
		os.orExtensions[value3] = newNumber
		return newNumber
	}
}
func (os *orSegment) registerOr(value3 segment, extension segment) {
	os.orExtensions[value3] = extension
	os.value1.registerOr(os.value2.or(value3), extension)
	os.value1.registerOr(extension, extension)
	os.value2.registerOr(os.value1.or(value3), extension)
	os.value2.registerOr(extension, extension)
}
func (os *orSegment) orBranches() int {
	return os.value1.orBranches() + os.value2.orBranches()
}
func (os *orSegment) selfString(reverse bool) string {
	return os.value1.selfString(reverse) + "|" + os.value2.selfString(reverse)
}
func (s *orSegment) neatString(reverse bool, tail *stack) string {
	if reverse {
		return "(" + s.selfString(reverse) + ")" + tail.NeatString(reverse)
	} else {
		return tail.NeatString(reverse) + "(" + s.selfString(reverse) + ")"
	}
}
func (sl *stackList) newOrSegment(value1 segment, value2 segment) segment {
	var NewSegment segment
	if repeater1, isRepeater := value1.(*repeater); isRepeater && repeater1.unrepeat().or(value2.unrepeat()) == repeater1.unrepeat() {
		NewSegment = value1
	} else if repeater2, isRepeater := value2.(*repeater); isRepeater && repeater2.unrepeat().or(value1.unrepeat()) == repeater2.unrepeat() {
		NewSegment = value2
	} else {
		newOrDouble := &orSegment{
			list:              sl,
			value1:            value1,
			value2:            value2,
			segmentExtensions: map[segment]segment{},
			repeater:          nil,
			orExtensions:      map[segment]segment{},
		}
		newOrDouble.orExtensions[newOrDouble] = newOrDouble
		newOrDouble.orExtensions[value1] = newOrDouble
		newOrDouble.orExtensions[value2] = newOrDouble
		NewSegment = newOrDouble
	}
	value1.registerOr(value2, NewSegment)
	value2.registerOr(value1, NewSegment)

	return NewSegment
}
