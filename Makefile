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

test-in-docker: build-container
	docker run --rm \
	  --network host \
		-v "${PWD}":/go/src/github.com/grafana/${DSNAME}\
		-w /go/src/github.com/grafana/${DSNAME} \
		plugin-builder make test

copy-pluginj:
	mkdir -p dist
	cp ./plugin.json ./dist/plugin.json

# TODO: This should build for the current arch, not linux
build: copy-pluginj
	$(GO) build -o ./dist/${DSNAME}_linux_amd64 -tags netgo -ldflags '-w' ./pkg

build-darwin:
	$(GO) build -o ./dist/${DSNAME}_darwin_amd64 -tags netgo -ldflags '-w' ./pkg

build-dev:
	$(GO) build -o ./dist/${DSNAME}_linux_amd64 ./pkg

build-win:
	$(GO) build -o ./dist/${DSNAME}_windows_amd64.exe -tags netgo -ldflags '-w' ./pkg

build-in-circleci: build-in-circleci-linux build-in-circleci-windows

build-in-circleci-linux:
	$(GO) build -o /output/${DSNAME}_linux_amd64 -a -tags netgo -ldflags '-w' ./pkg

build-in-circleci-windows:
	CGO_ENABLED=1 GOOS=windows CC=/usr/bin/x86_64-w64-mingw32-gcc \
	PKG_CONFIG_PATH=/usr/lib/pkgconfig_win \
	$(GO) build -o /output/${DSNAME}_windows_amd64.exe -a -tags netgo -ldflags '-w' ./pkg

build-in-docker: build-container
	docker run --rm \
		-v "${PWD}":/go/src/github.com/grafana/${DSNAME} \
		-w /go/src/github.com/grafana/${DSNAME} \
		plugin-builder make build

build-in-docker-win: build-container
	docker run --rm \
		-e "CGO_ENABLED=1" -e "GOOS=windows" \
		-e "CC=/usr/bin/x86_64-w64-mingw32-gcc" -e "PKG_CONFIG_PATH=/usr/lib/pkgconfig_win" \
		-v "${PWD}":/go/src/github.com/grafana/${DSNAME} \
		-w /go/src/github.com/grafana/${DSNAME} \
		plugin-builder make build-win

build-container:
	docker build --tag plugin-builder .

build-container-rebuild:
	docker build --tag plugin-builder --no-cache=true .
