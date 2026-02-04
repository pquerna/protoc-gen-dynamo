module github.com/pquerna/protoc-gen-dynamo

go 1.24.0

require (
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.55.0
	github.com/cespare/xxhash/v2 v2.3.0
	github.com/dave/jennifer v1.7.1
	github.com/klauspost/compress v1.18.3
	github.com/lyft/protoc-gen-star/v2 v2.0.4
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/aws/smithy-go v1.24.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	golang.org/x/mod v0.32.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	golang.org/x/tools v0.41.0 // indirect
)

// https://github.com/lyft/protoc-gen-star/pull/132
replace github.com/lyft/protoc-gen-star/v2 => github.com/pquerna/protoc-gen-star/v2 v2.0.0-20250415201647-653a078eb414
