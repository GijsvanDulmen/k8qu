#!/usr/bin/env bash
source settings.sh

cd ..
KUBECONFIG=${CONFIG} helm -n k8qu install --set LOG_LEVEL=DEBUG --create-namespace k8qu helm/