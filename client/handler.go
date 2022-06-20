package client

import (
	"bytes"
	"github.com/dyntek-services-inc/goconfigure/inventory"
	"golang.org/x/crypto/ssh"
)

type Handler interface {
	Send(string) (string, error)
	Close() error
}

type SSHHandler struct {
	client *ssh.Client
}

func BasicConnect(host inventory.Host) (Handler, error) {
	config := &ssh.ClientConfig{
		User: host.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(host.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", host.Hostname+":22", config)
	if err != nil {
		return nil, err
	}
	return &SSHHandler{client: client}, nil
}

// Send opens a new session to the SSH server and sends the passed string.
// The standard output from the server is returned.
func (h *SSHHandler) Send(command string) (string, error) {
	session, err := h.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	var outBuffer bytes.Buffer
	session.Stdout = &outBuffer
	if err := session.Run(command); err != nil {
		return "", err
	}
	return outBuffer.String(), nil
}

func (h *SSHHandler) Close() error {
	return h.client.Close()
}
