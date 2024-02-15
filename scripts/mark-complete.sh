#!/usr/bin/env bash
source settings.sh

kubectl --kubeconfig ${CONFIG} patch queuejobs.k8qu.io $1 --type merge --patch-file ./patch-file.json
