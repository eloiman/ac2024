package day3

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
)

func readInputMuls1(filename string) int64 {
	var result int64 = 0

	file, ferr := os.Open(filename)
	if ferr != nil {
		log.Fatal(ferr)
		os.Exit(1)
	}

	defer file.Close()

	re := regexp.MustCompile("mul\\(([0-9]*),([0-9]*)\\)")

	defaultBuffSize := 128
	var currBuff []byte
	buff := make([]byte, defaultBuffSize)
	buffLen := 0
	var readErr error = nil
	rd := bufio.NewReaderSize(file, defaultBuffSize)
	for readErr == nil {
		buffLen, readErr = rd.Read(buff)
		if readErr != nil || buffLen == 0 {
			fmt.Printf("Stopped reading input\n")
			break
		}
		currBuff = append(currBuff, buff[:buffLen]...)
		strBuff := string(currBuff)
		submatches := re.FindAllStringSubmatch(strBuff, -1)
		for _, submatch := range submatches {
			if len(submatch) == 3 {
				m1, _ := strconv.Atoi(submatch[1])
				m2, _ := strconv.Atoi(submatch[2])
				result += int64(m1) * int64(m2)
			}
		}
		indexes := re.FindAllIndex(currBuff, -1)
		unprocessedLen := len(currBuff) - indexes[len(indexes)-1][1]
		currBuffClone := make([]byte, unprocessedLen)
		copy(currBuffClone, currBuff[len(currBuff)-unprocessedLen:])
		currBuff = currBuffClone
	}

	return result
}

type parserState struct {
	currBuff         []byte
	isParsingEnabled bool
	result           int64
}

type location struct {
	loc0 int
	loc1 int
}

type parser interface {
	visit(state *parserState)
	getElementLocation() location
}

type mulLocation struct {
	location
}

type doLocation struct {
	location
}

type dontLocation struct {
	location
}

func makeMulRegexp() func() *regexp.Regexp {
	re := regexp.MustCompile("mul\\([0-9]*,[0-9]*\\)")
	return func() *regexp.Regexp {
		return re
	}
}

func makeDoRegexp() func() *regexp.Regexp {
	re := regexp.MustCompile("do\\(\\)")
	return func() *regexp.Regexp {
		return re
	}
}

func makeDontRegexp() func() *regexp.Regexp {
	re := regexp.MustCompile("don't\\(\\)")
	return func() *regexp.Regexp {
		return re
	}
}

func (loc mulLocation) visit(state *parserState) {
	if state.isParsingEnabled {
		re := regexp.MustCompile("mul\\(([0-9]*),([0-9]*)\\)")
		subbuff := state.currBuff[loc.loc0:loc.loc1]
		submatch := re.FindStringSubmatch(string(subbuff))
		if len(submatch) == 3 {
			m1, _ := strconv.Atoi(submatch[1])
			m2, _ := strconv.Atoi(submatch[2])
			state.result += int64(m1) * int64(m2)
			//fmt.Printf("%d %d\n", m1, m2)
		}
	}
}

func (loc doLocation) visit(state *parserState) {
	state.isParsingEnabled = true
}

func (loc dontLocation) visit(state *parserState) {
	state.isParsingEnabled = false
}

func (loc mulLocation) getElementLocation() location {
	return loc.location
}

func (loc dontLocation) getElementLocation() location {
	return loc.location
}

func (loc doLocation) getElementLocation() location {
	return loc.location
}

type parserArray struct {
	pr []parser
}

func (prArr parserArray) Len() int {
	return len(prArr.pr)
}

func (prArr parserArray) Less(i, j int) bool {
	return prArr.pr[i].getElementLocation().loc0 < prArr.pr[j].getElementLocation().loc0
}

func (prArr parserArray) Swap(i, j int) {
	prArr.pr[i], prArr.pr[j] = prArr.pr[j], prArr.pr[i]
}

func makeParserState() *parserState {
	state := &parserState{}
	state.isParsingEnabled = true
	state.result = 0

	return state
}

func mergeUnique(pr1 []parser, pr2 []parser) []parser {
	result := append(pr1, pr2...)
	var sortIntf sort.Interface = parserArray{pr: result}
	sort.Sort(sortIntf)

	return result
}

func readInputMuls2(filename string) int64 {
	file, ferr := os.Open(filename)
	if ferr != nil {
		log.Fatal(ferr)
		os.Exit(1)
	}

	defer file.Close()

	reMul := makeMulRegexp()
	reDont := makeDontRegexp()
	reDo := makeDoRegexp()

	state := makeParserState()

	defaultBuffSize := 128
	buff := make([]byte, defaultBuffSize)
	buffLen := 0
	var readErr error = nil
	rd := bufio.NewReaderSize(file, defaultBuffSize)
	for readErr == nil {
		buffLen, readErr = rd.Read(buff)
		if readErr != nil || buffLen == 0 {
			fmt.Printf("Stopped reading input\n")
			break
		}
		state.currBuff = append(state.currBuff, buff[:buffLen]...)
		seqs := []parser{}
		mulLocs := reMul().FindAllIndex(state.currBuff, -1)
		for _, loc := range mulLocs {
			seqs = append(seqs, mulLocation{location: location{loc[0], loc[1]}})
		}
		doLocs := reDo().FindAllIndex(state.currBuff, -1)
		for _, loc := range doLocs {
			seqs = append(seqs, doLocation{location: location{loc[0], loc[1]}})
		}
		dontLocs := reDont().FindAllIndex(state.currBuff, -1)
		for _, loc := range dontLocs {
			seqs = append(seqs, dontLocation{location: location{loc[0], loc[1]}})
		}

		var sortIntf sort.Interface = parserArray{pr: seqs}
		sort.Sort(sortIntf)

		lastLocation := 0
		for _, seq := range seqs {
			seq.visit(state)
			lastLocation = seq.getElementLocation().loc1
		}

		unprocessedLen := len(state.currBuff) - lastLocation
		currBuffClone := make([]byte, unprocessedLen)
		copy(currBuffClone, state.currBuff[len(state.currBuff)-unprocessedLen:])
		state.currBuff = currBuffClone
	}

	return state.result
}

func Execute() {
	summ1 := readInputMuls1("day3_1.txt")
	fmt.Printf("summ1 = %d\n", summ1)

	summ2 := readInputMuls2("day3_1.txt")
	fmt.Printf("summ2 = %d\n", summ2)

	/*

		seqs1 := []parser{
			&dontLocation{location: location{3, 4}},
			&mulLocation{location: location{0, 1}}}
		seqs2 := []parser{
			&doLocation{location: location{20, 120}},
			&mulLocation{location: location{10, 11}},
			&mulLocation{location: location{30, 1000}}}

		var seq = mergeUnique(seqs1, seqs2)

		for _, p := range seq {
			loc := p.getElementLocation()
			fmt.Printf("%d %d\n", loc.loc0, loc.loc1)
			p.visit(state)
		}

		fmt.Printf("result = %d", state.result)
	*/
}
