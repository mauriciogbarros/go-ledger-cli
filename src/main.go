package main

import (
	"fmt"
	"log"
	"os"

	"go.mod/internal/run"
)

func main() {
	fmt.Println()
	err := run.Run()
	if err != nil {
		log.Printf("Error => %v\n", err)
		fmt.Println()
		os.Exit(1)
	}
	fmt.Println()
	os.Exit(0)
}
