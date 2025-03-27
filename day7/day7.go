package day6

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	myutils "strelox.com/ac2024/utils"
)

type equation struct {
	target uint64
	values []int
}

func readInput(filename string) []equation {
	fh, ferr := os.Open(filename)
	if ferr != nil {
		log.Fatal("Failed to open input " + ferr.Error())
		panic("ferr")
	}

	defer fh.Close()

	result := []equation{}
	scaner := bufio.NewScanner(fh)
	for scaner.Scan() {
		line := scaner.Text()
		lineReader := strings.NewReader(line)

		var target uint64
		nred, err := fmt.Fscanf(lineReader, "%d:", &target)
		if nred != 1 || err != nil {
			log.Fatal("Failed to read a target")
			panic("Failed to read a target")
		}

		output := equation{target: target}

		for {
			var value int
			nred, err := fmt.Fscanf(lineReader, "%d", &value)
			if nred != 1 || err != nil {
				break
			}

			output.values = append(output.values, value)
		}

		result = append(result, output)
	}

	return result
}

// TODO: optimize with cache, add benchmark tests
func (eq *equation) canBeCalibrated() bool {
	nvalues := len(eq.values)

	var bitset uint32 = 0
	var maxBitset uint32 = 1 << nvalues

	for bitset < maxBitset {
		var evalResult uint64 = uint64(eq.values[0])
		var mask uint32 = 1
		var ivalue int = 1
		for ivalue < nvalues {
			if bitset&mask == 0 {
				evalResult += uint64(eq.values[ivalue])
			} else {
				evalResult *= uint64(eq.values[ivalue])
			}
			mask <<= 1
			ivalue++
		}

		if eq.target == uint64(evalResult) {
			return true
		}

		bitset++
	}

	return false
}

func testCalibration(eq equation, results chan uint64, sem *myutils.Semaphore) {
	defer sem.Release()

	if eq.canBeCalibrated() {
		results <- eq.target
	} else {
		results <- 0
	}
}

func resultsCollector(maxResults int, results <-chan uint64, resultsSumm chan uint64) {
	var summ uint64 = 0
	nresult := 0
	for nresult < maxResults {
		select {
		case eval, ok := <-results:
			if ok {
				summ += eval
				nresult++
			}
		default:
		}
	}

	resultsSumm <- summ
}

func testAllCalibration(eqs []equation) uint64 {
	sem := myutils.NewSemaphore(runtime.NumCPU() * 4)

	results := make(chan uint64, sem.Size()+1)
	defer close(results)
	resultsSumm := make(chan uint64, 1)
	defer close(resultsSumm)

	go resultsCollector(len(eqs), results, resultsSumm)

	for _, eq := range eqs {
		sem.Acquire()
		go testCalibration(eq, results, sem)
	}

	summ := <-resultsSumm

	return summ
}

func Execute() {
	input := readInput("input.txt")

	fmt.Printf("len=%d\n", len(input))

	resultSumm := testAllCalibration(input)
	fmt.Printf("resultSumm=%d\n", resultSumm)
}
