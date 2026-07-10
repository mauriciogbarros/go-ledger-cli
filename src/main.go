package main

import (
	"log"
	"os"

	"go.mod/internal/cli"
)


func main() {
	msg, err := cli.Run()
	if err != nil {
		log.Printf("Error => %v\n", err)
		os.Exit(1)
	}
	log.Printf("Status =>  %s\n", msg)
	os.Exit(0)
}
