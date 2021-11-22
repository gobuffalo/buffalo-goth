TAGS ?= "sqlite"
GO_BIN ?= go

install:
	$(GO_BIN) install -tags ${TAGS} -v ./.
	make tidy

tidy:
	$(GO_BIN) mod tidy

build:
	$(GO_BIN) build -v .
	make tidy

test:
	$(GO_BIN) test -tags ${TAGS} ./...
	make tidy
