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
	status := internal.NewServiceCheck(legacyClient(), "internal-registry-cleanup-cd-images").DoCheck()

	if status != 0 {
		t.Fail()
	}
}

func TestLegacyServiceUnkown(t *testing.T) {
	status := internal.NewServiceCheck(legacyClient(), "some-non-existing-job").DoCheck()

	if status != 3 {
		t.Fail()
	}
}

func TestLegacyPlugin(t *testing.T) {
	status := internal.NewCsiPluginCheck(legacyClient(), "gluster-csi-plugin", "csi.serviceware.gluster").DoCheck()

	if status != 0 {
		t.Fail()
	}
}