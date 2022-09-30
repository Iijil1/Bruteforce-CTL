package prooffileproduction

import (
	"fmt"
	"os"

	"github.com/Iijil1/Bruteforce-CTL/decider"
	TM "github.com/Iijil1/Bruteforce-CTL/turingmachines"
)

func CreateProofFile(tm TM.TuringMachine, options decider.DeciderOptions) {
	if decided, info := decider.BruteforceCTL(tm, options, decider.SilentPrintOptions()); decided {
		tmString := tm.NameString()
		start := info.StartConfig
		printStack := []*decider.Configuration{start}
		i := 1
		start.PrintHelper = i

		file, err := os.Create("prooffiles/" + tmString + ".txt")
		if err != nil {
			fmt.Println("File does not exists or cannot be created")
			os.Exit(1)
		}
		defer file.Close()

		file.WriteString(tmString + "\n\n")

		for len(printStack) > 0 {
			config := printStack[0]
			printStack = printStack[1:]
			file.WriteString(fmt.Sprintf("%v:\t%v\n -> ", config.PrintHelper, config))
			for j, successor := range config.Successors {
				if successor.Status == decider.DONE {
					if successor.PrintHelper == 0 {
						i += 1
						successor.PrintHelper = i
						printStack = append(printStack, successor)
					}
					if j != 0 {
						file.WriteString(",")
					}
					file.WriteString(fmt.Sprint(successor.PrintHelper))
				}
			}
			file.WriteString("\n\n")
		}
	}
}
