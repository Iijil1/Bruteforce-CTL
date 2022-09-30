package abstracttapestack

import (
	"fmt"
	"testing"

	TM "github.com/Iijil1/Bruteforce-CTL/turingmachines"
)

func TestSimpleStacks(t *testing.T) {
	tm := TM.GetTuringMachineFromString("1RB0LD_1LC1RC_1LA0RC_---0LE_0RB1LD")
	stackList := NewStackCollection(tm, NoAbstraction())
	emptyStack := stackList.EmptyStack()
	t.Run("emptyStackStringIsEmpty", func(t *testing.T) {
		if emptyStack.NeatString(true) != "(0)@" || emptyStack.NeatString(false) != "(0)@" {
			t.Fail()
		}
	})
	t.Run("EmptyStackPlusZeroStaysEmpty", func(t *testing.T) {
		testStack := emptyStack.Push(tm.GetSymbol(0))[0]
		if emptyStack != testStack {
			t.Fail()
		}
	})
	t.Run("EmptyStackPlusOneIsDifferent", func(t *testing.T) {
		testStack := emptyStack.Push(tm.GetSymbol(1))[0]
		if emptyStack == testStack {
			t.Fail()
		}
	})
	t.Run("PoppingEmptyStackGivesZero", func(t *testing.T) {
		if Symbol, _ := emptyStack.Pop(); Symbol[0] != tm.Zero() {
			t.Fail()
		}
	})
	t.Run("PoppingEmptyStackGivesEmptyStack", func(t *testing.T) {
		if _, Stack := emptyStack.Pop(); Stack[0] != stackList.EmptyStack() {
			t.Fail()
		}
	})
	testStack1 := emptyStack.Push(tm.GetSymbol(1))[0]
	testStack10 := testStack1.Push(tm.GetSymbol(0))[0]
	testStack101 := testStack10.Push(tm.GetSymbol(1))[0]
	testStack1011 := testStack101.Push(tm.GetSymbol(1))[0]
	t.Run("StringsForSimpleStacks", func(t *testing.T) {
		if testStack1011.NeatString(false) != "(0)@1011" {
			fmt.Println("Expected (0)@1011 but got", emptyStack.NeatString(false))
			t.Fail()
		}
	})
	t.Run("ReverseStringsForSimpleStacks", func(t *testing.T) {
		if testStack1011.NeatString(true) != "1101(0)@" {
			fmt.Println("Expected 1101(0)@ but got", emptyStack.NeatString(true))
			t.Fail()
		}
	})
	t.Run("PoppingOneGivesCorrectSymbol", func(t *testing.T) {
		if Symbol, _ := testStack1011.Pop(); Symbol[0].String() != "1" {
			t.Fail()
		}
	})
	t.Run("PoppingZeroGivesCorrectSymbol", func(t *testing.T) {
		if Symbol, _ := testStack10.Pop(); Symbol[0].String() != "0" {
			t.Fail()
		}
	})
	t.Run("PoppingGivesCorrectStack", func(t *testing.T) {
		if _, Stack := testStack1011.Pop(); Stack[0] != testStack101 {
			t.Fail()
		}
	})
	testStack1011copy := testStack101.Push(tm.GetSymbol(1))[0]
	t.Run("PushingOneAgainReusesStack", func(t *testing.T) {
		if testStack1011copy != testStack1011 {
			t.Fail()
		}
	})
	testStack10copy := testStack1.Push(tm.GetSymbol(0))[0]
	t.Run("PushingOneAgainReusesStack", func(t *testing.T) {
		if testStack10copy != testStack10 {
			t.Fail()
		}
	})
}

func TestAnyStacks(t *testing.T) {
	tm := TM.GetTuringMachineFromString("1RB0LD_1LC1RC_1LA0RC_---0LE_0RB1LD")
	stackList := NewStackCollection(tm, Heuristics{SoftLengthLimitAny: 1, AdditionalSymbolsToAny: 0, LengthLimit: 0})
	emptyStack := stackList.EmptyStack()
	testStacks1 := emptyStack.Push(tm.GetSymbol(1))
	t.Run("MultivaluePushOnEmpty", func(t *testing.T) {
		if len(testStacks1) != 2 {
			t.Fail()
		}
	})
	anyEndStack := testStacks1[0]
	t.Run("StringForAnyStack", func(t *testing.T) {
		if anyEndStack.NeatString(false) != "(.)@" {
			fmt.Println("Expected (.)@ but got", anyEndStack.NeatString(false))
			t.Fail()
		}
	})
	t.Run("ReverseStringForAnyStack", func(t *testing.T) {
		if anyEndStack.NeatString(true) != "(.)@" {
			fmt.Println("Expected (.)@ but got", anyEndStack.NeatString(true))
			t.Fail()
		}
	})
	simpleStack := testStacks1[1]
	t.Run("SimpleAlternativeCorrect", func(t *testing.T) {
		if simpleStack.NeatString(false) != "(0)@1" {
			fmt.Println("Expected (0)@1 but got", simpleStack.NeatString(false))
			t.Fail()
		}
	})
	testStacks2 := anyEndStack.Push(tm.GetSymbol(0))
	t.Run("MultivaluePushZeroOnAny", func(t *testing.T) {
		if len(testStacks2) != 2 {
			t.Fail()
		}
	})
	t.Run("PushZeroOnAnyGivesAny", func(t *testing.T) {
		if testStacks2[0] != anyEndStack {
			fmt.Println("Expected (.)@ but got", testStacks2[0].NeatString(true))
			t.Fail()
		}
	})
	zeroAnyStack := testStacks2[1]
	t.Run("PushZeroOnAnyGivesCorrectAlternative", func(t *testing.T) {
		if zeroAnyStack.NeatString(false) != "(.)@0" {
			fmt.Println("Expected (.)@0 but got", zeroAnyStack.NeatString(false))
			t.Fail()
		}
	})
	testStacks3 := anyEndStack.Push(tm.GetSymbol(1))
	t.Run("MultivaluePushOneOnAny", func(t *testing.T) {
		if len(testStacks3) != 2 {
			t.Fail()
		}
	})
	t.Run("PushOneOnAnyGivesAny", func(t *testing.T) {
		if testStacks3[0] != anyEndStack {
			t.Fail()
		}
	})
	oneAnyStack := testStacks3[1]
	t.Run("PushOneOnAnyGivesCorrectAlternative", func(t *testing.T) {
		if oneAnyStack.NeatString(false) != "(.)@1" {
			fmt.Println("Expected (.)@1 but got", zeroAnyStack.NeatString(false))
			t.Fail()
		}
	})
	t.Run("ReverseStringsForLongerAnyStacks", func(t *testing.T) {
		if oneAnyStack.NeatString(true) != "1(.)@" {
			fmt.Println("Expected 1(.)@ but got", oneAnyStack.NeatString(true))
			t.Fail()
		}
	})
	t.Run("MultiValuePopAny", func(t *testing.T) {
		if Values, _ := anyEndStack.Pop(); len(Values) != 2 {
			t.Fail()
		}
	})
	t.Run("PopAnyGivesAny", func(t *testing.T) {
		if _, Stacks := anyEndStack.Pop(); len(Stacks) != 2 || Stacks[0] != anyEndStack || Stacks[1] != anyEndStack {
			t.Fail()
		}
	})
	t.Run("PopOneAnyGivesAny", func(t *testing.T) {
		if _, Stacks := oneAnyStack.Pop(); len(Stacks) != 1 || Stacks[0] != anyEndStack {
			t.Fail()
		}
	})
	t.Run("PopOneAnyGivesOne", func(t *testing.T) {
		if Values, _ := oneAnyStack.Pop(); len(Values) != 1 || Values[0].String() != "1" {
			t.Fail()
		}
	})
	testStacks4 := oneAnyStack.Push(tm.GetSymbol(0))
	t.Run("MultivaluePushZeroOnOneAny", func(t *testing.T) {
		if len(testStacks4) != 2 {
			t.Fail()
		}
	})
	t.Run("PushZeroOnOneAnyGivesZeroAny", func(t *testing.T) {
		if testStacks4[0] != zeroAnyStack {
			fmt.Println(testStacks4[0])
			fmt.Println(zeroAnyStack)
			t.Fail()
		}
	})
	zeroOneAnyStack := testStacks4[1]
	t.Run("PushZeroOnOneAnyGivesCorrectAlternative", func(t *testing.T) {
		if zeroOneAnyStack.NeatString(false) != "(.)@10" {
			fmt.Println("Expected (.)@10 but got", zeroAnyStack.NeatString(false))
			t.Fail()
		}
	})
}

func TestLengthLimit(t *testing.T) {
	tm := TM.GetTuringMachineFromString("1RB0LD_1LC1RC_1LA0RC_---0LE_0RB1LD")
	stackList := NewStackCollection(tm, SoftLimitAny(2))

	lengthTestStack := stackList.EmptyStack()
	t.Run("EmptyStackHaseZeroLength", func(t *testing.T) {
		if lengthTestStack.Length() != 0 {
			t.Fail()
		}
	})
	lengthTestStack = lengthTestStack.Push(tm.GetSymbol(1))[0]
	t.Run("LengthGoesToOne", func(t *testing.T) {
		if lengthTestStack.Length() != 1 {
			t.Fail()
		}
	})
	lengthTestStack = lengthTestStack.Push(tm.GetSymbol(0))[0]
	t.Run("LengthGoesToTwo", func(t *testing.T) {
		if lengthTestStack.Length() != 2 {
			t.Fail()
		}
	})
	lengthTestStack = lengthTestStack.Push(tm.GetSymbol(1))[0]
	t.Run("CutOffLengthStaysTwo", func(t *testing.T) {
		if lengthTestStack.Length() != 2 {
			t.Fail()
		}
	})
}

//TODO: tests for sequences, repeaters, orsegments
