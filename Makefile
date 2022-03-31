STAGE ?= poc # set default stage for deployment if not passed from cmdline

.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/connectHandler lambda/connect/main.go
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/disconnectHandler lambda/disconnect/main.go
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/manageHandler lambda/manage/main.go
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/postHandler lambda/post/main.go
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/broadcastHandler lambda/broadcast/main.go
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/listHandler lambda/list/main.go
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/chatHandler lambda/chat/main.go


clean:
	rm -rf ./bin ./vendor go.sum

deploy: build # clean build
	@echo "Deploying stage $(STAGE)"
	npx sls deploy --verbose --stage=$(STAGE)

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh

test:
	go test -v github.com/darren-reddick/go-apigw-webchat/internal/...
