package day6

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	"strelox.com/ac2024/utils"
	myutils "strelox.com/ac2024/utils"
)

func findStart(s string, y int, result chan []int, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for x, ch := range s {
			if len(result) != 0 {
				return
			}
			if ch == '^' {
				result <- []int{x, y}
				return
			}
		}
	}()
}

func readInput(filename string) ([][]byte, int, int) {
	file, ferr := os.Open(filename)
	if ferr != nil {
		log.Fatal(ferr)
		os.Exit(1)
	}

	defer file.Close()

	wg := &sync.WaitGroup{}

	y := 0
	result := [][]byte{}
	scanner := bufio.NewScanner(file)
	startResult := make(chan []int, 1)
	var start []int = nil
	for scanner.Scan() {
		s := scanner.Text()
		if len(startResult) == 0 {
			findStart(s, y, startResult, wg)
		}
		result = append(result, []byte(s))
		y++
	}

	wg.Wait()

	select {
	case start = <-startResult:
	default:
		panic("start haven't been found")
	}

	return result, start[0], start[1]
}

func makeBoundsChecker(szx, szy int) func(x, y int) bool {
	return func(x, y int) bool {
		return x < 0 || x >= szx || y < 0 || y >= szy
	}
}

const (
	STATE_UP = iota
	STATE_RIGHT
	STATE_DOWN
	STATE_LEFT
)

type pathElement struct {
	x, y int
}

type fullPathElement struct {
	pathElement
	state int
}

type fullPath struct {
	elements []fullPathElement
}

type path struct {
	elements []pathElement
}

func outputPole(pole [][]byte) {
	file, _ := os.Create("output.txt")
	writer := bufio.NewWriter(file)
	for y := 0; y < len(pole); y++ {
		writer.Write(pole[y])
		writer.WriteByte('\n')
	}
	writer.Flush()
	defer file.Close()
}

type trackerInput struct {
	input [][]byte
}

type totalCalc struct {
	trackerInput
	total    int
	path     path
	pathMemo map[pathElement]bool
}

type loopDetector struct {
	trackerInput
	isLoopDetected bool
	pathMemo       map[fullPathElement]int
	order          int
}

type fieldActions interface {
	getInput() [][]byte
	actSpace(x, y, status int) bool
}

func newTotalCalc(input [][]byte) *totalCalc {
	tc := &totalCalc{
		trackerInput: trackerInput{input: input},
		total:        0,
		pathMemo:     map[pathElement]bool{}}
	return tc
}

func newLoopDetector(input [][]byte) *loopDetector {
	ld := &loopDetector{
		isLoopDetected: false,
		pathMemo:       map[fullPathElement]int{},
		trackerInput:   trackerInput{input: input},
		order:          0}

	return ld
}

func (ti *trackerInput) getInput() [][]byte {
	return ti.input
}

func (tc *totalCalc) actSpace(x, y, state int) bool {
	pathElemenent := pathElement{x, y}
	_, ok := tc.pathMemo[pathElemenent]
	if !ok {
		tc.total++
		tc.pathMemo[pathElemenent] = true
		tc.path.elements = append(tc.path.elements, pathElement{x, y})
	}

	return false
}

func (ld *loopDetector) actSpace(x, y, state int) bool {
	pathElement := fullPathElement{pathElement{x, y}, state}
	_, ok := ld.pathMemo[pathElement]
	if !ok {
		ld.pathMemo[pathElement] = ld.order
		ld.order++
		return false
	}

	ld.isLoopDetected = true

	return true
}

func makeTrackPathFunc(xc int, yc int, boundsChecker func(x, y int) bool) func(fa fieldActions) {
	return func(fa fieldActions) {
		input := fa.getInput()
		x, y := xc, yc
		xp, yp := xc, yc-1
		state := STATE_UP
		for !boundsChecker(x, y) {
			ch := input[y][x]

			if ch == '#' {
				state = (state + 1) % 4
				x, y = xp, yp
			}

			if ch == '.' {
				shouldStopTracking := fa.actSpace(x, y, state)
				if shouldStopTracking {
					break
				}
			}

			xp, yp = x, y
			switch state {
			case STATE_UP:
				y--
			case STATE_DOWN:
				y++
			case STATE_LEFT:
				x--
			case STATE_RIGHT:
				x++
			default:
				log.Fatal("Wrong flow")
				os.Exit(1)
			}
		}
	}
}

func copyInput(input [][]byte) [][]byte {
	newInput := make([][]byte, len(input))
	for i := range input {
		newInput[i] = make([]byte, len(input[0]))
		copy(newInput[i], input[i])
	}

	return newInput
}

func tryMakeLoops(input [][]byte, path path, trackPathFunc func(fa fieldActions)) int {
	totalLoops := 0
	for _, pe := range path.elements {
		newInput := copyInput(input)
		newInput[pe.y][pe.x] = '#'
		loopDetector := newLoopDetector(newInput)
		trackPathFunc(loopDetector)
		if loopDetector.isLoopDetected {
			totalLoops++
		}
	}

	return totalLoops
}

func tryMakeLoopsSem(input [][]byte, path path, trackPathFunc func(fa fieldActions)) int {
	fmt.Printf("runtime.NumCPU()=%d\n", runtime.NumCPU())
	sem := utils.NewSemaphore(runtime.NumCPU() * 4) // 4 gives very good result somehow
	totalLoops := 0
	result := make(chan bool, len(path.elements))
	wg := &sync.WaitGroup{}
	for _, pe := range path.elements {
		wg.Add(1)
		sem.Acquire()
		go func(pe pathElement, input [][]byte, result chan bool, sem *myutils.Semaphore, wg *sync.WaitGroup) {
			defer wg.Done()
			newInput := copyInput(input)
			newInput[pe.y][pe.x] = '#'
			loopDetector := newLoopDetector(newInput)
			trackPathFunc(loopDetector)
			result <- loopDetector.isLoopDetected
			sem.Release()
		}(pe, input, result, sem, wg)
	}

	wg.Wait()
	close(result)

	for isLoopDetected := range result {
		if isLoopDetected {
			totalLoops++
		}
	}

	sem.Close()

	return totalLoops
}

func Execute() {
	input, xc, yc := readInput("input.txt")

	boundsChecker := makeBoundsChecker(len(input[0]), len(input))

	fmt.Printf("%dx%d %d %d\n", len(input[0]), len(input), xc, yc)
	input[yc][xc] = '.'

	trackPathFunc := makeTrackPathFunc(xc, yc, boundsChecker)
	totalCalc := newTotalCalc(input)
	trackPathFunc(totalCalc)
	fmt.Printf("total=%d\n", totalCalc.total)

	timeStart := time.Now()
	totalLoops := tryMakeLoops(input, totalCalc.path, trackPathFunc)
	myutils.TimeTrack(timeStart, "tryMakeLoops")
	fmt.Printf("1 totalLoops=%d\n", totalLoops)

	timeStart2 := time.Now()
	totalLoops2 := tryMakeLoopsSem(input, totalCalc.path, trackPathFunc)
	myutils.TimeTrack(timeStart2, "tryMakeLoopsSem")
	fmt.Printf("2 totalLoops=%d\n", totalLoops2)
}
