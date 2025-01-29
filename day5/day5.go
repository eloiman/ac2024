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
	"sync"
	"time"

	myutils "strelox.com/ac2024/utils"
)

type PageOrder struct {
	data map[int][]int
}

type ManualPages struct {
	data [][]int
}

type IntSlicePair struct {
	key    int
	values []int
}

type BySliceSize []IntSlicePair

func (b BySliceSize) Len() int {
	return len(b)
}

func (b BySliceSize) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b BySliceSize) Less(i, j int) bool {
	return len(b[i].values) < len(b[j].values)
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

func checkPagesOrder(pageOrder *PageOrder, pages []int) bool {
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

func checkPagesOrderGoroutine(wg *sync.WaitGroup, results chan int, pageOrder *PageOrder, pages []int) {
	defer wg.Done()
	ok := checkPagesOrder(pageOrder, pages)
	if ok {
		results <- pages[len(pages)/2]
	}
}

func calcAnswer(pageOrder *PageOrder, manualPages *ManualPages) (int, []int, []int) {
	defer myutils.TimeTrack(time.Now(), "calcAnswer")
	var summ int = 0
	indexes := []int{}
	wrongIndexes := []int{}
	for i, pages := range manualPages.data {
		ok := checkPagesOrder(pageOrder, pages)
		if ok {
			summ += pages[len(pages)/2]
			indexes = append(indexes, i)
		} else {
			wrongIndexes = append(wrongIndexes, i)
		}
	}

	return summ, indexes, wrongIndexes
}

func calcAnswerParallel(pageOrder *PageOrder, manualPages *ManualPages) int {
	defer myutils.TimeTrack(time.Now(), "calcAnswerParallel")
	results := make(chan int, len(manualPages.data))

	wg := &sync.WaitGroup{}
	wg.Add(len(manualPages.data))
	for _, pages := range manualPages.data {
		go checkPagesOrderGoroutine(wg, results, pageOrder, pages)
	}

	var summ int = 0
	wg.Wait()

	results <- -1

	isStopped := false
	n := 0
	for !isStopped {
		select {
		case x, ok := <-results:
			if !ok {
				continue
			}
			if x == -1 {
				isStopped = true
				break
			}
			summ += x
			n++
		}
	}

	return summ
}

func fixFailedPages(index int, pageOrder *PageOrder, manualPages *ManualPages) []int {
	pages := manualPages.data[index]
	availableLess := make([]IntSlicePair, len(pages))
	for i, v := range pages {
		availableLess[i] = IntSlicePair{key: v, values: []int{}}
		for j, u := range pages {
			if i != j {
				_, ok := slices.BinarySearch(pageOrder.data[v], u)
				if ok {
					availableLess[i].values = append(availableLess[i].values, u)
				}
			}
		}
	}

	result := make([]int, len(availableLess))
	sort.Sort(BySliceSize(availableLess))
	for i, a := range availableLess {
		result[i] = a.key
	}

	return result
}

func Execute() {
	pageOrder, manualPages := readInput("input.txt")

	summ0, _, wrongIndexes := calcAnswer(&pageOrder, &manualPages)
	fmt.Printf("ans1=%d\n", summ0)

	summFixed := 0
	for _, wi := range wrongIndexes {
		fixedOrder := fixFailedPages(wi, &pageOrder, &manualPages)
		summFixed += fixedOrder[len(fixedOrder)/2]
		//fmt.Printf("index=%d %v\n", wi, fixedOrder)
		//break
	}

	fmt.Printf("fixed summ = %d", summFixed)

	//summ := calcAnswerParallel(&pageOrder, &manualPages)
	//fmt.Printf("ans1p=%d\n", summ)
}
