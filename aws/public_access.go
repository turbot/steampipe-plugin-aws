package aws

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws/arn"

	"github.com/turbot/go-kit/helpers"
)

var (
	// AWS Global Keys to be checked for the trusted access for AWS Principal
	trustedAWSPrincipalConditionKeys = []string{
		"aws:principalaccount",
		"aws:principalarn",
		"aws:principalorgid",
		"aws:principalorgpaths", //["o-a1b2c3d4e5/*"]  , ["o-a1b2c3d4e5/r-ab12/ou-ab12-11111111/ou-ab12-22222222/ou-*"]
	}
	// AWS Global Keys to be checked for the trusted access for Service Principal
	trustedServicePrincipalConditionKeys = []string{
		"aws:sourcearn",
		"aws:sourceaccount", // SourceAccount is used for giving IAM roles access from an account to the topic.
		"aws:sourceowner",   // SourceOwner is used for giving access to other AWS Services from a specific account
	}

	// https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_condition_operators.html#Conditions_String
	StringConditionalOperators = []string{
		"StringEquals",
		"StringNotEquals",
		"StringEqualsIgnoreCase",
		"StringNotEqualsIgnoreCase",
		"StringLike",
		"StringNotLike",
	}

	// Numeric condition operators
	NumericConditionalOperators = []string{
		"NumericEquals",
		"NumericNotEquals",
		"NumericLessThan",
		"NumericLessThanEquals",
		"NumericGreaterThan",
		"NumericGreaterThanEquals",
	}

	// Date condition operators
	DateConditionalOperators = []string{
		"DateEquals",
		"DateNotEquals",
		"DateLessThan",
		"DateLessThanEquals",
		"DateGreaterThan",
		"DateGreaterThanEquals",
	}

	// Bool condition operators
	BoolConditionalOperators = []string{"Bool"}

	// Binary condition operators
	BinaryConditionalOperators = []string{"BinaryEquals"}

	// IP address condition operators
	IPAddressConditionalOperators = []string{
		"IpAddress",
		"NotIpAddress",
	}

	// ARNConditionalOperators = []string{
	// 	"ArnEquals",
	// 	"ArnLike",
	// 	"ArnNotEquals",
	// 	"ArnNotLike",
	// }

	ARNConditionalOperators = []string{
		// "ArnEquals",
		// "ArnLike",
		// "ArnNotEquals",
		// "ArnNotLike",
	}

	NotOperators = []string{"ArnNotEquals", "ArnNotLike", "StringNotEquals", "StringNotEqualsIgnoreCase", "StringNotLike", "NotIpAddress"}
)

type ConditionMap struct {
	And map[string][]string `type:"and"`
	Not map[string][]string `type:"not"`
	Or  map[string][]string `type:"or"`
}

type ConditionAndPrincipalMap struct {
	Principal struct {
		AWS, Service, Federated []string
	}
	Condition ConditionMap
}

// StringSliceDistinct returns a slice with the unique elements the input string slice
func StringSliceDistinct(slice []string) []string {
	res := []string{}
	countMap := make(map[string]int)
	for _, item := range slice {
		countMap[item]++
		// if this is the first time we have come across this item, add to res
		if countMap[item] == 1 {
			res = append(res, item)
		}
	}
	return res
}

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

func (policy *Policy) EvaluatePolicy() (*PolicyEvaluation, error) {
	//TODO - bring source account information for getting public, private or shared level access info
	re := regexp.MustCompile(`[0-9]{12}`)
	var policyEvaluation PolicyEvaluation
	// var public bool
	if policy.Statements == nil {
		return &policyEvaluation, nil
	}

	actions := []string{}
	allowedAccounts := []string{}
	allowedOrgs := []string{}
	allowedServices := []string{}
	allowedFederatedIdentities := []string{}
	allowedPrincipals := []string{}
	// publicAccessLevels := []string{}
	publicStatementIds := []string{}
	// deniedActions := []string{}
	// deniedAccounts := []string{}
	resourceAccountId := []string{}

	for index, stmt := range policy.Statements {
		// if stmt.Effect == "Allow" {
		// 	actions = append(actions, stmt.Action...)
		// }
		if stmt.Resource != nil {
			for _, item := range stmt.Resource {
				if arn.IsARN(item) {
					awsARN, _ := arn.Parse(item)
					if awsARN.AccountID != "" {
						resourceAccountId = append(resourceAccountId, awsARN.AccountID)
					}
				}
				// TODO - Should we add principals which doesn't have account ids
			}
		}

		public, evaluation := stmt.EvaluateStatement()
		// fmt.Println("evaluation:", evaluation)
		if public {
			policyEvaluation.IsPublic = true
			// publicAccessLevels = append(publicAccessLevels, evaluation.PublicAccessLevels...)
			actions = append(actions, stmt.Action...)
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

	policyEvaluation.AllowedPrincipalAccountIds = StringSliceDistinct(allowedAccounts)
	accountIds := []string{}
	for _, item := range StringSliceDistinct(allowedPrincipals) {
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
	// policyEvaluation.AllowedPrincipals = StringSliceDistinct(policyEvaluation.AllowedPrincipals)
	policyEvaluation.AllowedPrincipals = append(allowedPrincipals, allowedFederatedIdentities...)
	policyEvaluation.AllowedPrincipals = append(policyEvaluation.AllowedPrincipals, allowedOrgs...)
	policyEvaluation.AllowedPrincipals = append(policyEvaluation.AllowedPrincipals, allowedServices...)
	policyEvaluation.AllowedPrincipals = StringSliceDistinct(policyEvaluation.AllowedPrincipals)

	policyEvaluation.AllowedPrincipalFederatedIdentities = StringSliceDistinct(allowedFederatedIdentities)
	policyEvaluation.AllowedOrganizationIds = StringSliceDistinct(allowedOrgs)
	policyEvaluation.AllowedPrincipalServices = StringSliceDistinct(allowedServices)
	// policyEvaluation.PublicActions = StringSliceDistinct(publicAccessLevels)
	policyEvaluation.PublicStatementIds = publicStatementIds

	permissionsData := getParliamentIamPermissions()
	if len(actions) > 0 {
		policyEvaluation.PublicAccessLevels = GetAccessLevelsFromActions(permissionsData, actions)
	}

	resourceAccountId = StringSliceDistinct(resourceAccountId)
	policyEvaluation.AccessLevel = "private"
	if policyEvaluation.IsPublic {
		policyEvaluation.AccessLevel = "public"
	} else {
		if len(resourceAccountId) > 0 {
			for _, item := range policyEvaluation.AllowedPrincipals {
				if arn.IsARN(item) {
					arnParts, _ := arn.Parse(item)
					if arnParts.AccountID != resourceAccountId[0] {
						policyEvaluation.AccessLevel = "shared"
					}
				}
			}
		}
		// if len(policyEvaluation.AllowedPrincipalAccountIds) == 0 && len(policyEvaluation.AllowedOrganizationIds) == 0 {
		// 	policyEvaluation.AccessLevel = "private"
		// } else {
		// 	policyEvaluation.AccessLevel = "shared"
		// }
	}

	sort.Strings(policyEvaluation.AllowedOrganizationIds)
	sort.Strings(policyEvaluation.AllowedPrincipalAccountIds)
	sort.Strings(policyEvaluation.AllowedPrincipalFederatedIdentities)
	sort.Strings(policyEvaluation.AllowedPrincipalServices)
	sort.Strings(policyEvaluation.AllowedPrincipals)
	sort.Strings(policyEvaluation.PublicAccessLevels)
	sort.Strings(policyEvaluation.PublicStatementIds)

	return &policyEvaluation, nil
}

func (stmt *Statement) EvaluateStatement() (bool, PolicyEvaluation) {
	allowedPrincipals := []string{}
	stmtEvaluation := PolicyEvaluation{}
	// Check for the deny statements separately
	if stmt.Effect == "Deny" {
		// TODO
		return stmt.DenyStatementEvaluation(&stmtEvaluation)
		// return false, stmtEvaluation
	}

	// https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_notprincipal.html#specifying-notprincipal-allow
	if stmt.NotPrincipal != nil {
		if data, ok := stmt.NotPrincipal["AWS"]; ok {
			awsNotPrincipals := data.([]string)
			if helpers.StringSliceContains(awsNotPrincipals, "*") {
				return false, stmtEvaluation
			} else {
				fmt.Println("NotPrincipal With Allow")
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

	if helpers.StringSliceContains(awsPrincipals, "*") {
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
			// hasForAnyValue := strings.Contains(operatorKey, "ForAnyValue:")
			operatorKey = strings.ReplaceAll(operatorKey, "ForAnyValue:", "")
			// hasForAllValues := strings.Contains(operatorKey, "ForAllValues:")
			operatorKey = strings.ReplaceAll(operatorKey, "ForAllValues:", "")

			var typeOfOperator string = "Unknown"
			// fmt.Println("operatorKey:", operatorKey)
			switch operatorKey {
			case "StringEquals", "StringNotEquals", "StringEqualsIgnoreCase", "StringNotEqualsIgnoreCase", "StringLike", "StringNotLike":
				typeOfOperator = "String"
				operatorKey = strings.ReplaceAll(operatorKey, "IgnoreCase", "")
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

			// fmt.Println("typeOfOperator:", typeOfOperator)

			hasNotInOperator := CheckNotInOperator(operatorKey)
			hasLikeOperator := strings.Contains(operatorKey, "Like")

			if conditionOperatorValueMap, ok := operatorValue.(map[string]interface{}); ok {
				internalPublicPrincipalKey := true
				for conditionKey, conditionValue := range conditionOperatorValueMap {
					if hasPublicPrincipal {
						// hasPrincipalConditionKey := hasAWSPrincipalConditionKey(conditionKey)
						// hasServiceConditionKey := hasServicePrincipalConditionKey(conditionKey)
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
											// TODO - Normalizing account IDs, e.g., 23* becomes 23?????????, or remains 23*, or something else?
											// AllowedAccountsForPrincipals = append(AllowedAccountsForPrincipals, acctId)
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
										allowedPrincipalsForService = append(allowedPrincipalsForService, conditionValue.([]string)...)
										internalPublicPrincipalKey = false
									}
								} else if hasIfExistsSuffix {
									allowedPrincipalsForService = append(allowedPrincipalsForService, conditionValue.([]string)...)
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
										allowedPrincipalsForService = append(allowedPrincipalsForService, conditionValue.([]string)...)
										allowedPrincipals = helpers.RemoveFromStringSlice(allowedPrincipals, "*")
										internalPublicPrincipalKey = false
									}
								} else if hasIfExistsSuffix {
									allowedPrincipalsForService = append(allowedPrincipalsForService, conditionValue.([]string)...)
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
	}

	stmtEvaluation.AllowedOrganizationIds = append(stmtEvaluation.AllowedOrganizationIds, allowedOrganizationIds...)
	stmtEvaluation.AllowedPrincipals = append(stmtEvaluation.AllowedPrincipals, allowedPrincipals...)
	stmtEvaluation.AllowedPrincipalServices = append(stmtEvaluation.AllowedPrincipalServices, allowedPrincipalsForService...)
	stmtEvaluation.AllowedPrincipalFederatedIdentities = append(stmtEvaluation.AllowedPrincipalFederatedIdentities, federatedPrincipals...)

	if hasPublicConditionPrincipals && isPublic {
		return true, stmtEvaluation
	}
	return false, stmtEvaluation
}

// In a way "Effect" = "Deny" never allows grants but only explicitely denies the rights
func (stmt *Statement) DenyStatementEvaluation(evaluation *PolicyEvaluation) (bool, PolicyEvaluation) {
	if stmt.NotPrincipal != nil {
		// makes policy unsolvable as it denies access to only principals mentioned in `NotPrincipal` but allows access to everyone else.
		return false, *evaluation
	}
	// evaluation.Actions = stmt.Action
	return false, *evaluation
	// TODO: Instead of returning false should return an analysis to negate the allowed actions and principals from other allowed statements in the policy - more useful for the case of per principal analysis
}

/*
https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_multi-value-conditions.html

// ForAllValues – Use with multivalued condition keys. Tests whether the value of every member of the request set is a subset of the condition key set. The condition returns true if every key value in the request matches at least one value in the policy. It also returns true if there are no keys in the request, or if the key values resolve to a null data set, such as an empty string. Do not use ForAllValues with an Allow effect because it can be overly permissive.

ForAnyValue – Use with multivalued condition keys. Tests whether at least one member of the set of request values matches at least one member of the set of condition key values. The condition returns true if any one of the key values in the request matches any one of the condition values in the policy. For no matching key or a null dataset, the condition returns false.

The the difference between single-valued and multivalued condition keys depends on the number of values in the request context, not the number of values in the policy condition.
*/
func CheckForAnyValuePrefix(key string) bool {
	return strings.HasPrefix(key, "ForAnyValue")
}

func CheckForAllValuesPrefix(key string) bool {
	return strings.HasPrefix(key, "ForAllValues")
}
func CheckIfExistsSuffix(key string) bool {
	return strings.HasSuffix(key, "IfExists")
}
func CheckNotInOperator(operator string) bool {
	return strings.Contains(operator, "Not")
}
func hasAWSPrincipalConditionKey(conditionKey string) bool {
	return helpers.StringSliceContains(trustedAWSPrincipalConditionKeys, strings.ToLower(conditionKey))
}
func hasServicePrincipalConditionKey(conditionKey string) bool {
	return helpers.StringSliceContains(trustedServicePrincipalConditionKeys, strings.ToLower(conditionKey))
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
	return StringSliceDistinct(accessLevels)
}

/*
func (stmt *Statement) EvaluateStatement2(evaluation *PolicyEvaluation) bool {

	stmtEvaluation := PolicyEvaluation{}
	// Check for the deny statements separately
	// if stmt.Effect == "Deny" {
	// 	// TODO
	// 	return stmt.DenyStatementEvaluation(evaluation)
	// }

	// https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_notprincipal.html#specifying-notprincipal-allow
	if stmt.NotPrincipal != nil {
		if data, ok := stmt.NotPrincipal["AWS"]; ok {
			awsNotPrincipals := data.([]string)
			if helpers.StringSliceContains(awsNotPrincipals, "*") {
				return false
			} else {
				fmt.Println("NotPrincipal With Allow")
				return true
			}
		}
	}

	// Check for the allowed statement - TODO
	if stmt.NotAction != nil {
		return true
	}

	var awsPrincipals, servicePrincipals, federatedPrincipals []string
	var hasPublicPrincipal = false
	var isPublic = false

	if stmt.Principal != nil {
		if data, ok := stmt.Principal["AWS"]; ok {
			awsPrincipals = data.([]string)
			stmtEvaluation.AllowedPrincipals = awsPrincipals
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

	if helpers.StringSliceContains(awsPrincipals, "*") {
		hasPublicPrincipal = true
		isPublic = true
	}

	// If there is no restriction from the condition side, then policy depends completely on the statement principals section
	hasPublicConditionPrincipals := true

	if stmt.Condition != nil {
		internalPublicPrincipalOperator := true

		// Code to detect public
		for operatorKey, operatorValue := range stmt.Condition {
			hasIfExistsSuffix := CheckIfExistsSuffix(operatorKey)
			operatorKey = strings.ReplaceAll(operatorKey, "IfExists", "")

			hasNotInOperator := CheckNotInOperator(operatorKey)
			operatorKey = strings.ReplaceAll(operatorKey, "Not", "")

			if conditionOperatorValueMap, ok := operatorValue.(map[string]interface{}); ok {
				internalPublicPrincipalKey := true
				for conditionKey, conditionValue := range conditionOperatorValueMap {
					if hasAWSPrincipalConditionKey(conditionKey) {
						switch strings.ToLower(conditionKey) {
						case "aws:principalaccount":
							if !hasIfExistsSuffix {
								internalPublicPrincipalKey = false
							}
						case "aws:principalarn":
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
								internalPublicPrincipalKey = false
							}
						case "aws:principalorgid", "aws:principalorgpaths":
							if !hasIfExistsSuffix {
								internalPublicPrincipalKey = false
							}
						}
					}

					if hasServicePrincipalConditionKey(conditionKey) {
						switch strings.ToLower(conditionKey) {
						case "aws:sourcearn":
							sourceArnPublic := false
							for _, pARN := range conditionValue.([]string) {
								if arn.IsARN(pARN) {
									arnParts, _ := arn.Parse(pARN)
									if arnParts.AccountID == "*" {
										sourceArnPublic = true
									}
								}
							}
							if !sourceArnPublic {
								internalPublicPrincipalKey = false
							}
						case "aws:sourceaccount", "aws:sourceowner":
							if !hasIfExistsSuffix {
								internalPublicPrincipalKey = false
							}
						}
					}

					switch strings.ToLower(conditionKey) {
					case "aws:sourceip":
						if !hasNotInOperator {
							internalPublicPrincipalOperator = false
						}
					}
					if !internalPublicPrincipalKey {
						internalPublicPrincipalOperator = false
						break
					}
				}
			}
			if !internalPublicPrincipalOperator {
				hasPublicConditionPrincipals = false
				break
			}
		}

		// OLD CODE - to collect info
		for operatorKey, operatorValue := range stmt.Condition {
			hasIfExistsSuffix := CheckIfExistsSuffix(operatorKey)
			// hasForAnyValuePrefix := CheckForAnyValuePrefix(operatorKey)

			if conditionOperatorValueMap, ok := operatorValue.(map[string]interface{}); ok {
				// Check if the Principals contain * principals, in that case it is public but if there is a restriction using conditions then it will not remain public
				for conditionKey, conditionValue := range conditionOperatorValueMap {
					if hasPublicPrincipal {
						if hasAWSPrincipalConditionKey(conditionKey) {
							if !hasIfExistsSuffix {
								stmtEvaluation.AllowedPrincipals = helpers.RemoveFromStringSlice(stmtEvaluation.AllowedPrincipals, "*")
								stmtEvaluation.AllowedPrincipals = append(stmtEvaluation.AllowedPrincipals, conditionValue.([]string)...)

								if conditionKey == "aws:principalarn" {
									accounts := []string{}
									for _, pARN := range conditionValue.([]string) {
										if arn.IsARN(pARN) {
											arnParts, _ := arn.Parse(pARN)
											accounts = append(accounts, arnParts.AccountID)
										}
									}
									stmtEvaluation.AllowedPrincipalAccountIds = append(stmtEvaluation.AllowedPrincipalAccountIds, accounts...)
								}
							} else {
								switch conditionKey {
								case "aws:principalaccount":
									stmtEvaluation.AllowedPrincipalAccountIds = append(stmtEvaluation.AllowedPrincipalAccountIds, conditionValue.([]string)...)
									stmtEvaluation.AllowedPrincipals = append(stmtEvaluation.AllowedPrincipals, conditionValue.([]string)...)
								case "aws:principalarn":
									accounts := []string{}
									for _, pARN := range conditionValue.([]string) {
										if arn.IsARN(pARN) {
											arnParts, _ := arn.Parse(pARN)
											accounts = append(accounts, arnParts.AccountID)
										}
									}
									stmtEvaluation.AllowedPrincipalAccountIds = append(stmtEvaluation.AllowedPrincipalAccountIds, accounts...)
								case "aws:principalorgid":
									stmtEvaluation.AllowedOrganizationIds = append(stmtEvaluation.AllowedOrganizationIds, conditionValue.([]string)...)
								case "aws:principalorgpaths":
									orgs := []string{}
									for _, paths := range conditionValue.([]string) {
										orgs = append(orgs, strings.Split(paths, "/")[0])
									}
									stmtEvaluation.AllowedOrganizationIds = append(stmtEvaluation.AllowedOrganizationIds, orgs...)
								}
							}
						}
						if hasServicePrincipalConditionKey(conditionKey) && !hasIfExistsSuffix {
							stmtEvaluation.AllowedPrincipals = helpers.RemoveFromStringSlice(stmtEvaluation.AllowedPrincipals, "*")
							stmtEvaluation.AllowedPrincipals = append(stmtEvaluation.AllowedPrincipals, conditionValue.([]string)...)

							switch strings.ToLower(conditionKey) {
							case "aws:sourcearn":
								for _, pARN := range conditionValue.([]string) {
									if arn.IsARN(pARN) {
										arnParts, _ := arn.Parse(pARN)
										stmtEvaluation.AllowedPrincipalServices = append(stmtEvaluation.AllowedPrincipalServices, fmt.Sprintf("%s.amazonaws.com", arnParts.Service))
										stmtEvaluation.AllowedPrincipalAccountIds = append(stmtEvaluation.AllowedPrincipalAccountIds, arnParts.AccountID)
									}
								}
							case "aws:sourceaccount", "aws:sourceowner":
								stmtEvaluation.AllowedPrincipalAccountIds = append(stmtEvaluation.AllowedPrincipalAccountIds, conditionValue.([]string)...)
							}
						}
					}

					if len(awsPrincipals) == 0 &&
						len(servicePrincipals) != 0 &&
						hasServicePrincipalConditionKey(conditionKey) &&
						!hasIfExistsSuffix {
						stmtEvaluation.AllowedPrincipals = helpers.RemoveFromStringSlice(stmtEvaluation.AllowedPrincipals, "*")
						stmtEvaluation.AllowedPrincipals = append(stmtEvaluation.AllowedPrincipals, conditionValue.([]string)...)
						isPublic = false
						switch strings.ToLower(conditionKey) {
						case "aws:sourcearn":
							for _, pARN := range conditionValue.([]string) {
								if arn.IsARN(pARN) {
									arnParts, _ := arn.Parse(pARN)
									stmtEvaluation.AllowedPrincipalServices = append(stmtEvaluation.AllowedPrincipalServices, fmt.Sprintf("%s.amazonaws.com", arnParts.Service))
									stmtEvaluation.AllowedPrincipalAccountIds = append(stmtEvaluation.AllowedPrincipalAccountIds, arnParts.AccountID)
								}
							}
						case "aws:sourceaccount", "aws:sourceowner":
							stmtEvaluation.AllowedPrincipalAccountIds = append(stmtEvaluation.AllowedPrincipalAccountIds, conditionValue.([]string)...)
						}
					}

					// If the policy have principal org or path to org need to add that in the evaluation
					if helpers.StringSliceContains([]string{"aws:principalorgid", "aws:principalorgpaths"}, strings.ToLower(conditionKey)) {
						if val, ok := conditionValue.([]string); ok {
							stmtEvaluation.AllowedOrganizationIds = val
						}
					}
				}
			}
		}
	}

	evaluation.AllowedOrganizationIds = append(evaluation.AllowedOrganizationIds, stmtEvaluation.AllowedOrganizationIds...)
	evaluation.AllowedPrincipals = append(evaluation.AllowedPrincipals, stmtEvaluation.AllowedPrincipals...)
	evaluation.AllowedPrincipalServices = append(evaluation.AllowedPrincipalServices, stmtEvaluation.AllowedPrincipalServices...)
	evaluation.AllowedPrincipalFederatedIdentities = append(evaluation.AllowedPrincipalFederatedIdentities, stmtEvaluation.AllowedPrincipalFederatedIdentities...)

	if hasPublicConditionPrincipals && isPublic {
		return true
	}
	return false
}
*/
