package main

import (
	"os"

	dayPackage "strelox.com/ac2024/day6"
)

type dayExecuter interface {
	Execute()
}

type dayPackageOptions struct {
	inputsFolder string
}

func (options *dayPackageOptions) Execute() {
	os.Chdir(options.inputsFolder)
	dayPackage.Execute()
}

func main() {
	dayOptions := dayPackageOptions{inputsFolder: "day6"}
	var task dayExecuter = &dayOptions
	task.Execute()
}
