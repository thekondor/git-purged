.PHONY: default build full_build test

export GOPATH:=${PWD}

BUILD_DIR=${PWD}/build
APP_NAME=git-purged
BUILD_BIN=${BUILD_DIR}/${APP_NAME}
VERSION_AUTOGEN_FN=version_autogen.go
BUILD_VERSION_FN=${BUILD_DIR}/version_autogen.go

default: build

build: _build_dir _version_file
	@echo - Build "${APP_NAME}" as '${BUILD_BIN}'
	@go build -v -o ${BUILD_BIN} ./*.go

full_build: test build

test: _test_deps
	@echo Run tests for supplimentary code
	@go test -v git_purged


_test_deps:
	@echo Download test dependencies
	@go get -v -d -t ./src/git_purged

_build_dir:
	@echo - Prepare build directory
	@test -d ${BUILD_DIR} || mkdir ${BUILD_DIR}

_version_file: _build_dir
	@echo - Prepare build version
	@cp -f version.go.template ${BUILD_VERSION_FN}

	@echo - Collecting build version details
	$(eval GIT_ORIGIN:=$(shell git config --get remote.origin.url || (echo -n local/; git rev-parse --abbrev-ref HEAD) ))
	@sed -i -e 's^%BUILD_GIT_ORIGIN%^${GIT_ORIGIN}^g' ${BUILD_VERSION_FN}
	$(eval GIT_COMMIT:=$(shell git log -1 --format=%h))
	@sed -i -e 's^%BUILD_GIT_COMMIT%^${GIT_COMMIT}^g' ${BUILD_VERSION_FN}
	$(eval BUILD_DATE:=$(shell date +%s))
	@sed -i -e 's^%BUILD_DATE%^${BUILD_DATE}^g' ${BUILD_VERSION_FN}
	@echo - Build version ORIGIN: ${GIT_ORIGIN}, COMMIT: ${GIT_COMMIT}, DATE: ${BUILD_DATE}

	@ln -sf ${BUILD_VERSION_FN} ./${VERSION_AUTOGEN_FN}

clean:
	@echo - Clean build artefacts up
	@rm -rf ${BUILD_DIR}
	@test -f ./${VERSION_AUTOGEN_FN} && unlink ./${VERSION_AUTOGEN_FN} || /bin/true
