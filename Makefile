MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --no-builtin-variables

.PHONY: build
build:
	mkdir -p build
	go build -mod=readonly -v -o build/ .

.PHONY: generate
generate:
	buf generate proto

.PHONY: example
example: build
	buf --debug generate --template buf.example.gen.yaml --path examplepb/v1

.PHONY: fmt
fmt:
	buf format -w 

.PHONY: lint
lint:
	buf lint ./proto


.PHONY: adddep
adddep:
	go mod tidy -v
	go mod vendor

.PHONY: updatedeps
updatedeps:
	go get -u ./...
	go mod tidy -v
	go mod vendor
