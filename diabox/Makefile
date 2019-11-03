

build:
	docker build -t diabox_base -f diabox_base.dockerfile .
	docker build -t diabox_wx -f diabox_wx.dockerfile .

run:
	docker network create --driver bridge gophers
	docker volume create gopher01 \
	&& docker run \
		-d \
		-p 1022:22 \
		-p 18080:8080 \
		--network gophers \
		--rm \
		--name gopher01 \
		--env U_NAME=gopher01 \
		-v gopher01:/home \
		diabox_wx \
	&& xdg-open http://localhost:18080