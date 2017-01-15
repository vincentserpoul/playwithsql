package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"

	"time"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

// TokenSource contains the token, used for oauth2
type TokenSource struct {
	AccessToken string
}

func main() {

	pat, err := ioutil.ReadFile("./dotoken")
	if err != nil {
		log.Fatal("no dotoken file found")
	}

	SSHKeyFingerprint, err := ioutil.ReadFile("./sshkey_fingerprint")
	if err != nil {
		log.Fatal("no ssh fingerprint file")
	}

	tokenSource := &TokenSource{
		AccessToken: string(pat),
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)

	// Delete old Droplets
	err = DeletePreviousCoreOSDroplets(client)
	if err != nil {
		log.Fatalf("destroy previous droplets: %v", err)
	}

	// Create DB1
	CoreOSDB1, err := CreateCoreOSDroplet(client, "CoreOSDB1", string(SSHKeyFingerprint))
	if err != nil {
		log.Fatalf("could not create db droplet: %v", err)
	}

	// Create Bench
	CoreOSBench, err := CreateCoreOSDroplet(client, "CoreOSBench", string(SSHKeyFingerprint))
	if err != nil {
		log.Fatalf("could not create bench droplet: %v", err)
	}

	CoreOSDB1Active := false
	CoreOSBenchActive := false
	var errGet error
	for !CoreOSDB1Active && !CoreOSBenchActive {
		time.Sleep(time.Second * 15)
		fmt.Println("waiting for droplet to be created")
		CoreOSDB1, _, errGet = client.Droplets.Get(CoreOSDB1.ID)
		if errGet != nil {
			log.Fatalf("could not refresh db1 droplet status: %v", errGet)
		}
		CoreOSBench, _, errGet = client.Droplets.Get(CoreOSDB1.ID)
		if errGet != nil {
			log.Fatalf("could not refresh bench droplet status: %v", errGet)
		}

		CoreOSDB1Active = CoreOSDB1.Status == "active"
		fmt.Println("CoreOSDB1 status: ", CoreOSDB1.Status)
		CoreOSBenchActive = CoreOSBench.Status == "active"
		fmt.Println("CoreOSBench status: ", CoreOSBench.Status)
	}

	ipdb1, err := CoreOSDB1.PublicIPv4()
	if err != nil {
		log.Fatalf("could not get db1 ip: %v", err)
	}

	ipbench, err := CoreOSBench.PublicIPv4()
	if err != nil {
		log.Fatalf("could not get bench ip: %v", err)
	}

	fmt.Printf("%s db1\n", ipdb1)
	fmt.Printf("%s bench\n", ipbench)

	cmd := exec.Command("./test.sh", ipdb1, ipbench)

	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Fatalf("bash err %v", err)
	}
	fmt.Printf("test cmd: %q\n", out.String())

	_, err = client.Droplets.Delete(CoreOSDB1.ID)
	if err != nil {
		log.Fatalf("could not destroy db1: %v", err)
	}
	_, err = client.Droplets.Delete(CoreOSBench.ID)
	if err != nil {
		log.Fatalf("could not destroy bench ip: %v", err)
	}

}

// Token create oauth2 token
func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

// CreateCoreOSDroplet will create a droplet with a specific name
func CreateCoreOSDroplet(client *godo.Client, name string, SSHKeyFingerprint string) (*godo.Droplet, error) {
	createRequest := &godo.DropletCreateRequest{
		Name:   name,
		Region: "sgp1",
		Size:   "4gb",
		Image: godo.DropletCreateImage{
			Slug: "coreos-stable",
		},
		SSHKeys: []godo.DropletCreateSSHKey{
			godo.DropletCreateSSHKey{
				Fingerprint: SSHKeyFingerprint,
			},
		},
		Tags: []string{"playwithsql"},
	}

	newDroplet, _, err := client.Droplets.Create(createRequest)
	if err != nil {
		return nil, fmt.Errorf("CreateCoreOSDroplet client.Droplets.Create: %s", err)
	}

	return newDroplet, nil
}

// DeletePreviousCoreOSDroplets delete droplets used for the bench
func DeletePreviousCoreOSDroplets(client *godo.Client) error {
	list, _, err := client.Droplets.ListByTag(
		"playwithsql",
		&godo.ListOptions{
			Page:    1,
			PerPage: 10,
		},
	)
	if err != nil {
		return err
	}

	for _, droplet := range list {
		_, errDel := client.Droplets.Delete(droplet.ID)
		if errDel != nil {
			return errDel
		}
	}

	return nil
}

func getCoreOSStableSlug(client *godo.Client) (string, error) {
	list, _, err := client.Images.List(
		&godo.ListOptions{
			Page:    1,
			PerPage: 1000,
		},
	)
	if err != nil {
		return "", err
	}

	for _, droplet := range list {
		if droplet.Distribution == "CoreOS" && strings.Contains(droplet.Name, "(stable)") {
			return droplet.Slug, nil
		}
	}

	return "", nil
}
