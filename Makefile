# Name of the resulting binary
BINARY=piccolo

GIT_VERSION=`git describe --always --dirty`
GIT_BRANCH=`git rev-parse --abbrev-ref HEAD`

# Create -ldflags for go build, inject Git version info
LDFLAGS=-ldflags "-X github.com/shawnsilva/piccolo/version.gitVersion=${GIT_VERSION} -X github.com/shawnsilva/piccolo/version.gitBranch=${GIT_BRANCH}"

deps:
	go get -v ./...

# Build piccolo using your current systems architecture as the target
build: deps
	go build -v ${LDFLAGS} -o build/${BINARY} cmd/piccolo/piccolo.go

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

docker:
	docker build --tag "shawnsilva/piccolo:latest" --tag "shawnsilva/piccolo:${GIT_VERSION}" .

install: deps
	go install ${LDFLAGS} github.com/shawnsilva/piccolo/cmd/piccolo

test: deps
	go test -v -race ./...

clean:
	if [ -f build/${BINARY} ] ; then rm build/${BINARY} ; fi
	go clean -x ./...

.PHONY: clean build
