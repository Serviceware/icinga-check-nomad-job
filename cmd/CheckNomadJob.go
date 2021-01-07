package main

import (
	"bitbucket.org/sabio-it/icinga-check-nomad-job/internal"
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
		checkCsiPlugin()
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
	serviceCheck := internal.ServiceCheck{
		Client: client,
		Job:    opts.Job,
	}

	os.Exit(serviceCheck.Check())
}

func checkCsiPlugin() {
	pluginCheck := internal.CsiPluginCheck{
		Client: client,
		Job:    opts.Job,
	}
	os.Exit(pluginCheck.Check())
}
