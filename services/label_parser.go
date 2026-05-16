package services

import (
	"errors"
	"strings"
)

func ParseLabels(labelsStr string) (map[string]string, error) {
	result := make(map[string]string)

	labelsStr = strings.TrimSpace(labelsStr)

	if !strings.HasPrefix(labelsStr, "(") || !strings.HasSuffix(labelsStr, ")") {
		return nil, errors.New("nevalidan format labela, fale zagrade")
	}

	content := labelsStr[1 : len(labelsStr)-1]

	if content == "" {
		return result, nil
	}

	pairs := strings.Split(content, ";")

	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		kv := strings.Split(pair, ":")
		if len(kv) != 2 {
			return nil, errors.New("nevalidan par labela, mora biti u formatu kljuc:vrednost")
		}

		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])

		if key == "" || value == "" {
			return nil, errors.New("kljuc ili vrednost ne mogu biti prazni")
		}

		result[key] = value
	}

	return result, nil
}
