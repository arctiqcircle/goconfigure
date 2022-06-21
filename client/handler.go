package client

import (
	"bytes"
	"github.com/dyntek-services-inc/goconfigure/inventory"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
)

type Handler interface {
	Send(string) (string, error)
	io.Closer
}

type SSHHandler struct {
	client *ssh.Client
}

type Authentication func(host inventory.Host) (Handler, error)

//func BasicConnect(host inventory.Host) (Handler, error) {
//	config := &ssh.ClientConfig{
//		User: host.Username,
//		Auth: []ssh.AuthMethod{
//			ssh.Password(host.Password),
//		},
//		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
//	}
//	client, err := ssh.Dial("tcp", host.Hostname+":22", config)
//	if err != nil {
//		return nil, err
//	}
//	return &SSHHandler{client: client}, nil
//}

func BasicConnect() (Authentication, error) {
	return func(host inventory.Host) (Handler, error) {
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
	}, nil
}

func KeyConnect(keyfile string) (Authentication, error) {
	pemBytes, err := os.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(pemBytes)
	return func(host inventory.Host) (Handler, error) {
		config := &ssh.ClientConfig{
			User: host.Username,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
				//ssh.Password(host.Password),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}
		c, err := ssh.Dial("tcp", host.Hostname+":22", config)
		if err != nil {
			return nil, err
		}
		h := &SSHHandler{client: c}
		return h, nil
	}, nil
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
