package service

import (
	"errors"
	"strings"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

func AutoConvert(input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", errors.New("пустой ввод")
	}

	if strings.ContainsAny(input, ".-/") {
		return morse.ToText(input), nil
	}

	return morse.ToMorse(input), nil
}
