package decider

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	ATS "github.com/Iijil1/Bruteforce-CTL/abstracttapestack"
	TM "github.com/Iijil1/Bruteforce-CTL/turingmachines"
)

func TestIndividualMachine(t *testing.T) {
	tm := TM.GetTuringMachineFromString("1RB0LD_1LC1RC_0RB1RA_1RC0RE_1RD---")
	fmt.Println(TM.GetTMTable(tm))
	fmt.Println("--------------")
	BruteforceCTL(tm, DeciderOptions{
		StepLimit:     1000,
		InitialDepth:  100,
		DepthIncrease: 0,
		StackHeuristics: ATS.Heuristics{
			SoftLengthLimitAny:           0,
			AdditionalSymbolsToAny:       0,
			LengthLimit:                  0,
			RepetitionsBeforeAbstraction: 1,
			RepetitionLimit:              1,
			NestedRepetitions:            false,
			NoRepeaterUntilLength:        0,
			FixedRepetitionLength:        0,
			RepetitionLengthLowerLimit:   0,
			ConsecutiveRepeatersToOr:     false,
			MaxOrBranches:                0,
			UseStarInsteadOfPlus:         false,
			Buffer:                       0,
		},
	}, ResultsPrintOptions())
}
func TestOptionsEfficiency(t *testing.T) {
	undecided_indices, err := ioutil.ReadFile("../remaining_index")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	totalIndices := len(undecided_indices) / 4
	tests := 100
	decided := 0

	options := DeciderOptions{
		StepLimit:     50000,
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

	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	for i := 0; i < tests; i += 1 {
		randomNumber := rand.Intn(totalIndices)
		index := binary.BigEndian.Uint32(undecided_indices[randomNumber*4 : (randomNumber+1)*4])
		tm := TM.GetTMFromDBIndex("../all_5_states_undecided_machines_with_global_header", index)
		if result, _ := BruteforceCTL(tm, options, SilentPrintOptions()); result {
			decided += 1
		}
	}
	duration := time.Since(start)
	fmt.Printf("Options:\n%v\n\n", options.IndentedString("\t"))
	fmt.Printf("Decided %v out of %v machines in %v\n", decided, tests, duration)
	fmt.Printf("Might decide %v out of %v machines in %v\n", decided*totalIndices/tests, totalIndices, time.Duration(int(duration)*totalIndices/tests))
}

func TestDeciderCTLIijil(t *testing.T) {

	t.Run("NoAbstractionStacksSimulateCorrectly", func(t *testing.T) {
		OldBB5champion := TM.GetTuringMachineFromString("1RB0LC_1RC1RD_1LA0RB_0RE1RH_1LC1RA")
		if result, info := BruteforceCTL(OldBB5champion, DeciderOptions{500000, 1000000, 0, ATS.NoAbstraction()}, SilentPrintOptions()); result || info.Steps != 268932 {
			t.Fail()
		}
	})

	t.Run("MultiSymbolSimulateCorrectly", func(t *testing.T) {
		BB2x3Champion := TM.GetTuringMachineFromString("1RB2LB1RH_2LA2RB1LB")
		if result, info := BruteforceCTL(BB2x3Champion, DeciderOptions{100, 1000, 0, ATS.NoAbstraction()}, SilentPrintOptions()); result || info.Steps != 74 {
			t.Fail()
		}
	})

	t.Run("AbstractionsDoNotDecideBB5Champion", func(t *testing.T) {
		BB5champion := TM.GetTuringMachineFromString("1RB1LC_1RC1RB_1RD0LE_1LA1LD_1RZ0LA")
		if result, _ := BruteforceCTL(BB5champion, DeciderOptions{100000, 10, 1, ATS.SoftLimitAny(0)}, SilentPrintOptions()); result {
			t.Fail()
		}
	})

	t.Run("NoAbstractionsDecideACycler", func(t *testing.T) {
		tm := TM.GetTuringMachineFromString("1RB0RC_0LA0LA_1LD---_1RE1LB_1LB0RD")
		if result, _ := BruteforceCTL(tm, DeciderOptions{1000, 10, 1, ATS.NoAbstraction()}, SilentPrintOptions()); !result {
			t.Fail()
		}
	})

	t.Run("AnyAbstractionsDecideATranslatedCycler", func(t *testing.T) {
		tm := TM.GetTuringMachineFromString("1RB0RB_1LC1LD_0LE0RB_0RA1RD_0RA---")
		if result, _ := BruteforceCTL(tm, DeciderOptions{1000, 10, 2, ATS.SoftLimitAny(2)}, SilentPrintOptions()); !result {
			t.Fail()
		}
	})

	t.Run("AnyAbstractionsWithLimitDecideACounter", func(t *testing.T) {
		tm := TM.GetTuringMachineFromString("1RB1LA_0LC0RB_0LD0LB_1RE---_1LE1LA")
		if result, _ := BruteforceCTL(tm, DeciderOptions{1000, 10, 1, ATS.Heuristics{
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
			Buffer:                       0,
		}}, SilentPrintOptions()); !result {
			t.Fail()
		}
	})

	t.Run("RepitionsWithLimitDecideABouncer", func(t *testing.T) {
		tm := TM.GetTuringMachineFromString("1RB1LD_1RC1LC_0LA0RE_0LC0RB_0RB---")
		if result, _ := BruteforceCTL(tm, DeciderOptions{1000, 50, 10, ATS.RepetitionsWithLimit(2, 1, 100)}, SilentPrintOptions()); !result {
			t.Fail()
		}
	})

	t.Run("RepitionsWithLimitDecideABell", func(t *testing.T) {
		tm := TM.GetTuringMachineFromString("1RB0LE_1RC1LA_0LD1RB_1RC1LD_---1LC")
		if result, _ := BruteforceCTL(tm, DeciderOptions{200, 50, 10, ATS.RepetitionsWithLimit(1, 2, 500)}, SilentPrintOptions()); !result {
			t.Fail()
		}
	})

	t.Run("RepitionsWithLimitDecideAStitchedBouncer", func(t *testing.T) {
		tm := TM.GetTuringMachineFromString("1RB1LA_0RC0RE_1LD---_1LE1RD_0RA0LD")
		if result, _ := BruteforceCTL(tm, DeciderOptions{100000, 200, 20, ATS.RepetitionsWithLimit(2, 2, 100)}, SilentPrintOptions()); !result {
			t.Fail()
		}
	})

	t.Run("OrSegmentsDecideACounter", func(t *testing.T) {
		tm := TM.GetTuringMachineFromString("1RB1LC_0LA1RD_0RD1RE_1LA1RC_---0LA")
		if result, _ := BruteforceCTL(tm, DeciderOptions{100000, 100, 20, ATS.Heuristics{
			SoftLengthLimitAny:           0,
			AdditionalSymbolsToAny:       0,
			LengthLimit:                  100,
			RepetitionsBeforeAbstraction: 1,
			RepetitionLimit:              2,
			NestedRepetitions:            false,
			NoRepeaterUntilLength:        2,
			FixedRepetitionLength:        3,
			RepetitionLengthLowerLimit:   0,
			ConsecutiveRepeatersToOr:     true,
			MaxOrBranches:                0,
			UseStarInsteadOfPlus:         true,
		}}, SilentPrintOptions()); !result {
			t.Fail()
		}
	})

	t.Run("ExtremeAbstractionWorksForScaryMachine", func(t *testing.T) {
		tm := TM.GetTuringMachineFromString("1RB0LD_1LC1RC_1LA0RC_---0LE_0RB1LD")
		if result, _ := BruteforceCTL(tm, DeciderOptions{500000, 20, 5, ATS.Heuristics{
			SoftLengthLimitAny:           1,
			AdditionalSymbolsToAny:       4,
			LengthLimit:                  0,
			RepetitionsBeforeAbstraction: 0,
			RepetitionLimit:              2,
			NestedRepetitions:            false,
			NoRepeaterUntilLength:        2,
			FixedRepetitionLength:        0,
			RepetitionLengthLowerLimit:   0,
			ConsecutiveRepeatersToOr:     true,
			MaxOrBranches:                2,
			UseStarInsteadOfPlus:         true,
			Buffer:                       1,
		}}, SilentPrintOptions()); !result {
			t.Fail()
		}
	})

	t.Run("OrRepeaterForStitchedCounterBouncer", func(t *testing.T) {
		tm := TM.GetTuringMachineFromString("1LC0LC_---0RD_1RD1LA_0RB1RE_1LC0RD")
		if result, _ := BruteforceCTL(tm, DeciderOptions{100000, 20, 5, ATS.Heuristics{
			SoftLengthLimitAny:           0,
			AdditionalSymbolsToAny:       0,
			LengthLimit:                  5,
			RepetitionsBeforeAbstraction: 0,
			RepetitionLimit:              2,
			NestedRepetitions:            false,
			NoRepeaterUntilLength:        0,
			FixedRepetitionLength:        2,
			RepetitionLengthLowerLimit:   0,
			ConsecutiveRepeatersToOr:     true,
			MaxOrBranches:                2,
			UseStarInsteadOfPlus:         false,
			Buffer:                       0,
		}}, SilentPrintOptions()); !result {
			t.Fail()
		}
	})
}
func TestSkeletMachines(t *testing.T) {
	// skelet_machines := []uint32{68329601, 55767995, 5950405, 6897876, 60581745, 58211439, 7196989, 7728246, 12554268, 3810716, 3810169, 4982511, 7566785, 31357173, 2204428, 20569060, 1365166, 15439451, 14536286, 347505, 9980689, 45615747, 6237150, 60658955, 47260245, 13134219, 7163434, 5657318, 6626162, 4986661, 56967673, 6957734, 11896833, 11896832, 11896831, 13609549, 7512832, 35771936, 9914965, 3841616, 5915217, 57874080, 5878998}
	skelet_machines := []uint32{68329601, 55767995, 5950405}

	for _, index := range skelet_machines {
		t.Run("Skelet Machine", func(t *testing.T) {
			tm := TM.GetTMFromDBIndex("../all_5_states_undecided_machines_with_global_header", index)
			if result, _ := BruteforceCTL(tm, DeciderOptions{500000, 10, 1, ATS.Heuristics{
				SoftLengthLimitAny:           0,
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
			}}, SilentPrintOptions()); result {
				fmt.Println("Decided Skelet Machine", index, ":")
				TM.GetStandardTMFormat(tm)
				t.Fail()
			}
		})
	}
}

//TODO: more tests can't hurt. Tests for using the decider with macro machines?
