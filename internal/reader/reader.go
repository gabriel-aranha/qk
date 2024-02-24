package reader

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

type Reader struct {
	logger *zap.Logger
}

func NewReader(logger *zap.Logger) Reader {
	var reader Reader
	reader.logger = logger

	return reader
}

func (r *Reader) Read(filePath string) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		r.logger.Error("error reading input file", zap.Error(err))
		return nil, err
	}
	lines := strings.Split(string(content), "\n")

	return lines, nil
}
