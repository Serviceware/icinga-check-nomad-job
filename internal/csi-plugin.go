package internal

import (
	nomad "github.com/hashicorp/nomad/api"
	"strconv"
	"strings"
)

type CheckCsiPluginOpts struct {
	Job                    string `short:"j" long:"job" description:"Job to check"`
	Plugin                 string `short:"p" long:"plugin" description:"Plugin to check"`
	UnhealthyNodesWarning  int    `short:"w" long:"unhealthy-nodes-warning" default:"0" description:"Number of nodes which can be unhealthy until check returns warning"`
	UnhealthyNodesCritical int    `short:"c" long:"unhealthy-nodes-critical" default:"0" description:"Number of nodes which can be unhealthy until check returns critical"`
}

// Check if csi plugin is running on expected number of nodes
func CheckCsiPlugin(client *nomad.Client, opts *CheckCsiPluginOpts) Status {
	pluginInfo, _, err := client.CSIPlugins().Info(opts.Plugin, &nomad.QueryOptions{})

	if err != nil {
		if strings.Contains(err.Error(), "404") {
			println("plugin '", opts.Plugin, "' not found")
			return CRITICAL
		}

		println(err.Error())
		return UNKNOWN
	}

	status := OK

	healthyNodes := pluginInfo.NodesHealthy
	expectedNodes := pluginInfo.NodesExpected
	unhealthyNodes := expectedNodes - healthyNodes
	println(strconv.Itoa(healthyNodes) + " out of " + strconv.Itoa(expectedNodes) + " nodes available")
	println()
	println(createJobLink(client.Address(), opts.Job))

	if unhealthyNodes > opts.UnhealthyNodesWarning {
		status = WARNING
	}

	if unhealthyNodes > opts.UnhealthyNodesCritical {
		status = CRITICAL
	}

	return status
}
