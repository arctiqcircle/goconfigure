package inventory

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/dyntek-services-inc/goconfigure/render"
	"github.com/dyntek-services-inc/goconfigure/ssh"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Inventory struct {
	DefaultUsername string `yaml:"default_username" json:"default_username"`
	DefaultPassword string `yaml:"default_password" json:"default_password"`
	Hosts           []Host `yaml:"hosts" json:"hosts"`
}

type Host struct {
	Hostname string                 `yaml:"hostname" json:"hostname"`
	Username string                 `yaml:"username" json:"username"`
	Password string                 `yaml:"password" json:"password"`
	Data     map[string]interface{} `yaml:"data" json:"data"`
}

func LoadFromYAML(filename string) (*Inventory, error) {
	inv := Inventory{}
	yFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	yContent, err := io.ReadAll(yFile)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(yContent, &inv); err != nil {
		return nil, err
	}
	return &inv, nil
}

func LoadFromCSV(path string) (*Inventory, error) {
	inv := Inventory{}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	records := make([]map[string]interface{}, 1)
	// here we read the header line
	records[0] = make(map[string]interface{})
	reader := csv.NewReader(file)
	hrec, err := reader.Read()
	if err != nil {
		return nil, err
	}
	headers := make([]string, len(hrec))
	for i, field := range hrec {
		headers[i] = field
	}
	// TODO: convert to a string builder to handle larger documents
	for {
		host := Host{}
		rec, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		for i, field := range rec {
			switch headers[i] {
			case "hostname":
				host.Hostname = field
			case "username":
				host.Username = field
			case "password":
				host.Password = field
			default:
				host.Data[headers[i]] = field
			}
		}
		if len(host.Hostname) == 0 {
			return nil, errors.New(fmt.Sprintf("inventory file %s does not contain hostname field", path))
		}
		if len(host.Username) == 0 {
			return nil, errors.New(fmt.Sprintf("inventory file %s does not contain username field", path))
		}
		if len(host.Password) == 0 {
			return nil, errors.New(fmt.Sprintf("inventory file %s does not contain password field", path))
		}
	}
	return &inv, nil
}

func (inv *Inventory) Deploy(tplString, logDir string) error {
	rc := make([]chan []string, len(inv.Hosts)) // The response channels.
	for ih, host := range inv.Hosts {
		log.Printf("starting deployment for %s", host.Hostname)
		rc[ih] = make(chan []string)
		handler, err := ssh.Connect(host.Hostname, host.Username, host.Password)
		log.Printf("finished connecting too %s", host.Hostname)
		if err != nil {
			return err
		}
		go func(ro chan []string, hdlr *ssh.Handler, host Host) {
			rtplc := render.RenderCommands(host.Data, tplString)
			cc := make([]chan string, len(rtplc))
			for i, c := range rtplc {
				cc[i] = make(chan string)
				go func(co chan string, ci string) {
					r, err := hdlr.Send(ci)
					if err != nil {
						log.Fatal(err)
					}
					co <- r
				}(cc[i], c)
			}
			cco := make([]string, len(rtplc))
			for i, co := range cc {
				cco[i] = <-co
			}
			log.Printf("finished deployment of %s", host.Hostname)
			ro <- cco
		}(rc[ih], handler, host)
	}
	for ri, ro := range rc {
		rro := <-ro
		tr := strings.Join(rro, "\n")
		of := filepath.Join(logDir, inv.Hosts[ri].Hostname)
		log.Printf("writing output to %s.txt", of)
		if err := os.WriteFile(of+".txt", []byte(tr), 0666); err != nil {
			return err
		}
	}
	return nil
}

func (h *Host) GetLogin() (string, string, string) {
	return h.Hostname, h.Username, h.Password
}
