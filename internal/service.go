package internal

import (
	nomad "github.com/hashicorp/nomad/api"
	"strconv"
	"strings"
)

type CheckServiceOpts struct {
	Job string `long:"job"' short:"j" description:"The job to check"`
}

// Checks if the given job is running and if it has unhealthy allocation.
func CheckService(client *nomad.Client, opts *CheckServiceOpts) int {
	jobInfo, _, err := client.Jobs().Info(opts.Job, &nomad.QueryOptions{})

	if err != nil {
		if strings.Contains(err.Error(), "404") {
			println("job '", opts.Job, "' not found")
			return CRITICAL
		}

		println(err.Error())
		return UNKNOWN
	}

	deployment, _, _ := client.Jobs().LatestDeployment(*jobInfo.ID, &nomad.QueryOptions{})

	if deployment == nil {
		deployment = &nomad.Deployment{}
	}

	code := OK
	println("status=" + *jobInfo.Status)

	if *jobInfo.Status != "running" {
		code = CRITICAL
	}

	for key, value := range deployment.TaskGroups {
		println(key + ".unhealthyAllocs=" + strconv.Itoa(value.UnhealthyAllocs))

		if value.UnhealthyAllocs > 0 {
			code = CRITICAL
		}
	}

	println(createJobLink(client.Address(), *jobInfo.ID))
	for key, value := range jobInfo.Meta {
		if strings.HasPrefix(key, "OWNER") {
			println("owner=" + value)
		}
	}

	return code
}



