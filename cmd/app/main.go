package main

import (
	"main/bootstrapper"
	"log"
)

func main() {
	err := bootstrapper.New().RunAPI()
	if err != nil {
		log.Fatalf("failed to start: %v", err)
	}
 }
