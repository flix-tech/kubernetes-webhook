

kubernetes-webhook-server: $(shell find src -f)
	cd src/kubernetes-webhook-server && docker build -t nielsole/kubernetes-webhook-server .

push: kubernetes-webhook-server
	docker push nielsole/kubernetes-webhook-server