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

func isLoopExists(path Path) bool {
	p1 := 0
	p2 := 0
	for p1 = 1; p1 < len(path.elements); p1++ {
		if p1 == p2 {
			return true
		}
		if p1&1 != 0 {
			p2 += 1
		}
	}

	return false
}

func getPath(xc int, yc int, input [][]byte, boundsChecker func(x, y int) bool) (int, Path) {
	result := Path{}
	total := 0
	x, y := xc, yc
	xp, yp := xc, yc-1
	state := STATE_UP
	path := make(map[PathElement]bool)
	for !boundsChecker(x, y) {
		ch0 := input[y][x]
		if ch0 == '#' {
			state = (state + 1) % 4
			x, y = xp, yp
		}
		ch := &input[y][x]
		if *ch == '.' {
			pathElemenent := PathElement{x, y}
			_, ok := path[pathElemenent]
			if !ok {
				total++
			}
			path[pathElemenent] = true
			result.elements = append(result.elements, FullPathElement{PathElement{x, y}, state})
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

	return total, result
}

func Execute() {
	input, xc, yc := readInput("input.txt")

	boundsChecker := makeBoundsChecker(len(input[0]), len(input))

	fmt.Printf("%dx%d %d %d\n", len(input[0]), len(input), xc, yc)
	input[yc][xc] = '.'

	//outputPole(input)
	total, path := getPath(xc, yc, input, boundsChecker)
	fmt.Printf("total=%d\n", total)

	loopExists := isLoopExists(path)
	fmt.Printf("loop=%t\n", loopExists)
}
