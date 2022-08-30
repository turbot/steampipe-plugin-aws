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

func checkPolicyValidity(policy Policy) (bool, error) {
	if policy.Version == "" {
		return false, fmt.Errorf("policy element Version is missing")
	}

	if policy.Version != "2012-10-17" && policy.Version != "2008-10-17" {
		return false, fmt.Errorf("unsupported value for policy element Version: '%s' - values supported are '2012-10-17' or '2008-10-17'", policy.Version)
	}

	for _, statement := range policy.Statements {
		if statement.Effect == "" {
			return false, fmt.Errorf("policy element Effect is missing")
		}

		if statement.Effect != "Deny" && statement.Effect != "Allow" {
			return false, fmt.Errorf("unsupported value for policy element Effect: '%s' - values supported are 'Allow' or 'Deny'", statement.Effect)
		}

		if len(statement.Principal) == 0 {
			return false, fmt.Errorf("policy element Principal is missing")
		}

		// if len(statement.Resource) == 0 {
		// 	return false, fmt.Errorf("policy element Resource is missing")
		// }
	}

	return true, nil
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

	if policyContent == "" || policyContent == "{}" {
		return policySummary, nil
	}

	policyInterface, err := canonicalPolicy(policyContent)
	if err != nil {
		return policySummary, err
	}

	allAvailablePermissions := loadAllAvailablePermissions()

	policy := policyInterface.(Policy)

	if valid, err := checkPolicyValidity(policy); !valid {
		return policySummary, err
	}

	allowedEvaluatedStatements, deniedEvaluatedStatements, err := evaluateStatements(policy.Statements, userAccountId, allAvailablePermissions)
	if err != nil {
		return policySummary, err
	}

	reducedStatements := reduceStatements(allowedEvaluatedStatements, deniedEvaluatedStatements)
	statementsSummary := generateStatementsSummary(reducedStatements, allAvailablePermissions)

	policySummary.AccessLevel = evaluateAccessLevel(statementsSummary)
	policySummary.AllowedPrincipalFederatedIdentities = setToSortedSlice(statementsSummary.allowedPrincipalFederatedIdentitiesSet)
	policySummary.AllowedPrincipalServices = setToSortedSlice(statementsSummary.allowedPrincipalServicesSet)
	policySummary.AllowedPrincipals = setToSortedSlice(statementsSummary.allowedPrincipalsSet)
	policySummary.AllowedPrincipalAccountIds = setToSortedSlice(statementsSummary.allowedPrincipalAccountIdsSet)
	policySummary.AllowedOrganizationIds = setToSortedSlice(statementsSummary.allowedOrganizationIdsSet)
	policySummary.PublicStatementIds = setToSortedSlice(statementsSummary.publicStatementIds)
	policySummary.SharedStatementIds = setToSortedSlice(statementsSummary.sharedStatementIds)
	policySummary.PublicAccessLevels = statementsSummary.publicAccessLevels
	policySummary.SharedAccessLevels = statementsSummary.sharedAccessLevels
	policySummary.PrivateAccessLevels = statementsSummary.privateAccessLevels
	policySummary.IsPublic = statementsSummary.isPublic

	return policySummary, nil
}

type EvaluatedStatement struct {
	availablePermissions AvailablePermissions
	isPrivate            bool
	isPublic             bool
	isShared             bool
	principal            string
	principalType        string
	resource             string
	sid                  string
}

func (evaluatedStatement *EvaluatedStatement) ApplyDenyStatement(denyStatement EvaluatedStatement) {
	if denyStatement.principal == "*" && denyStatement.principalType == "arn" {
		evaluatedStatement.availablePermissions.permissions = map[string]bool{}
		evaluatedStatement.availablePermissions.isAllPermissions = false
		return
	}

	if denyStatement.principalType == "account" && evaluatedStatement.principalType == "arn" {
		account := extractAccountInPlaceFromArn(evaluatedStatement.principal)

		if account == denyStatement.principal {
			evaluatedStatement.availablePermissions.RemovePermissions(denyStatement.availablePermissions)
		}

		return
	}

	if denyStatement.principalType == evaluatedStatement.principalType {
		denyPrincipalValue := MakePolicyValue(denyStatement.principal)

		if denyPrincipalValue.Contains(evaluatedStatement.principal) {
			denyResourceValue := MakePolicyValue(denyStatement.resource)
			if denyResourceValue.Contains(evaluatedStatement.resource) {
				evaluatedStatement.availablePermissions.RemovePermissions(denyStatement.availablePermissions)
				evaluatedStatement.availablePermissions.isAllPermissions = false
			}
		}
	}
}

func evaluateStatements(statements []Statement, userAccountId string, allAvailablePermissions AllAvailablePermissions) ([]EvaluatedStatement, []EvaluatedStatement, error) {
	var currentEvaluatedStatements *[]EvaluatedStatement

	allowedEvaluatedStatements := make([]EvaluatedStatement, 0, len(statements))
	deniedEvaluatedStatements := make([]EvaluatedStatement, 0, len(statements))

	uniqueStatementIds := map[string]bool{}

	for statementIndex, statement := range statements {
		if statement.Effect == "Deny" {
			currentEvaluatedStatements = &deniedEvaluatedStatements
		} else {
			currentEvaluatedStatements = &allowedEvaluatedStatements
		}

		// Principals
		evaluatedPrincipal, err := evaluatePrincipal(statement.Principal, userAccountId)
		if err != nil {
			return allowedEvaluatedStatements, deniedEvaluatedStatements, err
		}

		// Conditions
		evaluatedCondition, err := refineUsingConditions(evaluatedPrincipal, statement.Condition)
		if err != nil {
			return allowedEvaluatedStatements, deniedEvaluatedStatements, err
		}

		// Before using Sid, let's check to see if it is unique
		sid := evaluatedSid(statement, statementIndex)
		if _, exists := uniqueStatementIds[sid]; exists {
			return allowedEvaluatedStatements, deniedEvaluatedStatements, fmt.Errorf("duplicate Sid found: %s", sid)
		}
		uniqueStatementIds[sid] = true

		allowedOrganizationIdsSet := evaluatedCondition.allowedOrganizationIdsSet
		allowedPrincipalFederatedIdentitiesSet := evaluatedCondition.allowedPrincipalFederatedIdentitiesSet
		allowedPrincipalServicesSet := evaluatedCondition.allowedPrincipalServicesSet
		allowedPrincipalsArnsSet := evaluatedCondition.allowedPrincipalsArnsSet
		allowedPrincipalsAccountsSet := evaluatedCondition.allowedPrincipalsAccountsSet
		isPublic := evaluatedCondition.isPublic
		isShared := evaluatedCondition.isShared
		isPrivate := evaluatedCondition.isPrivate

		actionSet := map[string]bool{}
		for _, action := range statement.Action {
			actionSet[action] = true
		}

		availablePermissions := allAvailablePermissions.FindAvailablePermissions(actionSet)

		// Resources
		var resources []string
		if len(statement.Resource) > 0 {
			resources = statement.Resource
		} else {
			resources = []string{""}
		}

		for _, resource := range resources {
			// Create individual statements here
			for allowedOrganizationId := range allowedOrganizationIdsSet {
				newStatement := EvaluatedStatement{
					availablePermissions: availablePermissions.Copy(),
					isPrivate:            isPrivate,
					isPublic:             isPublic,
					isShared:             isShared,
					principal:            allowedOrganizationId,
					principalType:        "organization",
					resource:             resource,
					sid:                  sid,
				}

				(*currentEvaluatedStatements) = append(*currentEvaluatedStatements, newStatement)
			}

			for allowedPrincipalFederatedIdentity := range allowedPrincipalFederatedIdentitiesSet {
				newStatement := EvaluatedStatement{
					availablePermissions: availablePermissions.Copy(),
					isPrivate:            isPrivate,
					isPublic:             isPublic,
					isShared:             isShared,
					principal:            allowedPrincipalFederatedIdentity,
					principalType:        "federated",
					resource:             resource,
					sid:                  sid,
				}

				(*currentEvaluatedStatements) = append(*currentEvaluatedStatements, newStatement)
			}

			for allowedPrincipalService := range allowedPrincipalServicesSet {
				newStatement := EvaluatedStatement{
					availablePermissions: availablePermissions.Copy(),
					isPrivate:            isPrivate,
					isPublic:             isPublic,
					isShared:             isShared,
					principal:            allowedPrincipalService,
					principalType:        "service",
					resource:             resource,
					sid:                  sid,
				}

				(*currentEvaluatedStatements) = append(*currentEvaluatedStatements, newStatement)
			}

			for allowedPrincipalArn := range allowedPrincipalsArnsSet {
				newStatement := EvaluatedStatement{
					availablePermissions: availablePermissions.Copy(),
					isPrivate:            isPrivate,
					isPublic:             isPublic,
					isShared:             isShared,
					principal:            allowedPrincipalArn,
					principalType:        "arn",
					resource:             resource,
					sid:                  sid,
				}

				(*currentEvaluatedStatements) = append(*currentEvaluatedStatements, newStatement)
			}

			for allowedPrincipalAccount := range allowedPrincipalsAccountsSet {
				newStatement := EvaluatedStatement{
					availablePermissions: availablePermissions.Copy(),
					isPrivate:            isPrivate,
					isPublic:             isPublic,
					isShared:             isShared,
					principal:            allowedPrincipalAccount,
					principalType:        "account",
					resource:             resource,
					sid:                  sid,
				}

				(*currentEvaluatedStatements) = append(*currentEvaluatedStatements, newStatement)
			}
		}
	}

	return allowedEvaluatedStatements, deniedEvaluatedStatements, nil
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

func (availablePermissions AvailablePermissions) HasPermissions() bool {
	return availablePermissions.isAllPermissions || len(availablePermissions.permissions) > 0
}

func (availablePermissions AvailablePermissions) IsAllPermissions() bool {
	return availablePermissions.isAllPermissions
}

func (availablePermissions AvailablePermissions) RemovePermissions(permissionsToRemove AvailablePermissions) {
	for permission := range permissionsToRemove.permissions {
		delete(availablePermissions.permissions, permission)
	}
}

func (availablePermissions AvailablePermissions) Copy() AvailablePermissions {
	copy := AvailablePermissions{
		isAllPermissions: availablePermissions.isAllPermissions,
		permissions:      map[string]bool{},
	}

	for key, value := range availablePermissions.permissions {
		copy.permissions[key] = value
	}

	return copy
}

func (allAvailablePermissions AllAvailablePermissions) FindAvailablePermissions(actionSet map[string]bool) AvailablePermissions {
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

func (allAvailablePermissions AllAvailablePermissions) GetAccessLevels(availablePermissions AvailablePermissions) map[string]bool {
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
	matcher    string
	priviledge string
	process    bool
	service    string
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
	allowedOrganizationIdsSet              map[string]bool
	allowedPrincipalAccountIdsSet          map[string]bool
	allowedPrincipalFederatedIdentitiesSet map[string]bool
	allowedPrincipalServicesSet            map[string]bool
	allowedPrincipalsSet                   map[string]bool
	isPublic                               bool
	isShared                               bool
	privateAccessLevels                    []string
	publicAccessLevels                     []string
	publicStatementIds                     map[string]bool
	sharedAccessLevels                     []string
	sharedStatementIds                     map[string]bool
}

func generateStatementsSummary(statements []EvaluatedStatement, allAvailablePermissions AllAvailablePermissions) StatementsSummary {
	statementsSummary := StatementsSummary{
		allowedOrganizationIdsSet:              map[string]bool{},
		allowedPrincipalAccountIdsSet:          map[string]bool{},
		allowedPrincipalFederatedIdentitiesSet: map[string]bool{},
		allowedPrincipalServicesSet:            map[string]bool{},
		allowedPrincipalsSet:                   map[string]bool{},
		publicStatementIds:                     map[string]bool{},
		sharedStatementIds:                     map[string]bool{},
	}

	publicAccessLevelSet := map[string]bool{}
	sharedAccessLevelSet := map[string]bool{}
	privateAccessLevelSet := map[string]bool{}

	for _, reducedStatement := range statements {
		// Does this statement have any actions and are the actions valid?
		if !reducedStatement.availablePermissions.IsAllPermissions() && len(reducedStatement.availablePermissions.permissions) == 0 {
			continue
		}

		evaluatedAccessLevels := allAvailablePermissions.GetAccessLevels(reducedStatement.availablePermissions)

		if len(evaluatedAccessLevels) == 0 {
			continue
		}

		switch reducedStatement.principalType {
		case "federated":
			statementsSummary.allowedPrincipalFederatedIdentitiesSet[reducedStatement.principal] = true
		case "organization":
			statementsSummary.allowedOrganizationIdsSet[reducedStatement.principal] = true
		case "arn":
			account := extractAccountFromArn(reducedStatement.principal)
			statementsSummary.allowedPrincipalsSet[reducedStatement.principal] = true
			statementsSummary.allowedPrincipalAccountIdsSet[account] = true
		case "account":
			account := extractAccount(reducedStatement.principal)

			statementsSummary.allowedPrincipalsSet[reducedStatement.principal] = true
			statementsSummary.allowedPrincipalAccountIdsSet[account] = true
		case "service":
			statementsSummary.allowedPrincipalServicesSet[reducedStatement.principal] = true
		}

		statementsSummary.isPublic = statementsSummary.isPublic || reducedStatement.isPublic
		statementsSummary.isShared = statementsSummary.isShared || reducedStatement.isShared

		if reducedStatement.isPublic {
			publicAccessLevelSet = mergeSet(publicAccessLevelSet, evaluatedAccessLevels)
			statementsSummary.publicStatementIds[reducedStatement.sid] = true
		}

		if reducedStatement.isShared {
			sharedAccessLevelSet = mergeSet(sharedAccessLevelSet, evaluatedAccessLevels)
			statementsSummary.sharedStatementIds[reducedStatement.sid] = true
		}

		if reducedStatement.isPrivate {
			privateAccessLevelSet = mergeSet(privateAccessLevelSet, evaluatedAccessLevels)
		}
	}

	statementsSummary.publicAccessLevels = setToSortedSlice(publicAccessLevelSet)
	statementsSummary.sharedAccessLevels = setToSortedSlice(sharedAccessLevelSet)
	statementsSummary.privateAccessLevels = setToSortedSlice(privateAccessLevelSet)

	// NOTE: Removing this behaviour for now. We need further discussion on how the table will deal with services.
	// // finally enrich the account IDs and Principals if there are no account ids or principals with "*"
	// if statementsSummary.isPublic &&
	// 	(len(statementsSummary.allowedPrincipalAccountIdsSet) == 0 && len(statementsSummary.allowedPrincipalsSet) == 0) {
	// 	statementsSummary.allowedPrincipalsSet["*"] = true
	// 	statementsSummary.allowedPrincipalAccountIdsSet["*"] = true
	// }

	return statementsSummary
}

func reduceStatements(allowedStatements []EvaluatedStatement, deniedStatements []EvaluatedStatement) []EvaluatedStatement {
	reducedStatements := allowedStatements

	for _, deniedStatement := range deniedStatements {
		updatedActiveStatements := []EvaluatedStatement{}
		for _, activeStatement := range reducedStatements {
			activeStatement.ApplyDenyStatement(deniedStatement)
			if !activeStatement.availablePermissions.HasPermissions() {
				continue
			}

			updatedActiveStatements = append(updatedActiveStatements, activeStatement)
		}

		reducedStatements = updatedActiveStatements
	}

	return reducedStatements
}

type EvaluatedOperator struct {
	category   string
	isCaseless bool
	isLike     bool
	isNegated  bool
}

func evaulateOperator(operator string) (EvaluatedOperator, bool) {
	evaulatedOperator := EvaluatedOperator{}
	evaluated := true

	operator = strings.ToLower(operator)
	if strings.HasSuffix(operator, "ifexists") {
		return evaulatedOperator, false
	}

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
	accessLevel map[string]string
	privileges  []string
}

type EvaluatedCondition struct {
	allowedOrganizationIdsSet              map[string]bool
	allowedPrincipalFederatedIdentitiesSet map[string]bool
	allowedPrincipalServicesSet            map[string]bool
	allowedPrincipalsAccountsSet           map[string]bool
	allowedPrincipalsArnsSet               map[string]bool
	isPrivate                              bool
	isPublic                               bool
	isShared                               bool
}

func (evaluatedCondition *EvaluatedCondition) Merge(other EvaluatedCondition) {
	evaluatedCondition.allowedOrganizationIdsSet = mergeSet(evaluatedCondition.allowedOrganizationIdsSet, other.allowedOrganizationIdsSet)
	evaluatedCondition.allowedPrincipalFederatedIdentitiesSet = mergeSet(evaluatedCondition.allowedPrincipalFederatedIdentitiesSet, other.allowedPrincipalFederatedIdentitiesSet)
	evaluatedCondition.allowedPrincipalServicesSet = mergeSet(evaluatedCondition.allowedPrincipalServicesSet, other.allowedPrincipalServicesSet)
	evaluatedCondition.allowedPrincipalsArnsSet = mergeSet(evaluatedCondition.allowedPrincipalsArnsSet, other.allowedPrincipalsArnsSet)
	evaluatedCondition.allowedPrincipalsAccountsSet = mergeSet(evaluatedCondition.allowedPrincipalsAccountsSet, other.allowedPrincipalsAccountsSet)

	evaluatedCondition.isPublic = evaluatedCondition.isPublic || other.isPublic
	evaluatedCondition.isShared = evaluatedCondition.isShared || other.isShared
	evaluatedCondition.isPrivate = evaluatedCondition.isPrivate || other.isPrivate
}

func refineUsingConditions(evaluatedPrincipal EvaluatedPrincipal, conditions map[string]interface{}) (EvaluatedCondition, error) {
	processed := false
	evaluatedCondition := EvaluatedCondition{}

	for operator, conditionKey := range conditions {
		evaulatedOperator, evaluated := evaulateOperator(operator)
		if !evaluated || evaulatedOperator.isNegated {
			continue
		}

		for conditionName, conditionValues := range conditionKey.(map[string]interface{}) {
			processed = true
			switch conditionName {
			case "aws:principalaccount":
				partialEvaluatedCondition := evaluatePrincipalAccountTypeCondition(evaluatedPrincipal, conditionValues.([]string), evaulatedOperator)
				evaluatedCondition.Merge(partialEvaluatedCondition)
			case "aws:sourceaccount":
				partialEvaluatedCondition := evaluateSourceAccountTypeCondition(evaluatedPrincipal, conditionValues.([]string), evaulatedOperator)
				evaluatedCondition.Merge(partialEvaluatedCondition)
			case "aws:sourceowner":
				partialEvaluatedCondition := evaluateSourceAccountTypeCondition(evaluatedPrincipal, conditionValues.([]string), evaulatedOperator)
				evaluatedCondition.Merge(partialEvaluatedCondition)
			case "aws:principalorgid":
				partialEvaluatedCondition := evaluatePrincipalOrganizationIdCondition(evaluatedPrincipal, conditionValues.([]string), evaulatedOperator)
				evaluatedCondition.Merge(partialEvaluatedCondition)
			case "aws:principalarn":
				partialEvaluatedCondition := evaluatePrincipalArnTypeCondition(evaluatedPrincipal, conditionValues.([]string), evaulatedOperator)
				evaluatedCondition.Merge(partialEvaluatedCondition)
			case "aws:sourcearn":
				partialEvaluatedCondition := evaluateSourceArnTypeCondition(evaluatedPrincipal, conditionValues.([]string), evaulatedOperator)
				evaluatedCondition.Merge(partialEvaluatedCondition)
			}
		}
	}

	if !processed {
		return evaluatedPrincipal.toEvaluatedCondition(), nil
	}

	return evaluatedCondition, nil
}

// TODO: We have a problem with the following code as it evaluates the Principal which is incorrect
func evaluateSourceArnTypeCondition(evaluatedPrincipal EvaluatedPrincipal, conditionValues []string, evaulatedOperator EvaluatedOperator) EvaluatedCondition {
	processed := false
	allowedPrincipalsAccountsSet := map[string]bool{}
	allowedPrincipalsArnsSet := map[string]bool{}

	isPublic := evaluatedPrincipal.isAwsPublic || evaluatedPrincipal.isServicePublic || evaluatedPrincipal.isFederatedPublic
	isShared := evaluatedPrincipal.isFederatedShared
	isPrivate := false

	for _, conditionValue := range conditionValues {
		if evaulatedOperator.category != "string" && evaulatedOperator.category != "arn" {
			continue
		}

		// value "*" means conditionAccount was invalid
		var conditionAccount string

		if evaulatedOperator.category == "arn" {
			conditionAccount = extractAccountInPlaceFromArn(conditionValue)
			if conditionAccount == "" {
				continue
			}
		} else if evaulatedOperator.category == "string" && !evaulatedOperator.isLike {
			conditionAccount = extractAccountInPlaceFromArn(conditionValue)
			if strings.Contains(conditionValue, "*") || strings.Contains(conditionValue, "?") || conditionAccount == "" {
				continue
			}
		} else {
			conditionAccount = extractAccountFromArn(conditionValue)
			if conditionAccount == "" {
				conditionAccount = "*"
			}
		}

		processed = true

		conditionPolicyValue := MakePolicyValue(conditionValue)

		// Simple direct comparison here
		for principalArns := range evaluatedPrincipal.allowedPrincipalsArnsSet {
			if principalArns == "*" {
				principalPolicyValue := MakePolicyValue("*")
				resolved := principalPolicyValue.Intersection(conditionPolicyValue)

				allowedPrincipalsArnsSet[resolved] = true
				if resolved == "*" || strings.Contains(conditionAccount, "*") || strings.Contains(conditionAccount, "?") {
					isPublic = true
					continue
				}

				isPublic = false
				if conditionAccount != evaluatedPrincipal.userAccountId {
					isShared = true
				} else {
					isPrivate = true
				}
			} else if conditionPolicyValue.Contains(principalArns) {
				if !conditionPolicyValue.Contains(principalArns) {
					continue
				}
				allowedPrincipalsArnsSet[principalArns] = true

				principalAccount := extractAccountFromArn(principalArns)
				if principalAccount != evaluatedPrincipal.userAccountId {
					isShared = true
				} else {
					isPrivate = true
				}
			}
		}

		accountPolicyValue := MakePolicyValue(conditionAccount)

		for principalAccount := range evaluatedPrincipal.allowedPrincipalsAccountsSet {
			if accountPolicyValue.Contains(principalAccount) {
				if evaulatedOperator.category == "arn" {
					replacementArn := updateAccountInArn(conditionValue, principalAccount)
					allowedPrincipalsArnsSet[replacementArn] = true
				} else {
					if conditionAccount == "*" {
						if evaulatedOperator.category != "arn" {
							allowedPrincipalsAccountsSet[principalAccount] = true
						} else {
							allowedPrincipalsAccountsSet[principalAccount] = true
						}
					} else {
						resolved := strings.Replace(conditionValue, conditionAccount, principalAccount, 1)

						allowedPrincipalsArnsSet[resolved] = true
					}
				}

				if principalAccount != evaluatedPrincipal.userAccountId {
					isShared = true
				} else {
					isPrivate = true
				}
			}
		}
	}

	if processed {
		return EvaluatedCondition{
			allowedOrganizationIdsSet:              evaluatedPrincipal.allowedOrganizationIdsSet,
			allowedPrincipalFederatedIdentitiesSet: evaluatedPrincipal.allowedPrincipalFederatedIdentitiesSet,
			allowedPrincipalsAccountsSet:           allowedPrincipalsAccountsSet,
			allowedPrincipalsArnsSet:               allowedPrincipalsArnsSet,
			allowedPrincipalServicesSet:            evaluatedPrincipal.allowedPrincipalServicesSet,
			isPublic:                               isPublic,
			isShared:                               isShared,
			isPrivate:                              isPrivate,
		}
	}

	return evaluatedPrincipal.toEvaluatedCondition()
}

func evaluatePrincipalArnTypeCondition(evaluatedPrincipal EvaluatedPrincipal, conditionValues []string, evaulatedOperator EvaluatedOperator) EvaluatedCondition {
	processed := false
	allowedPrincipalsAccountsSet := map[string]bool{}
	allowedPrincipalsArnsSet := map[string]bool{}

	isPublic := evaluatedPrincipal.isAwsPublic || evaluatedPrincipal.isServicePublic || evaluatedPrincipal.isFederatedPublic
	isShared := evaluatedPrincipal.isFederatedShared
	isPrivate := false

	for _, conditionValue := range conditionValues {
		if evaulatedOperator.category != "string" && evaulatedOperator.category != "arn" {
			continue
		}

		// value "*" means conditionAccount was invalid
		var conditionAccount string

		if evaulatedOperator.category == "arn" {
			conditionAccount = extractAccountInPlaceFromArn(conditionValue)
			if conditionAccount == "" {
				continue
			}
		} else if evaulatedOperator.category == "string" && !evaulatedOperator.isLike {
			conditionAccount = extractAccountInPlaceFromArn(conditionValue)
			if strings.Contains(conditionValue, "*") || strings.Contains(conditionValue, "?") || conditionAccount == "" {
				continue
			}
		} else {
			conditionAccount = extractAccountFromArn(conditionValue)
			if conditionAccount == "" {
				conditionAccount = "*"
			}
		}

		processed = true

		conditionPolicyValue := MakePolicyValue(conditionValue)

		// Simple direct comparison here
		for principalArns := range evaluatedPrincipal.allowedPrincipalsArnsSet {
			if principalArns == "*" {
				principalPolicyValue := MakePolicyValue("*")
				resolved := principalPolicyValue.Intersection(conditionPolicyValue)

				allowedPrincipalsArnsSet[resolved] = true
				if resolved == "*" || strings.Contains(conditionAccount, "*") || strings.Contains(conditionAccount, "?") {
					isPublic = true
					continue
				}

				isPublic = false
				if conditionAccount != evaluatedPrincipal.userAccountId {
					isShared = true
				} else {
					isPrivate = true
				}
			} else if conditionPolicyValue.Contains(principalArns) {
				if !conditionPolicyValue.Contains(principalArns) {
					continue
				}
				allowedPrincipalsArnsSet[principalArns] = true

				principalAccount := extractAccountFromArn(principalArns)
				if principalAccount != evaluatedPrincipal.userAccountId {
					isShared = true
				} else {
					isPrivate = true
				}
			}
		}

		accountPolicyValue := MakePolicyValue(conditionAccount)

		for principalAccount := range evaluatedPrincipal.allowedPrincipalsAccountsSet {
			if accountPolicyValue.Contains(principalAccount) {
				if evaulatedOperator.category == "arn" {
					replacementArn := updateAccountInArn(conditionValue, principalAccount)
					allowedPrincipalsArnsSet[replacementArn] = true
				} else {
					if conditionAccount == "*" {
						if evaulatedOperator.category != "arn" {
							allowedPrincipalsAccountsSet[principalAccount] = true
						} else {
							allowedPrincipalsAccountsSet[principalAccount] = true
						}
					} else {
						resolved := strings.Replace(conditionValue, conditionAccount, principalAccount, 1)

						allowedPrincipalsArnsSet[resolved] = true
					}
				}

				if principalAccount != evaluatedPrincipal.userAccountId {
					isShared = true
				} else {
					isPrivate = true
				}
			}
		}
	}

	if processed {
		return EvaluatedCondition{
			allowedOrganizationIdsSet:              evaluatedPrincipal.allowedOrganizationIdsSet,
			allowedPrincipalFederatedIdentitiesSet: evaluatedPrincipal.allowedPrincipalFederatedIdentitiesSet,
			allowedPrincipalsAccountsSet:           allowedPrincipalsAccountsSet,
			allowedPrincipalsArnsSet:               allowedPrincipalsArnsSet,
			allowedPrincipalServicesSet:            evaluatedPrincipal.allowedPrincipalServicesSet,
			isPublic:                               isPublic,
			isShared:                               isShared,
			isPrivate:                              isPrivate,
		}
	}

	return evaluatedPrincipal.toEvaluatedCondition()
}

func updateAccountInArn(arn string, account string) string {
	splitArn := strings.Split(arn, ":")

	// There should always be an account
	if len(splitArn) < 6 {
		return arn
	}

	splitArn[4] = account

	return strings.Join(splitArn, ":")
}

func evaluatePrincipalOrganizationIdCondition(evaluatedPrincipal EvaluatedPrincipal, conditionValues []string, evaulatedOperator EvaluatedOperator) EvaluatedCondition {
	processed := false
	allowedOrganizationIdsSet := map[string]bool{}
	isPublic := evaluatedPrincipal.isAwsPublic || evaluatedPrincipal.isServicePublic || evaluatedPrincipal.isFederatedPublic
	isShared := evaluatedPrincipal.isFederatedShared
	isPrivate := false

	for _, principal := range conditionValues {
		if evaulatedOperator.category != "string" {
			continue
		}

		organization := principal
		if evaulatedOperator.isLike {
			if organization == "*" || organization == "o-*" {
				allowedOrganizationIdsSet["o-*"] = true
				isPublic = true
				processed = true
				continue
			}

			if !strings.HasPrefix(organization, "o-") {
				continue
			}

			allowedOrganizationIdsSet[organization] = true
			isShared = true
			processed = true
			continue
		}

		if !strings.HasPrefix(organization, "o-") || strings.Contains(organization, "*") || strings.Contains(organization, "?") {
			continue
		}

		allowedOrganizationIdsSet[organization] = true
		isShared = true
		processed = true
	}

	if processed {
		return EvaluatedCondition{
			allowedOrganizationIdsSet:              allowedOrganizationIdsSet,
			allowedPrincipalFederatedIdentitiesSet: evaluatedPrincipal.allowedPrincipalFederatedIdentitiesSet,
			allowedPrincipalsAccountsSet:           evaluatedPrincipal.allowedPrincipalsAccountsSet,
			allowedPrincipalsArnsSet:               evaluatedPrincipal.allowedPrincipalsArnsSet,
			allowedPrincipalServicesSet:            evaluatedPrincipal.allowedPrincipalServicesSet,
			isPublic:                               isPublic,
			isShared:                               isShared,
			isPrivate:                              isPrivate,
		}
	}

	return evaluatedPrincipal.toEvaluatedCondition()
}

// TODO: We have a problem with the following code as it evaluates the Principal which is incorrect
func evaluateSourceAccountTypeCondition(evaluatedPrincipal EvaluatedPrincipal, conditionValues []string, evaulatedOperator EvaluatedOperator) EvaluatedCondition {
	processed := false
	allowedPrincipalsAccountsSet := map[string]bool{}
	allowedPrincipalsArnsSet := map[string]bool{}

	isPublic := evaluatedPrincipal.isFederatedPublic
	isShared := evaluatedPrincipal.isFederatedShared
	isPrivate := false

	if len(evaluatedPrincipal.allowedPrincipalServicesSet) > 0 {
		for _, conditionValue := range conditionValues {
			if evaulatedOperator.category != "string" {
				continue
			}

			if evaulatedOperator.isLike {
				// Sanity check the input first

				// Regex to allow for account: ["222244446666", "22224444666?", "22224?446666", ...] and must be exactly 12
				reAccountFormat := regexp.MustCompile(`^[0-9\?]{12}$`)

				// Not OK
				if !(strings.Contains(conditionValue, "*") && len(conditionValue) <= 12) && !reAccountFormat.MatchString(conditionValue) {
					continue
				}
			} else {
				// Regex to allow for account: ["222244446666", ...] and must be exactly 12
				reAccountFormat := regexp.MustCompile(`^[0-9]{12}$`)
				// Not OK
				if !reAccountFormat.MatchString(conditionValue) {
					continue
				}
			}

			processed = true

			if conditionValue == "*" {
				isPublic = true
				allowedPrincipalsAccountsSet[conditionValue] = true
			} else if strings.Contains(conditionValue, "*") || strings.Contains(conditionValue, "?") {
				isPublic = true
				allowedPrincipalsAccountsSet[conditionValue] = true
			} else {
				allowedPrincipalsAccountsSet[conditionValue] = true
				if conditionValue != evaluatedPrincipal.userAccountId {
					isShared = true
				} else {
					isPrivate = true
				}
			}

		}
	}

	if processed {
		return EvaluatedCondition{
			allowedOrganizationIdsSet:              evaluatedPrincipal.allowedOrganizationIdsSet,
			allowedPrincipalFederatedIdentitiesSet: evaluatedPrincipal.allowedPrincipalFederatedIdentitiesSet,
			allowedPrincipalsAccountsSet:           allowedPrincipalsAccountsSet,
			allowedPrincipalsArnsSet:               allowedPrincipalsArnsSet,
			allowedPrincipalServicesSet:            evaluatedPrincipal.allowedPrincipalServicesSet,
			isPublic:                               isPublic,
			isShared:                               isShared,
			isPrivate:                              isPrivate,
		}
	}

	return evaluatedPrincipal.toEvaluatedCondition()
}

func evaluatePrincipalAccountTypeCondition(evaluatedPrincipal EvaluatedPrincipal, conditionValues []string, evaulatedOperator EvaluatedOperator) EvaluatedCondition {
	processed := false
	allowedPrincipalsAccountsSet := map[string]bool{}
	allowedPrincipalsArnsSet := map[string]bool{}

	isPublic := evaluatedPrincipal.isAwsPublic || evaluatedPrincipal.isServicePublic || evaluatedPrincipal.isFederatedPublic
	isShared := evaluatedPrincipal.isFederatedShared
	isPrivate := false

	for _, conditionValue := range conditionValues {
		if evaulatedOperator.category != "string" {
			continue
		}

		if evaulatedOperator.isLike {
			// Sanity check the input first

			// Regex to allow for account: ["222244446666", "22224444666?", "22224?446666", ...] and must be exactly 12
			reAccountFormat := regexp.MustCompile(`^[0-9\?]{12}$`)

			// Not OK
			if !(strings.Contains(conditionValue, "*") && len(conditionValue) <= 12) && !reAccountFormat.MatchString(conditionValue) {
				continue
			}
		} else {
			// Regex to allow for account: ["222244446666", ...] and must be exactly 12
			reAccountFormat := regexp.MustCompile(`^[0-9]{12}$`)
			// Not OK
			if !reAccountFormat.MatchString(conditionValue) {
				continue
			}
		}

		processed = true

		conditionPolicyValue := MakePolicyValue(conditionValue)

		for principalArns := range evaluatedPrincipal.allowedPrincipalsArnsSet {
			if principalArns == "*" {
				principalPolicyValue := MakePolicyValue("*")
				resolved := principalPolicyValue.Intersection(conditionPolicyValue)

				if resolved == "*" {
					allowedPrincipalsArnsSet[resolved] = true
					isPublic = true
					continue
				}

				allowedPrincipalsAccountsSet[resolved] = true

				if strings.Contains(resolved, "*") || strings.Contains(resolved, "?") {
					isPublic = true
					continue
				}

				isPublic = false
				if resolved != evaluatedPrincipal.userAccountId {
					isShared = true
				} else {
					isPrivate = true
				}
			} else {
				account := extractAccountFromArn(principalArns)
				if !conditionPolicyValue.Contains(account) {
					continue
				}

				allowedPrincipalsArnsSet[principalArns] = true
				if account != evaluatedPrincipal.userAccountId {
					isShared = true
				} else {
					isPrivate = true
				}
			}
		}

		// BUG: Equals Looks like it does wildcards
		for principalAccount := range evaluatedPrincipal.allowedPrincipalsAccountsSet {
			if conditionPolicyValue.Contains(principalAccount) {
				principalPolicyValue := MakePolicyValue(principalAccount)

				resolved := conditionPolicyValue.Intersection(principalPolicyValue)
				allowedPrincipalsAccountsSet[resolved] = true

				if strings.Contains(resolved, "*") || strings.Contains(resolved, "?") {
					isPublic = true
				} else if resolved != evaluatedPrincipal.userAccountId {
					isShared = true
				} else {
					isPrivate = true
				}
			}
		}
	}

	if processed {
		return EvaluatedCondition{
			allowedOrganizationIdsSet:              evaluatedPrincipal.allowedOrganizationIdsSet,
			allowedPrincipalFederatedIdentitiesSet: evaluatedPrincipal.allowedPrincipalFederatedIdentitiesSet,
			allowedPrincipalsAccountsSet:           allowedPrincipalsAccountsSet,
			allowedPrincipalsArnsSet:               allowedPrincipalsArnsSet,
			allowedPrincipalServicesSet:            evaluatedPrincipal.allowedPrincipalServicesSet,
			isPublic:                               isPublic,
			isShared:                               isShared,
			isPrivate:                              isPrivate,
		}
	}

	return evaluatedPrincipal.toEvaluatedCondition()
}

type EvaluatedPrincipal struct {
	allowedOrganizationIdsSet              map[string]bool
	allowedPrincipalFederatedIdentitiesSet map[string]bool
	allowedPrincipalsAccountsSet           map[string]bool
	allowedPrincipalsArnsSet               map[string]bool
	allowedPrincipalServicesSet            map[string]bool
	isAwsPublic                            bool
	isAwsShared                            bool
	isAwsPrivate                           bool
	isFederatedPublic                      bool
	isFederatedShared                      bool
	isServicePublic                        bool
	userAccountId                          string
}

func (evaluatedPrincipal EvaluatedPrincipal) toEvaluatedCondition() EvaluatedCondition {
	return EvaluatedCondition{
		allowedOrganizationIdsSet:              evaluatedPrincipal.allowedOrganizationIdsSet,
		allowedPrincipalFederatedIdentitiesSet: evaluatedPrincipal.allowedPrincipalFederatedIdentitiesSet,
		allowedPrincipalsAccountsSet:           evaluatedPrincipal.allowedPrincipalsAccountsSet,
		allowedPrincipalsArnsSet:               evaluatedPrincipal.allowedPrincipalsArnsSet,
		allowedPrincipalServicesSet:            evaluatedPrincipal.allowedPrincipalServicesSet,
		isPublic:                               evaluatedPrincipal.isAwsPublic || evaluatedPrincipal.isServicePublic || evaluatedPrincipal.isFederatedPublic,
		isShared:                               evaluatedPrincipal.isAwsShared || evaluatedPrincipal.isFederatedShared,
		isPrivate:                              evaluatedPrincipal.isAwsPrivate,
	}
}

func evaluatePrincipal(principal Principal, userAccountId string) (EvaluatedPrincipal, error) {
	evaluatedPrincipal := EvaluatedPrincipal{
		allowedPrincipalFederatedIdentitiesSet: map[string]bool{},
		allowedPrincipalServicesSet:            map[string]bool{},
		allowedPrincipalsAccountsSet:           map[string]bool{},
		allowedPrincipalsArnsSet:               map[string]bool{},
		userAccountId:                          userAccountId,
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
					evaluatedPrincipal.isAwsPublic = true
					evaluatedPrincipal.allowedPrincipalsArnsSet[principalItem] = true
					continue
				}

				if strings.Contains(principalItem, "*") || strings.Contains(principalItem, "?") {
					continue
				}

				var currentSet *map[string]bool

				if reIsAwsAccount.MatchString(principalItem) {
					account = principalItem
					currentSet = &evaluatedPrincipal.allowedPrincipalsAccountsSet
				} else if reIsAwsResource.MatchString(principalItem) {
					account = reIsAwsResource.FindStringSubmatch(principalItem)[1]
					currentSet = &evaluatedPrincipal.allowedPrincipalsArnsSet
				} else {
					return evaluatedPrincipal, fmt.Errorf("unabled to parse arn or account: %s", principalItem)
				}

				if userAccountId != account {
					evaluatedPrincipal.isAwsShared = true
				} else {
					evaluatedPrincipal.isAwsPrivate = true
				}

				(*currentSet)[principalItem] = true
			case "Service":
				if strings.Contains(principalItem, "*") || strings.Contains(principalItem, "?") {
					continue
				}

				evaluatedPrincipal.allowedPrincipalServicesSet[principalItem] = true
				evaluatedPrincipal.isServicePublic = true
			case "Federated":
				if strings.Contains(principalItem, "*") || strings.Contains(principalItem, "?") {
					continue
				}

				if principalItem == "cognito-identity.amazonaws.com" ||
					principalItem == "www.amazon.com" ||
					principalItem == "graph.facebook.com" ||
					principalItem == "accounts.google.com" {
					evaluatedPrincipal.isFederatedPublic = true
				} else {
					evaluatedPrincipal.isFederatedShared = true
				}

				evaluatedPrincipal.allowedPrincipalFederatedIdentitiesSet[principalItem] = true
			}
		}
	}

	return evaluatedPrincipal, nil
}

func mergeSet(set1 map[string]bool, set2 map[string]bool) map[string]bool {
	if set1 == nil {
		if set2 == nil {
			return map[string]bool{}
		}

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

func extractAccountInPlaceFromArn(arn string) string {
	splitArn := strings.Split(arn, ":")

	// There should always be an account
	if len(splitArn) < 6 {
		return ""
	}

	return extractAccount(splitArn[4])
}

func extractAccountFromArn(arn string) string {
	reIsAwsResourceOrWildcard := regexp.MustCompile(`^.*:([0-9\?\*]{1,12}):.*$`)
	reIsAwsResource := regexp.MustCompile(`^arn:[a-z]*:[a-z]*:[a-z]*:([0-9]{12}):.*$`)

	if reIsAwsResourceOrWildcard.MatchString(arn) {
		var accountCheck string
		if reIsAwsResource.MatchString(arn) {
			arnAccount := reIsAwsResource.FindStringSubmatch(arn)
			accountCheck = arnAccount[1]
		} else {
			arnAccount := reIsAwsResourceOrWildcard.FindStringSubmatch(arn)
			accountCheck = arnAccount[1]
		}

		return extractAccount(accountCheck)
	}

	return "*"
}

func extractAccount(account string) string {
	reIsAwsAccountOrWildcard := regexp.MustCompile(`^([0-9\?\*]{1,12})`)
	reIsAwsAccount := regexp.MustCompile(`^[0-9]{12}$`)

	if reIsAwsAccountOrWildcard.MatchString(account) &&
		(reIsAwsAccount.MatchString(account) || strings.Contains(account, "*") || len(account) == 12) {
		return account
	}

	return ""
}
