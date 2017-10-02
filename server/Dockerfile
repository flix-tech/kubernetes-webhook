FROM golang:1.9-stretch
WORKDIR /go/src/app
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-s' -o main
FROM debian:stretch
COPY --from=0 /go/src/app/main /main
ENTRYPOINT /main
CMD ["-config", "/config.yml"]