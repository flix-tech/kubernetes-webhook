# Kubernetes Webhook Authentication via JWT

![Build status](https://travis-ci.org/flix-tech/kubernetes-webhook.svg?branch=master)

The Kubernetes webhook authentication via JWT allows you to authenticate to Kubernetes via JWT.
This allows you to have a trusted component that gives out temporary tokens with limited permissions in the cluster. It is useful for CI pipelines.

This repo is split up into two components:


* [server](server/readme.md)

    Authentication plugin for Kubernetes. Validates that the JWT is valid provides Kubernetes with the allowed groups.
* jwt-generator

    Reference implementation of JWT generator with a user and her groups.

[Getting started](server/readme.md)