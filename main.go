package main

import (
	"fmt"
	"log"

	"healerdb/config"
	"healerdb/dbq"
)

func main() {
	fmt.Println("--------------------------------------------------")
	fmt.Println("Hello World")

	// Read the config file and get the connection string
	connString, err := config.GetConnStr()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection String:", connString)

	// Connect to the database
	session, err := dbq.Connect(connString)
	if err != nil {
		log.Fatal(err)
	}

	// Check the session to the database is alive or not and ping the database
	err = session.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Defer to close the session to the database
	defer session.Close()

	// Drop all the databases
	// err = dbq.DropAllDatabases(session)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("All the databases dropped")

	// print a seperator
	fmt.Println("--------------------------------------------------")

	// Run DBFirstSetup function to create the necessary databases and collections
	err = dbq.DBFirstSetup(session)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[+] Database setup complete")
	// print a seperator
	fmt.Println("--------------------------------------------------")

	// Add a new target to the database called "test"
	err = dbq.AddTargetToDB(session, "test")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[INFO] Finished job adding a new target to the database")

	// print a seperator
	fmt.Println("--------------------------------------------------")

	// Add a new target to the database called "test2"
	err = dbq.AddTargetToDB(session, "test2")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[INFO] Finished job adding a new target to the database\nNew Target: test2")

	// Add a domain to the database enum, to target "test", domain: "test.com"
	err = dbq.AddDomainToDB(session, "test", "test.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[INFO] Finished job adding a new domain to the database\nNew Domain: test.com")

	// print a seperator
	fmt.Println("--------------------------------------------------")
	fmt.Println("Done")
}
