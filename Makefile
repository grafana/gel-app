DSNAME=gel
GO = GO111MODULE=on go
all: build

# Called by circle-ci task
backend-plugin-ci: test build-in-circleci

test:
	mkdir -p coverage
	$(GO) test ./pkg/... -v -cover -coverprofile=coverage/cover.out
	$(GO) tool cover -html=coverage/cover.out -o coverage/coverage.html

vendor:
	$(GO) mod vendor

copy-artifacts:
	mkdir -p dist
	cp ./plugin.json ./dist/plugin.json
	cp ./README.md ./dist/README.md

# TODO: This should build for the current arch, not linux
build: copy-artifacts
	$(GO) build -o ./dist/${DSNAME}_linux_amd64 -tags netgo -ldflags '-w' ./pkg

build-darwin: copy-artifacts
	$(GO) build -o ./dist/${DSNAME}_darwin_amd64 -tags netgo -ldflags '-w' ./pkg

build-win: copy-artifacts
	$(GO) build -o ./dist/${DSNAME}_windows_amd64.exe -tags netgo -ldflags '-w' ./pkg

build-in-circleci: build-in-circleci-linux build-in-circleci-windows

build-in-circleci-linux:
	$(GO) build -o /output/${DSNAME}_linux_amd64 -a -tags netgo -ldflags '-w' ./pkg

build-in-circleci-windows:
	CGO_ENABLED=1 GOOS=windows CC=/usr/bin/x86_64-w64-mingw32-gcc \
	PKG_CONFIG_PATH=/usr/lib/pkgconfig_win \
	$(GO) build -o /output/${DSNAME}_windows_amd64.exe -a -tags netgo -ldflags '-w' ./pkg

