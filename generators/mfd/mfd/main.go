package main

import (
	"os"

	"github.com/dizzyfool/genna/generators/mfd"

	"go.uber.org/zap"
)

func main() {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.Encoding = "console"
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	root := mfd.CreateCommand(logger)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
