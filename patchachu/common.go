package patchachu

// A GCP instance
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

func (deployment *Deployment) fetchInstances() []Instance {
	// Get the instances that match all the filters
	instances := []Instance{}
	// For each instance, add the deployment to the instance
	for _, instance := range instances {
		instance.Deployments = append(instance.Deployments, *deployment)
	}
	return instances
}
