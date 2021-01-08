package test

import (
	"bitbucket.org/sabio-it/icinga-check-nomad-job/internal"
	nomad "github.com/hashicorp/nomad/api"
	"log"
	"os"
	"testing"
)

func gretaClient() *nomad.Client {
	tlsConfig := &nomad.TLSConfig{
		CACert:     "greta/ca.pem",
		ClientCert: "greta/cert.pem",
		ClientKey:  "greta/key.pem",
		Insecure:   false,
	}

	config := &nomad.Config{Address: "https://nomad01.greta-internal.hc.swops.cloud:4646", TLSConfig: tlsConfig}
	client, err := nomad.NewClient(config)

	println(os.Getwd())

	if err != nil {
		log.Fatal(err.Error())
	}

	return client
}

func TestGretaService(t *testing.T) {
	status := internal.NewServiceCheck(gretaClient(), "internal-nexus").DoCheck()

	if status != 0 {
		t.Fail()
	}
}

func TestGretaServiceUnkown(t *testing.T) {
	status := internal.NewServiceCheck(gretaClient(), "some-non-existing-job").DoCheck()

	if status != 3 {
		t.Fail()
	}
}
