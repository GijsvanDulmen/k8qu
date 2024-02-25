#!/usr/bin/env bash
source settings.sh

kubectl --kubeconfig ${CONFIG} apply -f ../helm/crds/queuejob.yaml
kubectl --kubeconfig ${CONFIG} apply -f ../helm/crds/queuesettings.yaml
kubectl --kubeconfig ${CONFIG} apply -f ../helm/crds/markqueuejobcomplete.yaml
LOG_LEVEL=debug go run ../ -kubeconfig=${CONFIG}
#go run ../ -kubeconfig=${CONFIG}