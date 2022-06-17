package inventory

import (
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type Inventory map[string]Host

type Host struct {
	Hostname string                 `yaml:"hostname"`
	Username string                 `yaml:"username"`
	Password string                 `yaml:"password"`
	Data     map[string]interface{} `yaml:"data"`
}

func LoadInventory(filename string) (Inventory, error) {
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
	return inv, nil
}
