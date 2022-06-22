package aws

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws/arn"

	"github.com/turbot/go-kit/helpers"
)

var (
	// AWS condition operators to be checked for the trusted access
	conditionOperatorsToCheck = []string{
		"ArnEquals",
		"ArnEqualsIfExists",
		"ArnLike",
		"ArnLikeIfExists",
		"StringEquals",
		"StringEqualsIfExists",
		"StringEqualsIgnoreCase",
		"StringEqualsIgnoreCaseIfExists",
		"StringLike",
		"StringLikeIfExists",
	}
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

type PolicyEvaluation struct {
	// Policy                               Policy   `json:"policy"`
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
	var sourceAccountId string = "444455556666"
	re := regexp.MustCompile(`[0-9]{12}`)

	evaluation := PolicyEvaluation{}

	if policy.Statements == nil {
		return &evaluation, nil
	}

	for index, stmt := range policy.Statements {
		public := stmt.EvaluateStatement(&evaluation)
		if public {
			evaluation.IsPublic = true
			if stmt.Sid == "" {
				evaluation.PublicStatementIds = append(evaluation.PublicStatementIds, fmt.Sprintf("Statement[%s]", strconv.Itoa(index+1)))
			} else {
				evaluation.PublicStatementIds = append(evaluation.PublicStatementIds, stmt.Sid)
			}
		}
	}

	evaluation.AllowedPrincipalAccountIds = StringSliceDistinct(evaluation.AllowedPrincipalAccountIds)
	accountIds := []string{}
	for _, item := range StringSliceDistinct(evaluation.AllowedPrincipals) {
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
	evaluation.AllowedPrincipalAccountIds = accountIds

	// Add all types of principals into allowed principals
	evaluation.AllowedPrincipals = StringSliceDistinct(evaluation.AllowedPrincipals)
	evaluation.AllowedPrincipals = append(evaluation.AllowedPrincipals, evaluation.AllowedPrincipalServices...)
	evaluation.AllowedPrincipals = StringSliceDistinct(append(evaluation.AllowedPrincipals, evaluation.AllowedPrincipalFederatedIdentities...))

	evaluation.AllowedPrincipalFederatedIdentities = StringSliceDistinct(evaluation.AllowedPrincipalFederatedIdentities)
	evaluation.AllowedOrganizationIds = StringSliceDistinct(evaluation.AllowedOrganizationIds)
	evaluation.AllowedPrincipalServices = StringSliceDistinct(evaluation.AllowedPrincipalServices)
	evaluation.PublicAccessLevels = StringSliceDistinct(evaluation.PublicAccessLevels)
	evaluation.PublicStatementIds = StringSliceDistinct(evaluation.PublicStatementIds)

	evaluation.AccessLevel = "shared"
	if evaluation.IsPublic {
		evaluation.AccessLevel = "public"
	} else {
		if (len(evaluation.AllowedPrincipalAccountIds) == 0 && len(evaluation.AllowedOrganizationIds) == 0) || (len(evaluation.AllowedPrincipalAccountIds) > 0 && evaluation.AllowedPrincipalAccountIds[0] == sourceAccountId) {
			evaluation.AccessLevel = "private"
		}
	}

	sort.Strings(evaluation.AllowedOrganizationIds)
	sort.Strings(evaluation.AllowedPrincipalAccountIds)
	sort.Strings(evaluation.AllowedPrincipalFederatedIdentities)
	sort.Strings(evaluation.AllowedPrincipalServices)
	sort.Strings(evaluation.AllowedPrincipals)
	sort.Strings(evaluation.PublicAccessLevels)
	sort.Strings(evaluation.PublicStatementIds)

	return &evaluation, nil
}

func (stmt *Statement) EvaluateStatement(evaluation *PolicyEvaluation) bool {

	stmtEvaluation := PolicyEvaluation{}
	// Check for the deny statements separately
	if stmt.Effect == "Deny" {
		// TODO
		return stmt.DenyStatementEvaluation(evaluation)
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

		for operatorKey, operatorValue := range stmt.Condition {
			hasIfExistsSuffix := CheckIfExistsSuffix(operatorKey)
			// hasForAnyValuePrefix := CheckForAnyValuePrefix(operatorKey)

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
								if !hasIfExistsSuffix {
									internalPublicPrincipalKey = false
								}
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
								if !hasIfExistsSuffix {
									internalPublicPrincipalKey = false
								}
							}
						case "aws:sourceaccount", "aws:sourceowner":
							if !hasIfExistsSuffix {
								internalPublicPrincipalKey = false
							}
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
								} else {
									isPublic = false
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
									fmt.Println("aws:principalorgid", orgs)
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

							isPublic = false
						}
					}

					if stmt.Principal != nil &&
						stmt.Principal["AWS"] == nil &&
						stmt.Principal["Service"] != nil &&
						len(stmt.Principal["Service"].([]string)) != 0 && hasServicePrincipalConditionKey(conditionKey) && !hasIfExistsSuffix {
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

func (stmt *Statement) DenyStatementEvaluation(evaluation *PolicyEvaluation) bool {
	if stmt.NotPrincipal != nil {
		// makes policy unsolvable as it denies access to only principals mentioned in `NotPrincipal` but allows access to everyone else.
		return false
	}
	if stmt.Principal != nil && stmt.Principal["AWS"] != nil && helpers.StringSliceContains((stmt.Principal["AWS"]).([]string), "*") {
		return true
	}
	return false
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
func hasAWSPrincipalConditionKey(conditionKey string) bool {
	return helpers.StringSliceContains(trustedAWSPrincipalConditionKeys, strings.ToLower(conditionKey))
}
func hasServicePrincipalConditionKey(conditionKey string) bool {
	return helpers.StringSliceContains(trustedServicePrincipalConditionKeys, strings.ToLower(conditionKey))
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

// StringSliceDistinct returns a slice with the unique elements the input string slice
func StringSliceDistinct(slice []string) []string {
	var res = []string{}
	countMap := make(map[string]int)
	for _, item := range slice {
		countMap[item]++
	}
	for item := range countMap {
		res = append(res, item)
	}
	return res
}

/* USEFUL DATA - required while coding

[
  {
    "ForAnyValue:StringEquals": {
      "aws:PrincipalOrgPaths": [
        "o-a1b2c3d4e5/r-ab12/ou-ab12-11111111/ou-ab12-22222222/"
      ]
    }
  },
  {
    "ForAnyValue:StringLike": {
      "aws:PrincipalOrgPaths": [
        "o-a1b2c3d4e5/r-ab12/ou-ab12-11111111/ou-ab12-22222222/*"
      ]
    }
  },
  {
    "ForAnyValue:StringLike": {
      "aws:PrincipalOrgPaths": [
        "o-a1b2c3d4e5/r-ab12/ou-ab12-11111111/ou-ab12-22222222/*"
      ]
    }
  },
  {
    "ForAnyValue:StringLike": {
      "aws:PrincipalOrgPaths": ["o-a1b2c3d4e5/*"]
    }
  },
  {
    "ForAnyValue:StringLike": {
      "aws:PrincipalOrgPaths": [
        "o-a1b2c3d4e5/r-ab12/ou-ab12-33333333/*",
        "o-a1b2c3d4e5/r-ab12/ou-ab12-22222222/*"
      ]
    }
  }
]

[
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Action": "dynamodb:GetItem",
        "Resource": "arn:aws:dynamodb:*:*:table/Thread",
        "Condition": {
          "ForAllValues:StringEquals": {
            "dynamodb:Attributes": ["ID", "Message", "Tags"]
          }
        }
      }
    ]
  },
  {
    "Version": "2012-10-17",
    "Statement": {
      "Effect": "Deny",
      "Action": "dynamodb:PutItem",
      "Resource": "arn:aws:dynamodb:*:*:table/Thread",
      "Condition": {
        "ForAnyValue:StringEquals": {
          "dynamodb:Attributes": ["ID", "PostDateTime"]
        }
      }
    }
  }
]


				// log.Println("[INFO] operator key:", operatorKey)

				// TODO with multiple conditions as they behave like an and operator

				"Condition": {
					"StringEquals": {
						"aws:PrincipalAccount": "999988887777"
					},
					"ArnLike": {
						"aws:SourceArn": "arn:aws:cloudwatch:us-east-1:560741234067:alarm:*"
					}
				}

			"Condition": {
				"ForAnyValue:StringEquals": {
					"dynamodb:Attributes": ["ID", "PostDateTime"]
				}
			}

*/

// func listIAMActionFromParliament() {
// 	permissionsData = getParliamentIamPermissions()
// 	for _, service := range permissionsData {
// 		for _, privilege := range service.Privileges {
// 			a := strings.ToLower(service.Prefix + ":" + privilege.Privilege)
// 			awsIamPermissionData{
// 				AccessLevel: privilege.AccessLevel,
// 				Action:      a,
// 				Description: privilege.Description,
// 				Prefix:      service.Prefix,
// 				Privilege:   privilege.Privilege,
// 			}
// 		}
// 	}
// }

// func evaluateCondition(stmt Statement) *ConditionAssessment {

// 	return nil
// }
