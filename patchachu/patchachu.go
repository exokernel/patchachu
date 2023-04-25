package patchachu

import "time"

// An interface for a datastore
// The concrete implementation of this interface could be SQLite, flat file, etc.
type DataStore interface {
	// Check if the data is empty
	IsEmpty() bool

	// Seed the data
	// After this call the datastore should be ready to answer queries
	// and isEmpty() should return false
	Build(instances []Instance, deployments []Deployment) error

	// Clear the data
	// After this call the datastore should be empty and isEmpty() should return true
	Clear() error

	// Get all the instances that are covered by the deployment
	// Note that if a deployment has no instances it means that the filters aren't
	// matching any instances. This is is not ideal and should be fixed.
	InstancesForDeployment(deployment Deployment) []Instance

	// Get all the deployments whose filters match the instance
	// Note that if an instance is covered by multiple deployments it generally
	// means that we have two different deployments that are overlapping, which is
	// not ideal.
	DeploymentsForInstance(instance Instance) []Deployment

	// Get all the instances that have no deployments
	InstancesWithNoDeployments() []Instance

	// Get all the deployments that have no instances
	DeploymentsWithNoInstances() []Deployment

	// Set the expiration time
	setExpiresAt(time time.Time)

	// Get the time the data expires
	//expiresAt() time

	// Check if the data is expired
	IsExpired() bool

	// Generate a report in CSV format to the given file
	// The file will be created if it doesn't exist
	// The file might be a pipe to stdout
	//csvReport(reportFile File) error

	// Generate a report in JSON format to the given file
	// The file will be created if it doesn't exist
	// The file might be a pipe to stdout
	//jsonReport(reportFile File) error
}

type Patchastore struct {
	// The datastore
	store DataStore

	instances   []Instance
	deployments []Deployment
}

// Create a new Patchastore
func NewPatchastore() *Patchastore {
	return &Patchastore{
		store: nil,
	}
}

// Do the API calls to get the data into instances and deployments
func (pdb *Patchastore) fetch() error {
	// Get the all the deployments and then all the instances
	pdb.fetchDeployments()
	pdb.fetchInstances()
	// For each deployment, get the instances covered by the deployment
	for _, deployment := range pdb.deployments {
		instances := deployment.fetchInstances()
		for _, instance := range instances {
			// if this instance is already in the list of instances, add a pointer to it to the deployment
			// if this instance is not in the list of instances, add it to the list and add a pointer to it to the deployment
			// also it's very strange to find an instance not already in the list of all instances so we should log that
			if !instanceInList(instance, pdb.instances) {
				pdb.instances = append(pdb.instances, instance)
				deployment.Instances = append(deployment.Instances, instance)
				// log that we found an instance not already in the list of all instances
				println("Found an instance not already in the list of all instances: %v", instance.Name)
			} else {
				// add a pointer to the instance to the deployment's list of instances
				deployment.Instances = append(deployment.Instances, instance)
			}
		}
	}
	return nil
}

func instanceInList(instance Instance, list []Instance) bool {
	for _, i := range list {
		if i.Name == instance.Name {
			return true
		}
	}
	return false
}

// Do the API calls to get the deployment data
func (pdb *Patchastore) fetchDeployments() error {
	pdb.deployments = []Deployment{}
	return nil
}

func (pdb *Patchastore) fetchInstances() error {
	pdb.instances = []Instance{}
	return nil
}

// Populate the datastore if necessary
func (pdb *Patchastore) Populate() error {
	// If the datastore is not empty and has expired, clear it and reset the expiration time
	if !pdb.store.IsEmpty() && pdb.store.IsExpired() {
		pdb.store.Clear()
		pdb.store.setExpiresAt(time.Now().Add(24 * time.Hour))
	}
	// If the datastore is empty and has not expired, populate it
	if pdb.store.IsEmpty() && !pdb.store.IsExpired() {
		// Fetch the data
		pdb.fetch()

		// Build the datastore from instances and deployments
		pdb.store.Build(pdb.instances, pdb.deployments)
	}
	return nil
}

func (pdb *Patchastore) Init(config *Config) {
	println("Patcha! Initializing the Patchastore!")

	if config.StoreType == "sqlite" {
		pdb.store = NewSQLiteDataStore()
	}
}

func (pdb *Patchastore) InstancesWithNoDeployments() []Instance {
	return nil
}
