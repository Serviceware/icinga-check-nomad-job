package internal

import (
	"fmt"
	nomad "github.com/hashicorp/nomad/api"
)

type CsiPluginCheck struct {
	Client *nomad.Client
	Job    string
}

func (c *CsiPluginCheck) Check() int {
	pluginInfo, _, err := c.Client.CSIPlugins().Info(c.Job, &nomad.QueryOptions{})

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
	fmt.Printf("nodes-healthy  = %d\n", pluginInfo.NodesHealthy)
	fmt.Printf("nodes-expected = %d\n", pluginInfo.NodesExpected)
}
