module github.com/wxio/lcmd

go 1.13

replace docker.io/go-docker => github.com/millergarym/go-docker v1.0.1

require (
	docker.io/go-docker v1.0.0
	github.com/docker/go-connections v0.4.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.3.2
	github.com/jpillora/opts v1.1.2
	github.com/kr/pretty v0.1.0 // indirect
	google.golang.org/grpc v1.24.0
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)
