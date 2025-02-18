// Copyright © 2018 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

package controller

import (
	c "github.com/joseret/k8s-endpoints-sync-controller/src/config"
	"github.com/joseret/k8s-endpoints-sync-controller/src/handlers"
	log "github.com/joseret/k8s-endpoints-sync-controller/src/log"
	"k8s.io/apimachinery/pkg/util/wait"
	informercorev1 "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func StartController(kubeconfigPath string, eventHandler handlers.Handler, config *c.Config) error {
	kubeClient, err := getkubeclient(kubeconfigPath)
	if err != nil {
		return err
	}
	if config.WatchNamespaces {
		watchNamespaces(kubeClient, eventHandler, config)
	}
	if config.WatchEndpoints {
		watchEndpoints(kubeClient, eventHandler, config)
	}
	if config.WatchServices {
		watchServices(kubeClient, eventHandler, config)
	}
	//TODO - better error handling
	config.K8sClient[kubeconfigPath], _ = getkubeclient(kubeconfigPath)
	return nil
}

func getkubeclient(kubeconfigPath string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	log.Infof("building kubeclient")
	if err != nil {
		log.Errorf("Error with kubeconfig %s", err)
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func watchNamespaces(client *kubernetes.Clientset, eventHandler handlers.Handler, config *c.Config) cache.Store {

	indexers := cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}
	informer := informercorev1.NewNamespaceInformer(client, 0, indexers)

	informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    eventHandler.ObjectCreated,
			UpdateFunc: eventHandler.ObjectUpdated,
			DeleteFunc: eventHandler.ObjectDeleted,
		},
	)
	go informer.Run(wait.NeverStop)
	log.Infof("Waiting for namespaces to be synced")
	cache.WaitForCacheSync(wait.NeverStop, informer.HasSynced)
	log.Infof("synced namespaces")

	return nil
}

func watchEndpoints(client *kubernetes.Clientset, eventHandler handlers.Handler, config *c.Config) cache.Store {

	log.Infof("watchEndpoints-config.ResyncPeriod[%v]", config.ResyncPeriod)
	indexers := cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}
	informer := informercorev1.NewFilteredEndpointsInformer(client, "test-1", config.ResyncPeriod, indexers, nil)

	informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    eventHandler.ObjectCreated,
			UpdateFunc: eventHandler.ObjectUpdated,
			DeleteFunc: eventHandler.ObjectDeleted,
		},
	)
	go informer.Run(wait.NeverStop)
	return nil
}

func watchServices(client *kubernetes.Clientset, eventHandler handlers.Handler, config *c.Config) cache.Store {
	log.Infof("watchServices-config.ResyncPeriod[%v]", config.ResyncPeriod)
	indexers := cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}
	informer := informercorev1.NewFilteredServiceInformer(client, "test-1", config.ResyncPeriod, indexers, nil)

	informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    eventHandler.ObjectCreated,
			UpdateFunc: eventHandler.ObjectUpdated,
			DeleteFunc: eventHandler.ObjectDeleted,
		},
	)
	go informer.Run(wait.NeverStop)
	return nil
}
