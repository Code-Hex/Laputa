tcpserver = ./cmd/laputa/main.go
unixserver = ./cmd/balus/main.go
tcppid = $(PWD)/$(tcpserver).pid
unixpid = $(PWD)/$(unixserver).pid
port = 8080
path = /tmp/laputa.sock

tcp-run:
	@$(GOPATH)/bin/start_server --port=$(port) --pid-file=$(tcppid) -- ./laputa

tcp-restart:
	@cat $(tcppid) | xargs kill -HUP

tcp-stop:
	@cat $(tcppid) | xargs kill -TERM

unix-run:
	@$(GOPATH)/bin/start_server --path=$(path) --pid-file=$(unixpid) -- ./balus

unix-restart:
	@cat $(unixpid) | xargs kill -HUP

unix-stop:
	@cat $(unixpid) | xargs kill -TERM

build-dev:
	@echo building laputa
	@go build -o laputa -ldflags='-X main.mode=develop' $(tcpserver)
	@echo building balus
	@go build -o balus -ldflags='-X main.mode=develop' $(unixserver)

build-staging:
	@echo building laputa
	@go build -o laputa -ldflags='-X main.mode=staging' $(tcpserver)
	@echo building balus
	@go build -o balus -ldflags='-X main.mode=staging' $(unixserver)

