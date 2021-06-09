package test

import (
	"github.com/Serviceware/icinga-check-nomad-job/internal"
	nomad "github.com/hashicorp/nomad/api"
	"log"
	"testing"
)

func client() *nomad.Client {
	config := &nomad.Config{Address: "http://localhost:4646"}
	client, err := nomad.NewClient(config)

	if err != nil {
		log.Fatal(err.Error())
	}

	return client
}

func TestService(t *testing.T) {
	status := internal.CheckService(client(), &internal.CheckServiceOpts{Job: "production-capi-standalone-sokadach"})

	if status != internal.OK {
		t.Fail()
	}
}

func TestSystem(t *testing.T) {
	status := internal.CheckService(client(), &internal.CheckServiceOpts{Job: "swops-production-suricata-fsn", Type: "system"})

	if status != internal.OK {
		t.Fail()
	}
}
