package main

import (
	"flag"
	nomad "github.com/hashicorp/nomad/api"
	"log"
	"os"
)

var (
	address    string
	caCert     string
	clientCert string
	clientKey  string

	job     string
	jobType string
)

func main() {
	parseFlags()

	switch jobType {
	case "service":
		checkService()
	case "csi-plugin":
		checkSystem()
	default:
		println("job type '", jobType, "' not supported")
	}
}

func parseFlags() {
	flag.StringVar(&address, "address", "", "Address of nomad server")

	flag.StringVar(&caCert, "ca-cert", "", "Path to ca cert")
	flag.StringVar(&clientCert, "client-cert", "", "Path to client cert")
	flag.StringVar(&clientKey, "client-key", "", "Path to client key")

	flag.StringVar(&job, "job", "", "The nomad job to check")
	flag.StringVar(&jobType, "type", "service", "The type of the job (service, csi-plugin)")

	flag.Parse()
}

func nomadClient() *nomad.Client {
	tlsConfig := &nomad.TLSConfig{
		CACert:     caCert,
		ClientCert: clientCert,
		ClientKey:  clientKey,
		Insecure:   false,
	}

	config := &nomad.Config{Address: address, TLSConfig: tlsConfig}
	client, err := nomad.NewClient(config)

	if err != nil {
		log.Fatal(err)
	}

	return client
}

func checkService() {
	client := nomadClient()

	jobInfo, _, err := client.Jobs().Info(job, &nomad.QueryOptions{})

	if err != nil {
		println(err.Error())
		os.Exit(3)
	}

	if jobInfo == nil {
		println("job '", job, "' not found")
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

	//jobInfo, _, err := client.Jobs().Info(job, &nomad.QueryOptions{})
	pluginInfo, _, err := client.CSIPlugins().Info(job, &nomad.QueryOptions{})

	if err != nil {
		println(err.Error())
		os.Exit(3)
	}

	if pluginInfo.NodesHealthy == 0 {
		println("nodes-healthy=0")
		os.Exit(2)
	}

	if pluginInfo.NodesExpected != pluginInfo.NodesHealthy {
		println("nodes-healthy=", pluginInfo.NodesHealthy)
		os.Exit(1)
	}

	os.Exit(0)
}
