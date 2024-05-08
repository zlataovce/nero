package main

import (
	"go.uber.org/zap"
)

// appContext is the context of the CLI application.
type appContext struct {
	logger *zap.Logger
}
