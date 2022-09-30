package abstracttapestack

import (
	"fmt"

	TM "github.com/Iijil1/Bruteforce-CTL/turingmachines"
)

type StackCollection interface {
	EmptyStack() AbstractStack
	AnyEndStack() AbstractStack
}

type AbstractStack interface {
	Push(TM.Symbol) []AbstractStack
	Pop() ([]TM.Symbol, []AbstractStack)
	String() string
	NeatString(bool) string
	Length() int
}

type Heuristics struct {
	SoftLengthLimitAny           int
	AdditionalSymbolsToAny       int
	LengthLimit                  int
	RepetitionsBeforeAbstraction int
	RepetitionLimit              int
	NestedRepetitions            bool
	NoRepeaterUntilLength        int
	FixedRepetitionLength        int
	RepetitionLengthLowerLimit   int
	ConsecutiveRepeatersToOr     bool
	MaxOrBranches                int
	UseStarInsteadOfPlus         bool
	Buffer                       int
}

func NoAbstraction() Heuristics {
	return Heuristics{}
}
func SoftLimitAny(limit int) Heuristics {
	return Heuristics{SoftLengthLimitAny: limit + 1, LengthLimit: limit}
}
func RepetitionsWithLimit(RepetitionsUntilTry int, repetitionLimit int, lengthLimit int) Heuristics {
	return Heuristics{LengthLimit: lengthLimit, RepetitionsBeforeAbstraction: RepetitionsUntilTry, RepetitionLimit: repetitionLimit}
}
func FullAbstraction(repetitionLimit int, lengthLimit int) Heuristics {
	return Heuristics{SoftLengthLimitAny: 1, LengthLimit: lengthLimit, RepetitionsBeforeAbstraction: 0, NestedRepetitions: true, RepetitionLimit: repetitionLimit}
}

func (heuristics Heuristics) IndentedString(indent string) string {
	indentedString := fmt.Sprintf("%vSoftLengthLimitAny: %v\n", indent, heuristics.SoftLengthLimitAny)
	indentedString += fmt.Sprintf("%vAdditionalSymbolsToAny: %v\n", indent, heuristics.AdditionalSymbolsToAny)
	indentedString += fmt.Sprintf("%vLengthLimit: %v\n", indent, heuristics.LengthLimit)
	indentedString += fmt.Sprintf("%vRepetitionsBeforeAbstraction: %v\n", indent, heuristics.RepetitionsBeforeAbstraction)
	indentedString += fmt.Sprintf("%vRepetitionLimit: %v\n", indent, heuristics.RepetitionLimit)
	indentedString += fmt.Sprintf("%vNestedRepetitions: %v\n", indent, heuristics.NestedRepetitions)
	indentedString += fmt.Sprintf("%vNoRepeaterUntilLength: %v\n", indent, heuristics.NoRepeaterUntilLength)
	indentedString += fmt.Sprintf("%vFixedRepetitionLength: %v\n", indent, heuristics.FixedRepetitionLength)
	indentedString += fmt.Sprintf("%vRepetitionLengthLowerLimit: %v\n", indent, heuristics.RepetitionLengthLowerLimit)
	indentedString += fmt.Sprintf("%vConsecutiveRepeatersToOr: %v\n", indent, heuristics.ConsecutiveRepeatersToOr)
	indentedString += fmt.Sprintf("%vMaxOrBranches: %v\n", indent, heuristics.MaxOrBranches)
	indentedString += fmt.Sprintf("%vUseStarInsteadOfPlus: %v\n", indent, heuristics.UseStarInsteadOfPlus)
	indentedString += fmt.Sprintf("%vBuffer: %v", indent, heuristics.Buffer)
	return indentedString
}
