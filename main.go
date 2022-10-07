package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/Iijil1/Bruteforce-CTL/decider"
	TM "github.com/Iijil1/Bruteforce-CTL/turingmachines"
)

func main() {
	var err error

	argMaxMachines := flag.Int("m", 0, "maximum number of machines tested. 0 to test entire index file.")
	argIndexFile := flag.String("f", "remaining_index", "undecided index file to use. Empty to test entire database")
	argOptions := flag.String("o", "", "the option preset to use.")

	flag.Parse()

	maxMachines := *argMaxMachines
	indexFileName := *argIndexFile

	var options decider.DeciderOptions
	switch *argOptions {
	case "bouncer":
		options = decider.BouncerSearch()
	case "bruteforce":
		options = decider.AggressiveAbstraction()
	case "counter":
		options = decider.CounterSearch()
	case "forcecounter":
		options = decider.ForceCounterSearch()
	case "stitched":
		options = decider.StitchedBouncerSearch()
	case "nonbinary":
		options = decider.NonBinaryCounters()
	case "buffer":
		options = decider.Buffer()
	case "counterTwo":
		options = decider.CounterSize2Special()
	case "counterThree":
		options = decider.CounterSize(3)
	case "counterFour":
		options = decider.CounterSize(4)
	case "counterFive":
		options = decider.CounterSize(5)
	default:
		fmt.Println("no options specified")
		os.Exit(-1)
	}

	runName := "ctl-" + fmt.Sprint(time.Now().Unix())
	runFileName := "output/" + runName + "-run"
	f, _ := os.OpenFile(runFileName,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	var undecidedIndex []byte
	if indexFileName != "" {
		undecidedIndex, err = ioutil.ReadFile(indexFileName)

		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}

	numMachines := TM.DBENTRIES
	if maxMachines != 0 && maxMachines < numMachines {
		numMachines = maxMachines
	}
	if indexFileName != "" && len(undecidedIndex)/4 < numMachines {
		numMachines = len(undecidedIndex) / 4
	}

	decided := 0
	maxDepth := 0
	var maxDepthIndex uint32
	maxSteps := 0
	var maxStepsIndex uint32

	startTime := time.Now()
	for i := 0; i < numMachines; i += 1 {
		index := uint32(i)
		if indexFileName != "" {
			index = binary.BigEndian.Uint32(undecidedIndex[i*4 : (i+1)*4])
		}
		tm := TM.GetTMFromDBIndex(TM.DBFILEPATH, index)
		if result, info := decider.BruteforceCTL(tm, options, decider.SilentPrintOptions()); result {
			decided += 1
			if info.MaxDepth > maxDepth {
				maxDepth = info.MaxDepth
				maxDepthIndex = index
			}
			if info.Steps > maxSteps {
				maxSteps = info.Steps
				maxStepsIndex = index
			}
			var arr [4]byte
			binary.BigEndian.PutUint32(arr[0:4], index)
			f.Write(arr[:])
		}

		if i%25 == 24 {
			fmt.Println(time.Since(startTime), "\tDone:", i+1, "of", numMachines, "=", 100*(i+1)/numMachines, "%\tSuccess rate:", 100*decided/(i+1), "%")
		}
	}
	f.Close()
	duration := time.Since(startTime)

	infoFileName := "output/" + runName + "-info.txt"

	file, err := os.Create(infoFileName)
	if err != nil {
		os.Exit(-1)
	}
	defer file.Close()

	file.WriteString(runName + "\n\n")

	file.WriteString("IndexFile: " + indexFileName + "\n")
	file.WriteString("Duration: " + fmt.Sprint(duration) + "\n")
	file.WriteString("Machines: " + fmt.Sprint(numMachines) + "\n")
	file.WriteString("Decided: " + fmt.Sprint(decided) + "\n\n")

	file.WriteString("Options:\n" + options.IndentedString("\t") + "\n\n")

	file.WriteString("MaxSteps: " + fmt.Sprint(maxSteps) + "\n")
	file.WriteString("MaxStepsIndex: " + fmt.Sprint(maxStepsIndex) + "\n")

	file.WriteString("MaxDepth: " + fmt.Sprint(maxDepth) + "\n")
	file.WriteString("MaxDepthIndex: " + fmt.Sprint(maxDepthIndex) + "\n")
}
