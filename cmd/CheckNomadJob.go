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
	Plugin  string `short:"p" long:"plugin" description:"The plugin to check"`
}

func main() {
	parseFlags()

	exitCode := 3
	switch opts.JobType {
	case "service":
		exitCode = internal.NewServiceCheck(nomadClient(), opts.Job).DoCheck()
	case "csi-plugin":
		exitCode = internal.NewCsiPluginCheck(nomadClient(), opts.Job, opts.Plugin).DoCheck()
	default:
		println("job type '", opts.JobType, "' not supported")
	}

	os.Exit(exitCode)
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
