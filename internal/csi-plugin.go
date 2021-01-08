package internal

import (
	"fmt"
	nomad "github.com/hashicorp/nomad/api"
	"strconv"
)

type CsiPluginCheck struct {
	client *nomad.Client
	job    string
	plugin string
}

func NewCsiPluginCheck(client *nomad.Client, job string, plugin string) Check {
	return &CsiPluginCheck{client: client, job: job, plugin: plugin}
}

func (c *CsiPluginCheck) DoCheck() int {
	pluginInfo, _, err := c.client.CSIPlugins().Info(c.job, &nomad.QueryOptions{})

	if err != nil {
		println(err.Error())
		return 2
	}

	c.printPluginStatus(pluginInfo)
	return c.determineStatus(pluginInfo)
}

func (c *CsiPluginCheck) determineStatus(pluginInfo *nomad.CSIPlugin) int {
	if pluginInfo.NodesHealthy == 0 {
		return 2
	}

	if pluginInfo.NodesExpected != pluginInfo.NodesHealthy {
		return 1
	}

	return 0
}

func (c *CsiPluginCheck) printPluginStatus(pluginInfo *nomad.CSIPlugin) {
	healthyNodes := strconv.Itoa(pluginInfo.NodesHealthy)
	expectedNodes := strconv.Itoa(pluginInfo.NodesExpected)
	println(healthyNodes + " out of " + expectedNodes + " available")
	println()
	println(c.createJobLink())
}

func (c *CsiPluginCheck) createJobLink() string {
	link := fmt.Sprintf("%s/ui/jobs/%s", c.client.Address(), c.job)

	return "<a href=\"" + link + "\" target=\"_blank\">" + link + "</a>"
}
