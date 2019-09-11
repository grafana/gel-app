DSNAME=gel
GO = GO111MODULE=on go
all: build

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


# TODO: This should build for the current arch, not linux
build:
	$(GO) build -mod=vendor -o ./dist/datasource/${DSNAME}_linux_amd64 -a -tags netgo -ldflags '-w' ./pkg

build-darwin:
	$(GO) build -mod=vendor -o ./dist/datasource/${DSNAME}_darwin_amd64 -a -tags netgo -ldflags '-w' ./pkg

build-dev:
	$(GO) build -mod=vendor -o ./dist/datasource/${DSNAME}_linux_amd64 -a ./pkg

build-win:
	$(GO) build -mod=vendor -o ./dist/datasource/${DSNAME}_windows_amd64.exe -a -tags netgo -ldflags '-w' ./pkg

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