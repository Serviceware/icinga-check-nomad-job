package main

import (
	"bitbucket.org/sabio-it/icinga-check-nomad-job/internal"
	"fmt"
	nomad "github.com/hashicorp/nomad/api"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
)

var opts struct {
	Address string `long:"address" description:"Address of the nomad server"`

	CaCert     string `long:"ca" description:"Path to ca cert"`
	ClientCert string `long:"cert" description:"Path to client cert"`
	ClientKey  string `long:"key" description:"Path to client key"`

	Job     string `short:"j" long:"job" description:"Job to check"`
	JobType string `short:"t" long:"type" description:"Type of the job (service, csi-plugin)"`
}

var client *nomad.Client

func main() {
	parseFlags()
	client = nomadClient()

	switch opts.JobType {
	case "service":
		checkService()
	case "csi-plugin":
		checkSystem()
	default:
		println("job type '", opts.JobType, "' not supported")
	}
}

func parseFlags() {
	_, err := flags.Parse(&opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}

func nomadClient() *nomad.Client {
	tlsConfig := &nomad.TLSConfig{
		CACert:     opts.CaCert,
		ClientCert: opts.ClientCert,
		ClientKey:  opts.ClientKey,
		Insecure:   false,
	}

	config := &nomad.Config{Address: opts.Address, TLSConfig: tlsConfig}
	client, err := nomad.NewClient(config)

	if err != nil {
		log.Fatal(err)
	}

	return client
}

func checkService() {
}

func TestMe() {}

func checkSystem() {
	pluginInfo, _, err := client.CSIPlugins().Info(opts.Job, &nomad.QueryOptions{})

	if err != nil {
		println(err.Error())
		os.Exit(2)
	}

	printPluginStatus(pluginInfo)

	if pluginInfo.NodesHealthy == 0 {
		os.Exit(2)
	}

	if pluginInfo.NodesExpected != pluginInfo.NodesHealthy {
		os.Exit(1)
	}

	os.Exit(0)
}

func printPluginStatus(pluginInfo *nomad.CSIPlugin) {
	fmt.Printf("nodes-healthy  = %d\n", pluginInfo.NodesHealthy)
	fmt.Printf("nodes-expected = %d\n", pluginInfo.NodesExpected)
}
