package client

import (
	"github.com/dyntek-services-inc/goconfigure/inventory"
	"github.com/dyntek-services-inc/goconfigure/render"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Deploy(inv *inventory.Inventory, tplString, logDir string) error {
	rc := make([]chan []string, len(inv.Hosts)) // The response channels.
	for ih, host := range inv.Hosts {
		log.Printf("starting deployment for %s", host.Hostname)
		rc[ih] = make(chan []string)
		handler, err := BasicConnect(host)
		log.Printf("finished connecting to %s", host.Hostname)
		if err != nil {
			return err
		}
		go func(ro chan []string, h Handler, host inventory.Host) {
			rtplc := render.RenderCommands(host.Data, tplString)
			cc := make([]chan string, len(rtplc))
			for i, c := range rtplc {
				cc[i] = make(chan string)
				go func(co chan string, ci string) {
					r, err := h.Send(ci)
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
