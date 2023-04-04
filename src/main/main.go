// Copyright © 2018 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

package main

import (
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	c "github.com/joseret/k8s-endpoints-sync-controller/src/config"
	cc "github.com/joseret/k8s-endpoints-sync-controller/src/controller"
	"github.com/joseret/k8s-endpoints-sync-controller/src/handlers"
	log "github.com/joseret/k8s-endpoints-sync-controller/src/log"
)

func main() {

	log.Initialize()
	log.Infof("Starting clusterdiscovery controller")
	config, err := loadConfig()
	if err != nil {
		return
	}

	handler := &handlers.ClusterDiscoveryHandler{}
	if handlerErr := handler.Init(config); handlerErr != nil {
		log.Errorf("failed to initialize handler %v", handlerErr)
		return
	}
	for _, cluster := range config.ClustersToWatch {

		go cc.StartController(cluster, handler, config)

	}

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm
}

func loadConfig() (*c.Config, error) {

	conf := &c.Config{}

	if n, nexists := os.LookupEnv("NSTOWATCH"); nexists {
		conf.NamespaceToWatch = n
	} else {
		conf.NamespaceToWatch = ""
	}

	if e, eexists := os.LookupEnv("EXCLUDE"); eexists {
		log.Infof("Namespaces to exclude %s", e)
		conf.NamespacesToExclude = strings.Split(e, ",")
	}

	if cidr, eexists := os.LookupEnv("CIDR"); eexists {
		log.Infof("CIDR of local cluster to exclude %s", eexists)
		conf.CIDR = cidr
	}

	searchDir := "/etc/kubeconfigs"

	files, err := ioutil.ReadDir(searchDir)
	if err != nil {
		log.Errorf("Error reading dir %v", err)
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() && !strings.Contains(file.Name(), "data") {
			log.Infof("Kubeconfig of cluster to watch %s", file.Name())
			conf.ClustersToWatch = append(conf.ClustersToWatch, searchDir+"/"+file.Name())
		}
	}

	conf.ReplicatedLabelVal = "true"

	conf.WatchNamespaces = true
	conf.WatchEndpoints = true
	conf.WatchServices = true
	conf.ResyncPeriod = 30 * time.Second

	return conf, nil
}
