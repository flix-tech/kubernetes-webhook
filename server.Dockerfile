FROM golang:1.12-alpine3.9
RUN apk add --no-cache git
WORKDIR /go/src/app
ADD . /customapp/
ENV GOPATH /go/src/app
WORKDIR /customapp/server
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-s' -o main
FROM alpine:3.9
COPY --from=0 /customapp/server/main /main
ENTRYPOINT /main
CMD ["-config", "/config.yml"]