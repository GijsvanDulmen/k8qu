# K8qu - Queue your K8s stuff

[![Go Report Card](https://goreportcard.com/badge/github.com/GijsvanDulmen/k8qu)](https://goreportcard.com/report/github.com/GijsvanDulmen/k8qu)

K8qu is a Kubernetes
[Custom Resource Definition](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
(CRD) controller that allows you to create simple Queues for other resource creations on your cluster. The perfect
companion in your CICD pipelines (like Tekton) and plain Jobs for Kubernetes.

* [Background](#background)
* [Getting Started](#getting-started)
* [Want to contribute?](#want-to-contribute)

## Background

K8qu came to be what it is today because of the lack of a simple solution of queueing in Tekton and other places
within the Kubernetes ecosystem. It tries to keep a simple stupid model for queuing resource creations in your cluster
without the need for a complete queuing solution. Just a simple Kubernetes controller and as little as CRD's as possible.
No databases or other persistency needed. It doesn't want to be a perfect queue, but rather a simple one.

In combination with other frameworks this can become a powerfull addition.

## Cut to the.... CRD!
So how does it look like? Well... like this!

```yaml
apiVersion: k8qu.io/v1alpha1
kind: QueueJob
metadata:
  name: simple-job
spec:
  queue: "abc"
  executionTimeout: "20s"
  
  templates:
    - apiVersion: v1
      kind: Pod
      # ....

```

## Getting Started

To get started with K8qu please see the examples and install the helm chart provided into your cluster.

* [Examples](./examples)
* [Helm chart](./helm/)

The helm chart is also published in `https://raw.githubusercontent.com/GijsvanDulmen/k8qu-helm/main/`.
You can use `helm repo add https://raw.githubusercontent.com/GijsvanDulmen/k8qu-helm/main/` to add it.

## Want to contribute?

Great! Let us know by sending in a Github issue and/or PR. 
