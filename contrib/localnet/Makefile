all: terrad-env terrad-upgrade-env

terrad-env:
	docker build --platform linux/amd64 --no-cache --build-arg UID=$(shell id -u) --build-arg GID=$(shell id -g) --tag classic-terra/terrad-env terrad-env

terrad-upgrade-env:
	docker build --platform linux/amd64 --build-arg UID=$(shell id -u) --build-arg GID=$(shell id -g) --build-arg entrypoint=cosmovisor-entrypoint.sh --tag classic-terra/terrad-upgrade-env terrad-env

.PHONY: all terrad-env terrad-upgrade-env
