package aws

import (
	"testing"
)

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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error was returned from EvaluatePolicy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 0
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 0
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
		t.Fail()
	}
}

func testPolicyCreatedByStringEvaluates(t *testing.T) {
	// Set up
	userAccountId := "123456789012"

	// Test
	evaluated, err := EvaluatePolicy("", userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error was returned from EvaluatePolicy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 0
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 0
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
		t.Fail()
	}
}

func testPolicyCreatedByEmptyJsonStringEvaluates(t *testing.T) {
	// Set up
	userAccountId := "123456789012"

	// Test
	evaluated, err := EvaluatePolicy("{}", userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error was returned from EvaluatePolicy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 0
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 0
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
		t.Fail()
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error was returned from EvaluatePolicy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 0
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 0
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy("", userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error was returned from EvaluatePolicy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 0
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 0
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
		t.Fail()
	}
}

func testIfSourceAccountIdContainsCorrectAmountOfNumericalValuesAndStartsWithZeroItEvaluates(t *testing.T) {
	// Set up
	userAccountId := "012345678901"

	// Test
	evaluated, err := EvaluatePolicy("", userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error was returned from EvaluatePolicy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 0
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 0
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
		t.Fail()
	}
}

func TestPolicyPrincipalElement(t *testing.T) {
	t.Run("TestWhenPricipalIsAMisformedArnFails", testWhenPricipalIsAMisformedArnFails)
	t.Run("TestWhenPrincipalIsWildcarded", testWhenPrincipalIsWildcarded)
	t.Run("TestWhenAwsPrincipalIsWildcarded", testWhenAwsPrincipalIsWildcarded)
	t.Run("TestWhenPrincipalIsAUserAccountId", testWhenPrincipalIsAUserAccountId)
	t.Run("TestWhenPrincipalIsAUserAccountArn", testWhenPrincipalIsAUserAccountArn)
	t.Run("TestWhenPrincipalIsACrossAccountId", testWhenPrincipalIsACrossAccountId)
	t.Run("TestWhenPrincipalIsACrossAccountArn", testWhenPrincipalIsACrossAccountArn)
	t.Run("TestWhenPrincipalIsMultipleUserAccounts", testWhenPrincipalIsMultipleUserAccounts)
	t.Run("TestWhenPrincipalIsMultipleCrossAccountsInAscendingOrder", testWhenPrincipalIsMultipleCrossAccountsInAscendingOrder)
	t.Run("TestWhenPrincipalIsMultipleCrossAccountsInDescendingOrder", testWhenPrincipalIsMultipleCrossAccountsInDescendingOrder)
	t.Run("TestWhenPrincipalIsMultipleAccountsPrincipalsAcrossMultipleStatements", testWhenPrincipalIsMultipleAccountsPrincipalsAcrossMultipleStatements)
	t.Run("TestWhenPrincipalIsMultipleMixedAccounts", testWhenPrincipalIsMultipleMixedAccounts)
	// TODO: Questions
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
	t.Run("TestWhenPrincipalIsMultipleAccountAssumedRolePrincipalsAcrossMultipleStatements", testWhenPrincipalIsMultipleFederatedUserPrincipalsAcrossMultipleStatements)

	t.Run("TestWhenPricipalIsAService", testWhenPricipalIsAService)
	t.Run("TestWhenPrincipalIsMulitpleServicesInAscendingOrder", testWhenPrincipalIsMulitpleServicesInAscendingOrder)
	t.Run("TestWhenPrincipalIsMulitpleServicesInDescendingOrder", testWhenPrincipalIsMulitpleServicesInDescendingOrder)
	t.Run("TestWhenPrincipalIsMultipleServicePrincipalsAcrossMultipleStatements", testWhenPrincipalIsMultipleServicePrincipalsAcrossMultipleStatements)

	t.Run("TestWhenPrincipalIsMultipleTypes", testWhenPrincipalIsMultipleTypes)
	t.Run("TestWhenPrincipalIsMultipleTypesWithWildcard", testWhenPrincipalIsMultipleTypesWithWildcard)
	t.Run("TestWhenPrincipalIsMultipleTypesAcrossMultipleStatements", testWhenPrincipalIsMultipleTypesAcrossMultipleStatements)
	t.Run("TestWhenPrincipalIsMultipleTypesAcrossMultipleStatementsWithWildcard", testWhenPrincipalIsMultipleTypesAcrossMultipleStatementsWithWildcard)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 1
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals := evaluated.AllowedPrincipals[0]
	expectedAllowedPrincipals := "012345678901"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 0
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 1
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals := evaluated.AllowedPrincipals[0]
	expectedAllowedPrincipals := "*"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 1
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds := evaluated.AllowedPrincipalAccountIds[0]
	expectedAllowedPrincipalAccountIds := "*"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := true
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 1
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals := evaluated.AllowedPrincipals[0]
	expectedAllowedPrincipals := "*"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 1
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds := evaluated.AllowedPrincipalAccountIds[0]
	expectedAllowedPrincipalAccountIds := "*"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := true
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 1
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals := evaluated.AllowedPrincipals[0]
	expectedAllowedPrincipals := "arn:aws:iam::012345678901:root"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 0
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 1
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals := evaluated.AllowedPrincipals[0]
	expectedAllowedPrincipals := "arn:aws:iam::444455554444:root"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 1
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds := evaluated.AllowedPrincipalAccountIds[0]
	expectedAllowedPrincipalAccountIds := "444455554444"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 1
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals := evaluated.AllowedPrincipals[0]
	expectedAllowedPrincipals := "444455554444"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 1
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds := evaluated.AllowedPrincipalAccountIds[0]
	expectedAllowedPrincipalAccountIds := "444455554444"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 2
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals := evaluated.AllowedPrincipals[0]
	expectedAllowedPrincipals := "012345678901"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals = evaluated.AllowedPrincipals[1]
	expectedAllowedPrincipals = "arn:aws:iam::012345678901:root"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 0
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 2
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals := evaluated.AllowedPrincipals[0]
	expectedAllowedPrincipals := "arn:aws:iam::012345678901:root"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals = evaluated.AllowedPrincipals[1]
	expectedAllowedPrincipals = "arn:aws:iam::444455554444:root"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 1
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds := evaluated.AllowedPrincipalAccountIds[0]
	expectedAllowedPrincipalAccountIds := "444455554444"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 2
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals := evaluated.AllowedPrincipals[0]
	expectedAllowedPrincipals := "444455554444"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals = evaluated.AllowedPrincipals[1]
	expectedAllowedPrincipals = "arn:aws:iam::555544445555:root"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 2
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds := evaluated.AllowedPrincipalAccountIds[0]
	expectedAllowedPrincipalAccountIds := "444455554444"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds = evaluated.AllowedPrincipalAccountIds[1]
	expectedAllowedPrincipalAccountIds = "555544445555"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 2
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals := evaluated.AllowedPrincipals[0]
	expectedAllowedPrincipals := "444455554444"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals = evaluated.AllowedPrincipals[1]
	expectedAllowedPrincipals = "arn:aws:iam::555544445555:root"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 2
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds := evaluated.AllowedPrincipalAccountIds[0]
	expectedAllowedPrincipalAccountIds := "444455554444"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds = evaluated.AllowedPrincipalAccountIds[1]
	expectedAllowedPrincipalAccountIds = "555544445555"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 4
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals := evaluated.AllowedPrincipals[0]
	expectedAllowedPrincipals := "012345678901"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals = evaluated.AllowedPrincipals[1]
	expectedAllowedPrincipals = "444455554444"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals = evaluated.AllowedPrincipals[2]
	expectedAllowedPrincipals = "arn:aws:iam::012345678901:root"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals = evaluated.AllowedPrincipals[3]
	expectedAllowedPrincipals = "arn:aws:iam::444455554444:root"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 1
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds := evaluated.AllowedPrincipalAccountIds[0]
	expectedAllowedPrincipalAccountIds := "444455554444"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 5
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals := evaluated.AllowedPrincipals[0]
	expectedAllowedPrincipals := "*"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals = evaluated.AllowedPrincipals[1]
	expectedAllowedPrincipals = "012345678901"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals = evaluated.AllowedPrincipals[2]
	expectedAllowedPrincipals = "444455554444"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals = evaluated.AllowedPrincipals[3]
	expectedAllowedPrincipals = "arn:aws:iam::012345678901:root"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals = evaluated.AllowedPrincipals[4]
	expectedAllowedPrincipals = "arn:aws:iam::444455554444:root"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 2
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds := evaluated.AllowedPrincipalAccountIds[0]
	expectedAllowedPrincipalAccountIds := "*"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds = evaluated.AllowedPrincipalAccountIds[1]
	expectedAllowedPrincipalAccountIds = "444455554444"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := true
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 1
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals := evaluated.AllowedPrincipals[0]
	expectedAllowedPrincipals := "arn:aws:iam::012345678901:role/role-name"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 0
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 1
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals := evaluated.AllowedPrincipals[0]
	expectedAllowedPrincipals := "arn:aws:iam::444455554444:role/role-name"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 1
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds := evaluated.AllowedPrincipalAccountIds[0]
	expectedAllowedPrincipalAccountIds := "444455554444"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 2
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals := evaluated.AllowedPrincipals[0]
	expectedAllowedPrincipals := "arn:aws:iam::444455554444:role/role-name"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals = evaluated.AllowedPrincipals[1]
	expectedAllowedPrincipals = "arn:aws:iam::555544445555:role/role-name"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 2
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds := evaluated.AllowedPrincipalAccountIds[0]
	expectedAllowedPrincipalAccountIds := "444455554444"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds = evaluated.AllowedPrincipalAccountIds[1]
	expectedAllowedPrincipalAccountIds = "555544445555"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 2
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals := evaluated.AllowedPrincipals[0]
	expectedAllowedPrincipals := "arn:aws:iam::444455554444:role/role-name"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals = evaluated.AllowedPrincipals[1]
	expectedAllowedPrincipals = "arn:aws:iam::555544445555:role/role-name"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 2
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds := evaluated.AllowedPrincipalAccountIds[0]
	expectedAllowedPrincipalAccountIds := "444455554444"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds = evaluated.AllowedPrincipalAccountIds[1]
	expectedAllowedPrincipalAccountIds = "555544445555"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 2
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals := evaluated.AllowedPrincipals[0]
	expectedAllowedPrincipals := "arn:aws:iam::012345678901:role/role-name"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals = evaluated.AllowedPrincipals[1]
	expectedAllowedPrincipals = "arn:aws:iam::444455554444:role/role-name"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 1
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds := evaluated.AllowedPrincipalAccountIds[0]
	expectedAllowedPrincipalAccountIds := "444455554444"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 2
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals := evaluated.AllowedPrincipals[0]
	expectedAllowedPrincipals := "arn:aws:iam::012345678901:role/role-name-1"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals = evaluated.AllowedPrincipals[1]
	expectedAllowedPrincipals = "arn:aws:iam::012345678901:role/role-name-2"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 0
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 4
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals := evaluated.AllowedPrincipals[0]
	expectedAllowedPrincipals := "arn:aws:iam::012345678901:role/role-name-1"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals = evaluated.AllowedPrincipals[1]
	expectedAllowedPrincipals = "arn:aws:iam::012345678901:role/role-name-2"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals = evaluated.AllowedPrincipals[2]
	expectedAllowedPrincipals = "arn:aws:iam::444455554444:role/role-name-1"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals = evaluated.AllowedPrincipals[3]
	expectedAllowedPrincipals = "arn:aws:iam::444455554444:role/role-name-2"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 1
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds := evaluated.AllowedPrincipalAccountIds[0]
	expectedAllowedPrincipalAccountIds := "444455554444"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 1
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals := evaluated.AllowedPrincipals[0]
	expectedAllowedPrincipals := "arn:aws:sts::012345678901:assumed-role/role-name/role-session-name"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 0
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 1
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals := evaluated.AllowedPrincipals[0]
	expectedAllowedPrincipals := "arn:aws:sts::444455554444:assumed-role/role-name/role-session-name"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 1
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds := evaluated.AllowedPrincipalAccountIds[0]
	expectedAllowedPrincipalAccountIds := "444455554444"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 2
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals := evaluated.AllowedPrincipals[0]
	expectedAllowedPrincipals := "arn:aws:sts::012345678901:assumed-role/role-name/role-session-name-1"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals = evaluated.AllowedPrincipals[1]
	expectedAllowedPrincipals = "arn:aws:sts::012345678901:assumed-role/role-name/role-session-name-2"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 0
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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
			"AWS": ["arn:aws:sts::444455554444:assumed-role/role-name/role-session-name", "arn:aws:sts::555544445555:assumed-role/role-name/role-session-name"]
          }
        }
      ]
    }
	`

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 2
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals := evaluated.AllowedPrincipals[0]
	expectedAllowedPrincipals := "arn:aws:sts::444455554444:assumed-role/role-name/role-session-name"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals = evaluated.AllowedPrincipals[1]
	expectedAllowedPrincipals = "arn:aws:sts::555544445555:assumed-role/role-name/role-session-name"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 2
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds := evaluated.AllowedPrincipalAccountIds[0]
	expectedAllowedPrincipalAccountIds := "444455554444"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds = evaluated.AllowedPrincipalAccountIds[1]
	expectedAllowedPrincipalAccountIds = "555544445555"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 2
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals := evaluated.AllowedPrincipals[0]
	expectedAllowedPrincipals := "arn:aws:sts::444455554444:assumed-role/role-name/role-session-name"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals = evaluated.AllowedPrincipals[1]
	expectedAllowedPrincipals = "arn:aws:sts::555544445555:assumed-role/role-name/role-session-name"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 2
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds := evaluated.AllowedPrincipalAccountIds[0]
	expectedAllowedPrincipalAccountIds := "444455554444"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds = evaluated.AllowedPrincipalAccountIds[1]
	expectedAllowedPrincipalAccountIds = "555544445555"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 2
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals := evaluated.AllowedPrincipals[0]
	expectedAllowedPrincipals := "arn:aws:sts::012345678901:assumed-role/role-name/role-session-name"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals = evaluated.AllowedPrincipals[1]
	expectedAllowedPrincipals = "arn:aws:sts::444455554444:assumed-role/role-name/role-session-name"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 1
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds := evaluated.AllowedPrincipalAccountIds[0]
	expectedAllowedPrincipalAccountIds := "444455554444"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 4
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals := evaluated.AllowedPrincipals[0]
	expectedAllowedPrincipals := "arn:aws:sts::012345678901:assumed-role/role-name/role-session-name-1"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals = evaluated.AllowedPrincipals[1]
	expectedAllowedPrincipals = "arn:aws:sts::012345678901:assumed-role/role-name/role-session-name-2"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals = evaluated.AllowedPrincipals[2]
	expectedAllowedPrincipals = "arn:aws:sts::444455554444:assumed-role/role-name/role-session-name-1"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals = evaluated.AllowedPrincipals[3]
	expectedAllowedPrincipals = "arn:aws:sts::444455554444:assumed-role/role-name/role-session-name-2"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 1
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds := evaluated.AllowedPrincipalAccountIds[0]
	expectedAllowedPrincipalAccountIds := "444455554444"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 0
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 0
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 1
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	currentAllowedPrincipalFederatedIdentities := evaluated.AllowedPrincipalFederatedIdentities[0]
	expectedAllowedPrincipalFederatedIdentities := "cognito-identity.amazonaws.com"
	if currentAllowedPrincipalFederatedIdentities != expectedAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities: `%s` AllowedPrincipalFederatedIdentities expected: `%s`", currentAllowedPrincipalFederatedIdentities, expectedAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 0
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 0
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 2
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	currentAllowedPrincipalFederatedIdentities := evaluated.AllowedPrincipalFederatedIdentities[0]
	expectedAllowedPrincipalFederatedIdentities := "accounts.google.com"
	if currentAllowedPrincipalFederatedIdentities != expectedAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities: `%s` AllowedPrincipalFederatedIdentities expected: `%s`", currentAllowedPrincipalFederatedIdentities, expectedAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	currentAllowedPrincipalFederatedIdentities = evaluated.AllowedPrincipalFederatedIdentities[1]
	expectedAllowedPrincipalFederatedIdentities = "graph.facebook.com"
	if currentAllowedPrincipalFederatedIdentities != expectedAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities: `%s` AllowedPrincipalFederatedIdentities expected: `%s`", currentAllowedPrincipalFederatedIdentities, expectedAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 0
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 0
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 2
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	currentAllowedPrincipalFederatedIdentities := evaluated.AllowedPrincipalFederatedIdentities[0]
	expectedAllowedPrincipalFederatedIdentities := "accounts.google.com"
	if currentAllowedPrincipalFederatedIdentities != expectedAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities: `%s` AllowedPrincipalFederatedIdentities expected: `%s`", currentAllowedPrincipalFederatedIdentities, expectedAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	currentAllowedPrincipalFederatedIdentities = evaluated.AllowedPrincipalFederatedIdentities[1]
	expectedAllowedPrincipalFederatedIdentities = "graph.facebook.com"
	if currentAllowedPrincipalFederatedIdentities != expectedAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities: `%s` AllowedPrincipalFederatedIdentities expected: `%s`", currentAllowedPrincipalFederatedIdentities, expectedAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 0
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 0
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 2
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	currentAllowedPrincipalFederatedIdentities := evaluated.AllowedPrincipalFederatedIdentities[0]
	expectedAllowedPrincipalFederatedIdentities := "accounts.google.com"
	if currentAllowedPrincipalFederatedIdentities != expectedAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities: `%s` AllowedPrincipalFederatedIdentities expected: `%s`", currentAllowedPrincipalFederatedIdentities, expectedAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	currentAllowedPrincipalFederatedIdentities = evaluated.AllowedPrincipalFederatedIdentities[1]
	expectedAllowedPrincipalFederatedIdentities = "graph.facebook.com"
	if currentAllowedPrincipalFederatedIdentities != expectedAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities: `%s` AllowedPrincipalFederatedIdentities expected: `%s`", currentAllowedPrincipalFederatedIdentities, expectedAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 0
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := false
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 0
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 0
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 1
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentAllowedPrincipalServices := evaluated.AllowedPrincipalServices[0]
	expectedAllowedPrincipalServices := "ec2.amazonaws.com"
	if currentAllowedPrincipalServices != expectedAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices: `%s` AllowedPrincipalServices expected: `%s`", currentAllowedPrincipalServices, expectedAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := true
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 0
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 0
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 2
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentAllowedPrincipalServices := evaluated.AllowedPrincipalServices[0]
	expectedAllowedPrincipalServices := "ecs.amazonaws.com"
	if currentAllowedPrincipalServices != expectedAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices: `%s` AllowedPrincipalServices expected: `%s`", currentAllowedPrincipalServices, expectedAllowedPrincipalServices)
		t.Fail()
	}

	currentAllowedPrincipalServices = evaluated.AllowedPrincipalServices[1]
	expectedAllowedPrincipalServices = "elasticloadbalancing.amazonaws.com"
	if currentAllowedPrincipalServices != expectedAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices: `%s` AllowedPrincipalServices expected: `%s`", currentAllowedPrincipalServices, expectedAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := true
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 0
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 0
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 2
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentAllowedPrincipalServices := evaluated.AllowedPrincipalServices[0]
	expectedAllowedPrincipalServices := "ecs.amazonaws.com"
	if currentAllowedPrincipalServices != expectedAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices: `%s` AllowedPrincipalServices expected: `%s`", currentAllowedPrincipalServices, expectedAllowedPrincipalServices)
		t.Fail()
	}

	currentAllowedPrincipalServices = evaluated.AllowedPrincipalServices[1]
	expectedAllowedPrincipalServices = "elasticloadbalancing.amazonaws.com"
	if currentAllowedPrincipalServices != expectedAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices: `%s` AllowedPrincipalServices expected: `%s`", currentAllowedPrincipalServices, expectedAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := true
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 0
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 0
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 0
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 2
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentAllowedPrincipalServices := evaluated.AllowedPrincipalServices[0]
	expectedAllowedPrincipalServices := "ecs.amazonaws.com"
	if currentAllowedPrincipalServices != expectedAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices: `%s` AllowedPrincipalServices expected: `%s`", currentAllowedPrincipalServices, expectedAllowedPrincipalServices)
		t.Fail()
	}

	currentAllowedPrincipalServices = evaluated.AllowedPrincipalServices[1]
	expectedAllowedPrincipalServices = "elasticloadbalancing.amazonaws.com"
	if currentAllowedPrincipalServices != expectedAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices: `%s` AllowedPrincipalServices expected: `%s`", currentAllowedPrincipalServices, expectedAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := true
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 1
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals := evaluated.AllowedPrincipals[0]
	expectedAllowedPrincipals := "arn:aws:iam::444455554444:root"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 1
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds := evaluated.AllowedPrincipalAccountIds[0]
	expectedAllowedPrincipalAccountIds := "444455554444"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 1
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	currentAllowedPrincipalFederatedIdentities := evaluated.AllowedPrincipalFederatedIdentities[0]
	expectedAllowedPrincipalFederatedIdentities := "cognito-identity.amazonaws.com"
	if currentAllowedPrincipalFederatedIdentities != expectedAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities: `%s` AllowedPrincipalFederatedIdentities expected: `%s`", currentAllowedPrincipalFederatedIdentities, expectedAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 1
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentAllowedPrincipalServices := evaluated.AllowedPrincipalServices[0]
	expectedAllowedPrincipalServices := "ecs.amazonaws.com"
	if currentAllowedPrincipalServices != expectedAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices: `%s` AllowedPrincipalServices expected: `%s`", currentAllowedPrincipalServices, expectedAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := true
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 2
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals := evaluated.AllowedPrincipals[0]
	expectedAllowedPrincipals := "*"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals = evaluated.AllowedPrincipals[1]
	expectedAllowedPrincipals = "arn:aws:iam::444455554444:root"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 2
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds := evaluated.AllowedPrincipalAccountIds[0]
	expectedAllowedPrincipalAccountIds := "*"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds = evaluated.AllowedPrincipalAccountIds[1]
	expectedAllowedPrincipalAccountIds = "444455554444"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 1
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	currentAllowedPrincipalFederatedIdentities := evaluated.AllowedPrincipalFederatedIdentities[0]
	expectedAllowedPrincipalFederatedIdentities := "cognito-identity.amazonaws.com"
	if currentAllowedPrincipalFederatedIdentities != expectedAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities: `%s` AllowedPrincipalFederatedIdentities expected: `%s`", currentAllowedPrincipalFederatedIdentities, expectedAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 1
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentAllowedPrincipalServices := evaluated.AllowedPrincipalServices[0]
	expectedAllowedPrincipalServices := "ecs.amazonaws.com"
	if currentAllowedPrincipalServices != expectedAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices: `%s` AllowedPrincipalServices expected: `%s`", currentAllowedPrincipalServices, expectedAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := true
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 3
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals := evaluated.AllowedPrincipals[0]
	expectedAllowedPrincipals := "arn:aws:iam::012345678901:root"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals = evaluated.AllowedPrincipals[1]
	expectedAllowedPrincipals = "arn:aws:iam::444455554444:root"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals = evaluated.AllowedPrincipals[2]
	expectedAllowedPrincipals = "arn:aws:iam::555544445555:root"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 2
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds := evaluated.AllowedPrincipalAccountIds[0]
	expectedAllowedPrincipalAccountIds := "444455554444"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds = evaluated.AllowedPrincipalAccountIds[1]
	expectedAllowedPrincipalAccountIds = "555544445555"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 3
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	currentAllowedPrincipalFederatedIdentities := evaluated.AllowedPrincipalFederatedIdentities[0]
	expectedAllowedPrincipalFederatedIdentities := "accounts.google.com"
	if currentAllowedPrincipalFederatedIdentities != expectedAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities: `%s` AllowedPrincipalFederatedIdentities expected: `%s`", currentAllowedPrincipalFederatedIdentities, expectedAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	currentAllowedPrincipalFederatedIdentities = evaluated.AllowedPrincipalFederatedIdentities[1]
	expectedAllowedPrincipalFederatedIdentities = "cognito-identity.amazonaws.com"
	if currentAllowedPrincipalFederatedIdentities != expectedAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities: `%s` AllowedPrincipalFederatedIdentities expected: `%s`", currentAllowedPrincipalFederatedIdentities, expectedAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	currentAllowedPrincipalFederatedIdentities = evaluated.AllowedPrincipalFederatedIdentities[2]
	expectedAllowedPrincipalFederatedIdentities = "graph.facebook.com"
	if currentAllowedPrincipalFederatedIdentities != expectedAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities: `%s` AllowedPrincipalFederatedIdentities expected: `%s`", currentAllowedPrincipalFederatedIdentities, expectedAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 3
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentAllowedPrincipalServices := evaluated.AllowedPrincipalServices[0]
	expectedAllowedPrincipalServices := "dynamodb.amazonaws.com"
	if currentAllowedPrincipalServices != expectedAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices: `%s` AllowedPrincipalServices expected: `%s`", currentAllowedPrincipalServices, expectedAllowedPrincipalServices)
		t.Fail()
	}

	currentAllowedPrincipalServices = evaluated.AllowedPrincipalServices[1]
	expectedAllowedPrincipalServices = "ecs.amazonaws.com"
	if currentAllowedPrincipalServices != expectedAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices: `%s` AllowedPrincipalServices expected: `%s`", currentAllowedPrincipalServices, expectedAllowedPrincipalServices)
		t.Fail()
	}

	currentAllowedPrincipalServices = evaluated.AllowedPrincipalServices[2]
	expectedAllowedPrincipalServices = "elasticloadbalancing.amazonaws.com"
	if currentAllowedPrincipalServices != expectedAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices: `%s` AllowedPrincipalServices expected: `%s`", currentAllowedPrincipalServices, expectedAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := true
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
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

	// Test
	evaluated, err := EvaluatePolicy(policyContent, userAccountId)

	// Evaluate
	if err != nil {
		t.Fatalf("Unexpected error while evaluating policy: %s", err)
	}

	currentAccessLevel := evaluated.AccessLevel
	expectedAccessLevel := ""
	if currentAccessLevel != expectedAccessLevel {
		t.Logf("Unexpected AccessLevel: `%s` AccessLevel expected: `%s`", currentAccessLevel, expectedAccessLevel)
		t.Fail()
	}

	countAllowedOrganizationIds := len(evaluated.AllowedOrganizationIds)
	expectedCountAllowedOrganizationIds := 0
	if countAllowedOrganizationIds != expectedCountAllowedOrganizationIds {
		t.Logf("Unexpected AllowedOrganizationIds has: `%d` entries but: `%d` expected", countAllowedOrganizationIds, expectedCountAllowedOrganizationIds)
		t.Fail()
	}

	countAllowedPrincipals := len(evaluated.AllowedPrincipals)
	expectedCountAllowedPrincipals := 4
	if countAllowedPrincipals != expectedCountAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals has: `%d` entries but: `%d` expected", countAllowedPrincipals, expectedCountAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals := evaluated.AllowedPrincipals[0]
	expectedAllowedPrincipals := "*"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals = evaluated.AllowedPrincipals[1]
	expectedAllowedPrincipals = "arn:aws:iam::012345678901:root"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals = evaluated.AllowedPrincipals[2]
	expectedAllowedPrincipals = "arn:aws:iam::444455554444:root"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	currentAllowedPrincipals = evaluated.AllowedPrincipals[3]
	expectedAllowedPrincipals = "arn:aws:iam::555544445555:root"
	if currentAllowedPrincipals != expectedAllowedPrincipals {
		t.Logf("Unexpected AllowedPrincipals: `%s` AllowedPrincipals expected: `%s`", currentAllowedPrincipals, expectedAllowedPrincipals)
		t.Fail()
	}

	countAllowedPrincipalAccountIds := len(evaluated.AllowedPrincipalAccountIds)
	expectedCountAllowedPrincipalAccountIds := 3
	if countAllowedPrincipalAccountIds != expectedCountAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds has: `%d` entries but: `%d` expected", countAllowedPrincipalAccountIds, expectedCountAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds := evaluated.AllowedPrincipalAccountIds[0]
	expectedAllowedPrincipalAccountIds := "*"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds = evaluated.AllowedPrincipalAccountIds[1]
	expectedAllowedPrincipalAccountIds = "444455554444"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	currentAllowedPrincipalAccountIds = evaluated.AllowedPrincipalAccountIds[2]
	expectedAllowedPrincipalAccountIds = "555544445555"
	if currentAllowedPrincipalAccountIds != expectedAllowedPrincipalAccountIds {
		t.Logf("Unexpected AllowedPrincipalAccountIds: `%s` AllowedPrincipalAccountIds expected: `%s`", currentAllowedPrincipalAccountIds, expectedAllowedPrincipalAccountIds)
		t.Fail()
	}

	countAllowedPrincipalFederatedIdentities := len(evaluated.AllowedPrincipalFederatedIdentities)
	expectedCountAllowedPrincipalFederatedIdentities := 3
	if countAllowedPrincipalFederatedIdentities != expectedCountAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities has: `%d` entries but: `%d` expected", countAllowedPrincipalFederatedIdentities, expectedCountAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	currentAllowedPrincipalFederatedIdentities := evaluated.AllowedPrincipalFederatedIdentities[0]
	expectedAllowedPrincipalFederatedIdentities := "accounts.google.com"
	if currentAllowedPrincipalFederatedIdentities != expectedAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities: `%s` AllowedPrincipalFederatedIdentities expected: `%s`", currentAllowedPrincipalFederatedIdentities, expectedAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	currentAllowedPrincipalFederatedIdentities = evaluated.AllowedPrincipalFederatedIdentities[1]
	expectedAllowedPrincipalFederatedIdentities = "cognito-identity.amazonaws.com"
	if currentAllowedPrincipalFederatedIdentities != expectedAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities: `%s` AllowedPrincipalFederatedIdentities expected: `%s`", currentAllowedPrincipalFederatedIdentities, expectedAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	currentAllowedPrincipalFederatedIdentities = evaluated.AllowedPrincipalFederatedIdentities[2]
	expectedAllowedPrincipalFederatedIdentities = "graph.facebook.com"
	if currentAllowedPrincipalFederatedIdentities != expectedAllowedPrincipalFederatedIdentities {
		t.Logf("Unexpected AllowedPrincipalFederatedIdentities: `%s` AllowedPrincipalFederatedIdentities expected: `%s`", currentAllowedPrincipalFederatedIdentities, expectedAllowedPrincipalFederatedIdentities)
		t.Fail()
	}

	countAllowedPrincipalServices := len(evaluated.AllowedPrincipalServices)
	expectedCountAllowedPrincipalServices := 3
	if countAllowedPrincipalServices != expectedCountAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices has: `%d` entries but: `%d` expected", countAllowedPrincipalServices, expectedCountAllowedPrincipalServices)
		t.Fail()
	}

	currentAllowedPrincipalServices := evaluated.AllowedPrincipalServices[0]
	expectedAllowedPrincipalServices := "dynamodb.amazonaws.com"
	if currentAllowedPrincipalServices != expectedAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices: `%s` AllowedPrincipalServices expected: `%s`", currentAllowedPrincipalServices, expectedAllowedPrincipalServices)
		t.Fail()
	}

	currentAllowedPrincipalServices = evaluated.AllowedPrincipalServices[1]
	expectedAllowedPrincipalServices = "ecs.amazonaws.com"
	if currentAllowedPrincipalServices != expectedAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices: `%s` AllowedPrincipalServices expected: `%s`", currentAllowedPrincipalServices, expectedAllowedPrincipalServices)
		t.Fail()
	}

	currentAllowedPrincipalServices = evaluated.AllowedPrincipalServices[2]
	expectedAllowedPrincipalServices = "elasticloadbalancing.amazonaws.com"
	if currentAllowedPrincipalServices != expectedAllowedPrincipalServices {
		t.Logf("Unexpected AllowedPrincipalServices: `%s` AllowedPrincipalServices expected: `%s`", currentAllowedPrincipalServices, expectedAllowedPrincipalServices)
		t.Fail()
	}

	currentIsPublic := evaluated.IsPublic
	expectedIsPublic := true
	if currentIsPublic != expectedIsPublic {
		t.Logf("Unexpected IsPublic: `%t` IsPublic expected: `%t`", currentIsPublic, expectedIsPublic)
		t.Fail()
	}

	countPublicAccessLevels := len(evaluated.PublicAccessLevels)
	expectedCountPublicAccessLevels := 0
	if countPublicAccessLevels != expectedCountPublicAccessLevels {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicAccessLevels, expectedCountPublicAccessLevels)
		t.Fail()
	}

	countPublicStatementIds := len(evaluated.PublicStatementIds)
	expectedCountPublicStatementIds := 0
	if countPublicStatementIds != expectedCountPublicStatementIds {
		t.Logf("Unexpected PublicAccessLevels has: `%d` entries but: `%d` expected", countPublicStatementIds, expectedCountPublicStatementIds)
		t.Fail()
	}
}
