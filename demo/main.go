package main

import (
	"fmt"

	"github.com/duggaraju/c2pa-go/c2pa"
)

func main() {
	version := c2pa.Version()
	fmt.Printf("c2pa version: %s\n",  version)
}