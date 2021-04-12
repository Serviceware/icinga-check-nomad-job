package internal

import (
	nomad "github.com/hashicorp/nomad/api"
	"strconv"
	"strings"
)

type CheckServiceOpts struct {
	Job string `long:"job" short:"j" description:"The job to check"`
}

// Checks if the given job is running and if it has unhealthy allocation.
func CheckService(client *nomad.Client, opts *CheckServiceOpts) Status {
	jobInfo, _, err := client.Jobs().Info(opts.Job, &nomad.QueryOptions{})

	if err != nil {
		if strings.Contains(err.Error(), "404") {
			println("job '", opts.Job, "' not found")
			return CRITICAL
		}

		println(err.Error())
		return UNKNOWN
	}

	summary, _, _ := client.Jobs().Summary(*jobInfo.ID, nil)
	deployment, _, _ := client.Jobs().LatestDeployment(*jobInfo.ID, &nomad.QueryOptions{})

	if deployment == nil {
		deployment = &nomad.Deployment{}
	}

	status := OK
	println("status=" + *jobInfo.Status)

	if *jobInfo.Status != "running" {
		status = CRITICAL
	}

	for key, value := range deployment.TaskGroups {
		running := summary.Summary[key].Running
		desired := value.DesiredTotal

		if running < desired {
			println(key+".summary.running="+strconv.Itoa(running), "<", key+".desiredTotal="+strconv.Itoa(desired))
			status = CRITICAL
		} else {
			println(key+".summary.running="+strconv.Itoa(running), ">=", key+".desiredTotal="+strconv.Itoa(desired))
		}

	}

	println(createJobLink(client.Address(), *jobInfo.ID))
	for key, value := range jobInfo.Meta {
		if strings.HasPrefix(key, "OWNER") {
			println("owner=" + value)
		}
	}

	return status
}
