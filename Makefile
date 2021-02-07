noscl: $(shell find . -name "*.go")
	go build -ldflags="-s -w" -o ./noscl

dist: $(shell find . -name "*.go")
	mkdir -p dist
	gox -ldflags="-s -w" -osarch="windows/amd64 darwin/amd64 linux/386 linux/amd64 linux/arm freebsd/amd64" -output="dist/noscl_{{.OS}}_{{.Arch}}"
