module github.com/wxio/lcmd

go 1.13

replace docker.io/go-docker => github.com/millergarym/go-docker v1.0.1

require (
	docker.io/go-docker v1.0.0
	github.com/docker/go-connections v0.4.0
	github.com/golang/glog v1.0.0
	github.com/golang/protobuf v1.5.2
	github.com/jpillora/opts v1.1.2
	google.golang.org/grpc v1.53.0
)
