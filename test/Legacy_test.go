package test

import (
	"bitbucket.org/sabio-it/icinga-check-nomad-job/internal"
	nomad "github.com/hashicorp/nomad/api"
	"log"
	"os"
	"testing"
)

func legacyClient() *nomad.Client {
	tlsConfig := &nomad.TLSConfig{
		CACert:     "legacy/ca.pem",
		ClientCert: "legacy/cert.pem",
		ClientKey:  "legacy/key.pem",
		Insecure:   false,
	}

	config := &nomad.Config{Address: "https://server01.nomadserver-internal.hc.sabio.de:4646", TLSConfig: tlsConfig}
	client, err := nomad.NewClient(config)

	println(os.Getwd())

	if err != nil {
		log.Fatal(err.Error())
	}

	return client
}

func TestLegacyService(t *testing.T) {
	serviceCheck := internal.ServiceCheck{
		Client: legacyClient(),
		Job:    "internal-registry-cleanup-cd-images",
	}

	status := serviceCheck.Check()

	if status != 0 {
		t.Fail()
	}
}

func TestLegacyServiceUnkown(t *testing.T) {
	serviceCheck := internal.ServiceCheck{
		Client: legacyClient(),
		Job:    "some-not-existing-job",
	}

	status := serviceCheck.Check()

	if status != 3 {
		t.Fail()
	}
}
