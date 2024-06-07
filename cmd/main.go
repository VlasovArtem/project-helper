package main

import (
	"context"

	"github.com/rs/zerolog/log"
	"project-helper/internal/service/arg"
	"project-helper/internal/service/arg/enhance"
	"project-helper/internal/service/arg/predefined"
	"project-helper/internal/service/config"
	"project-helper/internal/service/flag"
	"project-helper/internal/service/flag/parser"
	"project-helper/internal/service/operation"
	"project-helper/internal/service/projecthelper"
	"project-helper/internal/service/tag"
	"project-helper/internal/service/tag/extractor"
)

func main() {
	configService, err := config.NewService()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create config service")
	}

	flagParserService := parser.NewService(configService)

	flags, err := flagParserService.ParseFlags()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read flags")
	}

	flagsService := flag.NewFlagsService(flags)
	tagExtractorService := extractor.NewService()
	predefinedArgService := predefined.NewService(configService)
	tagService := tag.NewService(configService)
	enhanceArgService := enhance.NewService(tagExtractorService, tagService, predefinedArgService)

	argService := arg.NewService(flagsService, enhanceArgService, predefinedArgService)

	operationService := operation.NewService(configService)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create operation service")
	}

	service := projecthelper.NewService(operationService, flagsService, argService)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create operation service")
	}

	err = service.Run(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to run operation")
	}
}
