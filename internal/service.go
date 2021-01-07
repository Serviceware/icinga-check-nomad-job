package internal

import (
	nomad "github.com/hashicorp/nomad/api"
	"log"
	"strings"
)

type ServiceCheck struct {
	client *nomad.Client
	job    string
}

func (c *ServiceCheck) CheckService() int {
	jobInfo, _, err := c.client.Jobs().Info(c.job, &nomad.QueryOptions{})

	if err != nil {
		println(err.Error())
		return 3
	}

	if jobInfo == nil {
		println("job '", c.job, "' not found")
		return 3
	}

	deployment, _, _ := c.client.Jobs().LatestDeployment(*jobInfo.ID, &nomad.QueryOptions{})

	c.printJobInfo(jobInfo, deployment)
	return determineStatus(jobInfo, deployment)
}

func determineStatus(jobInfo *nomad.Job, deployment *nomad.Deployment) int {
	if *jobInfo.Status != "running" {
		return 2
	}

	for _, value := range deployment.TaskGroups {
		if value.UnhealthyAllocs > 0 {
			return 2
		}
	}

	return 0
}

func (c *ServiceCheck) printJobInfo(jobInfo *nomad.Job, deployment *nomad.Deployment) {
	log.Printf("%s/ui/jobs/%s\n", c.client.Address(), jobInfo.ID)
	log.Println()
	log.Printf("status=%s\n", *jobInfo.Status)

	for key, value := range deployment.TaskGroups {
		log.Printf("%s.unhealthyAllocs=%d", key, value.UnhealthyAllocs)
	}

	for key, value := range jobInfo.Meta {
		if strings.HasPrefix(key, "OWNER") {
			log.Printf("owner=%s\n", value)
		}
	}
}
