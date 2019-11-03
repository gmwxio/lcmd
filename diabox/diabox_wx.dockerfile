FROM diabox_base as builder

FROM ubuntu:19.04

ARG DEBIAN_FRONTEND=noninteractive

RUN echo "Australia/Sydney" > /etc/timezone

RUN apt-get update \
	&& apt-get install -y apt-transport-https \
	&& apt-get install -y apt-utils \
	&& apt-get install -y ca-certificates \
	&& apt-get install -y curl \
	&& apt-get install -y gnupg2 \
	&& apt-get install -y software-properties-common

# Yarn
RUN  curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - \
	&& echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list
# Add source repository for recent node
RUN curl -sL https://deb.nodesource.com/setup_8.x | bash -
# Docker
RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add - \
	&& add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu cosmic stable"
# docker-compose
RUN curl -L "https://github.com/docker/compose/releases/download/1.24.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose \
	&& chmod a+x /usr/local/bin/docker-compose
# git-lfs
RUN curl -s https://packagecloud.io/install/repositories/github/git-lfs/script.deb.sh | bash

# 
# 
# 

RUN apt-get update -y && apt-get install -y \
	apt-utils \
	build-essential \
	ca-certificates \
	curl \
	dbus \
	dnsutils \
	docker-ce \
	g++ \
	git \
	git-lfs \
	gnupg2 \
	iputils-ping \
	libbz2-dev \
	libffi-dev \
	libpcre3-dev \
	libssl-dev \
	libzip-dev \
	net-tools \
	netcat \
	nfs-common \
	nodejs \
	openjdk-8-jdk \
	pkg-config \
	python \
	python-dev \
	python-pip \
	python3 \
	python3-pip \
	rpcbind \
	ssh \
	sudo \
	tcpdump \
	vim \
	unzip \
	wget \
	yarn \
	zip \
	zlib1g-dev \
	zlibc

# 
# 
# 

RUN service ssh start

RUN pip3 install \
	doit \
	pystache \
	awscli
RUN python3 -m pip install -U \
	pylint
RUN pip install \
	supervisor

RUN VERSION=0.11.1 \
	&& INSTALLER=bazel-$VERSION-installer-linux-x86_64.sh \
	&& wget -q https://github.com/bazelbuild/bazel/releases/download/$VERSION/$INSTALLER \
	&& chmod +x $INSTALLER \
	&& ./$INSTALLER \
	&& rm $INSTALLER \
	&& /usr/local/bin/bazel version \
	&& rm -rf /root/.cache/bazel/

# Install go from a binary release
RUN wget -q https://dl.google.com/go/go1.13.1.linux-amd64.tar.gz \
	&& tar -C /usr/local -xzf go1.13.1.linux-amd64.tar.gz \
	&& rm go1.13.1.linux-amd64.tar.gz

# 
# 
# 

ENV CODE_SERVER_VERSION=2.1650-vsc1.39.2
RUN wget -q -O /code-server2.tgz https://github.com/cdr/code-server/releases/download/${CODE_SERVER_VERSION}/code-server${CODE_SERVER_VERSION}-linux-x86_64.tar.gz
RUN tar -C /usr/local -xzf /code-server2.tgz \
	&& rm /code-server2.tgz \
	&& ln -s /usr/local/code-server${CODE_SERVER_VERSION}-linux-x86_64 /usr/local/code-server

ENV WX_VERION=0.0.6
RUN wget -q -O /wx.tgz https://github.com/wxio/wx/releases/download/v${WX_VERION}/wx_${WX_VERION}_linux_x86_64.tar.gz \
	&& tar -C /usr/local/go/bin -xzf /wx.tgz \
	&& rm /wx.tgz

# 
# 
# 

COPY root /

#
COPY --from=builder /go/bin/* /go/bin/
COPY --from=builder /go/include /go/include
# GOPRIVATE="bitbucket.org"

RUN echo 'complete -C /go/bin/godna godna' >> /etc/bash.bashrc

EXPOSE 22 8080

ENV U_NAME=gopher01
#echo '20091110' | openssl passwd -1 -stdin
ENV U_PASSWDHASH="$1$M97Wg8Ay$8T8LA4XgdMjYs7W1Gs1p.."
ENV U_GID=433
ENV U_UID=431

CMD ["/startup.sh"]
