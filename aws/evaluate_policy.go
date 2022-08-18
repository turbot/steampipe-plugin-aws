package aws

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"
)

type PolicySummary struct {
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

func EvaluatePolicy(policyContent string, userAccountId string) (PolicySummary, error) {
	policySummary := PolicySummary{
		AccessLevel: "private",
	}

	// Check source account id which should be valid.
	re := regexp.MustCompile(`^[0-9]{12}$`)

	if !re.MatchString(userAccountId) {
		return policySummary, fmt.Errorf("source account id is invalid: %s", userAccountId)
	}

	if policyContent == "" {
		return policySummary, nil
	}

	policyInterface, err := canonicalPolicy(policyContent)
	if err != nil {
		return policySummary, err
	}

	allAvailablePermissions := loadAllAvailablePermissions()

	policy := policyInterface.(Policy)

	allowedEvaluatedStatements, deniedEvaluatedStatements, err := evaluateStatements(policy.Statements, userAccountId, allAvailablePermissions)
	if err != nil {
		return policySummary, err
	}

	statementsSummary := generateStatementsSummary(allowedEvaluatedStatements, deniedEvaluatedStatements, allAvailablePermissions)

	policySummary.AccessLevel = evaluateAccessLevel(statementsSummary)
	policySummary.AllowedPrincipalFederatedIdentities = setToSortedSlice(statementsSummary.allowedPrincipalFederatedIdentitiesSet)
	policySummary.AllowedPrincipalServices = setToSortedSlice(statementsSummary.allowedPrincipalServicesSet)
	policySummary.AllowedPrincipals = setToSortedSlice(statementsSummary.allowedPrincipalsSet)
	policySummary.AllowedPrincipalAccountIds = setToSortedSlice(statementsSummary.allowedPrincipalAccountIdsSet)
	policySummary.AllowedOrganizationIds = setToSortedSlice(statementsSummary.allowedOrganizationIds)
	policySummary.PublicStatementIds = setToSortedSlice(statementsSummary.publicStatementIds)
	policySummary.SharedStatementIds = setToSortedSlice(statementsSummary.sharedStatementIds)
	policySummary.PublicAccessLevels = statementsSummary.publicAccessLevels
	policySummary.SharedAccessLevels = statementsSummary.sharedAccessLevels
	policySummary.PrivateAccessLevels = statementsSummary.privateAccessLevels
	policySummary.IsPublic = statementsSummary.isPublic

	return policySummary, nil
}

type EvaluatedStatement struct {
	principal            EvaluatedPrincipal
	condition            EvaluatedCondition
	sid                  string
	availablePermissions AvailablePermissions
}

func evaluateStatements(statements []Statement, userAccountId string, allAvailablePermissions AllAvailablePermissions) ([]EvaluatedStatement, []EvaluatedStatement, error) {
	var currentEvaluatedStatements *[]EvaluatedStatement
	allowedEvaluatedStatements := make([]EvaluatedStatement, 0, len(statements))
	deniedEvaluatedStatements := make([]EvaluatedStatement, 0, len(statements))

	uniqueStatementIds := map[string]bool{}

	for statementIndex, statement := range statements {
		if !checkEffectValid(statement.Effect) {
			return allowedEvaluatedStatements, deniedEvaluatedStatements, fmt.Errorf("element Effect is invalid - valid choices are 'Allow' or 'Deny'")
		}

		if statement.Effect == "Deny" {
			currentEvaluatedStatements = &deniedEvaluatedStatements
		} else {
			currentEvaluatedStatements = &allowedEvaluatedStatements
		}

		// Conditions
		evaluatedCondition, err := evaluateCondition(statement.Condition, userAccountId)
		if err != nil {
			return allowedEvaluatedStatements, deniedEvaluatedStatements, err
		}

		// Principals
		hasResources := len(statement.Resource) > 0
		evaluatedPrincipal, err := evaluatePrincipal(statement.Principal, userAccountId, hasResources, evaluatedCondition.hasConditions)
		if err != nil {
			return allowedEvaluatedStatements, deniedEvaluatedStatements, err
		}

		// Before using Sid, let's check to see if it is unique
		sid := evaluatedSid(statement, statementIndex)
		if _, exists := uniqueStatementIds[sid]; exists {
			return allowedEvaluatedStatements, deniedEvaluatedStatements, fmt.Errorf("duplicate Sid found: %s", sid)
		}
		uniqueStatementIds[sid] = true

		evaluatedStatement := EvaluatedStatement{
			principal: evaluatedPrincipal,
			condition: evaluatedCondition,
			sid:       sid,
		}

		actionSet := map[string]bool{}
		for _, action := range statement.Action {
			actionSet[action] = true
		}
		//evaluatedStatement.availablePermissions = findAvailablePermissions(actionSet, allAvailablePermissions)
		evaluatedStatement.availablePermissions = allAvailablePermissions.findAvailablePermissions(actionSet)

		(*currentEvaluatedStatements) = append(*currentEvaluatedStatements, evaluatedStatement)
	}

	return allowedEvaluatedStatements, deniedEvaluatedStatements, nil
}

func checkEffectValid(effect string) bool {
	if effect == "Deny" || effect == "Allow" {
		return true
	}

	return false
}

func evaluateAccessLevel(statements StatementsSummary) string {
	if statements.isPublic {
		return "public"
	}

	if statements.isShared {
		return "shared"
	}

	return "private"
}

func evaluatedSid(statement Statement, statementIndex int) string {
	if statement.Sid == "" {
		return fmt.Sprintf("Statement[%d]", statementIndex+1)
	}

	return statement.Sid
}

type AllAvailablePermissions struct {
	servicePermissions map[string]Permissions
}

func loadAllAvailablePermissions() AllAvailablePermissions {
	allAvailablePermissions := AllAvailablePermissions{}

	servicePermissions := map[string]Permissions{}
	parliamentPermissions := getParliamentIamPermissions()

	for _, parliamentPermission := range parliamentPermissions {
		if _, exist := servicePermissions[parliamentPermission.Prefix]; !exist {
			privileges := []string{}
			accessLevel := map[string]string{}

			for _, priviledge := range parliamentPermission.Privileges {
				lowerPriviledge := strings.ToLower(priviledge.Privilege)
				privileges = append(privileges, lowerPriviledge)
				accessLevel[lowerPriviledge] = priviledge.AccessLevel
			}

			sort.Strings(privileges)

			servicePermissions[parliamentPermission.Prefix] = Permissions{
				privileges:  privileges,
				accessLevel: accessLevel,
			}
		}
	}

	allAvailablePermissions.servicePermissions = servicePermissions

	return allAvailablePermissions
}

type AvailablePermissions struct {
	isAllPermissions bool
	permissions      map[string]bool
}

func (allAvailablePermissions AllAvailablePermissions) findAvailablePermissions(actionSet map[string]bool) AvailablePermissions {
	if _, exists := actionSet["*"]; exists {
		return AvailablePermissions{isAllPermissions: true}
	}

	permissions := map[string]bool{}

	for action := range actionSet {
		permissionSummary := createPermissionSummary(action)

		if !permissionSummary.process {
			continue
		}

		// Find service
		if _, exists := allAvailablePermissions.servicePermissions[permissionSummary.service]; !exists {
			continue
		}

		servicePermissions := allAvailablePermissions.servicePermissions[permissionSummary.service]

		if permissionSummary.matcher == "" {
			if _, exists := servicePermissions.accessLevel[permissionSummary.priviledge]; exists {
				permission := fmt.Sprintf("%s:%s", permissionSummary.service, permissionSummary.priviledge)
				permissions[permission] = true
			}
			continue
		}

		// Find API Call
		privilegesLen := len(servicePermissions.privileges)
		checkIndex := sort.SearchStrings(servicePermissions.privileges, permissionSummary.priviledge)
		if checkIndex >= privilegesLen {
			continue
		}

		evaluatedPriviledgeLen := len(permissionSummary.priviledge)
		matcher := regexp.MustCompile(permissionSummary.matcher)
		for ; checkIndex < privilegesLen; checkIndex++ {
			currentPrivilege := servicePermissions.privileges[checkIndex]
			currentPrivilegeLen := len(currentPrivilege)

			splitIndex := int(math.Min(float64(currentPrivilegeLen), float64(evaluatedPriviledgeLen)))
			partialPriviledge := currentPrivilege[0:splitIndex]

			if partialPriviledge != permissionSummary.priviledge {
				break
			}
			if !matcher.MatchString(currentPrivilege) {
				continue
			}

			permission := fmt.Sprintf("%s:%s", permissionSummary.service, currentPrivilege)
			permissions[permission] = true
		}
	}

	return AvailablePermissions{permissions: permissions}
}

func (allAvailablePermissions AllAvailablePermissions) getAccessLevels(availablePermissions AvailablePermissions) map[string]bool {
	if availablePermissions.isAllPermissions {
		return map[string]bool{
			"List":                   true,
			"Permissions management": true,
			"Read":                   true,
			"Tagging":                true,
			"Write":                  true,
		}
	}

	accessLevels := map[string]bool{}

	for permission := range availablePermissions.permissions {
		actionSummary := createPermissionSummary(permission)

		if !actionSummary.process {
			continue
		}

		// Find service
		if _, exists := allAvailablePermissions.servicePermissions[actionSummary.service]; !exists {
			continue
		}

		servicePermissions := allAvailablePermissions.servicePermissions[actionSummary.service]

		if actionSummary.matcher == "" {
			if accessLevel, exists := servicePermissions.accessLevel[actionSummary.priviledge]; exists {
				accessLevels[accessLevel] = true
			}
			continue
		}
		// Find API Call
		privilegesLen := len(servicePermissions.privileges)
		checkIndex := sort.SearchStrings(servicePermissions.privileges, actionSummary.priviledge)
		if checkIndex >= privilegesLen {
			continue
		}

		evaluatedPriviledgeLen := len(actionSummary.priviledge)
		matcher := regexp.MustCompile(actionSummary.matcher)
		for ; checkIndex < privilegesLen; checkIndex++ {
			currentPrivilege := servicePermissions.privileges[checkIndex]
			currentPrivilegeLen := len(currentPrivilege)

			splitIndex := int(math.Min(float64(currentPrivilegeLen), float64(evaluatedPriviledgeLen)))
			partialPriviledge := currentPrivilege[0:splitIndex]

			if partialPriviledge != actionSummary.priviledge {
				break
			}
			if !matcher.MatchString(currentPrivilege) {
				continue
			}
			accessLevel := servicePermissions.accessLevel[currentPrivilege]
			accessLevels[accessLevel] = true
		}
	}

	return accessLevels
}

type PermissionSummary struct {
	process    bool
	service    string
	priviledge string
	matcher    string
}

func createPermissionSummary(action string) PermissionSummary {
	evaluated := PermissionSummary{}

	lowerAction := strings.ToLower(action)
	actionParts := strings.Split(lowerAction, ":")
	evaluated.service = actionParts[0]

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

type StatementsSummary struct {
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

func createStatementsSummary(statements []EvaluatedStatement, allAvailablePermissions AllAvailablePermissions) StatementsSummary {
	statementsSummary := StatementsSummary{
		publicStatementIds: map[string]bool{},
		sharedStatementIds: map[string]bool{},
	}

	publicAccessLevelSet := map[string]bool{}
	sharedAccessLevelSet := map[string]bool{}
	privateAccessLevelSet := map[string]bool{}

	for _, reducedStatement := range statements {
		// Does this statement have any actions and are the actions valid?
		// TODO: Actions valid?
		if !reducedStatement.availablePermissions.isAllPermissions && len(reducedStatement.availablePermissions.permissions) == 0 {
			continue
		}

		evaluatedAccessLevels := allAvailablePermissions.getAccessLevels(reducedStatement.availablePermissions)

		if len(evaluatedAccessLevels) == 0 {
			continue
		}

		statementsSummary.allowedOrganizationIds = mergeSets(
			statementsSummary.allowedOrganizationIds,
			reducedStatement.principal.allowedOrganizationIds,
			reducedStatement.condition.allowedOrganizationIds,
		)
		statementsSummary.allowedPrincipalAccountIdsSet = mergeSets(
			statementsSummary.allowedPrincipalAccountIdsSet,
			reducedStatement.principal.allowedPrincipalAccountIdsSet,
			reducedStatement.condition.allowedPrincipalAccountIdsSet,
		)
		statementsSummary.allowedPrincipalFederatedIdentitiesSet = mergeSets(
			statementsSummary.allowedPrincipalFederatedIdentitiesSet,
			reducedStatement.principal.allowedPrincipalFederatedIdentitiesSet,
			reducedStatement.condition.allowedPrincipalFederatedIdentitiesSet,
		)
		statementsSummary.allowedPrincipalServicesSet = mergeSets(
			statementsSummary.allowedPrincipalServicesSet,
			reducedStatement.principal.allowedPrincipalServicesSet,
			reducedStatement.condition.allowedPrincipalServicesSet,
		)
		statementsSummary.allowedPrincipalsSet = mergeSets(
			statementsSummary.allowedPrincipalsSet,
			reducedStatement.principal.allowedPrincipalsSet,
			reducedStatement.condition.allowedPrincipalsSet,
		)
		isPublic := reducedStatement.principal.isPublic || reducedStatement.condition.isPublic
		isShared := reducedStatement.principal.isShared || reducedStatement.condition.isShared
		isPrivate := reducedStatement.principal.isPrivate || reducedStatement.condition.isPrivate
		statementsSummary.isPublic = statementsSummary.isPublic || reducedStatement.principal.isPublic || reducedStatement.condition.isPublic
		statementsSummary.isShared = statementsSummary.isShared || reducedStatement.principal.isShared || reducedStatement.condition.isShared

		if isPublic {
			publicAccessLevelSet = mergeSet(publicAccessLevelSet, evaluatedAccessLevels)
			statementsSummary.publicStatementIds[reducedStatement.sid] = true
		}

		if isShared {
			sharedAccessLevelSet = mergeSet(sharedAccessLevelSet, evaluatedAccessLevels)
			//if len(sharedActionSet) > 0 {
			statementsSummary.sharedStatementIds[reducedStatement.sid] = true
			//}
		}

		if isPrivate {
			privateAccessLevelSet = mergeSet(privateAccessLevelSet, evaluatedAccessLevels)
		}
	}

	statementsSummary.publicAccessLevels = setToSortedSlice(publicAccessLevelSet)
	statementsSummary.sharedAccessLevels = setToSortedSlice(sharedAccessLevelSet)
	statementsSummary.privateAccessLevels = setToSortedSlice(privateAccessLevelSet)

	return statementsSummary
}

func generateStatementsSummary(allowedStatements []EvaluatedStatement, deniedStatements []EvaluatedStatement, allAvailablePermissions AllAvailablePermissions) StatementsSummary {
	reducedStatements := allowedStatements

	for _, deniedStatement := range deniedStatements {
		for reducedStatementIndex := range reducedStatements {

			for deniedPermission := range deniedStatement.availablePermissions.permissions {
				delete(reducedStatements[reducedStatementIndex].availablePermissions.permissions, deniedPermission)
			}
		}
	}

	return createStatementsSummary(reducedStatements, allAvailablePermissions)
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

type Permissions struct {
	privileges  []string
	accessLevel map[string]string
}

type EvaluatedCondition struct {
	allowedPrincipalFederatedIdentitiesSet map[string]bool
	allowedPrincipalServicesSet            map[string]bool
	allowedPrincipalsSet                   map[string]bool
	allowedPrincipalAccountIdsSet          map[string]bool
	allowedOrganizationIds                 map[string]bool
	isPublic                               bool
	isShared                               bool
	isPrivate                              bool
	hasConditions                          bool
}

func evaluateCondition(conditions map[string]interface{}, userAccountId string) (EvaluatedCondition, error) {
	var evaluatedCondition EvaluatedCondition

	for operator, conditionKey := range conditions {
		evaulatedOperator, evaluated := evaulateOperator(operator)
		if !evaluated {
			continue
		}

		if evaulatedOperator.isNegated {
			return evaluatedCondition, fmt.Errorf("TODO: Implement")
			// NOTE: Here we have an issue with the table.
			// 		 The problem is that if we say some principal is NOT an account, this means everything but.
			// 		 I do not know how to represent this in the current table design.
		}

		for conditionName, conditionValues := range conditionKey.(map[string]interface{}) {
			switch conditionName {
			case "aws:principalaccount":
				evaluatedCondition = evaluateAccountTypeCondition(conditionValues.([]string), evaulatedOperator, userAccountId)
				evaluatedCondition.hasConditions = true
			case "aws:sourceaccount":
				evaluatedCondition = evaluateAccountTypeCondition(conditionValues.([]string), evaulatedOperator, userAccountId)
				evaluatedCondition.hasConditions = true
			case "aws:sourceowner":
				evaluatedCondition = evaluateAccountTypeCondition(conditionValues.([]string), evaulatedOperator, userAccountId)
				evaluatedCondition.hasConditions = true
			case "aws:principalarn":
				evaluatedCondition = evaluateArnTypeCondition(conditionValues.([]string), evaulatedOperator, userAccountId)
				evaluatedCondition.hasConditions = true
			case "aws:sourcearn":
				evaluatedCondition = evaluateArnTypeCondition(conditionValues.([]string), evaulatedOperator, userAccountId)
				evaluatedCondition.hasConditions = true
			case "aws:principalorgid":
				evaluatedCondition = evaluateOrganizationCondition(conditionValues.([]string), evaulatedOperator, userAccountId)
				evaluatedCondition.hasConditions = true
			}
		}
	}

	return evaluatedCondition, nil
}

func evaluateArnTypeCondition(conditionValues []string, evaulatedOperator EvaluatedOperator, userAccountId string) EvaluatedCondition {
	evaluatedCondition := EvaluatedCondition{
		allowedPrincipalsSet:          map[string]bool{},
		allowedPrincipalAccountIdsSet: map[string]bool{},
	}

	for _, principal := range conditionValues {
		if evaulatedOperator.category != "string" && evaulatedOperator.category != "arn" {
			continue
		}

		if evaulatedOperator.isLike {
			if evaulatedOperator.category == "string" {
				evaluatedCondition.allowedPrincipalsSet[principal] = true
				// We need to pull the account out of a wildcard type
				// Assume that account is before any other numeric
				// There should always be 12 digits
				reAccountExtractor := regexp.MustCompile(`^.*[:\*\?]([0-9]{12})[:\*\?].*$`)
				arnAccount := reAccountExtractor.FindStringSubmatch(principal)
				if len(arnAccount) > 0 {
					account := arnAccount[1]
					if account != userAccountId {
						evaluatedCondition.isShared = true
					} else {
						evaluatedCondition.isPrivate = true
					}
					evaluatedCondition.allowedPrincipalAccountIdsSet[account] = true
				} else {
					evaluatedCondition.isPublic = true
					evaluatedCondition.allowedPrincipalAccountIdsSet["*"] = true
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
					evaluatedCondition.allowedPrincipalsSet[principal] = true
					evaluatedCondition.allowedPrincipalAccountIdsSet[account] = true
					evaluatedCondition.isPublic = true
					continue
				}

				if accountLength == 0 || accountLength != 12 {
					continue
				}

				if strings.Contains(account, "?") {
					evaluatedCondition.allowedPrincipalsSet[principal] = true
					evaluatedCondition.allowedPrincipalAccountIdsSet[account] = true
					evaluatedCondition.isPublic = true
					continue
				}

				re := regexp.MustCompile(`^[0-9]{12}$`)
				if !re.MatchString(account) {
					continue
				}

				evaluatedCondition.allowedPrincipalsSet[principal] = true
				evaluatedCondition.allowedPrincipalAccountIdsSet[account] = true

				if account != userAccountId {
					evaluatedCondition.isShared = true
					continue
				}

				evaluatedCondition.isPrivate = true
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

		evaluatedCondition.allowedPrincipalsSet[principal] = true
		evaluatedCondition.allowedPrincipalAccountIdsSet[account] = true

		if account == userAccountId {
			evaluatedCondition.isPrivate = true
		} else {
			evaluatedCondition.isShared = true
		}
	}

	return evaluatedCondition
}

func evaluateOrganizationCondition(conditionValues []string, evaulatedOperator EvaluatedOperator, userAccountId string) EvaluatedCondition {
	evaluatedCondition := EvaluatedCondition{
		allowedOrganizationIds: map[string]bool{},
	}

	for _, principal := range conditionValues {
		if evaulatedOperator.category != "string" {
			continue
		}

		organization := principal
		if evaulatedOperator.isLike {
			if organization == "*" || organization == "o-*" {
				evaluatedCondition.allowedOrganizationIds["o-*"] = true
				evaluatedCondition.isPublic = true
				continue
			}

			if !strings.HasPrefix(organization, "o-") {
				continue
			}

			evaluatedCondition.allowedOrganizationIds[organization] = true
			evaluatedCondition.isShared = true

			continue
		}

		if !strings.HasPrefix(organization, "o-") || strings.Contains(organization, "*") || strings.Contains(organization, "?") {
			continue
		}

		evaluatedCondition.allowedOrganizationIds[organization] = true
		evaluatedCondition.isShared = true
	}

	return evaluatedCondition
}

func evaluateAccountTypeCondition(conditionValues []string, evaulatedOperator EvaluatedOperator, userAccountId string) EvaluatedCondition {
	evaluatedCondition := EvaluatedCondition{
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
				evaluatedCondition.allowedPrincipalsSet[principal] = true
				evaluatedCondition.allowedPrincipalAccountIdsSet[account] = true
				evaluatedCondition.isPublic = true
				continue
			}

			if accountLength == 0 || accountLength != 12 {
				continue
			}

			if strings.Contains(account, "?") {
				evaluatedCondition.allowedPrincipalsSet[principal] = true
				evaluatedCondition.allowedPrincipalAccountIdsSet[account] = true
				evaluatedCondition.isPublic = true
				continue
			}

			reAccountFormat := regexp.MustCompile(`^[0-9]{12}$`)
			if !reAccountFormat.MatchString(account) {
				continue
			}

			evaluatedCondition.allowedPrincipalsSet[principal] = true
			evaluatedCondition.allowedPrincipalAccountIdsSet[account] = true
			if account != userAccountId {
				evaluatedCondition.isShared = true
				continue
			}

			evaluatedCondition.isPrivate = true
			continue
		}

		// Check if principal doesn't match an account ID, ignore
		reAccountFormat := regexp.MustCompile(`^[0-9]{12}$`)
		if !reAccountFormat.MatchString(principal) {
			continue
		}

		evaluatedCondition.allowedPrincipalsSet[principal] = true
		evaluatedCondition.allowedPrincipalAccountIdsSet[principal] = true

		if principal == userAccountId {
			evaluatedCondition.isPrivate = true
		} else {
			evaluatedCondition.isShared = true
		}
	}

	return evaluatedCondition
}

type EvaluatedPrincipal struct {
	allowedPrincipalFederatedIdentitiesSet map[string]bool
	allowedPrincipalServicesSet            map[string]bool
	allowedPrincipalsSet                   map[string]bool
	allowedPrincipalAccountIdsSet          map[string]bool
	allowedOrganizationIds                 map[string]bool
	isPublic                               bool
	isShared                               bool
	isPrivate                              bool
}

func evaluatePrincipal(principal Principal, userAccountId string, hasResources bool, hasConditions bool) (EvaluatedPrincipal, error) {
	evaluatedPrincipal := EvaluatedPrincipal{
		allowedPrincipalFederatedIdentitiesSet: map[string]bool{},
		allowedPrincipalServicesSet:            map[string]bool{},
		allowedPrincipalsSet:                   map[string]bool{},
		allowedPrincipalAccountIdsSet:          map[string]bool{},
	}

	if len(principal) == 0 && hasResources && !hasConditions {
		evaluatedPrincipal.allowedPrincipalsSet[userAccountId] = true
		evaluatedPrincipal.allowedPrincipalAccountIdsSet[userAccountId] = true
		evaluatedPrincipal.isPrivate = true

		return evaluatedPrincipal, nil
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
					evaluatedPrincipal.isPublic = true
					evaluatedPrincipal.allowedPrincipalAccountIdsSet[account] = true
				} else {
					if reIsAwsAccount.MatchString(principalItem) {
						account = principalItem
					} else if reIsAwsResource.MatchString(principalItem) {
						arnAccount := reIsAwsResource.FindStringSubmatch(principalItem)
						account = arnAccount[1]
					} else {
						return evaluatedPrincipal, fmt.Errorf("unabled to parse arn or account: %s", principalItem)
					}

					if userAccountId != account {
						evaluatedPrincipal.isShared = true
					} else {
						evaluatedPrincipal.isPrivate = true
					}
					evaluatedPrincipal.allowedPrincipalAccountIdsSet[account] = true
				}

				evaluatedPrincipal.allowedPrincipalsSet[principalItem] = true
			case "Service":
				evaluatedPrincipal.allowedPrincipalServicesSet[principalItem] = true
				evaluatedPrincipal.isPublic = true
			case "Federated":
				evaluatedPrincipal.allowedPrincipalFederatedIdentitiesSet[principalItem] = true
				evaluatedPrincipal.isPrivate = true
			}
		}
	}

	return evaluatedPrincipal, nil
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
