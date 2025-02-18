// Copyright © 2018 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: MIT

package config

import (
	"time"

	"k8s.io/client-go/kubernetes"
)

type Config struct {
	ClustersToWatch     []string
	CIDRToWatch         []string
	K8sClient           map[string]*kubernetes.Clientset
	ClusterToApply      string
	CIDRToApply         string
	NamespaceToWatch    string
	NamespacesToExclude []string
	ReplicatedLabelVal  string
	WatchNamespaces     bool
	WatchEndpoints      bool
	WatchServices       bool
	ResyncPeriod        time.Duration
	CIDR                string
}

const REPLICATED_LABEL_KEY = "replicated"
const KUBERNETES = "kubernetes"
const SVC_ANNOTATION_SYNDICATE_KEY = "vmware.com/syndicate-mode"
const SVC_ANNOTATION_UNION = "union"
const SVC_ANNOTATION_SOURCE = "source"
const SVC_ANNOTATION_RECEIVER = "receiver"
const SVC_ANNOTATION_SINGULAR = "singular"
