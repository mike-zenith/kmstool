package kmstool

import (
	"errors"
	"regexp"
)

// embed regexp.Regexp in a new type so we can extend it
type NamedGroupRegexp struct {
	*regexp.Regexp
}

func (r *NamedGroupRegexp) GetNamedGroupsFromMatch(match []string) map[string]string {
	captures := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}
		captures[name] = match[i]

	}
	return captures
}

func (r *NamedGroupRegexp) FindAllByteSubmatchMap(c string, n int) map[int]map[string]string {
	captures := make(map[int]map[string]string)

	match := r.FindAllStringSubmatch(c, n)
	if match == nil {
		return captures
	}

	for i, foundMap := range match {
		captures[i] = r.GetNamedGroupsFromMatch(foundMap)
	}

	return captures
}

func NewRegexpStringParser(r string) func(string) (map[string]string, error) {
	return func(c string) (map[string]string, error) {
		re := NamedGroupRegexp{regexp.MustCompile(r)}
		foundSubsets := re.FindAllByteSubmatchMap(c, -1)

		result := make(map[string]string)
		for _, foundSubset := range foundSubsets {
			if len(foundSubset["hit"]) == 0 || len(foundSubset["key"]) == 0 {
				return nil, errors.New("Subset does not contain 'hit' and 'key'")
			}
			result[string(foundSubset["hit"])] = foundSubset["key"]
		}

		return result, nil
	}
}
