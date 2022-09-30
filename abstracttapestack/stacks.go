package abstracttapestack

import TM "github.com/Iijil1/Bruteforce-CTL/turingmachines"

type stack struct {
	list              *stackList
	top               segment
	tail              *stack
	cutOff            *stack
	length            int
	numRepeaters      int
	extensions        map[TM.Symbol][]AbstractStack
	segmentExtensions map[segment]*stack
}

func (s *stack) Push(pushedSymbol TM.Symbol) []AbstractStack {
	if extensions, exists := s.extensions[pushedSymbol]; exists {
		return extensions
	} else {
		possibleStacks := []AbstractStack{}
		if s.list.abstractions.RepetitionLimit > s.numRepeaters {
			startSegment := s.list.getSegmentBySymbol(pushedSymbol)
			currentTail := s
			for i := 1; (s.list.abstractions.RepetitionsBeforeAbstraction)*i+s.list.abstractions.NoRepeaterUntilLength <= currentTail.length; i += 1 {
				fastTail := currentTail
				repetitionFound := true
				for repeats := 0; repeats < s.list.abstractions.RepetitionsBeforeAbstraction; repeats += 1 {
					endSegment := fastTail.top
					fastTail = fastTail.tail
					for j := 2; j <= i; j += 1 {
						endSegment = endSegment.extend(fastTail.top)
						fastTail = fastTail.tail
					}
					if startSegment != endSegment {
						repetitionFound = false
						break
					}
				}
				if repetitionFound && fastTail != s.list.anyEndStack &&
					(s.list.abstractions.FixedRepetitionLength == 0 || s.list.abstractions.FixedRepetitionLength == i) &&
					(s.list.abstractions.RepetitionLengthLowerLimit == 0 || s.list.abstractions.RepetitionLengthLowerLimit <= i) {
					repeatedSegment := startSegment.repeat()
					if s.list.abstractions.ConsecutiveRepeatersToOr {
						if oldRepeater, isRepeater := fastTail.top.(*repeater); isRepeater {
							repeatedOrSegment := oldRepeater.or(repeatedSegment).repeat()
							if s.list.abstractions.MaxOrBranches == 0 || repeatedOrSegment.orBranches() <= s.list.abstractions.MaxOrBranches {
								possibleStacks = append(possibleStacks, fastTail.tail.pushSegment(repeatedOrSegment))
							}
						}
					}
					possibleStacks = append(possibleStacks, fastTail.pushSegment(repeatedSegment))
				}
				if currentTail.length == 0 {
					break
				}
				if _, isRepeater := currentTail.top.(*repeater); !s.list.abstractions.NestedRepetitions && isRepeater {
					break
				}
				startSegment = startSegment.extend(currentTail.top)
				currentTail = currentTail.tail
			}
		}
		simpleExtension := s.pushSegment(s.list.getSegmentBySymbol(pushedSymbol))
		if s.list.abstractions.SoftLengthLimitAny != 0 && s.list.abstractions.SoftLengthLimitAny <= s.length+1 {
			cutStack := simpleExtension.cutOff
			for i := 0; i < s.list.abstractions.AdditionalSymbolsToAny; i += 1 {
				cutStack = cutStack.cutOff
			}
			possibleStacks = append(possibleStacks, cutStack)
		}
		if s.list.abstractions.LengthLimit == 0 || s.length < s.list.abstractions.LengthLimit {
			possibleStacks = append(possibleStacks, simpleExtension)
		}
		s.extensions[pushedSymbol] = possibleStacks
		return possibleStacks
	}
}
func (s *stack) Pop() ([]TM.Symbol, []AbstractStack) {
	return s.top.pop(s.tail)
}
func (s *stack) pushSegment(newSegment segment) *stack {
	if extension, exists := s.segmentExtensions[newSegment]; exists {
		return extension
	} else {
		newStack := &stack{
			list:              s.list,
			top:               newSegment,
			tail:              s,
			cutOff:            s.list.anyEndStack,
			length:            s.length + 1,
			numRepeaters:      s.numRepeaters + newSegment.numRepeaters(),
			extensions:        map[TM.Symbol][]AbstractStack{},
			segmentExtensions: map[segment]*stack{},
		}
		s.segmentExtensions[newSegment] = newStack
		if s.list.abstractions.SoftLengthLimitAny != 0 && newStack.length > 1 {
			if _, isRepeater := newSegment.(*repeater); s.cutOff == s.list.AnyEndStack() && isRepeater {
				newStack.cutOff = s.list.anyEndStack
			} else {
				newStack.cutOff = s.cutOff.pushSegment(newSegment)
			}
		}
		if newRepeater, isRepeater := newSegment.(*repeater); isRepeater {
			newRepeater.setupRepeater(newStack, newStack)
		}
		return newStack
	}
}

func (s *stack) NeatString(reverse bool) string {
	return s.top.neatString(reverse, s.tail)
}

func (s *stack) String() string {
	return s.NeatString(false)
}

func (s *stack) Length() int {
	return s.length
}
