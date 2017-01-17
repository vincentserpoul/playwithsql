package digitalocean

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

// DOProvider will implement the provider interface
type DOProvider struct {
	DOToken           string
	SSHKeyFingerprint string
	Client            *godo.Client
	BenchServer       *godo.Droplet
	DBServers         []*godo.Droplet
}

// New creates DO provider with the required params
func New(doToken string, sshKeyFingerprint string) *DOProvider {
	doProv := &DOProvider{
		DOToken:           doToken,
		SSHKeyFingerprint: sshKeyFingerprint,
	}
	doProv.Client = getDOClient(doToken)

	return doProv
}

// TokenSource contains the token, used for oauth2
type TokenSource struct {
	AccessToken string
}

// CreatePlayWithSQLServers creates all necessary coreos servers for the bench
func (d *DOProvider) CreatePlayWithSQLServers(howManyDBServers int) (err error) {

	// Create bench server
	d.BenchServer, err = createPlayWithSQLServer(d.Client, "CoreOSBench", d.SSHKeyFingerprint)
	if err != nil {
		return fmt.Errorf("CreatePlayWithSQLServers could not create bench droplet: %v", err)
	}

	for i := 0; i < howManyDBServers; i++ {
		// Create DB server
		CoreOSDB, errDB := createPlayWithSQLServer(
			d.Client,
			"CoreOSDB"+strconv.Itoa(i),
			d.SSHKeyFingerprint,
		)
		if errDB != nil {
			log.Fatalf("could not create db droplet: %v", errDB)
		}
		d.DBServers = append(d.DBServers, CoreOSDB)
	}

	AllCreated := false
	// Wait until they are all created
	for !AllCreated {
		time.Sleep(time.Second * 10)
		log.Println("waiting for droplet to be created")
		isCreated, errStatus := isDropletActive(d.Client, d.BenchServer.ID)
		if errStatus != nil {
			return fmt.Errorf("CreatePlayWithSQLServers(%d): %v", howManyDBServers, errStatus)
		}
		AllCreated = isCreated

		for _, dbServer := range d.DBServers {
			isCreated, errStatus := isDropletActive(d.Client, dbServer.ID)
			if errStatus != nil {
				return fmt.Errorf("CreatePlayWithSQLServers(%d): %v", howManyDBServers, errStatus)
			}
			AllCreated = AllCreated && isCreated
		}
	}

	return nil
}

// GetBenchServerIP returns the IP of the bench server
func (d *DOProvider) GetBenchServerIP() (string, error) {
	return d.BenchServer.PublicIPv4()
}

// GetMasterDBServerIP returns the IP of the master db server
func (d *DOProvider) GetMasterDBServerIP() (string, error) {
	return d.DBServers[0].PublicIPv4()
}

// DeleteAllPlayWithSQLServers destroys all instances of play with SQL
func (d *DOProvider) DeleteAllPlayWithSQLServers() error {
	return deleteAllPlayWithSQLServers(d.Client)
}

// isDropletActive check if a droplet is active or not
func isDropletActive(client *godo.Client, id int) (bool, error) {
	droplet, _, err := client.Droplets.Get(id)
	if err != nil {
		return false, fmt.Errorf("isDropletActive(%d): %v", id, err)
	}
	return droplet.Status == "active", nil
}

// Token create oauth2 token
func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

// createPlayWithSQLServer create a droplet with a specific name
func createPlayWithSQLServer(
	client *godo.Client,
	name string,
	SSHKeyFingerprint string,
) (*godo.Droplet, error) {
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
		return nil, fmt.Errorf("CreateCoreOSDclient.Droplets.Create: %s", err)
	}

	return newDroplet, nil
}

// deleteAllPlayWithSQLServers delete droplets used for the bench
func deleteAllPlayWithSQLServers(client *godo.Client) error {

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

// getDOClient gets a client for digital ocean
func getDOClient(doToken string) *godo.Client {
	tokenSource := &TokenSource{AccessToken: doToken}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	return godo.NewClient(oauthClient)
}
