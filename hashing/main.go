package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		usage()
	}

	switch os.Args[1] {
	default:
		usage()
	case "verify":
		hashed := sha256.Sum256([]byte(os.Args[2]))
		if hex.EncodeToString(hashed[:]) == os.Args[3] {
			fmt.Println("OK")
		}
		fmt.Printf("MISMATCH\n  expected: %s\n  actual:   %s\n", os.Args[3], hashed)
		os.Exit(1)
	case "hash":
		hashed := sha256.Sum256([]byte(os.Args[2]))
		fmt.Printf("hashed data: %x \n", hex.EncodeToString(hashed[:]))
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage:
  hashlab hash <string>
  hashlab verify <string> <expected-hex>`)
	os.Exit(2)
}
