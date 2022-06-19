// GoConfigure is an application for templating configurations meant to be pushed via SSH.
// It can be consumed as a command line tool. Passing arguments of the form `goconfigure_0.1_arm64
// -t TEMPLATE_FILE_NAME -i DATA_FILE_NAME` will render and push the render to the devices
// defined in the inventory file.
package main

import (
	"flag"
	"github.com/dyntek-services-inc/goconfigure/inventory"
	"github.com/dyntek-services-inc/goconfigure/render"
	"log"
	"os"
)

var pwd, pwdErr = os.Getwd()

func main() {
	if pwdErr != nil {
		log.Fatal(pwdErr)
	}
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
			inv, err := inventory.LoadFromYAML(*invFilename)
			if err != nil {
				log.Fatal(err)
			}
			tplString, err := render.FileToString(*tplFilename)
			if err != nil {
				log.Fatal(err)
			}
			if err := inv.Deploy(tplString, pwd); err != nil {
				log.Fatal(err)
			}
		}
	}
}
