package main

import (
	"bitbucket.org/sabio-it/icinga-check-nomad-job/internal"
	nomad "github.com/hashicorp/nomad/api"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
)

type args struct {
	Address string `long:"address" description:"Address of the nomad server" group:"connection" default:"http://localhost:4646"`

	CaCert     string `long:"ca" description:"Path to ca cert" group:"connection"`
	ClientCert string `long:"cert" description:"Path to client cert" group:"connection"`
	ClientKey  string `long:"key" description:"Path to client key" group:"connection"`

	Service   internal.CheckServiceOpts   `command:"service" description:"Checks a service"`
	CsiPlugin internal.CheckCsiPluginOpts `command:"csi-plugin" description:"Checks a csi plugin"`
}

var opts = new(args)

func main() {
	parser := parseFlags()

	exitCode := internal.CRITICAL
	switch parser.Active.Name {
	case "service":
		exitCode = internal.CheckService(nomadClient(), &opts.Service)
	case "csi-plugin":
		exitCode = internal.CheckCsiPlugin(nomadClient(), &opts.CsiPlugin)
	}

	os.Exit(exitCode)
}

// Parses the command line
func parseFlags() *flags.Parser {
	p := flags.NewParser(opts, flags.Default)
	_, err := p.Parse()

	if err != nil {
		log.Fatal(err.Error())
	}

	return p
}

// Creates a Nomad client from the given opts
func nomadClient() *nomad.Client {
	var tlsConfig *nomad.TLSConfig

	if opts.ClientCert != "" {
		tlsConfig = &nomad.TLSConfig{
			CACert:     opts.CaCert,
			ClientCert: opts.ClientCert,
			ClientKey:  opts.ClientKey,
			Insecure:   false,
		}

	}

	config := &nomad.Config{Address: opts.Address, TLSConfig: tlsConfig}
	client, err := nomad.NewClient(config)

	if err != nil {
		log.Fatal(err)
	}

	return client
}
