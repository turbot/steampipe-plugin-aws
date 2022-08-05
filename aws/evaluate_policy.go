package aws

import (
	"fmt"
	"regexp"
	"sort"
)

type EvaluatedPrincipal struct {
	allowedPrincipalFederatedIdentitiesSet map[string]bool
	allowedPrincipalServicesSet            map[string]bool
	allowedPrincipalsSet                   map[string]bool
	allowedPrincipalAccountIdsSet          map[string]bool
	isPublic                               bool
}

type EvaluatedStatement struct {
	allowedPrincipalFederatedIdentitiesSet map[string]bool
	allowedPrincipalServicesSet            map[string]bool
	allowedPrincipalsSet                   map[string]bool
	allowedPrincipalAccountIdsSet          map[string]bool
	isPublic                               bool
}

type EvaluatedPolicy struct {
	AccessLevel                         string   `json:"access_level"`
	AllowedOrganizationIds              []string `json:"allowed_organization_ids"`
	AllowedPrincipals                   []string `json:"allowed_principals"`
	AllowedPrincipalAccountIds          []string `json:"allowed_principal_account_ids"`
	AllowedPrincipalFederatedIdentities []string `json:"allowed_principal_federated_identities"`
	AllowedPrincipalServices            []string `json:"allowed_principal_services"`
	IsPublic                            bool     `json:"is_public"`
	PublicAccessLevels                  []string `json:"public_access_levels"`
	PublicStatementIds                  []string `json:"public_statement_ids"`
}

func EvaluatePolicy(policyContent string, userAccountId string) (EvaluatedPolicy, error) {
	evaluatedPolicy := EvaluatedPolicy{}

	// Check source account id which should be valid.
	re := regexp.MustCompile(`^[0-9]{12}$`)

	if !re.MatchString(userAccountId) {
		return evaluatedPolicy, fmt.Errorf("source account id is invalid: %s", userAccountId)
	}

	if policyContent == "" {
		return evaluatedPolicy, nil
	}

	policyInterface, err := canonicalPolicy(policyContent)
	if err != nil {
		return evaluatedPolicy, err
	}

	policy := policyInterface.(Policy)

	evaluatedStatement, err := evaluateStatements(policy.Statements, userAccountId)
	if err != nil {
		return evaluatedPolicy, err
	}

	evaluatedPolicy.AllowedPrincipalFederatedIdentities = setToSortedSlice(evaluatedStatement.allowedPrincipalFederatedIdentitiesSet)
	evaluatedPolicy.AllowedPrincipalServices = setToSortedSlice(evaluatedStatement.allowedPrincipalServicesSet)
	evaluatedPolicy.AllowedPrincipals = setToSortedSlice(evaluatedStatement.allowedPrincipalsSet)
	evaluatedPolicy.AllowedPrincipalAccountIds = setToSortedSlice(evaluatedStatement.allowedPrincipalAccountIdsSet)
	evaluatedPolicy.IsPublic = evaluatedStatement.isPublic

	return evaluatedPolicy, nil
}

func evaluateStatements(statements []Statement, userAccountId string) (EvaluatedStatement, error) {
	evaluatedStatement := EvaluatedStatement{}
	for _, statement := range statements {
		if !checkEffectValid(statement.Effect) {
			return evaluatedStatement, fmt.Errorf("element Effect is invalid - valid choices are 'Allow' or 'Deny'")
		}

		// TODO: For phase 1 - we are only interested in allow else continue with next
		if statement.Effect == "Deny" {
			continue
		}

		// Check principal
		evaluatedPrinciple, err := evaluatePrincipal(statement.Principal, userAccountId)
		if err != nil {
			return evaluatedStatement, err
		}

		evaluatedStatement.allowedPrincipalFederatedIdentitiesSet = mergeSet(
			evaluatedStatement.allowedPrincipalFederatedIdentitiesSet,
			evaluatedPrinciple.allowedPrincipalFederatedIdentitiesSet,
		)

		evaluatedStatement.allowedPrincipalServicesSet = mergeSet(
			evaluatedStatement.allowedPrincipalServicesSet,
			evaluatedPrinciple.allowedPrincipalServicesSet,
		)

		evaluatedStatement.allowedPrincipalsSet = mergeSet(
			evaluatedStatement.allowedPrincipalsSet,
			evaluatedPrinciple.allowedPrincipalsSet,
		)

		evaluatedStatement.allowedPrincipalAccountIdsSet = mergeSet(
			evaluatedStatement.allowedPrincipalAccountIdsSet,
			evaluatedPrinciple.allowedPrincipalAccountIdsSet,
		)

		evaluatedStatement.isPublic = evaluatedPrinciple.isPublic
	}

	return evaluatedStatement, nil
}

func evaluatePrincipal(principal Principal, userAccountId string) (EvaluatedPrincipal, error) {
	evaluatedPrinciple := EvaluatedPrincipal{
		allowedPrincipalFederatedIdentitiesSet: map[string]bool{},
		allowedPrincipalServicesSet:            map[string]bool{},
		allowedPrincipalsSet:                   map[string]bool{},
		allowedPrincipalAccountIdsSet:          map[string]bool{},
	}

	for principalKey, rawPrincipalItem := range principal {
		principalItems := rawPrincipalItem.([]string)

		reIsAwsAccount := regexp.MustCompile(`^[0-9]{12}$`)
		reIsAwsResource := regexp.MustCompile(`^arn:[a-z]*:[a-z]*:[a-z]*:([0-9]{12}):.*$`)

		for _, principalItem := range principalItems {
			switch principalKey {
			case "AWS":

				var account string

				if reIsAwsAccount.MatchString(principalItem) {
					account = principalItem
				} else if reIsAwsResource.MatchString(principalItem) {
					arnAccount := reIsAwsResource.FindStringSubmatch(principalItem)
					account = arnAccount[1]
				} else if principalItem == "*" {
					evaluatedPrinciple.isPublic = true
					account = principalItem
				} else {
					return evaluatedPrinciple, fmt.Errorf("unabled to parse arn: %s", principalItem)
				}

				if userAccountId != account {
					evaluatedPrinciple.allowedPrincipalAccountIdsSet[account] = true
				}

				evaluatedPrinciple.allowedPrincipalsSet[principalItem] = true
			case "Service":
				evaluatedPrinciple.allowedPrincipalServicesSet[principalItem] = true
			case "Federated":
				evaluatedPrinciple.allowedPrincipalFederatedIdentitiesSet[principalItem] = true
			}
		}
	}

	if len(evaluatedPrinciple.allowedPrincipalServicesSet) > 0 {
		evaluatedPrinciple.isPublic = true
	}

	return evaluatedPrinciple, nil
}

func checkEffectValid(effect string) bool {
	if effect == "Deny" || effect == "Allow" {
		return true
	}

	return false
}

func mergeSet(set1 map[string]bool, set2 map[string]bool) map[string]bool {
	if set1 == nil {
		return set2
	}
	if set2 == nil {
		return set1
	}

	for key, value := range set2 {
		set1[key] = value
	}

	return set1
}

func setToSortedSlice(set map[string]bool) []string {
	slice := make([]string, 0, len(set))
	for index := range set {
		slice = append(slice, index)
	}

	sort.Strings(slice)

	return slice
}
