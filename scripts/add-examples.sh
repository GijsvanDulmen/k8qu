#!/usr/bin/env bash
source settings.sh

cd ..

#kubectl --kubeconfig ${CONFIG} delete queuejobs.k8qu.io --all

#kubectl --kubeconfig ${CONFIG} apply -f ./examples/auto-complete-me/
#
#kubectl --kubeconfig ${CONFIG} apply -f ./examples/other-queues/
#
#kubectl --kubeconfig ${CONFIG} apply -f ./examples/job.yaml
#kubectl --kubeconfig ${CONFIG} apply -f ./examples/job2.yaml
#
#kubectl --kubeconfig ${CONFIG} apply -f ./examples/multi-resources.yaml


#kubectl --kubeconfig ${CONFIG} apply -f ./examples/settings.yaml
#kubectl --kubeconfig ${CONFIG} apply -f ./examples/max-time-inqueue-1.yaml
#sleep 1
#kubectl --kubeconfig ${CONFIG} apply -f ./examples/max-time-inqueue-2.yaml
#kubectl --kubeconfig ${CONFIG} apply -f ./examples/timeout.yaml

#kubectl --kubeconfig ${CONFIG} apply -f ./examples/timeout-on-queue-settings.yaml

#kubectl --kubeconfig ${CONFIG} apply -f ./examples/templates-on-timeout.yaml
#kubectl --kubeconfig ${CONFIG} apply -f ./examples/templates-on-success.yaml
#kubectl --kubeconfig ${CONFIG} apply -f ./examples/templates-on-failure.yaml

#kubectl --kubeconfig ${CONFIG} apply -f ./examples/multi-completions.yaml


#kubectl --kubeconfig ${CONFIG} apply -f ./examples/markcomplete/job.yaml
#sleep 10
#kubectl --kubeconfig ${CONFIG} apply -f ./examples/markcomplete/mark-complete.yaml

#kubectl --kubeconfig ${CONFIG} create -f ./examples/markcomplete-in-parts/job.yaml
#kubectl --kubeconfig ${CONFIG} create -f ./examples/markcomplete-in-parts/job.yaml

#sleep 5
#kubectl --kubeconfig ${CONFIG} apply -f ./examples/markcomplete-in-parts/mark-complete-part-one.yaml
#sleep 5
#kubectl --kubeconfig ${CONFIG} apply -f ./examples/markcomplete-in-parts/mark-complete-part-two.yaml
#kubectl --kubeconfig ${CONFIG} apply -f ./examples/markcomplete-in-parts/mark-complete-part-three.yaml



kubectl --kubeconfig ${CONFIG} create -f ./examples/markcomplete-failure/job.yaml
kubectl --kubeconfig ${CONFIG} create -f ./examples/markcomplete-failure/mark-complete.yaml


# test getting from queue settings
#helm template ./examples/queue-chart/ \
#    --set instance="1" \
#    --set queue="parallel-two" \
#    --set enableTts="no" | kubectl --kubeconfig ${CONFIG} apply -f -

#for i in {1..10}
#do
#  helm template ./examples/queue-chart/ --set instance=${i} --set queue=parallel-two | kubectl --kubeconfig ${CONFIG} apply -f -
#  sleep 2
#done

#for i in {1..100}
#do
#  helm template ./examples/queue-chart/ --set instance=${i} --set queue=q${i} | kubectl --kubeconfig ${CONFIG} apply -f -
#done
