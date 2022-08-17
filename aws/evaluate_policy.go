package aws

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"
)

type EvaluatedAction struct {
	process    bool
	prefix     string
	priviledge string
	matcher    string
}

type CompletedEvaluation struct {
	allowedPrincipalFederatedIdentitiesSet map[string]bool
	allowedPrincipalServicesSet            map[string]bool
	allowedPrincipalsSet                   map[string]bool
	allowedPrincipalAccountIdsSet          map[string]bool
	allowedOrganizationIds                 map[string]bool
	isPublic                               bool
	isShared                               bool
	isPrivate                              bool
}

type EvaluatedStatements struct {
	allowedPrincipalFederatedIdentitiesSet map[string]bool
	allowedPrincipalServicesSet            map[string]bool
	allowedPrincipalsSet                   map[string]bool
	allowedPrincipalAccountIdsSet          map[string]bool
	allowedOrganizationIds                 map[string]bool
	publicStatementIds                     map[string]bool
	sharedStatementIds                     map[string]bool
	publicAccessLevels                     []string
	sharedAccessLevels                     []string
	privateAccessLevels                    []string
	isPublic                               bool
	isShared                               bool
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
	SharedAccessLevels                  []string `json:"shared_access_levels"`
	PrivateAccessLevels                 []string `json:"private_access_levels"`
	PublicStatementIds                  []string `json:"public_statement_ids"`
	SharedStatementIds                  []string `json:"shared_statement_ids"`
}

type Permissions struct {
	privileges  []string
	accessLevel map[string]string
}

func EvaluatePolicy(policyContent string, userAccountId string) (EvaluatedPolicy, error) {
	evaluatedPolicy := EvaluatedPolicy{
		AccessLevel: "private",
	}

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

	permissions := getSortedPermissions()

	policy := policyInterface.(Policy)

	evaluatedStatements, err := evaluateStatements(policy.Statements, userAccountId, permissions)
	if err != nil {
		return evaluatedPolicy, err
	}

	evaluatedPolicy.AccessLevel = evaluateAccessLevel(evaluatedStatements)
	evaluatedPolicy.AllowedPrincipalFederatedIdentities = setToSortedSlice(evaluatedStatements.allowedPrincipalFederatedIdentitiesSet)
	evaluatedPolicy.AllowedPrincipalServices = setToSortedSlice(evaluatedStatements.allowedPrincipalServicesSet)
	evaluatedPolicy.AllowedPrincipals = setToSortedSlice(evaluatedStatements.allowedPrincipalsSet)
	evaluatedPolicy.AllowedPrincipalAccountIds = setToSortedSlice(evaluatedStatements.allowedPrincipalAccountIdsSet)
	evaluatedPolicy.AllowedOrganizationIds = setToSortedSlice(evaluatedStatements.allowedOrganizationIds)
	evaluatedPolicy.PublicStatementIds = setToSortedSlice(evaluatedStatements.publicStatementIds)
	evaluatedPolicy.SharedStatementIds = setToSortedSlice(evaluatedStatements.sharedStatementIds)
	evaluatedPolicy.PublicAccessLevels = evaluatedStatements.publicAccessLevels
	evaluatedPolicy.SharedAccessLevels = evaluatedStatements.sharedAccessLevels
	evaluatedPolicy.PrivateAccessLevels = evaluatedStatements.privateAccessLevels
	evaluatedPolicy.IsPublic = evaluatedStatements.isPublic

	return evaluatedPolicy, nil
}

func evaluateAccessLevel(statements EvaluatedStatements) string {
	if statements.isPublic {
		return "public"
	}

	if statements.isShared {
		return "shared"
	}

	return "private"
}

type EvaluateStatements struct {
	statements             EvaluatedStatements
	publicAccessLevelsSet  map[string]bool
	sharedAccessLevelsSet  map[string]bool
	privateAccessLevelsSet map[string]bool
}

func evaluateStatements(statements []Statement, userAccountId string, permissions map[string]Permissions) (EvaluatedStatements, error) {
	var evaluatedStatement EvaluatedStatements

	allowedStatements := EvaluateStatements{
		statements: EvaluatedStatements{
			publicStatementIds: map[string]bool{},
			sharedStatementIds: map[string]bool{},
		},
		publicAccessLevelsSet:  map[string]bool{},
		sharedAccessLevelsSet:  map[string]bool{},
		privateAccessLevelsSet: map[string]bool{},
	}
	deniedStatements := EvaluateStatements{
		statements: EvaluatedStatements{
			publicStatementIds: map[string]bool{},
			sharedStatementIds: map[string]bool{},
		},
		publicAccessLevelsSet:  map[string]bool{},
		sharedAccessLevelsSet:  map[string]bool{},
		privateAccessLevelsSet: map[string]bool{},
	}

	uniqueStatementIds := map[string]bool{}

	var currentStatements *EvaluateStatements

	for statementIndex, statement := range statements {
		if !checkEffectValid(statement.Effect) {
			return evaluatedStatement, fmt.Errorf("element Effect is invalid - valid choices are 'Allow' or 'Deny'")
		}

		// TODO: For phase 1 - we are only interested in allow else continue with next
		if statement.Effect == "Deny" {
			currentStatements = &deniedStatements
		} else {
			currentStatements = &allowedStatements
		}

		// Principal
		evaluatedPrinciple, err := evaluatePrincipal(statement.Principal, userAccountId)
		if err != nil {
			return evaluatedStatement, err
		}

		evaluatedCondition, err := evaluateCondition(statement.Condition, userAccountId)
		if err != nil {
			return evaluatedStatement, err
		}

		currentStatements.statements.allowedPrincipalFederatedIdentitiesSet = mergeSets(
			currentStatements.statements.allowedPrincipalFederatedIdentitiesSet,
			evaluatedPrinciple.allowedPrincipalFederatedIdentitiesSet,
			evaluatedCondition.allowedPrincipalFederatedIdentitiesSet,
		)

		currentStatements.statements.allowedPrincipalServicesSet = mergeSets(
			currentStatements.statements.allowedPrincipalServicesSet,
			evaluatedPrinciple.allowedPrincipalServicesSet,
			evaluatedCondition.allowedPrincipalServicesSet,
		)

		currentStatements.statements.allowedPrincipalsSet = mergeSets(
			currentStatements.statements.allowedPrincipalsSet,
			evaluatedPrinciple.allowedPrincipalsSet,
			evaluatedCondition.allowedPrincipalsSet,
		)

		currentStatements.statements.allowedPrincipalAccountIdsSet = mergeSets(
			currentStatements.statements.allowedPrincipalAccountIdsSet,
			evaluatedPrinciple.allowedPrincipalAccountIdsSet,
			evaluatedCondition.allowedPrincipalAccountIdsSet,
		)

		currentStatements.statements.allowedOrganizationIds = mergeSets(
			currentStatements.statements.allowedOrganizationIds,
			evaluatedPrinciple.allowedOrganizationIds,
			evaluatedCondition.allowedOrganizationIds,
		)

		// Visibility
		isStatementPublic := evaluatedPrinciple.isPublic || evaluatedCondition.isPublic
		isStatementShared := evaluatedPrinciple.isShared || evaluatedCondition.isShared
		isStatementPrivate := evaluatedPrinciple.isPrivate || evaluatedCondition.isPrivate

		// Before using Sid, let's check to see if it is unique
		sid := evaluatedSid(statement, statementIndex)
		if _, exists := uniqueStatementIds[sid]; exists {
			return evaluatedStatement, fmt.Errorf("duplicate Sid found: %s", sid)
		}
		uniqueStatementIds[sid] = true

		if isStatementPublic {
			currentStatements.statements.isPublic = true
			currentStatements.statements.publicStatementIds[sid] = true
			for _, action := range statement.Action {
				if _, exists := currentStatements.publicAccessLevelsSet[action]; !exists {
					currentStatements.publicAccessLevelsSet[action] = true
				}
			}
		}

		if isStatementShared {
			currentStatements.statements.isShared = true
			currentStatements.statements.sharedStatementIds[sid] = true
			for _, action := range statement.Action {
				if _, exists := currentStatements.sharedAccessLevelsSet[action]; !exists {
					currentStatements.sharedAccessLevelsSet[action] = true
				}
			}
		}

		if isStatementPrivate {
			// Actions
			for _, action := range statement.Action {
				if _, exists := currentStatements.privateAccessLevelsSet[action]; !exists {
					currentStatements.privateAccessLevelsSet[action] = true
				}
			}
		}
	}

	evaluatedStatement = evaluateOverallStatements(
		allowedStatements,
		deniedStatements,
		permissions,
	)

	return evaluatedStatement, nil
}

func evaluateOverallStatements(
	allowedStatements EvaluateStatements,
	deniedStatements EvaluateStatements,
	permissions map[string]Permissions,
) EvaluatedStatements {
	overallStatements := EvaluatedStatements{}

	if deniedStatements.statements.isPublic {
		return overallStatements
	}

	overallStatements.allowedPrincipalFederatedIdentitiesSet = allowedStatements.statements.allowedPrincipalFederatedIdentitiesSet
	overallStatements.allowedPrincipalServicesSet = allowedStatements.statements.allowedPrincipalServicesSet
	overallStatements.allowedPrincipalsSet = allowedStatements.statements.allowedPrincipalsSet
	overallStatements.allowedPrincipalAccountIdsSet = allowedStatements.statements.allowedPrincipalAccountIdsSet
	overallStatements.allowedOrganizationIds = allowedStatements.statements.allowedOrganizationIds
	overallStatements.publicStatementIds = allowedStatements.statements.publicStatementIds
	overallStatements.sharedStatementIds = allowedStatements.statements.sharedStatementIds
	overallStatements.publicAccessLevels = evaluateActionSet(allowedStatements.publicAccessLevelsSet, permissions)
	overallStatements.sharedAccessLevels = evaluateActionSet(allowedStatements.sharedAccessLevelsSet, permissions)
	overallStatements.privateAccessLevels = evaluateActionSet(allowedStatements.privateAccessLevelsSet, permissions)
	overallStatements.isPublic = allowedStatements.statements.isPublic
	overallStatements.isShared = allowedStatements.statements.isShared

	return overallStatements
}

func evaluateAction(action string) EvaluatedAction {
	evaluated := EvaluatedAction{}

	lowerAction := strings.ToLower(action)
	actionParts := strings.Split(lowerAction, ":")
	evaluated.prefix = actionParts[0]

	if len(actionParts) < 2 || actionParts[1] == "" {
		return evaluated
	}

	evaluated.process = true

	raw := actionParts[1]

	wildcardLocator := regexp.MustCompile(`[0-9a-z:]*(\*|\?)`)
	located := wildcardLocator.FindString(raw)

	if located == "" {
		evaluated.priviledge = raw
		return evaluated
	}

	evaluated.priviledge = located[:len(located)-1]

	// Convert Wildcards to regexp
	matcher := fmt.Sprintf("^%s$", raw)
	matcher = strings.Replace(matcher, "*", "[a-z0-9]*", len(matcher))
	matcher = strings.Replace(matcher, "?", "[a-z0-9]{1}", len(matcher))

	evaluated.matcher = matcher

	return evaluated
}

func evaluateActionSet(allowedActionSet map[string]bool, permissions map[string]Permissions) []string {
	if _, exists := allowedActionSet["*"]; exists {
		return []string{
			"List",
			"Permissions management",
			"Read",
			"Tagging",
			"Write",
		}
	}

	accessLevels := map[string]bool{}

	for action := range allowedActionSet {
		evaluatedAction := evaluateAction(action)

		if !evaluatedAction.process {
			continue
		}

		// Find service
		if _, exists := permissions[evaluatedAction.prefix]; !exists {
			continue
		}

		permission := permissions[evaluatedAction.prefix]

		// Find API Call
		privilegesLen := len(permission.privileges)
		checkIndex := sort.SearchStrings(permission.privileges, evaluatedAction.priviledge)
		if checkIndex >= privilegesLen {
			continue
		}

		if evaluatedAction.matcher == "" {
			accessLevel := permission.accessLevel[evaluatedAction.priviledge]

			if _, exists := accessLevels[accessLevel]; !exists {
				accessLevels[accessLevel] = true
			}
			continue
		}

		evaluatedPriviledgeLen := len(evaluatedAction.priviledge)
		matcher := regexp.MustCompile(evaluatedAction.matcher)
		for ; checkIndex < privilegesLen; checkIndex++ {
			currentPrivilege := permission.privileges[checkIndex]
			currentPrivilegeLen := len(currentPrivilege)

			splitIndex := int(math.Min(float64(currentPrivilegeLen), float64(evaluatedPriviledgeLen)))
			partialPriviledge := currentPrivilege[0:splitIndex]

			if partialPriviledge != evaluatedAction.priviledge {
				break
			}
			if !matcher.MatchString(currentPrivilege) {
				continue
			}
			accessLevel := permission.accessLevel[currentPrivilege]

			if _, exists := accessLevels[accessLevel]; !exists {
				accessLevels[accessLevel] = true
			}
		}
	}

	return setToSortedSlice(accessLevels)
}

func evaluatedSid(statement Statement, statementIndex int) string {
	if statement.Sid == "" {
		return fmt.Sprintf("Statement[%d]", statementIndex+1)
	}

	return statement.Sid
}

type EvaluatedOperator struct {
	category   string
	isNegated  bool
	isLike     bool
	isCaseless bool
}

func evaulateOperator(operator string) (EvaluatedOperator, bool) {
	// Check if there is an IfExists and then strip it.
	operator = strings.ToLower(operator)
	operator = strings.TrimSuffix(operator, "ifexists")

	evaulatedOperator := EvaluatedOperator{}
	evaluated := true
	switch operator {
	case "stringequals":
		evaulatedOperator.category = "string"
		evaulatedOperator.isNegated = false
		evaulatedOperator.isLike = false
		evaulatedOperator.isCaseless = false
	case "stringnotequals":
		evaulatedOperator.category = "string"
		evaulatedOperator.isNegated = true
		evaulatedOperator.isLike = false
		evaulatedOperator.isCaseless = false
	case "stringequalsignorecase":
		evaulatedOperator.category = "string"
		evaulatedOperator.isNegated = false
		evaulatedOperator.isLike = false
		evaulatedOperator.isCaseless = true
	case "stringnotequalsignorecase":
		evaulatedOperator.category = "string"
		evaulatedOperator.isNegated = true
		evaulatedOperator.isLike = false
		evaulatedOperator.isCaseless = true
	case "stringlike":
		evaulatedOperator.category = "string"
		evaulatedOperator.isNegated = false
		evaulatedOperator.isLike = true
		evaulatedOperator.isCaseless = false
	case "stringnotlike":
		evaulatedOperator.category = "string"
		evaulatedOperator.isNegated = false
		evaulatedOperator.isLike = true
		evaulatedOperator.isCaseless = false
	case "arnequals":
		evaulatedOperator.category = "arn"
		evaulatedOperator.isNegated = false
		evaulatedOperator.isLike = false
		evaulatedOperator.isCaseless = true
	case "arnlike":
		evaulatedOperator.category = "arn"
		evaulatedOperator.isNegated = false
		evaulatedOperator.isLike = true
		evaulatedOperator.isCaseless = true
	case "arnnotequals":
		evaulatedOperator.category = "arn"
		evaulatedOperator.isNegated = true
		evaulatedOperator.isLike = false
		evaulatedOperator.isCaseless = true
	case "arnnotlike":
		evaulatedOperator.category = "arn"
		evaulatedOperator.isNegated = true
		evaulatedOperator.isLike = true
		evaulatedOperator.isCaseless = true
	default:
		evaluated = false
	}

	return evaulatedOperator, evaluated
}

func evaluateArnTypeCondition(conditionValues []string, evaulatedOperator EvaluatedOperator, userAccountId string) CompletedEvaluation {
	completedEvaluation := CompletedEvaluation{
		allowedPrincipalsSet:          map[string]bool{},
		allowedPrincipalAccountIdsSet: map[string]bool{},
	}

	for _, principal := range conditionValues {
		if evaulatedOperator.category != "string" && evaulatedOperator.category != "arn" {
			continue
		}

		if evaulatedOperator.isLike {
			if evaulatedOperator.category == "string" {
				completedEvaluation.allowedPrincipalsSet[principal] = true
				// We need to pull the account out of a wildcard type
				// Assume that account is before any other numeric
				// There should never be any more that 13 digits
				reAccountExtractor := regexp.MustCompile(`^.*[^0-9]+([0-9]{12})[^0-9]+.*$`)
				arnAccount := reAccountExtractor.FindStringSubmatch(principal)
				if len(arnAccount) > 0 {
					account := arnAccount[1]
					if account != userAccountId {
						completedEvaluation.isShared = true
					} else {
						completedEvaluation.isPrivate = true
					}
					completedEvaluation.allowedPrincipalAccountIdsSet[account] = true
				} else {
					completedEvaluation.isPublic = true
					completedEvaluation.allowedPrincipalAccountIdsSet["*"] = true
				}
			} else if evaulatedOperator.category == "arn" {
				splitPrincipal := strings.Split(principal, ":")
				// There should always be an account
				if len(splitPrincipal) < 5 {
					continue
				}

				account := splitPrincipal[4]
				accountLength := len(account)

				if strings.Contains(account, "*") && accountLength <= 12 {
					completedEvaluation.allowedPrincipalsSet[principal] = true
					completedEvaluation.allowedPrincipalAccountIdsSet["*"] = true
					completedEvaluation.isPublic = true
					continue
				}

				if accountLength == 0 || accountLength != 12 {
					continue
				}

				if strings.Contains(account, "?") {
					completedEvaluation.allowedPrincipalsSet[principal] = true
					completedEvaluation.allowedPrincipalAccountIdsSet["*"] = true
					completedEvaluation.isPublic = true
					continue
				}

				re := regexp.MustCompile(`^[0-9]{12}$`)
				if !re.MatchString(account) {
					continue
				}

				completedEvaluation.allowedPrincipalsSet[principal] = true
				completedEvaluation.allowedPrincipalAccountIdsSet[account] = true

				if account != userAccountId {
					completedEvaluation.isShared = true
					continue
				}

				completedEvaluation.isPrivate = true
			}

			continue
		}

		// Check if principal doesn't match an the ARN format, ignore
		reIsAwsResource := regexp.MustCompile(`^arn:[a-z]*:[a-z]*:[a-z]*:([0-9]{12}):.*$`)
		if !reIsAwsResource.MatchString(principal) {
			continue
		}

		arnAccount := reIsAwsResource.FindStringSubmatch(principal)
		account := arnAccount[1]

		// Check if principal doesn't match an account ID, ignore
		reAccount := regexp.MustCompile(`^[0-9]{12}$`)
		if !reAccount.MatchString(account) {
			continue
		}

		completedEvaluation.allowedPrincipalsSet[principal] = true
		completedEvaluation.allowedPrincipalAccountIdsSet[account] = true

		if account == userAccountId {
			completedEvaluation.isPrivate = true
		} else {
			completedEvaluation.isShared = true
		}
	}

	return completedEvaluation
}

func evaluateOrganizationCondition(conditionValues []string, evaulatedOperator EvaluatedOperator, userAccountId string) CompletedEvaluation {
	completedEvaluation := CompletedEvaluation{
		allowedOrganizationIds: map[string]bool{},
	}

	for _, principal := range conditionValues {
		if evaulatedOperator.category != "string" {
			continue
		}

		organization := principal
		if evaulatedOperator.isLike {
			if organization == "*" || organization == "o-*" {
				completedEvaluation.allowedOrganizationIds["o-*"] = true
				completedEvaluation.isPublic = true
				continue
			}

			if !strings.HasPrefix(organization, "o-") {
				continue
			}

			completedEvaluation.allowedOrganizationIds[organization] = true
			completedEvaluation.isShared = true

			continue
		}

		if !strings.HasPrefix(organization, "o-") || strings.Contains(organization, "*") {
			continue
		}

		completedEvaluation.allowedOrganizationIds[organization] = true
		completedEvaluation.isShared = true
	}

	return completedEvaluation
}

func evaluateAccountTypeCondition(conditionValues []string, evaulatedOperator EvaluatedOperator, userAccountId string) CompletedEvaluation {
	completedEvaluation := CompletedEvaluation{
		allowedPrincipalsSet:          map[string]bool{},
		allowedPrincipalAccountIdsSet: map[string]bool{},
	}

	for _, principal := range conditionValues {
		if evaulatedOperator.category != "string" {
			continue
		}

		if evaulatedOperator.isLike {
			account := principal
			accountLength := len(account)

			if strings.Contains(account, "*") && accountLength <= 12 {
				completedEvaluation.allowedPrincipalsSet[principal] = true
				completedEvaluation.allowedPrincipalAccountIdsSet["*"] = true
				completedEvaluation.isPublic = true
				continue
			}

			if accountLength == 0 || accountLength != 12 {
				continue
			}

			if strings.Contains(account, "?") {
				completedEvaluation.allowedPrincipalsSet[principal] = true
				completedEvaluation.allowedPrincipalAccountIdsSet["*"] = true
				completedEvaluation.isPublic = true
				continue
			}

			re := regexp.MustCompile(`^[0-9]{12}$`)
			if !re.MatchString(account) {
				continue
			}

			completedEvaluation.allowedPrincipalsSet[principal] = true
			completedEvaluation.allowedPrincipalAccountIdsSet[account] = true
			if account != userAccountId {
				completedEvaluation.isShared = true
				continue
			}

			completedEvaluation.isPrivate = true
			continue
		}

		// Check if principal doesn't match an account ID, ignore
		re := regexp.MustCompile(`^[0-9]{12}$`)
		if !re.MatchString(principal) {
			continue
		}

		completedEvaluation.allowedPrincipalsSet[principal] = true
		completedEvaluation.allowedPrincipalAccountIdsSet[principal] = true

		if principal == userAccountId {
			completedEvaluation.isPrivate = true
		} else {
			completedEvaluation.isShared = true
		}
	}

	return completedEvaluation
}

func evaluateCondition(conditions map[string]interface{}, userAccountId string) (CompletedEvaluation, error) {
	var completedEvaluation CompletedEvaluation

	for operator, conditionKey := range conditions {
		evaulatedOperator, evaluated := evaulateOperator(operator)
		if !evaluated {
			continue
		}

		if evaulatedOperator.isNegated {
			return completedEvaluation, fmt.Errorf("TODO: Implement")
			// NOTE: Here we have an issue with the table.
			// 		 The problem is that if we say some principal is NOT an account, this means everything but.
			// 		 I do not know how to represent this in the current table design.
		}

		for conditionName, conditionValues := range conditionKey.(map[string]interface{}) {
			switch conditionName {
			case "aws:principalaccount":
				completedEvaluation = evaluateAccountTypeCondition(conditionValues.([]string), evaulatedOperator, userAccountId)
			case "aws:sourceaccount":
				completedEvaluation = evaluateAccountTypeCondition(conditionValues.([]string), evaulatedOperator, userAccountId)
			case "aws:principalarn":
				completedEvaluation = evaluateArnTypeCondition(conditionValues.([]string), evaulatedOperator, userAccountId)
			case "aws:sourcearn":
				completedEvaluation = evaluateArnTypeCondition(conditionValues.([]string), evaulatedOperator, userAccountId)
			case "aws:principalorgid":
				completedEvaluation = evaluateOrganizationCondition(conditionValues.([]string), evaulatedOperator, userAccountId)
			}
		}
	}

	return completedEvaluation, nil
}

func evaluatePrincipal(principal Principal, userAccountId string) (CompletedEvaluation, error) {
	completedPrinciple := CompletedEvaluation{
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

				if principalItem == "*" {
					account = principalItem
					completedPrinciple.isPublic = true
					completedPrinciple.allowedPrincipalAccountIdsSet[account] = true
				} else {
					if reIsAwsAccount.MatchString(principalItem) {
						account = principalItem
					} else if reIsAwsResource.MatchString(principalItem) {
						arnAccount := reIsAwsResource.FindStringSubmatch(principalItem)
						account = arnAccount[1]
					} else {
						return completedPrinciple, fmt.Errorf("unabled to parse arn or account: %s", principalItem)
					}

					if userAccountId != account {
						completedPrinciple.isShared = true
					} else {
						completedPrinciple.isPrivate = true
					}
					completedPrinciple.allowedPrincipalAccountIdsSet[account] = true
				}

				completedPrinciple.allowedPrincipalsSet[principalItem] = true
			case "Service":
				completedPrinciple.allowedPrincipalServicesSet[principalItem] = true
				completedPrinciple.isPublic = true
			case "Federated":
				completedPrinciple.allowedPrincipalFederatedIdentitiesSet[principalItem] = true
				completedPrinciple.isPrivate = true
			}
		}
	}

	return completedPrinciple, nil
}

func checkEffectValid(effect string) bool {
	if effect == "Deny" || effect == "Allow" {
		return true
	}

	return false
}

func mergeSets(dest map[string]bool, source1 map[string]bool, source2 map[string]bool) map[string]bool {
	dest = mergeSet(dest, source1)
	dest = mergeSet(dest, source2)

	return dest
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

//func getSortedPermissions() map[string][]ParliamentPrivilege {
func getSortedPermissions() map[string]Permissions {
	sorted := map[string]Permissions{}
	unsorted := getParliamentIamPermissions()

	for _, parliamentService := range unsorted {
		if _, exist := sorted[parliamentService.Prefix]; !exist {
			privileges := []string{}
			accessLevel := map[string]string{}

			for _, priviledge := range parliamentService.Privileges {
				lowerPriviledge := strings.ToLower(priviledge.Privilege)
				privileges = append(privileges, lowerPriviledge)
				accessLevel[lowerPriviledge] = priviledge.AccessLevel
			}

			sort.Strings(privileges)

			sorted[parliamentService.Prefix] = Permissions{
				privileges:  privileges,
				accessLevel: accessLevel,
			}
		}
	}

	return sorted
}
