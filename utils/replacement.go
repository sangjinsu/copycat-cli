package utils

import (
	"fmt"
	"strings"
)

func ParseReplacements(arg string) *strings.Replacer {
	fmt.Println(arg)

	replacer := strings.NewReplacer(
		"{", "",
		"}", "",
		"[", "",
		"]", "",
		"(", "",
		")", "")
	inputs := replacer.Replace(arg)

	var replacements []string
	pairs := strings.Split(inputs, ",")
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		parts := strings.SplitN(pair, "=", 2)

		if len(parts) != 2 {
			continue
		}

		if len(parts) == 2 {
			oldStr := strings.TrimSpace(parts[0])
			newStr := strings.TrimSpace(parts[1])
			replacements = append(replacements, oldStr, newStr)
		}
	}

	fmt.Println(replacements)
	return strings.NewReplacer(replacements...)
}
