install:
	glide install -v

test:
	pwd
	go test
	exit `gofmt -d *.go | wc -l`

build-ci: test
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X main.__VERSION__=${TRAVIS_TAG} -s -w -extldflags '-static'" -o ./dist/kubernetes-webhook-linux
