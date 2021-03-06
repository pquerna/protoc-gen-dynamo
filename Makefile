MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --no-builtin-variables

.PHONY: install
install:
	go install -mod=vendor -v .

generate:
	protoc \
	-I . \
	--go_out="plugins=grpc,paths=source_relative:." \
	dynamo/*.proto

example:
	DEBUG_PGD=true protoc example.proto --proto_path=. --proto_path=examplepb --go_out="plugins=grpc,paths=source_relative:examplepb" --dynamo_out="lang=go,paths=source_relative:examplepb"

.PHONY: adddep
adddep:
	go mod tidy -v
	go mod vendor

.PHONY: updatedeps
updatedeps:
	go get -u ./...
	go mod tidy -v
	go mod vendor
