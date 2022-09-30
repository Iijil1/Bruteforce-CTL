package abstracttapestack

import TM "github.com/Iijil1/Bruteforce-CTL/turingmachines"

type sequence struct {
	list              *stackList
	value             segment
	tail              segment
	segmentExtensions map[segment]segment
	repeater          segment
	orExtensions      map[segment]segment
}

func (s *sequence) pop(tail *stack) ([]TM.Symbol, []AbstractStack) {
	return s.tail.pop(tail.pushSegment(s.value))
}

func (s *sequence) extend(value segment) segment {
	if extension, exists := s.segmentExtensions[value]; exists {
		return extension
	} else {
		newNumber := s.list.newSequence(value, s)
		s.segmentExtensions[value] = newNumber
		return newNumber
	}
}
func (s *sequence) repeat() segment {
	if s.repeater != nil {
		return s.repeater
	} else {
		return s.list.newRepeater(s)
	}
}
func (s *sequence) unrepeat() segment {
	return s
}
func (s *sequence) setupRepeater(currentStack *stack, baseStack *stack) {
	s.tail.setupRepeater(currentStack.pushSegment(s.value), baseStack)
}
func (s *sequence) registerRepeater(entry segment) {
	s.repeater = entry
	s.registerOr(entry, entry)
}
func (s *sequence) numRepeaters() int {
	return s.value.numRepeaters() + s.tail.numRepeaters()
}
func (s *sequence) or(value2 segment) segment {
	if extension, exists := s.orExtensions[value2]; exists {
		return extension
	} else {
		return s.list.newOrSegment(s, value2)
	}
}
func (s *sequence) registerOr(value2 segment, extension segment) {
	s.orExtensions[value2] = extension
	s.orExtensions[extension] = extension
}
func (*sequence) orBranches() int {
	return 1
}
func (s *sequence) selfString(reverse bool) string {
	segmentString := s.value.selfString(reverse)

	if reverse {
		return s.tail.selfString(reverse) + segmentString
	} else {
		return segmentString + s.tail.selfString(reverse)
	}
}

func (s *sequence) neatString(reverse bool, tail *stack) string {
	segmentString := s.selfString(reverse)
	if reverse {
		return segmentString + tail.NeatString(reverse)
	} else {
		return tail.NeatString(reverse) + segmentString
	}
}

func (sl *stackList) newSequence(newValue segment, tail segment) segment {
	newSegment := &sequence{
		list:              sl,
		value:             newValue,
		tail:              tail,
		segmentExtensions: map[segment]segment{},
		repeater:          nil,
		orExtensions:      map[segment]segment{},
	}
	newSegment.orExtensions[newSegment] = newSegment
	return newSegment
}
