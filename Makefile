.DEFAULT_GOAL := defaultTarget

# Name of the resulting binary
BINARY=piccolo

GIT_VERSION=$(shell git describe --always --dirty)
GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
GIT_TAG=$(shell git describe --exact-match HEAD 2>/dev/null || true)

# Create -ldflags for go build, inject Git version info
LDFLAGS=-ldflags "-X github.com/shawnsilva/piccolo/version.gitVersion=${GIT_VERSION} -X github.com/shawnsilva/piccolo/version.gitBranch=${GIT_BRANCH}"

GO_PKG_FILES=$(shell go list ./... | grep -v vendor)
GO_FILES_NO_VENDOR = $(shell find . \( ! -regex '.*/\..*' \) -type f -name '*.go' -not -path "./vendor/*")

defaultTarget: clean deps check test build

deps:
	@echo
	@echo "[deps]"
	@echo
	@echo "Installing golint"
	@go get -u golang.org/x/lint/golint

# Build piccolo using your current systems architecture as the target
build: deps
	@echo
	@echo "[build]"
	@echo
	@go build -v ${LDFLAGS} -o build/${BINARY} cmd/piccolo/piccolo.go

buildAll: deps
	env GOOS=linux GOARCH=amd64 go build -v ${LDFLAGS} -o build/linux_amd64/${BINARY} cmd/piccolo/main.go
	env GOOS=linux GOARCH=386 go build -v ${LDFLAGS} -o build/linux_386/${BINARY} cmd/piccolo/main.go
	env GOOS=linux GOARCH=arm go build -v ${LDFLAGS} -o build/linux_arm/${BINARY} cmd/piccolo/main.go
	env GOOS=windows GOARCH=amd64 go build -v ${LDFLAGS} -o build/windows_amd64/${BINARY}.exe cmd/piccolo/main.go
	env GOOS=windows GOARCH=386 go build -v ${LDFLAGS} -o build/windows_386/${BINARY}.exe cmd/piccolo/main.go
	env GOOS=darwin GOARCH=amd64 go build -v ${LDFLAGS} -o build/darwin_amd64/${BINARY} cmd/piccolo/main.go
	env GOOS=darwin GOARCH=386 go build -v ${LDFLAGS} -o build/darwin_386/${BINARY} cmd/piccolo/main.go
	env GOOS=freebsd GOARCH=amd64 go build -v ${LDFLAGS} -o build/freebsd_amd64/${BINARY} cmd/piccolo/main.go
	env GOOS=freebsd GOARCH=386 go build -v ${LDFLAGS} -o build/freebsd_386/${BINARY} cmd/piccolo/main.go
	env GOOS=freebsd GOARCH=arm go build -v ${LDFLAGS} -o build/freebsd_arm/${BINARY} cmd/piccolo/main.go
	env GOOS=solaris GOARCH=amd64 go build -v ${LDFLAGS} -o build/solaris_amd64/${BINARY} cmd/piccolo/main.go

check: deps
	@echo
	@echo "[check]"
	@echo
	@echo "Running go vet..."
	@go vet ${GO_PKG_FILES}
	@echo "Running gofmt (not modifying)..."
	@gofmt -d -l ${GO_FILES_NO_VENDOR} | read && echo "ERROR: gofmt's style checks didn't pass" 1>&2 && exit 1 || true
	@echo "Running golint..."
	@golint -set_exit_status ${GO_PKG_FILES} || (echo "ERROR: golint found errors" 1>&2 && exit 1)

docker-build:
	docker build --tag "shawnsilva/piccolo:latest" .
	@if [ "${GIT_TAG}" != "" ] ; then docker tag "shawnsilva/piccolo:latest" "shawnsilva/piccolo:${GIT_TAG}"; fi

docker-push:
	docker push "shawnsilva/piccolo"

install: deps
	go install ${LDFLAGS} github.com/shawnsilva/piccolo/cmd/piccolo

test: deps
	@echo
	@echo "[test]"
	@echo
	@go test -v -race -cover ${GO_PKG_FILES}

fmt:
	@echo
	@echo "[gofmt]"
	@echo
	@gofmt -l -w ${GO_FILES_NO_VENDOR}

lint:
	@echo
	@echo "[golint]"
	@echo
	@golint ${GO_PKG_FILES}

vet:
	@echo
	@echo "[go vet]"
	@echo
	@go vet ${GO_PKG_FILES}

clean:
	@echo
	@echo "[cleaning]"
	@echo
	@if [ -f build/${BINARY} ] ; then rm build/${BINARY} ; fi
	@go clean -x ${GO_PKG_FILES}

.PHONY: clean check
