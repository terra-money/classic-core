all: terrad-env

terrad-env:
	docker build --build-arg UID=$(shell id -u) --build-arg GID=$(shell id -g) --tag terramoney/terrad-env terrad-env

.PHONY: all terrad-env
