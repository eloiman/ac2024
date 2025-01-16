package day4

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func readInput(filename string) [][]byte {
	file, ferr := os.Open(filename)
	if ferr != nil {
		log.Fatal(ferr)
		os.Exit(1)
	}

	defer file.Close()

	var result [][]byte

	scaner := bufio.NewScanner(file)

	for scaner.Scan() {
		ss := scaner.Text()
		result = append(result, []byte(ss))
	}

	return result
}

type buffer struct {
	field [][]byte
}

type location struct {
	x, y int
}

func (buff *buffer) checkConstrains(loc location) bool {
	return loc.x >= 0 && loc.x < len(buff.field[0]) && loc.y >= 0 && loc.y < len(buff.field)
}

type locationGen func(location, int) location
type locationIndexedGen func(location, int, int) location

func getDirections() []locationIndexedGen {
	dirs := []locationGen{
		func(loc location, i int) location { return location{loc.x, loc.y + i} },
		func(loc location, i int) location { return location{loc.x, loc.y - i} },
		func(loc location, i int) location { return location{loc.x + i, loc.y} },
		func(loc location, i int) location { return location{loc.x - i, loc.y} },
		func(loc location, i int) location { return location{loc.x + i, loc.y + i} },
		func(loc location, i int) location { return location{loc.x - i, loc.y - i} },
		func(loc location, i int) location { return location{loc.x + i, loc.y - i} },
		func(loc location, i int) location { return location{loc.x - i, loc.y + i} }}

	output := []locationIndexedGen{}
	for _, f := range dirs {
		indexedGen := func(loc location, index int, len int) location {
			return f(loc, index)
		}
		output = append(output, indexedGen)
	}

	return output
}

func getDirectionsX() []locationIndexedGen {
	dirs := []locationIndexedGen{
		func(loc location, i int, len int) location {
			return location{loc.x - len/2 + i, loc.y + len/2 - i}
		},
		func(loc location, i int, len int) location {
			return location{loc.x + len/2 - i, loc.y - len/2 + i}
		},
		func(loc location, i int, len int) location {
			return location{loc.x - len/2 + i, loc.y - len/2 + i}
		},
		func(loc location, i int, len int) location {
			return location{loc.x + len/2 - i, loc.y + len/2 - i}
		}}

	return dirs
}

func (buff *buffer) checkDirectionContent(content string, cloc location, locGen locationIndexedGen) bool {
	cindex := 0
	loc := cloc
	for ; cindex < len(content); cindex++ {
		loc = locGen(cloc, cindex, len(content))
		if !buff.checkConstrains(loc) {
			break
		}

		if buff.field[loc.y][loc.x] != content[cindex] {
			break
		}
	}

	if cindex == len(content) {
		return true
	}

	return false
}

func Execute() {
	input := readInput("input.txt")
	fmt.Printf("ylen = %d\n", len(input[0]))

	buff := buffer{field: input}

	cntr := uint64(0)
	dirs := getDirections()
	for y := 0; y < len(buff.field); y++ {
		for x := 0; x < len(buff.field[0]); x++ {
			for idirs := 0; idirs < len(dirs); idirs++ {
				if buff.checkDirectionContent("XMAS", location{x, y}, dirs[idirs]) {
					cntr++
				}
			}
		}
	}

	fmt.Printf("cntr=%d\n", cntr)

	cntrX := uint64(0)
	dirsX := getDirectionsX()
	for y := 0; y < len(buff.field); y++ {
		for x := 0; x < len(buff.field[0]); x++ {
			if (buff.checkDirectionContent("MAS", location{x, y}, dirsX[0]) || buff.checkDirectionContent("MAS", location{x, y}, dirsX[1])) &&
				(buff.checkDirectionContent("MAS", location{x, y}, dirsX[2]) || buff.checkDirectionContent("MAS", location{x, y}, dirsX[3])) {
				cntrX++
			}
		}
	}

	fmt.Printf("cntrX=%d\n", cntrX)
}
