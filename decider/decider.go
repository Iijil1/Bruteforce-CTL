package decider

import (
	"fmt"

	ATS "github.com/Iijil1/Bruteforce-CTL/abstracttapestack"
	TM "github.com/Iijil1/Bruteforce-CTL/turingmachines"
)

type Status int

const TODO Status = 0
const DONE Status = 1
const HALTSATANYDEPTH Status = 2
const ABORTED Status = 3
const HALTSBECAUSEDEPTHLIMIT Status = 4

type Configuration struct {
	State                 TM.State
	Direction             TM.Direction
	StackList             ATS.StackCollection
	TM                    TM.TuringMachine
	LeftTape              ATS.AbstractStack
	RightTape             ATS.AbstractStack
	Status                Status
	Predecessors          []*Configuration
	Successors            []*Configuration
	PrintHelper           int
	Depth                 int
	DepthLimitCausingHalt int
}

func (c *Configuration) String() string {
	stateString := c.State.String()
	if c.Direction == TM.RIGHT {
		stateString = stateString + ">"
	} else {
		stateString = "<" + stateString
	}
	return c.LeftTape.NeatString(false) + stateString + c.RightTape.NeatString(true)
}

type ConfigurationKey struct {
	State     TM.State
	Direction TM.Direction
	LeftTape  ATS.AbstractStack
	RightTape ATS.AbstractStack
}

type RunInfo struct {
	Steps            int
	ConfigurationMap map[ConfigurationKey]*Configuration
	StartConfig      *Configuration
	MaxDepth         int
}

func BruteforceCTL(tm TM.TuringMachine, options DeciderOptions, printOptions PrintOptions) (bool, RunInfo) {
	stackList := ATS.NewStackCollection(tm, options.StackHeuristics)
	startConfiguration := &Configuration{State: tm.StartState(), TM: tm, Direction: TM.RIGHT, StackList: stackList, LeftTape: stackList.EmptyStack(), RightTape: stackList.EmptyStack(), Status: TODO, Depth: 0}
	todoStack := []*Configuration{startConfiguration}
	redoStack := []*Configuration{}
	forcedLines := []*Configuration{}
	configurationMap := make(map[ConfigurationKey]*Configuration)
	configurationMap[ConfigurationKey{tm.StartState(), TM.RIGHT, stackList.EmptyStack(), stackList.EmptyStack()}] = startConfiguration
	depthlimit := options.InitialDepth
	maxDepth := 0
	// specialSights := 0

	for steps := 0; steps < options.StepLimit; steps += 1 {
		var currentConfiguration *Configuration
		if len(redoStack) > 0 {
			currentConfiguration = redoStack[len(redoStack)-1]
			redoStack = redoStack[:len(redoStack)-1]
		} else if len(forcedLines) > 0 {
			currentConfiguration = forcedLines[0]
			forcedLines = forcedLines[1:]
		} else if len(todoStack) > 0 {
			currentConfiguration = todoStack[0]
			todoStack = todoStack[1:]
		} else {
			if printOptions.Success {
				fmt.Printf("After %v steps with maximum depth %v, we found a closed set of configurations. This machine never halts.\n", steps, maxDepth)
				if printOptions.StatusCountInfo {
					todo, done, halts, aborted, haltsDepth := statusCounting(configurationMap)
					fmt.Printf("TODO: %v, DONE: %v, HALTS: %v, ABORTED: %v, HALTSBECAUSEDEPTH: %v\n", todo, done, halts, aborted, haltsDepth)
				}
			}
			return true, RunInfo{Steps: steps, ConfigurationMap: configurationMap, StartConfig: startConfiguration, MaxDepth: maxDepth}
		}

		if currentConfiguration.Status == DONE {
			steps -= 1
			continue
		}
		if currentConfiguration.Status == HALTSATANYDEPTH || currentConfiguration.Status == ABORTED || (currentConfiguration.Status == HALTSBECAUSEDEPTHLIMIT && currentConfiguration.DepthLimitCausingHalt == depthlimit) {
			steps -= 1
			for _, successor := range currentConfiguration.Successors {
				if (successor.Status == DONE || successor.Status == TODO) && successor != startConfiguration {
					abort := true
					for _, predessor := range successor.Predecessors {
						if predessor.Status == DONE {
							abort = false
							break
						}
					}
					if abort {
						successor.Status = ABORTED
						redoStack = append(redoStack, successor)
						if printOptions.SettingAborted {
							fmt.Printf("\t\tIt is now unnecessary to check %v as it has no predecessors among DONE configurations\n", successor)
						}

					}
				}
			}
			continue
		}
		if printOptions.CurrentConfiguration {
			fmt.Printf("Checking %v\n", currentConfiguration)
		}
		currentConfiguration.Status = DONE
		forcedLine := false
		if options.ForcedLines {
			forcedLine = true
		}
		var possibleSymbols []TM.Symbol
		var popStackOptions []ATS.AbstractStack
		if currentConfiguration.Direction == TM.RIGHT {
			possibleSymbols, popStackOptions = currentConfiguration.RightTape.Pop()
		} else {
			possibleSymbols, popStackOptions = currentConfiguration.LeftTape.Pop()
		}
		if len(possibleSymbols) > 1 {
			forcedLine = false
		}
		possibleSuccessors := make([]ConfigurationKey, len(possibleSymbols))
		successorPriority := make([]int, len(possibleSymbols))
		overallFailure := false
		overallFailureBecauseOfDepthLimit := true
		for i := range possibleSymbols {
			popOptionFailure := false
			popOptionFailureBecauseOfDepthLimit := false
			var nextRightStack, nextLeftStack ATS.AbstractStack
			if currentConfiguration.Direction == TM.RIGHT {
				nextRightStack = popStackOptions[i]
				nextLeftStack = currentConfiguration.LeftTape
			} else {
				nextRightStack = currentConfiguration.RightTape
				nextLeftStack = popStackOptions[i]
			}
			if printOptions.DetailedHaltDecision {
				fmt.Printf("\tOption %v%v%v%v", nextLeftStack.NeatString(false), currentConfiguration.State, possibleSymbols[i], nextRightStack.NeatString(true))
			}
			nextSymbol, nextDirection, nextState := tm.Transition(currentConfiguration.State, possibleSymbols[i])
			if nextState.IsHalt() {
				if printOptions.DetailedHaltDecision {
					fmt.Printf(" halts immediatly\n")
				}
				popOptionFailure = true
			} else if nextState.IsLoop() {
				if printOptions.DetailedHaltDecision {
					possibleSuccessors[i] = ConfigurationKey{}
					fmt.Printf(" loops in place forever\n")
				}
			} else {
				if printOptions.DetailedHaltDecision {
					fmt.Printf(" can turn into:\n")
				}
				var pushStackOptions []ATS.AbstractStack
				if nextDirection == TM.RIGHT {
					pushStackOptions = nextLeftStack.Push(nextSymbol)
				} else {
					pushStackOptions = nextRightStack.Push(nextSymbol)
				}

				if len(pushStackOptions) > 1 {
					forcedLine = false
				}
				popOptionFailure = true
				for _, stack := range pushStackOptions {
					if nextDirection == TM.RIGHT {
						nextLeftStack = stack
					} else {
						nextRightStack = stack
					}
					if printOptions.DetailedHaltDecision {
						stateString := nextState.String()
						if nextDirection == TM.RIGHT {
							stateString = stateString + ">"
						} else {
							stateString = "<" + stateString
						}
						fmt.Printf("\t\tOption %v%v%v", nextLeftStack.NeatString(false), stateString, nextRightStack.NeatString(true))
					}
					successorKey := ConfigurationKey{nextState, nextDirection, nextLeftStack, nextRightStack}
					if successor, exists := configurationMap[ConfigurationKey{nextState, nextDirection, nextLeftStack, nextRightStack}]; !exists {
						if printOptions.DetailedHaltDecision {
							fmt.Printf(" which can work because it was never seen before\n")
						}
						if currentConfiguration.Depth < depthlimit {
							if successorPriority[i] < 2 {
								possibleSuccessors[i] = successorKey
								successorPriority[i] = 2
							}
							popOptionFailure = false
							continue
						} else {
							if printOptions.DetailedHaltDecision {
								fmt.Printf("\t\tBut the depth is too big, so we will not consider that yet\n")
							}
							popOptionFailureBecauseOfDepthLimit = true
						}
					} else {
						if printOptions.DetailedHaltDecision {
							fmt.Printf(" which was seen before")
						}
						if successor.Status != HALTSATANYDEPTH && successor.Status != HALTSBECAUSEDEPTHLIMIT {
							if printOptions.DetailedHaltDecision {
								fmt.Printf(" but isn't known to halt yet, so we'll take it\n")
							}
							if successorPriority[i] < 3 {
								possibleSuccessors[i] = successorKey
								successorPriority[i] = 3
							}
							popOptionFailure = false
							break
						}
						if successor.Status == HALTSBECAUSEDEPTHLIMIT {
							if printOptions.DetailedHaltDecision {
								fmt.Printf(" but was found to halt because of the depthlimit\n")
							}
							if successor.DepthLimitCausingHalt < depthlimit {
								if printOptions.DetailedHaltDecision {
									fmt.Printf("\t\tHowever, that was at a lower depthlimit, so we might try it again\n")
								}
								if successorPriority[i] < 1 ||
									(successorPriority[i] == 1 && successor.DepthLimitCausingHalt < configurationMap[possibleSuccessors[i]].DepthLimitCausingHalt) {
									possibleSuccessors[i] = successorKey
									successorPriority[i] = 1
								}
								popOptionFailure = false
								continue
							}
							popOptionFailureBecauseOfDepthLimit = true
						}
						if printOptions.DetailedHaltDecision {
							fmt.Printf(" and been found to halt\n")
						}
					}
				}
			}

			if popOptionFailure {
				overallFailure = true
				if !popOptionFailureBecauseOfDepthLimit {
					overallFailureBecauseOfDepthLimit = false
					break
				}
			}

		}

		if overallFailure {
			if overallFailureBecauseOfDepthLimit == false {
				currentConfiguration.Status = HALTSATANYDEPTH
				if currentConfiguration == startConfiguration {
					if printOptions.FoundStartHalting {
						fmt.Printf("\tWe found that the starting configuration has the possibility to halt. We can't prove this machine doesn't halt. That took %v steps and %v configurations.\n", steps, len(configurationMap))
					}
					return false, RunInfo{Steps: steps, StartConfig: startConfiguration, MaxDepth: maxDepth}
				}
				if printOptions.FoundHalting {
					fmt.Println("\tWe found that this configuration has the possibility to halt. We will mark that and try the predecessors again.")
				}
			} else {
				currentConfiguration.Status = HALTSBECAUSEDEPTHLIMIT
				currentConfiguration.DepthLimitCausingHalt = depthlimit
				if currentConfiguration == startConfiguration {
					if options.DepthIncrease > 0 {

						if printOptions.DepthLimitIncrease {
							fmt.Printf("\tWe found that the starting configuration has the possibility to halt with depthlimit %v. That took %v steps. Increasing the limit to %v.\n", depthlimit, steps, depthlimit+options.DepthIncrease)
							if printOptions.StatusCountInfo {
								todo, done, halts, aborted, haltsDepth := statusCounting(configurationMap)
								fmt.Printf("TODO: %v, DONE: %v, HALTS: %v, ABORTED: %v, HALTSBECAUSEDEPTH: %v\n", todo, done, halts, aborted, haltsDepth)
							}
						}
						depthlimit += options.DepthIncrease
						currentConfiguration.Status = TODO
						todoStack = append(todoStack, currentConfiguration)
						continue
					} else {
						if printOptions.FoundStartHalting {
							fmt.Printf("\tWe found that the starting configuration has the possibility to halt with the given depth limit. We can't prove this machine doesn't halt. That took %v steps.\n", steps)
							if printOptions.StatusCountInfo {
								todo, done, halts, aborted, haltsDepth := statusCounting(configurationMap)
								fmt.Printf("TODO: %v, DONE: %v, HALTS: %v, ABORTED: %v, HALTSBECAUSEDEPTH: %v\n", todo, done, halts, aborted, haltsDepth)
							}
						}
						return false, RunInfo{Steps: steps, StartConfig: startConfiguration, MaxDepth: maxDepth}

					}
				}
				if printOptions.FoundHalting {
					fmt.Printf("\tWe found that this configuration has the possibility to halt at depthlimit %v. We will mark that and try the predecessors again.\n", depthlimit)
				}
			}
			for _, predecessor := range currentConfiguration.Predecessors {

				if predecessor.Status == DONE || predecessor.Status == TODO {
					predecessor.Status = TODO
					redoStack = append(redoStack, predecessor)
					if printOptions.FoundHalting && printOptions.AddedConfiguration {
						fmt.Printf("\t\tAdding %v\n", predecessor)
					}
				} else {
					if printOptions.FoundHalting && printOptions.AddedConfiguration {
						fmt.Printf("\t\t%v is a predecessor, but is not added again because of the status %v\n", predecessor, predecessor.Status)
					}
				}
			}

			redoStack = append(redoStack, currentConfiguration)
			continue
		}
		currentConfiguration.Successors = []*Configuration{}
		for _, successorKey := range possibleSuccessors {
			if successorKey == (ConfigurationKey{}) {
				continue
			}
			if successor, exists := configurationMap[successorKey]; !exists {
				nextConfiguration := &Configuration{State: successorKey.State, Direction: successorKey.Direction, StackList: currentConfiguration.StackList, TM: currentConfiguration.TM, LeftTape: successorKey.LeftTape, RightTape: successorKey.RightTape, Status: TODO, Depth: currentConfiguration.Depth + 1, Predecessors: []*Configuration{currentConfiguration}, Successors: []*Configuration{}}
				if forcedLine {
					nextConfiguration.Depth = currentConfiguration.Depth
					forcedLines = append(forcedLines, nextConfiguration)
				} else {
					todoStack = append(todoStack, nextConfiguration)
				}
				configurationMap[successorKey] = nextConfiguration
				currentConfiguration.Successors = append(currentConfiguration.Successors, nextConfiguration)
				if nextConfiguration.Depth > maxDepth {
					maxDepth = nextConfiguration.Depth
				}
				if printOptions.AddedConfiguration {
					fmt.Printf("\tAdding %v\n", nextConfiguration)
				}
			} else {
				if successor.Status == ABORTED || successor.Status == HALTSBECAUSEDEPTHLIMIT {
					successor.Status = TODO
					todoStack = append(todoStack, successor)
					if printOptions.AddedConfiguration {
						fmt.Printf("\tAdding %v, checking it is now necessary again\n", successor)
					}
				} else if printOptions.AddedConfiguration {
					fmt.Printf("\tWe reach %v and it was already seen\n", successor)
				}
				successor.Predecessors = append(successor.Predecessors, currentConfiguration)
				currentConfiguration.Successors = append(currentConfiguration.Successors, successor)
			}
		}

	}
	if printOptions.OutOfSteps {
		fmt.Printf("Even after %v steps, we keep finding new configurations. We can't decide this machine.\n", options.StepLimit)
		if printOptions.StatusCountInfo {
			todo, done, halts, aborted, haltsDepth := statusCounting(configurationMap)
			fmt.Printf("TODO: %v, DONE: %v, HALTS: %v, ABORTED: %v, HALTSBECAUSEDEPTH: %v\n", todo, done, halts, aborted, haltsDepth)
		}
	}
	return false, RunInfo{Steps: options.StepLimit, ConfigurationMap: configurationMap, StartConfig: startConfiguration, MaxDepth: maxDepth}
}

func statusCounting(configurationMap map[ConfigurationKey]*Configuration) (todo int, done int, halts int, aborted int, haltsDepth int) {
	for _, conf := range configurationMap {
		switch conf.Status {
		case TODO:
			todo += 1
		case DONE:
			done += 1
		case HALTSATANYDEPTH:
			halts += 1
		case ABORTED:
			aborted += 1
		case HALTSBECAUSEDEPTHLIMIT:
			haltsDepth += 1
		}
	}
	return
}
