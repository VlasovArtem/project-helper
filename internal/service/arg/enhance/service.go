package enhance

//go:generate mockgen -destination=mocks/mock_service.go -package=mocks -source=service.go

import (
	"strings"

	"github.com/pkg/errors"
	"project-helper/internal/domain/dto"
	"project-helper/internal/domain/entity"
	"project-helper/internal/utils"
)

type (
	ExtractorService interface {
		ExtractTags(arg entity.Arg) entity.Tags
		ExtractTag(tag entity.Tag) (string, error)
	}
	TagService interface {
		GetTagValue(request *dto.GetTagValueRequest) (string, error)
	}
	PredefinedArgService interface {
		TryToFindPredefinedArgValue(request *dto.TryToFindPredefinedArgRequest) (string, error)
		GetPredefinedArgValues(request *dto.GetPredefinedArgsRequest) ([]string, error)
	}
)

type Service struct {
	extractorService     ExtractorService
	tagService           TagService
	predefinedArgService PredefinedArgService
}

func NewService(
	extractorService ExtractorService,
	tagService TagService,
	predefinedArgService PredefinedArgService,
) *Service {
	return &Service{
		extractorService:     extractorService,
		tagService:           tagService,
		predefinedArgService: predefinedArgService,
	}
}

func (s *Service) EnhanceArgs(request *dto.EnhanceArgsRequest) ([]string, error) {
	err := utils.Validate.Struct(request)
	if err != nil {
		return make([]string, 0), errors.Wrap(err, "request is not valid")
	}

	args := make([]string, len(request.Args))

	for i, arg := range request.Args {
		enhanceTags := s.extractorService.ExtractTags(entity.Arg(arg))

		if len(enhanceTags) == 0 {
			args[i] = arg
			continue
		}

		for _, enhanceTag := range enhanceTags {
			extractedEnhanceTag, err := s.extractorService.ExtractTag(enhanceTag)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to extract tag")
			}

			tagValue, err := s.tagService.GetTagValue(&dto.GetTagValueRequest{
				Operation:    request.Operation,
				Flags:        request.Flags,
				ExtractedTag: extractedEnhanceTag,
			})
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get tag value")
			}

			tagValue, err = s.predefinedArgService.TryToFindPredefinedArgValue(&dto.TryToFindPredefinedArgRequest{
				ParsedTag: extractedEnhanceTag,
				Value:     tagValue,
			})
			if err != nil {
				return nil, errors.Wrapf(err, "failed to try to find predefined arg")
			}

			tagValue = strings.ReplaceAll(arg, string(enhanceTag), tagValue)

			escapeArg, err := utils.EscapeValue(tagValue)
			if err != nil {
				return nil, errors.Wrap(err, "failed to escape arg")
			}

			args[i] = escapeArg
		}
	}

	return args, nil
}

func (s *Service) GetEnhancedOperationArgs(request *dto.GetEnhancedOperationArgs) ([]string, error) {
	if err := utils.Validate.Struct(request); err != nil {
		return nil, errors.Wrap(err, "request is not valid")
	}

	var newArgs []string

	for _, arg := range request.Operation.Args {
		enhanceTags := s.extractorService.ExtractTags(entity.Arg(arg))

		if len(enhanceTags) == 0 {
			newArgs = append(newArgs, arg)

			continue
		}

		var predefinedArgs []string

		for _, enhanceTag := range enhanceTags {
			extractEnhanceTag, err := s.extractorService.ExtractTag(enhanceTag)
			if err != nil {
				return nil, errors.Wrap(err, "failed to extract tag")
			}

			if request.Operation.PredefinedArgsTag.Name == extractEnhanceTag {
				predefinedArgs, err = s.predefinedArgService.GetPredefinedArgValues(&dto.GetPredefinedArgsRequest{
					Flags:             request.Flags,
					PredefinedArgsTag: request.Operation.PredefinedArgsTag,
				},
				)
				if err != nil {
					return nil, errors.Wrap(err, "failed to get predefined args")
				}

				continue
			}
		}

		if len(predefinedArgs) != 0 {
			newArgs = append(newArgs, predefinedArgs...)

			continue
		}

		newArgs = append(newArgs, arg)
	}

	return newArgs, nil
}
