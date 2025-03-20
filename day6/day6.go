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

type PathElement struct {
	x, y int
}

type FullPathElement struct {
	PathElement
	state int
}

type FullPath struct {
	elements []FullPathElement
}

type Path struct {
	elements []PathElement
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

type TrackerInput struct {
	input [][]byte
}

type TotalCalc struct {
	TrackerInput
	total    int
	path     Path
	pathMemo map[PathElement]bool
}

type LoopDetector struct {
	TrackerInput
	isLoopDetected bool
	pathMemo       map[FullPathElement]int
	order          int
}

type FieldActions interface {
	GetInput() [][]byte
	ActSpace(x, y, status int) bool
}

func NewTotalCalc(input [][]byte) *TotalCalc {
	tc := &TotalCalc{
		TrackerInput: TrackerInput{input: input},
		total:        0,
		pathMemo:     map[PathElement]bool{}}
	return tc
}

func NewLoopDetector(input [][]byte) *LoopDetector {
	ld := &LoopDetector{
		isLoopDetected: false,
		pathMemo:       map[FullPathElement]int{},
		TrackerInput:   TrackerInput{input: input},
		order:          0}

	return ld
}

func (ti *TrackerInput) GetInput() [][]byte {
	return ti.input
}

func (tc *TotalCalc) ActSpace(x, y, state int) bool {
	pathElemenent := PathElement{x, y}
	_, ok := tc.pathMemo[pathElemenent]
	if !ok {
		tc.total++
		tc.pathMemo[pathElemenent] = true
		tc.path.elements = append(tc.path.elements, PathElement{x, y})
	}

	return false
}

func (ld *LoopDetector) ActSpace(x, y, state int) bool {
	pathElement := FullPathElement{PathElement{x, y}, state}
	_, ok := ld.pathMemo[pathElement]
	if !ok {
		ld.pathMemo[pathElement] = ld.order
		ld.order++
		return false
	}

	ld.isLoopDetected = true

	return true
}

func makeTrackPathFunc(xc int, yc int, boundsChecker func(x, y int) bool) func(fa FieldActions) {
	return func(fa FieldActions) {
		input := fa.GetInput()
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
				shouldStopTracking := fa.ActSpace(x, y, state)
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

func tryMakeLoops(input [][]byte, path Path, trackPathFunc func(fa FieldActions)) int {
	totalLoops := 0
	for _, pe := range path.elements {
		newInput := copyInput(input)
		newInput[pe.y][pe.x] = '#'
		loopDetector := NewLoopDetector(newInput)
		trackPathFunc(loopDetector)
		if loopDetector.isLoopDetected {
			totalLoops++
		}
	}

	return totalLoops
}

func tryMakeLoopsSem(input [][]byte, path Path, trackPathFunc func(fa FieldActions)) int {
	fmt.Printf("runtime.NumCPU()=%d\n", runtime.NumCPU())
	sem := utils.NewSemaphore(runtime.NumCPU() * 4) // 4 gives very good result somehow
	totalLoops := 0
	result := make(chan bool, len(path.elements))
	wg := &sync.WaitGroup{}
	for _, pe := range path.elements {
		wg.Add(1)
		sem.Acquire()
		go func(pe PathElement, input [][]byte, result chan bool, sem *myutils.Semaphore, wg *sync.WaitGroup) {
			defer wg.Done()
			newInput := copyInput(input)
			newInput[pe.y][pe.x] = '#'
			loopDetector := NewLoopDetector(newInput)
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
	totalCalc := NewTotalCalc(input)
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
