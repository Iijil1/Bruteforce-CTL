package abstracttapestack

import (
	"fmt"

	TM "github.com/Iijil1/Bruteforce-CTL/turingmachines"
)

type symbol struct {
	list              *stackList
	value             TM.Symbol
	segmentExtensions map[segment]segment
	repeater          segment
	orExtensions      map[segment]segment
}

func (s *symbol) pop(tail *stack) ([]TM.Symbol, []AbstractStack) {
	return []TM.Symbol{s.value}, []AbstractStack{tail}
}
func (s *symbol) extend(value segment) segment {
	if extension, exists := s.segmentExtensions[value]; exists {
		return extension
	} else {
		newNumber := s.list.newSequence(value, s)
		s.segmentExtensions[value] = newNumber
		return newNumber
	}
}
func (s *symbol) repeat() segment {
	if s.repeater != nil {
		return s.repeater
	} else {
		return s.list.newRepeater(s)
	}
}
func (s *symbol) unrepeat() segment {
	return s
}
func (s *symbol) setupRepeater(currentStack *stack, baseStack *stack) {
	if list, exists := currentStack.extensions[s.value]; exists {
		currentStack.extensions[s.value] = append(list, baseStack)
	} else {
		if _, isOrSegment := baseStack.top.unrepeat().(*orSegment); isOrSegment {
			currentStack.extensions[s.value] = []AbstractStack{baseStack, currentStack.pushSegment(s)}
		} else {
			currentStack.extensions[s.value] = []AbstractStack{baseStack}
		}
	}
}
func (s *symbol) registerRepeater(entry segment) {
	s.repeater = entry
	s.registerOr(entry, entry)
}
func (*symbol) numRepeaters() int {
	return 0
}
func (s *symbol) or(value2 segment) segment {
	if extension, exists := s.orExtensions[value2]; exists {
		return extension
	} else {
		return s.list.newOrSegment(s, value2)
	}
}
func (s *symbol) registerOr(value2 segment, extension segment) {
	s.orExtensions[value2] = extension
	s.orExtensions[extension] = extension
}
func (*symbol) orBranches() int {
	return 1
}
func (s *symbol) selfString(bool) string {
	return fmt.Sprint(s.value)
}
func (s *symbol) neatString(reverse bool, tail *stack) string {
	if reverse {
		return s.selfString(reverse) + tail.NeatString(reverse)
	} else {
		return tail.NeatString(reverse) + s.selfString(reverse)
	}
}
