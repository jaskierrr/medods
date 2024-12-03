package main

import (
	"main/bootstrapper"
	"log"
)

func main() {
	err := bootstrapper.New().RunAPI()
	if err != nil {
		log.Printf("failed to start: %v", err)
	}
 }
