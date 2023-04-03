# Copyright © 2018 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: MIT

FROM ubuntu:focal-20230308
ADD dist/k8s-endpoints-sync-controller /k8s-endpoints-sync-controller
RUN chmod +x /k8s-endpoints-sync-controller
CMD "/k8s-endpoints-sync-controller"
