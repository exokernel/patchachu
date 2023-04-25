package main

import (
	"github.com/exokernel/patchachu/patchachu"
)

func main() {
	pdb := patchachu.NewPatchastore()
	pdb.Init(&patchachu.Config{
		StoreType: "sqlite",
	})

	// Populate the datastore
	pdb.Populate()

	// Do stuff with the datastore
	instances := pdb.InstancesWithNoDeployments()

	for _, instance := range instances {
		println(instance.Name)
	}
}
