package decider

import (
	"github.com/Iijil1/Bruteforce-CTL/abstracttapestack"
)

func heuristicList(level int) []abstracttapestack.Heuristics {
	heuristicList := []abstracttapestack.Heuristics{
		{ //cycler variants
			SoftLengthLimitAny:           1,
			AdditionalSymbolsToAny:       0,
			LengthLimit:                  0,
			RepetitionsBeforeAbstraction: 0,
			RepetitionLimit:              0,
			NestedRepetitions:            false,
			NoRepeaterUntilLength:        0,
			FixedRepetitionLength:        0,
			RepetitionLengthLowerLimit:   0,
			ConsecutiveRepeatersToOr:     false,
			MaxOrBranches:                0,
			UseStarInsteadOfPlus:         false,
			Buffer:                       10 * (level + 1),
		},
		{ //simple bouncers
			SoftLengthLimitAny:           0,
			AdditionalSymbolsToAny:       0,
			LengthLimit:                  0,
			RepetitionsBeforeAbstraction: 1 + level,
			RepetitionLimit:              1,
			NestedRepetitions:            false,
			NoRepeaterUntilLength:        2 * level,
			FixedRepetitionLength:        0,
			RepetitionLengthLowerLimit:   1 + level,
			ConsecutiveRepeatersToOr:     false,
			MaxOrBranches:                0,
			UseStarInsteadOfPlus:         false,
			Buffer:                       2 * level,
		},
		{ //bouncer variants
			SoftLengthLimitAny:           0,
			AdditionalSymbolsToAny:       0,
			LengthLimit:                  0,
			RepetitionsBeforeAbstraction: 1 + (level / 2),
			RepetitionLimit:              2,
			NestedRepetitions:            false,
			NoRepeaterUntilLength:        (level / 2),
			FixedRepetitionLength:        0,
			RepetitionLengthLowerLimit:   1 + (level / 2),
			ConsecutiveRepeatersToOr:     false,
			MaxOrBranches:                0,
			UseStarInsteadOfPlus:         false,
			Buffer:                       (level / 2),
		},
		{ //binary counters
			SoftLengthLimitAny:           1,
			AdditionalSymbolsToAny:       0,
			LengthLimit:                  0,
			RepetitionsBeforeAbstraction: 0,
			RepetitionLimit:              2,
			NestedRepetitions:            false,
			NoRepeaterUntilLength:        0,
			FixedRepetitionLength:        0,
			RepetitionLengthLowerLimit:   0,
			ConsecutiveRepeatersToOr:     true,
			MaxOrBranches:                2,
			UseStarInsteadOfPlus:         true,
			Buffer:                       0,
		},
		{ //small CTLs
			SoftLengthLimitAny:           1,
			AdditionalSymbolsToAny:       0,
			LengthLimit:                  0,
			RepetitionsBeforeAbstraction: 0,
			RepetitionLimit:              2,
			NestedRepetitions:            false,
			NoRepeaterUntilLength:        0,
			FixedRepetitionLength:        0,
			RepetitionLengthLowerLimit:   0,
			ConsecutiveRepeatersToOr:     true,
			MaxOrBranches:                0,
			UseStarInsteadOfPlus:         true,
			Buffer:                       0,
		},
		{ //other counters
			SoftLengthLimitAny:           0,
			AdditionalSymbolsToAny:       0,
			LengthLimit:                  0,
			RepetitionsBeforeAbstraction: 0,
			RepetitionLimit:              2,
			NestedRepetitions:            false,
			NoRepeaterUntilLength:        0,
			FixedRepetitionLength:        0,
			RepetitionLengthLowerLimit:   2,
			ConsecutiveRepeatersToOr:     true,
			MaxOrBranches:                0,
			UseStarInsteadOfPlus:         true,
			Buffer:                       0,
		},
	}
	//add counters with specific symbol size
	for size := 2; size < level+3; size++ {
		heuristicList = append(heuristicList, abstracttapestack.Heuristics{
			SoftLengthLimitAny:           0,
			AdditionalSymbolsToAny:       0,
			LengthLimit:                  0,
			RepetitionsBeforeAbstraction: 0,
			RepetitionLimit:              2,
			NestedRepetitions:            false,
			NoRepeaterUntilLength:        0,
			FixedRepetitionLength:        size,
			RepetitionLengthLowerLimit:   0,
			ConsecutiveRepeatersToOr:     true,
			MaxOrBranches:                0,
			UseStarInsteadOfPlus:         true,
			Buffer:                       size,
		})
	}
	return heuristicList
}

func optionsList(level int) []DeciderOptions {
	var stepLimit, initialDepth, depthIncrease int
	switch level {
	case 0:
		stepLimit, initialDepth, depthIncrease = 1000, 5, 1
	case 1:
		stepLimit, initialDepth, depthIncrease = 10000, 10, 2
	case 2:
		stepLimit, initialDepth, depthIncrease = 100000, 15, 3
	case 3:
		stepLimit, initialDepth, depthIncrease = 1000000, 20, 4
	}
	optionsList := []DeciderOptions{}
	for _, heuristic := range heuristicList(level) {
		optionsList = append(optionsList, DeciderOptions{
			StepLimit:       stepLimit,
			InitialDepth:    initialDepth,
			DepthIncrease:   depthIncrease,
			ForcedLines:     true,
			StackHeuristics: heuristic,
		})
	}
	return optionsList
}
