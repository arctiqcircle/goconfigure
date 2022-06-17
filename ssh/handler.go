package ssh

import (
	"bytes"
	"github.com/dyntek-services-inc/goconfigure/inventory"
	"golang.org/x/crypto/ssh"
)

type Handler struct {
	host   inventory.Host
	client *ssh.Client
}

func Connect(host inventory.Host) (*Handler, error) {
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
	return &Handler{client: client, host: host}, nil
}

func (h Handler) GetHost() inventory.Host {
	return h.host
}

// Send opens a new session to the SSH server and sends the passed string.
// The standard output from the server is returned.
func (h Handler) Send(command string) (string, error) {
	session, err := h.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	var outBuffer bytes.Buffer
	session.Stdout = &outBuffer
	session.Run(command)
	return outBuffer.String(), nil
}

func (h Handler) Close() error {
	return h.client.Close()
}
