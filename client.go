package main

import (
	"context"
	"google.golang.org/api/logging/v2"
	"google.golang.org/api/option"
)

type Logging struct {
	svc *logging.Service
}

type Tag string

const (
	LEVEL      Tag = "LEVEL"
	IsInternal Tag = "I"
)

type (
	Level string
)

const (
	WARNING Level = "Warning"
)

func NewLogging(conf []byte) (*Logging, error) {
	ctx := context.Background()
	service, err := logging.NewService(ctx, option.WithCredentialsJSON(conf))
	if err != nil {
		return nil, err
	}
	return &Logging{svc: service}, err
}
