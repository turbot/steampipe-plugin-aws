package aws

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws/arn"

	"github.com/turbot/go-kit/helpers"
)

type PolicyEvaluation struct {
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

func (policy *Policy) EvaluatePolicy(sourceAccountID string) (*PolicyEvaluation, error) {
	//TODO - bring source account information for getting public, private or shared level access info
	re := regexp.MustCompile(`[0-9]{12}`)

	var policyEvaluation PolicyEvaluation
	if !re.Match([]byte(sourceAccountID)) {
		return &policyEvaluation, fmt.Errorf("%s is not a valid. Please enter a valid account id", sourceAccountID)
	}
	// var public bool
	if policy.Statements == nil {
		return &policyEvaluation, nil
	}

	var publicActions, allowedOrgs, allowedAccounts, allowedServices, allowedFederatedIdentities, allowedPrincipals, publicStatementIds []string
	deniedActions := []string{}
	// deniedAccounts := []string{}

	// Only for allow statements
	for index, stmt := range policy.Statements {
		if stmt.Effect == "Deny" {
			continue
		}
		// Check for the deny statements separately
		public, evaluation := stmt.EvaluateStatement()
		if public {
			policyEvaluation.IsPublic = true
			publicActions = append(publicActions, stmt.Action...)
			if stmt.Sid == "" {
				publicStatementIds = append(publicStatementIds, fmt.Sprintf("Statement[%d]", index+1))
			} else {
				publicStatementIds = append(publicStatementIds, stmt.Sid)
			}
		}

		allowedAccounts = append(allowedAccounts, evaluation.AllowedPrincipalAccountIds...)
		allowedOrgs = append(allowedOrgs, evaluation.AllowedOrganizationIds...)
		allowedServices = append(allowedServices, evaluation.AllowedPrincipalServices...)
		allowedFederatedIdentities = append(allowedFederatedIdentities, evaluation.AllowedPrincipalFederatedIdentities...)
		allowedPrincipals = append(allowedPrincipals, evaluation.AllowedPrincipals...)
	}

	// Only for denied statements
	for _, stmt := range policy.Statements {
		if stmt.Effect == "Allow" {
			continue
		}

		deniedEvaluation := stmt.DenyStatementEvaluation()
		if helpers.StringSliceContains(deniedEvaluation.DeniedPrincipals, "*") {
			deniedActions = append(deniedActions, stmt.Action...)
		}
	}

	policyEvaluation.AllowedPrincipalAccountIds = helpers.StringSliceDistinct(allowedAccounts)
	accountIds := []string{}
	for _, item := range helpers.StringSliceDistinct(allowedPrincipals) {
		if arn.IsARN(item) {
			awsARN, _ := arn.Parse(item)
			if awsARN.AccountID != "" {
				accountIds = append(accountIds, awsARN.AccountID)
			}
		} else if item == "*" || re.Match([]byte(item)) {
			accountIds = append(accountIds, item)
		}
		// TODO - Should we add principals which doesn't have account ids
	}
	policyEvaluation.AllowedPrincipalAccountIds = append(allowedAccounts, accountIds...)

	// Add all types of principals into allowed principals
	policyEvaluation.AllowedPrincipals = append(allowedPrincipals, allowedFederatedIdentities...)
	policyEvaluation.AllowedPrincipals = append(policyEvaluation.AllowedPrincipals, allowedServices...)
	policyEvaluation.AllowedPrincipals = helpers.StringSliceDistinct(policyEvaluation.AllowedPrincipals)

	policyEvaluation.AllowedPrincipalFederatedIdentities = helpers.StringSliceDistinct(allowedFederatedIdentities)
	policyEvaluation.AllowedOrganizationIds = helpers.StringSliceDistinct(allowedOrgs)
	policyEvaluation.AllowedPrincipalServices = helpers.StringSliceDistinct(allowedServices)
	policyEvaluation.AllowedPrincipalAccountIds = helpers.StringSliceDistinct(policyEvaluation.AllowedPrincipalAccountIds)
	policyEvaluation.PublicStatementIds = helpers.StringSliceDistinct(publicStatementIds)

	// Action access level will be avaialble only for the policies granting public access
	publicActions = helpers.StringSliceDiff(publicActions, deniedActions)
	if len(publicActions) > 0 {
		permissionsData := getParliamentIamPermissions()
		policyEvaluation.PublicAccessLevels = GetAccessLevelsFromActions(permissionsData, publicActions)
	}

	policyEvaluation.AccessLevel = "private"
	if policyEvaluation.IsPublic {
		policyEvaluation.AccessLevel = "public"
		policyEvaluation.AllowedPrincipalAccountIds = helpers.RemoveFromStringSlice(policyEvaluation.AllowedPrincipalAccountIds, sourceAccountID)
	} else {
		if len(policyEvaluation.AllowedOrganizationIds) > 0 {
			policyEvaluation.AccessLevel = "shared"
		} else if len(policyEvaluation.AllowedPrincipalAccountIds) > 0 {
			for _, item := range policyEvaluation.AllowedPrincipalAccountIds {

				if arn.IsARN(item) {
					arnParts, _ := arn.Parse(item)
					if arnParts.AccountID != sourceAccountID {
						policyEvaluation.AccessLevel = "shared"
					}
				} else if item != sourceAccountID {
					policyEvaluation.AccessLevel = "shared"
				}
			}
		} else if len(policyEvaluation.AllowedPrincipals) > 0 {
			for _, item := range policyEvaluation.AllowedPrincipals {
				if arn.IsARN(item) {
					arnParts, _ := arn.Parse(item)
					if arnParts.AccountID != sourceAccountID {
						policyEvaluation.AccessLevel = "shared"
					}
				}
			}
		}
	}
	policyEvaluation.AllowedPrincipals = helpers.StringSliceDistinct(policyEvaluation.AllowedPrincipals)

	if helpers.StringSliceContains([]string{"private", "shared"}, policyEvaluation.AccessLevel) {
		policyEvaluation.AllowedPrincipalAccountIds = helpers.RemoveFromStringSlice(policyEvaluation.AllowedPrincipalAccountIds, []string{"*", sourceAccountID}...)
	}

	sort.Strings(helpers.StringSliceDistinct(policyEvaluation.AllowedOrganizationIds))
	sort.Strings(helpers.StringSliceDistinct(policyEvaluation.AllowedPrincipalFederatedIdentities))
	sort.Strings(policyEvaluation.AllowedPrincipalAccountIds)
	sort.Strings(policyEvaluation.AllowedPrincipalServices)
	sort.Strings(policyEvaluation.AllowedPrincipals)
	sort.Strings(policyEvaluation.PublicAccessLevels)
	sort.Strings(policyEvaluation.PublicStatementIds)

	return &policyEvaluation, nil
}

func (stmt *Statement) EvaluateStatement() (bool, PolicyEvaluation) {
	allowedPrincipals := []string{}
	stmtEvaluation := PolicyEvaluation{}

	// https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_notprincipal.html#specifying-notprincipal-allow
	if stmt.NotPrincipal != nil {
		if data, ok := stmt.NotPrincipal["AWS"]; ok {
			awsNotPrincipals := data.([]string)
			if helpers.StringSliceContains(awsNotPrincipals, "*") {
				return false, stmtEvaluation
			} else {

				return true, stmtEvaluation
			}
		}
	}

	// Check for the allowed statement - TODO
	// if stmt.NotAction != nil {
	// 	return true
	// }

	var awsPrincipals, servicePrincipals, federatedPrincipals []string
	var hasPublicPrincipal = false
	var isPublic = false

	if stmt.Principal != nil {
		if data, ok := stmt.Principal["AWS"]; ok {
			awsPrincipals = data.([]string)
			allowedPrincipals = awsPrincipals
		}
		if data, ok := stmt.Principal["Service"]; ok {
			servicePrincipals = data.([]string)
			stmtEvaluation.AllowedPrincipalServices = servicePrincipals
		}
		if data, ok := stmt.Principal["Federated"]; ok {
			federatedPrincipals = data.([]string)
			stmtEvaluation.AllowedPrincipalFederatedIdentities = federatedPrincipals
		}
	}

	// if helpers.StringSliceContains(awsPrincipals, "*") {
	if helpers.StringSliceContains(awsPrincipals, "*") || (len(awsPrincipals) == 0 && len(servicePrincipals) > 0) {
		hasPublicPrincipal = true
		isPublic = true
	}

	// If there is no restriction from the condition side, then policy depends completely on the statement principals section
	hasPublicConditionPrincipals := true
	allowedAccountsForPrincipals := []string{}
	allowedOrganizationIds := []string{}
	allowedServicesForPrincipals := []string{}
	allowedPrincipalsForService := []string{}

	if stmt.Condition != nil {
		internalPublicPrincipalOperator := true

		// Code to detect public
		for operatorKey, operatorValue := range stmt.Condition {
			hasIfExistsSuffix := CheckIfExistsSuffix(operatorKey)
			operatorKey = strings.ReplaceAll(operatorKey, "IfExists", "")
			// hasForAnyValue := strings.HasPrefix(operatorKey, "ForAnyValue:")
			// hasForAllValues := strings.HasPrefix(operatorKey, "ForAllValues:")
			operatorKey = strings.ReplaceAll(operatorKey, "ForAnyValue:", "")
			operatorKey = strings.ReplaceAll(operatorKey, "ForAllValues:", "")

			var typeOfOperator string = "Unknown"
			typeOfOperator, operatorKey = getOperatorType(operatorKey)

			hasNotInOperator := CheckNotInOperator(operatorKey)
			hasLikeOperator := strings.Contains(operatorKey, "Like")

			if conditionOperatorValueMap, ok := operatorValue.(map[string]interface{}); ok {
				internalPublicPrincipalKey := true
				for conditionKey, conditionValue := range conditionOperatorValueMap {
					if hasPublicPrincipal || len(awsPrincipals) == 0 {

						if typeOfOperator == "String" {
							switch conditionKey {
							case "aws:principalaccount": // Works with String operators
								if !hasIfExistsSuffix && !hasNotInOperator && !hasLikeOperator {
									allowedPrincipals = helpers.RemoveFromStringSlice(allowedPrincipals, "*")
									allowedAccountsForPrincipals = append(allowedAccountsForPrincipals, conditionValue.([]string)...)
									internalPublicPrincipalKey = false
								} else if hasLikeOperator && !hasNotInOperator && !hasIfExistsSuffix {
									var wildcardAccountIds []string
									for _, acctId := range conditionValue.([]string) {
										if !strings.ContainsAny(acctId, "*?") {
											wildcardAccountIds = append(wildcardAccountIds, acctId)
										} else {
											allowedAccountsForPrincipals = append(allowedAccountsForPrincipals, acctId)
										}
									}
									if len(wildcardAccountIds) == 0 {
										internalPublicPrincipalKey = false
										allowedPrincipals = helpers.RemoveFromStringSlice(allowedPrincipals, "*")
									}
									// else {// TODO - How to add the wildcard account Ids to the list}
								} else if hasIfExistsSuffix {
									allowedAccountsForPrincipals = append(allowedAccountsForPrincipals, conditionValue.([]string)...)
								}
								// else if hasNotInOperator {
								// 	// TODO
								// 	// Shall I add * into AllowedAccounts as it means all accounts other than the accounts mentioned in the condition
								// }
							case "aws:principalorgid", "aws:principalorgpaths":
								if !hasIfExistsSuffix && !hasNotInOperator && !hasLikeOperator {
									allowedOrganizationIds = append(allowedOrganizationIds, conditionValue.([]string)...)
									allowedPrincipals = helpers.RemoveFromStringSlice(allowedPrincipals, "*")
									internalPublicPrincipalKey = false
								} else if hasLikeOperator && !hasNotInOperator && !hasIfExistsSuffix {
									for _, orgIdOrPath := range conditionValue.([]string) {
										// o-[a-z0-9]{10,32} - regex for organization id
										org := strings.Split(orgIdOrPath, "/")[0]
										if !strings.ContainsAny(org, "*?") {
											// Public as if org id is having wildcards(*,?) it will allow a number of organizations which could only be determined after expanding the org id based on the pattern
											allowedOrganizationIds = append(allowedOrganizationIds, org)
										} else {
											internalPublicPrincipalKey = false
											allowedPrincipals = helpers.RemoveFromStringSlice(allowedPrincipals, "*")
											allowedOrganizationIds = append(allowedOrganizationIds, org)
										}
									}
								}
								// else if hasIfExistsSuffix || hasNotInOperator { // means public
								// 	// public for hasIfExistsSuffix and hasNotInOperator because
								// 	// hasIfExistsSuffix - if `aws:principalorgid` or `aws:principalorgpaths` doesn't exists in the request context it evalute to true
								// 	// internalPublicPrincipalKey = false
								// }
							case "aws:principalarn": // Works with both ARN and String operators
								if hasLikeOperator && !hasNotInOperator {
									principalArnPublic := false
									for _, pARN := range conditionValue.([]string) {
										if arn.IsARN(pARN) {
											arnParts, _ := arn.Parse(pARN)
											if arnParts.AccountID == "*" {
												principalArnPublic = true
											}
										}
									}
									if !principalArnPublic {
										allowedPrincipals = append(allowedPrincipals, conditionValue.([]string)...)
										allowedPrincipals = helpers.RemoveFromStringSlice(allowedPrincipals, "*")
										internalPublicPrincipalKey = false
									}
								} else if hasIfExistsSuffix {
									allowedPrincipals = append(allowedPrincipals, conditionValue.([]string)...)
									// This key is included in the request context for all signed requests. Anonymous requests do not include this key.
								} else if hasNotInOperator {
									//TODO What to add in allowed account and allowed principals in this case
								}
							case "aws:sourceaccount", "aws:sourceowner": // Works with String operators
								if !hasIfExistsSuffix && !hasNotInOperator && !hasLikeOperator {
									allowedAccountsForPrincipals = append(allowedAccountsForPrincipals, conditionValue.([]string)...)
									allowedPrincipals = helpers.RemoveFromStringSlice(allowedPrincipals, "*")
									internalPublicPrincipalKey = false
								} else if hasLikeOperator && !hasNotInOperator && !hasIfExistsSuffix {
									var wildcardAccountIds []string
									for _, acctId := range conditionValue.([]string) {
										if !strings.ContainsAny(acctId, "*?") {
											// TODO - Normalizing account IDs, e.g., 23* becomes 23?????????, or remains 23*, or something else?
											// AllowedAccounts = append(AllowedAccounts, acctId)
											wildcardAccountIds = append(wildcardAccountIds, acctId)
										} else {
											allowedServicesForPrincipals = append(allowedServicesForPrincipals, acctId)
										}
									}
									if len(wildcardAccountIds) == 0 {
										internalPublicPrincipalKey = false
										allowedPrincipals = helpers.RemoveFromStringSlice(allowedPrincipals, "*")
									}
									// else {// TODO - How to add the wildcard account Ids to the list}
								} else if hasIfExistsSuffix {
									allowedPrincipals = append(allowedPrincipals, conditionValue.([]string)...)
								} else if hasNotInOperator {
									// TODO
									// Shall I add * into AllowedAccountsForPrincipals as it means all accounts other than the accounts mentioned in the condition
								}
							case "aws:sourcearn": // Works with both ARN and String operators
								if hasLikeOperator && !hasNotInOperator {
									principalArnPublic := false
									for _, pARN := range conditionValue.([]string) {
										if arn.IsARN(pARN) {
											arnParts, _ := arn.Parse(pARN)
											if arnParts.AccountID == "*" {
												principalArnPublic = true
											}
										}
									}
									if !principalArnPublic {
										allowedPrincipals = append(allowedPrincipals, conditionValue.([]string)...)
										internalPublicPrincipalKey = false
									}
								} else if hasIfExistsSuffix {
									allowedPrincipals = append(allowedPrincipals, conditionValue.([]string)...)
									// This key is included in the request context for all signed requests. Anonymous requests do not include this key.
								} else if hasNotInOperator {
									//TODO What to add in allowed account and allowed principals in this case
								}
							}
						} else if typeOfOperator == "Arn" {
							switch conditionKey {
							case "aws:principalarn": // Works with both ARN and String operators
								if !hasNotInOperator {
									principalArnPublic := false
									for _, pARN := range conditionValue.([]string) {
										if arn.IsARN(pARN) {
											arnParts, _ := arn.Parse(pARN)
											if arnParts.AccountID == "*" {
												principalArnPublic = true
											}
										}
									}
									if !principalArnPublic {
										// fmt.Println("AllowedAccountsForPrincipals:", allowedAccountsForPrincipals)
										allowedPrincipals = append(allowedPrincipals, conditionValue.([]string)...)
										allowedPrincipals = helpers.RemoveFromStringSlice(allowedPrincipals, "*")
										internalPublicPrincipalKey = false
									}
								} else if hasIfExistsSuffix {
									allowedPrincipals = append(allowedPrincipals, conditionValue.([]string)...)
									// This key is included in the request context for all signed requests. Anonymous requests do not include this key.
								} else if hasNotInOperator {
									//TODO What to add in allowed account and allowed principals in this case
								}
							case "aws:sourcearn": // Works with both ARN and String operators
								if !hasNotInOperator {
									principalArnPublic := false
									for _, pARN := range conditionValue.([]string) {
										allowedPrincipals = append(allowedPrincipals, pARN)
										if arn.IsARN(pARN) {
											arnParts, _ := arn.Parse(pARN)
											if arnParts.AccountID == "*" {
												principalArnPublic = true
											}
										}
									}
									if !principalArnPublic {
										allowedPrincipals = append(allowedPrincipals, conditionValue.([]string)...)
										allowedPrincipals = helpers.RemoveFromStringSlice(allowedPrincipals, "*")
										internalPublicPrincipalKey = false
									}
								} else if hasIfExistsSuffix {
									allowedPrincipals = append(allowedPrincipals, conditionValue.([]string)...)
									// This key is included in the request context for all signed requests. Anonymous requests do not include this key.
								} else if hasNotInOperator {
									//TODO What to add in allowed account and allowed principals in this case
								}
							}
						} else if typeOfOperator == "IP" {
							if conditionKey == "aws:sourceip" {
								if !hasNotInOperator {
									internalPublicPrincipalOperator = false
								}
							}
						}
					}

					if !internalPublicPrincipalKey {
						internalPublicPrincipalOperator = false
						break
					}
				}
			}
			// Check for ip address
			if !internalPublicPrincipalOperator {
				hasPublicConditionPrincipals = false
				break
			}
		}

		// Code to collect info
		for operatorKey, operatorValue := range stmt.Condition {
			hasIfExistsSuffix := CheckIfExistsSuffix(operatorKey)
			operatorKey = strings.ReplaceAll(operatorKey, "IfExists", "")
			// hasForAnyValue := strings.HasPrefix(operatorKey, "ForAnyValue:")
			operatorKey = strings.ReplaceAll(operatorKey, "ForAnyValue:", "")
			// hasForAllValues := strings.HasPrefix(operatorKey, "ForAllValues:")
			operatorKey = strings.ReplaceAll(operatorKey, "ForAllValues:", "")

			var typeOfOperator string = "Unknown"
			typeOfOperator, operatorKey = getOperatorType(operatorKey)

			hasNotInOperator := CheckNotInOperator(operatorKey)
			hasLikeOperator := strings.Contains(operatorKey, "Like")

			if conditionOperatorValueMap, ok := operatorValue.(map[string]interface{}); ok {
				// Check if the Principals contain * principals, in that case it is public but if there is a restriction using conditions then it will not remain public
				for conditionKey, conditionValue := range conditionOperatorValueMap {
					if hasPublicPrincipal || len(awsPrincipals) == 0 {
						// hasPrincipalConditionKey := hasAWSPrincipalConditionKey(conditionKey)
						// hasServiceConditionKey := hasServicePrincipalConditionKey(conditionKey)
						if typeOfOperator == "String" {
							switch conditionKey {
							case "aws:principalaccount": // Works with String operators
								if !hasNotInOperator && !hasLikeOperator {
									allowedAccountsForPrincipals = append(allowedAccountsForPrincipals, conditionValue.([]string)...)
								}
							case "aws:principalorgid", "aws:principalorgpaths":
								if !hasIfExistsSuffix && !hasNotInOperator && !hasLikeOperator {
									allowedOrganizationIds = append(allowedOrganizationIds, conditionValue.([]string)...)
									allowedPrincipals = helpers.RemoveFromStringSlice(allowedPrincipals, "*")
								} else if hasLikeOperator && !hasNotInOperator && !hasIfExistsSuffix {
									for _, orgIdOrPath := range conditionValue.([]string) {
										// o-[a-z0-9]{10,32} - regex for organization id
										org := strings.Split(orgIdOrPath, "/")[0]
										if !strings.ContainsAny(org, "*?") {
											// Public as if org id is having wildcards(*,?) it will allow a number of organizations which could only be determined after expanding the org id based on the pattern
											allowedOrganizationIds = append(allowedOrganizationIds, org)
										} else {
											allowedPrincipals = helpers.RemoveFromStringSlice(allowedPrincipals, "*")
											allowedOrganizationIds = append(allowedOrganizationIds, org)
										}
									}
								}
							case "aws:principalarn", "aws:sourcearn": // Works with both ARN and String operators
								if (hasLikeOperator && !hasNotInOperator) || hasIfExistsSuffix {
									allowedPrincipals = append(allowedPrincipals, conditionValue.([]string)...)
								}
							case "aws:sourceaccount", "aws:sourceowner": // Works with String operators
								if !hasIfExistsSuffix && !hasNotInOperator && !hasLikeOperator {
									allowedAccountsForPrincipals = append(allowedAccountsForPrincipals, conditionValue.([]string)...)
								} else if hasIfExistsSuffix {
									allowedPrincipals = append(allowedPrincipals, conditionValue.([]string)...)
								}
							}
						} else if typeOfOperator == "Arn" {
							switch conditionKey {
							case "aws:principalarn", "aws:sourcearn":
								if !hasNotInOperator || hasIfExistsSuffix {
									allowedPrincipals = append(allowedPrincipals, conditionValue.([]string)...)
								}
							}
						}
					}
				}
			}
		}
	}

	stmtEvaluation.AllowedOrganizationIds = append(stmtEvaluation.AllowedOrganizationIds, allowedOrganizationIds...)
	stmtEvaluation.AllowedPrincipalAccountIds = append(stmtEvaluation.AllowedPrincipalAccountIds, allowedAccountsForPrincipals...)
	stmtEvaluation.AllowedPrincipals = append(stmtEvaluation.AllowedPrincipals, allowedPrincipals...)
	stmtEvaluation.AllowedPrincipalServices = append(stmtEvaluation.AllowedPrincipalServices, allowedPrincipalsForService...)
	stmtEvaluation.AllowedPrincipalFederatedIdentities = append(stmtEvaluation.AllowedPrincipalFederatedIdentities, federatedPrincipals...)

	if hasPublicConditionPrincipals && isPublic {
		return true, stmtEvaluation
	}
	return false, stmtEvaluation
}

type DeniedStmtEvaluation struct {
	Actions                            []string `json:"actions"`
	DeniedOrganizationIds              []string `json:"denied_organization_ids"`
	DeniedPrincipals                   []string `json:"denied_principals"`
	DeniedPrincipalAccountIds          []string `json:"denied_principal_account_ids"`
	DeniedPrincipalFederatedIdentities []string `json:"denied_principal_federated_identities"`
	DeniedPrincipalServices            []string `json:"denied_principal_services"`
}

// In a way "Effect" = "Deny" never allows grants but only explicitely denies the rights
func (stmt *Statement) DenyStatementEvaluation() DeniedStmtEvaluation {
	deniedEvaluation := DeniedStmtEvaluation{
		Actions: stmt.Action,
	}
	if stmt.NotPrincipal != nil {
		// makes policy unsolvable as it denies access to only principals mentioned in `NotPrincipal` but allows access to everyone else.
		return deniedEvaluation
	}

	var awsPrincipals, deniedPrincipals, deniedOrgPaths, deniedServices, deniedFederatedPrincipal []string
	if stmt.Principal != nil {
		if data, ok := stmt.Principal["AWS"]; ok {
			awsPrincipals = data.([]string)
		}
		if data, ok := stmt.Principal["Service"]; ok {
			deniedServices = data.([]string)
		}
		if data, ok := stmt.Principal["Federated"]; ok {
			deniedFederatedPrincipal = data.([]string)
		}
	}

	deniedEvaluation.DeniedPrincipalServices = deniedServices
	deniedEvaluation.DeniedPrincipalFederatedIdentities = deniedFederatedPrincipal

	hasPublicPrincipal := false

	// Action denied to all principals until there is a explicit condition to limit its impact
	if helpers.StringSliceContains(awsPrincipals, "*") {
		hasPublicPrincipal = true
		deniedPrincipals = []string{"*"}
	}

	if stmt.Condition != nil {

		// Code to detect principals which are denied
		for operatorKey, operatorValue := range stmt.Condition {
			hasIfExistsSuffix := CheckIfExistsSuffix(operatorKey)
			operatorKey = strings.ReplaceAll(operatorKey, "IfExists", "")
			operatorKey = strings.ReplaceAll(operatorKey, "ForAnyValue:", "")
			operatorKey = strings.ReplaceAll(operatorKey, "ForAllValues:", "")

			var typeOfOperator string = "Unknown"
			typeOfOperator, operatorKey = getOperatorType(operatorKey)
			hasNotInOperator := CheckNotInOperator(operatorKey)
			hasLikeOperator := strings.Contains(operatorKey, "Like")

			if conditionOperatorValueMap, ok := operatorValue.(map[string]interface{}); ok {
				for conditionKey, conditionValue := range conditionOperatorValueMap {
					if hasPublicPrincipal || len(awsPrincipals) == 0 {

						if typeOfOperator == "String" {
							switch conditionKey {
							case "aws:principalaccount": // Works with String operators
								if !hasIfExistsSuffix && !hasNotInOperator && !hasLikeOperator {
									deniedPrincipals = helpers.RemoveFromStringSlice(deniedPrincipals, "*")
									deniedPrincipals = append(deniedPrincipals, conditionValue.([]string)...)
								} else if hasLikeOperator && !hasNotInOperator && !hasIfExistsSuffix {
									deniedPrincipals = helpers.RemoveFromStringSlice(deniedPrincipals, "*")
									deniedPrincipals = append(deniedPrincipals, conditionValue.([]string)...)
								}
							case "aws:principalorgid", "aws:principalorgpaths":
								if !hasIfExistsSuffix && !hasNotInOperator && !hasLikeOperator {
									deniedOrgPaths = append(deniedOrgPaths, conditionValue.([]string)...)
									deniedPrincipals = helpers.RemoveFromStringSlice(deniedPrincipals, "*")
								} else if hasLikeOperator && !hasNotInOperator && !hasIfExistsSuffix {
									deniedOrgPaths = append(deniedOrgPaths, conditionValue.([]string)...)
									deniedPrincipals = helpers.RemoveFromStringSlice(deniedPrincipals, "*")
								}
							case "aws:principalarn": // Works with both ARN and String operators
								if hasLikeOperator && !hasNotInOperator {
									deniedPrincipals = helpers.RemoveFromStringSlice(deniedPrincipals, "*")
									deniedPrincipals = append(deniedPrincipals, conditionValue.([]string)...)
								}
							case "aws:sourceaccount", "aws:sourceowner": // Works with String operators
								if !hasIfExistsSuffix && !hasNotInOperator && !hasLikeOperator {
									deniedPrincipals = helpers.RemoveFromStringSlice(deniedPrincipals, "*")
									deniedPrincipals = append(deniedPrincipals, conditionValue.([]string)...)
								} else if hasLikeOperator && !hasNotInOperator && !hasIfExistsSuffix {
									var wildcardAccountIds []string
									for _, acctId := range conditionValue.([]string) {
										if !strings.ContainsAny(acctId, "*?") {
											// TODO - Normalizing account IDs, e.g., 23* becomes 23?????????, or remains 23*, or something else?
											// AllowedAccounts = append(AllowedAccounts, acctId)
											wildcardAccountIds = append(wildcardAccountIds, acctId)
										}
									}
									if len(wildcardAccountIds) == 0 {
										deniedPrincipals = helpers.RemoveFromStringSlice(deniedPrincipals, "*")
										deniedPrincipals = append(deniedPrincipals, conditionValue.([]string)...)
									}
									// else {// TODO - How to add the wildcard account Ids to the list}
								}
							case "aws:sourcearn": // Works with both ARN and String operators
								if hasLikeOperator && !hasNotInOperator {
									principalArnPublic := false
									for _, pARN := range conditionValue.([]string) {
										if arn.IsARN(pARN) {
											arnParts, _ := arn.Parse(pARN)
											if arnParts.AccountID == "*" {
												principalArnPublic = true
											}
										}
									}
									if !principalArnPublic {
										deniedPrincipals = helpers.RemoveFromStringSlice(deniedPrincipals, "*")
										deniedPrincipals = append(deniedPrincipals, conditionValue.([]string)...)
									}
								}
							}
						} else if typeOfOperator == "Arn" {
							switch conditionKey {
							case "aws:principalarn": // Works with both ARN and String operators
								if !hasNotInOperator {
									principalArnPublic := false
									for _, pARN := range conditionValue.([]string) {
										if arn.IsARN(pARN) {
											arnParts, _ := arn.Parse(pARN)
											if arnParts.AccountID == "*" {
												principalArnPublic = true
											}
										}
									}
									if !principalArnPublic {
										deniedPrincipals = helpers.RemoveFromStringSlice(deniedPrincipals, "*")
										deniedPrincipals = append(deniedPrincipals, conditionValue.([]string)...)
									}
								}
							case "aws:sourcearn": // Works with both ARN and String operators
								if !hasNotInOperator {
									principalArnPublic := false
									for _, pARN := range conditionValue.([]string) {
										if arn.IsARN(pARN) {
											arnParts, _ := arn.Parse(pARN)
											if arnParts.AccountID == "*" {
												principalArnPublic = true
											}
										}
									}
									if !principalArnPublic {
										deniedPrincipals = helpers.RemoveFromStringSlice(deniedPrincipals, "*")
										deniedPrincipals = append(deniedPrincipals, conditionValue.([]string)...)
									}
								}
							}
						}
					}
				}
			}
		}
	}
	deniedEvaluation.DeniedPrincipals = deniedPrincipals
	deniedEvaluation.DeniedOrganizationIds = deniedOrgPaths

	return deniedEvaluation
	// TODO: Instead of returning false should return an analysis to negate the allowed actions and principals from other allowed statements in the policy - more useful for the case of per principal analysis
}

/*
https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_multi-value-conditions.html

// ForAllValues – Use with multivalued condition keys. Tests whether the value of every member of the request set is a subset of the condition key set. The condition returns true if every key value in the request matches at least one value in the policy. It also returns true if there are no keys in the request, or if the key values resolve to a null data set, such as an empty string. Do not use ForAllValues with an Allow effect because it can be overly permissive.

ForAnyValue – Use with multivalued condition keys. Tests whether at least one member of the set of request values matches at least one member of the set of condition key values. The condition returns true if any one of the key values in the request matches any one of the condition values in the policy. For no matching key or a null dataset, the condition returns false.

The the difference between single-valued and multivalued condition keys depends on the number of values in the request context, not the number of values in the policy condition.
*/
func CheckIfExistsSuffix(key string) bool {
	return strings.HasSuffix(key, "IfExists")
}
func CheckNotInOperator(operator string) bool {
	return strings.Contains(operator, "Not")
}

func removeNotFromOperator(operatorKey string) string {
	return strings.ReplaceAll(operatorKey, "Not", "")
}

func intersection(s1, s2 []string) (inter []string) {
	hash := make(map[string]bool)
	for _, e := range s1 {
		hash[e] = true
	}
	for _, e := range s2 {
		// If elements present in the hashmap then append intersection list.
		if hash[e] {
			inter = append(inter, e)
		}
	}
	//Remove dups from slice.
	inter = removeDups(inter)
	return
}

//Remove dups from slice.
func removeDups(elements []string) (nodups []string) {
	encountered := make(map[string]bool)
	for _, element := range elements {
		if !encountered[element] {
			nodups = append(nodups, element)
			encountered[element] = true
		}
	}
	return
}

func GetAccessLevelsFromActions(permissionsData ParliamentPermissions, actions []string) []string {
	accessLevels := make([]string, 0)
	if helpers.StringSliceContains(actions, "*") {
		accessLevels = []string{"List", "Permissions management", "Read", "Tagging", "Write"}
	} else {
		for _, action := range actions {
			actionParts := strings.Split(action, ":")
			service := actionParts[0]
			actionPart := ""
			if len(actionParts) == 2 {
				actionPart = actionParts[1]
			}

			re := regexp.MustCompile(strings.ReplaceAll(fmt.Sprintf("^(?i)%s$", actionPart), "*", ".*"))

			for _, parliamentService := range permissionsData {
				if strings.ToLower(service) == strings.ToLower(parliamentService.Prefix) {
					for _, privilege := range parliamentService.Privileges {
						if re.Match([]byte(privilege.Privilege)) {
							accessLevels = append(accessLevels, privilege.AccessLevel)
						}
					}
				}
			}
		}
	}
	return helpers.StringSliceDistinct(accessLevels)
}

func getOperatorType(key string) (operatorType string, newKey string) {
	var typeOfOperator string = "Unknown"
	switch key {
	case "StringEquals", "StringNotEquals", "StringEqualsIgnoreCase", "StringNotEqualsIgnoreCase", "StringLike", "StringNotLike":
		typeOfOperator = "String"
		key = strings.ReplaceAll(key, "IgnoreCase", "")
	case "ArnEquals", "ArnLike", "ArnNotEquals", "ArnNotLike":
		typeOfOperator = "Arn"
	case "NumericEquals", "NumericNotEquals", "NumericLessThan", "NumericLessThanEquals", "NumericGreaterThan", "NumericGreaterThanEquals":
		typeOfOperator = "Numeric"
	case "DateEquals", "DateNotEquals", "DateLessThan", "DateLessThanEquals", "DateGreaterThan", "DateGreaterThanEquals":
		typeOfOperator = "Date"
	case "IpAddress", "NotIpAddress":
		typeOfOperator = "IP"
	case "Bool":
		typeOfOperator = "Bool"
	case "BinaryEquals":
		typeOfOperator = "Binary"
	}
	return typeOfOperator, key
}
