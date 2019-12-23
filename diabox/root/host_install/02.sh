export DEBIAN_FRONTEND=noninteractive
# Yarn
curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - \
	&& echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list
# Add source repository for recent node
curl -sL https://deb.nodesource.com/setup_8.x | bash -
# Docker
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add - \
	&& add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu cosmic stable"
# docker-compose
curl -L "https://github.com/docker/compose/releases/download/1.24.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose \
	&& chmod a+x /usr/local/bin/docker-compose
# git-lfs
curl -s https://packagecloud.io/install/repositories/github/git-lfs/script.deb.sh | bash

# 
# 
# 

apt-get update -y && apt-get install -y \
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

echo "## add user to docker group"
echo "sudo usermod -a -G docker USER"