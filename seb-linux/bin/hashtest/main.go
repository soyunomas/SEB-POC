package main

import (
	"fmt"
	"log"
	"os"

	"seb-linux/internal/crypto"
)

func main() {
	path := "test_exam.seb"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	ck, err := crypto.DeriveConfigKey(path)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("ConfigKey: %s\n", ck)
}
