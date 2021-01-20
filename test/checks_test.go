package test

import (
	"bitbucket.org/sabio-it/icinga-check-nomad-job/internal"
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

func Test(t *testing.T) {
	status := internal.CheckService(client(), &internal.CheckServiceOpts{Job: "example"})

	if status != internal.OK {
		t.Fail()
	}
}
