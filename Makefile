TAGS ?= "sqlite"
GO_BIN ?= go

install:
	packr2
	$(GO_BIN) install -tags ${TAGS} -v ./.
	make tidy

tidy:
	$(GO_BIN) mod tidy

build:
	packr2
	$(GO_BIN) build -v .
	make tidy

test:
	packr2 clean
	$(GO_BIN) test -tags ${TAGS} ./...
	make tidy

