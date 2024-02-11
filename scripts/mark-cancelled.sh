#!/usr/bin/env bash
source settings.sh

kubectl --kubeconfig ${CONFIG} patch jobs.k8qu.io $1 --type merge --patch-file ./patch-file-failed.json
