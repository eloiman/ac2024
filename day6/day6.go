package day6

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func grFindStart(s string, y int, output chan []int, status chan int) {
	guardSeacher := func() {
		for x, ch := range s {
			select {
			case st := <-status:
				status <- st
				return
			default:
				if ch == '^' {
					output <- []int{x, y}
					status <- '1'
					return
				}
			}
		}
	}

	guardSeacher()
}

func readInput(filename string) ([][]byte, int, int) {
	file, ferr := os.Open(filename)
	if ferr != nil {
		log.Fatal(ferr)
		os.Exit(1)
	}

	defer file.Close()

	var res []int = nil
	y := 0
	result := [][]byte{}
	scanner := bufio.NewScanner(file)
	outchan := make(chan []int, 1)
	status := make(chan int, 0)
	for scanner.Scan() {
		s := scanner.Text()
		if res == nil {
			select {
			case st := <-status:
				res = <-outchan
				status <- st
			default:
				go grFindStart(s, y, outchan, status)
			}
		}
		result = append(result, []byte(s))
		y++
	}

	return result, res[0], res[1]
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

func Execute() {
	input, xc, yc := readInput("input.txt")

	boundsChecker := makeBoundsChecker(len(input[0]), len(input))

	fmt.Printf("%dx%d %d %d\n", len(input[0]), len(input), xc, yc)
	input[yc][xc] = '.'

	x, y := xc, yc
	xp, yp := xc, yc-1
	state := STATE_UP
	total := 0
	for !boundsChecker(x, y) {
		ch0 := input[y][x]
		if ch0 == '#' {
			state = (state + 1) % 4
			x, y = xp, yp
		}
		ch := &input[y][x]
		if *ch == '.' {
			*ch = 'X'
			total++
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

	outputPole(input)
	fmt.Printf("total=%d\n", total)
}
