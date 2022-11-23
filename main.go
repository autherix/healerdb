package main

import (
	"fmt"

	"healerdb/conndb"
)

func main() {
	fmt.Println("Hello World")
	conndb.Connect()
}
