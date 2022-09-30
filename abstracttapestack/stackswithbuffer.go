package abstracttapestack

import (
	"math"

	TM "github.com/Iijil1/Bruteforce-CTL/turingmachines"
)

type stackWithBuffer struct {
	list               *stackList
	bufferEncoding     int
	bufferLength       int
	stackWithoutBuffer AbstractStack
}

func (swb stackWithBuffer) Push(symbol TM.Symbol) []AbstractStack {
	maxBufferLength := swb.list.abstractions.Buffer
	base := swb.list.symbolsInUse.NumSymbols()
	if swb.bufferLength < maxBufferLength {
		nextStack := stackWithBuffer{
			list:               swb.list,
			bufferEncoding:     swb.bufferEncoding*base + symbol.GetInt(),
			bufferLength:       swb.bufferLength + 1,
			stackWithoutBuffer: swb.stackWithoutBuffer,
		}
		return []AbstractStack{nextStack}
	}
	extractionHelper := int(math.Pow(float64(base), float64(maxBufferLength-1)))
	lastSymbolInt := swb.bufferEncoding / extractionHelper
	lastSymbol := swb.list.symbolsInUse.GetSymbol(lastSymbolInt)
	nextEncoding := swb.bufferEncoding % extractionHelper
	nextEncoding = nextEncoding*base + symbol.GetInt()
	possibleStacksWithoutBuffer := swb.stackWithoutBuffer.Push(lastSymbol)
	possibleStacks := make([]AbstractStack, len(possibleStacksWithoutBuffer))
	for i := 0; i < len(possibleStacks); i += 1 {
		possibleStacks[i] = stackWithBuffer{
			list:               swb.list,
			bufferEncoding:     nextEncoding,
			bufferLength:       maxBufferLength,
			stackWithoutBuffer: possibleStacksWithoutBuffer[i],
		}
	}
	return possibleStacks

}
func (swb stackWithBuffer) Pop() ([]TM.Symbol, []AbstractStack) {
	if swb.bufferLength == 0 {
		possibleSymbols, possibleStacksWithoutBuffer := swb.stackWithoutBuffer.Pop()
		possibleStacks := make([]AbstractStack, len(possibleStacksWithoutBuffer))
		for i := 0; i < len(possibleStacks); i += 1 {
			possibleStacks[i] = stackWithBuffer{
				list:               swb.list,
				bufferEncoding:     0,
				bufferLength:       0,
				stackWithoutBuffer: possibleStacksWithoutBuffer[i],
			}
		}
		return possibleSymbols, possibleStacks
	} else {
		base := swb.list.symbolsInUse.NumSymbols()
		nextSymbol := swb.list.symbolsInUse.GetSymbol(swb.bufferEncoding % base)
		nextStack := stackWithBuffer{
			list:               swb.list,
			bufferEncoding:     swb.bufferEncoding / base,
			bufferLength:       swb.bufferLength - 1,
			stackWithoutBuffer: swb.stackWithoutBuffer,
		}
		return []TM.Symbol{nextSymbol}, []AbstractStack{nextStack}
	}
}
func (swb stackWithBuffer) String() string {
	return swb.NeatString(false)
}
func (swb stackWithBuffer) NeatString(reverse bool) string {
	base := swb.list.symbolsInUse.NumSymbols()
	bufferString := ""
	restEncoding := swb.bufferEncoding
	for i := 0; i < swb.bufferLength; i += 1 {
		nextSymbol := swb.list.symbolsInUse.GetSymbol(restEncoding % base).String()
		if reverse {
			bufferString = bufferString + nextSymbol
		} else {
			bufferString = nextSymbol + bufferString
		}
		restEncoding = restEncoding / base
	}

	if reverse {
		return bufferString + swb.stackWithoutBuffer.NeatString(reverse)
	} else {
		return swb.stackWithoutBuffer.NeatString(reverse) + bufferString
	}
}
func (swb stackWithBuffer) Length() int {
	return swb.bufferLength + swb.stackWithoutBuffer.Length()
}
