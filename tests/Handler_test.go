package tests

import (
	"errors"
	"fmt"
	"github.com/dyntek-services-inc/goconfigure/client"
	"github.com/dyntek-services-inc/goconfigure/inventory"
	"github.com/dyntek-services-inc/goconfigure/render"
	"log"
	"strings"
	"testing"
)

func TestConnectHandler(t *testing.T) {
	t.Log("Testing connections from YAML inventory")
	inv, err := inventory.LoadFromYAML("secrets/hosts.yml")
	if err != nil {
		panic(err)
	}
	// BasicConnect to Hosts
	var handlers []client.Handler
	for _, host := range inv.Hosts {
		t.Logf("Connecting to %s", host.Hostname)
		h, err := client.BasicConnect(host)
		handlers = append(handlers, h)
		if err != nil {
			panic(err)
		}
	}
	// Cleanup
	for _, h := range handlers {
		if err := h.Close(); err != nil {
			log.Fatal("handlers failed to close properly", err)
		}
	}
	t.Log("Testing connections from CSV inventory")
	inv, err = inventory.LoadFromCSV("secrets/hosts.csv")
	if err != nil {
		panic(err)
	}
	// BasicConnect to Hosts
	handlers = []client.Handler{}
	for _, host := range inv.Hosts {
		t.Logf("Connecting to %s", host.Hostname)
		h, err := client.BasicConnect(host)
		handlers = append(handlers, h)
		if err != nil {
			panic(err)
		}
	}
	// Cleanup
	for _, h := range handlers {
		if err := h.Close(); err != nil {
			log.Fatal("handlers failed to close properly", err)
		}
	}
}

func TestEchoHandler(t *testing.T) {
	t.Log("Testing connections from YAML inventory")
	inv, err := inventory.LoadFromYAML("secrets/hosts.yml")
	if err != nil {
		panic(err)
	}
	// BasicConnect to Hosts
	var handlers []client.Handler
	for _, host := range inv.Hosts {
		t.Logf("Connecting to %s", host.Hostname)
		h, err := client.BasicConnect(host)
		handlers = append(handlers, h)
		if err != nil {
			panic(err)
		}
	}
	// Send Command to Host
	for _, h := range handlers {
		response, err := h.Send("echo \"hello world!\"")
		if err != nil {
			panic(err)
		}
		response = strings.TrimSpace(response)
		if response != "hello world!" {
			panic(errors.New(fmt.Sprintf("response %s not equal to %s", response, "hello world!")))
		} else {
			t.Logf("response %s succesfully matches %s", response, "hello world!")
		}
	}
	// Cleanup
	for _, h := range handlers {
		if err := h.Close(); err != nil {
			log.Fatal("handlers failed to close properly", err)
		}
	}
	t.Log("Testing connections from CSV inventory")
	inv, err = inventory.LoadFromCSV("secrets/hosts.csv")
	if err != nil {
		panic(err)
	}
	// BasicConnect to Hosts
	handlers = []client.Handler{}
	for _, host := range inv.Hosts {
		t.Logf("Connecting to %s", host.Hostname)
		h, err := client.BasicConnect(host)
		handlers = append(handlers, h)
		if err != nil {
			panic(err)
		}
	}
	// Send Command to Host
	for _, h := range handlers {
		response, err := h.Send("echo \"hello world!\"")
		if err != nil {
			panic(err)
		}
		response = strings.TrimSpace(response)
		if response != "hello world!" {
			panic(errors.New(fmt.Sprintf("response %s not equal to %s", response, "hello world!")))
		} else {
			t.Logf("response %s succesfully matches %s", response, "hello world!")
		}
	}
	// Cleanup
	for _, h := range handlers {
		if err := h.Close(); err != nil {
			log.Fatal("handlers failed to close properly", err)
		}
	}
}

func TestTemplateHandler(t *testing.T) {
	t.Log("Testing connections from YAML inventory")
	inv, err := inventory.LoadFromYAML("secrets/hosts.yml")
	tplString, err := render.FileToString("secrets/example.txt")
	if err != nil {
		panic(err)
	}
	// BasicConnect to Hosts and Render Template
	var handlers []client.Handler
	var tpls [][]string
	for _, host := range inv.Hosts {
		t.Logf("Connecting to %s", host.Hostname)
		h, err := client.BasicConnect(host)
		if err != nil {
			panic(err)
		}
		handlers = append(handlers, h)
		tpls = append(tpls, render.RenderCommands(host.Data, tplString))
	}
	// Send Commands to Host
	for i, h := range handlers {
		for _, command := range tpls[i] {
			response, err := h.Send(command)
			if err != nil {
				panic(err)
			}
			response = strings.TrimSpace(response)
		}
	}
	// Cleanup
	for _, h := range handlers {
		if err := h.Close(); err != nil {
			log.Fatal("handlers failed to close properly", err)
		}
	}
	t.Log("Testing connections from CSV inventory")
	inv, err = inventory.LoadFromCSV("secrets/hosts.csv")
	if err != nil {
		panic(err)
	}
	// BasicConnect to Hosts and Render Template
	handlers = []client.Handler{}
	tpls = [][]string{{}}
	for _, host := range inv.Hosts {
		t.Logf("Connecting to %s", host.Hostname)
		h, err := client.BasicConnect(host)
		if err != nil {
			panic(err)
		}
		handlers = append(handlers, h)
		tpls = append(tpls, render.RenderCommands(host.Data, tplString))
	}
	// Send Commands to Host
	for i, h := range handlers {
		for _, command := range tpls[i] {
			response, err := h.Send(command)
			if err != nil {
				panic(err)
			}
			response = strings.TrimSpace(response)
		}
	}
	// Cleanup
	for _, h := range handlers {
		if err := h.Close(); err != nil {
			log.Fatal("handlers failed to close properly", err)
		}
	}
}
