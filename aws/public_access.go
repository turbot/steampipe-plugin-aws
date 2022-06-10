package aws

import (
	"fmt"
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
		} else {
			// TODO - Should we add principals which doesn't have account ids
			accountIds = append(accountIds, item)
		}
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

	if evaluation.IsPublic {
		evaluation.AccessLevel = "public"
	} else {
		if helpers.StringSliceContains(evaluation.AllowedPrincipals, "*") {
			evaluation.AccessLevel = "shared"
		}
		// else if (evaluation.AllowedPrincipalAccountIds) {

		// }
	}

	return &evaluation, nil
}

func (stmt *Statement) EvaluateStatement(evaluation *PolicyEvaluation) bool {

	stmtEvaluation := PolicyEvaluation{}
	// Check for the deny statements separately
	if stmt.Effect == "Deny" {
		// TODO
		return stmt.DenyStatementEvaluation(evaluation)
	}

	// Check for the allowed statement
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

	var assessment ConditionAssessment
	if stmt.Condition != nil {
		// log.Println("[INFO] AM I REACHING HERE")
		// { "StringEquals": { "aws:PrincipalAccount": "999988887777"}}
		// { "operatorKey": "operatorValue" }
		// { "operatorKey": { "conditionKey": "conditionValue"}}
		// var conditionAssessment := map[string]map[string]interface{}
		assessment.And = []*ConditionAssessment{}
		for operatorKey, operatorValue := range stmt.Condition {

			// intAssessment := ConditionAssessment{Operator: operatorKey}
			// conditionAssessment[operatorKey] = map[string]interface{}
			// hasAnyValuePrefix := CheckForAnyValuePrefix(operatorKey)
			// hasAllValuesPrefix := CheckForAllValuesPrefix(operatorKey)
			hasIfExistsSuffix := CheckIfExistsSuffix(operatorKey)
			// log.Println("[INFO] operator key:", operatorKey)

			// if helpers.StringSliceContains(conditionOperatorsToCheck, key) {

			// }
			// log.Println("[INFO] operator value:", operatorValue)
			// log.Printf("[INFO] operator value type: %T\n", operatorValue)

			if conditiionOperatorValueMap, ok := operatorValue.(map[string]interface{}); ok {

				// log.Println("[INFO] operator key:", operatorKey)

				// TODO with multiple conditions as they behave like an and operator
				/*
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

				for conditionKey, conditionValue := range conditiionOperatorValueMap {
					// if hasAnyValuePrefix {
					// 	intAssessment.Or = []*ConditionAssessment{
					// 		{Value: conditionValue.([]string), Key: conditionKey},
					// 	}
					// }

					// Check if the Principals contain * principals, in that case it is public but if there is a restriction using conditions then it will not remain public
					if hasPublicPrincipal {
						if hasAWSPrincipalConditionKey(conditionKey) && !hasIfExistsSuffix {
							stmtEvaluation.AllowedPrincipals = helpers.RemoveFromStringSlice(stmtEvaluation.AllowedPrincipals, "*")
							stmtEvaluation.AllowedPrincipals = append(stmtEvaluation.AllowedPrincipals, conditionValue.([]string)...)
							isPublic = false
							// conditionAssessment["isPublic"] = false
						}
						if hasServicePrincipalConditionKey(conditionKey) && !hasIfExistsSuffix {
							stmtEvaluation.AllowedPrincipals = helpers.RemoveFromStringSlice(stmtEvaluation.AllowedPrincipals, "*")
							stmtEvaluation.AllowedPrincipals = append(stmtEvaluation.AllowedPrincipals, conditionValue.([]string)...)
							isPublic = false
							// conditionAssessment["isPublic"] = false
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
	return isPublic
}

func (stmt *Statement) DenyStatementEvaluation(evaluation *PolicyEvaluation) bool {
	return false
}

// https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_multi-value-conditions.html
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

type ConditionAssessment struct {
	And      []*ConditionAssessment `type:"and"`
	Not      *ConditionAssessment   `type:"not"`
	Or       []*ConditionAssessment `type:"or"`
	Operator string                 `type:"operator"`
	Key      string                 `type:"key"`
	Value    []string               `type:"value"`
}
