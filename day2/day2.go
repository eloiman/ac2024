package day2

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func readInput(filename string) [][]int {
	file, ferr := os.Open(filename)
	if ferr != nil {
		log.Fatal(ferr)
		os.Exit(1)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var result [][]int
	for scanner.Scan() {
		line := scanner.Text()
		lineReader := strings.NewReader(line)
		var lineInts []int
		for true {
			var n int
			_, err := fmt.Fscanf(lineReader, "%d", &n)
			if err != nil {
				break
			}

			lineInts = append(lineInts, n)
		}

		result = append(result, lineInts)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return result
}

func intAbs(a int, b int) int {
	if a > b {
		return a - b
	}

	return b - a
}

func Execute() {
	input := readInput("day2_1.txt")

	var nOkReports int = 0
	for _, report := range input {
		lenReport := len(report)
		reportTest := report
		for i := 0; i <= lenReport; i++ {
			reportReversed := sort.Reverse(sort.IntSlice(reportTest))
			if sort.IntsAreSorted(reportTest) || sort.IsSorted(reportReversed) {
				var prevLevel int = reportTest[0]
				for k, level := range reportTest {
					if k > 0 {
						abs := intAbs(prevLevel, level)
						if abs >= 1 && abs <= 3 {
							prevLevel = level
						} else {
							prevLevel = -1
							break
						}
					}
				}

				if prevLevel != -1 {
					nOkReports++
					break
				} else {
					reportTest = nil
					for j := 0; j < lenReport; j++ {
						if j != i {
							reportTest = append(reportTest, report[j])
						}
					}
				}
			} else if i != lenReport {
				reportTest = nil
				for j := 0; j < lenReport; j++ {
					if j != i {
						reportTest = append(reportTest, report[j])
					}
				}
			}
		}
	}

	fmt.Printf("%d", nOkReports)
}
