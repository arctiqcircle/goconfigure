package ssh

import (
	"bytes"
	"golang.org/x/crypto/ssh"
)

type Handler struct {
	client *ssh.Client
}

func Connect(hostname, username, password string) (*Handler, error) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", hostname+":22", config)
	if err != nil {
		return nil, err
	}
	return &Handler{client: client}, nil
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
