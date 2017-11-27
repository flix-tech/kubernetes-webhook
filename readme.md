# Kubernetes Webhook Authentication via JWT

![Build status](https://travis-ci.org/flix-tech/kubernetes-webhook.svg?branch=master)

This repo is split up into two components:

* jwt-generator

    Generates JWT with a user and her groups. Is only contained in this repo as reference implementation.

* [server](server/readme.md)

    Validates that the JWT is valid and that the client is allowed to authenticate the user and the groups.

Please take a look at the subfolders for more information.

# Further docs
See the [documentation](server/readme.md)