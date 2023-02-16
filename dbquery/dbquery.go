package dbquery

import (
	"runtime"

	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/tidwall/gjson"
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
func CreateClient(connectionString string) (*mongo.Client, error) {
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
func CloseClient(client *mongo.Client) {
	if err := client.Disconnect(nil); err != nil {
		panic(err)
	}
	fmt.Println("[+] Disconnected from database successfully")
}

// function to get the pointer to the client to the database, fetches all database names and returns a slice of strings containing the names of the databases and an error
func GetDatabases(client *mongo.Client) ([]string, error) {
	// Get all the databases
	databases, err := client.ListDatabaseNames(context.TODO(), bson.M{}, nil)
	if err != nil {
		return nil, fmt.Errorf("[-] Error getting databases: %v", err)
	}
	fmt.Println("[+] Got databases successfully")

	return databases, nil
}

// function to check if a database with the given name exists, returns a boolean and an error
func CheckDatabase(client *mongo.Client, database string) (bool, error) {
	// Check if the database exists
	databases, err := GetDatabases(client)
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

// function PurgeDatabases to delete all the databases in the database except the admin and config databases, returns an error
func PurgeDatabases(client *mongo.Client) error {
	// Get all the databases
	databases, err := GetDatabases(client)
	if err != nil {
		return fmt.Errorf("[-] Error purging databases: %v", err)
	}

	// Delete all the databases except the admin and config databases
	for _, db := range databases {
		if db != "admin" && db != "config" {
			err = client.Database(db).Drop(context.Background())
			if err != nil {
				return fmt.Errorf("[-] Error purging databases: %v", err)
			}
		}
	}
	fmt.Println("[+] Purged databases successfully")

	return nil
}

// function to check if a collection with the given name exists in the given database, returns a boolean and an error
func CheckCollection(client *mongo.Client, database string, collection string) (bool, error) {
	// Check if the collection exists
	collections, err := GetCollections(client, database)
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
func CheckDocument(client *mongo.Client, database string, collection string, id string) (bool, error) {
	// Check if the document exists
	documents, err := GetDocuments(client, database, collection)
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
func CreateDocument(client *mongo.Client, database string, collection string, document string) error {
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
func CreateCollection(client *mongo.Client, database string, collection string) error {
	// Check if the collection exists
	exists, err2 := CheckCollection(client, database, collection)
	if err2 != nil {
		return fmt.Errorf("[-] Error creating collection: %v", err2)
	}
	if exists {
		return fmt.Errorf("[-] Error creating collection: collection already exists")
	}

	// Create a collection, use CreateDocument function to create a document in the collection
	err := CreateDocument(client, database, collection, `{"exists": true}`)
	if err != nil {
		return fmt.Errorf("[-] Error creating collection: %v", err)
	}
	fmt.Println("[+] Created collection successfully")

	return nil
}

// function to create a database, with the provided name, returns an error -> for creating a database, it just creates a collection in the database with the name 'exists' and the content as 'exists': true
func CreateDatabase(client *mongo.Client, database string) error {
	// Check if the database exists
	exists, err2 := CheckDatabase(client, database)
	if err2 != nil {
		return fmt.Errorf("[-] Error creating database: %v", err2)
	}
	if exists {
		return fmt.Errorf("[-] Error creating database: database already exists")
	}

	// Create a database, use CreateCollection function to create a collection in the database
	err := CreateCollection(client, database, "exists")
	if err != nil {
		return fmt.Errorf("[-] Error creating database: %v", err)
	}
	fmt.Println("[+] Created database successfully")

	return nil
}

// function to drop a collection, with the provided name(removes if exists), in the given database, returns an error
func DropCollection(client *mongo.Client, database string, collection string) error {
	// Drop a collection
	err := client.Database(database).Collection(collection).Drop(nil)
	if err != nil {
		return fmt.Errorf("[-] Error dropping collection: %v", err)
	}
	fmt.Println("[+] Dropped collection successfully")

	return nil
}

// function to drop a database, with the provided name(removes if exists), returns an error
func DropDatabase(client *mongo.Client, database string) error {
	// Drop a database
	err := client.Database(database).Drop(nil)
	if err != nil {
		return fmt.Errorf("[-] Error dropping database: %v", err)
	}
	fmt.Println("[+] Dropped database successfully")

	return nil
}

// function to get the pointer to the client to the database, a database name, fetches all collection names in the given database and returns a slice of strings containing the names of the collections and an error
func GetCollections(client *mongo.Client, database string) ([]string, error) {
	// Get all the collections in the database
	collections, err := client.Database(database).ListCollectionNames(context.TODO(), bson.M{}, nil)
	if err != nil {
		return nil, fmt.Errorf("[-] Error getting collections: %v", err)
	}
	fmt.Println("[+] Got collections successfully")

	return collections, nil
}

// function to get the pointer to the client to the database, a database name and a collection name, fetches all documents in the given collection, converts each document to a string of json object and returns a slice of strings containing the json objects and an error
func GetDocuments(client *mongo.Client, database string, collection string) ([]string, error) {
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
func GetDocument(client *mongo.Client, database string, collection string, id string) (string, error) {
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
func DeleteDocument(client *mongo.Client, database string, collection string, id string) error {
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

// function to get client, db, coll, and document string, inserts the document into the collection in the database, returns an error
func InsertDocument(client *mongo.Client, database string, collection string, document string) error {
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
func InsertDocuments(client *mongo.Client, database string, collection string, documents string) error {
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
func AddUniqueIndex(client *mongo.Client, database string, collection string, index string) error {
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

// function AddTarget to add a collection to the database, returns an error
func AddTarget(client *mongo.Client, database string, target string) error {
	// Check if the target already exists
	exists, err := CheckTarget(client, database, target)
	if err != nil {
		return fmt.Errorf("[-] Error checking target: %v", err)
	}
	if exists {
		return fmt.Errorf("[-] Target already exists")
	}

	// Add a collection with target name to the database, use CreateCollection() to create a collection
	err = client.Database(database).CreateCollection(nil, target)
	if err != nil {
		return fmt.Errorf("[-] Error adding target: %v", err)
	}
	fmt.Println("[+] Added target successfully")

	return nil
}

// function CheckTarget to check if a collection exists in the database, returns a boolean and an error
func CheckTarget(client *mongo.Client, database string, target string) (bool, error) {
	// Query the database for the collection, use CheckCollection() to check if a collection exists
	exists, err := CheckCollection(client, database, target)
	if err != nil {
		return false, fmt.Errorf("[-] Error checking target: %v", err)
	}

	return exists, nil
}

// function qdb to query the database, parameters are database name, collection name, query object, returns a string of json objects and an error
func QueryDocuments(client *mongo.Client, database string, collection string, query bson.M) (string, error) {
	// Query the database
	cursor, err := client.Database(database).Collection(collection).Find(context.TODO(), query)
	if err != nil {
		return "", fmt.Errorf("[-] Error querying database: %v", err)
	}

	// Print the documents seperated by a newline and a '------------' line
	for cursor.Next(context.TODO()) {
		var document bson.M
		err = cursor.Decode(&document)
		if err != nil {
			return "", fmt.Errorf("[-] Error decoding document: %v", err)
		}
		fmt.Println(document)
		fmt.Println("------------")
	}

	// Convert the documents to json
	var documents []bson.M
	err = cursor.All(nil, &documents)
	if err != nil {
		return "", fmt.Errorf("[-] Error converting documents to json: %v", err)
	}

	// Convert the documents to json
	jsondocuments, err := json.Marshal(documents)
	if err != nil {
		return "", fmt.Errorf("[-] Error converting documents to json: %v", err)
	}

	return string(jsondocuments), nil
}

// function CheckDomain to check if the domain exists in the provided DB name, coll name, Use QueryDocuments() to query the database
func CheckDomain(client *mongo.Client, database string, collection string, domain string) (bool, error) {
	// Create a bson object with the domain `{"domain": domain}`
	bsondomain := bson.M{"domain": domain}

	// Query the database for the domain which returns a string of json objects and an error
	jsondocuments, err := QueryDocuments(client, database, collection, bsondomain)
	if err != nil {
		return false, fmt.Errorf("[-] Error checking domain: %v", err)
	}

	// Check if the json string is empty, if not the domain exists, so return true and nil
	if jsondocuments != "null" {
		return true, nil
	} else {
		return false, nil
	}
}

// function AddDomain to add a domain to the database, returns an error -> get client, db, collection(target), domain string
func AddDomain(client *mongo.Client, database string, target string, domain string) error {
	// Check if the domain already exists, use CheckDomain() to check if a domain exists
	exists, err := CheckDomain(client, database, target, domain)
	if err != nil {
		return fmt.Errorf("[-] Error checking domain: %v", err)
	}
	if exists {
		return fmt.Errorf("[-] Domain already exists")
	}
	// Create a bson object with the domain `{"domain": domain}`
	bsondomain := bson.M{"domain": domain}
	// convert to json string
	jsondomain, err := json.Marshal(bsondomain)
	if err != nil {
		return fmt.Errorf("[-] Error converting bson to json: %v", err)
	}

	// Check if the domain exists in the collection, use CheckDomain() to check if a domain exists
	// Insert the bson object into the collection use InsertDocument() to insert a document
	err = InsertDocument(client, database, target, string(jsondomain))
	if err != nil {
		return fmt.Errorf("[-] Error adding domain: %v", err)
	}
	fmt.Println("[+] Added domain successfully")

	return nil
}

// function AddSubdomain to add a subdomain to the provided database name, coll name, inside the document with the provided domain name, returns an error, create the document with the provided domain name if it doesn't exist
func CheckSubdomain(client *mongo.Client, database string, collection string, domain string, subdomain string) (bool, error) {
	// Get the document with the domain name, use QueryDocuments() to query the database
	jsondocuments, err := QueryDocuments(client, database, collection, bson.M{"domain": domain})
	if err != nil {
		return false, fmt.Errorf("[-] Error checking subdomain: %v", err)
	}

	// iterate over jsondocuments.#.subdomains.#.subdomain and check if the subdomain already exists
	fmt.Println(gjson.Get(jsondocuments, "#.subdomains"))
	subdomainsList := gjson.Get(jsondocuments, "#.subdomains.#.subdomain")

	// if subdomainsList is empty, the subdomain doesn't exist
	if subdomainsList.String() == "" {
		return false, nil
	}

	// print a seperator
	fmt.Println("--------------------------------------")

	// iterate over the subdomainsList and check if the subdomain already exists
	subdomainAddrList := subdomainsList.Array()[0].Array()
	for _, subdomainAddr := range subdomainAddrList {
		fmt.Println(subdomainAddr.String())
		if subdomainAddr.String() == subdomain {
			return true, nil
		}
	}

	return false, nil
}

// function AddSubdomain to add a subdomain to the provided database name, coll name, inside the document with the provided domain name, returns an error, create the document with the provided domain name if it doesn't exist
func AddSubdomain(client *mongo.Client, database string, collection string, domain string, subdomain string) error {
	// Check if the subdomain already exists, use CheckSubdomain() to check if a subdomain exists
	exists, err := CheckSubdomain(client, database, collection, domain, subdomain)
	if err != nil {
		return fmt.Errorf("[-] Error checking subdomain: %v", err)
	}
	if exists {
		return fmt.Errorf("[-] Subdomain already exists")
	}
