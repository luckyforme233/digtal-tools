package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
)

type SshClient struct {
	session *ssh.Session
	client  *ssh.Client
}

func NewSshClient(user string, host string, port int, privateKeyPath string, privateKeyPassword string) (*SshClient, error) {
	fmt.Println("privateKeyPath:", privateKeyPath)

	key, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	// build SSH client config
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			// use OpenSSH's known_hosts file if you care about host validation
			return nil
		},
	}
	server := fmt.Sprintf("%v:%v", host, port)
	conn, err := ssh.Dial("tcp", server, config)
	if err != nil {
		return nil, fmt.Errorf("Dial to %v failed %v", server, err)
	}

	// open session
	session, err := conn.NewSession()
	if err != nil {
		return nil, fmt.Errorf("Create session for %v failed %v", server, err)
	}

	client := &SshClient{
		session: session,
		client:  conn,
	}

	return client, nil
}

func (s SshClient) Close() {
	s.session.Close()
	s.client.Close()
}

// Opens a new SSH connection and runs the specified command
// Returns the combined output of stdout and stderr
func (s *SshClient) RunCommand(cmd string) (string, error) {
	// open connection
	// run command and capture stdout/stderr
	output, err := s.session.Output(cmd)
	fmt.Println("runCommand : ", string(output), err)
	return fmt.Sprintf("%s", output), err
}
