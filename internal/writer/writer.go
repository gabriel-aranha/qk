package writer

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/gabriel-aranha/qk/internal/types"
	"go.uber.org/zap"
)

type Writer struct {
	logger *zap.Logger
}

func NewWriter(logger *zap.Logger) Writer {
	var writer Writer
	writer.logger = logger

	return writer
}

func (w *Writer) Write(games types.Games) error {
	cwd, err := os.Getwd()
	if err != nil {
		w.logger.Error("error getting cwd", zap.Error(err))
		return err
	}

	filePath := filepath.Join(cwd, "output", "report.json")

	dirPath := filepath.Dir(filePath)

	_, err = os.Stat(dirPath)

	if os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			w.logger.Error("error creating directory", zap.Error(err))
			return err
		}
	}

	file, err := os.Create(filePath)
	if err != nil {
		w.logger.Error("error creating output file", zap.Error(err))
		return err
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(games, "", "  ")
	if err != nil {
		w.logger.Error("error marshalling output file", zap.Error(err))
		return err
	}

	_, err = file.Write(jsonData)
	if err != nil {
		w.logger.Error("error writing to output file", zap.Error(err))
		return err
	}

	return nil
}
