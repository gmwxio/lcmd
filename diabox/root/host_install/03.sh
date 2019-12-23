pip3 install \
	doit \
	pystache \
	awscli
python3 -m pip install -U \
	pylint
pip install \
	supervisor

VERSION=0.11.1 \
	&& INSTALLER=bazel-$VERSION-installer-linux-x86_64.sh \
	&& wget -q https://github.com/bazelbuild/bazel/releases/download/$VERSION/$INSTALLER \
	&& chmod +x $INSTALLER \
	&& ./$INSTALLER \
	&& rm $INSTALLER \
	&& /usr/local/bin/bazel version \
	&& rm -rf /root/.cache/bazel/
