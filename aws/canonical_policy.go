package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//
// https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_grammar.html#policies-grammar-bnf
//

// Policy represents an IAM Policy document
// It would be nice if we could sort the fields (json keys) but postgres jsonb
// "does not preserve the order of object keys",
// per https://www.postgresql.org/docs/9.4/datatype-json.html
type Policy struct {
	Id         string     `json:"Id,omitempty"` // Optional, case sensitive
	Statements Statements `json:"Statement"`    // Required, array of Statements or single statement
	Version    string     `json:"Version"`      // Required, version date string, 2012-10-17 for latest policy model, previous version was 2008-10-17
}

// Statement represents a Statement in an IAM Policy.
// It would be nice if we could sort the fields (json keys) but postgres jsonb
// "does not preserve the order of object keys",
// per https://www.postgresql.org/docs/9.4/datatype-json.html
type Statement struct {
	Action       Value                  `json:"Action,omitempty"`       // Required if NotAction element is missing, string or array of strings, case insensitive <- TODO: Add check
	Condition    map[string]interface{} `json:"Condition,omitempty"`    // Optional, map of conditions
	Effect       string                 `json:"Effect"`                 // Required, Allow or Deny, case sensitive <- TODO: Add check
	NotAction    Value                  `json:"NotAction,omitempty"`    // Required if Action element is missing, string or array of strings, case insensitive  <- TODO: Add check
	NotPrincipal Principal              `json:"NotPrincipal,omitempty"` // Optional for Identity Policies, Required for Resource Policies if NotResource element is missing, string (*) or map of strings/arrays
	NotResource  CaseSensitiveValue     `json:"NotResource,omitempty"`  // Required if Resource element is missing, string or array of strings, case sensitive  <- TODO: Add check
	Principal    Principal              `json:"Principal,omitempty"`    // Optional for Identity Policies, Required for Resource Policies if NotResource element is missing, string (*) or map of strings/arrays
	Resource     CaseSensitiveValue     `json:"Resource,omitempty"`     // Required if NotResource element is missing, string or array of strings, case sensitive  <- TODO: Add check
	Sid          string                 `json:"Sid,omitempty"`          // Optional, case sensitive
}

// tempStatement is used unmarshall to this struct, then copy to Statement to change string case
type tempStatement struct {
	Action       Value                  `json:"Action,omitempty"`       // Optional, string or array of strings, case insensitive
	Condition    map[string]interface{} `json:"Condition,omitempty"`    // Optional, map of conditions
	Effect       string                 `json:"Effect"`                 // Required, Allow or Deny, case sensitive
	NotAction    Value                  `json:"NotAction,omitempty"`    // Optional, string or array of strings, case insensitive
	NotPrincipal Principal              `json:"NotPrincipal,omitempty"` // Optional, string (*) or map of strings/arrays
	NotResource  CaseSensitiveValue     `json:"NotResource,omitempty"`  // Optional, string or array of strings, case sensitive
	Principal    Principal              `json:"Principal,omitempty"`    // Optional, string (*) or map of strings/arrays
	Resource     CaseSensitiveValue     `json:"Resource,omitempty"`     // Optional, string or array of strings, case sensitive
	Sid          string                 `json:"Sid,omitempty"`          // Optional, case sensitive
}

// Statements is an array of statements from an IAM policy
type Statements []Statement

// UnmarshalJSON for the Policy struct.  A policy can contain a single Statement or an
// array of statements, we always convert to array.  Currently, we do not sort these
// but we probably should....
func (statement *Statements) UnmarshalJSON(b []byte) error {
	var raw interface{}
	err := json.Unmarshal(b, &raw)
	if err != nil {
		return fmt.Errorf("UnmarshalJSON failed for Statements (raw): %s", url.QueryEscape(string(b)))
	}

	newStatements := make([]Statement, 0)

	switch raw.(type) {
	// Single Statement case
	case map[string]interface{}:
		var stmt Statement
		if err := json.Unmarshal(b, &stmt); err != nil {
			return fmt.Errorf("UnmarshalJSON failed for Statements (Single Statement): %s", url.QueryEscape(string(b)))
		}
		newStatements = append(newStatements, stmt)
		*statement = newStatements
	// Array of Statements case
	case []interface{}:
		var stmts []Statement
		if err := json.Unmarshal(b, &stmts); err != nil {
			return fmt.Errorf("UnmarshalJSON failed for Statements (Array of Statement): %s", url.QueryEscape(string(b)))
		}
		*statement = stmts

	default:
		return fmt.Errorf("invalid %s value element: allowed is only string or map[]interface{}", reflect.TypeOf(raw))
	}

	return nil
}

// UnmarshalJSON for the Statement struct
func (statement *Statement) UnmarshalJSON(b []byte) error {
	var newStatement tempStatement

	if err := json.Unmarshal(b, &newStatement); err != nil {
		return err
	}

	statement.Sid = newStatement.Sid
	statement.Effect = newStatement.Effect
	statement.Principal = newStatement.Principal
	statement.NotPrincipal = newStatement.NotPrincipal
	statement.Action = newStatement.Action
	statement.NotAction = newStatement.NotAction
	statement.Resource = newStatement.Resource
	statement.NotResource = newStatement.NotResource

	c, err := canonicalCondition(newStatement.Condition)
	if err != nil {
		return fmt.Errorf("error unmarshalling / converting condition: %s", err)
	}
	statement.Condition = c

	return nil
}

// canonicalCondition converts the conditions to a standard format for easier matching
// Note that:
//   - conditions keys are CASE INSENSITIVE - we convert them to lower case.
//   - Like other fields in IAM policies, the condition values can either be a string
//     or an array of strings - we always convery them to arrays for easier searching
//     and we remove duplicates
//   - condition values can be string, boolean, or numeric depending on the operator
//     key,  but whereever the a bool or int is accepted, a string representation is
//     also accepted - e.g. you can use `true` or `"true"`.  While it would probably
//     be ideal to cast to the ACTUAL type based on the operator, we currently cast
//     them all to strings - Its simpler, and the net effect is pretty much the same;
//     since postgres json functions only return text or jsonb, you need to cast
//     them explicitly in your query anyway....
func canonicalCondition(src map[string]interface{}) (map[string]interface{}, error) {
	newConditions := make(map[string]interface{})

	for operator, condition := range src {
		newCondition := make(map[string]interface{})

		for conditionKey, conditionValue := range condition.(map[string]interface{}) {
			// convert the condition key to lower case
			newKey := strings.ToLower(conditionKey)

			// convert the value to a slice of string....)
			newSlice, err := toSliceOfStrings(conditionValue)
			if err != nil {
				return nil, err
			}

			newCondition[newKey] = newSlice
		}

		newConditions[operator] = newCondition
	}

	return newConditions, nil
}

// Principal may be string '*' or a map of principaltype:value.  If '*', we add as an
// array element to the AWS principal type.
// Each value in the map may be a string or []string, we convert everything to []string
// and sort it and remove duplicates
// NOTE: We have the tests in place and we could change this type
type Principal map[string]interface{}

// UnmarshalJSON for the Principal struct
func (principal *Principal) UnmarshalJSON(b []byte) error {
	var raw interface{}

	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	switch principalValue := raw.(type) {
	case string:
		p := make(map[string]interface{})
		p["AWS"] = []string{principalValue}
		*principal = p

	case map[string]interface{}:
		// convert each sub item to array of string
		p := make(map[string]interface{})
		for k, v := range principalValue {
			newSlice, err := toSliceOfStrings(v)
			if err != nil {
				return nil
			}

			p[k] = newSlice
		}

		*principal = p

	default:
		return fmt.Errorf("only string or map[]interface{} values are allowed for the Principal type")
	}

	return nil
}

// Value is an AWS IAM value string or array.  AWS allows string or []string as value,
// we convert everything to []string to avoid casting.  We also sort these - order does
// not matter for arrays/lists in IAM policies, so we sort them for easier diffing,
// and remove duplicates since they're ignored anyway
type Value []string

// UnmarshalJSON for the Value struct
func (result *Value) UnmarshalJSON(data []byte) error {
	var raw interface{}

	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	// convert the value to an array of strings
	newSlice, err := toSliceOfLowerStrings(raw)
	if err != nil {
		return err
	}

	*result = newSlice

	return nil
}

// CaseSensitiveValue is used for value arrays that care about case
// AWS allows string or []string as value, we convert everything to []string to
// avoid casting. We also sort these - order does not matter for arrays/lists
// in IAM policies, so we sort them for easier diffing and remove duplicates
// since they're ignored anyway
type CaseSensitiveValue []string

// UnmarshalJSON for the CaseSensitiveValue struct
func (result *CaseSensitiveValue) UnmarshalJSON(data []byte) error {
	var raw interface{}

	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	// convert the value to an array of strings
	newSlice, err := toSliceOfStrings(raw)
	if err != nil {
		return err
	}

	*result = newSlice
	return nil
}

// canonicalPolicy converts a (unescaped) policy string to canonical format
func canonicalPolicy(src string) (interface{}, error) {
	var policy Policy

	if err := json.Unmarshal([]byte(src), &policy); err != nil {
		return nil, fmt.Errorf("canonicalPolicy failed unmarshalling source data: %+v.  src: %s", err, url.QueryEscape(src))
	}

	return policy, nil
}

//// UTILITY FUNCTIONS

// toSliceOfStrings converts a string or array value to an array of strings
func toSliceOfStrings(scalarOrSlice interface{}) ([]string, error) {
	slice := []string{}
	set := map[string]bool{}

	if reflect.TypeOf(scalarOrSlice).Kind() == reflect.Slice {
		switch item := scalarOrSlice.(type) {
		case []string:
			for _, value := range item {
				if _, exists := set[value]; !exists {
					slice = append(slice, value)
					set[value] = true
				}
			}
		default:
			for _, rawValue := range scalarOrSlice.([]interface{}) {
				value := types.ToString(rawValue)
				if _, exists := set[value]; !exists {
					slice = append(slice, value)
					set[value] = true
				}
			}
		}

		sort.Strings(slice)
	} else {
		slice = append(slice, types.ToString(scalarOrSlice))
	}

	return slice, nil
}

// toSliceOfLowerStrings converts a string or array value to an array of strings and lowercase's the values
func toSliceOfLowerStrings(scalarOrSlice interface{}) ([]string, error) {
	slice := []string{}
	set := map[string]bool{}

	if reflect.TypeOf(scalarOrSlice).Kind() == reflect.Slice {
		switch item := scalarOrSlice.(type) {
		case []string:
			for _, value := range item {
				value = strings.ToLower(value)
				if _, exists := set[value]; !exists {
					slice = append(slice, value)
					set[value] = true
				}
			}
		default:
			for _, rawValue := range scalarOrSlice.([]interface{}) {
				value := types.ToString(rawValue)
				value = strings.ToLower(value)
				if _, exists := set[value]; !exists {
					slice = append(slice, value)
					set[value] = true
				}
			}
		}

		sort.Strings(slice)
	} else {
		value := types.ToString(scalarOrSlice)
		value = strings.ToLower(value)
		slice = append(slice, value)
	}

	return slice, nil
}

// uniqueStrings removes duplicate items from a slice of strings
func uniqueStrings(arr []string) []string {
	occured := map[string]bool{}
	result := []string{}
	for e := range arr {
		// check if already the mapped (if true)
		if !occured[arr[e]] {
			occured[arr[e]] = true

			// Append to result slice.
			result = append(result, arr[e])
		}
	}

	return result
}

//// TRANSFORM FUNCTIONS

// unescape a string.  Often (but not always), a policy doc is an escaped string,
// and it must be unescaped beofre converting to canonical form
func unescape(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("unescape")

	// get the value of policy safely
	inputStr := types.SafeString(d.Value)

	data, err := url.QueryUnescape(inputStr)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// policyToCanonical converts a (unescaped) IAM policy to a standardized form
func policyToCanonical(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("policyStringToCanonical")

	data := types.SafeString(d.Value)
	if data == "" {
		return nil, nil
	}

	newPolicy, err := canonicalPolicy(data)
	if err != nil {
		logger.Error("policyStringToCanonical", "err", err)
		return nil, err
	}

	return newPolicy, nil
}

// Inline policies in canonical form
func inlinePoliciesToStd(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	inlinePolicies := d.HydrateItem.([]map[string]interface{})

	var inlinePoliciesStd []map[string]interface{}
	if inlinePolicies == nil {
		return nil, nil
	}

	for _, inlinePolicy := range inlinePolicies {
		strPolicy, err := json.Marshal(inlinePolicy["PolicyDocument"])
		if err != nil {
			plugin.Logger(ctx).Error("inlinePoliciesToStd", fmt.Sprintf("transform_error for %s", d.ColumnName), err)
			return nil, err
		}
		policyStd, errStd := canonicalPolicy(string(strPolicy))
		if errStd != nil {
			return nil, errStd
		}

		inlinePoliciesStd = append(inlinePoliciesStd, map[string]interface{}{
			"PolicyDocument": policyStd,
			"PolicyName":     inlinePolicy["PolicyName"],
		})
	}

	return inlinePoliciesStd, nil
}
