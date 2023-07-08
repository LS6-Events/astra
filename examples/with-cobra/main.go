package main

import "withcobra/cmd"

func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
