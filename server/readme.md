# Webhook authentication Server

This server responds to Kubernetes Authentication challenges and responds with the groups of the user after checking whether the signer is allowed to log the user in.

## Installation

### Binary
Compiled binaries are available as [releases](https://github.com/flix-tech/kubernetes-webhook/releases).

### From source

Go modules are used. This package cannot be in the $GOPATH, but must be in a separate directory.
Use a recent Go version (e.g. >=1.12)

    # In this directory:
    make build-ci
    # or in the root of this repository:
    docker build -t kubernetes-webhook -f server.Dockerfile .

## Configuration

### Webhook Server

Run the binary with `-h` to get all available CLI options. Most notable pass the path to your `config`. which contains the keys and globs that the server accepts.
See [config.yml](config.yml) for an example.

    ./kubernetes-webhook -config /path/to/config

### Kubernetes

The API-server needs to be configured to use this server for webhooks.

    apiserver.Authentication.WebHook.ConfigFile="/path/to/WebHookConfigFile"

The file has the same serializer as the Kubernetes client config.

Example for the WebHookConfigFile:

    clusters:
    - name: token-authenticator
      cluster:
        server: https://address-of-jwt-authentication-server/authenticate
        certificate-authority: /path/to/trusted/authority.crt
    
    users:
    - name: webhook-authentication
    
    contexts:
    - name: webhook
      context:
        cluster: token-authenticator
        user: webhook-authentication
    
    current-context: webhook

For further info see the [K8s documentation](https://kubernetes.io/docs/admin/authentication/#webhook-token-authentication)

## Usage

To run this container on the APIserver you could run:

    docker run -p 127.0.0.1:3443:3000 --name kubernetes-webhook   \
    -v /etc/kube-webhook-auth/config.yml:/config.yml:ro   \
    -v /etc/kube-webhook-auth/trusted-signer.key:/trusted-signer.key:ro   \
    kubernetes-webhook