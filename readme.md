# Kubernetes Webhook Authentication via JWT

This repo is split up into two components:

* jwt-generator

    Generates JWT with a user and her groups

* server

    Validates that the JWT is valid and that the client is allowed to authenticate the user and the groups.

Please take a look at the subfolders for more information.

## Initialization

This repo expects to be at `$GOPATH/src/flix-tech/kubernetes-webhook`.
Glide is used for dependency management. 

Run `glide install -v` in both subfolders.
