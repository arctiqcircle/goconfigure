// GoConfigure is an application for templating configurations meant to be pushed via SSH.
// It can be consumed as a command line tool. Passing arguments of the form `goconfigure
// -t TEMPLATE_FILE_NAME -i DATA_FILE_NAME` will render and push the render to the devices
// defined in the inventory file.
package main

import (
	"flag"
	"github.com/dyntek-services-inc/goconfigure/inventory"
	"github.com/dyntek-services-inc/goconfigure/render"
	"github.com/dyntek-services-inc/goconfigure/ssh"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var pwd, pwdErr = os.Getwd()

func deploy(inventory inventory.Inventory, tplString string) error {
	rc := make(map[string]chan []string, len(inventory)) // The response channels.
	for name, host := range inventory {
		log.Printf("starting deployment for %s", host.Hostname)
		rc[name] = make(chan []string)
		handler, err := ssh.Connect(host)
		log.Printf("finished connecting too %s", host.Hostname)
		if err != nil {
			return err
		}
		go func(ro chan []string, hdlr *ssh.Handler) {
			rtplc := render.RenderCommands(hdlr.GetHost().Data, tplString)
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
			log.Printf("finished deployment of %s", hdlr.GetHost().Hostname)
			ro <- cco
		}(rc[name], handler)
	}
	for name, ro := range rc {
		if pwdErr != nil {
			return pwdErr
		}
		rro := <-ro
		tr := strings.Join(rro, "\n")
		of := filepath.Join(pwd, name)
		log.Printf("writing output to %s.txt", of)
		if err := os.WriteFile(of+".txt", []byte(tr), 0666); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	invFilename := flag.String("i", "", "inventory filename")
	tplFilename := flag.String("t", "", "template filename")
	flag.Parse()
	if len(*invFilename) == 0 && len(*tplFilename) == 0 {
		// No inventory or template flags were passed, start manual mode
		// TODO: implement manual mode
	} else {
		if len(*invFilename) == 0 || len(*tplFilename) == 0 {
			// One of the flags was passed but not the other
			log.Fatal("one flag was passed, but not both")
		} else {
			// Both flags were passed, begin deployment
			inv, err := inventory.LoadInventory(*invFilename)
			if err != nil {
				log.Fatal(err)
			}
			tplString, err := render.FileToString(*tplFilename)
			if err != nil {
				log.Fatal(err)
			}
			if err := deploy(inv, tplString); err != nil {
				log.Fatal(err)
			}
		}
	}
}
