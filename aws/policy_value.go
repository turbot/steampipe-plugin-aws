package aws

import (
	"fmt"
	"regexp"
	"strings"
)

type PolicyValue struct {
	value string
	regex string
}

func MakePolicyValue(value string) PolicyValue {
	regex := strings.Replace(value, "*", ".*", len(value))
	regex = strings.Replace(regex, "?", "[^*]{1}", len(regex))
	regex = fmt.Sprintf("^%s$", regex)

	return PolicyValue{value, regex}
}

func (policyValue PolicyValue) Contains(other PolicyValue) bool {
	re := regexp.MustCompile(policyValue.regex)
	return re.MatchString(other.value)
}

func (policyValue PolicyValue) Intersection(other PolicyValue) string {
	if policyValue.Contains(other) {
		return other.value
	}

	if other.Contains(policyValue) {
		return policyValue.value
	}

	return ""
}
