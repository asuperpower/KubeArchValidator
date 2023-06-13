# Kubernetes Architecture Validator
Runs an admission controller to validate that the architecture of the image being pulled matches the architecture of the cluster. Aims to prevent the pulling of ARM images (built on M1 Macs) into amd64 clusters, as this can be a pain to remove without elevated permissions.
Only works for images using the V2 docker endpoint that allows the reading of metadata to look up the image architecture


# KubeArchValidator

KubeArchValidator is a Kubernetes admission controller that validates the architecture of Docker images against the node architecture of your Kubernetes cluster.

## Overview

This application is designed to prevent the deployment of Docker images that do not match the architecture of your Kubernetes nodes. This is useful to prevent issues when deploying images built for different architectures, such as images built on Apple M1 chips (arm64) being deployed to a cluster of x86-64 nodes.

## Gotchas
Your images all _need_ to be on a Docker Registry HTTP API V2 to read the metadata and get the image architecture. If the V2 endpoint is not available, the image will get rejected. If you're using something like an Azure Container Registry, you should be fine.

## Deployment

KubeArchValidator should be deployed as a service within your Kubernetes cluster. You will also need to register it as a `ValidatingWebhookConfiguration` in your Kubernetes cluster to let the API server know to send admission requests to it. Detailed instructions for deployment will be provided at a later stage.

## Testing

For testing purposes, you can run KubeArchValidator locally using a tool like Minikube. We will provide instructions on how to set this up in the future.

Please note that KubeArchValidator is currently in development and not ready for production use.

## Contributing

Contributions to KubeArchValidator are welcome. Please make sure to test your changes thoroughly before submitting a pull request.
