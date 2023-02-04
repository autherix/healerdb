package dbquery

import (
	"runtime"

	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"          // ignore this error
	"go.mongodb.org/mongo-driver/mongo"         // ignore this error
	"go.mongodb.org/mongo-driver/mongo/options" // ignore this error

	// primitive
	"go.mongodb.org/mongo-driver/bson/primitive" // ignore this error
)

// function trace to print the function name and the line number info
func trace() {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	fmt.Printf("%s:%d %s ", file, line, f.Name())
}

// Function to create a session to the mongoDB database using the given connection string, returns a pointer to the session and an error
func Createclient(connectionString string) (*mongo.Client, error) {
	// set options for the client to the database like the connection timeout
	options := options.Client().ApplyURI(connectionString).SetConnectTimeout(10 * time.Second)
	// Create a client to the database
	client, err := mongo.NewClient(options)
	if err != nil {
		return nil, fmt.Errorf("[-] Error creating client to database: %v", err)
	}
	fmt.Println("[+] client created successfully")

	// Connect to the database
	err = client.Connect(nil)
	if err != nil {
		return nil, fmt.Errorf("[-] Error connecting to database: %v", err)
	}
	fmt.Println("[+] Connected to database successfully")

	// Ping the database to check if the connection is successful
	err = client.Ping(nil, nil)
	if err != nil {
		return nil, fmt.Errorf("[-] Error pinging database: %v", err)
	}
	fmt.Println("[+] Pinged database successfully")

	return client, nil
}

// Function to close the session to the database
func Closeclient(client *mongo.Client) {
	if err := client.Disconnect(nil); err != nil {
		panic(err)
	}
	fmt.Println("[+] Disconnected from database successfully")
}

// function to get the pointer to the client to the database, fetches all database names and returns a slice of strings containing the names of the databases and an error
func Getdatabases(client *mongo.Client) ([]string, error) {
	// Get all the databases
	databases, err := client.ListDatabaseNames(context.TODO(), bson.M{}, nil)
	if err != nil {
		return nil, fmt.Errorf("[-] Error getting databases: %v", err)
	}
	fmt.Println("[+] Got databases successfully")

	return databases, nil
}

// function to check if a database with the given name exists, returns a boolean and an error
func Checkdatabase(client *mongo.Client, database string) (bool, error) {
	// Check if the database exists
	databases, err := Getdatabases(client)
	if err != nil {
		return false, fmt.Errorf("[-] Error checking database: %v", err)
	}
	for _, db := range databases {
		if db == database {
			return true, nil
		}
	}

	return false, nil
}

// function to check if a collection with the given name exists in the given database, returns a boolean and an error
func Checkcollection(client *mongo.Client, database string, collection string) (bool, error) {
	// Check if the collection exists
	collections, err := Getcollections(client, database)
	if err != nil {
		return false, fmt.Errorf("[-] Error checking collection: %v", err)
	}
	for _, col := range collections {
		if col == collection {
			return true, nil
		}
	}

	return false, nil
}

// function to check if a document with the given id exists in the given collection in the given database, returns a boolean and an error
func Checkdocument(client *mongo.Client, database string, collection string, id string) (bool, error) {
	// Check if the document exists
	documents, err := Getdocuments(client, database, collection)
	if err != nil {
		return false, fmt.Errorf("[-] Error checking document: %v", err)
	}
	for _, doc := range documents {
		if doc == id {
			return true, nil
		}
	}

	return false, nil
}

// function to create a document in the given collection in the given database, with the provided json string of the document, returns an error -> convert the json string to an interface or bson object and insert it into the collection
func Createdocument(client *mongo.Client, database string, collection string, document string) error {
	// convert the json string to an interface using unmashal
	var doc interface{}
	err := json.Unmarshal([]byte(document), &doc)
	if err != nil {
		return fmt.Errorf("[-] error creating document: %v", err)

	}

	// insert the document into the collection
	_, err = client.Database(database).Collection(collection).InsertOne(context.TODO(), doc)
	if err != nil {
		return fmt.Errorf("[-] error creating document: %v", err)

	}
	fmt.Println("[+] Created document successfully")

	return nil
}

// function to create a collection, with the provided name, in the given database, returns an error -> for creating a collection, it just creates a document in the collection with the content as 'exists': true
func Createcollection(client *mongo.Client, database string, collection string) error {
	// Check if the collection exists
	exists, err2 := Checkcollection(client, database, collection)
	if err2 != nil {
		return fmt.Errorf("[-] Error creating collection: %v", err2)

	}
	if exists {
		return fmt.Errorf("[-] Error creating collection: collection already exists")

	}

	// Create a collection, use CreateDocument function to create a document in the collection
	err := Createdocument(client, database, collection, `{"exists": true}`)
	if err != nil {
		return fmt.Errorf("[-] Error creating collection: %v", err)

	}
	fmt.Println("[+] Created collection successfully")

	return nil
}

// function to create a database, with the provided name, returns an error -> for creating a database, it just creates a collection in the database with the name 'exists' and the content as 'exists': true
func Createdatabase(client *mongo.Client, database string) error {
	// Check if the database exists
	exists, err2 := Checkdatabase(client, database)
	if err2 != nil {
		return fmt.Errorf("[-] Error creating database: %v", err2)

	}
	if exists {
		return fmt.Errorf("[-] Error creating database: database already exists")

	}

	// Create a database, use CreateCollection function to create a collection in the database
	err := Createcollection(client, database, "exists")
	if err != nil {
		return fmt.Errorf("[-] Error creating database: %v", err)

	}
	fmt.Println("[+] Created database successfully")

	return nil
}

// function to drop a collection, with the provided name(removes if exists), in the given database, returns an error
func Dropcollection(client *mongo.Client, database string, collection string) error {
	// Drop a collection
	err := client.Database(database).Collection(collection).Drop(nil)
	if err != nil {
		return fmt.Errorf("[-] Error dropping collection: %v", err)

	}
	fmt.Println("[+] Dropped collection successfully")

	return nil
}

// function to drop a database, with the provided name(removes if exists), returns an error
func Dropdatabase(client *mongo.Client, database string) error {
	// Drop a database
	err := client.Database(database).Drop(nil)
	if err != nil {
		return fmt.Errorf("[-] Error dropping database: %v", err)

	}
	fmt.Println("[+] Dropped database successfully")

	return nil
}

// function to get the pointer to the client to the database, a database name, fetches all collection names in the given database and returns a slice of strings containing the names of the collections and an error
func Getcollections(client *mongo.Client, database string) ([]string, error) {
	// Get all the collections in the database
	collections, err := client.Database(database).ListCollectionNames(context.TODO(), bson.M{}, nil)
	if err != nil {
		return nil, fmt.Errorf("[-] Error getting collections: %v", err)
	}
	fmt.Println("[+] Got collections successfully")

	return collections, nil
}

// function to get the pointer to the client to the database, a database name and a collection name, fetches all documents in the given collection, converts each document to a string of json object and returns a slice of strings containing the json objects and an error
func Getdocuments(client *mongo.Client, database string, collection string) ([]string, error) {
	// Get all the documents in the collection
	cursor, err := client.Database(database).Collection(collection).Find(nil, nil)
	if err != nil {
		return nil, fmt.Errorf("[-] Error getting documents: %v", err)
	}
	fmt.Println("[+] Got documents successfully")

	// Convert the documents to json objects
	var documents []string
	for cursor.Next(nil) {
		var document string
		err = cursor.Decode(&document)
		if err != nil {
			return nil, fmt.Errorf("[-] Error converting documents to json: %v", err)
		}
		documents = append(documents, document)
	}

	return documents, nil
}

// function to get pointer to client to database, database name, collection name and document id, fetches the document with the given id from the given collection in the given database, converts the document to a string of json object and returns the json object and an error
func Getdocument(client *mongo.Client, database string, collection string, id string) (string, error) {
	// Get the document with the given id from the collection
	var document string = ""

	// Convert the id to an object id
	objectid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", fmt.Errorf("[-] Error converting id to object id: %v", err)
	}

	// Get the document
	err = client.Database(database).Collection(collection).FindOne(nil, bson.M{"_id": objectid}).Decode(&document)
	if err != nil {
		return "", fmt.Errorf("[-] Error getting document: %v", err)
	}
	fmt.Println("[+] Got document successfully")

	return document, nil
}

// function to get pointer to client to database, database name, collection name and document id, deletes the document with the given id from the given collection in the given database, returns an error
func Deletedocument(client *mongo.Client, database string, collection string, id string) error {
	// Delete the document with the given id from the collection if exists
	var document string = ""

	// Convert the id to an object id
	objectid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("[-] Error converting id to object id: %v", err)

	}

	filter := bson.M{"_id": objectid} // filter to find the document with the given id

	// Delete the document
	err = client.Database(database).Collection(collection).FindOneAndDelete(nil, filter).Decode(&document)
	if err != nil {
		return fmt.Errorf("[-] Error deleting document: %v", err)

	}
	fmt.Println("[+] Deleted document successfully")

	return nil
}

// function to get client, db, coll, and query string, returns a slice of strings containing the json objects of the documents that match the query and an error
func Querydocuments(client *mongo.Client, database string, collection string, query string) ([]string, error) {
	// Query the documents in the collection with the given query
	var documents []string

	// Convert the query to a bson object
	bsonquery, err := primitive.ObjectIDFromHex(query)
	if err != nil {
		return nil, fmt.Errorf("[-] Error converting query to bson: %v", err)
	}

	// Query the documents
	cursor, err := client.Database(database).Collection(collection).Find(nil, bsonquery)
	if err != nil {
		return nil, fmt.Errorf("[-] Error querying documents: %v", err)
	}
	fmt.Println("[+] Queried documents successfully")

	// Convert the documents to json objects
	for cursor.Next(nil) {
		var document string
		err = cursor.Decode(&document)
		if err != nil {
			return nil, fmt.Errorf("[-] Error converting documents to json: %v", err)
		}
		documents = append(documents, document)
	}

	return documents, nil
}

// function to get client, db, coll, and document string, inserts the document into the collection in the database, returns an error
func Insertdocument(client *mongo.Client, database string, collection string, document string) error {
	// convert the json document to a bson object
	var bsondocument bson.M = bson.M{}
	err := json.Unmarshal([]byte(document), &bsondocument)
	if err != nil {
		return fmt.Errorf("[-] Error converting document to bson: %v", err)

	}

	// Insert the document into the collection
	_, err = client.Database(database).Collection(collection).InsertOne(nil, bsondocument)
	if err != nil {
		return fmt.Errorf("[-] Error inserting document: %v", err)

	}
	fmt.Println("[+] Inserted document successfully")

	return nil
}

// function to get client, db, coll, document id and document string, inserts the documents into the collection in the database, returns an error -> don't forget the documents are passed as a string of json objects, convert them to bson objects before inserting
func Insertdocuments(client *mongo.Client, database string, collection string, documents string) error {
	// Convert the documents to interface objects using unmarshal
	var bsondocuments []interface{}
	err := json.Unmarshal([]byte(documents), &bsondocuments)
	if err != nil {
		return fmt.Errorf("[-] Error converting documents to bson: %v", err)

	}

	// Insert the documents into the collection
	_, err = client.Database(database).Collection(collection).InsertMany(nil, bsondocuments)
	if err != nil {
		return fmt.Errorf("[-] Error inserting documents: %v", err)

	}
	fmt.Println("[+] Inserted documents successfully")

	return nil
}

// function to add unique index to a collection in a database, returns an error
func Adduniqueindex(client *mongo.Client, database string, collection string, index string) error {
	// Add unique index to the collection
	_, err := client.Database(database).Collection(collection).Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.M{index: 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("[-] Error adding unique index: %v", err)

	}
	fmt.Println("[+] Added unique index successfully")

	return nil
}
