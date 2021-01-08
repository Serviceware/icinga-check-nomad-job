package internal

import (
	"fmt"
	nomad "github.com/hashicorp/nomad/api"
	"strconv"
	"strings"
)

type ServiceCheck struct {
	client *nomad.Client
	job    string

	jobInfo    *nomad.Job
	deployment *nomad.Deployment
}

func NewServiceCheck(client *nomad.Client, job string) Check {
	return &ServiceCheck{
		client: client,
		job:    job,
	}
}

func (c *ServiceCheck) DoCheck() int {
	jobInfo, _, err := c.client.Jobs().Info(c.job, &nomad.QueryOptions{})

	if err != nil {
		println(err.Error())
		return 3
	}

	if jobInfo == nil {
		println("job '", c.job, "' not found")
		return 3
	}

	c.jobInfo = jobInfo
	c.deployment, _, _ = c.client.Jobs().LatestDeployment(*jobInfo.ID, &nomad.QueryOptions{})

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
	println("status=" + *c.jobInfo.Status)

	for key, value := range c.deployment.TaskGroups {
		println(key + ".unhealthyAllocs=" + strconv.Itoa(value.UnhealthyAllocs))
	}

	println(c.createJobLink())
	for key, value := range c.jobInfo.Meta {
		if strings.HasPrefix(key, "OWNER") {
			println("owner=" + value)
		}
	}
}

func (c *ServiceCheck) createJobLink() string {
	link := fmt.Sprintf("%s/ui/jobs/%s", c.client.Address(), *c.jobInfo.ID)

	return "<a href=\"" + link + "\" target=\"_blank\">" + link + "</a>"
}
