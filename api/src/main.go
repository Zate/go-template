package main

import (
	"context"
	"flag"
	"log/slog"
)

//	@title			AIchemist AI Toolkit
//	@version		1.0
//	@description	Toolkit for performing some types of AI operations for security engineers

//	@host		localhost:8080
//	@BasePath	/

var (
	Environment string
	AWS_Region  string
	ServiceName string
	XHeaders    map[string]string
	URL         string
	Debug       bool = false
	S           *Service
)

func init() {
	Environment = getEnv("APP_ENV", "local")
	AWS_Region = getEnv("AWS_REGION", "local")
	ServiceName = getEnv("SERVICE_NAME", "no_service_name_set")

	XHeaders = map[string]string{
		"content-type": "application/json",
	}
}

func main() {
	listenPort := flag.Int("port", 8080, "port to listen on")
	flag.Parse()
	S, err := NewService(*listenPort)
	if err != nil {
		panic(err)
	}

	S.Logger = S.Logger.With(slog.Group(
		"service",
		slog.Any("name", ServiceName),
		slog.Any("environment", Environment),
		slog.Any("region", AWS_Region),
		slog.Any("port", *listenPort),
	))
	S.Logger.LogAttrs(
		context.Background(),
		slog.LevelInfo,
		"STARTING_SERVICE",
	)

	err = S.Run()
	if err != nil {
		panic(err)
	}
}
