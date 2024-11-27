package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

func PaginateVerses(text string, page, limit int, log *logrus.Logger) ([]string, error) {
	log.Infof("PaginateVerses called with page: %d, limit: %d", page, limit)

	if page < 1 {
		page = 1
		log.Warn("pagination: Page number is less than 1, defaulting to 1")
	}
	if limit < 1 {
		limit = 10
		log.Warn("pagination: Limit is less than 1, defaulting to 10")
	}

	verses := strings.Split(text, "\\n\\n")
	log.Infof("Total verses found: %d", len(verses))

	start := (page - 1) * limit
	end := start + limit
	log.Debugf("Calculated start index: %d, end index: %d", start, end)

	if start >= len(verses) {
		err := fmt.Errorf("pagination: Page number is out of range")
		log.Errorf("Error: %v", err)
		return nil, err
	}
	if end > len(verses) {
		end = len(verses)
		log.Debugf("Adjusted end index to: %d", end)
	}

	paginatedVerses := verses[start:end]
	log.Infof("Returning %d verses for page: %d", len(paginatedVerses), page)

	return paginatedVerses, nil
}
