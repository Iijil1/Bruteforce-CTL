package prooffileproduction

import (
	"fmt"
	"os"
	"sort"

	ATS "github.com/Iijil1/Bruteforce-CTL/abstracttapestack"
	"github.com/Iijil1/Bruteforce-CTL/decider"
	TM "github.com/Iijil1/Bruteforce-CTL/turingmachines"
)

func fransFaaseStrings(c *decider.Configuration) []string {
	stateString := c.State.String()
	leftString := c.LeftTape.NeatString(false)
	rightString := c.RightTape.NeatString(true)
	resultStrings := []string{}
	var symbols []TM.Symbol
	var stacks []ATS.AbstractStack
	if c.Direction == TM.RIGHT {
		symbols, stacks = c.RightTape.Pop()
	} else {
		symbols, stacks = c.LeftTape.Pop()
	}
	for i := range symbols {
		symbolString := fmt.Sprint(symbols[i])
		if c.Direction == TM.RIGHT {
			rightString = stacks[i].NeatString(true)
		} else {
			leftString = stacks[i].NeatString(false)
		}
		resultStrings = append(resultStrings, trimLeadingZeroes(stateString+": "+leftString+" "+symbolString+" "+rightString))
	}
	return resultStrings
}

func remove(stringMap map[string]bool, changes *bool, removedString string) {
	if print, exists := stringMap[removedString]; !exists {
		*changes = true
		stringMap[removedString] = false
	} else if print {
		stringMap[removedString] = false
	}
}

func trimLeadingZeroes(configString string) string {
	if configString[4] == '0' {
		firstNonZero := 7
		for configString[firstNonZero] == '0' {
			firstNonZero += 1
		}
		configString = configString[:7] + configString[firstNonZero:]
	}

	if configString[len(configString)-3] == '0' {
		lastNonZero := 4
		for configString[len(configString)-lastNonZero-1] == '0' {
			lastNonZero += 1
		}
		configString = configString[:len(configString)-lastNonZero] + configString[len(configString)-4:]
	}
	return configString
}

const NOBODYCARESABOUTRUNTIME = false

func CreateFransFaaseVerificationFile(tm TM.TuringMachine, options decider.DeciderOptions) {
	tmString := tm.NameString()
	if decided, info := decider.BruteforceCTL(tm, options, decider.ResultsPrintOptions()); decided {

		//----------------- Generate initial config strings --------------------

		start := info.StartConfig
		stringsToPrint := map[string]bool{}
		printStack := []*decider.Configuration{start}
		start.PrintHelper = 1
		maxLength := 0
		for len(printStack) > 0 {
			config := printStack[0]
			printStack = printStack[1:]
			for _, configString := range fransFaaseStrings(config) {
				stringsToPrint[configString] = true
				if len(configString) > maxLength {
					maxLength = len(configString)
				}
			}
			for _, successor := range config.Successors {
				if successor.Status == decider.DONE {
					if successor.PrintHelper == 0 {
						successor.PrintHelper = 1
						printStack = append(printStack, successor)
					}
				}
			}
		}
		//----------------- Remove unneccessary strings --------------------

		if NOBODYCARESABOUTRUNTIME {
			changes := true
			for changes == true {
				changes = false
				for configString := range stringsToPrint {
					openBracketsStack := []int{}
					for i := range configString {
						if configString[i] == '(' {
							openBracketsStack = append(openBracketsStack, i)
						}
						if configString[i] == ')' && len(openBracketsStack) > 0 && len(configString) > i+1 {
							j := openBracketsStack[len(openBracketsStack)-1]
							openBracketsStack = openBracketsStack[:len(openBracketsStack)-1]
							if configString[i+1] == '+' || configString[i+1] == '*' {
								preString := configString[:j]
								repValue := configString[j+1 : i]
								postString := configString[i+2:]
								postStringWithRep := repValue + postString
								preStringWithRep := preString + repValue

								if configString[i+1] == '+' {
									if noRepeaterPrint, exists := stringsToPrint[trimLeadingZeroes(preString+postString)]; exists {
										if _, exists := stringsToPrint[trimLeadingZeroes(preString+"("+repValue+")*"+postString)]; !exists {
											changes = true
											stringsToPrint[trimLeadingZeroes(preString+"("+repValue+")*"+postString)] = noRepeaterPrint || stringsToPrint[configString]
										}
										remove(stringsToPrint, &changes, trimLeadingZeroes(preString+postString))
										remove(stringsToPrint, &changes, configString)
									}
								}
								if configString[i+1] == '*' {
									remove(stringsToPrint, &changes, trimLeadingZeroes(preString+postString))
								}
								for k := 0; len(preStringWithRep)+len(postString) <= maxLength; k += 1 {
									remove(stringsToPrint, &changes, trimLeadingZeroes(preStringWithRep+postString))
									remove(stringsToPrint, &changes, trimLeadingZeroes(preStringWithRep+"("+repValue+")*"+postString))
									remove(stringsToPrint, &changes, trimLeadingZeroes(preString+"("+repValue+")*"+postStringWithRep))
									remove(stringsToPrint, &changes, trimLeadingZeroes(preStringWithRep+"("+repValue+")+"+postString))
									remove(stringsToPrint, &changes, trimLeadingZeroes(preString+"("+repValue+")+"+postStringWithRep))

									postStringWithRep = repValue + postStringWithRep
									preStringWithRep = preStringWithRep + repValue
								}
							}
						}
					}
				}
			}
		}
		//----------------- Writing the file --------------------

		file, err := os.Create("fransFaaseVerification/" + tmString + ".txt")
		if err != nil {
			fmt.Println("File does not exists or cannot be created")
			os.Exit(1)
		}
		defer file.Close()

		file.WriteString(tmString + "\n\n")

		stringSlice := []string{}
		for configString, print := range stringsToPrint {
			if print {
				stringSlice = append(stringSlice, configString)
			}
		}
		sort.Slice(stringSlice, func(i, j int) bool {
			return stringSlice[i] < stringSlice[j]
		})
		for _, configString := range stringSlice {
			file.WriteString(configString + "\n")
		}

	}
}
