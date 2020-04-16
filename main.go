package main

import (
	"github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"github.com/pquerna/protoc-gen-dynamo/internal/pgd"
)

func main() {
	pgs.Init(pgs.DebugEnv("DEBUG_PGD")).
		RegisterModule(pgd.New()).
		RegisterPostProcessor(pgsgo.GoFmt()).
		Render()
}
