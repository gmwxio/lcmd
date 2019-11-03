FROM golang:1.13 as builder

# protoc
ENV PROTOC_VERSION 3.9.0
RUN apt-get update && \
	apt-get upgrade -y -o Dpkg::Options::="--force-confold" && \
	apt-get install -y unzip
RUN wget -q https://github.com/google/protobuf/releases/download/v$PROTOC_VERSION/protoc-$PROTOC_VERSION-linux-x86_64.zip && \
   unzip protoc-$PROTOC_VERSION-linux-x86_64.zip && \
   rm -f protoc-$PROTOC_VERSION-linux-x86_64.zip
#
# protoc-gen-go
ENV GOLANG_PROTOBUF_VERSION v1.2.0
RUN mkdir -p /go/src/github.com/golang && \
   cd /go/src/github.com/golang && \
   git clone https://github.com/golang/protobuf.git && \
   cd protobuf && \
   git checkout $GOLANG_PROTOBUF_VERSION && \
   cd protoc-gen-go && \
   go install
#
# grpc-gateway & swagger
ENV GRPC_GATEWAY_VERSION v1.11.0
RUN git clone https://github.com/grpc-ecosystem/grpc-gateway.git && \
	cd grpc-gateway && \
	git checkout $GRPC_GATEWAY_VERSION && \
	cd protoc-gen-grpc-gateway && \
	go install && \
	cd ../protoc-gen-swagger && \
	go install
#
# micro
ENV MICRO_VERSION v0.8.0
RUN git clone https://github.com/micro/protoc-gen-micro.git && \
	cd protoc-gen-micro && \
	git checkout $MICRO_VERSION && \
	go install
#
# validate
ENV VALIDATE_VERSION v0.1.0
RUN git clone https://github.com/envoyproxy/protoc-gen-validate.git && \
	cd protoc-gen-validate && \
	git checkout $VALIDATE_VERSION && \
	go mod init github.com/envoyproxy/protoc-gen-validate && \
	go install
#
# godna
ENV GODNA_VER v1.15.0
RUN git clone https://github.com/wxio/godna.git  && \
	cd godna && \
	git checkout $GODNA_VER && \
	go install -ldflags  "-X main.version=$GODNA_VER -X main.commit=$(git log --pretty=format:\"%H\" -1) -X main.date=$(git log --pretty=format:\"%ad\" -1 --date=format:'%F-%T%z')"

# Go tools for vs-code extension
RUN go get github.com/stamblerre/gocode
RUN go get github.com/uudashr/gopkgs/cmd/gopkgs
RUN go get github.com/ramya-rao-a/go-outline
RUN go get github.com/acroca/go-symbols
RUN go get golang.org/x/tools/cmd/guru
RUN go get golang.org/x/tools/cmd/gorename
RUN go get github.com/fatih/gomodifytags
RUN go get github.com/haya14busa/goplay/cmd/goplay
RUN go get github.com/josharian/impl
RUN go get github.com/tylerb/gotype-live
RUN go get github.com/rogpeppe/godef
RUN go get github.com/zmb3/gogetdoc
RUN go get golang.org/x/tools/cmd/goimports
RUN go get github.com/sqs/goreturns
RUN go get winterdrache.de/goformat/goformat
RUN go get golang.org/x/lint/golint
RUN go get github.com/cweill/gotests/...
RUN go get honnef.co/go/tools/...
RUN go get github.com/golangci/golangci-lint/cmd/golangci-lint
RUN go get github.com/mgechev/revive
RUN go get github.com/sourcegraph/go-langserver
RUN go get golang.org/x/tools/gopls
RUN go get github.com/go-delve/delve/cmd/dlv
RUN go get github.com/davidrjenni/reftools/cmd/fillstruct
RUN go get github.com/godoctor/godoctor
  
#    \
#   && go get -u -v golang.org/x/tools/cmd/godoc
