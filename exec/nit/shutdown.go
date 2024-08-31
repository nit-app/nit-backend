package main

import (
	"go.uber.org/zap"
	"io"
	"os"
)

func shutdown(c chan os.Signal, s io.Closer) {
	<-c
	zap.S().Info("stopping")
	_ = s.Close()
}
