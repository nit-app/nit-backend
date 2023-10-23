package main

import (
	"go.uber.org/zap"
	"io"
	"os"
)

func shutdown(c chan os.Signal, s io.Closer) {
	zap.S().Info("stopping application")
	<-c
	_ = s.Close()
}
