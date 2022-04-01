STAGE ?= poc # set default stage for deployment if not passed from cmdline
FUNCTIONS = $(shell find lambda -type d -maxdepth 1 -mindepth 1 ! -name utils -exec basename {} \;)
GO111MODULE = on

.PHONY: build clean deploy gomodgen build-%

build: gomodgen $(addprefix build-,$(FUNCTIONS))

build-%:
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/$*Handler lambda/$*/main.go

clean:
	rm -rf ./bin ./vendor go.sum

remove:
	@echo "Removing stage $(STAGE)"
	npx sls remove --verbose --stage=$(STAGE)

deploy: build # clean build
	@echo "Deploying stage $(STAGE)"
	npx sls deploy --verbose --stage=$(STAGE)

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh

test:
	go test -v github.com/darren-reddick/go-apigw-webchat/internal/...
