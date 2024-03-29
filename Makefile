# set default stage for deployment if not passed from cmdline
STAGE ?= poc
REGION ?= eu-west-1
FUNCTIONS = $(shell find lambda -type d -maxdepth 1 -mindepth 1 ! -name utils -exec basename {} \;)
GO111MODULE = on

.PHONY: build clean deploy gomodgen build-% install

install:
	npm install

build: install gomodgen $(addprefix build-,$(FUNCTIONS))

build-%:
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/$*Handler lambda/$*/main.go

clean:
	rm -rf ./bin ./vendor go.sum

remove:
	@echo "Removing stage $(STAGE)"
	npx sls remove --verbose --stage=$(STAGE) --region $(REGION)
	@echo Cleaning up dangling log groups
	node ./scripts/delete-log-group.js --region $(REGION) --log-group '/aws/websocket/go-apigw-webchat-$(STAGE)' --delete

deploy: build
	@echo "Deploying stage $(STAGE)"
	npx sls deploy --verbose --stage=$(STAGE) --region $(REGION)

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh

test:
	go test -v github.com/darren-reddick/go-apigw-webchat/internal/...

e2etest:
	export WEBSOCKET_URL=$$(node ./scripts/aws4-sign.js --url $$(node ./scripts/get-websocket-url.js --stage $(STAGE))); go test ./tests -v



