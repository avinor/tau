package strings

import "regexp"

var (
	outputRegexp = regexp.MustCompile("(?m:^\\s*\"?([^\"=\\s]*)\"?\\s*=\\s*\"?([^\"\\n]*)\"?$)")
)

// ParseVars parses each line as key=value and returns a map of all variables.
// If a line cannot be parsed it will be ignored.
func ParseVars(output string) map[string]string {
	matches := outputRegexp.FindAllStringSubmatch(output, -1)
	values := map[string]string{}

	if len(matches) == 0 {
		return values
	}

	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		values[match[1]] = match[2]
	}

	return values
}
