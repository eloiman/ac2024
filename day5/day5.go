package day5

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"sort"
	"strings"
)

type PageOrder struct {
	data map[int][]int
}

type ManualPages struct {
	data [][]int
}

func readPageOrder(scanner *bufio.Scanner) (PageOrder, error) {
	var result PageOrder

	result.data = make(map[int][]int)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		lineReader := strings.NewReader(line)
		var x, y int
		n, err := fmt.Fscanf(lineReader, "%d|%d", &x, &y)
		if n != 2 || err != nil {
			log.Fatal(err)
			newError := errors.New("readPageOrder failed to read input")
			return PageOrder{}, errors.Join(newError, err)
		}

		result.data[y] = append(result.data[y], x)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		newError := errors.New("readPageOrder failed : scanner failed")
		return PageOrder{}, errors.Join(newError, err)
	}

	keys := make([]int, 0, len(result.data))
	for k := range result.data {
		keys = append(keys, k)
	}

	for i := 0; i < len(keys); i++ {
		sort.Ints(result.data[keys[i]])
	}

	return result, nil
}

func readManualPages(scanner *bufio.Scanner) (ManualPages, error) {
	var result ManualPages

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		lineReader := strings.NewReader(line)
		var pages []int
		var x int
		n, err := fmt.Fscanf(lineReader, "%d", &x)
		if n != 1 || err != nil {
			break
		}
		pages = append(pages, x)
		for {
			n, err = fmt.Fscanf(lineReader, ",%d", &x)
			if n != 1 || err != nil {
				copy([]int{}, pages)
				break
			}

			pages = append(pages, x)
		}
		result.data = append(result.data, pages)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		newError := errors.New("readPageOrder failed : scanner failed")
		return ManualPages{}, errors.Join(newError, err)
	}

	return result, nil
}

func readInput(filename string) (PageOrder, ManualPages) {
	file, ferr := os.Open(filename)
	if ferr != nil {
		log.Fatal(ferr)
		os.Exit(1)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var err error
	pageOrder, err := readPageOrder(scanner)
	if err != nil {
		log.Fatal(ferr)
		os.Exit(1)
	}

	manualPages, err := readManualPages(scanner)
	if err != nil {
		log.Fatal(ferr)
		os.Exit(1)
	}

	return pageOrder, manualPages
}

func checkPagesOrder(index int, pageOrder *PageOrder, pages []int) bool {
	for k, p := range pages {
		for j := k; j != 0; j-- {
			_, ok := slices.BinarySearch(pageOrder.data[p], pages[j-1])
			if !ok {
				return false
			}
		}
	}

	return true
}

func calcAnswer(pageOrder *PageOrder, manualPages *ManualPages) int {
	var summ int = 0
	for pindex, pages := range manualPages.data {
		ok := checkPagesOrder(pindex, pageOrder, pages)
		if ok {
			summ += pages[len(pages)/2]
		}
	}

	return summ
}

func Execute() {
	pageOrder, manualPages := readInput("input.txt")

	summ := calcAnswer(&pageOrder, &manualPages)
	fmt.Printf("%d", summ)
}
