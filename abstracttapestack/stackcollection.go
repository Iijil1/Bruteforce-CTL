package abstracttapestack

import TM "github.com/Iijil1/Bruteforce-CTL/turingmachines"

type stackList struct {
	abstractions   Heuristics
	symbolsInUse   TM.SymbolCollection
	symbolSegments map[TM.Symbol]segment
	emptyStack     *stack
	anyEndStack    *stack
}

func (sl *stackList) getSegmentBySymbol(query TM.Symbol) segment {
	if symbolSegment, exists := sl.symbolSegments[query]; exists {
		return symbolSegment
	} else {
		newSegment := &symbol{
			list:              sl,
			value:             query,
			segmentExtensions: map[segment]segment{},
			repeater:          nil,
			orExtensions:      map[segment]segment{},
		}
		newSegment.orExtensions[newSegment] = newSegment
		sl.symbolSegments[query] = newSegment
		return newSegment
	}
}

func (sl *stackList) AnyEndStack() AbstractStack {
	if sl.abstractions.Buffer > 0 {
		return stackWithBuffer{
			list:               sl,
			bufferEncoding:     0,
			bufferLength:       0,
			stackWithoutBuffer: sl.anyEndStack,
		}
	}
	return sl.anyEndStack
}

func (sl *stackList) EmptyStack() AbstractStack {
	if sl.abstractions.Buffer > 0 {
		return stackWithBuffer{
			list:               sl,
			bufferEncoding:     0,
			bufferLength:       0,
			stackWithoutBuffer: sl.emptyStack,
		}
	}
	return sl.emptyStack
}

func NewStackCollection(tm TM.TuringMachine, abstractions Heuristics) StackCollection {
	ns := &stackList{
		abstractions:   abstractions,
		symbolsInUse:   tm,
		symbolSegments: map[TM.Symbol]segment{},
	}

	infZeroSegment := &infZeroes{
		list:              ns,
		segmentExtensions: map[segment]segment{},
		orExtensions:      map[segment]segment{},
	}
	infZeroSegment.orExtensions[infZeroSegment] = infZeroSegment
	infAnySegment := &infAnys{
		list:              ns,
		segmentExtensions: map[segment]segment{},
		orExtensions:      map[segment]segment{},
	}
	infAnySegment.orExtensions[infAnySegment] = infAnySegment
	emptyStack := &stack{
		list:              ns,
		top:               infZeroSegment,
		length:            0,
		numRepeaters:      0,
		extensions:        map[TM.Symbol][]AbstractStack{},
		segmentExtensions: map[segment]*stack{},
	}
	anyEndStack := &stack{
		list:              ns,
		top:               infAnySegment,
		length:            0,
		numRepeaters:      1,
		extensions:        map[TM.Symbol][]AbstractStack{},
		segmentExtensions: map[segment]*stack{},
	}
	emptyStack.tail = emptyStack
	emptyStack.cutOff = anyEndStack
	emptyStack.extensions[tm.Zero()] = []AbstractStack{emptyStack}
	emptyStack.segmentExtensions[ns.getSegmentBySymbol(tm.Zero())] = emptyStack
	anyEndStack.tail = anyEndStack
	anyEndStack.cutOff = anyEndStack

	ns.emptyStack = emptyStack
	ns.anyEndStack = anyEndStack
	return ns
}
