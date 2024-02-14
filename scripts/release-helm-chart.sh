#!/bin/bash
cd ../helm
helm package .
cp ./*.tgz ../../k8qu-helm/
cd ../../k8qu-helm/
helm repo index .