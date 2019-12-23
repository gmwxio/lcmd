export DEBIAN_FRONTEND=noninteractive

apt-get update \
	&& apt-get install -y apt-transport-https \
	&& apt-get install -y apt-utils \
	&& apt-get install -y ca-certificates \
	&& apt-get install -y curl \
	&& apt-get install -y gnupg2 \
	&& apt-get install -y software-properties-common
