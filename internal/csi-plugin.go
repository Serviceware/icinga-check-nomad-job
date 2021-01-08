package internal

import (
	"fmt"
	nomad "github.com/hashicorp/nomad/api"
	"strconv"
	"strings"
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
	pluginInfo, _, err := c.client.CSIPlugins().Info(c.plugin, &nomad.QueryOptions{})

	if err != nil {
		if strings.Contains(err.Error(), "404") {
			println("plugin '", c.plugin, "' not found")
			return 2
		}

		println(err.Error())
		return 3
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
	println(healthyNodes + " out of " + expectedNodes + " nodes available")
	println()
	println(c.createJobLink())
}

func (c *CsiPluginCheck) createJobLink() string {
	link := fmt.Sprintf("%s/ui/jobs/%s", c.client.Address(), c.job)

	return "<a href=\"" + link + "\" target=\"_blank\">" + link + "</a>"
}
