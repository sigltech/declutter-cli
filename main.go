package main

import (
	"declutter/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		return
	}
}
