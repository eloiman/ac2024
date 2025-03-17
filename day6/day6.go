package day6

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
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
	order int
}

type Path struct {
	elements []FullPathElement
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

type TotalCalc struct {
	total    int
	path     Path
	pathMemo map[PathElement]bool
	input    [][]byte
}

type LoopDetector struct {
	isLoopDetected bool
	pathSoFar      Path
	input          [][]byte
}

type FieldActions interface {
	Init()
	GetInput() [][]byte
	ActSpace(x, y, status int)
}

func (tc *TotalCalc) Init() {
	tc.pathMemo = map[PathElement]bool{}
	tc.total = 0
}

func (tc *TotalCalc) GetInput() [][]byte {
	return tc.input
}

func (tc *TotalCalc) ActSpace(x, y, state int) {
	pathElemenent := PathElement{x, y}
	_, ok := tc.pathMemo[pathElemenent]
	if !ok {
		tc.total++
	}
	tc.pathMemo[pathElemenent] = true
	tc.path.elements = append(tc.path.elements, FullPathElement{PathElement{x, y}, state, len(tc.path.elements) + 1})
}

func makeTrackPathFunc(xc int, yc int, boundsChecker func(x, y int) bool) func(fa FieldActions) {
	return func(fa FieldActions) {
		input := fa.GetInput()
		x, y := xc, yc
		xp, yp := xc, yc-1
		state := STATE_UP
		for !boundsChecker(x, y) {
			ch0 := input[y][x]
			if ch0 == '#' {
				state = (state + 1) % 4
				x, y = xp, yp
			}
			ch := &input[y][x]
			if *ch == '.' {
				fa.ActSpace(x, y, state)
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

func Execute() {
	input, xc, yc := readInput("input.txt")

	boundsChecker := makeBoundsChecker(len(input[0]), len(input))

	fmt.Printf("%dx%d %d %d\n", len(input[0]), len(input), xc, yc)
	input[yc][xc] = '.'

	trackPathFunc := makeTrackPathFunc(xc, yc, boundsChecker)
	totalCalc := TotalCalc{input: input}
	totalCalc.Init()
	trackPathFunc(&totalCalc)
	fmt.Printf("total=%d\n", totalCalc.total)
}
