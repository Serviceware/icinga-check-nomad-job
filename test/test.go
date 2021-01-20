package test

import (
	nomad "github.com/hashicorp/nomad/api"
	"log"
	"os"
)

func client() *nomad.Client {
	config := &nomad.Config{Address: "http://localhost:4646"}
	client, err := nomad.NewClient(config)

	println(os.Getwd())

	if err != nil {
		log.Fatal(err.Error())
	}

	return client
}

//func TestGretaService(t *testing.T) {
//	status := internal.NewServiceCheck(gretaClient(), "internal-nexus").DoCheck()
//
//	if status != 0 {
//		t.Fail()
//	}
//}
//
//func TestGretaServiceUnkown(t *testing.T) {
//	status := internal.NewServiceCheck(gretaClient(), "some-non-existing-job").DoCheck()
//
//	if status != 3 {
//		t.Fail()
//	}
//}
