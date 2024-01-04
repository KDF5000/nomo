package utils

import (
	"unicode"
)

type ContentElement struct {
	Text  string
	IsTag bool
}

func ScanContent(content string) []ContentElement {
	var elements []ContentElement
	var lastTagEnd = 0
	var tagStart = -1
	for i, v := range content {
		// found tag end
		if tagStart != -1 && unicode.IsSpace(v) {
			if i == tagStart+1 {
				// no tag only #
				lastTagEnd = tagStart
				tagStart = -1
				continue
			}

			elements = append(elements, ContentElement{
				Text:  content[tagStart:i],
				IsTag: true,
			})

			tagStart = -1
			lastTagEnd = i
			continue
		}

		if v != '#' {
			i++
			continue
		}

		tagStart = i
		if i == 0 {
			continue
		}

		elements = append(elements, ContentElement{
			Text:  content[lastTagEnd:i],
			IsTag: false,
		})
	}

	if tagStart != -1 {
		elements = append(elements, ContentElement{
			Text:  content[tagStart:],
			IsTag: true,
		})
	} else if lastTagEnd < len(content) {
		elements = append(elements, ContentElement{
			Text:  content[lastTagEnd:],
			IsTag: false,
		})
	}

	return elements
}

func RetriveTags(content string) []string {
	var tags []string
	for i := 0; i < len(content); {
		if content[i] != '#' {
			i++
			continue
		}

		// find tag end
		var end = len(content)
		start := i + 1
		for start := i + 1; start < len(content); start++ {
			if content[start] == ' ' {
				end = start
				break
			}
		}

		if end != start {
			tags = append(tags, content[i+1:end])
		}
		i = end
	}

	return tags
}

func Truncate(content string, count int) string {
	s := []rune(content)
	if count > len(s) {
		count = len(s)
	}
	return string(s[:count])
}
