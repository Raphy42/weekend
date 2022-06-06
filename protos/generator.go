//go:build ignore
// +build ignore

package main

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/Raphy42/weekend/core/logger"
)

func main() {
	log := logger.New()

	if len(os.Args) != 2 {
		panic("not enough arguments given to `go:generate generator.go`")
	}

	path := os.Args[1]
	log.Info("starting grpc codegen", zap.String("dir.out", path))
	entries, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	protos := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".proto" {
			protos = append(entry.Name())
		}
	}
}
