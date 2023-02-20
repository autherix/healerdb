package main

import (
	"fmt"

	"healerdb/dbquery"
	"healerdb/myutils"
)

func main() {
	fmt.Println("Hello, World!")

	// Create a client to mongoDB server using connstr
	connstr := "mongodb://healerdb:hamidpapi@localhost:27017"
	client, err := dbquery.CreateClient(connstr)
	if err != nil {
		fmt.Println(err)
		return
	}
	// defer to disconnect from mongoDB server
	defer client.Disconnect(nil)

	fmt.Println("Connected to MongoDB!")

	// print a seperator
	fmt.Println("--------------------------------------------------")

	// list all the databases in the mongoDB server
	databases, err := dbquery.GetDatabases(client)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(databases)

	// print a seperator
	fmt.Println("--------------------------------------------------")

	// Create a database called 'safe-panel'
	dbname := "safe-panel"
	err = dbquery.CreateDatabase(client, dbname)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to create database")
		// return
	} else {
		fmt.Println("Database created!")
	}

	// print a seperator
	fmt.Println("--------------------------------------------------")

	// Create a collection called 'users' in the database 'safe-panel'
	collectionname := "users"
	err = dbquery.CreateCollection(client, dbname, collectionname)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to create collection")
		// return
	} else {
		fmt.Println("Collection created!")
	}

	// create an index on the field 'email' in the collection 'users' in the database 'safe-panel'
	indexname := "email"
	err = dbquery.AddUniqueIndex(client, dbname, collectionname, indexname)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to create index")
		return
	} else {
		fmt.Println("Index created!")
	}
	// Create an index on the field 'username' in the collection 'users' in the database 'safe-panel'
	indexname = "username"
	err = dbquery.AddUniqueIndex(client, dbname, collectionname, indexname)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to create index")
		return
	} else {
		fmt.Println("Index created!")
	}

	// print a seperator
	fmt.Println("--------------------------------------------------")

	passwd := "123456"
	passwd_hash := myutils.HashString(passwd)
	// Insert a document in	to the collection 'users' in the database 'safe-panel' using a struct type 'User'
	type User struct {
		Username   string `json:"username"`
		PasswdHash string `json:"passwd_hash"`
		Email      string `json:"email"`
	}
	admin_user := User{
		Username:   "admin",
		PasswdHash: passwd_hash,
		Email:      "admin@autherix.com",
	}
	// convert to json string
	user_json, err := myutils.Struct2json(admin_user)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to convert struct to json")
		return
	}
	fmt.Println(user_json)
	err = dbquery.InsertDocument(client, dbname, collectionname, user_json)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to insert document")
		// return
	} else {
		fmt.Println("Document inserted!")
	}

	// print a seperator
	fmt.Println("--------------------------------------------------")

	// Create a target called 'surf' in database 'enum'
	dbname = "enum"
	collectionname = "surf"
	err = dbquery.CreateDatabase(client, dbname)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to create database")
		// return
	} else {
		fmt.Println("Database created!")
	}
	err = dbquery.CreateCollection(client, dbname, collectionname)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to create collection")
		// return
	} else {
		fmt.Println("Collection created!")
	}

	// print a seperator
	fmt.Println("--------------------------------------------------")

	// Add a domain to the target 'surf' in database 'enum'
	domain := "test6.com"
	subdomain := "sub.test6.com"
	// Use AddDomain function to add a domain to the target 'surf' in database 'enum'
	err = dbquery.AddDomain(client, dbname, collectionname, domain)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to add domain")
		// return
	} else {
		fmt.Println("Domain added!")
	}

	// print a seperator
	fmt.Println("--------------------------------------------------")

	// Check if the subdomain subdomain is already present in the target 'surf' in database 'enum'
	subexists, err := dbquery.CheckSubdomain(client, dbname, collectionname, domain, subdomain)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to check subdomain")
		// return
	} else {
		fmt.Println("Subdomain exists:", subexists)
	}

	// print a seperator
	fmt.Println("--------------------------------------------------")
}
