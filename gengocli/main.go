package main

import (
	"fmt"
	"github.com/ls6-events/gengo/gengocli/cmd"
	"os"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Printf("Failed to execute command: %s\n", err.Error())
		os.Exit(1)
	}
}
