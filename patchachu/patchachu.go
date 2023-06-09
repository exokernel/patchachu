package patchachu

import (
	"context"
	"fmt"
	"time"

	osconfig "cloud.google.com/go/osconfig/apiv1"
	osconfigpb "cloud.google.com/go/osconfig/apiv1/osconfigpb"
	"google.golang.org/api/iterator"
)

// An interface for a datastore
// The concrete implementation of this interface could be SQLite, flat file, etc.
type DataStore interface {
	// Check if the data is empty
	IsEmpty() bool

	// Seed the data
	// After this call the datastore should be ready to answer queries
	// and isEmpty() should return false
	Build(instances []Instance, deployments []Deployment) error

	// Clear the data e.g. drop all tables in sqlite
	// After this call the datastore should be empty and isEmpty() should return true
	Clear() error

	// Completely destroy the underlying datastore e.g. rm sqlite db file
	// It's good to clean up after ourselves
	// After this call the datastore should be empty and isEmpty() should return true
	//Destroy() error

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
	store  DataStore
	config *Config

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
func (pdb *Patchastore) fetch(project string) error {
	// Get the all the deployments and then all the instances
	// These API calls could be done concurrently
	pdb.fetchDeployments(project)
	pdb.fetchInstances(project)

	// Now we link the instances and deployments together by making them reference each other
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
// https://github.com/googleapis/google-cloud-go/wiki/Iterator-Guidelines
// https://pkg.go.dev/cloud.google.com/go/osconfig@v1.11.1/apiv1#Client.ListPatchDeployments
// TODO: fully convert osconfigpb.PatchDeployment to Deployment and append to pdb.deployments
func (pdb *Patchastore) fetchDeployments(project string) error {
	ctx := context.Background()

	// Create a new client.
	c, err := osconfig.NewClient(ctx)
	if err != nil {
		fmt.Printf("Failed to create client: %v", err)
		return err
	}
	defer c.Close()

	// Create the initial request to list patch deployments.
	req := &osconfigpb.ListPatchDeploymentsRequest{
		Parent: project,
	}

	// Iterate over all paginated results.
	it := c.ListPatchDeployments(ctx, req)
	p := iterator.NewPager(it, 50, "") // TODO: make page size configurable
	for {
		var deployments []*osconfigpb.PatchDeployment
		nextPageToken, err := p.NextPage(&deployments)
		if err != nil {
			fmt.Println("Failed to get patch deployments:", err)
			return err
		}
		for _, deployment := range deployments {
			fmt.Println("Deployment:", deployment.GetName())
			// new deployment from osconfigpb.PatchDeployment
			pdbDeployment := Deployment{
				Name: deployment.GetName(),
			}
			// append to pdb.deployments
			pdb.deployments = append(pdb.deployments, pdbDeployment)
		}
		if nextPageToken == "" {
			break
		}
	}

	return nil
}

func (pdb *Patchastore) fetchInstances(project string) error {
	pdb.instances = []Instance{}
	return nil
}

// What?
// Populate the in-memory cache and build the persistent datastore if necessary
// Else build the in-memory cache from the datastore
//
// Why?
// This is the heavy lifting necessary to get patchachu ready to answer queries from its in-memory cache and persist that data
// to disk so that it doesn't have to be fetched from GCP on subsequent runs (at least until the data expires)
// It's a one-time cost that is paid in full before the first query
//
// Reasoning:
// We want to avoid making API calls if possible, because they're slow, hence the persistent datastore
// We want to use the in-memory cache if possible, because it's faster than queyring the datastore
//
// If the datastore is empty or expired our strategy is to do the API calls once building our cache as we go and then
// save the cached data in the persistent datastore. We don't need to do the API calls again until the data in the datastore expires.
//
// If the datastore is not empty and has not expired, we can use the data in the datastore to build our in-memory cache.
//
// TODO: Error handling
func (pdb *Patchastore) Populate() {
	// If the datastore is not empty and has expired, clear it and reset the expiration time
	if !pdb.store.IsEmpty() && pdb.store.IsExpired() {
		println("Datastore is not empty and has expired, clearing it and resetting the expiration time")
		pdb.store.Clear()
		pdb.store.setExpiresAt(time.Now().Add(24 * time.Hour))
	}
	// If the datastore is empty and has not expired, populate it
	if !pdb.store.IsEmpty() && pdb.store.IsExpired() {
		println("Datastore is empty and has not expired, populating it")
		// Fetch the data with API calls and build the in-memory cache from the fetched data
		for _, project := range pdb.config.Projects {
			pdb.fetch(project)
		}

		// Build the datastore from memory-cached instances and deployments
		pdb.store.Build(pdb.instances, pdb.deployments)
	} else if !pdb.store.IsEmpty() && !pdb.store.IsExpired() {
		// Great news everyone! We can use the locally stored data!
		// The datastore is not empty and has not expired, so we need to build our in-memory cache from the datastore
		//pdb.instances, pdb.deployments = pdb.store.InstancesAndDeployments()
		println("TODO: Build the in-memory cache from the datastore")
	}
	// we don't need to do anything if the datastore is empty and has expired
	println("nothing to do")
}

func (pdb *Patchastore) Init(config *Config) {
	println("Patcha! Initializing the Patchastore!")
	pdb.config = config

	if config.StoreType == "sqlite" {
		pdb.store = NewSQLiteDataStore()
	}
}

func (pdb *Patchastore) InstancesWithNoDeployments() []Instance {
	return nil
}
