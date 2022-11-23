package main

import (
	"fmt"

	"gomongo/conndb"
	"gomongo/gomongo"
)

func main() {
	fmt.Println("Hello World")
	gomongo.Usefmt()
	conndb.Connect()
}
