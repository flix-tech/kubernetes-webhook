# Key generation

    ssh-keygen -t rsa -b 2048 -f jwtRS256.key
    # Don't add passphrase
    openssl rsa -in jwtRS256.key -pubout -outform PEM -out jwtRS256.key.pub

# Usage

    cat jwtRS256.key | ./jwt-generator --user gitlab --group gitlab-group1 --group gitlab-group2