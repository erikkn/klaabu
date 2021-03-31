NAME := klaabu

BUILD  := ${CURDIR}/build
BIN    := ${BUILD}/bin/${NAME}

.PHONY: clean tidy compile test run
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
	env
	go build -mod=readonly -o ${BIN}-$(GOOS)-$(GOARCH) cli/*.go
endif
else
	go build -mod=readonly -o ${BIN} cli/*.go
endif

test:
	go test -v ./...

build: tidy fmt vendor compile

run: build
	${BIN}
