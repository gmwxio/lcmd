export DEBIAN_FRONTEND=noninteractive

apt-get -y install cabal-install
apt-get -y install zlibc zlib1g-dev
cabal update
curl -sSL https://get.haskellstack.org/ | sh
pip3 install awscli
