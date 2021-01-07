package internal

import (
	nomad "github.com/hashicorp/nomad/api"
	"log"
	"strings"
)

type ServiceCheck struct {
	Client *nomad.Client
	Job    string

	jobInfo    *nomad.Job
	deployment *nomad.Deployment
}

func (c *ServiceCheck) Check() int {
	jobInfo, _, err := c.Client.Jobs().Info(c.Job, &nomad.QueryOptions{})

	if err != nil {
		println(err.Error())
		return 3
	}

	if jobInfo == nil {
		println("job '", c.Job, "' not found")
		return 3
	}

	c.jobInfo = jobInfo
	c.deployment, _, _ = c.Client.Jobs().LatestDeployment(*jobInfo.ID, &nomad.QueryOptions{})

	if c.deployment == nil {
		c.deployment = &nomad.Deployment{}
	}

	c.printJobInfo()
	return c.determineStatus()
}

func (c *ServiceCheck) determineStatus() int {
	if *c.jobInfo.Status != "running" {
		return 2
	}

	for _, value := range c.deployment.TaskGroups {
		if value.UnhealthyAllocs > 0 {
			return 2
		}
	}

	return 0
}

func (c *ServiceCheck) printJobInfo() {
	log.Printf("%s/ui/jobs/%s\n", c.Client.Address(), *c.jobInfo.ID)
	for key, value := range c.jobInfo.Meta {
		if strings.HasPrefix(key, "OWNER") {
			log.Printf("owner=%s\n", value)
		}
	}
	log.Println()
	log.Printf("status=%s\n", *c.jobInfo.Status)

	for key, value := range c.deployment.TaskGroups {
		log.Printf("%s.unhealthyAllocs=%d", key, value.UnhealthyAllocs)
	}
}
