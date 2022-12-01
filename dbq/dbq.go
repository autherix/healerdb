package dbq

import (
	"fmt"

	"gopkg.in/mgo.v2"
	// log.warningf is used to print the warning messages

	// local config package
	"healerdb/config"
	"healerdb/myutils"
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

// Function CheckDatabaseExists to check if a database exists or not, return true if exists and false if not exists
func CheckDatabaseExists(session *mgo.Session, database string) bool {
	databases, err := GetDatabases(session)
	if err != nil {
		return false
	}
	for _, db := range databases {
		if db == database {
			return true
		}
	}
	return false
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

// Function to drop all databases except the default databases (admin, local, config), return error if any error occurs
func DropAllDatabases(session *mgo.Session) error {
	databases, err := GetDatabases(session)
	if err != nil {
		return err
	}
	for _, database := range databases {
		if database != "admin" && database != "local" && database != "config" {
			err = DropDatabase(session, database)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Function to create a new collection in a database, return error if any error occurs
func AddCollection(session *mgo.Session, database string, collection string) error {
	err := session.DB(database).C(collection).Insert(map[string]bool{"exists": true})
	return err
}

// Function to check if a collection exists in a database or not, return true if exists and false if not exists
func CheckCollectionExists(session *mgo.Session, database string, collection string) bool {
	collections, err := GetCollections(session, database)
	if err != nil {
		return false
	}
	for _, col := range collections {
		if col == collection {
			return true
		}
	}
	return false
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
	databases, err := config.GetDatabasesNames()
	if err != nil {
		return err
	}
	// Create the databases
	for _, database := range databases {
		// First check if the database exists
		exists := CheckDatabaseExists(session, database)
		if exists {
			fmt.Println("[INFO] Database already exists:\t'" + database + "'")
		} else {
			// Create the database
			fmt.Println("[+] Creating database:\t\t'" + database + "'")
			err = AddDatabase(session, database)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Function AddTargetToDB to add a target(as collection) to the database enum if it's not already there, return error if any error occurs
func AddTargetToDB(session *mgo.Session, target string) error {
	// Iterate over the databases in the config
	databases, err := config.GetDatabases()
	if err != nil {
		return err
	}
	for _, database := range databases {
		// If database.TargetBased is true, then
		if database.TargetBased {
			// First check if the target is already in the database
			collections, err := GetCollections(session, database.Name)
			if err != nil {
				return err
			}
			// If the target is not in the database, then add it
			if !myutils.ContainsString(collections, target) {
				// Add the target to the database
				err = AddCollection(session, database.Name, target)
				if err != nil {
					return err
				}
				fmt.Println("[+] Added target to database:\t'" + target + "'")
			} else {
				fmt.Println("[INFO] Target already exists in database:\t'" + target + "'")
			}
		}
	}
	return nil
}

// Remove target from the database
func RemoveTargetFromDB(session *mgo.Session, target string) error {
	// Iterate over the databases in the config
	databases, err := config.GetDatabases()
	if err != nil {
		return err
	}
	for _, database := range databases {
		// If database.TargetBased is true, then
		if database.TargetBased {
			// First check if the target is already in the database
			collections, err := GetCollections(session, database.Name)
			if err != nil {
				return err
			}
			// If the target is in the database, then remove it
			if myutils.ContainsString(collections, target) {
				// Remove the target from the database
				err = DropCollection(session, database.Name, target)
				if err != nil {
					return err
				}
				fmt.Println("[+] Removed target:\t'" + target + "' from database:\t'" + database.Name + "'")
			} else {
				fmt.Println("[INFO] Target:\t'" + target + "' does not exist in database:\t'" + database.Name + "'")
			}
		}
	}
	return nil
}
