package main

import (
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

func main() {
	parseFlags()

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
	client := nomadClient()

	jobInfo, _, err := client.Jobs().Info(opts.Job, &nomad.QueryOptions{})

	if err != nil {
		println(err.Error())
		os.Exit(3)
	}

	if jobInfo == nil {
		println("job '", opts.Job, "' not found")
		os.Exit(3)
	}

	switch *jobInfo.Status {
	case "running":
		println("status=running")
		os.Exit(0)
	case "pending":
		println("status=pending")
		os.Exit(1)
	case "dead":
		println("status=dead")
		os.Exit(2)
	default:
		println("unhandled status '", *jobInfo.Status, "'")
		os.Exit(3)
	}
}

func checkSystem() {
	client := nomadClient()

	pluginInfo, _, err := client.CSIPlugins().Info(opts.Job, &nomad.QueryOptions{})

	if err != nil {
		println(err.Error())
		os.Exit(2)
	}

	if pluginInfo.NodesHealthy == 0 {
		fmt.Printf("nodes-healthy  = 0\n")
		fmt.Printf("nodes-expected = %d\n", pluginInfo.NodesExpected)
		os.Exit(2)
	}

	if pluginInfo.NodesExpected != pluginInfo.NodesHealthy {
		fmt.Printf("nodes-healthy  = %d\n", pluginInfo.NodesHealthy)
		fmt.Printf("nodes-expected = %d\n", pluginInfo.NodesExpected)
		os.Exit(1)
	}

	os.Exit(0)
}
