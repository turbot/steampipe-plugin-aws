// TODO: Condition and Principal I am sure breaks
package aws

import (
	"fmt"
	"testing"
)

// / Evaluation Functions
func evaluatePublicAccessLevelsTest(t *testing.T, source PolicySummary, expected PolicySummary) []string {
	errors := []string{}

	countPublicAccessLevels := len(source.PublicAccessLevels)
	expectedCountPublicAccessLevels := len(expected.PublicAccessLevels)
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		errors = append(errors, fmt.Sprintf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels))
	} else {
		for index := range source.PublicAccessLevels {
			currentPublicAccessLevels := source.PublicAccessLevels[index]
			expectedPublicAccessLevels := expected.PublicAccessLevels[index]
			if currentPublicAccessLevels != expectedPublicAccessLevels {
				errors = append(errors, fmt.Sprintf("Unexpected PublicAccessLevels: `%s` PublicAccessLevels expected: `%s`", currentPublicAccessLevels, expectedPublicAccessLevels))
			}
		}
	}

	countSharedAccessLevels := len(source.SharedAccessLevels)
	expectedCountSharedAccessLevels := len(expected.SharedAccessLevels)
	if countSharedAccessLevels != expectedCountSharedAccessLevels {
		errors = append(errors, fmt.Sprintf("Unexpected SharedAccessLevels has: `%d` entries but: `%d` expected", countSharedAccessLevels, expectedCountSharedAccessLevels))
	} else {
		for index := range source.SharedAccessLevels {
			currentSharedAccessLevels := source.SharedAccessLevels[index]
			expectedSharedAccessLevels := expected.SharedAccessLevels[index]
			if currentSharedAccessLevels != expectedSharedAccessLevels {
				errors = append(errors, fmt.Sprintf("Unexpected SharedAccessLevels: `%s` SharedAccessLevels expected: `%s`", currentSharedAccessLevels, expectedSharedAccessLevels))
			}
		}
	}

	countPrivateAccessLevels := len(source.PrivateAccessLevels)
	expectedCountPrivateAccessLevels := len(expected.PrivateAccessLevels)
	if countPrivateAccessLevels != expectedCountPrivateAccessLevels {
		errors = append(errors, fmt.Sprintf("Unexpected PrivateAccessLevels has: `%d` entries but: `%d` expected", countPrivateAccessLevels, expectedCountPrivateAccessLevels))
	} else {
		for index := range source.PrivateAccessLevels {
			currentPrivateAccessLevels := source.PrivateAccessLevels[index]
			expectedPrivateAccessLevels := expected.PrivateAccessLevels[index]
			if currentPrivateAccessLevels != expectedPrivateAccessLevels {
				errors = append(errors, fmt.Sprintf("Unexpected PrivateAccessLevels: `%s` PrivateAccessLevels expected: `%s`", currentPrivateAccessLevels, expectedPrivateAccessLevels))
			}
		}
	}

	return errors
}

func evaluateAccessLevelTest(t *testing.T, source PolicySummary, expected PolicySummary) []string {
	errors := []string{}

	currentAccessLevel := source.AccessLevel
	expectedAccessLevel := expected.AccessLevel
	if currentAccessLevel != expectedAccessLevel {
		errors = append(errors, fmt.Sprintf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel))
	}

	return errors
}

func evaluateSidTest(t *testing.T, source PolicySummary, expected PolicySummary) []string {
	errors := []string{}

	countPublicStatementIds := len(source.PublicStatementIds)
	expectedCountPublicStatementIds := len(expected.PublicStatementIds)
	if countPublicStatementIds != expectedCountPublicStatementIds {
		errors = append(errors, fmt.Sprintf("Unexpected PublicStatementIds has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds))
	} else {
		for index := range source.PublicStatementIds {
			currentPublicStatementIds := source.PublicStatementIds[index]
			expectedPublicStatementIds := expected.PublicStatementIds[index]
			if currentPublicStatementIds != expectedPublicStatementIds {
				errors = append(errors, fmt.Sprintf("Unexpected PublicStatementIds: `%s` PublicStatementIds expected: `%s`", currentPublicStatementIds, expectedPublicStatementIds))
			}
		}
	}

	countSharedStatementIds := len(source.SharedStatementIds)
	expectedCountSharedStatementIds := len(expected.SharedStatementIds)
	if countSharedStatementIds != expectedCountSharedStatementIds {
		errors = append(errors, fmt.Sprintf("Unexpected SharedStatementIds has: `%d` entries but: `%d` expected", countSharedStatementIds, expectedCountSharedStatementIds))
	} else {
		for index := range source.SharedStatementIds {
			currentSharedStatementIds := source.SharedStatementIds[index]
			expectedSharedStatementIds := expected.SharedStatementIds[index]
			if currentSharedStatementIds != expectedSharedStatementIds {
				errors = append(errors, fmt.Sprintf("Unexpected SharedStatementIds: `%s` SharedStatementIds expected: `%s`", currentSharedStatementIds, expectedSharedStatementIds))
			}
		}
	}

	return errors
}

func evaluatePrincipalTest(t *testing.T, source PolicySummary, expected PolicySummary) []string {
	errors := []string{}

	currentIsPublic := source.IsPublic
	expectedIsPublic := expected.IsPublic
	if currentIsPublic != expectedIsPublic {
		errors = append(errors, fmt.Sprintf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic))
	}

	countAllowedPrincipals := len(source.AllowedPrincipals)
	expectedCountAllowedPrincipals := len(expected.AllowedPrincipals)
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		errors = append(errors, fmt.Sprintf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals))
	} else {
		for index := range source.AllowedPrincipals {
			currentAllowedPrincipals := source.AllowedPrincipals[index]
			expectedAllowedPrincipals := expected.AllowedPrincipals[index]
			if currentAllowedPrincipals != expectedAllowedPrincipals {
				errors = append(errors, fmt.Sprintf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals))
			}
		}
	}

	countAllowedPrincipalAccountIds := len(source.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := len(expected.AllowedPrincipalAccountIds)
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		errors = append(errors, fmt.Sprintf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds))
	} else {
		for index := range source.AllowedPrincipalAccountIds {
			currentAllowedPrincipalAccountIds := source.AllowedPrincipalAccountIds[index]
			expectedAllowedPrincipalAccountIds := expected.AllowedPrincipalAccountIds[index]
			if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
				errors = append(errors, fmt.Sprintf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds))
			}
		}
	}

	countAllowedPrincipalFederatedIdentities := len(source.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := len(expected.AllowedPrincipalFederatedIdentities)
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		errors = append(errors, fmt.Sprintf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities))
	} else {
		for index := range source.AllowedPrincipalFederatedIdentities {
			currentAllowedPrincipalFederatedIdentities := source.AllowedPrincipalFederatedIdentities[index]
			expectedAllowedPrincipalFederatedIdentities := expected.AllowedPrincipalFederatedIdentities[index]
			if currentAllowedPrincipalFederatedIdentities != expectedAllowedPrincipalFederatedIdentities {
				errors = append(errors, fmt.Sprintf("Unexpected AllowedPrincipalFederatedIdentities: `%s` AllowedPrincipalFederatedIdentities expected: `%s`", currentAllowedPrincipalFederatedIdentities, expectedAllowedPrincipalFederatedIdentities))
			}
		}
	}

	countAllowedPrincipalServices := len(source.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := len(expected.AllowedPrincipalServices)
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		errors = append(errors, fmt.Sprintf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices))
	} else {
		for index := range source.AllowedPrincipalServices {
			currentAllowedPrincipalServices := source.AllowedPrincipalServices[index]
			expectedAllowedPrincipalServices := expected.AllowedPrincipalServices[index]
			if currentAllowedPrincipalServices != expectedAllowedPrincipalServices {
				errors = append(errors, fmt.Sprintf("Unexpected AllowedPrincipalServices: `%s` AllowedPrincipalServices expected: `%s`", currentAllowedPrincipalServices, expectedAllowedPrincipalServices))
			}
		}
	}

	return errors
}

func evaluateIntegration(t *testing.T, source PolicySummary, expected PolicySummary) []string {
	errors := []string{}
	currentAccessLevel := source.AccessLevel
	expectedAccessLevel := expected.AccessLevel
	if currentAccessLevel != expectedAccessLevel {
		errors = append(errors, fmt.Sprintf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel))
	}

	currentIsPublic := source.IsPublic
	expectedIsPublic := expected.IsPublic
	if currentIsPublic != expectedIsPublic {
		errors = append(errors, fmt.Sprintf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic))
	}

	countAllowedOrganizationIds := len(source.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := len(expected.AllowedOrganizationIds)
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		errors = append(errors, fmt.Sprintf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds))
	} else {
		for index := range source.AllowedOrganizationIds {
			currentAllowedOrganizationIds := source.AllowedOrganizationIds[index]
			expectedAllowedOrganizationIds := expected.AllowedOrganizationIds[index]
			if currentAllowedOrganizationIds != expectedAllowedOrganizationIds {
				errors = append(errors, fmt.Sprintf("Unexpected AllowedOrganizationIds: `%s` AllowedOrganizationIds expected: `%s`", currentAllowedOrganizationIds, expectedAllowedOrganizationIds))
			}
		}
	}

	countAllowedPrincipals := len(source.AllowedPrincipals)
	expectedCountAllowedPrincipals := len(expected.AllowedPrincipals)
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		errors = append(errors, fmt.Sprintf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals))
	} else {
		for index := range source.AllowedPrincipals {
			currentAllowedPrincipals := source.AllowedPrincipals[index]
			expectedAllowedPrincipals := expected.AllowedPrincipals[index]
			if currentAllowedPrincipals != expectedAllowedPrincipals {
				errors = append(errors, fmt.Sprintf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals))
			}
		}
	}

	countAllowedPrincipalAccountIds := len(source.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := len(expected.AllowedPrincipalAccountIds)
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		errors = append(errors, fmt.Sprintf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds))
	} else {
		for index := range source.AllowedPrincipalAccountIds {
			currentAllowedPrincipalAccountIds := source.AllowedPrincipalAccountIds[index]
			expectedAllowedPrincipalAccountIds := expected.AllowedPrincipalAccountIds[index]
			if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
				errors = append(errors, fmt.Sprintf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds))
			}
		}
	}

	countAllowedPrincipalFederatedIdentities := len(source.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := len(expected.AllowedPrincipalFederatedIdentities)
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		errors = append(errors, fmt.Sprintf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities))
	} else {
		for index := range source.AllowedPrincipalFederatedIdentities {
			currentAllowedPrincipalFederatedIdentities := source.AllowedPrincipalFederatedIdentities[index]
			expectedAllowedPrincipalFederatedIdentities := expected.AllowedPrincipalFederatedIdentities[index]
			if currentAllowedPrincipalFederatedIdentities != expectedAllowedPrincipalFederatedIdentities {
				errors = append(errors, fmt.Sprintf("Unexpected AllowedPrincipalFederatedIdentities: `%s` AllowedPrincipalFederatedIdentities expected: `%s`", currentAllowedPrincipalFederatedIdentities, expectedAllowedPrincipalFederatedIdentities))
			}
		}
	}

	countAllowedPrincipalServices := len(source.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := len(expected.AllowedPrincipalServices)
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		errors = append(errors, fmt.Sprintf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices))
	} else {
		for index := range source.AllowedPrincipalServices {
			currentAllowedPrincipalServices := source.AllowedPrincipalServices[index]
			expectedAllowedPrincipalServices := expected.AllowedPrincipalServices[index]
			if currentAllowedPrincipalServices != expectedAllowedPrincipalServices {
				errors = append(errors, fmt.Sprintf("Unexpected AllowedPrincipalServices: `%s` AllowedPrincipalServices expected: `%s`", currentAllowedPrincipalServices, expectedAllowedPrincipalServices))
			}
		}
	}

	countPublicAccessLevels := len(source.PublicAccessLevels)
	expectedCountPublicAccessLevels := len(expected.PublicAccessLevels)
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		errors = append(errors, fmt.Sprintf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels))
	} else {
		for index := range source.PublicAccessLevels {
			currentPublicAccessLevels := source.PublicAccessLevels[index]
			expectedPublicAccessLevels := expected.PublicAccessLevels[index]
			if currentPublicAccessLevels != expectedPublicAccessLevels {
				errors = append(errors, fmt.Sprintf("Unexpected PublicAccessLevels: `%s` PublicAccessLevels expected: `%s`", currentPublicAccessLevels, expectedPublicAccessLevels))
			}
		}
	}

	countSharedAccessLevels := len(source.SharedAccessLevels)
	expectedCountSharedAccessLevels := len(expected.SharedAccessLevels)
	if countSharedAccessLevels != expectedCountSharedAccessLevels {
		errors = append(errors, fmt.Sprintf("Unexpected SharedAccessLevels has: `%d` entries but: `%d` expected", countSharedAccessLevels, expectedCountSharedAccessLevels))
	} else {
		for index := range source.SharedAccessLevels {
			currentSharedAccessLevels := source.SharedAccessLevels[index]
			expectedSharedAccessLevels := expected.SharedAccessLevels[index]
			if currentSharedAccessLevels != expectedSharedAccessLevels {
				errors = append(errors, fmt.Sprintf("Unexpected SharedAccessLevels: `%s` SharedAccessLevels expected: `%s`", currentSharedAccessLevels, expectedSharedAccessLevels))
			}
		}
	}

	countPrivateAccessLevels := len(source.PrivateAccessLevels)
	expectedCountPrivateAccessLevels := len(expected.PrivateAccessLevels)
	if countPrivateAccessLevels != expectedCountPrivateAccessLevels {
		errors = append(errors, fmt.Sprintf("Unexpected PrivateAccessLevels has: `%d` entries but: `%d` expected", countPrivateAccessLevels, expectedCountPrivateAccessLevels))
	} else {
		for index := range source.PrivateAccessLevels {
			currentPrivateAccessLevels := source.PrivateAccessLevels[index]
			expectedPrivateAccessLevels := expected.PrivateAccessLevels[index]
			if currentPrivateAccessLevels != expectedPrivateAccessLevels {
				errors = append(errors, fmt.Sprintf("Unexpected PrivateAccessLevels: `%s` PrivateAccessLevels expected: `%s`", currentPrivateAccessLevels, expectedPrivateAccessLevels))
			}
		}
	}

	countPublicStatementIds := len(source.PublicStatementIds)
	expectedCountPublicStatementIds := len(expected.PublicStatementIds)
	if countPublicStatementIds != expectedCountPublicStatementIds {
		errors = append(errors, fmt.Sprintf("Unexpected PublicStatementIds has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds))
	} else {
		for index := range source.PublicStatementIds {
			currentPublicStatementIds := source.PublicStatementIds[index]
			expectedPublicStatementIds := expected.PublicStatementIds[index]
			if currentPublicStatementIds != expectedPublicStatementIds {
				errors = append(errors, fmt.Sprintf("Unexpected PublicStatementIds: `%s` PublicStatementIds expected: `%s`", currentPublicStatementIds, expectedPublicStatementIds))
			}
		}
	}

	countSharedStatementIds := len(source.SharedStatementIds)
	expectedCountSharedStatementIds := len(expected.SharedStatementIds)
	if countSharedStatementIds != expectedCountSharedStatementIds {
		errors = append(errors, fmt.Sprintf("Unexpected SharedStatementIds has: `%d` entries but: `%d` expected", countSharedStatementIds, expectedCountSharedStatementIds))
	} else {
		for index := range source.SharedStatementIds {
			currentSharedStatementIds := source.SharedStatementIds[index]
			expectedSharedStatementIds := expected.SharedStatementIds[index]
			if currentSharedStatementIds != expectedSharedStatementIds {
				errors = append(errors, fmt.Sprintf("Unexpected SharedStatementIds: `%s` SharedStatementIds expected: `%s`", currentSharedStatementIds, expectedSharedStatementIds))
			}
		}
	}

	return errors
}

/// Test start here

func TestPolicyInstantiation(t *testing.T) {
	t.Run("TestPolicyInstantiationEmptyStringEvaluates", testPolicyInstantiationEmptyStringEvaluates)
	t.Run("TestPolicyInstantiationEmptyJsonStringEvaluates", testPolicyInstantiationEmptyJsonStringEvaluates)
	t.Run("TestPolicyInstantiationWithPrincipalAndWithoutStatementsPolicyEvaluates", testPolicyInstantiationWithPrincipalAndWithoutStatementsPolicyEvaluates)

	t.Run("TestPolicyValidityPolicyWithCorrectVersionEvaluates1", testPolicyValidityPolicyWithCorrectVersionEvaluates1)
	t.Run("TestPolicyValidityPolicyWithCorrectVersionEvaluates2", testPolicyValidityPolicyWithCorrectVersionEvaluates2)
	t.Run("TestPolicyValidityPolicyWithIncorrectVersionFails", testPolicyValidityPolicyWithIncorrectVersionFails)
	t.Run("TestPolicyValidityPolicyWithMissingVersionFails", testPolicyValidityPolicyWithMissingVersionFails)
}

func testPolicyInstantiationEmptyStringEvaluates(t *testing.T) {
	// Set up
	userAccountId := "123456789012"
	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy("", userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Empty policy is not in its expected format")
	}
}

func testPolicyInstantiationEmptyJsonStringEvaluates(t *testing.T) {
	// Set up
	userAccountId := "123456789012"
	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy("{}", userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Empty policy is not in its expected format")
	}
}

func testPolicyInstantiationWithPrincipalAndWithoutStatementsPolicyEvaluates(t *testing.T) {
	// Set up
	userAccountId := "123456789012"
	policyContent := `
	{
	  "Version": "2012-10-17"
	}
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Empty policy is not in its expected format")
	}
}

func TestPolicyVersionElement(t *testing.T) {
	t.Run("TestPolicyValidityPolicyWithCorrectVersionEvaluates1", testPolicyValidityPolicyWithCorrectVersionEvaluates1)
	t.Run("TestPolicyValidityPolicyWithCorrectVersionEvaluates2", testPolicyValidityPolicyWithCorrectVersionEvaluates2)
	t.Run("TestPolicyValidityPolicyWithIncorrectVersionFails", testPolicyValidityPolicyWithIncorrectVersionFails)
	t.Run("TestPolicyValidityPolicyWithMissingVersionFails", testPolicyValidityPolicyWithMissingVersionFails)
}

func testPolicyValidityPolicyWithCorrectVersionEvaluates1(t *testing.T) {
	// Set up
	userAccountId := "123456789012"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": {
            "AWS": "012345678901"
          },
          "Resource": "*"
        },
        {
          "Effect": "Deny",
          "Principal": {
            "AWS": "012345678901"
          },
          "Resource": "*"
        }
      ]
    }
	`
	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPolicyValidityPolicyWithCorrectVersionEvaluates2(t *testing.T) {
	// Set up
	userAccountId := "123456789012"
	policyContent := `
    {
      "Version": "2008-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": {
            "AWS": "012345678901"
          },
          "Resource": "*"
        },
        {
          "Effect": "Deny",
          "Principal": {
            "AWS": "012345678901"
          },
          "Resource": "*"
        }
      ]
    }
	`
	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPolicyValidityPolicyWithIncorrectVersionFails(t *testing.T) {
	// Set up
	userAccountId := "123456789012"
	policyContent := `
    {
      "Version": "2000-10-10",
      "Statement": [
        {
          "Effect": "Allow"
        },
        {
          "Effect": "Deny"
        }
      ]
    }
	`

	// Test
	_, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err == nil {
		t.Fatal("Expected error but no error was returned from EvaluatePolicy")
	}

	expectedErrorMsg := "unsupported value for policy element Version: '2000-10-10' - values supported are '2012-10-17' or '2008-10-17'"

	if errorMsg := err.Error(); errorMsg != expectedErrorMsg {
		t.Fatalf("The error message returned: '%s' but was expected to be: '%s'", errorMsg, expectedErrorMsg)
	}
}

func testPolicyValidityPolicyWithMissingVersionFails(t *testing.T) {
	// Set up
	userAccountId := "123456789012"
	policyContent := `
    {
      "Statement": [
        {
          "Effect": "Allow"
        },
        {
          "Effect": "Deny"
        }
      ]
    }
	`

	// Test
	_, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err == nil {
		t.Fatal("policy element Version is missing")
	}

	expectedErrorMsg := "policy element Version is missing"

	if errorMsg := err.Error(); errorMsg != expectedErrorMsg {
		t.Fatalf("The error message returned: '%s' but was expected to be: '%s'", errorMsg, expectedErrorMsg)
	}
}

// ABOVE
func TestPolicyEffectElement(t *testing.T) {
	t.Run("TestEffectElementWithValidValues", testEffectElementWithValidValues)
	t.Run("TestIfEffectElementWhenValueAllowHasWrongCasingFails", testIfEffectElementWhenValueAllowHasWrongCasingFails)
	t.Run("TestIfEffectElementWhenValueDenyHasWrongCasingFails", testIfEffectElementWhenValueDenyHasWrongCasingFails)
	t.Run("TestIfEffectElementWhenValueIsUnknownFails", testIfEffectElementWhenValueIsUnknownFails)
	t.Run("TestIfEffectElementIsMissingFails", testIfEffectElementIsMissingFails)
}

func testEffectElementWithValidValues(t *testing.T) {
	// Set up
	userAccountId := "123456789012"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": {
            "AWS": "012345678901"
          },
          "Resource": "*"
        },
        {
          "Effect": "Deny",
          "Principal": {
            "AWS": "012345678901"
          },
          "Resource": "*"
        }
      ]
    }
	`
	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testIfEffectElementWhenValueAllowHasWrongCasingFails(t *testing.T) {
	// Set up
	userAccountId := "123456789012"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "allow"
        }
      ]
    }
	`

	// Test
	_, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err == nil {
		t.Fatal("Expected error but no error was returned from EvaluatePolicy")
	}

	expectedErrorMsg := "unsupported value for policy element Effect: 'allow' - values supported are 'Allow' or 'Deny'"

	if errorMsg := err.Error(); errorMsg != expectedErrorMsg {
		t.Fatalf("The error message returned: '%s' but was expected to be: '%s'", errorMsg, expectedErrorMsg)
	}
}

func testIfEffectElementWhenValueDenyHasWrongCasingFails(t *testing.T) {
	// Set up
	userAccountId := "123456789012"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "deny"
        }
      ]
    }
	`

	// Test
	_, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err == nil {
		t.Fatal("Expected error but no error was returned from EvaluatePolicy")
	}

	expectedErrorMsg := "unsupported value for policy element Effect: 'deny' - values supported are 'Allow' or 'Deny'"

	if errorMsg := err.Error(); errorMsg != expectedErrorMsg {
		t.Fatalf("The error message returned: '%s' but was expected to be: '%s'", errorMsg, expectedErrorMsg)
	}
}

func testIfEffectElementWhenValueIsUnknownFails(t *testing.T) {
	// Set up
	userAccountId := "123456789012"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "random"
        }
      ]
    }
	`

	// Test
	_, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err == nil {
		t.Fatal("Expected error but no error was returned from EvaluatePolicy")
	}

	expectedErrorMsg := "unsupported value for policy element Effect: 'random' - values supported are 'Allow' or 'Deny'"

	if errorMsg := err.Error(); errorMsg != expectedErrorMsg {
		t.Fatalf("The error message returned: '%s' but was expected to be: '%s'", errorMsg, expectedErrorMsg)
	}
}

func testIfEffectElementIsMissingFails(t *testing.T) {
	// Set up
	userAccountId := "123456789012"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
        }
      ]
    }
	`

	// Test
	_, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err == nil {
		t.Fatal("Expected error but no error was returned from EvaluatePolicy")
	}

	expectedErrorMsg := "policy element Effect is missing"

	if errorMsg := err.Error(); errorMsg != expectedErrorMsg {
		t.Fatalf("The error message returned: '%s' but was expected to be: '%s'", errorMsg, expectedErrorMsg)
	}
}

func TestSourceAccountId(t *testing.T) {
	t.Run("TestIfSourceAccountIdContainsNonNumericalValuesItFails", testIfSourceAccountIdContainsNonNumericalValuesItFails)
	t.Run("TestIfSourceAccountIdContainsTooFewNumericalValuesItFails", testIfSourceAccountIdContainsTooFewNumericalValuesItFails)
	t.Run("TestIfSourceAccountIdContainsTooManyNumericalValuesItFails", testIfSourceAccountIdContainsTooManyNumericalValuesItFails)
	t.Run("TestIfSourceAccountIdContainsCorrectAmountOfNumericalValuesItEvaluates", testIfSourceAccountIdContainsCorrectAmountOfNumericalValuesItEvaluates)
	t.Run("TestIfSourceAccountIdContainsCorrectAmountOfNumericalValuesAndStartsWithZeroItEvaluates", testIfSourceAccountIdContainsCorrectAmountOfNumericalValuesAndStartsWithZeroItEvaluates)
}

func testIfSourceAccountIdContainsNonNumericalValuesItFails(t *testing.T) {
	// Set up
	userAccountId := "123A123123"

	policyContent := `
	{
	  "Version": "2012-10-17"
	}
	`

	// Test
	_, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err == nil {
		t.Fatal("Expected error but no error was returned from EvaluatePolicy")
	}

	expectedErrorMsg := "source account id is invalid: 123A123123"

	if errorMsg := err.Error(); errorMsg != expectedErrorMsg {
		t.Fatalf("The error message returned: '%s' but was expected to be: '%s'", errorMsg, expectedErrorMsg)
	}
}

func testIfSourceAccountIdContainsTooFewNumericalValuesItFails(t *testing.T) {
	// Set up
	userAccountId := "01234567890"

	policyContent := `
	{
	  "Version": "2012-10-17"
	}
	`

	// Test
	_, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err == nil {
		t.Fatal("Expected error but no error was returned from EvaluatePolicy")
	}

	expectedErrorMsg := "source account id is invalid: 01234567890"

	if errorMsg := err.Error(); errorMsg != expectedErrorMsg {
		t.Fatalf("The error message returned: '%s' but was expected to be: '%s'", errorMsg, expectedErrorMsg)
	}
}

func testIfSourceAccountIdContainsTooManyNumericalValuesItFails(t *testing.T) {
	// Set up
	userAccountId := "012345678901234"

	policyContent := `
	{
	  "Version": "2012-10-17"
	}
	`

	// Test
	_, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err == nil {
		t.Fatal("Expected error but no error was returned from EvaluatePolicy")
	}

	expectedErrorMsg := "source account id is invalid: 012345678901234"

	if errorMsg := err.Error(); errorMsg != expectedErrorMsg {
		t.Fatalf("The error message returned: '%s' but was expected to be: '%s'", errorMsg, expectedErrorMsg)
	}
}

func testIfSourceAccountIdContainsCorrectAmountOfNumericalValuesItEvaluates(t *testing.T) {
	// Set up
	userAccountId := "123456789012"
	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy("", userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testIfSourceAccountIdContainsCorrectAmountOfNumericalValuesAndStartsWithZeroItEvaluates(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy("", userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func TestPolicyPrincipalElementWildcard(t *testing.T) {
	t.Run("TestWhenPrincipalIsWildcarded", testWhenPrincipalIsWildcarded)
	t.Run("TestWhenAwsPrincipalIsWildcarded", testWhenAwsPrincipalIsWildcarded)
	t.Run("TestWhenAwsPrincipalIsWildcardedAndUnknownCondition", testWhenAwsPrincipalIsWildcardedAndUnknownCondition)
	t.Run("TestWhenPrincipalIsMultipleMixedAccountsWithWildcard", testWhenPrincipalIsMultipleMixedAccountsWithWildcard)
	t.Run("TestWhenStatementHasBothPublicAndSharedAccountThenTheEvaluationIsPublic", testWhenStatementHasBothPublicAndSharedAccountThenTheEvaluationIsPublic)
	t.Run("TestWhenStatementHasBothPublicAndPrivateAccountThenTheEvaluationIsPublic", testWhenStatementHasBothPublicAndPrivateAccountThenTheEvaluationIsPublic)
}

func testWhenPrincipalIsWildcarded(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": "012345678901"
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"Write"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenAwsPrincipalIsWildcarded(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": "*",
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"Write"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenAwsPrincipalIsWildcardedAndUnknownCondition(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": "*",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "sts:ExternalId": "shdf89awphfgiohigui;sahyf09238[478r["
            },
            "Bool": {
              "aws:MultiFactorAuthPresent": "true"
            }
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"Write"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsMultipleMixedAccountsWithWildcard(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["444455554444", "arn:aws:iam::444455554444:root", "*", "012345678901", "arn:aws:iam::012345678901:root"]
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "public",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"*",
			"012345678901",
			"444455554444",
			"arn:aws:iam::012345678901:root",
			"arn:aws:iam::444455554444:root",
		},
		AllowedPrincipalAccountIds: []string{
			"*",
			"012345678901",
			"444455554444",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"Write"},
		SharedAccessLevels:                  []string{"Write"},
		PrivateAccessLevels:                 []string{"Write"},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenStatementHasBothPublicAndSharedAccountThenTheEvaluationIsPublic(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": "*",
          "Resource": "*"
        },
        {
			"Effect": "Allow",
			"Action": "sts:AssumeRole",
			"Principal": "222244446666",
			"Resource": "*"
		  } 
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "public",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"*",
			"222244446666",
		},
		AllowedPrincipalAccountIds: []string{
			"*",
			"222244446666",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"Write"},
		SharedAccessLevels:                  []string{"Write"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{"Statement[2]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenStatementHasBothPublicAndPrivateAccountThenTheEvaluationIsPublic(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": "*",
          "Resource": "*"
        },
        {
			"Effect": "Allow",
			"Action": "sts:AssumeRole",
			"Principal": "012345678901",
			"Resource": "*"
		  } 
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "public",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"*",
			"012345678901",
		},
		AllowedPrincipalAccountIds: []string{
			"*",
			"012345678901",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"Write"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"Write"},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func TestPolicyPrincipalElementAccounts(t *testing.T) {
	// NOTE: Changed to silently failed
	// t.Run("TestWhenPrincipalIsAMisformedAccountWithOneDigitShortFails", testWhenPrincipalIsAMisformedAccountWithOneDigitShortFails)
	// t.Run("TestWhenPrincipalIsAMisformedAccountWithOneDigitExtraFails", testWhenPrincipalIsAMisformedAccountWithOneDigitExtraFails)
	t.Run("TestWhenPrincipalIsAUserAccountId", testWhenPrincipalIsAUserAccountId)
	t.Run("TestWhenPrincipalIsAUserAccountArn", testWhenPrincipalIsAUserAccountArn)
	t.Run("TestWhenPrincipalIsACrossAccountId", testWhenPrincipalIsACrossAccountId)
	t.Run("TestWhenPrincipalIsACrossAccountArn", testWhenPrincipalIsACrossAccountArn)
	t.Run("TestWhenPrincipalIsMultipleUserAccounts", testWhenPrincipalIsMultipleUserAccounts)
	t.Run("TestWhenPrincipalIsMultipleCrossAccountsInAscendingOrder", testWhenPrincipalIsMultipleCrossAccountsInAscendingOrder)
	t.Run("TestWhenPrincipalIsMultipleCrossAccountsInDescendingOrder", testWhenPrincipalIsMultipleCrossAccountsInDescendingOrder)
	t.Run("TestWhenPrincipalIsMultipleAccountsPrincipalsAcrossMultipleStatements", testWhenPrincipalIsMultipleAccountsPrincipalsAcrossMultipleStatements)
	t.Run("TestWhenPrincipalIsMultipleMixedAccounts", testWhenPrincipalIsMultipleMixedAccounts)
	t.Run("TestWhenPrincipalHasAWildcardInAccountThenIgnorePrincipal", testWhenPrincipalHasAWildcardInAccountThenIgnorePrincipal)
}

func testWhenPrincipalIsAUserAccountId(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": "*"
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"Write"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsAUserAccountArn(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": "arn:aws:iam::012345678901:root"
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::012345678901:root"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"Write"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsACrossAccountId(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": "444455554444"
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"444455554444"},
		AllowedPrincipalAccountIds:          []string{"444455554444"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"Write"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsACrossAccountArn(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": "arn:aws:iam::444455554444:root"
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::444455554444:root"},
		AllowedPrincipalAccountIds:          []string{"444455554444"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"Write"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsMultipleUserAccounts(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["012345678901", "arn:aws:iam::012345678901:root"]
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "private",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"012345678901",
			"arn:aws:iam::012345678901:root",
		},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"Write"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsMultipleCrossAccountsInAscendingOrder(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["444455554444", "arn:aws:iam::555544445555:root"]
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"444455554444",
			"arn:aws:iam::555544445555:root",
		},
		AllowedPrincipalAccountIds: []string{
			"444455554444",
			"555544445555",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"Write"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsMultipleCrossAccountsInDescendingOrder(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:iam::555544445555:root", "444455554444"]
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"444455554444",
			"arn:aws:iam::555544445555:root",
		},
		AllowedPrincipalAccountIds: []string{
			"444455554444",
			"555544445555",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"Write"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsMultipleAccountsPrincipalsAcrossMultipleStatements(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:iam::444455554444:root"]
          },
          "Resource": "*"
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:iam::012345678901:root"]
          },
          "Resource": "*"
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:iam::444455554444:root"]
          },
          "Resource": "*"
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:iam::012345678901:root"]
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:iam::012345678901:root",
			"arn:aws:iam::444455554444:root",
		},
		AllowedPrincipalAccountIds: []string{
			"012345678901",
			"444455554444",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"Write"},
		PrivateAccessLevels:                 []string{"Write"},
		PublicStatementIds:                  []string{},
		SharedStatementIds: []string{
			"Statement[1]",
			"Statement[3]",
		},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsMultipleMixedAccounts(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["444455554444", "arn:aws:iam::444455554444:root", "012345678901", "arn:aws:iam::012345678901:root"]
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"012345678901",
			"444455554444",
			"arn:aws:iam::012345678901:root",
			"arn:aws:iam::444455554444:root",
		},
		AllowedPrincipalAccountIds: []string{
			"012345678901",
			"444455554444",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"Write"},
		PrivateAccessLevels:                 []string{"Write"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalHasAWildcardInAccountThenIgnorePrincipal(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["44445555*", "arn:aws:iam::444455554444:*", "01234567????", "arn:aws:iam::012345678901:ro??"]
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func TestPolicyPrincipalElementArn(t *testing.T) {
	// NOTE: Changed to silently failed
	// t.Run("TestWhenPrincipalIsAMisformedArnFails", testWhenPrincipalIsAMisformedArnFails)
	t.Run("TestWhenPrincipalIsAUserAccountRole", testWhenPrincipalIsAUserAccountRole)
	t.Run("TestWhenPrincipalIsACrossAccountRole", testWhenPrincipalIsACrossAccountRole)
	t.Run("TestWhenPrincipalIsMultipleUserAccountRoles", testWhenPrincipalIsMultipleUserAccountRoles)
	t.Run("TestWhenPrincipalIsMultipleCrossAccountRolesInAscendingOrder", testWhenPrincipalIsMultipleCrossAccountRolesInAscendingOrder)
	t.Run("TestWhenPrincipalIsMultipleCrossAccountRolesInDescendingOrder", testWhenPrincipalIsMultipleCrossAccountRolesInDescendingOrder)
	t.Run("TestWhenPrincipalIsMultipleAccountRolePrincipalsAcrossMultipleStatements", testWhenPrincipalIsMultipleAccountRolePrincipalsAcrossMultipleStatements)
	t.Run("TestWhenPrincipalIsMultipleMixedAccountRoles", testWhenPrincipalIsMultipleMixedAccountRoles)
}

// NOTE: Changed to silently failed
// func testWhenPrincipalIsAMisformedArnFails(t *testing.T) {
// 	// Set up
// 	userAccountId := "012345678901"
// 	policyContent := `
//     {
//       "Version": "2012-10-17",
//       "Statement": [
//         {
//           "Effect": "Allow",
//           "Action": "sts:AssumeRole",
//           "Principal": {
//             "AWS": "arn:aws:sts::misformed:012345678901:assumed-role/role-name/role-session-name"
//           },
//           "Resource": "*"
//         }
//       ]
//     }
// 	`

// 	// Test
// 	_, err := EvaluatePolicy(policyContent, userAccountId)

// 	// Evaluate
// 	if err == nil {
// 		t.Fatal("Expected error but no error was returned from EvaluatePolicy")
// 	}

// 	expectedErrorMsg := "unabled to parse arn or account: arn:aws:sts::misformed:012345678901:assumed-role/role-name/role-session-name"

// 	if errorMsg := err.Error(); errorMsg != expectedErrorMsg {
// 		t.Fatalf("The error message returned: '%s' but was expected to be: '%s'", errorMsg, expectedErrorMsg)
// 	}
// }

func testWhenPrincipalIsAUserAccountRole(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": "arn:aws:iam::012345678901:role/role-name"
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::012345678901:role/role-name"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"Write"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsACrossAccountRole(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": "arn:aws:iam::444455554444:role/role-name"
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::444455554444:role/role-name"},
		AllowedPrincipalAccountIds:          []string{"444455554444"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"Write"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsMultipleUserAccountRoles(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": [
              "arn:aws:iam::012345678901:role/role-name-1", 
              "arn:aws:iam::012345678901:role/role-name-2"
            ]
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "private",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:iam::012345678901:role/role-name-1",
			"arn:aws:iam::012345678901:role/role-name-2",
		},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"Write"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsMultipleCrossAccountRolesInAscendingOrder(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:iam::444455554444:role/role-name", "arn:aws:iam::555544445555:role/role-name"]
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:iam::444455554444:role/role-name",
			"arn:aws:iam::555544445555:role/role-name",
		},
		AllowedPrincipalAccountIds: []string{
			"444455554444",
			"555544445555",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"Write"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsMultipleCrossAccountRolesInDescendingOrder(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:iam::555544445555:role/role-name", "arn:aws:iam::444455554444:role/role-name"]
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:iam::444455554444:role/role-name",
			"arn:aws:iam::555544445555:role/role-name",
		},
		AllowedPrincipalAccountIds: []string{
			"444455554444",
			"555544445555",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"Write"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsMultipleAccountRolePrincipalsAcrossMultipleStatements(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:iam::444455554444:role/role-name"]
          },
          "Resource": "*"
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:iam::012345678901:role/role-name"]
          },
          "Resource": "*"
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:iam::444455554444:role/role-name"]
          },
          "Resource": "*"
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:iam::012345678901:role/role-name"]
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:iam::012345678901:role/role-name",
			"arn:aws:iam::444455554444:role/role-name",
		},
		AllowedPrincipalAccountIds: []string{
			"012345678901",
			"444455554444",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"Write"},
		PrivateAccessLevels:                 []string{"Write"},
		PublicStatementIds:                  []string{},
		SharedStatementIds: []string{
			"Statement[1]",
			"Statement[3]",
		},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsMultipleMixedAccountRoles(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": [
              "arn:aws:iam::444455554444:role/role-name-2",
              "arn:aws:iam::444455554444:role/role-name-1",
              "arn:aws:iam::012345678901:role/role-name-2",
              "arn:aws:iam::012345678901:role/role-name-1"
            ]
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:iam::012345678901:role/role-name-1",
			"arn:aws:iam::012345678901:role/role-name-2",
			"arn:aws:iam::444455554444:role/role-name-1",
			"arn:aws:iam::444455554444:role/role-name-2",
		},
		AllowedPrincipalAccountIds: []string{
			"012345678901",
			"444455554444",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"Write"},
		PrivateAccessLevels:                 []string{"Write"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func TestPolicyPrincipalElementAssumedRole(t *testing.T) {
	t.Run("TestWhenPrincipalIsAUserAccountAssumedRole", testWhenPrincipalIsAUserAccountAssumedRole)
	t.Run("TestWhenPrincipalIsACrossAccountAssumedRole", testWhenPrincipalIsACrossAccountAssumedRole)
	t.Run("TestWhenPrincipalIsMultipleUserAccountAssumedRoles", testWhenPrincipalIsMultipleUserAccountAssumedRoles)
	t.Run("TestWhenPrincipalIsMultipleCrossAccountAssumedRolesInAscendingOrder", testWhenPrincipalIsMultipleCrossAccountAssumedRolesInAscendingOrder)
	t.Run("TestWhenPrincipalIsMultipleCrossAccountAssumedRolesInDescendingOrder", testWhenPrincipalIsMultipleCrossAccountAssumedRolesInDescendingOrder)
	t.Run("TestWhenPrincipalIsMultipleAccountAssumedRolePrincipalsAcrossMultipleStatements", testWhenPrincipalIsMultipleAccountAssumedRolePrincipalsAcrossMultipleStatements)
	t.Run("TestWhenPrincipalIsMultipleMixedAccountAssumedRoles", testWhenPrincipalIsMultipleMixedAccountAssumedRoles)
}

func testWhenPrincipalIsAUserAccountAssumedRole(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": "arn:aws:sts::012345678901:assumed-role/role-name/role-session-name"
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:sts::012345678901:assumed-role/role-name/role-session-name"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"Write"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsACrossAccountAssumedRole(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": "arn:aws:sts::444455554444:assumed-role/role-name/role-session-name"
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:sts::444455554444:assumed-role/role-name/role-session-name"},
		AllowedPrincipalAccountIds:          []string{"444455554444"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"Write"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsMultipleUserAccountAssumedRoles(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": [
              "arn:aws:sts::012345678901:assumed-role/role-name/role-session-name-1",
              "arn:aws:sts::012345678901:assumed-role/role-name/role-session-name-2"
            ]
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "private",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:sts::012345678901:assumed-role/role-name/role-session-name-1",
			"arn:aws:sts::012345678901:assumed-role/role-name/role-session-name-2",
		},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"Write"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsMultipleCrossAccountAssumedRolesInAscendingOrder(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": [
              "arn:aws:sts::444455554444:assumed-role/role-name/role-session-name",
              "arn:aws:sts::555544445555:assumed-role/role-name/role-session-name"
            ]
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:sts::444455554444:assumed-role/role-name/role-session-name",
			"arn:aws:sts::555544445555:assumed-role/role-name/role-session-name",
		},
		AllowedPrincipalAccountIds: []string{
			"444455554444",
			"555544445555",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"Write"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsMultipleCrossAccountAssumedRolesInDescendingOrder(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:sts::555544445555:assumed-role/role-name/role-session-name", "arn:aws:sts::444455554444:assumed-role/role-name/role-session-name"]
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:sts::444455554444:assumed-role/role-name/role-session-name",
			"arn:aws:sts::555544445555:assumed-role/role-name/role-session-name",
		},
		AllowedPrincipalAccountIds: []string{
			"444455554444",
			"555544445555",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"Write"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsMultipleAccountAssumedRolePrincipalsAcrossMultipleStatements(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:sts::444455554444:assumed-role/role-name/role-session-name"]
          },
          "Resource": "*"
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:sts::012345678901:assumed-role/role-name/role-session-name"]
          },
          "Resource": "*"
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:sts::444455554444:assumed-role/role-name/role-session-name"]
          },
          "Resource": "*"
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:sts::012345678901:assumed-role/role-name/role-session-name"]
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:sts::012345678901:assumed-role/role-name/role-session-name",
			"arn:aws:sts::444455554444:assumed-role/role-name/role-session-name",
		},
		AllowedPrincipalAccountIds: []string{
			"012345678901",
			"444455554444",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"Write"},
		PrivateAccessLevels:                 []string{"Write"},
		PublicStatementIds:                  []string{},
		SharedStatementIds: []string{
			"Statement[1]",
			"Statement[3]",
		},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsMultipleMixedAccountAssumedRoles(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": [
              "arn:aws:sts::444455554444:assumed-role/role-name/role-session-name-2",
              "arn:aws:sts::444455554444:assumed-role/role-name/role-session-name-1",
              "arn:aws:sts::012345678901:assumed-role/role-name/role-session-name-2",
              "arn:aws:sts::012345678901:assumed-role/role-name/role-session-name-1"
            ]
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:sts::012345678901:assumed-role/role-name/role-session-name-1",
			"arn:aws:sts::012345678901:assumed-role/role-name/role-session-name-2",
			"arn:aws:sts::444455554444:assumed-role/role-name/role-session-name-1",
			"arn:aws:sts::444455554444:assumed-role/role-name/role-session-name-2",
		},
		AllowedPrincipalAccountIds: []string{
			"012345678901",
			"444455554444",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"Write"},
		PrivateAccessLevels:                 []string{"Write"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func TestPolicyPrincipalElementFederated(t *testing.T) {
	t.Run("TestWhenPrincipalIsAFederatedUser", testWhenPrincipalIsAFederatedUser)
	t.Run("TestWhenPrincipalIsMulitpleFederatedUserInAscendingOrder", testWhenPrincipalIsMulitpleFederatedUserInAscendingOrder)
	t.Run("TestWhenPrincipalIsMulitpleFederatedUserInDescendingOrder", testWhenPrincipalIsMulitpleFederatedUserInDescendingOrder)
	t.Run("TestWhenPrincipalIsMultipleFederatedUserPrincipalsAcrossMultipleStatements", testWhenPrincipalIsMultipleFederatedUserPrincipalsAcrossMultipleStatements)
	t.Run("TestWhenPrincipalHasAWildcardInFederatedUserThenIgnorePrincipal", testWhenPrincipalHasAWildcardInFederatedUserThenIgnorePrincipal)

	t.Run("TestWhenPrincipalIsASamlUser", testWhenPrincipalIsASamlUser)
	t.Run("TestWhenPrincipalIsMulitpleSamlUserInAscendingOrder", testWhenPrincipalIsMulitpleSamlUserInAscendingOrder)
	t.Run("TestWhenPrincipalIsMulitpleSamlUserInDescendingOrder", testWhenPrincipalIsMulitpleSamlUserInDescendingOrder)
	t.Run("TestWhenPrincipalIsMultipleSamlUserPrincipalsAcrossMultipleStatements", testWhenPrincipalIsMultipleSamlUserPrincipalsAcrossMultipleStatements)
	t.Run("TestWhenPrincipalHasAWildcardInSamlUserThenIgnorePrincipal", testWhenPrincipalHasAWildcardInSamlUserThenIgnorePrincipal)

	t.Run("testWhenPrincipalIsMulitpleSamlProvidersWithAudienceConditions", testWhenPrincipalIsMulitpleSamlProvidersWithAudienceConditions)
	t.Run("testWhenPrincipalIsMulitpleSamlProvidersWithoutAudienceConditions", testWhenPrincipalIsMulitpleSamlProvidersWithoutAudienceConditions)
}

func testWhenPrincipalIsAFederatedUser(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Federated": "cognito-identity.amazonaws.com"
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{"cognito-identity.amazonaws.com"},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"Write"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsMulitpleFederatedUserInAscendingOrder(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	// NOTE: This case is not legit in practice as you cannot have two open providers
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Federated": ["accounts.google.com", "graph.facebook.com"]
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                "shared",
		AllowedOrganizationIds:     []string{},
		AllowedPrincipals:          []string{},
		AllowedPrincipalAccountIds: []string{},
		AllowedPrincipalFederatedIdentities: []string{
			"accounts.google.com",
			"graph.facebook.com",
		},
		AllowedPrincipalServices: []string{},
		IsPublic:                 false,
		PublicAccessLevels:       []string{},
		SharedAccessLevels:       []string{"Write"},
		PrivateAccessLevels:      []string{},
		PublicStatementIds:       []string{},
		SharedStatementIds:       []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsMulitpleFederatedUserInDescendingOrder(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	// NOTE: This case is not legit in practice as you cannot have two open providers
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Federated": ["graph.facebook.com", "accounts.google.com"]
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                "shared",
		AllowedOrganizationIds:     []string{},
		AllowedPrincipals:          []string{},
		AllowedPrincipalAccountIds: []string{},
		AllowedPrincipalFederatedIdentities: []string{
			"accounts.google.com",
			"graph.facebook.com",
		},
		AllowedPrincipalServices: []string{},
		IsPublic:                 false,
		PublicAccessLevels:       []string{},
		SharedAccessLevels:       []string{"Write"},
		PrivateAccessLevels:      []string{},
		PublicStatementIds:       []string{},
		SharedStatementIds:       []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsMultipleFederatedUserPrincipalsAcrossMultipleStatements(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	// NOTE: This case is not legit in practice as you cannot have two open providers
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Federated": ["graph.facebook.com"]
          },
          "Resource": "*"
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Federated": ["accounts.google.com"]
          },
          "Resource": "*"
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Federated": ["graph.facebook.com"]
          },
          "Resource": "*"
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Federated": ["accounts.google.com"]
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                "shared",
		AllowedOrganizationIds:     []string{},
		AllowedPrincipals:          []string{},
		AllowedPrincipalAccountIds: []string{},
		AllowedPrincipalFederatedIdentities: []string{
			"accounts.google.com",
			"graph.facebook.com",
		},
		AllowedPrincipalServices: []string{},
		IsPublic:                 false,
		PublicAccessLevels:       []string{},
		SharedAccessLevels:       []string{"Write"},
		PrivateAccessLevels:      []string{},
		PublicStatementIds:       []string{},
		SharedStatementIds: []string{
			"Statement[1]",
			"Statement[2]",
			"Statement[3]",
			"Statement[4]",
		},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalHasAWildcardInFederatedUserThenIgnorePrincipal(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Federated": ["graph.facebook.*"]
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsASamlUser(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Federated": "arn:aws:iam::111122223333:saml-provider-1/provider-name"
          },
          "Condition": { "StringEquals": { "SAML:aud": ["test"] } },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{"arn:aws:iam::111122223333:saml-provider-1/provider-name"},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"Write"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsMulitpleSamlUserInAscendingOrder(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Federated": ["arn:aws:iam::111122223333:saml-provider-1/provider-name", "arn:aws:iam::111122223333:saml-provider-2/provider-name"]
          },
          "Condition": { "StringEquals": { "SAML:aud": ["test"] } },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                "shared",
		AllowedOrganizationIds:     []string{},
		AllowedPrincipals:          []string{},
		AllowedPrincipalAccountIds: []string{},
		AllowedPrincipalFederatedIdentities: []string{
			"arn:aws:iam::111122223333:saml-provider-1/provider-name",
			"arn:aws:iam::111122223333:saml-provider-2/provider-name",
		},
		AllowedPrincipalServices: []string{},
		IsPublic:                 false,
		PublicAccessLevels:       []string{},
		SharedAccessLevels:       []string{"Write"},
		PrivateAccessLevels:      []string{},
		PublicStatementIds:       []string{},
		SharedStatementIds:       []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsMulitpleSamlUserInDescendingOrder(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Federated": ["arn:aws:iam::111122223333:saml-provider-1/provider-name", "arn:aws:iam::111122223333:saml-provider-2/provider-name"]
          },
          "Condition": { "StringEquals": { "SAML:aud": ["test"] } },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                "shared",
		AllowedOrganizationIds:     []string{},
		AllowedPrincipals:          []string{},
		AllowedPrincipalAccountIds: []string{},
		AllowedPrincipalFederatedIdentities: []string{
			"arn:aws:iam::111122223333:saml-provider-1/provider-name",
			"arn:aws:iam::111122223333:saml-provider-2/provider-name",
		},
		AllowedPrincipalServices: []string{},
		IsPublic:                 false,
		PublicAccessLevels:       []string{},
		SharedAccessLevels:       []string{"Write"},
		PrivateAccessLevels:      []string{},
		PublicStatementIds:       []string{},
		SharedStatementIds:       []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsMultipleSamlUserPrincipalsAcrossMultipleStatements(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Federated": ["arn:aws:iam::111122223333:saml-provider-1/provider-name"]
          },
          "Condition": { "StringEquals": { "SAML:aud": ["test"] } },
          "Resource": "*"
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Federated": ["arn:aws:iam::111122223333:saml-provider-2/provider-name"]
          },
          "Condition": { "StringEquals": { "SAML:aud": ["test"] } },
          "Resource": "*"
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Federated": ["arn:aws:iam::111122223333:saml-provider-1/provider-name"]
          },
          "Condition": { "StringEquals": { "SAML:aud": ["test"] } },
          "Resource": "*"
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Federated": ["arn:aws:iam::111122223333:saml-provider-2/provider-name"]
          },
          "Condition": { "StringEquals": { "SAML:aud": ["test"] } },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                "shared",
		AllowedOrganizationIds:     []string{},
		AllowedPrincipals:          []string{},
		AllowedPrincipalAccountIds: []string{},
		AllowedPrincipalFederatedIdentities: []string{
			"arn:aws:iam::111122223333:saml-provider-1/provider-name",
			"arn:aws:iam::111122223333:saml-provider-2/provider-name",
		},
		AllowedPrincipalServices: []string{},
		IsPublic:                 false,
		PublicAccessLevels:       []string{},
		SharedAccessLevels:       []string{"Write"},
		PrivateAccessLevels:      []string{},
		PublicStatementIds:       []string{},
		SharedStatementIds: []string{
			"Statement[1]",
			"Statement[2]",
			"Statement[3]",
			"Statement[4]",
		},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalHasAWildcardInSamlUserThenIgnorePrincipal(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Federated": ["arn:aws:iam::111122223333:saml-provider-1/*"]
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsMulitpleSamlProvidersWithAudienceConditions(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRoleWithSAML",
          "Principal": { "Federated": "arn:aws:iam::111122223333:saml-provider-1/provider-name" },
          "Condition": { "StringEquals": { "SAML:aud": ["test"] } }
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRoleWithSAML",
          "Principal": { "Federated": "arn:aws:iam::111122223333:saml-provider-2/provider-name" },
          "Condition": { "StringEquals": { "SAML:iss": ["test"] } }
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRoleWithSAML",
          "Principal": { "Federated": "arn:aws:iam::111122223333:saml-provider-3/provider-name" },
          "Condition": { "StringEquals": { "SAML:sub": ["test"] } }
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRoleWithSAML",
          "Principal": { "Federated": "arn:aws:iam::111122223333:saml-provider-4/provider-name" },
          "Condition": { "StringEquals": { "SAML:sub_type": ["test"] } }
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRoleWithSAML",
          "Principal": { "Federated": "arn:aws:iam::111122223333:saml-provider-5/provider-name" },
          "Condition": { "StringEquals": { "SAML:eduPersonOrgDN": ["test"] } }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                "shared",
		AllowedOrganizationIds:     []string{},
		AllowedPrincipals:          []string{},
		AllowedPrincipalAccountIds: []string{},
		AllowedPrincipalFederatedIdentities: []string{
			"arn:aws:iam::111122223333:saml-provider-1/provider-name",
			"arn:aws:iam::111122223333:saml-provider-2/provider-name",
			"arn:aws:iam::111122223333:saml-provider-3/provider-name",
			"arn:aws:iam::111122223333:saml-provider-4/provider-name",
			"arn:aws:iam::111122223333:saml-provider-5/provider-name",
		},
		AllowedPrincipalServices: []string{},
		IsPublic:                 false,
		PublicAccessLevels:       []string{},
		SharedAccessLevels:       []string{"Write"},
		PrivateAccessLevels:      []string{},
		PublicStatementIds:       []string{},
		SharedStatementIds: []string{
			"Statement[1]",
			"Statement[2]",
			"Statement[3]",
			"Statement[4]",
			"Statement[5]",
		},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsMulitpleSamlProvidersWithoutAudienceConditions(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRoleWithSAML",
          "Principal": { "Federated": "arn:aws:iam::111122223333:saml-provider-1/provider-name" }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                "public",
		AllowedOrganizationIds:     []string{},
		AllowedPrincipals:          []string{},
		AllowedPrincipalAccountIds: []string{},
		AllowedPrincipalFederatedIdentities: []string{
			"arn:aws:iam::111122223333:saml-provider-1/provider-name",
		},
		AllowedPrincipalServices: []string{},
		IsPublic:                 true,
		PublicAccessLevels:       []string{"Write"},
		SharedAccessLevels:       []string{},
		PrivateAccessLevels:      []string{},
		PublicStatementIds:       []string{"Statement[1]"},
		SharedStatementIds:       []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func TestPolicyPrincipalElementService(t *testing.T) {
	t.Run("TestWhenPrincipalIsAService", testWhenPrincipalIsAService)
	t.Run("TestWhenPrincipalIsMulitpleServicesInAscendingOrder", testWhenPrincipalIsMulitpleServicesInAscendingOrder)
	t.Run("TestWhenPrincipalIsMulitpleServicesInDescendingOrder", testWhenPrincipalIsMulitpleServicesInDescendingOrder)
	t.Run("TestWhenPrincipalIsMultipleServicePrincipalsAcrossMultipleStatements", testWhenPrincipalIsMultipleServicePrincipalsAcrossMultipleStatements)
	t.Run("TestWhenPrincipalHasAWildcardInServicePrincipalsThenIgnorePrincipal", testWhenPrincipalHasAWildcardInServicePrincipalsThenIgnorePrincipal)
}

func testWhenPrincipalIsAService(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Service": "ec2.amazonaws.com"
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ec2.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"Write"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsMulitpleServicesInAscendingOrder(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Service": ["ecs.amazonaws.com", "elasticloadbalancing.amazonaws.com"]
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices: []string{
			"ecs.amazonaws.com",
			"elasticloadbalancing.amazonaws.com",
		},
		IsPublic:            true,
		PublicAccessLevels:  []string{"Write"},
		SharedAccessLevels:  []string{},
		PrivateAccessLevels: []string{},
		PublicStatementIds:  []string{"Statement[1]"},
		SharedStatementIds:  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsMulitpleServicesInDescendingOrder(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Service": ["elasticloadbalancing.amazonaws.com", "ecs.amazonaws.com"]
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices: []string{
			"ecs.amazonaws.com",
			"elasticloadbalancing.amazonaws.com",
		},
		IsPublic:            true,
		PublicAccessLevels:  []string{"Write"},
		SharedAccessLevels:  []string{},
		PrivateAccessLevels: []string{},
		PublicStatementIds:  []string{"Statement[1]"},
		SharedStatementIds:  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsMultipleServicePrincipalsAcrossMultipleStatements(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Service": ["elasticloadbalancing.amazonaws.com"]
          },
          "Resource": "*"
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Service": ["ecs.amazonaws.com"]
          },
          "Resource": "*"
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Service": ["elasticloadbalancing.amazonaws.com"]
          },
          "Resource": "*"
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Service": ["ecs.amazonaws.com"]
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices: []string{
			"ecs.amazonaws.com",
			"elasticloadbalancing.amazonaws.com",
		},
		IsPublic:            true,
		PublicAccessLevels:  []string{"Write"},
		SharedAccessLevels:  []string{},
		PrivateAccessLevels: []string{},
		PublicStatementIds: []string{
			"Statement[1]",
			"Statement[2]",
			"Statement[3]",
			"Statement[4]",
		},
		SharedStatementIds: []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalHasAWildcardInServicePrincipalsThenIgnorePrincipal(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Service": ["ecs.amazonaws.*"]
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func TestPolicyPrincipalElementMultipleTypes(t *testing.T) {
	t.Run("TestWhenPrincipalIsMultipleTypes", testWhenPrincipalIsMultipleTypes)
	t.Run("TestWhenPrincipalIsMultipleTypesWithWildcard", testWhenPrincipalIsMultipleTypesWithWildcard)
	t.Run("TestWhenPrincipalIsMultipleTypesAcrossMultipleStatements", testWhenPrincipalIsMultipleTypesAcrossMultipleStatements)
	t.Run("TestWhenPrincipalIsMultipleTypesAcrossMultipleStatementsWithWildcard", testWhenPrincipalIsMultipleTypesAcrossMultipleStatementsWithWildcard)
}

func testWhenPrincipalIsMultipleTypes(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Service": "ecs.amazonaws.com",
            "AWS": ["arn:aws:iam::444455554444:root", "arn:aws:iam::012345678901:root"],
            "Federated": "cognito-identity.amazonaws.com"
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "public",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:iam::012345678901:root",
			"arn:aws:iam::444455554444:root",
		},
		AllowedPrincipalAccountIds: []string{
			"*",
			"012345678901",
			"444455554444",
		},
		AllowedPrincipalFederatedIdentities: []string{"cognito-identity.amazonaws.com"},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"Write"},
		SharedAccessLevels:                  []string{"Write"},
		PrivateAccessLevels:                 []string{"Write"},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsMultipleTypesWithWildcard(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Service": "ecs.amazonaws.com",
            "AWS": ["arn:aws:iam::444455554444:root", "*", "arn:aws:iam::012345678901:root"],
            "Federated": "cognito-identity.amazonaws.com"
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "public",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"*",
			"arn:aws:iam::012345678901:root",
			"arn:aws:iam::444455554444:root",
		},
		AllowedPrincipalAccountIds: []string{
			"*",
			"012345678901",
			"444455554444",
		},
		AllowedPrincipalFederatedIdentities: []string{"cognito-identity.amazonaws.com"},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"Write"},
		SharedAccessLevels:                  []string{"Write"},
		PrivateAccessLevels:                 []string{"Write"},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsMultipleTypesAcrossMultipleStatements(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:iam::444455554444:root", "arn:aws:iam::012345678901:root"],
            "Service": ["dynamodb.amazonaws.com"],
            "Federated": ["graph.facebook.com"]
          },
          "Resource": "*"
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:iam::555544445555:root"],
            "Service": ["elasticloadbalancing.amazonaws.com", "ecs.amazonaws.com"],
            "Federated": ["accounts.google.com"]
          },
          "Resource": "*"
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:iam::555544445555:root"],
            "Service": ["ecs.amazonaws.com"],
            "Federated": ["graph.facebook.com", "cognito-identity.amazonaws.com"]
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "public",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:iam::012345678901:root",
			"arn:aws:iam::444455554444:root",
			"arn:aws:iam::555544445555:root",
		},
		AllowedPrincipalAccountIds: []string{
			"*",
			"012345678901",
			"444455554444",
			"555544445555",
		},
		AllowedPrincipalFederatedIdentities: []string{
			"accounts.google.com",
			"cognito-identity.amazonaws.com",
			"graph.facebook.com",
		},
		AllowedPrincipalServices: []string{
			"dynamodb.amazonaws.com",
			"ecs.amazonaws.com",
			"elasticloadbalancing.amazonaws.com",
		},
		IsPublic:            true,
		PublicAccessLevels:  []string{"Write"},
		SharedAccessLevels:  []string{"Write"},
		PrivateAccessLevels: []string{"Write"},
		PublicStatementIds: []string{
			"Statement[1]",
			"Statement[2]",
			"Statement[3]",
		},
		SharedStatementIds: []string{
			"Statement[1]",
			"Statement[2]",
			"Statement[3]",
		},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testWhenPrincipalIsMultipleTypesAcrossMultipleStatementsWithWildcard(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:iam::444455554444:root", "arn:aws:iam::012345678901:root"],
            "Service": ["dynamodb.amazonaws.com"],
            "Federated": ["graph.facebook.com"]
          },
          "Resource": "*"
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:iam::555544445555:root"],
            "Service": ["elasticloadbalancing.amazonaws.com", "ecs.amazonaws.com"],
            "Federated": ["accounts.google.com"]
          },
          "Resource": "*"
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:iam::555544445555:root", "*"],
            "Service": ["ecs.amazonaws.com"],
            "Federated": ["graph.facebook.com", "cognito-identity.amazonaws.com"]
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "public",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"*",
			"arn:aws:iam::012345678901:root",
			"arn:aws:iam::444455554444:root",
			"arn:aws:iam::555544445555:root",
		},
		AllowedPrincipalAccountIds: []string{
			"*",
			"012345678901",
			"444455554444",
			"555544445555",
		},
		AllowedPrincipalFederatedIdentities: []string{
			"accounts.google.com",
			"cognito-identity.amazonaws.com",
			"graph.facebook.com",
		},
		AllowedPrincipalServices: []string{
			"dynamodb.amazonaws.com",
			"ecs.amazonaws.com",
			"elasticloadbalancing.amazonaws.com",
		},
		IsPublic:            true,
		PublicAccessLevels:  []string{"Write"},
		SharedAccessLevels:  []string{"Write"},
		PrivateAccessLevels: []string{"Write"},
		PublicStatementIds: []string{
			"Statement[1]",
			"Statement[2]",
			"Statement[3]",
		},
		SharedStatementIds: []string{
			"Statement[1]",
			"Statement[2]",
			"Statement[3]",
		},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func TestPolicyPrincipalElement(t *testing.T) {
	t.Run("TestPolicyPrincipalElementPrincipalMissingFails", testPolicyPrincipalElementPrincipalMissingFails)
	t.Run("TestPolicyPrincipalElementPrincipalPresent", testPolicyPrincipalElementPrincipalPresent)
}

func testPolicyPrincipalElementPrincipalMissingFails(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*"
        }
      ]
    }
	`

	// Test
	_, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err == nil {
		t.Fatal("Expected error but no error was returned from EvaluatePolicy")
	}

	expectedErrorMsg := "policy element Principal is missing"

	if errorMsg := err.Error(); errorMsg != expectedErrorMsg {
		t.Fatalf("The error message returned: '%s' but was expected to be: '%s'", errorMsg, expectedErrorMsg)
	}
}

func testPolicyPrincipalElementPrincipalPresent(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": "012345678901"
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"Write"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func TestPolicyResourceElement(t *testing.T) {
	//t.Run("TestPolicyResourceElementResourceMissingFails", testPolicyResourceElementResourceMissingFails)
	t.Run("TestPolicyResourceElementResourcePresent", testPolicyResourceElementResourcePresent)
}

// NOTE: Incorrect assumption, ECR policies do allow for missing Resources
// func testPolicyResourceElementResourceMissingFails(t *testing.T) {
// 	// Set up
// 	userAccountId := "012345678901"
// 	policyContent := `
//     {
//       "Version": "2012-10-17",
//       "Statement": [
//         {
//           "Effect": "Allow",
//           "Action": "ec2:DescribeVolumes",
//           "Principal": {
//             "AWS": "012345678901"
//           }
//         }
//       ]
//     }
// 	`

// 	// Test
// 	_, err := EvaluatePolicy(policyContent, userAccountId)

// 	// Evaluate
// 	if err == nil {
// 		t.Fatal("Expected error but no error was returned from EvaluatePolicy")
// 	}

// 	expectedErrorMsg := "policy element Resource is missing"

// 	if errorMsg := err.Error(); errorMsg != expectedErrorMsg {
// 		t.Fatalf("The error message returned: '%s' but was expected to be: '%s'", errorMsg, expectedErrorMsg)
// 	}
// }

func testPolicyResourceElementResourcePresent(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": "012345678901"
          },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"Write"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Principal Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func TestSidElement(t *testing.T) {
	t.Run("TestKnownSidInASingleStatementThatAllowsSharedAccess", testKnownSidInASingleStatementThatAllowsSharedAccess)
	t.Run("TestKnownSidInASingleStatementThatAllowsPrivateAccess", testKnownSidInASingleStatementThatAllowsPrivateAccess)

	t.Run("TestKnownSidInASingleStatementThatAllowsPublicAccess", testKnownSidInASingleStatementThatAllowsPublicAccess)
	t.Run("TestKnownSidsInMultipleStatementsThatAllowsPublicAccessInIncreasingOrder", testKnownSidsInMultipleStatementsThatAllowsPublicAccessInIncreasingOrder)
	t.Run("TestKnownSidsInMultipleStatementsThatAllowsPublicAccessInDecreasingOrder", testKnownSidsInMultipleStatementsThatAllowsPublicAccessInDecreasingOrder)
	t.Run("TestKnownSidsInMultipleStatementsThatHaveDuplicateNamesFails", testKnownSidsInMultipleStatementsThatHaveDuplicateNamesFails)
	t.Run("TestUnknownSidInASingleStatementThatAllowsPublicAccess", testUnknownSidInASingleStatementThatAllowsPublicAccess)
	t.Run("TestUnknownSidsInMultipleStatementsThatAllowsPublicAccess", testUnknownSidsInMultipleStatementsThatAllowsPublicAccess)
}

func testKnownSidInASingleStatementThatAllowsSharedAccess(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Sid": "Sid_Statement_1",
          "Effect": "Allow",
          "Principal": { "AWS": "444455554444" },
          "Action": "ec2:DescribeVolumes",
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"444455554444"},
		AllowedPrincipalAccountIds:          []string{"444455554444"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Sid_Statement_1"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluateSidTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Sid Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testKnownSidInASingleStatementThatAllowsPrivateAccess(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Sid": "Sid_Statement_1",
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2:DescribeVolumes",
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluateSidTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Sid Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testKnownSidInASingleStatementThatAllowsPublicAccess(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Sid": "Sid_Statement_1",
          "Effect": "Allow",
          "Principal": { "AWS": "*" },
          "Action": "ec2:DescribeVolumes",
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Sid_Statement_1"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluateSidTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Sid Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testKnownSidsInMultipleStatementsThatAllowsPublicAccessInIncreasingOrder(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Sid": "Sid_Statement_1",
          "Effect": "Allow",
          "Principal": { "AWS": "*" },
          "Action": "ec2:DescribeVolumes",
          "Resource": "*"
        },
        {
          "Sid": "Sid_Statement_2",
          "Effect": "Allow",
          "Principal": { "AWS": "*" },
          "Action": "ec2:DescribeVolumes",
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds: []string{
			"Sid_Statement_1",
			"Sid_Statement_2",
		},
		SharedStatementIds: []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluateSidTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Sid Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testKnownSidsInMultipleStatementsThatAllowsPublicAccessInDecreasingOrder(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Sid": "Sid_Statement_2",
          "Effect": "Allow",
          "Principal": { "AWS": "*" },
          "Action": "ec2:DescribeVolumes",
          "Resource": "*"
        },
        {
          "Sid": "Sid_Statement_1",
          "Effect": "Allow",
          "Principal": { "AWS": "*" },
          "Action": "ec2:DescribeVolumes",
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds: []string{
			"Sid_Statement_1",
			"Sid_Statement_2",
		},
		SharedStatementIds: []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluateSidTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Sid Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testKnownSidsInMultipleStatementsThatHaveDuplicateNamesFails(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Sid": "Sid_Statement_1",
          "Effect": "Allow",
          "Principal": { "AWS": "*" },
          "Resource": "*"
        },
        {
          "Sid": "Sid_Statement_1",
          "Effect": "Allow",
          "Principal": { "AWS": "*" },
          "Resource": "*"
        }
      ]
    }
	`

	// Test
	_, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err == nil {
		t.Fatal("Expected error but no error was returned from EvaluatePolicy")
	}

	expectedErrorMsg := "duplicate Sid found: Sid_Statement_1"

	if errorMsg := err.Error(); errorMsg != expectedErrorMsg {
		t.Fatalf("The error message returned: '%s' but was expected to be: '%s'", errorMsg, expectedErrorMsg)
	}
}

func testUnknownSidInASingleStatementThatAllowsPublicAccess(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "*" },
          "Action": "ec2:DescribeVolumes",
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluateSidTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Sid Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testUnknownSidsInMultipleStatementsThatAllowsPublicAccess(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "*" },
          "Action": "ec2:DescribeVolumes",
          "Resource": "*"
        },
        {
          "Effect": "Allow",
          "Principal": { "AWS": "*" },
          "Action": "ec2:DescribeVolumes",
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds: []string{
			"Statement[1]",
			"Statement[2]",
		},
		SharedStatementIds: []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluateSidTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Sid Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func TestAccessLevel(t *testing.T) {
	t.Run("TestPublicPrincipalIsPublicAccess", testPublicPrincipalIsPublicAccess)
	t.Run("TestServicePrincipalIsPublicAccess", testServicePrincipalIsPublicAccess)
	t.Run("TestCrossAccountPrincipalIsSharedAccess", testCrossAccountPrincipalIsSharedAccess)
	t.Run("TestUserAccountPrincipalIsPrivateAccess", testUserAccountPrincipalIsPrivateAccess)
	t.Run("TestAccessLevelSharedHasHigherPrecedenceThanAccessLevelPrivate", testAccessLevelSharedHasHigherPrecedenceThanAccessLevelPrivate)
	t.Run("TestAccessLevelPublicHasHigherPrecedenceThanAccessLevelPrivate", testAccessLevelPublicHasHigherPrecedenceThanAccessLevelPrivate)
	t.Run("TestAccessLevelPublicHasHigherPrecedenceThanAccessLevelShared", testAccessLevelPublicHasHigherPrecedenceThanAccessLevelShared)
}

func testPublicPrincipalIsPublicAccess(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "*" },
          "Action": "ec2:DescribeVolumes",
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluateAccessLevelTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("AccessLevel Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testServicePrincipalIsPublicAccess(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "Service": ["ecs.amazonaws.com"] },
          "Action": "ec2:DescribeVolumes",
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluateAccessLevelTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("AccessLevel Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testCrossAccountPrincipalIsSharedAccess(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "111122221111" },
          "Action": "ec2:DescribeVolumes",
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"111122221111"},
		AllowedPrincipalAccountIds:          []string{"111122221111"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluateAccessLevelTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("AccessLevel Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testUserAccountPrincipalIsPrivateAccess(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2:DescribeVolumes",
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluateAccessLevelTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("AccessLevel Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testAccessLevelSharedHasHigherPrecedenceThanAccessLevelPrivate(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "111122221111" },
          "Action": "ec2:DescribeVolumes",
          "Resource": "*"
        },
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2:DescribeVolumes",
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"012345678901",
			"111122221111",
		},
		AllowedPrincipalAccountIds: []string{
			"012345678901",
			"111122221111",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluateAccessLevelTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("AccessLevel Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testAccessLevelPublicHasHigherPrecedenceThanAccessLevelPrivate(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "*" },
          "Action": "ec2:DescribeVolumes",
          "Resource": "*"
        },
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2:DescribeVolumes",
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "public",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"*",
			"012345678901",
		},
		AllowedPrincipalAccountIds: []string{
			"*",
			"012345678901",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluateAccessLevelTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("AccessLevel Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testAccessLevelPublicHasHigherPrecedenceThanAccessLevelShared(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "*" },
          "Action": "ec2:DescribeVolumes",
          "Resource": "*"
        },
        {
          "Effect": "Allow",
          "Principal": { "AWS": "111122221111" },
          "Action": "ec2:DescribeVolumes",
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "public",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"*",
			"111122221111",
		},
		AllowedPrincipalAccountIds: []string{
			"*",
			"111122221111",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{"Statement[2]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluateAccessLevelTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("AccessLevel Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func TestPolicyActionElement(t *testing.T) {
	t.Run("TestPolicyActionUnknownApiService", testPolicyActionUnknownApiService)
	t.Run("TestPolicyActionUnknownApiFunction", testPolicyActionUnknownApiFunction)
	t.Run("TestPolicyActionKnownApiFunction", testPolicyActionKnownApiFunction)

	t.Run("TestPolicyActionMultipleStatementsWithKnownApiFunctions", testPolicyActionMultipleStatementsWithKnownApiFunctions)

	t.Run("TestPolicyActionFullWildcard", testPolicyActionFullWildcard)

	t.Run("TestPolicyActionSingleFullWildcardWithNoActionName", testPolicyActionSingleFullWildcardWithNoActionName)
	t.Run("TestPolicyActionSingleFullWildcardAtFrontOfAction", testPolicyActionSingleFullWildcardAtFrontOfAction)
	t.Run("TestPolicyActionSingleFullWildcardInMiddleOfAction", testPolicyActionSingleFullWildcardInMiddleOfAction)
	t.Run("TestPolicyActionSingleFullWildcardAtEndOfAction", testPolicyActionSingleFullWildcardAtEndOfAction)

	t.Run("TestPolicyActionSinglePartialWildcardAtFrontOfAction", testPolicyActionSinglePartialWildcardAtFrontOfAction)
	t.Run("TestPolicyActionSinglePartialWildcardInMiddleOfAction", testPolicyActionSinglePartialWildcardInMiddleOfAction)
	t.Run("TestPolicyActionSinglePartialWildcardAtEndOfAction", testPolicyActionSinglePartialWildcardAtEndOfAction)
	t.Run("TestPolicyActionMultipleWildcardsInAction", testPolicyActionMultipleWildcardsInAction)

	t.Run("TestPolicyActionSinglePartialWildcardAtEndOfKnownApiFunctionAction", testPolicyActionSinglePartialWildcardAtEndOfKnownApiFunctionAction)
	t.Run("TestPolicyActionSingleFullWildcardAtEndOfKnownApiFunctionAction", testPolicyActionSingleFullWildcardAtEndOfKnownApiFunctionAction)

	t.Run("TestPolicyActionIncompleteActionMissingFunctionPattern", testPolicyActionIncompleteActionMissingFunctionPattern)
	t.Run("TestPolicyActionWhenServiceNameIsGivenOnly", testPolicyActionWhenServiceNameIsGivenOnly)

	t.Run("TestPolicyActionWhenServiceNameIsMissing", testPolicyActionWhenServiceNameIsMissing)
}

func testPolicyActionUnknownApiService(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec:StartInstances",
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePublicAccessLevelsTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("PublicAccessLevels Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPolicyActionUnknownApiFunction(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2:PescribeVolumes",
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePublicAccessLevelsTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("PublicAccessLevels Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPolicyActionKnownApiFunction(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2:DescribeVolumesModifications",
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"Read"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePublicAccessLevelsTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("PublicAccessLevels Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPolicyActionMultipleStatementsWithKnownApiFunctions(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2:DescribeVolumes",
          "Resource": "*"
        },
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2:DescribeVolumesModifications",
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels: []string{
			"List",
			"Read",
		},
		PublicStatementIds: []string{},
		SharedStatementIds: []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePublicAccessLevelsTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("PublicAccessLevels Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPolicyActionFullWildcard(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "*",
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels: []string{
			"List",
			"Permissions management",
			"Read",
			"Tagging",
			"Write",
		},
		PublicStatementIds: []string{},
		SharedStatementIds: []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePublicAccessLevelsTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("PublicAccessLevels Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPolicyActionSingleFullWildcardWithNoActionName(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2:*",
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels: []string{
			"List",
			"Permissions management",
			"Read",
			"Tagging",
			"Write",
		},
		PublicStatementIds: []string{},
		SharedStatementIds: []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePublicAccessLevelsTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("PublicAccessLevels Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPolicyActionSingleFullWildcardAtFrontOfAction(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2:*VpcClassicLink",
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels: []string{
			"List",
			"Write",
		},
		PublicStatementIds: []string{},
		SharedStatementIds: []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePublicAccessLevelsTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("PublicAccessLevels Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPolicyActionSingleFullWildcardInMiddleOfAction(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2:Describe*Attributes",
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePublicAccessLevelsTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("PublicAccessLevels Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPolicyActionSingleFullWildcardAtEndOfAction(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2:Describe*",
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels: []string{
			"List",
			"Read",
		},
		PublicStatementIds: []string{},
		SharedStatementIds: []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePublicAccessLevelsTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("PublicAccessLevels Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPolicyActionSinglePartialWildcardAtFrontOfAction(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2:?escribeVolumesModifications",
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"Read"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePublicAccessLevelsTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("PublicAccessLevels Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPolicyActionSinglePartialWildcardInMiddleOfAction(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2:DescribeVolumes?odifications",
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"Read"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePublicAccessLevelsTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("PublicAccessLevels Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPolicyActionSinglePartialWildcardAtEndOfAction(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2:DescribeVolumesModification?",
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"Read"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePublicAccessLevelsTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("PublicAccessLevels Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPolicyActionMultipleWildcardsInAction(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2:*Volumes*?",
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels: []string{
			"List",
			"Read",
		},
		PublicStatementIds: []string{},
		SharedStatementIds: []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePublicAccessLevelsTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("PublicAccessLevels Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPolicyActionSinglePartialWildcardAtEndOfKnownApiFunctionAction(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2:StartInstances?",
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePublicAccessLevelsTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("PublicAccessLevels Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPolicyActionSingleFullWildcardAtEndOfKnownApiFunctionAction(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2:StartInstances*",
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"Write"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePublicAccessLevelsTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("PublicAccessLevels Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPolicyActionIncompleteActionMissingFunctionPattern(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2:",
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePublicAccessLevelsTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("PublicAccessLevels Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPolicyActionWhenServiceNameIsGivenOnly(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2",
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePublicAccessLevelsTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("PublicAccessLevels Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPolicyActionWhenServiceNameIsMissing(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Resource": "*"
        },
        {
          "Effect": "Allow",
          "Principal": { "AWS": "222233332222" },
          "Resource": "*"
        },
        {
          "Effect": "Allow",
          "Principal": { "AWS": "*" },
          "Resource": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePublicAccessLevelsTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("PublicAccessLevels Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func TestGlobalConditionSourceArn(t *testing.T) {
	// StringEquals
	t.Run("TestSourceArnConditionWhenValueIsAUserAccountUsingStringEquals", testSourceArnConditionWhenValueIsAUserAccountUsingStringEquals)
	t.Run("TestSourceArnConditionWhenValueIsACrossAccountUsingStringEquals", testSourceArnConditionWhenValueIsACrossAccountUsingStringEquals)
	t.Run("TestSourceArnConditionWhenValueIsFullWildcardUsingStringEquals", testSourceArnConditionWhenValueIsFullWildcardUsingStringEquals)
	t.Run("TestSourceArnConditionConditionWhenValueIsPartialWildcardUsingStringEquals", testSourceArnConditionConditionWhenValueIsPartialWildcardUsingStringEquals)
	t.Run("TestSourceArnConditionUsingStringEqualsIfExists", testSourceArnConditionUsingStringEqualsIfExists)
	t.Run("TestSourceArnConditionWhenArnIsAGlobalResourcedWithStringEquals", testSourceArnConditionWhenArnIsAGlobalResourcedWithStringEquals)
	// StringNotEquals
	// StringEqualsIgnoreCase
	t.Run("TestSourceArnConditionWhenValueIsAUserAccountUsingStringEqualsIgnoreCase", testSourceArnConditionWhenValueIsAUserAccountUsingStringEqualsIgnoreCase)
	t.Run("TestSourceArnConditionWhenValueIsACrossAccountUsingStringEqualsIgnoreCase", testSourceArnConditionWhenValueIsACrossAccountUsingStringEqualsIgnoreCase)
	t.Run("TestSourceArnConditionWhenValueIsFullWildcardUsingStringEqualsIgnoreCase", testSourceArnConditionWhenValueIsFullWildcardUsingStringEqualsIgnoreCase)
	t.Run("TestSourceArnConditionConditionWhenValueIsPartialWildcardUsingStringEqualsIgnoreCase", testSourceArnConditionConditionWhenValueIsPartialWildcardUsingStringEqualsIgnoreCase)
	t.Run("TestSourceArnConditionUsingStringEqualsIgnoreCaseIfExists", testSourceArnConditionUsingStringEqualsIgnoreCaseIfExists)
	t.Run("TestSourceArnConditionWhenArnIsAGlobalResourcedWithStringEqualsIgnoreCase", testSourceArnConditionWhenArnIsAGlobalResourcedWithStringEqualsIgnoreCase)
	// StringNotEqualsIgnoreCase
	// StringLike
	t.Run("TestSourceArnConditionWhenValueIsAUserAccountWithStringLike", testSourceArnConditionWhenValueIsAUserAccountWithStringLike)
	t.Run("TestSourceArnConditionWhenValueIsACrossAccountWithStringLike", testSourceArnConditionWhenValueIsACrossAccountWithStringLike)
	t.Run("TestSourceArnConditionWhenValueIsAnAccountWithOneDigitTooFewWithStringLike", testSourceArnConditionWhenValueIsAnAccountWithOneDigitTooFewWithStringLike)
	t.Run("TestSourceArnConditionWhenValueIsAnAccountWithOneDigitTooManyWithStringLike", testSourceArnConditionWhenValueIsAnAccountWithOneDigitTooManyWithStringLike)
	t.Run("TestSourceArnConditionWhenValueIsFullWildcardWithStringLike", testSourceArnConditionWhenValueIsFullWildcardWithStringLike)
	t.Run("TestSourceArnConditionWhenValueIsPartialWildcardWithStringLike", testSourceArnConditionWhenValueIsPartialWildcardWithStringLike)
	t.Run("TestSourceArnConditionUsingStringLikeIfExists", testSourceArnConditionUsingStringLikeIfExists)
	t.Run("TestSourceArnConditionWhenArnIsAGlobalResourcedWithStringLike", testSourceArnConditionWhenArnIsAGlobalResourcedWithStringLike)
	// StringNotLike
	// String Other
	t.Run("TestSourceArnConditionWhenValueWhenArnIsMalformedUsingStringOperators", testSourceArnConditionWhenValueWhenArnIsMalformedUsingStringOperators)
	t.Run("TestSourceArnConditionWhenValueWhenAccountIsOneDigitTooFewUsingStringOperators", testSourceArnConditionWhenValueWhenAccountIsOneDigitTooFewUsingStringOperators)
	t.Run("TestSourceArnConditionWhenValueWhenAccountIsOneDigitTooManyUsingStringOperators", testSourceArnConditionWhenValueWhenAccountIsOneDigitTooManyUsingStringOperators)
	t.Run("TestSourceArnConditionWithMultipleValuesUsingStringOperators", testSourceArnConditionWithMultipleValuesUsingStringOperators)

	// ArnEquals
	t.Run("TestSourceArnConditionWhenValueIsAUserAccountUsingArnEquals", testSourceArnConditionWhenValueIsAUserAccountUsingArnEquals)
	t.Run("TestSourceArnConditionWhenValueIsACrossAccountUsingArnEquals", testSourceArnConditionWhenValueIsACrossAccountUsingArnEquals)
	t.Run("TestSourceArnConditionWhenValueIsFullWildcardUsingArnEquals", testSourceArnConditionWhenValueIsFullWildcardUsingArnEquals)
	t.Run("TestSourceArnConditionUsingArnEqualsIfExists", testSourceArnConditionUsingArnEqualsIfExists)
	t.Run("TestSourceArnConditionWhenArnIsAGlobalResourcedWithArnEquals", testSourceArnConditionWhenArnIsAGlobalResourcedWithArnEquals)
	// ArnNotEquals
	// ArnLike
	t.Run("TestSourceArnConditionWhenValueIsAUserAccountWithArnLike", testSourceArnConditionWhenValueIsAUserAccountWithArnLike)
	t.Run("TestSourceArnConditionWhenValueIsACrossAccountWithArnLike", testSourceArnConditionWhenValueIsACrossAccountWithArnLike)
	t.Run("TestSourceArnConditionWhenValueStopsAtAccountSection", testSourceArnConditionWhenValueStopsAtAccountSection)
	t.Run("TestSourceArnConditionWhenValueIsMissingValueInAccountSection", testSourceArnConditionWhenValueIsMissingValueInAccountSection)
	t.Run("TestSourceArnConditionWhenValueIsAnAccountWithOneDigitTooFewWithArnLike", testSourceArnConditionWhenValueIsAnAccountWithOneDigitTooFewWithArnLike)
	t.Run("TestSourceArnConditionWhenValueIsAnAccountWithOneDigitTooManyWithArnLike", testSourceArnConditionWhenValueIsAnAccountWithOneDigitTooManyWithArnLike)
	t.Run("TestSourceArnConditionWhenValueIsFullWildcardWithArnLike", testSourceArnConditionWhenValueIsFullWildcardWithArnLike)
	t.Run("TestSourceArnConditionWhenValueIsInvalidValueWithArnLike", testSourceArnConditionWhenValueIsInvalidValueWithArnLike)
	t.Run("TestSourceArnConditionUsingArnLikeIfExists", testSourceArnConditionUsingArnLikeIfExists)
	t.Run("TestSourceArnConditionWhenValueIsPartialWildcardWithArnLike", testSourceArnConditionWhenValueIsPartialWildcardWithArnLike)
	t.Run("TestSourceArnConditionWhenValueWhenAccountIsWildcardedOneTooFewUsingArnLike", testSourceArnConditionWhenValueWhenAccountIsWildcardedOneTooFewUsingArnLike)
	t.Run("TestSourceArnConditionWhenValueWhenAccountIsWildcardedOneTooManyUsingArnLike", testSourceArnConditionWhenValueWhenAccountIsWildcardedOneTooManyUsingArnLike)
	t.Run("TestSourceArnConditionWhenArnIsAGlobalResourcedWithArnLike", testSourceArnConditionWhenArnIsAGlobalResourcedWithArnLike)
	// ArnNotLike
	// Arn Other
	t.Run("TestSourceArnConditionWhenValueWhenArnIsMalformedUsingArnOperators", testSourceArnConditionWhenValueWhenArnIsMalformedUsingArnOperators)
	t.Run("TestSourceArnConditionWhenValueWhenAccountIsOneDigitTooFewUsingArnOperators", testSourceArnConditionWhenValueWhenAccountIsOneDigitTooFewUsingArnOperators)
	t.Run("TestSourceArnConditionWhenValueWhenAccountIsOneDigitTooManyUsingArnOperators", testSourceArnConditionWhenValueWhenAccountIsOneDigitTooManyUsingArnOperators)
	t.Run("TestSourceArnConditionWithMultipleValuesUsingArnOperators", testSourceArnConditionWithMultipleValuesUsingArnOperators)

	// Others
	t.Run("TestSourceArnConditionIsNotAnArnOrStringType", testSourceArnConditionIsNotAnArnOrStringType)
	t.Run("TestSourceArnConditionWhenUnknownOperatorType", testSourceArnConditionWhenUnknownOperatorType)
	t.Run("TestSourceArnConditionWhenAcrossMultipleStatements", testSourceArnConditionWhenAcrossMultipleStatements)
}

func testSourceArnConditionWhenValueIsAUserAccountUsingStringEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceArn": ["arn:aws:iam::012345678901:root"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::012345678901:root"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenValueIsACrossAccountUsingStringEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceArn": ["arn:aws:iam::222233332222:root"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::222233332222:root"},
		AllowedPrincipalAccountIds:          []string{"222233332222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenValueIsFullWildcardUsingStringEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceArn": ["*"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionConditionWhenValueIsPartialWildcardUsingStringEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceArn": ["arn:*", "arn:?"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionUsingStringEqualsIfExists(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIfExists": {
              "aws:SourceArn": ["arn:aws:iam::012345678901:root"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenArnIsAGlobalResourcedWithStringEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceArn": ["arn:aws:s3:::bucket"]
            }
          },
          "Principal": {
            "Service": ["ecs.amazonaws.com"]
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:s3:::bucket"},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenValueIsAUserAccountUsingStringEqualsIgnoreCase(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIgnoreCase": {
              "aws:SourceArn": ["arn:aws:iam::012345678901:root"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::012345678901:root"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenValueIsACrossAccountUsingStringEqualsIgnoreCase(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIgnoreCase": {
              "aws:SourceArn": ["arn:aws:iam::222233332222:root"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::222233332222:root"},
		AllowedPrincipalAccountIds:          []string{"222233332222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenValueIsFullWildcardUsingStringEqualsIgnoreCase(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIgnoreCase": {
              "aws:SourceArn": ["*"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionConditionWhenValueIsPartialWildcardUsingStringEqualsIgnoreCase(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIgnoreCase": {
              "aws:SourceArn": ["arn:*", "arn:?"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionUsingStringEqualsIgnoreCaseIfExists(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIgnoreCaseIfExists": {
              "aws:SourceArn": ["arn:aws:iam::012345678901:root"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenArnIsAGlobalResourcedWithStringEqualsIgnoreCase(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
   {
     "Version": "2012-10-17",
     "Statement": [
       {
         "Effect": "Allow",
         "Action": "ec2:DescribeVolumes",
         "Resource": "*",
         "Condition": {
           "StringEqualsIgnoreCase": {
             "aws:SourceArn": ["arn:aws:s3:::bucket"]
           }
         },
         "Principal": {
           "Service": ["ecs.amazonaws.com"]
         }
       }
     ]
   }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:s3:::bucket"},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenValueIsAUserAccountWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:SourceArn": ["arn:*:012345678901:*"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:*:012345678901:*"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenValueIsACrossAccountWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:SourceArn": ["arn:*:222233332222:*"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:*:222233332222:*"},
		AllowedPrincipalAccountIds:          []string{"222233332222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenValueIsAnAccountWithOneDigitTooFewWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:SourceArn": ["arn:*:22223333222:*"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:*:22223333222:*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenValueIsAnAccountWithOneDigitTooManyWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:SourceArn": ["arn:*:2222333322222:*"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:*:2222333322222:*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenValueIsFullWildcardWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:SourceArn": ["*"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenValueIsPartialWildcardWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:SourceArn": ["*1234*"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*1234*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionUsingStringLikeIfExists(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLikeIfExists": {
              "aws:SourceArn": ["arn:aws:iam::222233332222:*"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenArnIsAGlobalResourcedWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:SourceArn": ["arn:aws:s3:::bucket"]
            }
          },
          "Principal": {
            "Service": ["ecs.amazonaws.com"]
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:s3:::bucket"},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenValueWhenArnIsMalformedUsingStringOperators(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceArn": ["arn:aws:iam:wrong:wrong:012345678901:root"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenValueWhenAccountIsOneDigitTooFewUsingStringOperators(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceArn": ["arn:aws:iam::01234567890:root"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenValueWhenAccountIsOneDigitTooManyUsingStringOperators(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceArn": ["arn:aws:iam::0123456789012:root"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWithMultipleValuesUsingStringOperators(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceArn": [
                "arn:aws:iam::012345678901:root",
                "*",
                "arn:aws:iam::222233332222:root"
              ]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:iam::012345678901:root",
			"arn:aws:iam::222233332222:root",
		},
		AllowedPrincipalAccountIds: []string{
			"012345678901",
			"222233332222",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenValueIsAUserAccountUsingArnEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnEquals": {
              "aws:SourceArn": ["arn:aws:iam::012345678901:root"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::012345678901:root"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenValueIsACrossAccountUsingArnEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnEquals": {
              "aws:SourceArn": ["arn:aws:iam::222233332222:root"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::222233332222:root"},
		AllowedPrincipalAccountIds:          []string{"222233332222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenValueIsFullWildcardUsingArnEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnEquals": {
              "aws:SourceArn": ["*"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionUsingArnEqualsIfExists(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnEqualsIfExists": {
              "aws:SourceArn": ["arn:aws:iam::012345678901:root"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenArnIsAGlobalResourcedWithArnEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnEquals": {
              "aws:SourceArn": ["arn:aws:s3:::bucket"]
            }
          },
          "Principal": {
            "Service": ["ecs.amazonaws.com"]
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:s3:::bucket"},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenValueIsAUserAccountWithArnLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLike": {
              "aws:SourceArn": ["arn:aws:iam::012345678901:*"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::012345678901:*"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenValueIsACrossAccountWithArnLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLike": {
              "aws:SourceArn": ["arn:aws:iam::222233332222:*"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::222233332222:*"},
		AllowedPrincipalAccountIds:          []string{"222233332222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenValueStopsAtAccountSection(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLike": {
              "aws:SourceArn": ["arn:aws:iam::*"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenValueIsMissingValueInAccountSection(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLike": {
              "aws:SourceArn": ["arn:aws:iam::"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenValueIsAnAccountWithOneDigitTooFewWithArnLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLike": {
              "aws:SourceArn": ["arn:aws:iam::22223333222:*"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenValueIsAnAccountWithOneDigitTooManyWithArnLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLike": {
              "aws:SourceArn": ["arn:aws:iam::2222333322223:*"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenValueIsFullWildcardWithArnLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLike": {
              "aws:SourceArn": ["arn:aws:iam::*:*"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::*:*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenValueIsInvalidValueWithArnLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLike": {
              "aws:SourceArn": ["arn:aws:iam::01234567890A"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionUsingArnLikeIfExists(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLikeIfExists": {
              "aws:SourceArn": ["arn:aws:iam::222233332222:*"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenValueIsPartialWildcardWithArnLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLike": {
              "aws:SourceArn": ["arn:aws:iam::0123456789??:root", "arn:aws:iam::2222333322*:root"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "public",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:iam::0123456789??:root",
			"arn:aws:iam::2222333322*:root",
		},
		AllowedPrincipalAccountIds: []string{
			"0123456789??",
			"2222333322*",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenValueWhenAccountIsWildcardedOneTooFewUsingArnLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLike": {
              "aws:SourceArn": ["arn:aws:iam::0123456789?:root"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenValueWhenAccountIsWildcardedOneTooManyUsingArnLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLike": {
              "aws:SourceArn": ["arn:aws:iam::0123456789???:root"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenArnIsAGlobalResourcedWithArnLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLike": {
              "aws:SourceArn": ["arn:aws:s3:::bucket"]
            }
          },
          "Principal": {
            "Service": ["ecs.amazonaws.com"]
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:s3:::bucket"},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenValueWhenArnIsMalformedUsingArnOperators(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnEquals": {
              "aws:SourceArn": ["arn:aws:iam:wrong:wrong:012345678901:root"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenValueWhenAccountIsOneDigitTooFewUsingArnOperators(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnEquals": {
              "aws:SourceArn": ["arn:aws:iam::01234567890:root"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenValueWhenAccountIsOneDigitTooManyUsingArnOperators(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnEquals": {
              "aws:SourceArn": ["arn:aws:iam::0123456789012:root"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWithMultipleValuesUsingArnOperators(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnEquals": {
              "aws:SourceArn": [
                "arn:aws:iam::012345678901:root",
                "*",
                "arn:aws:iam::222233332222:root"
              ]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:iam::012345678901:root",
			"arn:aws:iam::222233332222:root",
		},
		AllowedPrincipalAccountIds: []string{
			"012345678901",
			"222233332222",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionIsNotAnArnOrStringType(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "NumericEquals": {
              "aws:SourceArn": ["arn:aws:iam::222233332222:root"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenUnknownOperatorType(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringUnknown": {
              "aws:SourceArn": ["arn:aws:iam::222233332222:root"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceArnConditionWhenAcrossMultipleStatements(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceArn": ["arn:aws:iam::012345678901:root"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:SourceArn": ["*"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnEquals": {
              "aws:SourceArn": ["arn:aws:iam::222233332222:root"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "public",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"*",
			"arn:aws:iam::012345678901:root",
			"arn:aws:iam::222233332222:root",
		},
		AllowedPrincipalAccountIds: []string{
			"*",
			"012345678901",
			"222233332222",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{"Statement[2]"},
		SharedStatementIds:                  []string{"Statement[3]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func TestGlobalConditionPrincipalArn(t *testing.T) {
	// StringEquals
	t.Run("TestPrincipalArnConditionWhenValueIsAUserAccountArnUsingStringEquals", testPrincipalArnConditionWhenValueIsAUserAccountArnUsingStringEquals)
	t.Run("TestPrincipalArnConditionWhenValueIsACrossAccountArnUsingStringEquals", testPrincipalArnConditionWhenValueIsACrossAccountArnUsingStringEquals)
	t.Run("TestPrincipalArnConditionWhenValueIsAUserAccountUsingStringEquals", testPrincipalArnConditionWhenValueIsAUserAccountUsingStringEquals)
	t.Run("TestPrincipalArnConditionWhenValueIsACrossAccountUsingStringEquals", testPrincipalArnConditionWhenValueIsACrossAccountUsingStringEquals)
	t.Run("TestPrincipalArnConditionWhenValueIsFullWildcardUsingStringEquals", testPrincipalArnConditionWhenValueIsFullWildcardUsingStringEquals)
	t.Run("TestSourceAccountConditionWhenValueIsPartialWildcardUsingStringEquals", testPrincipalArnConditionWhenValueIsPartialWildcardUsingStringEquals)
	t.Run("TestPrincipalArnConditionUsingStringEqualsIfExists", testPrincipalArnConditionUsingStringEqualsIfExists)
	// StringNotEquals
	// StringEqualsIgnoreCase
	t.Run("TestPrincipalArnConditionWhenValueIsAUserAccountUsingStringEqualsIgnoreCase", testPrincipalArnConditionWhenValueIsAUserAccountUsingStringEqualsIgnoreCase)
	t.Run("TestPrincipalArnConditionWhenValueIsACrossAccountUsingStringEqualsIgnoreCase", testPrincipalArnConditionWhenValueIsACrossAccountUsingStringEqualsIgnoreCase)
	t.Run("TestPrincipalArnConditionWhenValueIsFullWildcardUsingStringEqualsIgnoreCase", testPrincipalArnConditionWhenValueIsFullWildcardUsingStringEqualsIgnoreCase)
	t.Run("TestPrincipalArnConditionWhenValueIsPartialWildcardUsingStringEqualsIgnoreCase", testPrincipalArnConditionWhenValueIsPartialWildcardUsingStringEqualsIgnoreCase)
	t.Run("TestPrincipalArnConditionUsingStringEqualsIgnoreCaseIfExists", testPrincipalArnConditionUsingStringEqualsIgnoreCaseIfExists)
	// StringNotEqualsIgnoreCase
	// StringLike
	t.Run("TestPrincipalArnConditionWhenValueIsAUserAccountWithStringLike", testPrincipalArnConditionWhenValueIsAUserAccountWithStringLike)
	t.Run("TestPrincipalArnConditionWhenValueIsACrossAccountWithStringLike", testPrincipalArnConditionWhenValueIsACrossAccountWithStringLike)
	t.Run("TestPrincipalArnConditionWhenPrincipalValueIsAFullWildcardWithStringLike", testPrincipalArnConditionWhenPrincipalValueIsAFullWildcardWithStringLike)
	t.Run("TestPrincipalArnConditionWhenValueIsAnAccountWithOneDigitTooFewWithStringLike", testPrincipalArnConditionWhenValueIsAnAccountWithOneDigitTooFewWithStringLike)
	t.Run("TestPrincipalArnConditionWhenValueIsAnAccountWithOneDigitTooManyWithStringLike", testPrincipalArnConditionWhenValueIsAnAccountWithOneDigitTooManyWithStringLike)
	t.Run("TestPrincipalArnConditionWhenValueIsFullWildcardWithStringLike", testPrincipalArnConditionWhenValueIsFullWildcardWithStringLike)
	t.Run("TestPrincipalArnConditionWhenValueIsPartialWildcardWithStringLike", testPrincipalArnConditionWhenValueIsPartialWildcardWithStringLike)
	t.Run("TestPrincipalArnConditionWhenPrincipalValueWildcardAndConditionIsPartialWildcardWithStringLike", testPrincipalArnConditionWhenPrincipalValueWildcardAndConditionIsPartialWildcardWithStringLike)
	t.Run("TestPrincipalArnConditionUsingStringLikeIfExists", testPrincipalArnConditionUsingStringLikeIfExists)
	// StringNotLike
	// String Other
	t.Run("TestPrincipalArnConditionWhenValueWhenArnIsMalformedUsingStringOperators", testPrincipalArnConditionWhenValueWhenArnIsMalformedUsingStringOperators)
	t.Run("TestPrincipalArnConditionWhenValueWhenAccountIsOneDigitTooFewUsingStringOperators", TestPrincipalArnConditionWhenValueWhenAccountIsOneDigitTooFewUsingStringOperators)
	t.Run("TestPrincipalArnConditionWhenValueWhenAccountIsOneDigitTooManyUsingStringOperators", testPrincipalArnConditionWhenValueWhenAccountIsOneDigitTooManyUsingStringOperators)
	t.Run("TestPrincipalArnConditionWithMultipleValuesUsingStringOperators", testPrincipalArnConditionWithMultipleValuesUsingStringOperators)

	// ArnEquals
	t.Run("TestPrincipalArnConditionWhenValueIsAUserAccountUsingArnEquals", testPrincipalArnConditionWhenValueIsAUserAccountUsingArnEquals)
	t.Run("TestPrincipalArnConditionWhenValueIsACrossAccountUsingArnEquals", testPrincipalArnConditionWhenValueIsACrossAccountUsingArnEquals)
	t.Run("TestPrincipalArnConditionWhenValueIsFullWildcardUsingArnEquals", testPrincipalArnConditionWhenValueIsFullWildcardUsingArnEquals)
	t.Run("TestPrincipalArnConditionUsingArnEqualsIfExists", testPrincipalArnConditionUsingArnEqualsIfExists)
	// ArnNotEquals
	// ArnLike
	t.Run("TestPrincipalArnConditionWhenValueIsAUserAccountArnWithArnLike", testPrincipalArnConditionWhenValueIsAUserAccountArnWithArnLike)
	t.Run("TestPrincipalArnConditionWhenValueIsACrossAccountArnWithArnLike", testPrincipalArnConditionWhenValueIsACrossAccountArnWithArnLike)
	t.Run("TestPrincipalArnConditionWhenValueIsAUserAccountWithArnLike", testPrincipalArnConditionWhenValueIsAUserAccountWithArnLike)
	t.Run("TestPrincipalArnConditionWhenValueIsACrossAccountWithArnLike", testPrincipalArnConditionWhenValueIsACrossAccountWithArnLike)
	t.Run("TestPrincipalArnConditionWhenPrincipalValueIsAFullWildcardWithArnLike", testPrincipalArnConditionWhenPrincipalValueIsAFullWildcardWithArnLike)
	t.Run("TestPrincipalArnConditionWhenValueIsAUserAccountArnWithWildcardAccountNumberUsingArnLike", testPrincipalArnConditionWhenValueIsAUserAccountArnWithWildcardAccountNumberUsingArnLike)
	t.Run("TestPrincipalArnConditionWhenValueIsACrossAccountArnWithWildcardAccountNumberUsingArnLike", testPrincipalArnConditionWhenValueIsACrossAccountArnWithWildcardAccountNumberUsingArnLike)
	t.Run("TestPrincipalArnConditionWhenValueIsAUserAccountWithWildcardAccountNumberUsingArnLike", testPrincipalArnConditionWhenValueIsAUserAccountWithWildcardAccountNumberUsingArnLike)
	t.Run("TestPrincipalArnConditionWhenValueIsACrossAccountWithWildcardAccountNumberUsingArnLike", testPrincipalArnConditionWhenValueIsACrossAccountWithWildcardAccountNumberUsingArnLike)
	t.Run("TestPrincipalArnConditionWhenValueIsMissingAccountSection", testPrincipalArnConditionWhenValueIsMissingAccountSection)
	t.Run("TestPrincipalArnConditionWhenValueIsMissingValueInAccountSection", testPrincipalArnConditionWhenValueIsMissingValueInAccountSection)
	t.Run("TestPrincipalArnConditionWhenValueIsAnAccountWithOneDigitTooFewWithArnLike", testPrincipalArnConditionWhenValueIsAnAccountWithOneDigitTooFewWithArnLike)
	t.Run("TestPrincipalArnConditionWhenValueIsAnAccountWithOneDigitTooManyWithArnLike", testPrincipalArnConditionWhenValueIsAnAccountWithOneDigitTooManyWithArnLike)
	t.Run("TestPrincipalArnConditionWhenValueIsFullWildcardWithArnLike", testPrincipalArnConditionWhenValueIsFullWildcardWithArnLike)
	t.Run("TestPrincipalArnConditionWhenPrincipalValueIsFullWildcardWithConditionAsFullWildcardWithArnLike", testPrincipalArnConditionWhenPrincipalValueIsFullWildcardWithConditionAsFullWildcardWithArnLike)
	t.Run("TestPrincipalArnConditionWhenValueIsInvalidValueWithArnLike", testPrincipalArnConditionWhenValueIsInvalidValueWithArnLike)
	t.Run("TestPrincipalArnConditionUsingArnLikeIfExists", testPrincipalArnConditionUsingArnLikeIfExists)
	t.Run("TestPrincipalArnConditionWhenValueIsPartialWildcardWithArnLike", testPrincipalArnConditionWhenValueIsPartialWildcardWithArnLike)
	t.Run("TestPrincipalArnConditionWhenPrincipalValueWildcardAndConditionIsPartialWildcardWithArnLike", testPrincipalArnConditionWhenPrincipalValueWildcardAndConditionIsPartialWildcardWithArnLike)
	t.Run("TestPrincipalArnConditionWhenValueWhenAccountIsWildcardedOneTooFewUsingArnLike", testPrincipalArnConditionWhenValueWhenAccountIsWildcardedOneTooFewUsingArnLike)
	t.Run("TestPrincipalArnConditionWhenValueWhenAccountIsWildcardedOneTooManyUsingArnLike", testPrincipalArnConditionWhenValueWhenAccountIsWildcardedOneTooManyUsingArnLike)

	// ArnNotLike
	// Arn Other
	t.Run("TestPrincipalArnConditionWhenValueWhenArnIsMalformedUsingArnOperators", testPrincipalArnConditionWhenValueWhenArnIsMalformedUsingArnOperators)
	t.Run("TestPrincipalArnConditionWhenValueWhenAccountIsOneDigitTooFewUsingArnOperators", testPrincipalArnConditionWhenValueWhenAccountIsOneDigitTooFewUsingArnOperators)
	t.Run("TestPrincipalArnConditionWhenValueWhenAccountIsOneDigitTooManyUsingArnOperators", testPrincipalArnConditionWhenValueWhenAccountIsOneDigitTooManyUsingArnOperators)
	t.Run("TestPrincipalArnConditionWithMultipleValuesUsingArnOperators", testPrincipalArnConditionWithMultipleValuesUsingArnOperators)

	// Others
	t.Run("TestPrincipalArnConditionIsNotAnArnOrStringType", testPrincipalArnConditionIsNotAnArnOrStringType)
	t.Run("TestPrincipalArnConditionWhenUnknownOperatorType", testPrincipalArnConditionWhenUnknownOperatorType)
	t.Run("TestPrincipalArnConditionWhenAcrossMultipleStatements", testPrincipalArnConditionWhenAcrossMultipleStatements)
}

func testPrincipalArnConditionWhenValueIsAUserAccountArnUsingStringEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::012345678901:root"]
            }
          },
          "Principal": {
            "AWS": "arn:aws:iam::012345678901:root"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::012345678901:root"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsACrossAccountArnUsingStringEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::222233332222:root"]
            }
          },
          "Principal": {
            "AWS": "arn:aws:iam::222233332222:root"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::222233332222:root"},
		AllowedPrincipalAccountIds:          []string{"222233332222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsAUserAccountUsingStringEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::012345678901:root"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::012345678901:root"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsACrossAccountUsingStringEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::222233332222:root"]
            }
          },
          "Principal": {
            "AWS": "222233332222"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::222233332222:root"},
		AllowedPrincipalAccountIds:          []string{"222233332222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsFullWildcardUsingStringEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["*"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsPartialWildcardUsingStringEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:*", "arn:?"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionUsingStringEqualsIfExists(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIfExists": {
              "aws:PrincipalArn": ["arn:aws:iam::012345678901:root"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsAUserAccountUsingStringEqualsIgnoreCase(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIgnoreCase": {
              "aws:PrincipalArn": ["arn:aws:iam::012345678901:root"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::012345678901:root"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsACrossAccountUsingStringEqualsIgnoreCase(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIgnoreCase": {
              "aws:PrincipalArn": ["arn:aws:iam::222233332222:root"]
            }
          },
          "Principal": {
            "AWS": "222233332222"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::222233332222:root"},
		AllowedPrincipalAccountIds:          []string{"222233332222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsFullWildcardUsingStringEqualsIgnoreCase(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIgnoreCase": {
              "aws:PrincipalArn": ["*"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsPartialWildcardUsingStringEqualsIgnoreCase(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIgnoreCase": {
              "aws:PrincipalArn": ["arn:*", "arn:?"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionUsingStringEqualsIgnoreCaseIfExists(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIgnoreCaseIfExists": {
              "aws:PrincipalArn": ["arn:aws:iam::012345678901:root"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsAUserAccountWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalArn": ["arn:*:012345678901:*"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:*:012345678901:*"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsACrossAccountWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalArn": ["arn:*:222233332222:*"]
            }
          },
          "Principal": {
            "AWS": "222233332222"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:*:222233332222:*"},
		AllowedPrincipalAccountIds:          []string{"222233332222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenPrincipalValueIsAFullWildcardWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalArn": ["arn:*:222233332222:*"]
            }
          },
          "Principal": {
            "AWS": "*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:*:222233332222:*"},
		AllowedPrincipalAccountIds:          []string{"222233332222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsAnAccountWithOneDigitTooFewWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalArn": ["arn:*:22223333222:*"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsAnAccountWithOneDigitTooManyWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalArn": ["arn:*:2222333322222:*"]
            }
          },
          "Principal": {
            "AWS": "222233332222"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"222233332222"},
		AllowedPrincipalAccountIds:          []string{"222233332222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsFullWildcardWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalArn": ["*"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsPartialWildcardWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalArn": ["*1234*"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenPrincipalValueWildcardAndConditionIsPartialWildcardWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalArn": ["*1234*"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionUsingStringLikeIfExists(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLikeIfExists": {
              "aws:PrincipalArn": ["arn:aws:iam::222233332222:*"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueWhenArnIsMalformedUsingStringOperators(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam:wrong:wrong:012345678901:root"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func TestPrincipalArnConditionWhenValueWhenAccountIsOneDigitTooFewUsingStringOperators(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::01234567890:root"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueWhenAccountIsOneDigitTooManyUsingStringOperators(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::0123456789012:root"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWithMultipleValuesUsingStringOperators(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": [
                "arn:aws:iam::012345678901:root",
                "*",
                "arn:aws:iam::222233332222:root"
              ]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "private",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:iam::012345678901:root",
		},
		AllowedPrincipalAccountIds: []string{
			"012345678901",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsAUserAccountUsingArnEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::012345678901:root"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::012345678901:root"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsACrossAccountUsingArnEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::222233332222:root"]
            }
          },
          "Principal": {
            "AWS": "222233332222"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::222233332222:root"},
		AllowedPrincipalAccountIds:          []string{"222233332222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsFullWildcardUsingArnEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnEquals": {
              "aws:PrincipalArn": ["*"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionUsingArnEqualsIfExists(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnEqualsIfExists": {
              "aws:PrincipalArn": ["arn:aws:iam::012345678901:root"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsAUserAccountArnWithArnLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
		{
		  "Version": "2012-10-17",
		  "Statement": [
			{
			  "Effect": "Allow",
			  "Action": "ec2:DescribeVolumes",
			  "Resource": "*",
			  "Condition": {
				"ArnLike": {
				  "aws:PrincipalArn": ["arn:aws:iam::012345678901:*"]
				}
			  },
			  "Principal": {
				"AWS": "arn:aws:iam::012345678901:root"
			  }
			}
		  ]
		}
		`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::012345678901:root"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsACrossAccountArnWithArnLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
		{
		  "Version": "2012-10-17",
		  "Statement": [
			{
			  "Effect": "Allow",
			  "Action": "ec2:DescribeVolumes",
			  "Resource": "*",
			  "Condition": {
				"ArnLike": {
				  "aws:PrincipalArn": ["arn:aws:iam::222233332222:*"]
				}
			  },
			  "Principal": {
				"AWS": "arn:aws:iam::222233332222:root"
			  }
			}
		  ]
		}
		`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::222233332222:root"},
		AllowedPrincipalAccountIds:          []string{"222233332222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsAUserAccountWithArnLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLike": {
              "aws:PrincipalArn": ["arn:aws:iam::012345678901:*"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::012345678901:*"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsACrossAccountWithArnLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLike": {
              "aws:PrincipalArn": ["arn:aws:iam::222233332222:*"]
            }
          },
          "Principal": {
            "AWS": "222233332222"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::222233332222:*"},
		AllowedPrincipalAccountIds:          []string{"222233332222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenPrincipalValueIsAFullWildcardWithArnLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLike": {
              "aws:PrincipalArn": ["arn:aws:iam::222233332222:*"]
            }
          },
          "Principal": {
            "AWS": "*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::222233332222:*"},
		AllowedPrincipalAccountIds:          []string{"222233332222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsAUserAccountArnWithWildcardAccountNumberUsingArnLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
		{
		  "Version": "2012-10-17",
		  "Statement": [
			{
			  "Effect": "Allow",
			  "Action": "ec2:DescribeVolumes",
			  "Resource": "*",
			  "Condition": {
				"ArnLike": {
				  "aws:PrincipalArn": ["arn:aws:iam::0123456789??:root"]
				}
			  },
			  "Principal": {
				"AWS": "arn:aws:iam::012345678901:root"
			  }
			}
		  ]
		}
		`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::012345678901:root"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsACrossAccountArnWithWildcardAccountNumberUsingArnLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
		{
		  "Version": "2012-10-17",
		  "Statement": [
			{
			  "Effect": "Allow",
			  "Action": "ec2:DescribeVolumes",
			  "Resource": "*",
			  "Condition": {
				"ArnLike": {
					"aws:PrincipalArn": ["arn:aws:iam::2222333322*:root"]
				}
			  },
			  "Principal": {
				"AWS": "arn:aws:iam::222233332222:root"
			  }
			}
		  ]
		}
		`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::222233332222:root"},
		AllowedPrincipalAccountIds:          []string{"222233332222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsAUserAccountWithWildcardAccountNumberUsingArnLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLike": {
				"aws:PrincipalArn": ["arn:aws:iam::0123456789??:root"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::012345678901:root"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsACrossAccountWithWildcardAccountNumberUsingArnLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLike": {
              "aws:PrincipalArn": ["arn:aws:iam::2222333322*:root"]
            }
          },
          "Principal": {
            "AWS": "222233332222"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::222233332222:root"},
		AllowedPrincipalAccountIds:          []string{"222233332222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsMissingAccountSection(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLike": {
              "aws:PrincipalArn": ["arn:aws:iam:*"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsMissingValueInAccountSection(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLike": {
              "aws:PrincipalArn": ["arn:aws:iam::"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsAnAccountWithOneDigitTooFewWithArnLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLike": {
              "aws:PrincipalArn": ["arn:aws:iam::22223333222:*"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsAnAccountWithOneDigitTooManyWithArnLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLike": {
              "aws:PrincipalArn": ["arn:aws:iam::2222333322223:*"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsFullWildcardWithArnLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLike": {
              "aws:PrincipalArn": ["arn:aws:iam::*:root"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::012345678901:root"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenPrincipalValueIsFullWildcardWithConditionAsFullWildcardWithArnLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLike": {
              "aws:PrincipalArn": ["arn:aws:iam::*:root"]
            }
          },
          "Principal": {
            "AWS": "*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::*:root"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsInvalidValueWithArnLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLike": {
              "aws:PrincipalArn": ["arn:aws:iam::01234567890A"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        },
		{
			"Effect": "Allow",
			"Action": "ec2:DescribeVolumes",
			"Resource": "*",
			"Condition": {
			  "ArnLike": {
				"aws:PrincipalArn": ["arn:aws:iam::*"]
			  }
			},
			"Principal": {
			  "AWS": "222233332222"
			}
		  }  
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"012345678901",
			"222233332222",
		},
		AllowedPrincipalAccountIds: []string{
			"012345678901",
			"222233332222",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[2]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionUsingArnLikeIfExists(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLikeIfExists": {
              "aws:PrincipalArn": ["arn:aws:iam::222233332222:*"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueIsPartialWildcardWithArnLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLike": {
              "aws:PrincipalArn": ["arn:aws:iam::0123456789??:root", "arn:aws:iam::2222333322*:root"]
            }
          },
          "Principal": {
            "AWS": ["012345678901", "222233332222"]
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:iam::012345678901:root",
			"arn:aws:iam::222233332222:root",
		},
		AllowedPrincipalAccountIds: []string{
			"012345678901",
			"222233332222",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenPrincipalValueWildcardAndConditionIsPartialWildcardWithArnLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLike": {
              "aws:PrincipalArn": ["arn:aws:iam::0123456789??:root", "arn:aws:iam::2222333322*:root"]
            }
          },
          "Principal": {
            "AWS": "*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "public",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:iam::0123456789??:root",
			"arn:aws:iam::2222333322*:root",
		},
		AllowedPrincipalAccountIds: []string{
			"0123456789??",
			"2222333322*",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueWhenAccountIsWildcardedOneTooFewUsingArnLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLike": {
              "aws:PrincipalArn": ["arn:aws:iam::0123456789?:root"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueWhenAccountIsWildcardedOneTooManyUsingArnLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLike": {
              "aws:PrincipalArn": ["arn:aws:iam::0123456789???:root"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueWhenArnIsMalformedUsingArnOperators(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnEquals": {
              "aws:PrincipalArn": ["arn:aws:iam:wrong:wrong:012345678901:root"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueWhenAccountIsOneDigitTooFewUsingArnOperators(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::01234567890:root"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenValueWhenAccountIsOneDigitTooManyUsingArnOperators(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::0123456789012:root"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWithMultipleValuesUsingArnOperators(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnEquals": {
              "aws:PrincipalArn": [
                "arn:aws:iam::012345678901:root",
                "*",
                "arn:aws:iam::222233332222:root"
              ]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "private",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:iam::012345678901:root",
		},
		AllowedPrincipalAccountIds: []string{
			"012345678901",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionIsNotAnArnOrStringType(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "NumericEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::222233332222:root"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenUnknownOperatorType(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringUnknown": {
              "aws:PrincipalArn": ["arn:aws:iam::222233332222:root"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalArnConditionWhenAcrossMultipleStatements(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::012345678901:root"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalArn": ["*"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::222233332222:root"]
            }
          },
          "Principal": {
            "AWS": "222233332222"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"012345678901",
			"arn:aws:iam::012345678901:root",
			"arn:aws:iam::222233332222:root",
		},
		AllowedPrincipalAccountIds: []string{
			"012345678901",
			"222233332222",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[3]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func TestGlobalConditionSourceAccount(t *testing.T) {
	t.Run("TestSourceAccountConditionWhenConditionIsSameAccountAsPrincipal", testSourceAccountConditionWhenConditionIsSameAccountAsPrincipal)
	t.Run("TestSourceAccountConditionWhenConditionIsDifferentAccountToPrincipal", testSourceAccountConditionWhenConditionIsDifferentAccountToPrincipal)
	t.Run("TestSourceAccountConditionWhenConditionIsWildcardedAndPrincipalIsFixed", testSourceAccountConditionWhenConditionIsWildcardedAndPrincipalIsFixed)
	t.Run("TestSourceAccountConditionWhenConditionIsFixedAndPrincipalIsWildcarded", testSourceAccountConditionWhenConditionIsFixedAndPrincipalIsWildcarded)

	t.Run("TestSourceAccountConditionWhenConditionIsSameAccountAsPrincipal", testSourceAccountConditionWhenConditionIsSameAccountAsPrincipalArn)
	t.Run("TestSourceAccountConditionWhenConditionIsSameAccountAsPrincipal", testSourceAccountConditionWhenConditionIsDifferentAccountToPrincipalArn)
	t.Run("TestSourceAccountConditionWhenPrincipalIsArnAndConditionIsWildcarded", testSourceAccountConditionIsWildcardedAgainstAPrincipalArn)

	// StringEquals
	t.Run("TestSourceAccountConditionWhenValueIsAUserAccountUsingStringEquals", testSourceAccountConditionWhenValueIsAUserAccountUsingStringEquals)
	t.Run("TestSourceAccountConditionWhenValueIsACrossAccountUsingStringEquals", testSourceAccountConditionWhenValueIsACrossAccountUsingStringEquals)
	t.Run("TestSourceAccountConditionWhenValueIsFullWildcardUsingStringEquals", testSourceAccountConditionWhenValueIsFullWildcardUsingStringEquals)
	t.Run("TestSourceAccountConditionWhenValueIsPartialWildcardUsingStringEquals", testSourceAccountConditionWhenValueIsPartialWildcardUsingStringEquals)
	t.Run("TestSourceAccountConditionUsingStringEqualsIfExists", testSourceAccountConditionUsingStringEqualsIfExists)
	// StringNotEquals
	// StringEqualsIgnoreCase
	t.Run("TestSourceAccountConditionWhenValueIsAUserAccountUsingStringEqualsIgnoreCase", testSourceAccountConditionWhenValueIsAUserAccountUsingStringEqualsIgnoreCase)
	t.Run("TestSourceAccountConditionWhenValueIsACrossAccountUsingStringEqualsIgnoreCase", testSourceAccountConditionWhenValueIsACrossAccountUsingStringEqualsIgnoreCase)
	t.Run("TestSourceAccountConditionWhenValueIsFullWildcardUsingStringEqualsIgnoreCase", testSourceAccountConditionWhenValueIsFullWildcardUsingStringEqualsIgnoreCase)
	t.Run("TestSourceAccountConditionWhenValueIsPartialWildcardUsingStringEqualsIgnoreCase", testSourceAccountConditionWhenValueIsPartialWildcardUsingStringEqualsIgnoreCase)
	t.Run("TestSourceAccountConditionUsingStringEqualsIgnoreCaseIfExists", testSourceAccountConditionUsingStringEqualsIgnoreCaseIfExists)
	// StringNotEqualsIgnoreCase
	// StringLike
	t.Run("TestSourceAccountConditionWhenValueIsAUserAccountWithStringLike", testSourceAccountConditionWhenValueIsAUserAccountWithStringLike)
	t.Run("TestSourceAccountConditionWhenValueIsACrossAccountWithStringLike", testSourceAccountConditionWhenValueIsACrossAccountWithStringLike)
	t.Run("TestSourceAccountConditionWhenValueIsAnAccountWithOneDigitTooFewWithStringLike", testSourceAccountConditionWhenValueIsAnAccountWithOneDigitTooFewWithStringLike)
	t.Run("TestSourceAccountConditionWhenValueIsAnAccountWithOneDigitTooManyWithStringLike", testSourceAccountConditionWhenValueIsAnAccountWithOneDigitTooManyWithStringLike)
	t.Run("TestSourceAccountConditionWhenValueIsFullWildcardWithStringLike", testSourceAccountConditionWhenValueIsFullWildcardWithStringLike)
	t.Run("TestSourceAccountConditionWhenValueIsPartialWildcardWithStringLike", testSourceAccountConditionWhenValueIsPartialWildcardWithStringLike)
	t.Run("TestSourceAccountConditionUsingStringLikeIfExists", testSourceAccountConditionUsingStringLikeIfExists)
	t.Run("TestSourceAccountConditionWhenValueWhenAccountIsSingleWildcardedUsingStringLike", testSourceAccountConditionWhenValueWhenAccountIsSingleWildcardedUsingStringLike)
	t.Run("TestSourceAccountConditionWhenValueWhenAccountIsWildcardedOneTooFewUsingStringLike", testSourceAccountConditionWhenValueWhenAccountIsWildcardedOneTooFewUsingStringLike)
	t.Run("TestSourceAccountConditionWhenValueWhenAccountIsWildcardedOneTooManyUsingStringLike", testSourceAccountConditionWhenValueWhenAccountIsWildcardedOneTooManyUsingStringLike)

	// StringNotLike
	// String Other
	t.Run("TestSourceAccountConditionWhenValueWhenArnIsMalformedUsingStringOperators", testSourceAccountConditionWhenValueWhenArnIsMalformedUsingStringOperators)
	t.Run("TestSourceAccountConditionWhenValueWhenAccountIsOneDigitTooFewUsingStringOperators", testSourceAccountConditionWhenValueWhenAccountIsOneDigitTooFewUsingStringOperators)
	t.Run("TestSourceAccountConditionWhenValueWhenAccountIsOneDigitTooManyUsingStringOperators", testSourceAccountConditionWhenValueWhenAccountIsOneDigitTooManyUsingStringOperators)
	t.Run("TestSourceAccountConditionWithMultipleValuesUsingStringOperators", testSourceAccountConditionWithMultipleValuesUsingStringOperators)

	// Others
	t.Run("TestSourceAccountConditionIsNotAnArnOrStringType", TestSourceAccountConditionIsNotAStringType)
	t.Run("TestSourceAccountConditionWhenUnknownOperatorType", testSourceAccountConditionWhenUnknownOperatorType)
	t.Run("TestSourceAccountConditionWhenAcrossMultipleStatements", testSourceAccountConditionWhenAcrossMultipleStatements)
}

func testSourceAccountConditionWhenConditionIsSameAccountAsPrincipal(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceAccount": ["012345678901"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceAccountConditionWhenConditionIsDifferentAccountToPrincipal(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceAccount": ["012345678901"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceAccountConditionWhenConditionIsWildcardedAndPrincipalIsFixed(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:SourceAccount": ["*"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceAccountConditionWhenConditionIsFixedAndPrincipalIsWildcarded(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceAccount": ["012345678901"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceAccountConditionWhenConditionIsSameAccountAsPrincipalArn(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceAccount": ["012345678901"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceAccountConditionWhenConditionIsDifferentAccountToPrincipalArn(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceAccount": ["012345678901"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceAccountConditionIsWildcardedAgainstAPrincipalArn(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:SourceAccount": ["*"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceAccountConditionWhenValueIsAUserAccountUsingStringEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceAccount": ["012345678901"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceAccountConditionWhenValueIsACrossAccountUsingStringEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceAccount": ["222233332222"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"222233332222"},
		AllowedPrincipalAccountIds:          []string{"222233332222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceAccountConditionWhenValueIsFullWildcardUsingStringEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceAccount": ["*"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceAccountConditionWhenValueIsPartialWildcardUsingStringEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceAccount": ["12345678*", "1234567890??"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceAccountConditionUsingStringEqualsIfExists(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIfExists": {
              "aws:SourceAccount": ["012345678901"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceAccountConditionWhenValueIsAUserAccountUsingStringEqualsIgnoreCase(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIgnoreCase": {
              "aws:SourceAccount": ["012345678901"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceAccountConditionWhenValueIsACrossAccountUsingStringEqualsIgnoreCase(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIgnoreCase": {
              "aws:SourceAccount": ["222233332222"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"222233332222"},
		AllowedPrincipalAccountIds:          []string{"222233332222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceAccountConditionWhenValueIsFullWildcardUsingStringEqualsIgnoreCase(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIgnoreCase": {
              "aws:SourceAccount": ["012345678901"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceAccountConditionWhenValueIsPartialWildcardUsingStringEqualsIgnoreCase(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIgnoreCase": {
              "aws:SourceAccount": ["12345678*", "1234567890??"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceAccountConditionUsingStringEqualsIgnoreCaseIfExists(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIgnoreCaseIfExists": {
              "aws:SourceAccount": ["012345678901"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceAccountConditionWhenValueIsAUserAccountWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:SourceAccount": ["012345678901"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceAccountConditionWhenValueIsACrossAccountWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:SourceAccount": ["222233332222"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"222233332222"},
		AllowedPrincipalAccountIds:          []string{"222233332222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceAccountConditionWhenValueIsAnAccountWithOneDigitTooFewWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:SourceAccount": ["22223333222"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceAccountConditionWhenValueIsAnAccountWithOneDigitTooManyWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:SourceAccount": ["2222333322222"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceAccountConditionWhenValueIsFullWildcardWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:SourceAccount": ["*"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceAccountConditionWhenValueIsPartialWildcardWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:SourceAccount": ["1234*"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"1234*"},
		AllowedPrincipalAccountIds:          []string{"1234*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceAccountConditionUsingStringLikeIfExists(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLikeIfExists": {
              "aws:SourceAccount": ["222233332222"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceAccountConditionWhenValueWhenAccountIsSingleWildcardedUsingStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:SourceAccount": ["0123456789??"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"0123456789??"},
		AllowedPrincipalAccountIds:          []string{"0123456789??"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceAccountConditionWhenValueWhenAccountIsWildcardedOneTooFewUsingStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:SourceAccount": ["0123456789?"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceAccountConditionWhenValueWhenAccountIsWildcardedOneTooManyUsingStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLike": {
              "aws:SourceAccount": ["0123456789???"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceAccountConditionWhenValueWhenArnIsMalformedUsingStringOperators(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceAccount": ["01234567890A"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceAccountConditionWhenValueWhenAccountIsOneDigitTooFewUsingStringOperators(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceAccount": ["01234567890"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceAccountConditionWhenValueWhenAccountIsOneDigitTooManyUsingStringOperators(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceAccount": ["0123456789012"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceAccountConditionWithMultipleValuesUsingStringOperators(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceAccount": [
                "012345678901",
                "222233332222"
              ]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"012345678901",
			"222233332222",
		},
		AllowedPrincipalAccountIds: []string{
			"012345678901",
			"222233332222",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func TestSourceAccountConditionIsNotAStringType(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnEquals": {
              "aws:SourceAccount": ["222233332222"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceAccountConditionWhenUnknownOperatorType(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringUnknown": {
              "aws:SourceAccount": ["222233332222"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceAccountConditionWhenAcrossMultipleStatements(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceAccount": ["012345678901"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:SourceAccount": ["*"]
            }
          },
          "Principal": {
            "AWS": "*"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceAccount": ["222233332222"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "public",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"*",
			"012345678901",
			"222233332222",
		},
		AllowedPrincipalAccountIds: []string{
			"*",
			"012345678901",
			"222233332222",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{"Statement[2]"},
		SharedStatementIds:                  []string{"Statement[3]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func TestGlobalConditionPrincipalAccount(t *testing.T) {
	t.Run("TestPrincipalAccountConditionWhenConditionIsSameAccountAsPrincipal", testPrincipalAccountConditionWhenConditionIsSameAccountAsPrincipal)
	t.Run("TestPrincipalAccountConditionWhenConditionIsDifferentAccountToPrincipal", testPrincipalAccountConditionWhenConditionIsDifferentAccountToPrincipal)
	t.Run("TestPrincipalAccountConditionWhenConditionIsWildcardedAndPrincipalIsFixed", testPrincipalAccountConditionWhenConditionIsWildcardedAndPrincipalIsFixed)
	t.Run("TestPrincipalAccountConditionWhenConditionIsFixedAndPrincipalIsWildcarded", testPrincipalAccountConditionWhenConditionIsFixedAndPrincipalIsWildcarded)

	t.Run("TestPrincipalAccountConditionWhenConditionIsSameAccountAsPrincipal", testPrincipalAccountConditionWhenConditionIsSameAccountAsPrincipalArn)
	t.Run("TestPrincipalAccountConditionWhenConditionIsSameAccountAsPrincipal", testPrincipalAccountConditionWhenConditionIsDifferentAccountToPrincipalArn)
	t.Run("TestPrincipalAccountConditionWhenPrincipalIsArnAndConditionIsWildcarded", testPrincipalAccountConditionIsWildcardedAgainstAPrincipalArn)

	// StringEquals
	t.Run("TestPrincipalAccountConditionWhenValueIsAUserAccountUsingStringEquals", testPrincipalAccountConditionWhenValueIsAUserAccountUsingStringEquals)
	t.Run("TestPrincipalAccountConditionWhenValueIsACrossAccountUsingStringEquals", testPrincipalAccountConditionWhenValueIsACrossAccountUsingStringEquals)
	t.Run("TestPrincipalAccountConditionWhenValueIsFullWildcardUsingStringEquals", testPrincipalAccountConditionWhenValueIsFullWildcardUsingStringEquals)
	t.Run("TestPrincipalAccountConditionWhenValueIsPartialWildcardUsingStringEquals", testPrincipalAccountConditionWhenValueIsPartialWildcardUsingStringEquals)
	t.Run("TestPrincipalAccountConditionUsingStringEqualsIfExists", testPrincipalAccountConditionUsingStringEqualsIfExists)
	// StringNotEquals
	// StringEqualsIgnoreCase
	t.Run("TestPrincipalAccountConditionWhenValueIsAUserAccountUsingStringEqualsIgnoreCase", testPrincipalAccountConditionWhenValueIsAUserAccountUsingStringEqualsIgnoreCase)
	t.Run("TestPrincipalAccountConditionWhenValueIsACrossAccountUsingStringEqualsIgnoreCase", testPrincipalAccountConditionWhenValueIsACrossAccountUsingStringEqualsIgnoreCase)
	t.Run("TestPrincipalAccountConditionWhenValueIsFullWildcardUsingStringEqualsIgnoreCase", testPrincipalAccountConditionWhenValueIsFullWildcardUsingStringEqualsIgnoreCase)
	t.Run("TestPrincipalAccountConditionWhenValueIsPartialWildcardUsingStringEqualsIgnoreCase", testPrincipalAccountConditionWhenValueIsPartialWildcardUsingStringEqualsIgnoreCase)
	t.Run("TestPrincipalAccountConditionUsingStringEqualsIgnoreCaseIfExists", testPrincipalAccountConditionUsingStringEqualsIgnoreCaseIfExists)
	// StringNotEqualsIgnoreCase
	// StringLike
	t.Run("TestPrincipalAccountConditionWhenValueIsAUserAccountWithStringLike", testPrincipalAccountConditionWhenValueIsAUserAccountWithStringLike)
	t.Run("TestPrincipalAccountConditionWhenValueIsACrossAccountWithStringLike", testPrincipalAccountConditionWhenValueIsACrossAccountWithStringLike)
	t.Run("TestPrincipalAccountConditionWhenValueIsAnAccountWithOneDigitTooFewWithStringLike", testPrincipalAccountConditionWhenValueIsAnAccountWithOneDigitTooFewWithStringLike)
	t.Run("TestPrincipalAccountConditionWhenValueIsAnAccountWithOneDigitTooManyWithStringLike", testPrincipalAccountConditionWhenValueIsAnAccountWithOneDigitTooManyWithStringLike)
	t.Run("TestPrincipalAccountConditionWhenValueIsFullWildcardWithStringLike", testPrincipalAccountConditionWhenValueIsFullWildcardWithStringLike)
	t.Run("TestPrincipalAccountConditionWhenValueIsPartialWildcardWithStringLike", testPrincipalAccountConditionWhenValueIsPartialWildcardWithStringLike)
	t.Run("TestPrincipalAccountConditionUsingStringLikeIfExists", testPrincipalAccountConditionUsingStringLikeIfExists)
	t.Run("TestPrincipalAccountConditionWhenValueWhenAccountIsSingleWildcardedUsingStringLike", testPrincipalAccountConditionWhenValueWhenAccountIsSingleWildcardedUsingStringLike)
	t.Run("TestPrincipalAccountConditionWhenValueWhenAccountIsWildcardedOneTooFewUsingStringLike", testPrincipalAccountConditionWhenValueWhenAccountIsWildcardedOneTooFewUsingStringLike)
	t.Run("TestPrincipalAccountConditionWhenValueWhenAccountIsWildcardedOneTooManyUsingStringLike", testPrincipalAccountConditionWhenValueWhenAccountIsWildcardedOneTooManyUsingStringLike)

	// StringNotLike
	// String Other
	t.Run("TestPrincipalAccountConditionWhenValueWhenArnIsMalformedUsingStringOperators", testPrincipalAccountConditionWhenValueWhenArnIsMalformedUsingStringOperators)
	t.Run("TestPrincipalAccountConditionWhenValueWhenAccountIsOneDigitTooFewUsingStringOperators", testPrincipalAccountConditionWhenValueWhenAccountIsOneDigitTooFewUsingStringOperators)
	t.Run("TestPrincipalAccountConditionWhenValueWhenAccountIsOneDigitTooManyUsingStringOperators", testPrincipalAccountConditionWhenValueWhenAccountIsOneDigitTooManyUsingStringOperators)
	t.Run("TestPrincipalAccountConditionWithMultipleValuesUsingStringOperators", testPrincipalAccountConditionWithMultipleValuesUsingStringOperators)

	// Others
	t.Run("TestPrincipalAccountConditionIsNotAnArnOrStringType", testPrincipalAccountConditionIsNotAStringType)
	t.Run("TestPrincipalAccountConditionWhenUnknownOperatorType", testPrincipalAccountConditionWhenUnknownOperatorType)
	t.Run("TestPrincipalAccountConditionWhenAcrossMultipleStatements", testPrincipalAccountConditionWhenAcrossMultipleStatements)
}

func testPrincipalAccountConditionWhenConditionIsSameAccountAsPrincipal(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["012345678901"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalAccountConditionWhenConditionIsDifferentAccountToPrincipal(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["012345678901"]
            }
          },
          "Principal": {
            "AWS": "222244442222"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalAccountConditionWhenConditionIsWildcardedAndPrincipalIsFixed(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalAccount": ["*"]
            }
          },
          "Principal": {
            "AWS": "222244442222"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"222244442222"},
		AllowedPrincipalAccountIds:          []string{"222244442222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalAccountConditionWhenConditionIsFixedAndPrincipalIsWildcarded(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["012345678901"]
            }
          },
          "Principal": {
            "AWS": "*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalAccountConditionWhenConditionIsSameAccountAsPrincipalArn(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["012345678901"]
            }
          },
          "Principal": {
            "AWS": "arn:aws:iam::012345678901:root"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::012345678901:root"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalAccountConditionWhenConditionIsDifferentAccountToPrincipalArn(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["012345678901"]
            }
          },
          "Principal": {
            "AWS": "arn:aws:iam::222244446666:root"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalAccountConditionIsWildcardedAgainstAPrincipalArn(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalAccount": ["*"]
            }
          },
          "Principal": {
            "AWS": "arn:aws:iam::012345678901:root"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::012345678901:root"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalAccountConditionWhenValueIsAUserAccountUsingStringEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["012345678901"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalAccountConditionWhenValueIsACrossAccountUsingStringEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["222233332222"]
            }
          },
          "Principal": {
            "AWS": "222233332222"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"222233332222"},
		AllowedPrincipalAccountIds:          []string{"222233332222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalAccountConditionWhenValueIsFullWildcardUsingStringEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["*"]
            }
          },
          "Principal": {
            "AWS": "*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalAccountConditionWhenValueIsPartialWildcardUsingStringEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["12345678*", "1234567890??"]
            }
          },
          "Principal": {
            "AWS": "*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalAccountConditionUsingStringEqualsIfExists(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIfExists": {
              "aws:PrincipalAccount": ["012345678901"]
            }
          },
          "Principal": {
            "AWS": "*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalAccountConditionWhenValueIsAUserAccountUsingStringEqualsIgnoreCase(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIgnoreCase": {
              "aws:PrincipalAccount": ["012345678901"]
            }
          },
          "Principal": {
            "AWS": "*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalAccountConditionWhenValueIsACrossAccountUsingStringEqualsIgnoreCase(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIgnoreCase": {
              "aws:PrincipalAccount": ["222233332222"]
            }
          },
          "Principal": {
            "AWS": "*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"222233332222"},
		AllowedPrincipalAccountIds:          []string{"222233332222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalAccountConditionWhenValueIsFullWildcardUsingStringEqualsIgnoreCase(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIgnoreCase": {
              "aws:PrincipalAccount": ["012345678901"]
            }
          },
          "Principal": {
            "AWS": "*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalAccountConditionWhenValueIsPartialWildcardUsingStringEqualsIgnoreCase(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIgnoreCase": {
              "aws:PrincipalAccount": ["12345678*", "1234567890??"]
            }
          },
          "Principal": {
            "AWS": "*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalAccountConditionUsingStringEqualsIgnoreCaseIfExists(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIgnoreCaseIfExists": {
              "aws:PrincipalAccount": ["012345678901"]
            }
          },
          "Principal": {
            "AWS": "*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalAccountConditionWhenValueIsAUserAccountWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalAccount": ["012345678901"]
            }
          },
          "Principal": {
            "AWS": "*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalAccountConditionWhenValueIsACrossAccountWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalAccount": ["222233332222"]
            }
          },
          "Principal": {
            "AWS": "*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"222233332222"},
		AllowedPrincipalAccountIds:          []string{"222233332222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

// NOTE: I think we should test this one, I believe this value is incorrect, should be private as there is no account with this number
func testPrincipalAccountConditionWhenValueIsAnAccountWithOneDigitTooFewWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalAccount": ["22223333222"]
            }
          },
          "Principal": {
            "AWS": "*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

// NOTE: I think we should test this one, I believe this value is incorrect, should be private as there is no account with this number
func testPrincipalAccountConditionWhenValueIsAnAccountWithOneDigitTooManyWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalAccount": ["2222333322222"]
            }
          },
          "Principal": {
            "AWS": "*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalAccountConditionWhenValueIsFullWildcardWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalAccount": ["*"]
            }
          },
          "Principal": {
            "AWS": "*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalAccountConditionWhenValueIsPartialWildcardWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalAccount": ["1234*"]
            }
          },
          "Principal": {
            "AWS": "*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"1234*"},
		AllowedPrincipalAccountIds:          []string{"1234*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalAccountConditionUsingStringLikeIfExists(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLikeIfExists": {
              "aws:PrincipalAccount": ["222233332222"]
            }
          },
          "Principal": {
            "AWS": "*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalAccountConditionWhenValueWhenAccountIsSingleWildcardedUsingStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalAccount": ["0123456789??"]
            }
          },
          "Principal": {
            "AWS": "*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"0123456789??"},
		AllowedPrincipalAccountIds:          []string{"0123456789??"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

// NOTE: I think we should test this one, I believe this value is incorrect, should be private as there is no account with this number
func testPrincipalAccountConditionWhenValueWhenAccountIsWildcardedOneTooFewUsingStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalAccount": ["0123456789?"]
            }
          },
          "Principal": {
            "AWS": "*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

// NOTE: I think we should test this one, I believe this value is incorrect, should be private as there is no account with this number
func testPrincipalAccountConditionWhenValueWhenAccountIsWildcardedOneTooManyUsingStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLike": {
              "aws:PrincipalAccount": ["0123456789???"]
            }
          },
          "Principal": {
            "AWS": "*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalAccountConditionWhenValueWhenArnIsMalformedUsingStringOperators(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["01234567890A"]
            }
          },
          "Principal": {
            "AWS": "*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalAccountConditionWhenValueWhenAccountIsOneDigitTooFewUsingStringOperators(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["01234567890"]
            }
          },
          "Principal": {
            "AWS": "*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalAccountConditionWhenValueWhenAccountIsOneDigitTooManyUsingStringOperators(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["0123456789012"]
            }
          },
          "Principal": {
            "AWS": "*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}
	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalAccountConditionWithMultipleValuesUsingStringOperators(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": [
                "012345678901",
                "222233332222"
              ]
            }
          },
          "Principal": {
            "AWS": "*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"012345678901",
			"222233332222",
		},
		AllowedPrincipalAccountIds: []string{
			"012345678901",
			"222233332222",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalAccountConditionIsNotAStringType(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnEquals": {
              "aws:PrincipalAccount": ["222233332222"]
            }
          },
          "Principal": {
            "AWS": "*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalAccountConditionWhenUnknownOperatorType(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringUnknown": {
              "aws:PrincipalAccount": ["222233332222"]
            }
          },
          "Principal": {
            "AWS": "*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalAccountConditionWhenAcrossMultipleStatements(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["012345678901"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalAccount": ["*"]
            }
          },
          "Principal": {
            "AWS": "*"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["222233332222"]
            }
          },
          "Principal": {
            "AWS": "222233332222"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "public",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"*",
			"012345678901",
			"222233332222",
		},
		AllowedPrincipalAccountIds: []string{
			"*",
			"012345678901",
			"222233332222",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{"Statement[2]"},
		SharedStatementIds:                  []string{"Statement[3]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func TestGlobalConditionPrincipalOrgID(t *testing.T) {
	// StringEquals
	t.Run("TestPrincipalOrgIDConditionWhenValueIsAValidOrgIDUsingStringEquals", testPrincipalOrgIDConditionWhenValueIsAValidOrgIDUsingStringEquals)
	t.Run("TestPrincipalOrgIDConditionWhenValueIsAnInvalidOrgIDUsingStringEquals", testPrincipalOrgIDConditionWhenValueIsAnInvalidOrgIDUsingStringEquals)
	t.Run("TestPrincipalOrgIDConditionWhenValueIsFullWildcardUsingStringEquals", testPrincipalOrgIDConditionWhenValueIsFullWildcardUsingStringEquals)
	t.Run("TestPrincipalOrgIDConditionWhenValueIsPartialWildcardUsingStringEquals", testPrincipalOrgIDConditionWhenValueIsPartialWildcardUsingStringEquals)
	t.Run("TestPrincipalOrgIDConditionUsingStringEqualsIfExists", testPrincipalOrgIDConditionUsingStringEqualsIfExists)
	// StringNotEquals
	// StringEqualsIgnoreCase
	t.Run("TestPrincipalOrgIDConditionWhenValueIsAValidOrgIDUsingStringEqualsIgnoreCase", testPrincipalOrgIDConditionWhenValueIsAValidOrgIDUsingStringEqualsIgnoreCase)
	t.Run("TestPrincipalOrgIDConditionWhenValueIsAnInvalidOrgIDUsingStringEqualsIgnoreCase", testPrincipalOrgIDConditionWhenValueIsAnInvalidOrgIDUsingStringEqualsIgnoreCase)
	t.Run("TestPrincipalOrgIDConditionWhenValueIsFullWildcardUsingStringEqualsIgnoreCase", testPrincipalOrgIDConditionWhenValueIsFullWildcardUsingStringEqualsIgnoreCase)
	t.Run("TestPrincipalOrgIDConditionWhenValueIsPartialWildcardUsingStringEqualsIgnoreCase", testPrincipalOrgIDConditionWhenValueIsPartialWildcardUsingStringEqualsIgnoreCase)
	t.Run("TestPrincipalOrgIDConditionUsingStringEqualsIgnoreCaseIfExists", testPrincipalOrgIDConditionUsingStringEqualsIgnoreCaseIfExists)
	// StringNotEqualsIgnoreCase
	// StringLike
	t.Run("TestPrincipalOrgIDConditionWhenValueIsAValidOrgIDUsingStringLike", testPrincipalOrgIDConditionWhenValueIsAValidOrgIDUsingStringLike)
	t.Run("TestPrincipalOrgIDConditionWhenValueIsAnInvalidOrgIDUsingStringLike", testPrincipalOrgIDConditionWhenValueIsAnInvalidOrgIDUsingStringLike)
	t.Run("TestPrincipalOrgIDConditionWhenValueIsAllOrganizationsUsingStringLike", testPrincipalOrgIDConditionWhenValueIsAllOrganizationsUsingStringLike)
	t.Run("TestPrincipalOrgIDConditionWhenValueIsFullWildcardUsingStringLike", testPrincipalOrgIDConditionWhenValueIsFullWildcardUsingStringLike)
	t.Run("TestPrincipalOrgIDConditionWhenValueIsWildcardOrganizationUsingStringLike", testPrincipalOrgIDConditionWhenValueIsWildcardOrganizationUsingStringLike)
	t.Run("TestPrincipalOrgIDConditionWhenValueIsInvalidWildcardOrganizationUsingStringLike", testPrincipalOrgIDConditionWhenValueIsInvalidWildcardOrganizationUsingStringLike)
	t.Run("TestPrincipalOrgIDConditionUsingStringLikeIfExists", testPrincipalOrgIDConditionUsingStringLikeIfExists)

	// StringNotLike
	// String Other
	t.Run("TestPrincipalOrgIDConditionWithMultipleValuesUsingStringOperators", testPrincipalOrgIDConditionWithMultipleValuesUsingStringOperators)

	// Others
	t.Run("TestPrincipalOrgIDConditionIsNotAStringType", testPrincipalOrgIDConditionIsNotAStringType)
	t.Run("TestPrincipalOrgIDConditionWhenUnknownOperatorType", testPrincipalOrgIDConditionWhenUnknownOperatorType)
	t.Run("TestPrincipalOrgIDConditionWhenAcrossMultipleStatements", testPrincipalOrgIDConditionWhenAcrossMultipleStatements)
	t.Run("TestPrincipalOrgIDConditionWhenPrincipalIsAnArn", testPrincipalOrgIDConditionWhenPrincipalIsAnArn)
	t.Run("TestPrincipalOrgIDConditionWhenPrincipalIsPublicAndHasAnOrganizationThenSharedAccess", testPrincipalOrgIDConditionWhenPrincipalIsPublicAndHasAnOrganizationThenSharedAccess)
}

func testPrincipalOrgIDConditionWhenValueIsAValidOrgIDUsingStringEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalOrgID": ["o-kkklllmmm555"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{"o-kkklllmmm555"},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalOrgIDConditionWhenValueIsAnInvalidOrgIDUsingStringEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalOrgID": ["invalid"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalOrgIDConditionWhenValueIsFullWildcardUsingStringEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalOrgID": ["*"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalOrgIDConditionWhenValueIsPartialWildcardUsingStringEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalOrgID": ["o-invalid*", "o-invalid?"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalOrgIDConditionUsingStringEqualsIfExists(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIfExists": {
              "aws:PrincipalOrgID": ["o-kkklllmmm555"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalOrgIDConditionWhenValueIsAValidOrgIDUsingStringEqualsIgnoreCase(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIgnoreCase": {
              "aws:PrincipalOrgID": ["o-kkklllmmm555"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{"o-kkklllmmm555"},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalOrgIDConditionWhenValueIsAnInvalidOrgIDUsingStringEqualsIgnoreCase(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIgnoreCase": {
              "aws:PrincipalOrgID": ["invalid"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalOrgIDConditionWhenValueIsFullWildcardUsingStringEqualsIgnoreCase(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIgnoreCase": {
              "aws:PrincipalOrgID": ["*"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalOrgIDConditionWhenValueIsPartialWildcardUsingStringEqualsIgnoreCase(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIgnoreCase": {
              "aws:PrincipalOrgID": ["o-invalid*", "o-invalid?"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalOrgIDConditionUsingStringEqualsIgnoreCaseIfExists(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIgnoreCaseIfExists": {
              "aws:PrincipalOrgID": ["o-kkklllmmm555"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalOrgIDConditionWhenValueIsAValidOrgIDUsingStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalOrgID": ["o-kkklllmmm555"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{"o-kkklllmmm555"},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalOrgIDConditionWhenValueIsAnInvalidOrgIDUsingStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalOrgID": ["invalid"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalOrgIDConditionWhenValueIsAllOrganizationsUsingStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalOrgID": ["o-*"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{"o-*"},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalOrgIDConditionWhenValueIsFullWildcardUsingStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalOrgID": ["*"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel: "public",
		AllowedOrganizationIds: []string{
			"o-*",
		},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalOrgIDConditionWhenValueIsWildcardOrganizationUsingStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalOrgID": ["o-kkklllmmm555*"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{"o-kkklllmmm555*"},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalOrgIDConditionWhenValueIsInvalidWildcardOrganizationUsingStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalOrgID": ["a*"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalOrgIDConditionUsingStringLikeIfExists(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLikeIfExists": {
              "aws:PrincipalOrgID": ["o-kkklllmmm555"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalOrgIDConditionWithMultipleValuesUsingStringOperators(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalOrgID": [
                "o-ggghhhiii555",
                "o-*",
                "o-kkklllmmm555"
              ]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel: "shared",
		AllowedOrganizationIds: []string{
			"o-ggghhhiii555",
			"o-kkklllmmm555",
		},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalOrgIDConditionIsNotAStringType(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnEquals": {
              "aws:PrincipalOrgID": ["222233332222"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalOrgIDConditionWhenUnknownOperatorType(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringUnknown": {
              "aws:PrincipalOrgID": ["222233332222"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalOrgIDConditionWhenAcrossMultipleStatements(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalOrgID": ["o-zzzyyyxxx987"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalOrgID": ["*"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalOrgID": ["o-aaabbbccc123"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel: "public",
		AllowedOrganizationIds: []string{
			"o-*",
			"o-aaabbbccc123",
			"o-zzzyyyxxx987",
		},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[2]"},
		SharedStatementIds: []string{
			"Statement[1]",
			"Statement[3]",
		},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalOrgIDConditionWhenPrincipalIsAnArn(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalOrgID": ["o-zzzyyyxxx987"]
            }
          },
          "Principal": {
            "AWS": "arn:aws:iam::012345678901:root"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalOrgID": ["*"]
            }
          },
          "Principal": {
            "AWS": "arn:aws:iam::012345678901:root"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalOrgID": ["o-aaabbbccc123"]
            }
          },
          "Principal": {
            "AWS": "arn:aws:iam::012345678901:root"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel: "public",
		AllowedOrganizationIds: []string{
			"o-*",
			"o-aaabbbccc123",
			"o-zzzyyyxxx987",
		},
		AllowedPrincipals:                   []string{"arn:aws:iam::012345678901:root"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[2]"},
		SharedStatementIds: []string{
			"Statement[1]",
			"Statement[3]",
		},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testPrincipalOrgIDConditionWhenPrincipalIsPublicAndHasAnOrganizationThenSharedAccess(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Sid": "AllowedOrganization",
          "Effect": "Allow",
          "Principal": {
            "AWS": "*"
          },
          "Action": "ec2:DescribeVolumes",
          "Resource": "arn:aws:s3:::omero-resource-policy-bucket",
          "Condition": { "StringEquals": { "aws:PrincipalOrgID": "o-aaabbbccc123" } }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel: "shared",
		AllowedOrganizationIds: []string{
			"o-aaabbbccc123",
		},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"AllowedOrganization"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}

}

func TestGlobalConditionSourceOwner(t *testing.T) {
	t.Run("TestSourceOwnerConditionWhenConditionIsSameAccountAsPrincipal", testSourceOwnerConditionWhenConditionIsSameAccountAsPrincipal)
	t.Run("TestSourceOwnerConditionWhenConditionIsDifferentAccountToPrincipal", testSourceOwnerConditionWhenConditionIsDifferentAccountToPrincipal)
	t.Run("TestSourceOwnerConditionWhenConditionIsWildcardedAndPrincipalIsFixed", testSourceOwnerConditionWhenConditionIsWildcardedAndPrincipalIsFixed)
	t.Run("TestSourceOwnerConditionWhenConditionIsFixedAndPrincipalIsWildcarded", testSourceOwnerConditionWhenConditionIsFixedAndPrincipalIsWildcarded)

	t.Run("TestSourceOwnerConditionWhenConditionIsSameAccountAsPrincipal", testSourceOwnerConditionWhenConditionIsSameAccountAsPrincipalArn)
	t.Run("TestSourceOwnerConditionWhenConditionIsSameAccountAsPrincipal", testSourceOwnerConditionWhenConditionIsDifferentAccountToPrincipalArn)
	t.Run("TestSourceOwnerConditionWhenPrincipalIsArnAndConditionIsWildcarded", testSourceOwnerConditionIsWildcardedAgainstAPrincipalArn)

	// StringEquals
	t.Run("TestSourceOwnerConditionWhenValueIsAUserAccountUsingStringEquals", testSourceOwnerConditionWhenValueIsAUserAccountUsingStringEquals)
	t.Run("TestSourceOwnerConditionWhenValueIsACrossAccountUsingStringEquals", testSourceOwnerConditionWhenValueIsACrossAccountUsingStringEquals)
	t.Run("TestSourceOwnerConditionWhenValueIsFullWildcardUsingStringEquals", testSourceOwnerConditionWhenValueIsFullWildcardUsingStringEquals)
	t.Run("TestSourceOwnerConditionWhenValueIsPartialWildcardUsingStringEquals", testSourceOwnerConditionWhenValueIsPartialWildcardUsingStringEquals)
	t.Run("TestSourceOwnerConditionUsingStringEqualsIfExists", testSourceOwnerConditionUsingStringEqualsIfExists)
	// StringNotEquals
	// StringEqualsIgnoreCase
	t.Run("TestSourceOwnerConditionWhenValueIsAUserAccountUsingStringEqualsIgnoreCase", testSourceOwnerConditionWhenValueIsAUserAccountUsingStringEqualsIgnoreCase)
	t.Run("TestSourceOwnerConditionWhenValueIsACrossAccountUsingStringEqualsIgnoreCase", testSourceOwnerConditionWhenValueIsACrossAccountUsingStringEqualsIgnoreCase)
	t.Run("TestSourceOwnerConditionWhenValueIsFullWildcardUsingStringEqualsIgnoreCase", testSourceOwnerConditionWhenValueIsFullWildcardUsingStringEqualsIgnoreCase)
	t.Run("TestSourceOwnerConditionWhenValueIsPartialWildcardUsingStringEqualsIgnoreCase", testSourceOwnerConditionWhenValueIsPartialWildcardUsingStringEqualsIgnoreCase)
	t.Run("TestSourceOwnerConditionUsingStringEqualsIgnoreCaseIfExists", testSourceOwnerConditionUsingStringEqualsIgnoreCaseIfExists)
	// StringNotEqualsIgnoreCase
	// StringLike
	t.Run("TestSourceOwnerConditionWhenValueIsAUserAccountWithStringLike", testSourceOwnerConditionWhenValueIsAUserAccountWithStringLike)
	t.Run("TestSourceOwnerConditionWhenValueIsACrossAccountWithStringLike", testSourceOwnerConditionWhenValueIsACrossAccountWithStringLike)
	t.Run("TestSourceOwnerConditionWhenValueIsAnAccountWithOneDigitTooFewWithStringLike", testSourceOwnerConditionWhenValueIsAnAccountWithOneDigitTooFewWithStringLike)
	t.Run("TestSourceOwnerConditionWhenValueIsAnAccountWithOneDigitTooManyWithStringLike", testSourceOwnerConditionWhenValueIsAnAccountWithOneDigitTooManyWithStringLike)
	t.Run("TestSourceOwnerConditionWhenValueIsFullWildcardWithStringLike", testSourceOwnerConditionWhenValueIsFullWildcardWithStringLike)
	t.Run("TestSourceOwnerConditionWhenValueIsPartialWildcardWithStringLike", testSourceOwnerConditionWhenValueIsPartialWildcardWithStringLike)
	t.Run("TestSourceOwnerConditionUsingStringLikeIfExists", testSourceOwnerConditionUsingStringLikeIfExists)
	t.Run("TestSourceOwnerConditionWhenValueWhenAccountIsSingleWildcardedUsingStringLike", testSourceOwnerConditionWhenValueWhenAccountIsSingleWildcardedUsingStringLike)
	t.Run("TestSourceOwnerConditionWhenValueWhenAccountIsWildcardedOneTooFewUsingStringLike", testSourceOwnerConditionWhenValueWhenAccountIsWildcardedOneTooFewUsingStringLike)
	t.Run("TestSourceOwnerConditionWhenValueWhenAccountIsWildcardedOneTooManyUsingStringLike", testSourceOwnerConditionWhenValueWhenAccountIsWildcardedOneTooManyUsingStringLike)

	// StringNotLike
	// String Other
	t.Run("TestSourceOwnerConditionWhenValueWhenArnIsMalformedUsingStringOperators", testSourceOwnerConditionWhenValueWhenArnIsMalformedUsingStringOperators)
	t.Run("TestSourceOwnerConditionWhenValueWhenAccountIsOneDigitTooFewUsingStringOperators", testSourceOwnerConditionWhenValueWhenAccountIsOneDigitTooFewUsingStringOperators)
	t.Run("TestSourceOwnerConditionWhenValueWhenAccountIsOneDigitTooManyUsingStringOperators", testSourceOwnerConditionWhenValueWhenAccountIsOneDigitTooManyUsingStringOperators)
	t.Run("TestSourceOwnerConditionWithMultipleValuesUsingStringOperators", testSourceOwnerConditionWithMultipleValuesUsingStringOperators)

	// Others
	t.Run("TestSourceOwnerConditionIsNotAnArnOrStringType", TestSourceOwnerConditionIsNotAStringType)
	t.Run("TestSourceOwnerConditionWhenUnknownOperatorType", testSourceOwnerConditionWhenUnknownOperatorType)
	t.Run("TestSourceOwnerConditionWhenAcrossMultipleStatements", testSourceOwnerConditionWhenAcrossMultipleStatements)
}

func testSourceOwnerConditionWhenConditionIsSameAccountAsPrincipal(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceOwner": ["012345678901"]
            }
          },
          "Principal": {
            "Service": ["ecs.amazonaws.com"]
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceOwnerConditionWhenConditionIsDifferentAccountToPrincipal(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceOwner": ["012345678901"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceOwnerConditionWhenConditionIsWildcardedAndPrincipalIsFixed(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:SourceOwner": ["*"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceOwnerConditionWhenConditionIsFixedAndPrincipalIsWildcarded(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceOwner": ["012345678901"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceOwnerConditionWhenConditionIsSameAccountAsPrincipalArn(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceOwner": ["012345678901"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceOwnerConditionWhenConditionIsDifferentAccountToPrincipalArn(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceOwner": ["012345678901"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceOwnerConditionIsWildcardedAgainstAPrincipalArn(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:SourceOwner": ["*"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceOwnerConditionWhenValueIsAUserAccountUsingStringEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceOwner": ["012345678901"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceOwnerConditionWhenValueIsACrossAccountUsingStringEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceOwner": ["222233332222"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"222233332222"},
		AllowedPrincipalAccountIds:          []string{"222233332222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceOwnerConditionWhenValueIsFullWildcardUsingStringEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceOwner": ["*"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceOwnerConditionWhenValueIsPartialWildcardUsingStringEquals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceOwner": ["12345678*", "1234567890??"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceOwnerConditionUsingStringEqualsIfExists(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIfExists": {
              "aws:SourceOwner": ["012345678901"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceOwnerConditionWhenValueIsAUserAccountUsingStringEqualsIgnoreCase(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIgnoreCase": {
              "aws:SourceOwner": ["012345678901"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceOwnerConditionWhenValueIsACrossAccountUsingStringEqualsIgnoreCase(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIgnoreCase": {
              "aws:SourceOwner": ["222233332222"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"222233332222"},
		AllowedPrincipalAccountIds:          []string{"222233332222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceOwnerConditionWhenValueIsFullWildcardUsingStringEqualsIgnoreCase(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIgnoreCase": {
              "aws:SourceOwner": ["012345678901"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceOwnerConditionWhenValueIsPartialWildcardUsingStringEqualsIgnoreCase(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIgnoreCase": {
              "aws:SourceOwner": ["12345678*", "1234567890??"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceOwnerConditionUsingStringEqualsIgnoreCaseIfExists(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEqualsIgnoreCaseIfExists": {
              "aws:SourceOwner": ["012345678901"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceOwnerConditionWhenValueIsAUserAccountWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:SourceOwner": ["012345678901"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{"012345678901"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceOwnerConditionWhenValueIsACrossAccountWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:SourceOwner": ["222233332222"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"222233332222"},
		AllowedPrincipalAccountIds:          []string{"222233332222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceOwnerConditionWhenValueIsAnAccountWithOneDigitTooFewWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:SourceOwner": ["22223333222"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceOwnerConditionWhenValueIsAnAccountWithOneDigitTooManyWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:SourceOwner": ["2222333322222"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceOwnerConditionWhenValueIsFullWildcardWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:SourceOwner": ["*"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceOwnerConditionWhenValueIsPartialWildcardWithStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:SourceOwner": ["1234*"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"1234*"},
		AllowedPrincipalAccountIds:          []string{"1234*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceOwnerConditionUsingStringLikeIfExists(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLikeIfExists": {
              "aws:SourceOwner": ["222233332222"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceOwnerConditionWhenValueWhenAccountIsSingleWildcardedUsingStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:SourceOwner": ["0123456789??"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"0123456789??"},
		AllowedPrincipalAccountIds:          []string{"0123456789??"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceOwnerConditionWhenValueWhenAccountIsWildcardedOneTooFewUsingStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:SourceOwner": ["0123456789?"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceOwnerConditionWhenValueWhenAccountIsWildcardedOneTooManyUsingStringLike(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnLike": {
              "aws:SourceOwner": ["0123456789???"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceOwnerConditionWhenValueWhenArnIsMalformedUsingStringOperators(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceOwner": ["01234567890A"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceOwnerConditionWhenValueWhenAccountIsOneDigitTooFewUsingStringOperators(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceOwner": ["01234567890"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceOwnerConditionWhenValueWhenAccountIsOneDigitTooManyUsingStringOperators(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceOwner": ["0123456789012"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}
	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceOwnerConditionWithMultipleValuesUsingStringOperators(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceOwner": [
                "012345678901",
                "222233332222"
              ]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"012345678901",
			"222233332222",
		},
		AllowedPrincipalAccountIds: []string{
			"012345678901",
			"222233332222",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func TestSourceOwnerConditionIsNotAStringType(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "ArnEquals": {
              "aws:SourceOwner": ["222233332222"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceOwnerConditionWhenUnknownOperatorType(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringUnknown": {
              "aws:SourceOwner": ["222233332222"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testSourceOwnerConditionWhenAcrossMultipleStatements(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceOwner": ["012345678901"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:SourceOwner": ["*"]
            }
          },
          "Principal": {
            "AWS": "*"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:SourceOwner": ["222233332222"]
            }
          },
		  "Principal": {
			"Service": ["ecs.amazonaws.com"]
		  }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "public",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"*",
			"012345678901",
			"222233332222",
		},
		AllowedPrincipalAccountIds: []string{
			"*",
			"012345678901",
			"222233332222",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{"Statement[2]"},
		SharedStatementIds:                  []string{"Statement[3]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func TestDenyPermissionsByAccount(t *testing.T) {
	t.Run("TestDenyPermissionsByAccountRemovesPrincipalWithRespectiveDeny", testDenyPermissionsByAccountRemovesPrincipalWithRespectiveDeny)
	t.Run("TestDenyPermissionsByAccountRemovesCorrectPrincipalWithRespectiveDeny", testDenyPermissionsByAccountRemovesCorrectPrincipalWithRespectiveDeny)
	t.Run("TestDenyPermissionsByAccountRemovesCorrectPrincipalsWithRespectiveDenies", testDenyPermissionsByAccountRemovesCorrectPrincipalsWithRespectiveDenies)
	t.Run("TestDenyPermissionsByAccountRemovesCorrectPrincipalsWhenDenyingMultiplePermissions", testDenyPermissionsByAccountRemovesCorrectPrincipalsWhenDenyingMultiplePermissions)

	t.Run("TestDenyPermissionsByAccountRemovesCorrectPrincipalsWhenDenyWildcardPermissions", testDenyPermissionsByAccountRemovesCorrectPrincipalsWhenDenyWildcardPermissions)
	t.Run("TestDenyPermissionsByAccountRemovesAllPrincipalsWhenDenyHasMultiplPrincipals", testDenyPermissionsByAccountRemovesAllPrincipalsWhenDenyHasMultiplPrincipals)
	t.Run("TestDenyPermissionsByAccountMultiplePermissionsWithMultiplePrincipalsAndDenyOnePermissionsFromEach", testDenyPermissionsByAccountMultiplePermissionsWithMultiplePrincipalsAndDenyOnePermissionsFromEach)
	t.Run("TestDenyPermissionsByAccountFullWildcardPrincipalThatFullyContainsAllAllowPermissionsDeniesAll", testDenyPermissionsByAccountFullWildcardPrincipalThatFullyContainsAllAllowPermissionsDeniesAll)
	t.Run("TestDenyPermissionsByAccountWhereDenyHasPartiallyWildcardedPrincipalsForAccounts", testDenyPermissionsByAccountWhereDenyHasPartiallyWildcardedPrincipalsForAccounts)

	t.Run("TestDenyPermissionsByAccountAllowPermissionsIsTheSupersetOfDenyResources", testDenyPermissionsByAccountAllowPermissionsIsTheSupersetOfDenyResources)
	t.Run("TestDenyPermissionsByAccountAllowPermissionsIsTheSubsetOfDenyResources", testDenyPermissionsByAccountAllowPermissionsIsTheSubsetOfDenyResources)
}

func testDenyPermissionsByAccountRemovesPrincipalWithRespectiveDeny(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "222244446666"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByAccountRemovesCorrectPrincipalWithRespectiveDeny(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "666644442222"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "666644442222"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"222244446666"},
		AllowedPrincipalAccountIds:          []string{"222244446666"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByAccountRemovesCorrectPrincipalsWithRespectiveDenies(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "666644442222"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "666644442222"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "222244446666"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByAccountRemovesCorrectPrincipalsWhenDenyingMultiplePermissions(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": ["222244446666", "666644442222"]
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumesModifications",
          "Resource": "*",
          "Principal": {
            "AWS": ["222244446666", "666644442222"]
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumesModifications",
          "Resource": "*",
          "Principal": {
            "AWS": "222244446666"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"666644442222"},
		AllowedPrincipalAccountIds:          []string{"666644442222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels: []string{
			"List",
			"Read",
		},
		PrivateAccessLevels: []string{},
		PublicStatementIds:  []string{},
		SharedStatementIds: []string{
			"Statement[1]",
			"Statement[2]",
		},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByAccountRemovesCorrectPrincipalsWhenDenyWildcardPermissions(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": ["222244446666", "666644442222"]
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumesModifications",
          "Resource": "*",
          "Principal": {
            "AWS": ["222244446666", "666644442222"]
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes*",
          "Resource": "*",
          "Principal": {
            "AWS": "222244446666"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"666644442222"},
		AllowedPrincipalAccountIds:          []string{"666644442222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels: []string{
			"List",
			"Read",
		},
		PrivateAccessLevels: []string{},
		PublicStatementIds:  []string{},
		SharedStatementIds: []string{
			"Statement[1]",
			"Statement[2]",
		},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByAccountRemovesAllPrincipalsWhenDenyHasMultiplPrincipals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": ["222244446666", "666644442222"]
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": ["222244446666", "666644442222"]
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByAccountMultiplePermissionsWithMultiplePrincipalsAndDenyOnePermissionsFromEach(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": ["222244446666", "666644442222"]
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumesModifications",
          "Resource": "*",
          "Principal": {
            "AWS": ["222244446666", "666644442222"]
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumesModifications",
          "Resource": "*",
          "Principal": {
            "AWS": "666644442222"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"222244446666",
			"666644442222",
		},
		AllowedPrincipalAccountIds: []string{
			"222244446666",
			"666644442222",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels: []string{
			"List",
			"Read",
		},
		PrivateAccessLevels: []string{},
		PublicStatementIds:  []string{},
		SharedStatementIds: []string{
			"Statement[1]",
			"Statement[2]",
		},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByAccountFullWildcardPrincipalThatFullyContainsAllAllowPermissionsDeniesAll(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByAccountWhereDenyHasPartiallyWildcardedPrincipalsForAccounts(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "012345678901"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "0123456789??"
          }
        },
		{
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "2222*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"012345678901",
			"222244446666",
		},
		AllowedPrincipalAccountIds: []string{
			"012345678901",
			"222244446666",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByAccountAllowPermissionsIsTheSupersetOfDenyResources(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "arn:*",
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Principal": {
            "AWS": "666644442222"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "arn:*",
          "Principal": {
            "AWS": "666644442222"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"222244446666",
			"666644442222",
		},
		AllowedPrincipalAccountIds: []string{
			"222244446666",
			"666644442222",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds: []string{
			"Statement[1]",
			"Statement[3]",
		},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByAccountAllowPermissionsIsTheSubsetOfDenyResources(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
		  "Resource": "arn:*",
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "arn:*",
          "Principal": {
            "AWS": "666644442222"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Principal": {
            "AWS": "666644442222"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func TestDenyPermissionsByArn(t *testing.T) {
	t.Run("TestDenyPermissionsByArnRemovesPrincipalWithRespectiveDeny", testDenyPermissionsByArnRemovesPrincipalWithRespectiveDeny)
	t.Run("TestDenyPermissionsByArnRemovesCorrectPrincipalWithRespectiveDeny", testDenyPermissionsByArnRemovesCorrectPrincipalWithRespectiveDeny)
	t.Run("TestDenyPermissionsByArnRemovesCorrectPrincipalsWithRespectiveDenies", testDenyPermissionsByArnRemovesCorrectPrincipalsWithRespectiveDenies)
	t.Run("TestDenyPermissionsByArnRemovesCorrectPrincipalsWhenDenyingMultiplePermissions", testDenyPermissionsByArnRemovesCorrectPrincipalsWhenDenyingMultiplePermissions)

	t.Run("TestDenyPermissionsByArnRemovesCorrectPrincipalsWhenDenyWildcardPermissions", testDenyPermissionsByArnRemovesCorrectPrincipalsWhenDenyWildcardPermissions)
	t.Run("TestDenyPermissionsByArnRemovesAllPrincipalsWhenDenyHasMultiplPrincipals", testDenyPermissionsByArnRemovesAllPrincipalsWhenDenyHasMultiplPrincipals)
	t.Run("TestDenyPermissionsByArnMultiplePermissionsWithMultiplePrincipalsAndDenyOnePermissionsFromEach", testDenyPermissionsByArnMultiplePermissionsWithMultiplePrincipalsAndDenyOnePermissionsFromEach)
	t.Run("TestDenyPermissionsByArnFullWildcardPrincipalThatFullyContainsAllAllowPermissionsDeniesAll", testDenyPermissionsByArnFullWildcardPrincipalThatFullyContainsAllAllowPermissionsDeniesAll)
	t.Run("TestDenyPermissionsByArnWhereDenyHasPartiallyWildcardedPrincipalsForAccounts", testDenyPermissionsByArnWhereDenyHasPartiallyWildcardedPrincipalsForAccounts)

	t.Run("TestDenyPermissionsByArnAllowPermissionsIsTheSupersetOfDenyResources", testDenyPermissionsByArnAllowPermissionsIsTheSupersetOfDenyResources)
	t.Run("TestDenyPermissionsByArnAllowPermissionsIsTheSubsetOfDenyResources", testDenyPermissionsByArnAllowPermissionsIsTheSubsetOfDenyResources)
}

func testDenyPermissionsByArnRemovesPrincipalWithRespectiveDeny(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "arn:aws:iam::222244446666:root"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "arn:aws:iam::222244446666:root"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByArnRemovesCorrectPrincipalWithRespectiveDeny(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "arn:aws:iam::222244446666:root"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "arn:aws:iam::666644442222:root"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "arn:aws:iam::666644442222:root"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::222244446666:root"},
		AllowedPrincipalAccountIds:          []string{"222244446666"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByArnRemovesCorrectPrincipalsWithRespectiveDenies(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "arn:aws:iam::222244446666:root"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "arn:aws:iam::666644442222:root"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "arn:aws:iam::666644442222:root"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "arn:aws:iam::222244446666:root"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByArnRemovesCorrectPrincipalsWhenDenyingMultiplePermissions(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": ["arn:aws:iam::222244446666:root", "arn:aws:iam::666644442222:root"]
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumesModifications",
          "Resource": "*",
          "Principal": {
            "AWS": ["arn:aws:iam::222244446666:root", "arn:aws:iam::666644442222:root"]
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "arn:aws:iam::222244446666:root"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumesModifications",
          "Resource": "*",
          "Principal": {
            "AWS": "arn:aws:iam::222244446666:root"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::666644442222:root"},
		AllowedPrincipalAccountIds:          []string{"666644442222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels: []string{
			"List",
			"Read",
		},
		PrivateAccessLevels: []string{},
		PublicStatementIds:  []string{},
		SharedStatementIds: []string{
			"Statement[1]",
			"Statement[2]",
		},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByArnRemovesCorrectPrincipalsWhenDenyWildcardPermissions(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": ["arn:aws:iam::222244446666:root", "arn:aws:iam::666644442222:root"]
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumesModifications",
          "Resource": "*",
          "Principal": {
            "AWS": ["arn:aws:iam::222244446666:root", "arn:aws:iam::666644442222:root"]
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes*",
          "Resource": "*",
          "Principal": {
            "AWS": "arn:aws:iam::222244446666:root"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::666644442222:root"},
		AllowedPrincipalAccountIds:          []string{"666644442222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels: []string{
			"List",
			"Read",
		},
		PrivateAccessLevels: []string{},
		PublicStatementIds:  []string{},
		SharedStatementIds: []string{
			"Statement[1]",
			"Statement[2]",
		},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByArnRemovesAllPrincipalsWhenDenyHasMultiplPrincipals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": ["arn:aws:iam::222244446666:root", "arn:aws:iam::666644442222:root"]
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": ["arn:aws:iam::222244446666:root", "arn:aws:iam::666644442222:root"]
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByArnMultiplePermissionsWithMultiplePrincipalsAndDenyOnePermissionsFromEach(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": ["arn:aws:iam::222244446666:root", "arn:aws:iam::666644442222:root"]
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumesModifications",
          "Resource": "*",
          "Principal": {
            "AWS": ["arn:aws:iam::222244446666:root", "arn:aws:iam::666644442222:root"]
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "arn:aws:iam::222244446666:root"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumesModifications",
          "Resource": "*",
          "Principal": {
            "AWS": "arn:aws:iam::666644442222:root"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:iam::222244446666:root",
			"arn:aws:iam::666644442222:root",
		},
		AllowedPrincipalAccountIds: []string{
			"222244446666",
			"666644442222",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels: []string{
			"List",
			"Read",
		},
		PrivateAccessLevels: []string{},
		PublicStatementIds:  []string{},
		SharedStatementIds: []string{
			"Statement[1]",
			"Statement[2]",
		},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByArnFullWildcardPrincipalThatFullyContainsAllAllowPermissionsDeniesAll(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "arn:aws:iam::222244446666:root"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByArnWhereDenyHasPartiallyWildcardedPrincipalsForAccounts(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "arn:aws:iam::222244446666:root"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "arn:aws:iam::012345678901:root"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "arn:aws:iam::012345678901:ro??"
          }
        },
		{
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "2222*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:iam::012345678901:root",
			"arn:aws:iam::222244446666:root",
		},
		AllowedPrincipalAccountIds: []string{
			"012345678901",
			"222244446666",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByArnAllowPermissionsIsTheSupersetOfDenyResources(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "arn:aws:iam::222244446666:root"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "arn:*",
          "Principal": {
            "AWS": "arn:aws:iam::222244446666:root"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Principal": {
            "AWS": "arn:aws:iam::666644442222:root"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "arn:*",
          "Principal": {
            "AWS": "arn:aws:iam::666644442222:root"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:iam::222244446666:root",
			"arn:aws:iam::666644442222:root",
		},
		AllowedPrincipalAccountIds: []string{
			"222244446666",
			"666644442222",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds: []string{
			"Statement[1]",
			"Statement[3]",
		},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByArnAllowPermissionsIsTheSubsetOfDenyResources(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "arn:*",
          "Principal": {
            "AWS": "arn:aws:iam::222244446666:root"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "AWS": "arn:aws:iam::222244446666:root"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "arn:*",
          "Principal": {
            "AWS": "arn:aws:iam::666644442222:root"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Principal": {
            "AWS": "arn:aws:iam::666644442222:root"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func TestDenyPermissionsByFederated(t *testing.T) {
	t.Run("TestDenyPermissionsByFederatedRemovesPrincipalWithRespectiveDeny", testDenyPermissionsByFederatedRemovesPrincipalWithRespectiveDeny)
	t.Run("TestDenyPermissionsByFederatedRemovesCorrectPrincipalWithRespectiveDeny", testDenyPermissionsByFederatedRemovesCorrectPrincipalWithRespectiveDeny)
	t.Run("TestDenyPermissionsByFederatedRemovesCorrectPrincipalsWithRespectiveDenies", testDenyPermissionsByFederatedRemovesCorrectPrincipalsWithRespectiveDenies)
	t.Run("TestDenyPermissionsByFederatedRemovesCorrectPrincipalsWhenDenyingMultiplePermissions", testDenyPermissionsByFederatedRemovesCorrectPrincipalsWhenDenyingMultiplePermissions)

	t.Run("TestDenyPermissionsByFederatedRemovesCorrectPrincipalsWhenDenyWildcardPermissions", testDenyPermissionsByFederatedRemovesCorrectPrincipalsWhenDenyWildcardPermissions)
	t.Run("TestDenyPermissionsByFederatedRemovesAllPrincipalsWhenDenyHasMultiplPrincipals", testDenyPermissionsByFederatedRemovesAllPrincipalsWhenDenyHasMultiplPrincipals)
	t.Run("TestDenyPermissionsByFederatedMultiplePermissionsWithMultiplePrincipalsAndDenyOnePermissionsFromEach", testDenyPermissionsByFederatedMultiplePermissionsWithMultiplePrincipalsAndDenyOnePermissionsFromEach)
	t.Run("TestDenyPermissionsByFederatedFullWildcardPrincipalThatFullyContainsAllAllowPermissionsDeniesAll", testDenyPermissionsByFederatedFullWildcardPrincipalThatFullyContainsAllAllowPermissionsDeniesAll)
	t.Run("TestDenyPermissionsByFederatedWhereDenyHasPartiallyWildcardedPrincipalsForAccounts", testDenyPermissionsByFederatedWhereDenyHasPartiallyWildcardedPrincipalsForAccounts)

	t.Run("TestDenyPermissionsByFederatedAllowPermissionsIsTheSupersetOfDenyResources", testDenyPermissionsByFederatedAllowPermissionsIsTheSupersetOfDenyResources)
	t.Run("TestDenyPermissionsByFederatedAllowPermissionsIsTheSubsetOfDenyResources", testDenyPermissionsByFederatedAllowPermissionsIsTheSubsetOfDenyResources)
}

func testDenyPermissionsByFederatedRemovesPrincipalWithRespectiveDeny(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Federated": "cognito-identity.amazonaws.com"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Federated": "cognito-identity.amazonaws.com"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByFederatedRemovesCorrectPrincipalWithRespectiveDeny(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Federated": "cognito-identity.amazonaws.com"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Federated": "accounts.google.com"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Federated": "accounts.google.com"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{"cognito-identity.amazonaws.com"},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByFederatedRemovesCorrectPrincipalsWithRespectiveDenies(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Federated": "cognito-identity.amazonaws.com"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Federated": "accounts.google.com"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Federated": "accounts.google.com"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Federated": "cognito-identity.amazonaws.com"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByFederatedRemovesCorrectPrincipalsWhenDenyingMultiplePermissions(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Federated": ["cognito-identity.amazonaws.com", "accounts.google.com"]
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumesModifications",
          "Resource": "*",
          "Principal": {
            "Federated": ["cognito-identity.amazonaws.com", "accounts.google.com"]
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Federated": "cognito-identity.amazonaws.com"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumesModifications",
          "Resource": "*",
          "Principal": {
            "Federated": "cognito-identity.amazonaws.com"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{"accounts.google.com"},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels: []string{
			"List",
			"Read",
		},
		PrivateAccessLevels: []string{},
		PublicStatementIds:  []string{},
		SharedStatementIds: []string{
			"Statement[1]",
			"Statement[2]",
		},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByFederatedRemovesCorrectPrincipalsWhenDenyWildcardPermissions(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Federated": ["cognito-identity.amazonaws.com", "accounts.google.com"]
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumesModifications",
          "Resource": "*",
          "Principal": {
            "Federated": ["cognito-identity.amazonaws.com", "accounts.google.com"]
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes*",
          "Resource": "*",
          "Principal": {
            "Federated": "cognito-identity.amazonaws.com"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{"accounts.google.com"},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels: []string{
			"List",
			"Read",
		},
		PrivateAccessLevels: []string{},
		PublicStatementIds:  []string{},
		SharedStatementIds: []string{
			"Statement[1]",
			"Statement[2]",
		},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByFederatedRemovesAllPrincipalsWhenDenyHasMultiplPrincipals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Federated": ["cognito-identity.amazonaws.com", "accounts.google.com"]
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Federated": ["cognito-identity.amazonaws.com", "accounts.google.com"]
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByFederatedMultiplePermissionsWithMultiplePrincipalsAndDenyOnePermissionsFromEach(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Federated": ["cognito-identity.amazonaws.com", "accounts.google.com"]
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumesModifications",
          "Resource": "*",
          "Principal": {
            "Federated": ["cognito-identity.amazonaws.com", "accounts.google.com"]
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Federated": "cognito-identity.amazonaws.com"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumesModifications",
          "Resource": "*",
          "Principal": {
            "Federated": "accounts.google.com"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                "shared",
		AllowedOrganizationIds:     []string{},
		AllowedPrincipals:          []string{},
		AllowedPrincipalAccountIds: []string{},
		AllowedPrincipalFederatedIdentities: []string{
			"accounts.google.com",
			"cognito-identity.amazonaws.com",
		},
		AllowedPrincipalServices: []string{},
		IsPublic:                 false,
		PublicAccessLevels:       []string{},
		SharedAccessLevels: []string{
			"List",
			"Read",
		},
		PrivateAccessLevels: []string{},
		PublicStatementIds:  []string{},
		SharedStatementIds: []string{
			"Statement[1]",
			"Statement[2]",
		},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByFederatedFullWildcardPrincipalThatFullyContainsAllAllowPermissionsDeniesAll(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Federated": "cognito-identity.amazonaws.com"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByFederatedWhereDenyHasPartiallyWildcardedPrincipalsForAccounts(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Federated": "cognito-identity.amazonaws.com"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Federated": "graph.facebook.com"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Federated": "graph.facebook.c??"
          }
        },
		{
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Federated": "cognito-*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                "shared",
		AllowedOrganizationIds:     []string{},
		AllowedPrincipals:          []string{},
		AllowedPrincipalAccountIds: []string{},
		AllowedPrincipalFederatedIdentities: []string{
			"cognito-identity.amazonaws.com",
			"graph.facebook.com",
		},
		AllowedPrincipalServices: []string{},
		IsPublic:                 false,
		PublicAccessLevels:       []string{},
		SharedAccessLevels:       []string{"List"},
		PrivateAccessLevels:      []string{},
		PublicStatementIds:       []string{},
		SharedStatementIds: []string{
			"Statement[1]",
			"Statement[2]",
		},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByFederatedAllowPermissionsIsTheSupersetOfDenyResources(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Federated": "cognito-identity.amazonaws.com"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "arn:*",
          "Principal": {
            "Federated": "cognito-identity.amazonaws.com"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Principal": {
            "Federated": "graph.facebook.com"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "arn:*",
          "Principal": {
            "Federated": "graph.facebook.com"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                "shared",
		AllowedOrganizationIds:     []string{},
		AllowedPrincipals:          []string{},
		AllowedPrincipalAccountIds: []string{},
		AllowedPrincipalFederatedIdentities: []string{
			"cognito-identity.amazonaws.com",
			"graph.facebook.com",
		},
		AllowedPrincipalServices: []string{},
		IsPublic:                 false,
		PublicAccessLevels:       []string{},
		SharedAccessLevels:       []string{"List"},
		PrivateAccessLevels:      []string{},
		PublicStatementIds:       []string{},
		SharedStatementIds: []string{
			"Statement[1]",
			"Statement[3]",
		},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByFederatedAllowPermissionsIsTheSubsetOfDenyResources(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "arn:*",
          "Principal": {
            "Federated": "cognito-identity.amazonaws.com"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Federated": "cognito-identity.amazonaws.com"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "arn:*",
          "Principal": {
            "Federated": "graph.facebook.com"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Principal": {
            "Federated": "graph.facebook.com"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func TestDenyPermissionsByService(t *testing.T) {
	t.Run("TestDenyPermissionsByServiceRemovesPrincipalWithRespectiveDeny", testDenyPermissionsByServiceRemovesPrincipalWithRespectiveDeny)
	t.Run("TestDenyPermissionsByServiceRemovesCorrectPrincipalWithRespectiveDeny", testDenyPermissionsByServiceRemovesCorrectPrincipalWithRespectiveDeny)
	t.Run("TestDenyPermissionsByServiceRemovesCorrectPrincipalsWithRespectiveDenies", testDenyPermissionsByServiceRemovesCorrectPrincipalsWithRespectiveDenies)
	t.Run("TestDenyPermissionsByServiceRemovesCorrectPrincipalsWhenDenyingMultiplePermissions", testDenyPermissionsByServiceRemovesCorrectPrincipalsWhenDenyingMultiplePermissions)

	t.Run("TestDenyPermissionsByServiceRemovesCorrectPrincipalsWhenDenyWildcardPermissions", testDenyPermissionsByServiceRemovesCorrectPrincipalsWhenDenyWildcardPermissions)
	t.Run("TestDenyPermissionsByServiceRemovesAllPrincipalsWhenDenyHasMultiplPrincipals", testDenyPermissionsByServiceRemovesAllPrincipalsWhenDenyHasMultiplPrincipals)
	t.Run("TestDenyPermissionsByServiceMultiplePermissionsWithMultiplePrincipalsAndDenyOnePermissionsFromEach", testDenyPermissionsByServiceMultiplePermissionsWithMultiplePrincipalsAndDenyOnePermissionsFromEach)
	t.Run("TestDenyPermissionsByServiceFullWildcardPrincipalThatFullyContainsAllAllowPermissionsDeniesAll", testDenyPermissionsByServiceFullWildcardPrincipalThatFullyContainsAllAllowPermissionsDeniesAll)
	t.Run("TestDenyPermissionsByServiceWhereDenyHasPartiallyWildcardedPrincipalsForAccounts", testDenyPermissionsByServiceWhereDenyHasPartiallyWildcardedPrincipalsForAccounts)

	t.Run("TestDenyPermissionsByServiceAllowPermissionsIsTheSupersetOfDenyResources", testDenyPermissionsByServiceAllowPermissionsIsTheSupersetOfDenyResources)
	t.Run("TestDenyPermissionsByServiceAllowPermissionsIsTheSubsetOfDenyResources", testDenyPermissionsByServiceAllowPermissionsIsTheSubsetOfDenyResources)
}

func testDenyPermissionsByServiceRemovesPrincipalWithRespectiveDeny(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Service": "ec2.amazonaws.com"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Service": "ec2.amazonaws.com"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByServiceRemovesCorrectPrincipalWithRespectiveDeny(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Service": "ec2.amazonaws.com"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Service": "elasticloadbalancing.amazonaws.com"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Service": "elasticloadbalancing.amazonaws.com"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ec2.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"List"},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByServiceRemovesCorrectPrincipalsWithRespectiveDenies(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Service": "ec2.amazonaws.com"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Service": "elasticloadbalancing.amazonaws.com"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Service": "elasticloadbalancing.amazonaws.com"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Service": "ec2.amazonaws.com"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByServiceRemovesCorrectPrincipalsWhenDenyingMultiplePermissions(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Service": ["ec2.amazonaws.com", "elasticloadbalancing.amazonaws.com"]
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumesModifications",
          "Resource": "*",
          "Principal": {
            "Service": ["ec2.amazonaws.com", "elasticloadbalancing.amazonaws.com"]
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Service": "ec2.amazonaws.com"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumesModifications",
          "Resource": "*",
          "Principal": {
            "Service": "ec2.amazonaws.com"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"elasticloadbalancing.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels: []string{
			"List",
			"Read",
		},
		SharedAccessLevels:  []string{},
		PrivateAccessLevels: []string{},
		PublicStatementIds: []string{
			"Statement[1]",
			"Statement[2]",
		},
		SharedStatementIds: []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByServiceRemovesCorrectPrincipalsWhenDenyWildcardPermissions(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Service": ["ec2.amazonaws.com", "elasticloadbalancing.amazonaws.com"]
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumesModifications",
          "Resource": "*",
          "Principal": {
            "Service": ["ec2.amazonaws.com", "elasticloadbalancing.amazonaws.com"]
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes*",
          "Resource": "*",
          "Principal": {
            "Service": "ec2.amazonaws.com"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"elasticloadbalancing.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels: []string{
			"List",
			"Read",
		},
		SharedAccessLevels:  []string{},
		PrivateAccessLevels: []string{},
		PublicStatementIds: []string{
			"Statement[1]",
			"Statement[2]",
		},
		SharedStatementIds: []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByServiceRemovesAllPrincipalsWhenDenyHasMultiplPrincipals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Service": ["ec2.amazonaws.com", "elasticloadbalancing.amazonaws.com"]
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Service": ["ec2.amazonaws.com", "elasticloadbalancing.amazonaws.com"]
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByServiceMultiplePermissionsWithMultiplePrincipalsAndDenyOnePermissionsFromEach(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Service": ["ec2.amazonaws.com", "elasticloadbalancing.amazonaws.com"]
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumesModifications",
          "Resource": "*",
          "Principal": {
            "Service": ["ec2.amazonaws.com", "elasticloadbalancing.amazonaws.com"]
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Service": "ec2.amazonaws.com"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumesModifications",
          "Resource": "*",
          "Principal": {
            "Service": "elasticloadbalancing.amazonaws.com"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices: []string{
			"ec2.amazonaws.com",
			"elasticloadbalancing.amazonaws.com",
		},
		IsPublic: true,
		PublicAccessLevels: []string{
			"List",
			"Read",
		},
		SharedAccessLevels:  []string{},
		PrivateAccessLevels: []string{},
		PublicStatementIds: []string{
			"Statement[1]",
			"Statement[2]",
		},
		SharedStatementIds: []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByServiceFullWildcardPrincipalThatFullyContainsAllAllowPermissionsDeniesAll(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Service": "ec2.amazonaws.com"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": "*"
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByServiceWhereDenyHasPartiallyWildcardedPrincipalsForAccounts(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Service": "ec2.amazonaws.com"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Service": "graph.facebook.com"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Service": "graph.facebook.c??"
          }
        },
		{
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Service": "cognito-*"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices: []string{
			"ec2.amazonaws.com",
			"graph.facebook.com",
		},
		IsPublic:            true,
		PublicAccessLevels:  []string{"List"},
		SharedAccessLevels:  []string{},
		PrivateAccessLevels: []string{},
		PublicStatementIds: []string{
			"Statement[1]",
			"Statement[2]",
		},
		SharedStatementIds: []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByServiceAllowPermissionsIsTheSupersetOfDenyResources(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Service": "ec2.amazonaws.com"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "arn:*",
          "Principal": {
            "Service": "ec2.amazonaws.com"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Principal": {
            "Service": "ecr.amazonaws.com"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "arn:*",
          "Principal": {
            "Service": "ecr.amazonaws.com"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices: []string{
			"ec2.amazonaws.com",
			"ecr.amazonaws.com",
		},
		IsPublic:            true,
		PublicAccessLevels:  []string{"List"},
		SharedAccessLevels:  []string{},
		PrivateAccessLevels: []string{},
		PublicStatementIds: []string{
			"Statement[1]",
			"Statement[3]",
		},
		SharedStatementIds: []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByServiceAllowPermissionsIsTheSubsetOfDenyResources(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "arn:*",
          "Principal": {
            "Service": "ec2.amazonaws.com"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Principal": {
            "Service": "ec2.amazonaws.com"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "arn:*",
          "Principal": {
            "Service": "ecr.amazonaws.com"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Principal": {
            "Service": "ecr.amazonaws.com"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func TestDenyPermissionsByGlobalConditionPrincipalAccount(t *testing.T) {
	t.Run("TestDenyPermissionsByGlobalConditionPrincipalAccountRemovesPrincipalWithRespectiveDeny", testDenyPermissionsByGlobalConditionPrincipalAccountRemovesPrincipalWithRespectiveDeny)
	t.Run("TestDenyPermissionsByGlobalConditionPrincipalAccountRemovesCorrectPrincipalWithRespectiveDeny", testDenyPermissionsByGlobalConditionPrincipalAccountRemovesCorrectPrincipalWithRespectiveDeny)
	t.Run("TestDenyPermissionsByGlobalConditionPrincipalAccountRemovesCorrectPrincipalsWithRespectiveDenies", testDenyPermissionsByGlobalConditionPrincipalAccountRemovesCorrectPrincipalsWithRespectiveDenies)
	t.Run("TestDenyPermissionsByGlobalConditionPrincipalAccountRemovesCorrectPrincipalsWhenDenyingMultiplePermissions", testDenyPermissionsByGlobalConditionPrincipalAccountRemovesCorrectPrincipalsWhenDenyingMultiplePermissions)

	t.Run("TestDenyPermissionsByGlobalConditionPrincipalAccountRemovesCorrectPrincipalsWhenDenyWildcardPermissions", testDenyPermissionsByGlobalConditionPrincipalAccountRemovesCorrectPrincipalsWhenDenyWildcardPermissions)
	t.Run("TestDenyPermissionsByGlobalConditionPrincipalAccountRemovesAllPrincipalsWhenDenyHasMultiplPrincipals", testDenyPermissionsByGlobalConditionPrincipalAccountRemovesAllPrincipalsWhenDenyHasMultiplPrincipals)
	t.Run("TestDenyPermissionsByGlobalConditionPrincipalAccountMultiplePermissionsWithMultiplePrincipalsAndDenyOnePermissionsFromEach", testDenyPermissionsByGlobalConditionPrincipalAccountMultiplePermissionsWithMultiplePrincipalsAndDenyOnePermissionsFromEach)
	t.Run("TestDenyPermissionsByGlobalConditionPrincipalAccountFullWildcardPrincipalThatFullyContainsAllAllowPermissionsDeniesAll", testDenyPermissionsByGlobalConditionPrincipalAccountFullWildcardPrincipalThatFullyContainsAllAllowPermissionsDeniesAll)
	t.Run("TestDenyPermissionsByGlobalConditionPrincipalAccountWhereDenyHasPartiallyWildcardedPrincipalsForAccounts", testDenyPermissionsByGlobalConditionPrincipalAccountWhereDenyHasPartiallyWildcardedPrincipalsForAccounts)
	t.Run("TestDenyPermissionsByGlobalConditionPrincipalAccountWhereDenyHasPartiallyWildcardedPrincipalsForOtherAccounts", testDenyPermissionsByGlobalConditionPrincipalAccountWhereDenyHasPartiallyWildcardedPrincipalsForOtherAccounts)

	t.Run("TestDenyPermissionsByGlobalConditionPrincipalAccountAllowPermissionsIsTheSupersetOfDenyPrincipals", testDenyPermissionsByGlobalConditionPrincipalAccountAllowPermissionsIsTheSupersetOfDenyPrincipals)
	t.Run("TestDenyPermissionsByGlobalConditionPrincipalAccountAllowPermissionsIsTheSubsetOfDenyPrincipals", testDenyPermissionsByGlobalConditionPrincipalAccountAllowPermissionsIsTheSubsetOfDenyPrincipals)
	t.Run("TestDenyPermissionsByGlobalConditionPrincipalAccountAllowPermissionsIsTheSupersetOfDenyResources", testDenyPermissionsByGlobalConditionPrincipalAccountAllowPermissionsIsTheSupersetOfDenyResources)
	t.Run("TestDenyPermissionsByGlobalConditionPrincipalAccountAllowPermissionsIsTheSubsetOfDenyResources", testDenyPermissionsByGlobalConditionPrincipalAccountAllowPermissionsIsTheSubsetOfDenyResources)
	t.Run("TestDenyPermissionsByGlobalConditionPrincipalAccountWithTheSameWildcard", testDenyPermissionsByGlobalConditionPrincipalAccountWithTheSameWildcard)
}

func testDenyPermissionsByGlobalConditionPrincipalAccountRemovesPrincipalWithRespectiveDeny(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["222244446666"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["222244446666"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByGlobalConditionPrincipalAccountRemovesCorrectPrincipalWithRespectiveDeny(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["222244446666"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["666644442222"]
            }
          },
          "Principal": {
            "AWS": "666644442222"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["666644442222"]
            }
          },
          "Principal": {
            "AWS": "666644442222"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"222244446666"},
		AllowedPrincipalAccountIds:          []string{"222244446666"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByGlobalConditionPrincipalAccountRemovesCorrectPrincipalsWithRespectiveDenies(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["222244446666"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["666644442222"]
            }
          },
          "Principal": {
            "AWS": "666644442222"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["666644442222"]
            }
          },
          "Principal": {
            "AWS": "666644442222"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["222244446666"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByGlobalConditionPrincipalAccountRemovesCorrectPrincipalsWhenDenyingMultiplePermissions(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["222244446666", "666644442222"]
            }
          },
          "Principal": {
            "AWS": ["222244446666", "666644442222"]
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumesModifications",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["222244446666", "666644442222"]
            }
          },
          "Principal": {
            "AWS": ["222244446666", "666644442222"]
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["222244446666"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumesModifications",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["222244446666"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"666644442222"},
		AllowedPrincipalAccountIds:          []string{"666644442222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels: []string{
			"List",
			"Read",
		},
		PrivateAccessLevels: []string{},
		PublicStatementIds:  []string{},
		SharedStatementIds: []string{
			"Statement[1]",
			"Statement[2]",
		},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByGlobalConditionPrincipalAccountRemovesCorrectPrincipalsWhenDenyWildcardPermissions(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["222244446666", "666644442222"]
            }
          },
          "Principal": {
            "AWS": ["222244446666", "666644442222"]
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumesModifications",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["222244446666", "666644442222"]
            }
          },
          "Principal": {
            "AWS": ["222244446666", "666644442222"]
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes*",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["222244446666"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"666644442222"},
		AllowedPrincipalAccountIds:          []string{"666644442222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels: []string{
			"List",
			"Read",
		},
		PrivateAccessLevels: []string{},
		PublicStatementIds:  []string{},
		SharedStatementIds: []string{
			"Statement[1]",
			"Statement[2]",
		},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByGlobalConditionPrincipalAccountRemovesAllPrincipalsWhenDenyHasMultiplPrincipals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["222244446666", "666644442222"]
            }
          },
          "Principal": {
            "AWS": ["222244446666", "666644442222"]
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["222244446666", "666644442222"]
            }
          },
          "Principal": {
            "AWS": ["222244446666", "666644442222"]
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByGlobalConditionPrincipalAccountMultiplePermissionsWithMultiplePrincipalsAndDenyOnePermissionsFromEach(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["222244446666", "666644442222"]
            }
          },
          "Principal": {
            "AWS": ["222244446666", "666644442222"]
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumesModifications",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["222244446666", "666644442222"]
            }
          },
          "Principal": {
            "AWS": ["222244446666", "666644442222"]
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["222244446666"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumesModifications",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["666644442222"]
            }
          },
          "Principal": {
            "AWS": "666644442222"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"222244446666",
			"666644442222",
		},
		AllowedPrincipalAccountIds: []string{
			"222244446666",
			"666644442222",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels: []string{
			"List",
			"Read",
		},
		PrivateAccessLevels: []string{},
		PublicStatementIds:  []string{},
		SharedStatementIds: []string{
			"Statement[1]",
			"Statement[2]",
		},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByGlobalConditionPrincipalAccountFullWildcardPrincipalThatFullyContainsAllAllowPermissionsDeniesAll(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["222244446666"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalAccount": ["*"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByGlobalConditionPrincipalAccountWhereDenyHasPartiallyWildcardedPrincipalsForAccounts(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["222244446666"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["012345678901"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalAccount": ["0123456789??"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalAccount": ["2222*"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByGlobalConditionPrincipalAccountWhereDenyHasPartiallyWildcardedPrincipalsForOtherAccounts(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["222244446666"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["012345678901"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalAccount": ["98765432310??"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalAccount": ["4444*"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"222244446666",
		},
		AllowedPrincipalAccountIds: []string{
			"222244446666",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByGlobalConditionPrincipalAccountAllowPermissionsIsTheSupersetOfDenyPrincipals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalAccount": ["22224444????"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalAccount": ["222244446666"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByGlobalConditionPrincipalAccountAllowPermissionsIsTheSubsetOfDenyPrincipals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["222244446666"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalAccount": ["22224444????"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByGlobalConditionPrincipalAccountAllowPermissionsIsTheSupersetOfDenyResources(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalAccount": ["222244446666"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "arn:*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalAccount": ["222244446666"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"222244446666"},
		AllowedPrincipalAccountIds:          []string{"222244446666"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByGlobalConditionPrincipalAccountAllowPermissionsIsTheSubsetOfDenyResources(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "arn:*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["222244446666"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalAccount": ["222244446666"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByGlobalConditionPrincipalAccountWithTheSameWildcard(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalAccount": ["22224444????"]
            }
          },
          "Principal": {
            "AWS": "222244442222"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalAccount": ["22224444????"]
            }
          },
          "Principal": {
            "AWS": "222244442222"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func TestDenyPermissionsByGlobalConditionPrincipalArn(t *testing.T) {
	t.Run("TestDenyPermissionsByGlobalConditionPrincipalArnRemovesPrincipalWithRespectiveDeny", testDenyPermissionsByGlobalConditionPrincipalArnRemovesPrincipalWithRespectiveDeny)
	t.Run("TestDenyPermissionsByGlobalConditionPrincipalArnRemovesCorrectPrincipalWithRespectiveDeny", testDenyPermissionsByGlobalConditionPrincipalArnRemovesCorrectPrincipalWithRespectiveDeny)
	t.Run("TestDenyPermissionsByGlobalConditionPrincipalArnRemovesCorrectPrincipalsWithRespectiveDenies", testDenyPermissionsByGlobalConditionPrincipalArnRemovesCorrectPrincipalsWithRespectiveDenies)
	t.Run("TestDenyPermissionsByGlobalConditionPrincipalArnRemovesCorrectPrincipalsWhenDenyingMultiplePermissions", testDenyPermissionsByGlobalConditionPrincipalArnRemovesCorrectPrincipalsWhenDenyingMultiplePermissions)

	t.Run("TestDenyPermissionsByGlobalConditionPrincipalArnRemovesCorrectPrincipalsWhenDenyWildcardPermissions", testDenyPermissionsByGlobalConditionPrincipalArnRemovesCorrectPrincipalsWhenDenyWildcardPermissions)
	t.Run("TestDenyPermissionsByGlobalConditionPrincipalArnRemovesAllPrincipalsWhenDenyHasMultiplPrincipals", testDenyPermissionsByGlobalConditionPrincipalArnRemovesAllPrincipalsWhenDenyHasMultiplPrincipals)
	t.Run("TestDenyPermissionsByGlobalConditionPrincipalArnMultiplePermissionsWithMultiplePrincipalsAndDenyOnePermissionsFromEach", testDenyPermissionsByGlobalConditionPrincipalArnMultiplePermissionsWithMultiplePrincipalsAndDenyOnePermissionsFromEach)
	t.Run("TestDenyPermissionsByGlobalConditionPrincipalArnFullWildcardPrincipalThatFullyContainsAllAllowPermissionsDeniesAll", testDenyPermissionsByGlobalConditionPrincipalArnFullWildcardPrincipalThatFullyContainsAllAllowPermissionsDeniesAll)
	t.Run("TestDenyPermissionsByGlobalConditionPrincipalArnWhereDenyHasPartiallyWildcardedPrincipalsForAccounts", testDenyPermissionsByGlobalConditionPrincipalArnWhereDenyHasPartiallyWildcardedPrincipalsForAccounts)
	t.Run("TestDenyPermissionsByGlobalConditionPrincipalArnWhereDenyHasPartiallyWildcardedPrincipalsForOtherAccounts", testDenyPermissionsByGlobalConditionPrincipalArnWhereDenyHasPartiallyWildcardedPrincipalsForOtherAccounts)

	t.Run("TestDenyPermissionsByGlobalConditionPrincipalArnAllowPermissionsIsTheSupersetOfDenyPrincipals", testDenyPermissionsByGlobalConditionPrincipalArnAllowPermissionsIsTheSupersetOfDenyPrincipals)
	t.Run("TestDenyPermissionsByGlobalConditionPrincipalArnAllowPermissionsIsTheSubsetOfDenyPrincipals", testDenyPermissionsByGlobalConditionPrincipalArnAllowPermissionsIsTheSubsetOfDenyPrincipals)
	t.Run("TestDenyPermissionsByGlobalConditionPrincipalArnAllowPermissionsIsTheSupersetOfDenyResources", testDenyPermissionsByGlobalConditionPrincipalArnAllowPermissionsIsTheSupersetOfDenyResources)
	t.Run("TestDenyPermissionsByGlobalConditionPrincipalArnAllowPermissionsIsTheSubsetOfDenyResources", testDenyPermissionsByGlobalConditionPrincipalArnAllowPermissionsIsTheSubsetOfDenyResources)

	t.Run("TestDenyPermissionsByGlobalConditionPrincipalArnWithTheSameWildcard", testDenyPermissionsByGlobalConditionPrincipalArnWithTheSameWildcard)
}

func testDenyPermissionsByGlobalConditionPrincipalArnRemovesPrincipalWithRespectiveDeny(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::222244446666:root"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::222244446666:root"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByGlobalConditionPrincipalArnRemovesCorrectPrincipalWithRespectiveDeny(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::222244446666:root"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::666644442222:root"]
            }
          },
          "Principal": {
            "AWS": "666644442222"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::666644442222:root"]
            }
          },
          "Principal": {
            "AWS": "666644442222"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::222244446666:root"},
		AllowedPrincipalAccountIds:          []string{"222244446666"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByGlobalConditionPrincipalArnRemovesCorrectPrincipalsWithRespectiveDenies(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::222244446666:root"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::666644442222:root"]
            }
          },
          "Principal": {
            "AWS": "666644442222"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::666644442222:root"]
            }
          },
          "Principal": {
            "AWS": "666644442222"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::222244446666:root"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByGlobalConditionPrincipalArnRemovesCorrectPrincipalsWhenDenyingMultiplePermissions(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::222244446666:root", "arn:aws:iam::666644442222:root"]
            }
          },
          "Principal": {
            "AWS": ["222244446666", "666644442222"]
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumesModifications",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::222244446666:root", "arn:aws:iam::666644442222:root"]
            }
          },
          "Principal": {
            "AWS": ["222244446666", "666644442222"]
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::222244446666:root"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumesModifications",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::222244446666:root"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::666644442222:root"},
		AllowedPrincipalAccountIds:          []string{"666644442222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels: []string{
			"List",
			"Read",
		},
		PrivateAccessLevels: []string{},
		PublicStatementIds:  []string{},
		SharedStatementIds: []string{
			"Statement[1]",
			"Statement[2]",
		},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByGlobalConditionPrincipalArnRemovesCorrectPrincipalsWhenDenyWildcardPermissions(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::222244446666:root", "arn:aws:iam::666644442222:root"]
            }
          },
          "Principal": {
            "AWS": ["222244446666", "666644442222"]
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumesModifications",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::222244446666:root", "arn:aws:iam::666644442222:root"]
            }
          },
          "Principal": {
            "AWS": ["222244446666", "666644442222"]
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes*",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::222244446666:root"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::666644442222:root"},
		AllowedPrincipalAccountIds:          []string{"666644442222"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels: []string{
			"List",
			"Read",
		},
		PrivateAccessLevels: []string{},
		PublicStatementIds:  []string{},
		SharedStatementIds: []string{
			"Statement[1]",
			"Statement[2]",
		},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByGlobalConditionPrincipalArnRemovesAllPrincipalsWhenDenyHasMultiplPrincipals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::222244446666:root", "arn:aws:iam::666644442222:root"]
            }
          },
          "Principal": {
            "AWS": ["222244446666", "666644442222"]
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::222244446666:root", "arn:aws:iam::666644442222:root"]
            }
          },
          "Principal": {
            "AWS": ["222244446666", "666644442222"]
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByGlobalConditionPrincipalArnMultiplePermissionsWithMultiplePrincipalsAndDenyOnePermissionsFromEach(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::222244446666:root", "arn:aws:iam::666644442222:root"]
            }
          },
          "Principal": {
            "AWS": ["222244446666", "666644442222"]
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumesModifications",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::222244446666:root", "arn:aws:iam::666644442222:root"]
            }
          },
          "Principal": {
            "AWS": ["222244446666", "666644442222"]
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::222244446666:root"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumesModifications",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::666644442222:root"]
            }
          },
          "Principal": {
            "AWS": "666644442222"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:iam::222244446666:root",
			"arn:aws:iam::666644442222:root",
		},
		AllowedPrincipalAccountIds: []string{
			"222244446666",
			"666644442222",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels: []string{
			"List",
			"Read",
		},
		PrivateAccessLevels: []string{},
		PublicStatementIds:  []string{},
		SharedStatementIds: []string{
			"Statement[1]",
			"Statement[2]",
		},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByGlobalConditionPrincipalArnFullWildcardPrincipalThatFullyContainsAllAllowPermissionsDeniesAll(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::222244446666:root"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalArn": ["*"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByGlobalConditionPrincipalArnWhereDenyHasPartiallyWildcardedPrincipalsForAccounts(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::222244446666:root"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::012345678901:root"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalArn": ["arn:aws:iam::0123456789??:root"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalArn": ["arn:aws:iam::2222*:root"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByGlobalConditionPrincipalArnWhereDenyHasPartiallyWildcardedPrincipalsForOtherAccounts(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::222244446666:root"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::012345678901:root"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalArn": ["arn:aws:iam::9876543231??:root"]
            }
          },
          "Principal": {
            "AWS": "012345678901"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalArn": ["arn:aws:iam::4444*:root"]
            }
          },
          "Principal": {
            "AWS": "222244442222"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:iam::012345678901:root",
			"arn:aws:iam::222244446666:root",
		},
		AllowedPrincipalAccountIds: []string{
			"012345678901",
			"222244446666",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{"List"},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByGlobalConditionPrincipalArnAllowPermissionsIsTheSupersetOfDenyPrincipals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalArn": ["arn:aws:iam::22224444????:root"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalArn": ["arn:aws:iam::222244446666:root"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByGlobalConditionPrincipalArnAllowPermissionsIsTheSubsetOfDenyPrincipals(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::222244446666:root"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalArn": ["arn:aws:iam::22224444????:root"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByGlobalConditionPrincipalArnAllowPermissionsIsTheSupersetOfDenyResources(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalArn": ["arn:aws:iam::222244446666:root"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "arn:*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalArn": ["arn:aws:iam::222244446666:root"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::222244446666:root"},
		AllowedPrincipalAccountIds:          []string{"222244446666"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{"List"},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByGlobalConditionPrincipalArnAllowPermissionsIsTheSubsetOfDenyResources(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "arn:*",
          "Condition": {
            "StringEquals": {
              "aws:PrincipalArn": ["arn:aws:iam::222244446666:root"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalArn": ["arn:aws:iam::222244446666:root"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}

func testDenyPermissionsByGlobalConditionPrincipalArnWithTheSameWildcard(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalArn": ["arn:aws:iam::22224444????:root/*"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        },
        {
          "Effect": "Deny",
          "Action": "ec2:DescribeVolumes",
          "Resource": "*",
          "Condition": {
            "StringLike": {
              "aws:PrincipalArn": ["arn:aws:iam::22224444????:root/*"]
            }
          },
          "Principal": {
            "AWS": "222244446666"
          }
        }
      ]
    }
	`

	expected := PolicySummary{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		SharedAccessLevels:                  []string{},
		PrivateAccessLevels:                 []string{},
		PublicStatementIds:                  []string{},
		SharedStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	errors := evaluatePrincipalTest(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Fatal("Conditions Unit Test error detected")
	}

	errors = evaluateIntegration(t, evaluated, expected)
	if len(errors) > 0 {
		for _, error := range errors {
			t.Log(error)
		}
		t.Log("Integration Test error detected - Find Unit Test error to resolve issue")
		t.Fail()
	}
}