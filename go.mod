module github.com/pquerna/protoc-gen-dynamo

go 1.24

require (
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.43.1
	github.com/dave/jennifer v1.7.1
	github.com/klauspost/compress v1.18.0
	github.com/lyft/protoc-gen-star/v2 v2.0.4
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/aws/smithy-go v1.22.3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/spf13/afero v1.14.0 // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	golang.org/x/tools v0.33.0 // indirect
)

// https://github.com/lyft/protoc-gen-star/pull/132
replace github.com/lyft/protoc-gen-star/v2 => github.com/pquerna/protoc-gen-star/v2 v2.0.0-20250415201647-653a078eb414
