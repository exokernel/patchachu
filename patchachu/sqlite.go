package patchachu

import "time"

// An SQLite datastore that implements the DataStore interface
type SQLiteDataStore struct {

	// The path to the SQLite db file on disk
	//db File
	// The path to the SQLite WAL file
	//wal File

	// In memory caches can be used to speed lookups/reporting and refreshed on demand
	// Could have separate Cache structure/interface that multiple DataStores can use since
	// they all use deployments and instances.

	// Cache of the deployments
	deployments []Deployment
	// Cache of the instances
	instances []Instance

	// DataStore has an expiration so we can automatically rebuild or
	// warn the user it's old data and prompt for rebuild
	//expiresAt time
}

func NewSQLiteDataStore() *SQLiteDataStore {
	return &SQLiteDataStore{}
}

func (db *SQLiteDataStore) Build(instances []Instance, deployments []Deployment) error {
	// Create the SQLite db file if it doesn't exist
	// Create the deployments table and the instances table if they don't exist
	// Create the table mapping deployments to instances if it doesn't exist
	// Insert the deployments into the deployments table
	// Insert the instances into the instances table along with a row in the table mapping deployments to instances
	return nil
}

func (db *SQLiteDataStore) IsEmpty() bool {
	return true
}

func (db *SQLiteDataStore) Clear() error {
	return nil
}

func (db *SQLiteDataStore) InstancesForDeployment(deployment Deployment) []Instance {
	return nil
}

func (db *SQLiteDataStore) DeploymentsForInstance(instance Instance) []Deployment {
	return nil
}

func (db *SQLiteDataStore) InstancesWithNoDeployments() []Instance {
	return nil
}

func (db *SQLiteDataStore) DeploymentsWithNoInstances() []Deployment {
	return nil
}

func (db *SQLiteDataStore) IsExpired() bool {
	return false
}

func (db *SQLiteDataStore) setExpiresAt(time time.Time) {

}
