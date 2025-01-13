package day1

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
)

func readInput(filename string) ([]int, []int) {
	file, ferr := os.Open(filename)
	if ferr != nil {
		log.Fatal(ferr)
		os.Exit(1)
	}

	var l1 []int
	var l2 []int

	for true {
		var n1 int
		var n2 int
		_, err := fmt.Fscanf(file, "%d%d\n", &n1, &n2)
		if err != nil {
			break
		}

		l1 = append(l1, n1)
		l2 = append(l2, n2)
	}

	return l1, l2
}

func Execute() {
	l1, l2 := readInput("day1_1.txt")
	if len(l1) != len(l2) || len(l1) != 1000 || len(l2) != 1000 {
		panic(errors.New("wrong input"))
	}

	l1original := l1

	sort.Sort(sort.IntSlice(l1))
	sort.Sort(sort.IntSlice(l2))

	if !sort.IntsAreSorted(l1) || !sort.IntsAreSorted(l2) {
		panic(errors.New("not sorted"))
	}

	var total int = 0
	for i, n1 := range l1 {
		n2 := l2[i]
		if n1 > n2 {
			total += (n1 - n2)
		} else {
			total += (n2 - n1)
		}
	}

	fmt.Printf("total=%d\n", total)

	var l2freq map[int]int
	l2freq = make(map[int]int)

	for _, n2 := range l2 {
		l2freq[n2]++
	}

	var score int = 0
	for _, n1 := range l1original {
		score += n1 * l2freq[n1]
	}

	fmt.Printf("score=%d\n", score)
}
