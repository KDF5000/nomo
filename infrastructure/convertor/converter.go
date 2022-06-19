package convertor

import (
	"context"
	"time"
)

type ConvConfig struct {
	Ctx    context.Context
	Url    string
	Params Html2ImageParams
}

type Converter interface {
	Convert(ctx context.Context) ([]byte, error)
	GetTimeElapsed() time.Duration
}
