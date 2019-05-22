module github.com/flix-tech/kubernetes-webhook/server

go 1.12

replace github.com/flix-tech/kubernetes-webhook v0.1.1 => ../

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/flix-tech/kubernetes-webhook v0.1.1
	github.com/gobwas/glob v0.2.3
	github.com/prometheus/common v0.4.1
	gopkg.in/yaml.v2 v2.2.2
	k8s.io/api v0.0.0-20190515023547-db5a9d1c40eb // indirect
	k8s.io/apimachinery v0.0.0-20190515023456-b74e4c97951f
	k8s.io/kubernetes v1.14.2
)
