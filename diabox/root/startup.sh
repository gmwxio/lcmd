#!/bin/sh

groupadd -r ${U_NAME} -g ${U_GID}
useradd --uid ${U_UID} \
	--gid ${U_NAME} \
	--create-home \
	--home-dir /home/${U_NAME} \
	--shell /bin/bash ${U_NAME} \
	--password ${U_PASSWDHASH}

usermod -a -G wx-sudoers ${U_NAME}
usermod -a -G docker ${U_NAME}

service ssh start

cat > /supervisord.conf << EOF
[unix_http_server]
file = /var/run/supervisord.sock

[supervisord]
nodaemon=true

[rpcinterface:supervisor]
supervisor.rpcinterface_factory = supervisor.rpcinterface:make_main_rpcinterface

[supervisorctl]
serverurl=unix:///var/run/supervisord.sock

[program:code-server]
command=su - $U_NAME -c '/usr/local/code-server/code-server --auth none'
EOF

rm /var/run/supervisord.sock
supervisord -c /supervisord.conf
