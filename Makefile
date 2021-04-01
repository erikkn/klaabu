NAME := klaabu

BUILD  := ${CURDIR}/build
BIN    := ${BUILD}/bin
VERSION := dev

.PHONY: clean tidy compile shasum test run
.DEFAULT_GOAL := build

clean:
	rm -rf ${BUILD}

tidy:
	go mod tidy -v

fmt:
	go fmt github.com/erikkn/klaabu/...

vendor:
	go mod vendor

compile:
ifdef GOOS
ifdef GOARCH
	go build -mod=readonly -ldflags="-X 'main.Version=$${VERSION}'" -o ${BIN}/${NAME}-$(GOOS)-$(GOARCH) cli/*.go
endif
else
	go build -mod=readonly -ldflags="-X 'main.Version=$${VERSION}'" -o ${BIN}/${NAME} cli/*.go
endif

shasum:
	cd ${BIN} ; for binary in ./* ; do sha512sum $$binary > ${BUILD}/$$binary.sha512 ; done

test:
	go test -v ./...

build: tidy fmt vendor compile shasum

run: build
	${BIN}
