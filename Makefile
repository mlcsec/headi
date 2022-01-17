before.build:
ifeq (,$(wildcard ./go.mod))
	go mod init v0 && go mod tidy
endif
	go mod download && go mod vendor

build.headi:
	@echo "build in ${PWD}";go build -o headi main.go