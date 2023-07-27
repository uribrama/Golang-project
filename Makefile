.PHONY: default copy-lambdas lint lint-lambdas

project := $(shell basename $(shell git rev-parse --show-toplevel))

# Find all packages under ./lambdas and trim to the path under ./lambdas
lambdas := $(shell go list ./lambdas/... | sed 's|.*/lambdas/||')

default: copy-lambdas

local:
	GO_ENV=development PROJECT=Golang-project go run ./lambdas/persons

build lambdas copy-lambdas:
	mkdir -p artifacts
	mkdir -p out
	$(MAKE) -j5 $(lambdas)

$(lambdas):
	mkdir -p artifacts
	PACKAGE_BIN=$(shell echo $@ | tr '/' '_'); \
		    CGO_ENABLED=0 \
		    GOOS=linux \
			GOARCH=amd64 \
		    go build \
		    -trimpath \
		    -a \
		    -installsuffix cgo \
		    -o out/"$${PACKAGE_BIN}" \
		    lambdas/$@/*.go \
		    && $(MAKE) zip-$${PACKAGE_BIN}
	cp aws/deployFunction.sh artifacts/

zip-%:
	touch -t 198001010000 out/$*
	zip -Xj artifacts/$*.zip out/$*
	-zip -ur artifacts/$*.zip config

lint-lambdas:
	go fmt ./lambdas/...
	git diff --exit-code

deploy:
	serverless deploy --stage dev --aws-profile serverless