LOGXY_BINARY=logxy
LOGXY_IMAGE?=dskard/logxy-build
LOGXY_LOG=logxy.log
LOGXY_PORT=8744
LOGXY_SRC=src/logxy
DOCKER_RUN_CMD_FLAGS=
NETWORK=${PROJECT}_default
PROJECT=logxy
WORKDIR=/opt/work
RESULT_SERVER_XML?=result_server.xml
RESULT_CLIENT_XML?=result_client.xml


DOCKER_RUN_CMD=docker run --rm --init \
	--name logxy \
	--volume=${CURDIR}:${WORKDIR} \
	--workdir=${WORKDIR}/src/logxy \
	-e GOPATH=${WORKDIR} \
	-e GOCACHE=${WORKDIR}/cache \
	-e GIT_COMMITTER_NAME='xxx' \
	-e GIT_COMMITTER_EMAIL='xxx@yyy.com' \
	${DOCKER_RUN_CMD_FLAGS} \
	${LOGXY_IMAGE}

# NOTE: This Makefile does not support running with concurrency (-j XX).
.NOTPARALLEL:


all: dep-check ${LOGXY_BINARY}

build:
	docker build -t ${LOGXY_IMAGE} docker

proxy:
	@$(eval DOCKER_RUN_CMD_FLAGS := --network=${NETWORK} --workdir=${WORKDIR} --user=`id -u`:`id -g`)
	${DOCKER_RUN_CMD} ./${LOGXY_BINARY} --forward-to http://selenium-hub:4444 --port ${LOGXY_PORT} --log ${LOGXY_LOG}

shell:
	@$(eval DOCKER_RUN_CMD_FLAGS := -it --user=`id -u`:`id -g`)
	${DOCKER_RUN_CMD} /bin/bash

${LOGXY_BINARY}:
	@$(eval DOCKER_RUN_CMD_FLAGS := --user=`id -u`:`id -g`)
	${DOCKER_RUN_CMD} make -C ../../ do-all

test: dep-check
	@$(eval DOCKER_RUN_CMD_FLAGS := --user=`id -u`:`id -g`)
	${DOCKER_RUN_CMD} go get github.com/tebeka/go2xunit
	${DOCKER_RUN_CMD} make -C ../../ do-test

dep-init:
	@$(eval DOCKER_RUN_CMD_FLAGS := --user=`id -u`:`id -g`)
	${DOCKER_RUN_CMD} dep init

dep-ensure:
	@$(eval DOCKER_RUN_CMD_FLAGS := --user=`id -u`:`id -g`)
	${DOCKER_RUN_CMD} dep ensure

dep-check:
	@$(eval DOCKER_RUN_CMD_FLAGS := --user=`id -u`:`id -g`)
	${DOCKER_RUN_CMD} dep check

do-all:
	go build -v -o ${LOGXY_BINARY} ${LOGXY_SRC}/main.go

do-test:
	cd ${LOGXY_SRC}; \
	2>&1 go test -gocheck.vv | tee ${RESULT_SERVER_XML}.log || true; \
	../../bin/go2xunit -gocheck -fail -input ${RESULT_SERVER_XML}.log -output ${RESULT_SERVER_XML}

clean:
	rm -rf \
	    client/${RESULT_CLIENT_XML} \
	    src/${RESULT_SERVER_XML} \
	    ${LOGXY_SRC}/${RESULT_SERVER_XML}.log

distclean: clean
	rm -rf \
	    bin \
	    cache \
	    pkg \
	    ${LOGXY_BINARY} \
	    ${LOGXY_LOG}

.PHONY: all clean dep-check dep-init dep-ensure distclean do-all do-test test ${LOGXY_BINARY}
