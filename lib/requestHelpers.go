package lib

import "regexp"

// NamedRegexp extends a Regexp object to add a method for returning named capture groups.
type NamedRegexp struct {
	*regexp.Regexp
}

// FindStringSubmatchMap returns a map of named capture groups in a regular expression.
func (r *NamedRegexp) FindStringSubmatchMap(s string) map[string]string {
	captures := make(map[string]string)

	match := r.FindStringSubmatch(s)
	if match == nil {
		return captures
	}

	for i, name := range r.SubexpNames() {
		// Ignore the whole regexp match and unnamed groups
		if i == 0 || name == "" {
			continue
		}

		captures[name] = match[i]

	}
	return captures
}
