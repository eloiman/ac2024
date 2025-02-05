package day6

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"strelox.com/ac2024/utils"
)

func grFindStart(s string, y int, output chan []int, status chan int) {
	guardSeacher := func() {
		for x, ch := range s {
			select {
			case st := <-status:
				utils.Unused(st)
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
	outchan := make(chan []int, 2)
	status := make(chan int, 2)
	for scanner.Scan() {
		s := scanner.Text()
		if res == nil {
			select {
			case st := <-status:
				utils.Unused(st)
				res = <-outchan
			default:
				go grFindStart(s, y, outchan, status)
			}
		}
		result = append(result, []byte(s))
		y++
	}

	return result, res[0], res[1]
}

func Execute() {
	input, x, y := readInput("input.txt")

	fmt.Printf("%dx%d %d %d", len(input[0]), len(input), x, y)
}
