package main

import (
	"flag"
	"fmt"

	"github.com/ICKelin/cframe/pkg/edagemanager"
	log "github.com/ICKelin/cframe/pkg/logs"
)

func main() {
	flgConf := flag.String("c", "", "config file path")
	flag.Parse()

	conf, err := ParseConfig(*flgConf)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Init(conf.Log.Path, conf.Log.Level, conf.Log.Days)
	log.Info("%v", conf)

	// create etcd storage
	store := edagemanager.NewEtcdStorage(conf.Etcd)

	// create edage manager
	edageManager := edagemanager.New(store)

	// create edage host manager
	edagemanager.NewEdageHostManager(store)

	r := NewRegistryServer(conf.ListenAddr)

	// watch for edage delete/put
	// tell registry edage change
	go edageManager.Watch(
		func(edg *edagemanager.Edage) {
			r.DelEdage(edg)
		},
		func(edg *edagemanager.Edage) {
			r.ModifyEdage(edg)
		})

	r.ListenAndServe()
}
