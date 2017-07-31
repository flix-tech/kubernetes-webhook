# JWT Server

This server responds to Kubernetes Authentication challenges and checks whether the signer is allowed to log the user in.

## Build

    glide install -v
    go install
    
## Test

    go test

## Configuration

### Webhook Server

Run the binary with `-h` to get all available CLI options. Most notable pass the path to your `config`. which contains the keys and globs that the server accepts.
See `config.yml` for an example

### Kubernetes

The API-server needs to be configured to use this server for webhooks.

    apiserver.Authentication.WebHook.ConfigFile="/path/to/WebHookConfigFile"

The file has the same serializer as the Kubernetes client config.
For further info see the [K8s documentation](https://kubernetes.io/docs/admin/authentication/#webhook-token-authentication)