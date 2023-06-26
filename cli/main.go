package main

import "github.com/ls6-events/gengo/cli/cmd"

func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
