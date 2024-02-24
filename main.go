package main

import (
	"github.com/gabriel-aranha/qk/internal/parser"
	"github.com/gabriel-aranha/qk/internal/reader"
	"github.com/gabriel-aranha/qk/internal/writer"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	reader := reader.NewReader(logger)
	arrayLines, err := reader.Read("./input/games.log")
	if err != nil {
		logger.Error("error reading file", zap.Error(err))
		return
	}

	parser := parser.NewParser(logger)
	games, err := parser.Parse(arrayLines)
	if err != nil {
		logger.Error("error parsing file", zap.Error(err))
		return
	}

	writer := writer.NewWriter(logger)
	err = writer.Write(games)
	if err != nil {
		logger.Error("error writing file", zap.Error(err))
		return
	}
}
