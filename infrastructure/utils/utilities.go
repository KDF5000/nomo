package utils

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
