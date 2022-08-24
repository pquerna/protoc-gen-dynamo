MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --no-builtin-variables

.PHONY: install
install:
	go install -mod=vendor -v .

.PHONY: generate
generate:
	protoc \
	-I . \
	--go_out="paths=source_relative:." \
	--go-grpc_out="paths=source_relative:." \
	dynamo/v1/*.proto

example:
	DEBUG_PGD=true protoc example.proto --proto_path=. --proto_path=examplepb --go_out="paths=source_relative:examplepb" --dynamo_out="lang=go,paths=source_relative:examplepb"

.PHONY: adddep
adddep:
	go mod tidy -v
	go mod vendor

.PHONY: updatedeps
updatedeps:
	go get -u ./...
	go mod tidy -v
	go mod vendor
