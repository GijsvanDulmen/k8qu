#!/bin/bash
cd ..
docker build --progress=plain -t k8qu .

docker tag k8qu ghcr.io/gijsvandulmen/k8qu:latest
docker tag k8qu ghcr.io/gijsvandulmen/k8qu:1.4

docker push ghcr.io/gijsvandulmen/k8qu:latest
docker push ghcr.io/gijsvandulmen/k8qu:1.4