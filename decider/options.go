package decider

import (
	"fmt"

	ATS "github.com/Iijil1/Bruteforce-CTL/abstracttapestack"
)

type DeciderOptions struct {
	StepLimit       int
	InitialDepth    int
	DepthIncrease   int
	StackHeuristics ATS.Heuristics
}

func (options DeciderOptions) IndentedString(indent string) string {
	indentedString := fmt.Sprintf("%vStepLimit: %v\n", indent, options.StepLimit)
	indentedString += fmt.Sprintf("%vInitialDepth: %v\n", indent, options.InitialDepth)
	indentedString += fmt.Sprintf("%vDepthIncrease: %v\n", indent, options.DepthIncrease)
	indentedString += fmt.Sprintf("%vStackHeuristics:\n%v", indent, options.StackHeuristics.IndentedString(indent+"\t"))
	return indentedString
}

func AggressiveAbstraction() DeciderOptions {
	return DeciderOptions{
		StepLimit:     100000,
		InitialDepth:  10,
		DepthIncrease: 1,
		StackHeuristics: ATS.Heuristics{
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
	}
}

func BouncerSearch() DeciderOptions {
	return DeciderOptions{
		StepLimit:     10000,
		InitialDepth:  100,
		DepthIncrease: 25,
		StackHeuristics: ATS.Heuristics{
			SoftLengthLimitAny:           0,
			AdditionalSymbolsToAny:       0,
			LengthLimit:                  0,
			RepetitionsBeforeAbstraction: 2,
			RepetitionLimit:              1,
			NestedRepetitions:            false,
			NoRepeaterUntilLength:        5,
			FixedRepetitionLength:        0,
			RepetitionLengthLowerLimit:   0,
			ConsecutiveRepeatersToOr:     false,
			MaxOrBranches:                0,
			UseStarInsteadOfPlus:         false,
			Buffer:                       5,
		},
	}
}

func CounterSearch() DeciderOptions {
	return DeciderOptions{
		StepLimit:     50000,
		InitialDepth:  100,
		DepthIncrease: 20,
		StackHeuristics: ATS.Heuristics{
			SoftLengthLimitAny:           0,
			AdditionalSymbolsToAny:       0,
			LengthLimit:                  0,
			RepetitionsBeforeAbstraction: 1,
			RepetitionLimit:              2,
			NestedRepetitions:            false,
			NoRepeaterUntilLength:        0,
			FixedRepetitionLength:        0,
			RepetitionLengthLowerLimit:   2,
			ConsecutiveRepeatersToOr:     true,
			MaxOrBranches:                2,
			UseStarInsteadOfPlus:         true,
			Buffer:                       0,
		},
	}
}
func ForceCounterSearch() DeciderOptions {
	return DeciderOptions{
		StepLimit:     500000,
		InitialDepth:  50,
		DepthIncrease: 25,
		StackHeuristics: ATS.Heuristics{
			SoftLengthLimitAny:           0,
			AdditionalSymbolsToAny:       0,
			LengthLimit:                  0,
			RepetitionsBeforeAbstraction: 2,
			RepetitionLimit:              2,
			NestedRepetitions:            false,
			NoRepeaterUntilLength:        0,
			FixedRepetitionLength:        0,
			RepetitionLengthLowerLimit:   2,
			ConsecutiveRepeatersToOr:     true,
			MaxOrBranches:                2,
			UseStarInsteadOfPlus:         true,
			Buffer:                       0,
		},
	}
}
func StitchedBouncerSearch() DeciderOptions {
	return DeciderOptions{
		StepLimit:     100000,
		InitialDepth:  300,
		DepthIncrease: 100,
		StackHeuristics: ATS.Heuristics{
			SoftLengthLimitAny:           0,
			AdditionalSymbolsToAny:       0,
			LengthLimit:                  0,
			RepetitionsBeforeAbstraction: 2,
			RepetitionLimit:              2,
			NestedRepetitions:            false,
			NoRepeaterUntilLength:        0,
			FixedRepetitionLength:        0,
			RepetitionLengthLowerLimit:   3,
			ConsecutiveRepeatersToOr:     false,
			MaxOrBranches:                0,
			UseStarInsteadOfPlus:         true,
			Buffer:                       0,
		},
	}
}
func NonBinaryCounters() DeciderOptions {
	return DeciderOptions{
		StepLimit:     500000,
		InitialDepth:  50,
		DepthIncrease: 25,
		StackHeuristics: ATS.Heuristics{
			SoftLengthLimitAny:           0,
			AdditionalSymbolsToAny:       0,
			LengthLimit:                  0,
			RepetitionsBeforeAbstraction: 1,
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
}

//----------------- Print Options -------------------

type PrintOptions struct {
	OutOfSteps           bool
	FoundStartHalting    bool
	FoundHalting         bool
	Success              bool
	CurrentConfiguration bool
	AddedConfiguration   bool
	DiscardedAbstraction bool
	SettingAborted       bool
	StatusCountInfo      bool
	DetailedHaltDecision bool
	DepthLimitIncrease   bool
}

func FullPrintOptions() PrintOptions {
	return PrintOptions{
		OutOfSteps:           true,
		FoundStartHalting:    true,
		FoundHalting:         true,
		Success:              true,
		CurrentConfiguration: true,
		AddedConfiguration:   true,
		DiscardedAbstraction: true,
		SettingAborted:       true,
		StatusCountInfo:      true,
		DetailedHaltDecision: true,
		DepthLimitIncrease:   true,
	}
}

func SilentPrintOptions() PrintOptions {
	return PrintOptions{}
}

func ResultsPrintOptions() PrintOptions {
	return PrintOptions{OutOfSteps: true, FoundStartHalting: true, Success: true, StatusCountInfo: true}
}

func TestingPrintOptions() PrintOptions {
	return PrintOptions{OutOfSteps: true, FoundStartHalting: true, Success: true, CurrentConfiguration: true, FoundHalting: true, StatusCountInfo: true, DepthLimitIncrease: true}
}
