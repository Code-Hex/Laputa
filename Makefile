tcpserver = ./cmd/laputa/main.go
unixserver = ./cmd/balus/main.go
tcppid = $(PWD)/laputa.pid
unixpid = $(PWD)/balus.pid
path = /tmp/laputa.sock
port = 8080

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

tcp-run:
	@carton exec start_server --port=$(port) --pid-file=$(tcppid) -- ./laputa

tcp-restart:
	@cat $(tcppid) | xargs kilgl -HUP

tcp-stop:
	@cat $(tcppid) | xargs kill -TERM

unix-run:
	@carton exec start_server --path=$(path) --pid-file=$(unixpid) -- ./balus

unix-restart:
	@cat $(unixpid) | xargs kill -HUP

unix-stop:
	@cat $(unixpid) | xargs kill -TERM

run-staging:
	@carton exec start_server --port=$(port) --pid-file=$(tcppid) -- ./laputa
	@carton exec start_server --path=$(path) --pid-file=$(unixpid) -- ./balus

restart-staging: tcp-restart unix-restart