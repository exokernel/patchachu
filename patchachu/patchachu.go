package patchachu

type Instance struct {
	Name string
	//IP          IpAddress
	Tags        []string
	Project     string
	Region      string
	Zone        string
	Deployments []Deployment
}

// A GCP patch deployment
type Deployment struct {
	Name string
	//Filters   []Filter
	Project   string
	Instances []Instance
}

// An interface for a datastore
// The concrete implementation of this interface could be SQLite, flat file, etc.
type DataStore interface {
	// Check if the data is empty
	isEmpty() bool

	// Seed the data
	// After this call the datastore should be ready to answer queries
	// and isEmpty() should return false
	build(instances []Instance, deployments []Deployment) error

	// Clear the data
	// After this call the datastore should be empty and isEmpty() should return true
	clear() error

	// Get all the instances that are covered by the deployment
	// Note that if a deployment has no instances it means that the filters aren't
	// matching any instances. This is is not ideal and should be fixed.
	instancesForDeployment(deployment Deployment) []Instance

	// Get all the deployments whose filters match the instance
	// Note that if an instance is covered by multiple deployments it generally
	// means that we have two different deployments that are overlapping, which is
	// not ideal.
	deploymentsForInstance(instance Instance) []Deployment

	// Get all the instances that have no deployments
	instancesWithNoDeployments() []Instance

	// Get all the deployments that have no instances
	deploymentsWithNoInstances() []Deployment

	// Set the expiration time
	//setExpiresAt(time time)

	// Get the time the data expires
	//expiresAt() time

	// Check if the data is expired
	isExpired() bool

	// Generate a report in CSV format to the given file
	// The file will be created if it doesn't exist
	// The file might be a pipe to stdout
	//csvReport(reportFile File) error

	// Generate a report in JSON format to the given file
	// The file will be created if it doesn't exist
	// The file might be a pipe to stdout
	//jsonReport(reportFile File) error
}

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

type Patchastore struct {
	// The datastore
	store DataStore
}

func (pdb *Patchastore) New() *Patchastore {
	return &Patchastore{
		store: nil,
	}
}

func (pdb *Patchastore) Init() {
	println("Patcha! Initializing the Patchastore!")
}
