package aws

import (
	"fmt"
	"regexp"
	"strings"
)

type PolicyValue struct {
	value  string
	regexp *regexp.Regexp
}

func MakePolicyValue(value string) PolicyValue {
	regex := strings.Replace(value, "*", ".*", len(value))
	regex = strings.Replace(regex, "?", "[^*]{1}", len(regex))
	regex = fmt.Sprintf("^%s$", regex)
	regexp := regexp.MustCompile(regex)

	return PolicyValue{value, regexp}
}

func (policyValue PolicyValue) Contains(other string) bool {
	return policyValue.regexp.MatchString(other)
}

func (policyValue PolicyValue) Intersection(other PolicyValue) string {
	if policyValue.Contains(other.value) {
		return other.value
	}

	if other.Contains(policyValue.value) {
		return policyValue.value
	}

	return ""
}
