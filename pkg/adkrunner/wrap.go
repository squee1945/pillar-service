package adkrunner

import "strings"

func wrap(input string, limit int, prefix string) []string {
	words := strings.Fields(input)
	if len(words) == 0 {
		return nil
	}

	contentLimit := limit - len(prefix)
	if contentLimit <= 0 {
		contentLimit = 1
	}

	var result []string
	currentLine := strings.Builder{}

	currentLine.WriteString(prefix)

	for _, word := range words {
		currentContentLength := currentLine.Len() - len(prefix)
		potentialNewLength := currentContentLength + 1 + len(word)
		isNewLine := currentContentLength == 0

		if isNewLine {
			if len(word) <= contentLimit {
				currentLine.WriteString(word)
			} else {
				currentLine.WriteString(word)
			}
			continue
		}

		if potentialNewLength <= contentLimit {
			currentLine.WriteString(" ")
			currentLine.WriteString(word)
			continue
		}

		result = append(result, currentLine.String())

		currentLine.Reset()
		currentLine.WriteString(prefix)
		currentLine.WriteString(word)
	}

	result = append(result, currentLine.String())
	return result
}
