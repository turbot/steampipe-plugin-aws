package aws

import (
	"fmt"
	"testing"
)

/// Evaluation Functions
func evaluatePublicAccessLevelsTest(t *testing.T, source EvaluatedPolicy, expected EvaluatedPolicy) []string {
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

	return errors
}

func evaluateAccessLevelTest(t *testing.T, source EvaluatedPolicy, expected EvaluatedPolicy) []string {
	errors := []string{}

	currentAccessLevel := source.AccessLevel
	expectedAccessLevel := expected.AccessLevel
	if currentAccessLevel != expectedAccessLevel {
		errors = append(errors, fmt.Sprintf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel))
	}

	return errors
}

func evaluateSidTest(t *testing.T, source EvaluatedPolicy, expected EvaluatedPolicy) []string {
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

	return errors
}

func evaluatePrincipalTest(t *testing.T, source EvaluatedPolicy, expected EvaluatedPolicy) []string {
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

func evaluateIntegration(t *testing.T, source EvaluatedPolicy, expected EvaluatedPolicy) []string {
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

	return errors
}

/// Test start here

func TestPolicyStatementElement(t *testing.T) {
	t.Run("TestPolicyCreatedWithCanonicaliseWithNoStatementsPolicyEvaluates", testPolicyCreatedWithCanonicaliseWithNoStatementsPolicyEvaluates)

	t.Run("TestPolicyCreatedByStructEvaluates", testPolicyCreatedByStringEvaluates)
	t.Run("TestPolicyCreatedByEmptyJsonStringEvaluates", testPolicyCreatedByEmptyJsonStringEvaluates)
}

func testPolicyCreatedWithCanonicaliseWithNoStatementsPolicyEvaluates(t *testing.T) {
	// Set up
	userAccountId := "123456789012"
	policyContent := `
	{
	  "Version": "2012-10-17"
	}
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		PublicStatementIds:                  []string{},
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

func testPolicyCreatedByStringEvaluates(t *testing.T) {
	// Set up
	userAccountId := "123456789012"
	expected := EvaluatedPolicy{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		PublicStatementIds:                  []string{},
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

func testPolicyCreatedByEmptyJsonStringEvaluates(t *testing.T) {
	// Set up
	userAccountId := "123456789012"
	expected := EvaluatedPolicy{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		PublicStatementIds:                  []string{},
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

func TestPolicyEffectElement(t *testing.T) {
	t.Run("TestEffectElementWithValidValues", testEffectElementWithValidValues)
	t.Run("TestIfEffectElementWhenValueAllowHasWrongCasingFails", testIfEffectElementWhenValueAllowHasWrongCasingFails)
	t.Run("TestIfEffectElementWhenValueDenyHasWrongCasingFails", testIfEffectElementWhenValueDenyHasWrongCasingFails)
	t.Run("TestIfEffectElementWhenValueIsUnknownFails", testIfEffectElementWhenValueIsUnknownFails)
}

func testEffectElementWithValidValues(t *testing.T) {
	// Set up
	userAccountId := "123456789012"
	policyContent := `
    {
      "Version": "2012-10-17",
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
	expected := EvaluatedPolicy{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		PublicStatementIds:                  []string{},
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

	expectedErrorMsg := "element Effect is invalid - valid choices are 'Allow' or 'Deny'"

	if errorMsg := err.Error(); errorMsg != expectedErrorMsg {
		t.Fatalf("The error message returned is expected to be: %s", expectedErrorMsg)
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

	expectedErrorMsg := "element Effect is invalid - valid choices are 'Allow' or 'Deny'"

	if errorMsg := err.Error(); errorMsg != expectedErrorMsg {
		t.Fatalf("The error message returned is expected to be: %s", expectedErrorMsg)
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

	expectedErrorMsg := "element Effect is invalid - valid choices are 'Allow' or 'Deny'"

	if errorMsg := err.Error(); errorMsg != expectedErrorMsg {
		t.Fatalf("The error message returned is expected to be: %s", expectedErrorMsg)
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

	// Test
	_, err := EvaluatePolicy("", userAccountId)

	// Evaluate
	if err == nil {
		t.Fatal("Expected error but no error was returned from EvaluatePolicy")
	}

	expectedErrorMsg := "source account id is invalid: 123A123123"

	if errorMsg := err.Error(); errorMsg != expectedErrorMsg {
		t.Fatalf("The error message returned is expected to be: %s", expectedErrorMsg)
	}
}

func testIfSourceAccountIdContainsTooFewNumericalValuesItFails(t *testing.T) {
	// Set up
	userAccountId := "01234567890"

	// Test
	_, err := EvaluatePolicy("", userAccountId)

	// Evaluate
	if err == nil {
		t.Fatal("Expected error but no error was returned from EvaluatePolicy")
	}

	expectedErrorMsg := "source account id is invalid: 01234567890"

	if errorMsg := err.Error(); errorMsg != expectedErrorMsg {
		t.Fatalf("The error message returned is expected to be: %s", expectedErrorMsg)
	}
}

func testIfSourceAccountIdContainsTooManyNumericalValuesItFails(t *testing.T) {
	// Set up
	userAccountId := "012345678901234"

	// Test
	_, err := EvaluatePolicy("", userAccountId)

	// Evaluate
	if err == nil {
		t.Fatal("Expected error but no error was returned from EvaluatePolicy")
	}

	expectedErrorMsg := "source account id is invalid: 012345678901234"

	if errorMsg := err.Error(); errorMsg != expectedErrorMsg {
		t.Fatalf("The error message returned is expected to be: %s", expectedErrorMsg)
	}
}

func testIfSourceAccountIdContainsCorrectAmountOfNumericalValuesItEvaluates(t *testing.T) {
	// Set up
	userAccountId := "123456789012"
	expected := EvaluatedPolicy{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		PublicStatementIds:                  []string{},
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
	expected := EvaluatedPolicy{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		PublicStatementIds:                  []string{},
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

func TestPolicyPrincipalElement(t *testing.T) {
	t.Run("TestWhenPricipalIsAMisformedArnFails", testWhenPricipalIsAMisformedArnFails)
	t.Run("TestWhenPrincipalIsWildcarded", testWhenPrincipalIsWildcarded)
	t.Run("TestWhenAwsPrincipalIsWildcarded", testWhenAwsPrincipalIsWildcarded)
	t.Run("TestWhenAwsPrincipalIsWildcardedFollowedByNormalStatementShouldKeepIsPublic", testWhenAwsPrincipalIsWildcardedFollowedByNormalStatementShouldKeepItPublic)

	t.Run("TestWhenPrincipalIsAUserAccountId", testWhenPrincipalIsAUserAccountId)
	t.Run("TestWhenPrincipalIsAUserAccountArn", testWhenPrincipalIsAUserAccountArn)
	t.Run("TestWhenPrincipalIsACrossAccountId", testWhenPrincipalIsACrossAccountId)
	t.Run("TestWhenPrincipalIsACrossAccountArn", testWhenPrincipalIsACrossAccountArn)
	t.Run("TestWhenPrincipalIsMultipleUserAccounts", testWhenPrincipalIsMultipleUserAccounts)
	t.Run("TestWhenPrincipalIsMultipleCrossAccountsInAscendingOrder", testWhenPrincipalIsMultipleCrossAccountsInAscendingOrder)
	t.Run("TestWhenPrincipalIsMultipleCrossAccountsInDescendingOrder", testWhenPrincipalIsMultipleCrossAccountsInDescendingOrder)
	t.Run("TestWhenPrincipalIsMultipleAccountsPrincipalsAcrossMultipleStatements", testWhenPrincipalIsMultipleAccountsPrincipalsAcrossMultipleStatements)
	t.Run("TestWhenPrincipalIsMultipleMixedAccounts", testWhenPrincipalIsMultipleMixedAccounts)

	t.Run("TestWhenPrincipalIsMultipleMixedAccountsWithWildcard", testWhenPrincipalIsMultipleMixedAccountsWithWildcard)

	t.Run("TestWhenPricipalIsAUserAccountRole", testWhenPricipalIsAUserAccountRole)
	t.Run("TestWhenPricipalIsACrossAccountRole", testWhenPricipalIsACrossAccountRole)
	t.Run("TestWhenPrincipalIsMultipleUserAccountRoles", testWhenPrincipalIsMultipleUserAccountRoles)
	t.Run("TestWhenPrincipalIsMultipleCrossAccountRolesInAscendingOrder", testWhenPrincipalIsMultipleCrossAccountRolesInAscendingOrder)
	t.Run("TestWhenPrincipalIsMultipleCrossAccountRolesInDescendingOrder", testWhenPrincipalIsMultipleCrossAccountRolesInDescendingOrder)
	t.Run("TestWhenPrincipalIsMultipleAccountRolePrincipalsAcrossMultipleStatements", testWhenPrincipalIsMultipleAccountRolePrincipalsAcrossMultipleStatements)
	t.Run("TestWhenPrincipalIsMultipleMixedAccountRoles", testWhenPrincipalIsMultipleMixedAccountRoles)

	t.Run("TestWhenPricipalIsAUserAccountAssumedRole", testWhenPricipalIsAUserAccountAssumedRole)
	t.Run("TestWhenPricipalIsACrossAccountAssumedRole", testWhenPricipalIsACrossAccountAssumedRole)
	t.Run("TestWhenPrincipalIsMultipleUserAccountAssumedRoles", testWhenPrincipalIsMultipleUserAccountAssumedRoles)
	t.Run("TestWhenPrincipalIsMultipleCrossAccountAssumedRolesInAscendingOrder", testWhenPrincipalIsMultipleCrossAccountAssumedRolesInAscendingOrder)
	t.Run("TestWhenPrincipalIsMultipleCrossAccountAssumedRolesInDescendingOrder", testWhenPrincipalIsMultipleCrossAccountAssumedRolesInDescendingOrder)
	t.Run("TestWhenPrincipalIsMultipleAccountAssumedRolePrincipalsAcrossMultipleStatements", testWhenPrincipalIsMultipleAccountAssumedRolePrincipalsAcrossMultipleStatements)
	t.Run("TestWhenPrincipalIsMultipleMixedAccountAssumedRoles", testWhenPrincipalIsMultipleMixedAccountAssumedRoles)

	t.Run("TestWhenPricipalIsAFederatedUser", testWhenPricipalIsAFederatedUser)
	t.Run("TestWhenPricipalIsMulitpleFederatedUserInAscendingOrder", testWhenPrincipalIsMulitpleFederatedUserInAscendingOrder)
	t.Run("TestWhenPrincipalIsMulitpleFederatedUserInDescendingOrder", testWhenPrincipalIsMulitpleFederatedUserInDescendingOrder)
	t.Run("TestWhenPrincipalIsMultipleFederatedUserPrincipalsAcrossMultipleStatements", testWhenPrincipalIsMultipleFederatedUserPrincipalsAcrossMultipleStatements)

	t.Run("TestWhenPricipalIsAService", testWhenPricipalIsAService)
	t.Run("TestWhenPrincipalIsMulitpleServicesInAscendingOrder", testWhenPrincipalIsMulitpleServicesInAscendingOrder)
	t.Run("TestWhenPrincipalIsMulitpleServicesInDescendingOrder", testWhenPrincipalIsMulitpleServicesInDescendingOrder)
	t.Run("TestWhenPrincipalIsMultipleServicePrincipalsAcrossMultipleStatements", testWhenPrincipalIsMultipleServicePrincipalsAcrossMultipleStatements)

	t.Run("TestWhenPrincipalIsMultipleTypes", testWhenPrincipalIsMultipleTypes)
	t.Run("TestWhenPrincipalIsMultipleTypesWithWildcard", testWhenPrincipalIsMultipleTypesWithWildcard)
	t.Run("TestWhenPrincipalIsMultipleTypesAcrossMultipleStatements", testWhenPrincipalIsMultipleTypesAcrossMultipleStatements)
	t.Run("TestWhenPrincipalIsMultipleTypesAcrossMultipleStatementsWithWildcard", testWhenPrincipalIsMultipleTypesAcrossMultipleStatementsWithWildcard)

	t.Run("TestWhenAwsPrincipalIsWildcardedAndEffectDenied", testWhenAwsPrincipalIsWildcardedAndEffectDenied)
	t.Run("TestWhenAwsPrincipalIsWildcardedDeniedButAnotherAccountIsAllowed", TestWhenAwsPrincipalIsWildcardedDeniedButAnotherAccountIsAllowed)
}

func testWhenPricipalIsAMisformedArnFails(t *testing.T) {
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
			"AWS": "arn:aws:sts::misformed:012345678901:assumed-role/role-name/role-session-name"
          }
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

	expectedErrorMsg := "unabled to parse arn: arn:aws:sts::misformed:012345678901:assumed-role/role-name/role-session-name"

	if errorMsg := err.Error(); errorMsg != expectedErrorMsg {
		t.Fatalf("The error message returned is expected to be: %s", expectedErrorMsg)
	}
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
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{"Write"},
		PublicStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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
          "Principal": "*"
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"Write"},
		PublicStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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

func testWhenAwsPrincipalIsWildcardedFollowedByNormalStatementShouldKeepItPublic(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": "*"
        },
        {
			"Effect": "Allow",
			"Action": "sts:AssumeRole",
			"Principal": "012345678901"
		  } 
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:            "public",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"*",
			"012345678901",
		},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"Write"},
		PublicStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"Write"},
		PublicStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::012345678901:root"},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{"Write"},
		PublicStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::444455554444:root"},
		AllowedPrincipalAccountIds:          []string{"444455554444"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{"Write"},
		PublicStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"444455554444"},
		AllowedPrincipalAccountIds:          []string{"444455554444"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{"Write"},
		PublicStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:            "private",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"012345678901",
			"arn:aws:iam::012345678901:root",
		},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{"Write"},
		PublicStatementIds:                  []string{},
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
          }
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:iam::012345678901:root"]
          }
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:iam::444455554444:root"]
          }
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:iam::012345678901:root"]
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:iam::012345678901:root",
			"arn:aws:iam::444455554444:root",
		},
		AllowedPrincipalAccountIds:          []string{"444455554444"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{"Write"},
		PublicStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
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
		PublicAccessLevels:                  []string{"Write"},
		PublicStatementIds:                  []string{},
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
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
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
		PublicAccessLevels:                  []string{"Write"},
		PublicStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"012345678901",
			"444455554444",
			"arn:aws:iam::012345678901:root",
			"arn:aws:iam::444455554444:root",
		},
		AllowedPrincipalAccountIds:          []string{"444455554444"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{"Write"},
		PublicStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
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
			"444455554444",
		},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"Write"},
		PublicStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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

func testWhenPricipalIsAUserAccountRole(t *testing.T) {
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
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::012345678901:role/role-name"},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{"Write"},
		PublicStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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

func testWhenPricipalIsACrossAccountRole(t *testing.T) {
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
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::444455554444:role/role-name"},
		AllowedPrincipalAccountIds:          []string{"444455554444"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{"Write"},
		PublicStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
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
		PublicAccessLevels:                  []string{"Write"},
		PublicStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
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
		PublicAccessLevels:                  []string{"Write"},
		PublicStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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
          }
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:iam::012345678901:role/role-name"]
          }
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:iam::444455554444:role/role-name"]
          }
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:iam::012345678901:role/role-name"]
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:iam::012345678901:role/role-name",
			"arn:aws:iam::444455554444:role/role-name",
		},
		AllowedPrincipalAccountIds:          []string{"444455554444"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{"Write"},
		PublicStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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
			"AWS": ["arn:aws:iam::012345678901:role/role-name-1", "arn:aws:iam::012345678901:role/role-name-2"]
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:            "private",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:iam::012345678901:role/role-name-1",
			"arn:aws:iam::012345678901:role/role-name-2",
		},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{"Write"},
		PublicStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:iam::012345678901:role/role-name-1",
			"arn:aws:iam::012345678901:role/role-name-2",
			"arn:aws:iam::444455554444:role/role-name-1",
			"arn:aws:iam::444455554444:role/role-name-2",
		},
		AllowedPrincipalAccountIds:          []string{"444455554444"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{"Write"},
		PublicStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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

func testWhenPricipalIsAUserAccountAssumedRole(t *testing.T) {
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
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:sts::012345678901:assumed-role/role-name/role-session-name"},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{"Write"},
		PublicStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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

func testWhenPricipalIsACrossAccountAssumedRole(t *testing.T) {
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
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:sts::444455554444:assumed-role/role-name/role-session-name"},
		AllowedPrincipalAccountIds:          []string{"444455554444"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{"Write"},
		PublicStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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
			"AWS": ["arn:aws:sts::012345678901:assumed-role/role-name/role-session-name-1", "arn:aws:sts::012345678901:assumed-role/role-name/role-session-name-2"]
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:            "private",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:sts::012345678901:assumed-role/role-name/role-session-name-1",
			"arn:aws:sts::012345678901:assumed-role/role-name/role-session-name-2",
		},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{"Write"},
		PublicStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
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
		PublicAccessLevels:                  []string{"Write"},
		PublicStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
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
		PublicAccessLevels:                  []string{"Write"},
		PublicStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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
          }
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:sts::012345678901:assumed-role/role-name/role-session-name"]
          }
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:sts::444455554444:assumed-role/role-name/role-session-name"]
          }
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:sts::012345678901:assumed-role/role-name/role-session-name"]
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:sts::012345678901:assumed-role/role-name/role-session-name",
			"arn:aws:sts::444455554444:assumed-role/role-name/role-session-name",
		},
		AllowedPrincipalAccountIds:          []string{"444455554444"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{"Write"},
		PublicStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:sts::012345678901:assumed-role/role-name/role-session-name-1",
			"arn:aws:sts::012345678901:assumed-role/role-name/role-session-name-2",
			"arn:aws:sts::444455554444:assumed-role/role-name/role-session-name-1",
			"arn:aws:sts::444455554444:assumed-role/role-name/role-session-name-2",
		},
		AllowedPrincipalAccountIds:          []string{"444455554444"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{"Write"},
		PublicStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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

func testWhenPricipalIsAFederatedUser(t *testing.T) {
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
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{"cognito-identity.amazonaws.com"},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{"Write"},
		PublicStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Federated": ["accounts.google.com", "graph.facebook.com"]
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                "private",
		AllowedOrganizationIds:     []string{},
		AllowedPrincipals:          []string{},
		AllowedPrincipalAccountIds: []string{},
		AllowedPrincipalFederatedIdentities: []string{
			"accounts.google.com",
			"graph.facebook.com",
		},
		AllowedPrincipalServices: []string{},
		IsPublic:                 false,
		PublicAccessLevels:       []string{"Write"},
		PublicStatementIds:       []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Federated": ["graph.facebook.com", "accounts.google.com"]
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                "private",
		AllowedOrganizationIds:     []string{},
		AllowedPrincipals:          []string{},
		AllowedPrincipalAccountIds: []string{},
		AllowedPrincipalFederatedIdentities: []string{
			"accounts.google.com",
			"graph.facebook.com",
		},
		AllowedPrincipalServices: []string{},
		IsPublic:                 false,
		PublicAccessLevels:       []string{"Write"},
		PublicStatementIds:       []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Federated": ["graph.facebook.com"]
          }
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Federated": ["accounts.google.com"]
          }
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Federated": ["graph.facebook.com"]
          }
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Federated": ["accounts.google.com"]
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                "private",
		AllowedOrganizationIds:     []string{},
		AllowedPrincipals:          []string{},
		AllowedPrincipalAccountIds: []string{},
		AllowedPrincipalFederatedIdentities: []string{
			"accounts.google.com",
			"graph.facebook.com",
		},
		AllowedPrincipalServices: []string{},
		IsPublic:                 false,
		PublicAccessLevels:       []string{"Write"},
		PublicStatementIds:       []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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

func testWhenPricipalIsAService(t *testing.T) {
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
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ec2.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"Write"},
		PublicStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices: []string{
			"ecs.amazonaws.com",
			"elasticloadbalancing.amazonaws.com",
		},
		IsPublic:           true,
		PublicAccessLevels: []string{"Write"},
		PublicStatementIds: []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices: []string{
			"ecs.amazonaws.com",
			"elasticloadbalancing.amazonaws.com",
		},
		IsPublic:           true,
		PublicAccessLevels: []string{"Write"},
		PublicStatementIds: []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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
          }
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Service": ["ecs.amazonaws.com"]
          }
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Service": ["elasticloadbalancing.amazonaws.com"]
          }
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "Service": ["ecs.amazonaws.com"]
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices: []string{
			"ecs.amazonaws.com",
			"elasticloadbalancing.amazonaws.com",
		},
		IsPublic:           true,
		PublicAccessLevels: []string{"Write"},
		PublicStatementIds: []string{
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
            "AWS": "arn:aws:iam::444455554444:root",
            "Federated": "cognito-identity.amazonaws.com"
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"arn:aws:iam::444455554444:root"},
		AllowedPrincipalAccountIds:          []string{"444455554444"},
		AllowedPrincipalFederatedIdentities: []string{"cognito-identity.amazonaws.com"},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"Write"},
		PublicStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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
            "AWS": ["arn:aws:iam::444455554444:root", "*"],
            "Federated": "cognito-identity.amazonaws.com"
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:            "public",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"*",
			"arn:aws:iam::444455554444:root",
		},
		AllowedPrincipalAccountIds: []string{
			"*",
			"444455554444",
		},
		AllowedPrincipalFederatedIdentities: []string{"cognito-identity.amazonaws.com"},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{"Write"},
		PublicStatementIds:                  []string{"Statement[1]"},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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
          }
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:iam::555544445555:root"],
            "Service": ["elasticloadbalancing.amazonaws.com", "ecs.amazonaws.com"],
            "Federated": ["accounts.google.com"]
          }
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:iam::555544445555:root"],
            "Service": ["ecs.amazonaws.com"],
            "Federated": ["graph.facebook.com", "cognito-identity.amazonaws.com"]
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:            "public",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"arn:aws:iam::012345678901:root",
			"arn:aws:iam::444455554444:root",
			"arn:aws:iam::555544445555:root",
		},
		AllowedPrincipalAccountIds: []string{
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
		IsPublic:           true,
		PublicAccessLevels: []string{"Write"},
		PublicStatementIds: []string{
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
          }
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:iam::555544445555:root"],
            "Service": ["elasticloadbalancing.amazonaws.com", "ecs.amazonaws.com"],
            "Federated": ["accounts.google.com"]
          }
        },
        {
          "Effect": "Allow",
          "Action": "sts:AssumeRole",
          "Principal": {
            "AWS": ["arn:aws:iam::555544445555:root", "*"],
            "Service": ["ecs.amazonaws.com"],
            "Federated": ["graph.facebook.com", "cognito-identity.amazonaws.com"]
          }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
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
		IsPublic:           true,
		PublicAccessLevels: []string{"Write"},
		PublicStatementIds: []string{
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

func testWhenAwsPrincipalIsWildcardedAndEffectDenied(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
		{
		  "Version": "2012-10-17",
		  "Statement": [
			{
			  "Effect": "Deny",
			  "Action": "sts:AssumeRole",
			  "Principal": "*"
			}
		  ]
		}
		`

	expected := EvaluatedPolicy{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		PublicStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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

func TestWhenAwsPrincipalIsWildcardedDeniedButAnotherAccountIsAllowed(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Deny",
          "Action": "sts:AssumeRole",
          "Principal": "*"
        },
        {
			"Effect": "Allow",
			"Action": "sts:AssumeRole",
			"Principal": "012345678901"
		  } 
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		PublicStatementIds:                  []string{},
	}

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

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
          "Principal": { "AWS": "444455554444" }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"444455554444"},
		AllowedPrincipalAccountIds:          []string{"444455554444"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		PublicStatementIds:                  []string{},
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
          "Principal": { "AWS": "012345678901" }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		PublicStatementIds:                  []string{},
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
          "Principal": { "AWS": "*" }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{},
		PublicStatementIds:                  []string{"Sid_Statement_1"},
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
          "Principal": { "AWS": "*" }
        },
        {
          "Sid": "Sid_Statement_2",
          "Effect": "Allow",
          "Principal": { "AWS": "*" }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{},
		PublicStatementIds: []string{
			"Sid_Statement_1",
			"Sid_Statement_2",
		},
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
          "Principal": { "AWS": "*" }
        },
        {
          "Sid": "Sid_Statement_1",
          "Effect": "Allow",
          "Principal": { "AWS": "*" }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{},
		PublicStatementIds: []string{
			"Sid_Statement_1",
			"Sid_Statement_2",
		},
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
          "Principal": { "AWS": "*" }
        },
        {
          "Sid": "Sid_Statement_1",
          "Effect": "Allow",
          "Principal": { "AWS": "*" }
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
		t.Fatalf("The error message returned is expected to be: %s", expectedErrorMsg)
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
          "Principal": { "AWS": "*" }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
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
          "Principal": { "AWS": "*" }
        },
        {
          "Effect": "Allow",
          "Principal": { "AWS": "*" }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{},
		PublicStatementIds: []string{
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
	t.Run("TestAccessLevelSharedHasHigherPrecidenceThanAccessLevelPrivate", testAccessLevelSharedHasHigherPrecidenceThanAccessLevelPrivate)
	t.Run("TestAccessLevelPublicHasHigherPrecidenceThanAccessLevelPrivate", testAccessLevelPublicHasHigherPrecidenceThanAccessLevelPrivate)
	t.Run("TestAccessLevelPublicHasHigherPrecidenceThanAccessLevelShared", testAccessLevelPublicHasHigherPrecidenceThanAccessLevelShared)
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
          "Principal": { "AWS": "*" }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"*"},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
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
          "Principal": { "Service": ["ecs.amazonaws.com"] }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "public",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{"ecs.amazonaws.com"},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
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
          "Principal": { "AWS": "111122221111" }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "shared",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"111122221111"},
		AllowedPrincipalAccountIds:          []string{"111122221111"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		PublicStatementIds:                  []string{},
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
          "Principal": { "AWS": "012345678901" }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		PublicStatementIds:                  []string{},
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

func testAccessLevelSharedHasHigherPrecidenceThanAccessLevelPrivate(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "111122221111" }
        },
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:            "shared",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"012345678901",
			"111122221111",
		},
		AllowedPrincipalAccountIds:          []string{"111122221111"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		PublicStatementIds:                  []string{},
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

func testAccessLevelPublicHasHigherPrecidenceThanAccessLevelPrivate(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "*" }
        },
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:            "public",
		AllowedOrganizationIds: []string{},
		AllowedPrincipals: []string{
			"*",
			"012345678901",
		},
		AllowedPrincipalAccountIds:          []string{"*"},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            true,
		PublicAccessLevels:                  []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
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

func testAccessLevelPublicHasHigherPrecidenceThanAccessLevelShared(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "*" }
        },
        {
          "Effect": "Allow",
          "Principal": { "AWS": "111122221111" }
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
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
		PublicAccessLevels:                  []string{},
		PublicStatementIds:                  []string{"Statement[1]"},
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
	t.Run("TestUnknownApiService", testUnknownApiService)
	t.Run("TestUnknownApiFunction", testUnknownApiFunction)
	t.Run("TestKnownApiFunction", testKnownApiFunction)

	t.Run("TestMultipleStatementsWithKnownApiFunctions", testMultipleStatementsWithKnownApiFunctions)

	t.Run("TestSingleFullWildcardWithNoActionName", testSingleFullWildcardWithNoActionName)
	t.Run("TestSingleFullWildcardAtEndOfAction", testSingleFullWildcardAtFrontOfAction)
	t.Run("TestSingleFullWildcardAtEndOfAction", testSingleFullWildcardInMiddleOfAction)
	t.Run("TestSingleFullWildcardAtEndOfAction", testSingleFullWildcardAtEndOfAction)

	t.Run("TestSinglePartialWildcardAtFrontOfAction", testSinglePartialWildcardAtFrontOfAction)
	t.Run("TestSinglePartialWildcardInMiddleOfAction", testSinglePartialWildcardInMiddleOfAction)
	t.Run("TestSinglePartialWildcardAtEndOfAction", testSinglePartialWildcardAtEndOfAction)
	t.Run("TestMultipleWildcardsInAction", TestMultipleWildcardsInAction)

	t.Run("TestSinglePartialWildcardAtEndOfKnownApiFunctionAction", testSinglePartialWildcardAtEndOfKnownApiFunctionAction)
	t.Run("TestSingleFullWildcardAtEndOfKnownApiFunctionAction", testSingleFullWildcardAtEndOfKnownApiFunctionAction)

	t.Run("TestIncompleteActionMissingFunctionPattern", testIncompleteActionMissingFunctionPattern)
	t.Run("TestActionWhenServiceNameIsGivenOnly", testActionWhenServiceNameIsGivenOnly)
}

func testUnknownApiService(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec:StartInstances"
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		PublicStatementIds:                  []string{},
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

func testUnknownApiFunction(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2:ZtartInztancez"
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		PublicStatementIds:                  []string{},
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

func testKnownApiFunction(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2:DescribeVolumesModifications"
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{"Read"},
		PublicStatementIds:                  []string{},
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

func testMultipleStatementsWithKnownApiFunctions(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2:DescribeVolumes"
        },
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2:DescribeVolumesModifications"
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels: []string{
			"List",
			"Read",
		},
		PublicStatementIds: []string{},
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

func testSingleFullWildcardWithNoActionName(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2:*"
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels: []string{
			"List",
			"Permissions management",
			"Read",
			"Tagging",
			"Write",
		},
		PublicStatementIds: []string{},
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

func testSingleFullWildcardAtFrontOfAction(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2:*VpcClassicLink"
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels: []string{
			"List",
			"Write",
		},
		PublicStatementIds: []string{},
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

func testSingleFullWildcardInMiddleOfAction(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2:Describe*Attributes"
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{"List"},
		PublicStatementIds:                  []string{},
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

func testSingleFullWildcardAtEndOfAction(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2:Describe*"
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels: []string{
			"List",
			"Read",
		},
		PublicStatementIds: []string{},
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

func testSinglePartialWildcardAtFrontOfAction(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2:?escribeVolumesModifications"
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{"Read"},
		PublicStatementIds:                  []string{},
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

func testSinglePartialWildcardInMiddleOfAction(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2:DescribeVolumes?odifications"
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{"Read"},
		PublicStatementIds:                  []string{},
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

func testSinglePartialWildcardAtEndOfAction(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
		  "Action": "ec2:DescribeVolumesModification?"
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{"Read"},
		PublicStatementIds:                  []string{},
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

func TestMultipleWildcardsInAction(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
		  "Action": "ec2:*Volumes*?"
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels: []string{
			"List",
			"Read",
		},
		PublicStatementIds: []string{},
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

func testSinglePartialWildcardAtEndOfKnownApiFunctionAction(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2:StartInstances?"
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		PublicStatementIds:                  []string{},
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

func testSingleFullWildcardAtEndOfKnownApiFunctionAction(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2:StartInstances*"
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{"Write"},
		PublicStatementIds:                  []string{},
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

func testIncompleteActionMissingFunctionPattern(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2:"
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		PublicStatementIds:                  []string{},
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

func testActionWhenServiceNameIsGivenOnly(t *testing.T) {
	// Set up
	userAccountId := "012345678901"
	policyContent := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": { "AWS": "012345678901" },
          "Action": "ec2"
        }
      ]
    }
	`

	expected := EvaluatedPolicy{
		AccessLevel:                         "private",
		AllowedOrganizationIds:              []string{},
		AllowedPrincipals:                   []string{"012345678901"},
		AllowedPrincipalAccountIds:          []string{},
		AllowedPrincipalFederatedIdentities: []string{},
		AllowedPrincipalServices:            []string{},
		IsPublic:                            false,
		PublicAccessLevels:                  []string{},
		PublicStatementIds:                  []string{},
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
