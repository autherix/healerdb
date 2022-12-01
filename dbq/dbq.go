package dbq

import (
	"gopkg.in/mgo.v2"

	// local config package
	"healerdb/config"
)

// Function Connect: to connect to the database with the given connection url (e.g. mongodb://localhost:27017), return the session and error if any occurs
func Connect(connstr string) (*mgo.Session, error) {
	// Create a session for the database
	session, err := mgo.Dial(connstr)
	if err != nil {
		panic(err)
	}

	// Check the session to the database is alive or not and ping the database
	err = session.Ping()
	if err != nil {
		panic(err)
	}

	// print that the session created and active
	println("Session created and active")

	// Return the session and error if any occurs
	return session, err
}

// Function to add a new database, and create a new collection in the database called "exists", and create a document in the collection called "exists" with the value "true", return error if any error occurs
func AddDatabase(session *mgo.Session, database string) error {
	err := session.DB(database).C("exists").Insert(map[string]bool{"exists": true})
	return err
}

// Function to drop a database, return error if any error occurs
func DropDatabase(session *mgo.Session, database string) error {
	err := session.DB(database).DropDatabase()
	return err
}

// Function to create a new collection in a database, return error if any error occurs
func AddCollection(session *mgo.Session, database string, collection string) error {
	err := session.DB(database).C(collection).Insert(map[string]bool{"exists": true})
	return err
}

// Function to drop a collection in a database, return error if any error occurs
func DropCollection(session *mgo.Session, database string, collection string) error {
	err := session.DB(database).C(collection).DropCollection()
	return err
}

// Function to drop all collections in a database, return error if any error occurs
func DropAllCollections(session *mgo.Session, database string) error {
	collections, err := GetCollections(session, database)
	if err != nil {
		return err
	}
	for _, collection := range collections {
		err = DropCollection(session, database, collection)
		if err != nil {
			return err
		}
	}
	return nil
}

// Function to add a document to a collection in a database, return error if any error occurs
func AddDocument(session *mgo.Session, database string, collection string, document interface{}) error {
	err := session.DB(database).C(collection).Insert(document)
	return err
}

// Function to delete all documents in a collection in a database with the provided query, delete all if no query was provided, return error if any error occurs
func DeleteDocumentsWithQuery(session *mgo.Session, database string, collection string, query interface{}) error {
	_, err := session.DB(database).C(collection).RemoveAll(query)
	return err
}

// Function to return all the databases in the database, return list of databases and error
func GetDatabases(session *mgo.Session) ([]string, error) {
	return session.DatabaseNames()
}

// Function to return all the collections in a database, return list of collections and error
func GetCollections(session *mgo.Session, database string) ([]string, error) {
	return session.DB(database).CollectionNames()
}

// Function to return all the documents in a collection, return list of documents and error
func GetDocuments(session *mgo.Session, database string, collection string) ([]interface{}, error) {
	var result []interface{}
	err := session.DB(database).C(collection).Find(nil).All(&result)
	return result, err
}

// Function to return all the documents with the provided query in a collection, return list of documents and error
func GetDocumentsWithQuery(session *mgo.Session, database string, collection string, query interface{}) ([]interface{}, error) {
	var result []interface{}
	err := session.DB(database).C(collection).Find(query).All(&result)
	return result, err
}

// Function to replace all documents in a collection in a database with the provided query, return error if any error occurs
func ReplaceDocumentsWithQuery(session *mgo.Session, database string, collection string, query interface{}, replacement interface{}) error {
	_, err := session.DB(database).C(collection).RemoveAll(query)
	if err != nil {
		return err
	}
	err = session.DB(database).C(collection).Insert(replacement)
	return err
}

// Function to update all documents in a collection in a database with the provided query, update all if no query was provided, return error if any error occurs
func UpdateDocumentsWithQuery(session *mgo.Session, database string, collection string, query interface{}, update interface{}) error {
	_, err := session.DB(database).C(collection).UpdateAll(query, update)
	return err
}

// Function to return the number of databases available
func CountDatabases(session *mgo.Session) (int, error) {
	databases, err := GetDatabases(session)
	return len(databases), err
}

// Function to return the number of collections in a database
func CountCollections(session *mgo.Session, database string) (int, error) {
	collections, err := GetCollections(session, database)
	return len(collections), err
}

// Function to return the number of documents in a collection in a database with the provided query, return number of documents and error
func CountDocumentsWithQuery(session *mgo.Session, database string, collection string, query interface{}) (int, error) {
	count, err := session.DB(database).C(collection).Find(query).Count()
	return count, err
}

// First data structure setup
/*
layer 1: database -> Each represents a special info type (e.g. enum, vuln, watch, etc.)

For Example in 'enum' database:
layer 2: collection -> Each represents a target handle as its username (e.g. "semrush", "uber", etc)
layer 3: document (JSON Object) -> Each represents a target's domain (e.g. "semrush.com", "uber.com", etc)
layer 4: Subdomain/RAW IP -> Each represents a target's subdomain (e.g. "www.semrush.com", "www.uber.com", etc)
Layer 5: directory -> Each represents a target's directory (e.g. "/admin", "/login", etc)
Layer 6: SubDirectory/File -> Each represents a target's subdirectory or file (e.g. "/admin/login", "/admin/login.php", etc)
Layer 7: Parameters -> Each represents a target's parameters (e.g. "/admin/login?username=admin&password=admin", etc) --> This is the final layer
*/

// Function DBFirstSetup to setup the first layer of the database (databases creation), Read the databases' names from the config file,->healerdb->dbs (is a list, select the 'name' field from the list), return error if any error occurs, Remove any other test or useless databases
func DBFirstSetup(session *mgo.Session) error {
	// Get the databases' names from the config file
	databases, err := config.GetDatabases()
	if err != nil {
		return err
	}
	// Create the databases
	for _, database := range databases {
		err = AddDatabase(session, database)
		if err != nil {
			return err
		}
	}
	return nil
}
