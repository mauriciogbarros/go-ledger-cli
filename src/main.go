package main

import (
	"fmt"
	"log"
	"os"

	"go.mod/internal/cli"
)

func main() {
	msg, err := cli.Run()
	if err != nil {
		log.Printf("Error => %v\n", err)
		fmt.Println()
		os.Exit(1)
	}
	log.Printf("Status =>  %s\n", msg)
	fmt.Println()
	os.Exit(0)
}
