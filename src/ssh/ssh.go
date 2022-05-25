package ssh

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
	"path"
)

type SshClient struct {
	client *ssh.Client
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

	client := &SshClient{
		client: conn,
	}

	return client, nil
}

func (s SshClient) Close() {
	s.client.Close()
}

// Opens a new SSH connection and runs the specified command
// Returns the combined output of stdout and stderr
func (s *SshClient) RunCommand(cmd string) (string, error) {
	// open connection
	// run command and capture stdout/stderr
	session, err := s.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("Create session for failed %v", err)
	}
	output, err := session.CombinedOutput(cmd)
	return string(output), err
}

func (s *SshClient) ScpCopy(localFilePath, remoteDir string) error {

	sftpClient, err := sftp.NewClient(s.client)
	if err != nil { //创建客户端
		fmt.Println("创建客户端失败", err)
		return err
	}
	defer sftpClient.Close()

	shellContent, err := ioutil.ReadFile(localFilePath)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
		return err
	}

	var remoteFileName = path.Base(localFilePath)

	dstFile, err := sftpClient.Create(path.Join(remoteDir, remoteFileName))
	if err != nil {
		log.Println("scpCopy:", err, "remote path:", path.Join(remoteDir, remoteFileName))
		return err
	}
	defer dstFile.Close()
	dstFile.Write(shellContent)
	return nil
}
