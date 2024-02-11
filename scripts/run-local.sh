#!/usr/bin/env bash
source settings.sh

kubectl --kubeconfig ${CONFIG} apply -f ../helm/crds/job.yaml
kubectl --kubeconfig ${CONFIG} apply -f ../helm/crds/queuesettings.yaml
LOG_LEVEL=debug go run ../ -kubeconfig=${CONFIG}