package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

// Provider represents a cloud server provider
type Provider interface {
	CreatePlayWithSQLServers(howManyDBServers int) error
	DeleteAllPlayWithSQLServers() error
	FetchAllPlayWithSQLServers() error
	GetBenchServerIP() (string, error)
	GetMasterDBServerIP() (string, error)
}

// LaunchBenches will trigger the shell script launching the bench on the servers
func LaunchBenches(provider Provider) (err error) {

	// err = provider.DeleteAllPlayWithSQLServers()
	// if err != nil {
	// 	return fmt.Errorf("LaunchBenches: %v", err)
	// }

	// err = provider.CreatePlayWithSQLServers(1)
	// if err != nil {
	// 	return fmt.Errorf("LaunchBenches: %v", err)
	// }

	err = provider.FetchAllPlayWithSQLServers()
	if err != nil {
		return fmt.Errorf("LaunchBenches: %v", err)
	}

	benchServerIP, err := provider.GetBenchServerIP()
	if err != nil {
		return fmt.Errorf("LaunchBenches: %v", err)
	}

	masterDBServerIP, err := provider.GetMasterDBServerIP()
	if err != nil {
		return fmt.Errorf("LaunchBenches: %v", err)
	}

	err = runBenches("digiocean_id_rsa", 1000, benchServerIP, masterDBServerIP)
	if err != nil {
		return fmt.Errorf("LaunchBenches: %v", err)
	}

	err = provider.DeleteAllPlayWithSQLServers()
	if err != nil {
		return fmt.Errorf("LaunchBenches: %v", err)
	}

	return nil
}

func runBenches(
	authorizedSSHKey string,
	loops int,
	benchServerIP string,
	masterDBServerIP string,
) error {
	pk, err := ioutil.ReadFile(os.Getenv("HOME") + "/.ssh/" + authorizedSSHKey)
	if err != nil {
		return fmt.Errorf("runBenches: ioutil.ReadFile sshkey%s", err)
	}

	signer, err := ssh.ParsePrivateKey(pk)
	if err != nil {
		return fmt.Errorf("runBenches: ssh.ParsePrivateKey %s", err)
	}

	config := &ssh.ClientConfig{
		User: "core",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	clientBench, err := ssh.Dial("tcp", benchServerIP+":22", config)
	if err != nil {
		return fmt.Errorf("runBenches: ssh.Dial benchServerIP %s", err)
	}

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	sessionBench, err := clientBench.NewSession()
	if err != nil {
		return fmt.Errorf("runBenches: clientBench.NewSession() %s", err)
	}
	defer func() {
		errSess := sessionBench.Close()
		if errSess != nil {
			log.Printf("closing sessionBench: %v", errSess)
		}
	}()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var bBench bytes.Buffer
	var bBenchErr bytes.Buffer
	sessionBench.Stdout = &bBench
	sessionBench.Stderr = &bBenchErr

	clientMasterDB, err := ssh.Dial("tcp", masterDBServerIP+":22", config)
	if err != nil {
		return fmt.Errorf("runBenches: ssh.Dial masterDBServerIP %s", err)
	}

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	sessionMasterDB, err := clientMasterDB.NewSession()
	if err != nil {
		return fmt.Errorf("runBenches: clientMasterDB.NewSession() %s", err)
	}
	defer func() {
		errSess := sessionMasterDB.Close()
		if errSess != nil {
			log.Printf("runBenches: sessionMasterDB.Close() %v", errSess)
		}
	}()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var bMasterDB bytes.Buffer
	var bMasterDBErr bytes.Buffer
	sessionMasterDB.Stdout = &bMasterDB
	sessionMasterDB.Stderr = &bMasterDBErr

	err = sessionMasterDB.Run(`
		PATH='/opt/bin:/usr/bin' && \
		cd /home/core/playwithsql && \
		./infra/databases/docker_local/mssql/container_launch.sh
	`)
	if err != nil {
		return fmt.Errorf("runBenches: sessionMasterDB.Run %s - %s", err, bMasterDBErr.String())
	}
	fmt.Println(bMasterDB.String())

	err = sessionBench.Run(`
		(docker rm -f pws-cmd  || true) && \
		docker run -t --name pws-cmd vincentserpoul/playwithsql-cmd -db=mssql -host=` + masterDBServerIP + ` -loops=$LOOPS && \
		docker rm -f pws-cmd \
	`)
	if err != nil {
		return fmt.Errorf("runBenches: sessionBench.Run %s - %s", err, bBenchErr.String())
	}

	return nil
}
