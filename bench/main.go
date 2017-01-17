package main

import (
	"io/ioutil"
	"log"

	"github.com/vincentserpoul/playwithsql/infra/server"
	"github.com/vincentserpoul/playwithsql/infra/server/digitalocean"
)

func main() {

	pat, err := ioutil.ReadFile("./dotoken")
	if err != nil {
		log.Fatalf("getDOClient: no digital ocean token file found")
	}

	sshKeyFingerprint, err := ioutil.ReadFile("./sshkey_fingerprint")
	if err != nil {
		log.Fatalf("getDOClient: no digital ocean sshkey_fingerprint file found")
	}

	provider := digitalocean.New(string(pat), string(sshKeyFingerprint))

	err = server.LaunchBenches(provider)
	if err != nil {
		log.Fatalf("%v", err)
	}
}
