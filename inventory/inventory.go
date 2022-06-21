package inventory

import (
	"encoding/csv"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
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

func LoadFromCSV(path string, panicOnEmpty bool) (*Inventory, error) {
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
	for {
		host := Host{Data: make(map[string]interface{})}
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
		if panicOnEmpty {
			if len(host.Password) == 0 {
				return nil, errors.New(fmt.Sprintf("inventory file %s does not contain password field", path))
			}
		}
		inv.Hosts = append(inv.Hosts, host)
	}
	return &inv, nil
}
