MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --no-builtin-variables


install:
	go install -mod=vendor -v .

generate:
	protoc dynamo.proto --proto_path=dynamo --go_out=import_path=dynamo,paths=source_relative:dynamo

example:
	DEBUG_PGD=true protoc example.proto --proto_path=. --proto_path=examplepb --go_out=examplepb --dynamo_out=examplepb

.PHONY: adddep
adddep:
	go mod tidy -v
	go mod vendor

.PHONY: updatedeps
updatedeps:
	go get -u ./...
	go mod tidy -v
	go mod vendor
