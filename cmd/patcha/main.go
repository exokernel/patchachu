package main

import (
	"github.com/exokernel/patchachu/patchachu"
)

func main() {
	pdb := patchachu.NewPatchastore()
	pdb.Init(&patchachu.Config{
		StoreType: "sqlite",
		Projects:  []string{"integration", "staging", "production"}, // TODO: Make this a command line argument
	})

	// Populate the datastore with patch info from each project
	pdb.Populate()

	// Do stuff with the datastore
	instances := pdb.InstancesWithNoDeployments()

	for _, instance := range instances {
		println(instance.Name)
	}
}
