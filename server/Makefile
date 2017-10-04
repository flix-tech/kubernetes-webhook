IMAGE            ?= dcr.mfb.io/devops/webhook-authentication
VERSION          ?= $(shell git describe --tags --always --dirty)
FULLIMAGE		 ?= $(IMAGE):$(VERSION)

default:
	docker build -t $(FULLIMAGE) .
	docker push $(FULLIMAGE)