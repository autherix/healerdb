package main

import (
	"fmt"

	"db/dbquery"
	"db/myutils"
)

func main() {
	fmt.Println("Hello, World!")

	// Create a client to mongoDB server using connstr
	connstr := "mongodb://localhost:27017"
	client, err := dbquery.Createclient(connstr)
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
	databases, err := dbquery.Getdatabases(client)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(databases)

	// print a seperator
	fmt.Println("--------------------------------------------------")

	// Create a database called 'safe-panel'
	dbname := "safe-panel"
	err = dbquery.Createdatabase(client, dbname)
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
	err = dbquery.Createcollection(client, dbname, collectionname)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to create collection")
		// return
	} else {
		fmt.Println("Collection created!")
	}

	// create an index on the field 'email' in the collection 'users' in the database 'safe-panel'
	indexname := "email"
	err = dbquery.Adduniqueindex(client, dbname, collectionname, indexname)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to create index")
		return
	} else {
		fmt.Println("Index created!")
	}
	// Create an index on the field 'username' in the collection 'users' in the database 'safe-panel'
	indexname = "username"
	err = dbquery.Adduniqueindex(client, dbname, collectionname, indexname)
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
	// Insert a document into the collection 'users' in the database 'safe-panel' using a struct type 'User'
	type User struct {
		Username   string `json:"username"`
		PasswdHash string `json:"passwd_hash"`
		Email      string `json:"email"`
	}
	admin_user := User{
		Username:   "admin2",
		PasswdHash: passwd_hash,
		Email:      "admin2@autherix.com",
	}
	// convert to json string
	user_json, err := myutils.Struct2json(admin_user)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to convert struct to json")
		return
	}
	fmt.Println(user_json)
	err = dbquery.Insertdocument(client, dbname, collectionname, user_json)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to insert document")
		// return
	} else {
		fmt.Println("Document inserted!")
	}

	// print a seperator
	fmt.Println("--------------------------------------------------")

}
