package extractor

import (
	"regexp"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"project-helper/internal/domain/entity"
	domainerrors "project-helper/internal/domain/errors"
)

type Service struct {
	extractorRegexp *regexp.Regexp
}

func NewService() *Service {
	return &Service{
		extractorRegexp: regexp.MustCompile("\\$\\{\\{([a-zA-Z-0-9]+)\\}\\}"),
	}
}

func (s *Service) ExtractTags(arg entity.Arg) entity.Tags {
	allString := s.extractorRegexp.FindAllString(string(arg), -1)

	tags := make(entity.Tags, len(allString))

	for i, tag := range allString {
		tags[i] = entity.Tag(tag)

	}

	return tags
}

func (s *Service) ExtractTag(tag entity.Tag) (string, error) {
	tagValues := s.extractorRegexp.FindStringSubmatch(string(tag))
	if len(tagValues) == 0 {
		return "", errors.Wrapf(domainerrors.ErrorTagValueNotFound, "tag '%s' is not valid", tag)
	}

	extractedTagValue := tagValues[1]

	log.Debug().Str("tag", string(tag)).Str("tag.extracted", extractedTagValue).Msg("Extracted tag")

	return extractedTagValue, nil
}
