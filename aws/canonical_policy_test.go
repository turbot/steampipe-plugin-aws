package aws

import (
	"testing"
)

func TestNoPolicyFails(t *testing.T) {
	// Set up
	var policyContent string

	// Test
	_, err := canonicalPolicy(policyContent)

	// Expected
	if err == nil {
		t.Fatal("An error is expected when parsing no content")
	}

	expectedErrorMsg := "canonicalPolicy failed unmarshalling source data: unexpected end of JSON input.  src: "

	if errorMsg := err.Error(); errorMsg != expectedErrorMsg {
		t.Fatalf("The error message returned is expected to be: %s", expectedErrorMsg)
	}
}

func TestEmptyPolicy(t *testing.T) {
	// Set up
	policyContent := `{}`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	expectedVersion := ""
	if policy.Version != expectedVersion {
		t.Logf("Unexpected Policy Version: `%s` Policy Version expected: `%s`", policy.Version, expectedVersion)
		t.Fail()
	}

	expectedId := ""
	if policy.Id != expectedId {
		t.Logf("Unexpected Policy ID: `%s` Policy ID expected: `%s`", policy.Id, expectedId)
		t.Fail()
	}

	expectedStatements := 0
	numberStatements := len(policy.Statements)
	if numberStatements != expectedStatements {
		t.Logf("Unexpected number of Statements: `%d` Statements expected: `%d`", numberStatements, expectedStatements)
		t.Fail()
	}
}

func TestVersionElement(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Version": "2012-10-17"
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	expectedVersion := "2012-10-17"
	if policy.Version != expectedVersion {
		t.Fatalf("Unexpected Policy Version: `%s` Policy Version expected: `%s`", policy.Version, expectedVersion)
	}
}

func TestIdElement(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Id": "cd3ad3d9-2776-4ef1-a904-4c229d1642ee"
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	expectedId := "cd3ad3d9-2776-4ef1-a904-4c229d1642ee"
	if policy.Id != expectedId {
		t.Fatalf("Unexpected Policy ID: `%s` Policy ID expected: `%s`", policy.Id, expectedId)
	}
}

func TestSingleEmptyStatement(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
        }
      ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	expectedStatements := 1
	numberStatements := len(policy.Statements)
	if numberStatements != expectedStatements {
		t.Fatalf("Unexpected number of Statements: `%d` Statements expected: `%d`", numberStatements, expectedStatements)
	}

	statement := policy.Statements[0]

	expectedEffect := ""
	if statement.Effect != expectedEffect {
		t.Logf("Unexpected Statement Effect: `%s` Statement Effect expected: `%s`", statement.Effect, expectedEffect)
		t.Fail()
	}

	expectedSid := ""
	if statement.Sid != expectedSid {
		t.Logf("Unexpected Statement SID: `%s` Statement SID expected: `%s`", statement.Sid, expectedSid)
		t.Fail()
	}

	if len(statement.Action) > 0 {
		t.Logf("Unexpected Statement Action: `%s`", statement.Action)
		t.Fail()
	}

	if len(statement.NotAction) > 0 {
		t.Logf("Unexpected Statement NotAction: `%s`", statement.NotAction)
		t.Fail()
	}

	if len(statement.Principal) > 0 {
		t.Logf("Unexpected Statement Principal: `%s`", statement.Principal)
		t.Fail()
	}

	if len(statement.NotPrincipal) > 0 {
		t.Logf("Unexpected Statement NotPrincipal: `%s`", statement.NotPrincipal)
		t.Fail()
	}

	if len(statement.NotResource) > 0 {
		t.Logf("Unexpected Statement NotResource: `%s`", statement.NotResource)
		t.Fail()
	}

	if len(statement.Resource) > 0 {
		t.Logf("Unexpected Statement Principal: `%s`", statement.Principal)
		t.Fail()
	}

	if len(statement.Condition) > 0 {
		t.Logf("Unexpected Statement Condition: `%s`", statement.Condition)
		t.Fail()
	}
}

func TestMultipleEmptyStatements(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
        },
        {
        },
        {
        }
      ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	expectedStatements := 3
	numberStatements := len(policy.Statements)
	if numberStatements != expectedStatements {
		t.Fatalf("Unexpected number of Statements: `%d` Statements expected: `%d`", numberStatements, expectedStatements)
	}

	for index, statement := range policy.Statements {
		expectedEffect := ""
		if statement.Effect != expectedEffect {
			t.Logf("Unexpected Statement Effect: `%s` Statement Effect expected: `%s` for index: %d", statement.Effect, expectedEffect, index)
			t.Fail()
		}

		expectedSid := ""
		if statement.Sid != expectedSid {
			t.Logf("Unexpected Statement SID: `%s` Statement SID expected: `%s`", statement.Sid, expectedSid)
			t.Fail()
		}

		if len(statement.Action) > 0 {
			t.Logf("Unexpected Statement Action: `%s`", statement.Action)
			t.Fail()
		}

		if len(statement.NotAction) > 0 {
			t.Logf("Unexpected Statement NotAction: `%s`", statement.NotAction)
			t.Fail()
		}

		if len(statement.Principal) > 0 {
			t.Logf("Unexpected Statement Principal: `%s`", statement.Principal)
			t.Fail()
		}

		if len(statement.NotPrincipal) > 0 {
			t.Logf("Unexpected Statement NotPrincipal: `%s`", statement.NotPrincipal)
			t.Fail()
		}

		if len(statement.NotResource) > 0 {
			t.Logf("Unexpected Statement NotResource: `%s`", statement.NotResource)
			t.Fail()
		}

		if len(statement.Resource) > 0 {
			t.Logf("Unexpected Statement Principal: `%s`", statement.Principal)
			t.Fail()
		}

		if len(statement.Condition) > 0 {
			t.Logf("Unexpected Statement Condition: `%s`", statement.Condition)
			t.Fail()
		}
	}
}

func TestStatementPrincipalElement(t *testing.T) {
	t.Run("TestStatementPrincipalElementSingleValue", testStatementPrincipalElementSingleValue)
	t.Run("TestStatementPrincipalElementMultipleValue", testStatementPrincipalElementMultipleValue)
	t.Run("TestStatementPrincipalElementMultipleValueAtPrincipalFails", testStatementPrincipalElementMultipleValueAtPrincipalFails)
	t.Run("TestStatementPrincipalElementRemoveDeplicateValues", testStatementPrincipalElementRemoveDuplicateValues)
	t.Run("TestStatementPrincipalElementSortValues", testStatementPrincipalElementSortValues)
	t.Run("TestStatementPrincipalElementValuesOtherThanStringsFails", testStatementPrincipalElementValuesOtherThanStringsFails)
}

func testStatementPrincipalElementValuesOtherThanStringsFails(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "Principal": 10
        }
      ]
    }
	`

	// Test
	_, err := canonicalPolicy(policyContent)

	// Expected
	if err == nil {
		t.Fatal("An error is expected when parsing no content")
	}

	expectedErrorMsg := "canonicalPolicy failed unmarshalling source data: UnmarshalJSON failed for Statements (Array of Statement): %5B%0A++++++++%7B%0A++++++++++%22Principal%22%3A+10%0A++++++++%7D%0A++++++%5D.  src: %0A++++%7B%0A++++++%22Statement%22%3A+%5B%0A++++++++%7B%0A++++++++++%22Principal%22%3A+10%0A++++++++%7D%0A++++++%5D%0A++++%7D%0A%09"

	if errorMsg := err.Error(); errorMsg != expectedErrorMsg {
		t.Fatalf("The error message returned is expected to be: %s", expectedErrorMsg)
	}
}

func testStatementPrincipalElementSingleValue(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "Principal": "*"
        },
        {
          "Principal": { "AWS": "*" }
        },
        {
          "Principal": { "CanonicalUser": "*" }
        },
        {
          "Principal": { "Service": "*" }
        },
        {
          "Principal": { "Any Value": "*" }
        }
      ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	expectedValue := "*"
	if policy.Statements[0].Principal["AWS"].([]string)[0] != expectedValue {
		t.Logf("Unexpected Statement Principal: `%s` Statement Principal expected: `*`", policy.Statements[0].Principal["AWS"].([]string)[0])
		t.Fail()
	}

	expectedEntries := 1
	if len(policy.Statements[0].Principal) != expectedEntries {
		t.Logf("Unexpected Statement Principal has: `%d` entries but: `%d` expected", len(policy.Statements[0].Principal), expectedEntries)
		t.Fail()
	}

	if policy.Statements[1].Principal["AWS"].([]string)[0] != expectedValue {
		t.Logf("Unexpected Statement Principal: `%s` Statement Principal expected: `*`", policy.Statements[1].Principal["AWS"].([]string)[0])
		t.Fail()
	}

	if len(policy.Statements[1].Principal) != expectedEntries {
		t.Logf("Unexpected Statement Principal has: `%d` entries but: `%d` expected", len(policy.Statements[1].Principal), expectedEntries)
		t.Fail()
	}

	if policy.Statements[2].Principal["CanonicalUser"].([]string)[0] != expectedValue {
		t.Logf("Unexpected Statement Principal: `%s` Statement Principal expected: `*`", policy.Statements[2].Principal["AWS"].([]string)[0])
		t.Fail()
	}

	if len(policy.Statements[2].Principal) != expectedEntries {
		t.Logf("Unexpected Statement Principal has: `%d` entries but: `%d` expected", len(policy.Statements[2].Principal), expectedEntries)
		t.Fail()
	}

	if policy.Statements[3].Principal["Service"].([]string)[0] != expectedValue {
		t.Logf("Unexpected Statement Principal: `%s` Statement Principal expected: `*`", policy.Statements[3].Principal["AWS"].([]string)[0])
		t.Fail()
	}

	if len(policy.Statements[3].Principal) != expectedEntries {
		t.Logf("Unexpected Statement Principal has: `%d` entries but: `%d` expected", len(policy.Statements[3].Principal), expectedEntries)
		t.Fail()
	}

	if policy.Statements[4].Principal["Any Value"].([]string)[0] != expectedValue {
		t.Logf("Unexpected Statement Principal: `%s` Statement Principal expected: `*`", policy.Statements[4].Principal["AWS"].([]string)[0])
		t.Fail()
	}

	if len(policy.Statements[4].Principal) != expectedEntries {
		t.Logf("Unexpected Statement Principal has: `%d` entries but: `%d` expected", len(policy.Statements[4].Principal), expectedEntries)
		t.Fail()
	}
}

func testStatementPrincipalElementMultipleValue(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "Principal": { "AWS": ["111199991111", "999911119999"] }
        },
        {
          "Principal": { "CanonicalUser": ["111199991111", "999911119999"] }
        },
        {
          "Principal": { "Service": ["111199991111", "999911119999"] }
        },
        {
          "Principal": { "Any Value": ["111199991111", "999911119999"] }
        }
      ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	expectedValueIndex0 := "111199991111"
	expectedValueIndex1 := "999911119999"
	if policy.Statements[0].Principal["AWS"].([]string)[0] != expectedValueIndex0 {
		t.Logf("Unexpected Statement Principal: `%s` Statement Principal expected: `*`", policy.Statements[0].Principal["AWS"].([]string)[0])
		t.Fail()
	}

	if policy.Statements[0].Principal["AWS"].([]string)[1] != expectedValueIndex1 {
		t.Logf("Unexpected Statement Principal: `%s` Statement Principal expected: `*`", policy.Statements[0].Principal["AWS"].([]string)[1])
		t.Fail()
	}

	expectedEntries := 2
	if len(policy.Statements[0].Principal["AWS"].([]string)) != expectedEntries {
		t.Logf("Unexpected Statement Principal has: `%d` entries but: `%d` expected", len(policy.Statements[1].Principal["AWS"].([]string)), expectedEntries)
		t.Fail()
	}

	if policy.Statements[1].Principal["CanonicalUser"].([]string)[0] != expectedValueIndex0 {
		t.Logf("Unexpected Statement Principal: `%s` Statement Principal expected: `*`", policy.Statements[1].Principal["CanonicalUser"].([]string)[0])
		t.Fail()
	}

	if policy.Statements[1].Principal["CanonicalUser"].([]string)[1] != expectedValueIndex1 {
		t.Logf("Unexpected Statement Principal: `%s` Statement Principal expected: `*`", policy.Statements[1].Principal["CanonicalUser"].([]string)[1])
		t.Fail()
	}

	if len(policy.Statements[1].Principal["CanonicalUser"].([]string)) != expectedEntries {
		t.Logf("Unexpected Statement Principal has: `%d` entries but: `%d` expected", len(policy.Statements[1].Principal["CanonicalUser"].([]string)), expectedEntries)
		t.Fail()
	}

	if policy.Statements[2].Principal["Service"].([]string)[0] != expectedValueIndex0 {
		t.Logf("Unexpected Statement Principal: `%s` Statement Principal expected: `*`", policy.Statements[2].Principal["Service"].([]string)[0])
		t.Fail()
	}

	if policy.Statements[2].Principal["Service"].([]string)[1] != expectedValueIndex1 {
		t.Logf("Unexpected Statement Principal: `%s` Statement Principal expected: `*`", policy.Statements[2].Principal["Service"].([]string)[1])
		t.Fail()
	}

	if len(policy.Statements[2].Principal["Service"].([]string)) != expectedEntries {
		t.Logf("Unexpected Statement Principal has: `%d` entries but: `%d` expected", len(policy.Statements[2].Principal["Service"].([]string)), expectedEntries)
		t.Fail()
	}

	if policy.Statements[3].Principal["Any Value"].([]string)[0] != expectedValueIndex0 {
		t.Logf("Unexpected Statement Principal: `%s` Statement Principal expected: `*`", policy.Statements[3].Principal["Any Value"].([]string)[0])
		t.Fail()
	}

	if policy.Statements[3].Principal["Any Value"].([]string)[1] != expectedValueIndex1 {
		t.Logf("Unexpected Statement Principal: `%s` Statement Principal expected: `*`", policy.Statements[3].Principal["Any Value"].([]string)[1])
		t.Fail()
	}

	if len(policy.Statements[3].Principal["Any Value"].([]string)) != expectedEntries {
		t.Logf("Unexpected Statement Principal has: `%d` entries but: `%d` expected", len(policy.Statements[3].Principal["Any Value"].([]string)), expectedEntries)
		t.Fail()
	}
}

func testStatementPrincipalElementSortValues(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "Principal": { "AWS": ["999911119999", "111199991111"] }
        },
        {
          "Principal": { "CanonicalUser": ["999911119999", "111199991111"] }
        },
        {
          "Principal": { "Service": ["999911119999", "111199991111"] }
        },
        {
          "Principal": { "Any Value": ["999911119999", "111199991111"] }
        }
      ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	expectedValueIndex0 := "111199991111"
	expectedValueIndex1 := "999911119999"
	if policy.Statements[0].Principal["AWS"].([]string)[0] != expectedValueIndex0 {
		t.Logf("Unexpected Statement Principal: `%s` Statement Principal expected: `*`", policy.Statements[0].Principal["AWS"].([]string)[0])
		t.Fail()
	}

	if policy.Statements[0].Principal["AWS"].([]string)[1] != expectedValueIndex1 {
		t.Logf("Unexpected Statement Principal: `%s` Statement Principal expected: `*`", policy.Statements[0].Principal["AWS"].([]string)[1])
		t.Fail()
	}

	expectedEntries := 2
	if len(policy.Statements[0].Principal["AWS"].([]string)) != expectedEntries {
		t.Logf("Unexpected Statement Principal has: `%d` entries but: `%d` expected", len(policy.Statements[1].Principal["AWS"].([]string)), expectedEntries)
		t.Fail()
	}

	if policy.Statements[1].Principal["CanonicalUser"].([]string)[0] != expectedValueIndex0 {
		t.Logf("Unexpected Statement Principal: `%s` Statement Principal expected: `*`", policy.Statements[1].Principal["CanonicalUser"].([]string)[0])
		t.Fail()
	}

	if policy.Statements[1].Principal["CanonicalUser"].([]string)[1] != expectedValueIndex1 {
		t.Logf("Unexpected Statement Principal: `%s` Statement Principal expected: `*`", policy.Statements[1].Principal["CanonicalUser"].([]string)[1])
		t.Fail()
	}

	if len(policy.Statements[1].Principal["CanonicalUser"].([]string)) != expectedEntries {
		t.Logf("Unexpected Statement Principal has: `%d` entries but: `%d` expected", len(policy.Statements[1].Principal["CanonicalUser"].([]string)), expectedEntries)
		t.Fail()
	}

	if policy.Statements[2].Principal["Service"].([]string)[0] != expectedValueIndex0 {
		t.Logf("Unexpected Statement Principal: `%s` Statement Principal expected: `*`", policy.Statements[2].Principal["Service"].([]string)[0])
		t.Fail()
	}

	if policy.Statements[2].Principal["Service"].([]string)[1] != expectedValueIndex1 {
		t.Logf("Unexpected Statement Principal: `%s` Statement Principal expected: `*`", policy.Statements[2].Principal["Service"].([]string)[1])
		t.Fail()
	}

	if len(policy.Statements[2].Principal["Service"].([]string)) != expectedEntries {
		t.Logf("Unexpected Statement Principal has: `%d` entries but: `%d` expected", len(policy.Statements[2].Principal["Service"].([]string)), expectedEntries)
		t.Fail()
	}

	if policy.Statements[3].Principal["Any Value"].([]string)[0] != expectedValueIndex0 {
		t.Logf("Unexpected Statement Principal: `%s` Statement Principal expected: `*`", policy.Statements[3].Principal["Any Value"].([]string)[0])
		t.Fail()
	}

	if policy.Statements[3].Principal["Any Value"].([]string)[1] != expectedValueIndex1 {
		t.Logf("Unexpected Statement Principal: `%s` Statement Principal expected: `*`", policy.Statements[3].Principal["Any Value"].([]string)[1])
		t.Fail()
	}

	if len(policy.Statements[3].Principal["Any Value"].([]string)) != expectedEntries {
		t.Logf("Unexpected Statement Principal has: `%d` entries but: `%d` expected", len(policy.Statements[3].Principal["Any Value"].([]string)), expectedEntries)
		t.Fail()
	}
}

func testStatementPrincipalElementRemoveDuplicateValues(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "Principal": { "AWS": ["111199991111", "999911119999", "111199991111"] }
        },
        {
          "Principal": { "CanonicalUser": ["111199991111", "999911119999", "111199991111"] }
        },
        {
          "Principal": { "Service": ["111199991111", "999911119999", "111199991111"] }
        },
        {
          "Principal": { "Any Value": ["111199991111", "999911119999", "111199991111"] }
        }
      ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	expectedValueIndex0 := "111199991111"
	expectedValueIndex1 := "999911119999"
	if policy.Statements[0].Principal["AWS"].([]string)[0] != expectedValueIndex0 {
		t.Logf("Unexpected Statement Principal: `%s` Statement Principal expected: `*`", policy.Statements[0].Principal["AWS"].([]string)[0])
		t.Fail()
	}

	if policy.Statements[0].Principal["AWS"].([]string)[1] != expectedValueIndex1 {
		t.Logf("Unexpected Statement Principal: `%s` Statement Principal expected: `*`", policy.Statements[0].Principal["AWS"].([]string)[1])
		t.Fail()
	}

	expectedEntries := 2
	if len(policy.Statements[0].Principal["AWS"].([]string)) != expectedEntries {
		t.Logf("Unexpected Statement Principal has: `%d` entries but: `%d` expected", len(policy.Statements[1].Principal["AWS"].([]string)), expectedEntries)
		t.Fail()
	}

	if policy.Statements[1].Principal["CanonicalUser"].([]string)[0] != expectedValueIndex0 {
		t.Logf("Unexpected Statement Principal: `%s` Statement Principal expected: `*`", policy.Statements[1].Principal["CanonicalUser"].([]string)[0])
		t.Fail()
	}

	if policy.Statements[1].Principal["CanonicalUser"].([]string)[1] != expectedValueIndex1 {
		t.Logf("Unexpected Statement Principal: `%s` Statement Principal expected: `*`", policy.Statements[1].Principal["CanonicalUser"].([]string)[1])
		t.Fail()
	}

	if len(policy.Statements[1].Principal["CanonicalUser"].([]string)) != expectedEntries {
		t.Logf("Unexpected Statement Principal has: `%d` entries but: `%d` expected", len(policy.Statements[1].Principal["CanonicalUser"].([]string)), expectedEntries)
		t.Fail()
	}

	if policy.Statements[2].Principal["Service"].([]string)[0] != expectedValueIndex0 {
		t.Logf("Unexpected Statement Principal: `%s` Statement Principal expected: `*`", policy.Statements[2].Principal["Service"].([]string)[0])
		t.Fail()
	}

	if policy.Statements[2].Principal["Service"].([]string)[1] != expectedValueIndex1 {
		t.Logf("Unexpected Statement Principal: `%s` Statement Principal expected: `*`", policy.Statements[2].Principal["Service"].([]string)[1])
		t.Fail()
	}

	if len(policy.Statements[2].Principal["Service"].([]string)) != expectedEntries {
		t.Logf("Unexpected Statement Principal has: `%d` entries but: `%d` expected", len(policy.Statements[2].Principal["Service"].([]string)), expectedEntries)
		t.Fail()
	}

	if policy.Statements[3].Principal["Any Value"].([]string)[0] != expectedValueIndex0 {
		t.Logf("Unexpected Statement Principal: `%s` Statement Principal expected: `*`", policy.Statements[3].Principal["Any Value"].([]string)[0])
		t.Fail()
	}

	if policy.Statements[3].Principal["Any Value"].([]string)[1] != expectedValueIndex1 {
		t.Logf("Unexpected Statement Principal: `%s` Statement Principal expected: `*`", policy.Statements[3].Principal["Any Value"].([]string)[1])
		t.Fail()
	}

	if len(policy.Statements[3].Principal["Any Value"].([]string)) != expectedEntries {
		t.Logf("Unexpected Statement Principal has: `%d` entries but: `%d` expected", len(policy.Statements[3].Principal["Any Value"].([]string)), expectedEntries)
		t.Fail()
	}
}

func testStatementPrincipalElementMultipleValueAtPrincipalFails(t *testing.T) {
	// Set up
	policyContent := `{ "Statement": [{ "Principal": ["111199991111", "999911119999"] }] }`

	// Test
	_, err := canonicalPolicy(policyContent)

	// Expected
	if err == nil {
		t.Fatal("An error is expected when parsing no content")
	}

	expectedErrorMsg := "canonicalPolicy failed unmarshalling source data: UnmarshalJSON failed for Statements (Array of Statement): %5B%7B+%22Principal%22%3A+%5B%22111199991111%22%2C+%22999911119999%22%5D+%7D%5D.  src: %7B+%22Statement%22%3A+%5B%7B+%22Principal%22%3A+%5B%22111199991111%22%2C+%22999911119999%22%5D+%7D%5D+%7D"

	if errorMsg := err.Error(); errorMsg != expectedErrorMsg {
		t.Fatalf("The error message returned is expected to be: %s", expectedErrorMsg)
	}
}

func TestStatementNotPrincipalElement(t *testing.T) {
	t.Run("TestStatementNotPrincipalElementSingleValue", testStatementNotPrincipalElementSingleValue)
	t.Run("TestStatementNotPrincipalElementMultipleValue", testStatementNotPrincipalElementMultipleValue)
	t.Run("TestStatementNotPrincipalElementMultipleValueAtNotPrincipalFails", testStatementNotPrincipalElementMultipleValueAtNotPrincipalFails)
	t.Run("TestStatementNotPrincipalElementRemoveDeplicateValues", testStatementNotPrincipalElementRemoveDuplicateValues)
	t.Run("TestStatementNotPrincipalElementSortValues", testStatementNotPrincipalElementSortValues)
	t.Run("TestStatementNotPrincipalElementValuesOtherThanStringsFails", testStatementNotPrincipalValuesOtherThanStringsFails)
}

func testStatementNotPrincipalValuesOtherThanStringsFails(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "NotPrincipal": 10
        }
      ]
    }
	`

	// Test
	_, err := canonicalPolicy(policyContent)

	// Expected
	if err == nil {
		t.Fatal("An error is expected when parsing no content")
	}

	expectedErrorMsg := "canonicalPolicy failed unmarshalling source data: UnmarshalJSON failed for Statements (Array of Statement): %5B%0A++++++++%7B%0A++++++++++%22NotPrincipal%22%3A+10%0A++++++++%7D%0A++++++%5D.  src: %0A++++%7B%0A++++++%22Statement%22%3A+%5B%0A++++++++%7B%0A++++++++++%22NotPrincipal%22%3A+10%0A++++++++%7D%0A++++++%5D%0A++++%7D%0A%09"

	if errorMsg := err.Error(); errorMsg != expectedErrorMsg {
		t.Fatalf("The error message returned is expected to be: %s", expectedErrorMsg)
	}
}

func testStatementNotPrincipalElementSingleValue(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "NotPrincipal": "*"
        },
        {
          "NotPrincipal": { "AWS": "*" }
        },
        {
          "NotPrincipal": { "CanonicalUser": "*" }
        },
        {
          "NotPrincipal": { "Service": "*" }
        },
        {
          "NotPrincipal": { "Any Value": "*" }
        }
      ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	expectedValue := "*"
	if policy.Statements[0].NotPrincipal["AWS"].([]string)[0] != expectedValue {
		t.Logf("Unexpected Statement NotPrincipal: `%s` Statement NotPrincipal expected: `*`", policy.Statements[0].NotPrincipal["AWS"].([]string)[0])
		t.Fail()
	}

	expectedEntries := 1
	if len(policy.Statements[0].NotPrincipal) != expectedEntries {
		t.Logf("Unexpected Statement NotPrincipal has: `%d` entries but: `%d` expected", len(policy.Statements[0].NotPrincipal), expectedEntries)
		t.Fail()
	}

	if policy.Statements[1].NotPrincipal["AWS"].([]string)[0] != expectedValue {
		t.Logf("Unexpected Statement NotPrincipal: `%s` Statement NotPrincipal expected: `*`", policy.Statements[1].NotPrincipal["AWS"].([]string)[0])
		t.Fail()
	}

	if len(policy.Statements[1].NotPrincipal) != expectedEntries {
		t.Logf("Unexpected Statement NotPrincipal has: `%d` entries but: `%d` expected", len(policy.Statements[1].NotPrincipal), expectedEntries)
		t.Fail()
	}

	if policy.Statements[2].NotPrincipal["CanonicalUser"].([]string)[0] != expectedValue {
		t.Logf("Unexpected Statement NotPrincipal: `%s` Statement NotPrincipal expected: `*`", policy.Statements[2].NotPrincipal["AWS"].([]string)[0])
		t.Fail()
	}

	if len(policy.Statements[2].NotPrincipal) != expectedEntries {
		t.Logf("Unexpected Statement NotPrincipal has: `%d` entries but: `%d` expected", len(policy.Statements[2].NotPrincipal), expectedEntries)
		t.Fail()
	}

	if policy.Statements[3].NotPrincipal["Service"].([]string)[0] != expectedValue {
		t.Logf("Unexpected Statement NotPrincipal: `%s` Statement NotPrincipal expected: `*`", policy.Statements[3].NotPrincipal["AWS"].([]string)[0])
		t.Fail()
	}

	if len(policy.Statements[3].NotPrincipal) != expectedEntries {
		t.Logf("Unexpected Statement NotPrincipal has: `%d` entries but: `%d` expected", len(policy.Statements[3].NotPrincipal), expectedEntries)
		t.Fail()
	}

	if policy.Statements[4].NotPrincipal["Any Value"].([]string)[0] != expectedValue {
		t.Logf("Unexpected Statement NotPrincipal: `%s` Statement NotPrincipal expected: `*`", policy.Statements[4].NotPrincipal["AWS"].([]string)[0])
		t.Fail()
	}

	if len(policy.Statements[4].NotPrincipal) != expectedEntries {
		t.Logf("Unexpected Statement NotPrincipal has: `%d` entries but: `%d` expected", len(policy.Statements[4].NotPrincipal), expectedEntries)
		t.Fail()
	}
}

func testStatementNotPrincipalElementMultipleValue(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "NotPrincipal": { "AWS": ["111199991111", "999911119999"] }
        },
        {
          "NotPrincipal": { "CanonicalUser": ["111199991111", "999911119999"] }
        },
        {
          "NotPrincipal": { "Service": ["111199991111", "999911119999"] }
        },
        {
          "NotPrincipal": { "Any Value": ["111199991111", "999911119999"] }
        }
      ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	expectedValueIndex0 := "111199991111"
	expectedValueIndex1 := "999911119999"
	if policy.Statements[0].NotPrincipal["AWS"].([]string)[0] != expectedValueIndex0 {
		t.Logf("Unexpected Statement NotPrincipal: `%s` Statement NotPrincipal expected: `*`", policy.Statements[0].NotPrincipal["AWS"].([]string)[0])
		t.Fail()
	}

	if policy.Statements[0].NotPrincipal["AWS"].([]string)[1] != expectedValueIndex1 {
		t.Logf("Unexpected Statement NotPrincipal: `%s` Statement NotPrincipal expected: `*`", policy.Statements[0].NotPrincipal["AWS"].([]string)[1])
		t.Fail()
	}

	expectedEntries := 2
	if len(policy.Statements[0].NotPrincipal["AWS"].([]string)) != expectedEntries {
		t.Logf("Unexpected Statement NotPrincipal has: `%d` entries but: `%d` expected", len(policy.Statements[1].NotPrincipal["AWS"].([]string)), expectedEntries)
		t.Fail()
	}

	if policy.Statements[1].NotPrincipal["CanonicalUser"].([]string)[0] != expectedValueIndex0 {
		t.Logf("Unexpected Statement NotPrincipal: `%s` Statement NotPrincipal expected: `*`", policy.Statements[1].NotPrincipal["CanonicalUser"].([]string)[0])
		t.Fail()
	}

	if policy.Statements[1].NotPrincipal["CanonicalUser"].([]string)[1] != expectedValueIndex1 {
		t.Logf("Unexpected Statement NotPrincipal: `%s` Statement NotPrincipal expected: `*`", policy.Statements[1].NotPrincipal["CanonicalUser"].([]string)[1])
		t.Fail()
	}

	if len(policy.Statements[1].NotPrincipal["CanonicalUser"].([]string)) != expectedEntries {
		t.Logf("Unexpected Statement NotPrincipal has: `%d` entries but: `%d` expected", len(policy.Statements[1].NotPrincipal["CanonicalUser"].([]string)), expectedEntries)
		t.Fail()
	}

	if policy.Statements[2].NotPrincipal["Service"].([]string)[0] != expectedValueIndex0 {
		t.Logf("Unexpected Statement NotPrincipal: `%s` Statement NotPrincipal expected: `*`", policy.Statements[2].NotPrincipal["Service"].([]string)[0])
		t.Fail()
	}

	if policy.Statements[2].NotPrincipal["Service"].([]string)[1] != expectedValueIndex1 {
		t.Logf("Unexpected Statement NotPrincipal: `%s` Statement NotPrincipal expected: `*`", policy.Statements[2].NotPrincipal["Service"].([]string)[1])
		t.Fail()
	}

	if len(policy.Statements[2].NotPrincipal["Service"].([]string)) != expectedEntries {
		t.Logf("Unexpected Statement NotPrincipal has: `%d` entries but: `%d` expected", len(policy.Statements[2].NotPrincipal["Service"].([]string)), expectedEntries)
		t.Fail()
	}

	if policy.Statements[3].NotPrincipal["Any Value"].([]string)[0] != expectedValueIndex0 {
		t.Logf("Unexpected Statement NotPrincipal: `%s` Statement NotPrincipal expected: `*`", policy.Statements[3].NotPrincipal["Any Value"].([]string)[0])
		t.Fail()
	}

	if policy.Statements[3].NotPrincipal["Any Value"].([]string)[1] != expectedValueIndex1 {
		t.Logf("Unexpected Statement NotPrincipal: `%s` Statement NotPrincipal expected: `*`", policy.Statements[3].NotPrincipal["Any Value"].([]string)[1])
		t.Fail()
	}

	if len(policy.Statements[3].NotPrincipal["Any Value"].([]string)) != expectedEntries {
		t.Logf("Unexpected Statement NotPrincipal has: `%d` entries but: `%d` expected", len(policy.Statements[3].NotPrincipal["Any Value"].([]string)), expectedEntries)
		t.Fail()
	}
}

func testStatementNotPrincipalElementSortValues(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "NotPrincipal": { "AWS": ["999911119999", "111199991111"] }
        },
        {
          "NotPrincipal": { "CanonicalUser": ["999911119999", "111199991111"] }
        },
        {
          "NotPrincipal": { "Service": ["999911119999", "111199991111"] }
        },
        {
          "NotPrincipal": { "Any Value": ["999911119999", "111199991111"] }
        }
      ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	expectedValueIndex0 := "111199991111"
	expectedValueIndex1 := "999911119999"
	if policy.Statements[0].NotPrincipal["AWS"].([]string)[0] != expectedValueIndex0 {
		t.Logf("Unexpected Statement NotPrincipal: `%s` Statement NotPrincipal expected: `*`", policy.Statements[0].NotPrincipal["AWS"].([]string)[0])
		t.Fail()
	}

	if policy.Statements[0].NotPrincipal["AWS"].([]string)[1] != expectedValueIndex1 {
		t.Logf("Unexpected Statement NotPrincipal: `%s` Statement NotPrincipal expected: `*`", policy.Statements[0].NotPrincipal["AWS"].([]string)[1])
		t.Fail()
	}

	expectedEntries := 2
	if len(policy.Statements[0].NotPrincipal["AWS"].([]string)) != expectedEntries {
		t.Logf("Unexpected Statement NotPrincipal has: `%d` entries but: `%d` expected", len(policy.Statements[1].NotPrincipal["AWS"].([]string)), expectedEntries)
		t.Fail()
	}

	if policy.Statements[1].NotPrincipal["CanonicalUser"].([]string)[0] != expectedValueIndex0 {
		t.Logf("Unexpected Statement NotPrincipal: `%s` Statement NotPrincipal expected: `*`", policy.Statements[1].NotPrincipal["CanonicalUser"].([]string)[0])
		t.Fail()
	}

	if policy.Statements[1].NotPrincipal["CanonicalUser"].([]string)[1] != expectedValueIndex1 {
		t.Logf("Unexpected Statement NotPrincipal: `%s` Statement NotPrincipal expected: `*`", policy.Statements[1].NotPrincipal["CanonicalUser"].([]string)[1])
		t.Fail()
	}

	if len(policy.Statements[1].NotPrincipal["CanonicalUser"].([]string)) != expectedEntries {
		t.Logf("Unexpected Statement NotPrincipal has: `%d` entries but: `%d` expected", len(policy.Statements[1].NotPrincipal["CanonicalUser"].([]string)), expectedEntries)
		t.Fail()
	}

	if policy.Statements[2].NotPrincipal["Service"].([]string)[0] != expectedValueIndex0 {
		t.Logf("Unexpected Statement NotPrincipal: `%s` Statement NotPrincipal expected: `*`", policy.Statements[2].NotPrincipal["Service"].([]string)[0])
		t.Fail()
	}

	if policy.Statements[2].NotPrincipal["Service"].([]string)[1] != expectedValueIndex1 {
		t.Logf("Unexpected Statement NotPrincipal: `%s` Statement NotPrincipal expected: `*`", policy.Statements[2].NotPrincipal["Service"].([]string)[1])
		t.Fail()
	}

	if len(policy.Statements[2].NotPrincipal["Service"].([]string)) != expectedEntries {
		t.Logf("Unexpected Statement NotPrincipal has: `%d` entries but: `%d` expected", len(policy.Statements[2].NotPrincipal["Service"].([]string)), expectedEntries)
		t.Fail()
	}

	if policy.Statements[3].NotPrincipal["Any Value"].([]string)[0] != expectedValueIndex0 {
		t.Logf("Unexpected Statement NotPrincipal: `%s` Statement NotPrincipal expected: `*`", policy.Statements[3].NotPrincipal["Any Value"].([]string)[0])
		t.Fail()
	}

	if policy.Statements[3].NotPrincipal["Any Value"].([]string)[1] != expectedValueIndex1 {
		t.Logf("Unexpected Statement NotPrincipal: `%s` Statement NotPrincipal expected: `*`", policy.Statements[3].NotPrincipal["Any Value"].([]string)[1])
		t.Fail()
	}

	if len(policy.Statements[3].NotPrincipal["Any Value"].([]string)) != expectedEntries {
		t.Logf("Unexpected Statement NotPrincipal has: `%d` entries but: `%d` expected", len(policy.Statements[3].NotPrincipal["Any Value"].([]string)), expectedEntries)
		t.Fail()
	}
}

func testStatementNotPrincipalElementRemoveDuplicateValues(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "NotPrincipal": { "AWS": ["111199991111", "999911119999", "111199991111"] }
        },
        {
          "NotPrincipal": { "CanonicalUser": ["111199991111", "999911119999", "111199991111"] }
        },
        {
          "NotPrincipal": { "Service": ["111199991111", "999911119999", "111199991111"] }
        },
        {
          "NotPrincipal": { "Any Value": ["111199991111", "999911119999", "111199991111"] }
        }
      ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	expectedValueIndex0 := "111199991111"
	expectedValueIndex1 := "999911119999"
	if policy.Statements[0].NotPrincipal["AWS"].([]string)[0] != expectedValueIndex0 {
		t.Logf("Unexpected Statement NotPrincipal: `%s` Statement NotPrincipal expected: `*`", policy.Statements[0].NotPrincipal["AWS"].([]string)[0])
		t.Fail()
	}

	if policy.Statements[0].NotPrincipal["AWS"].([]string)[1] != expectedValueIndex1 {
		t.Logf("Unexpected Statement NotPrincipal: `%s` Statement NotPrincipal expected: `*`", policy.Statements[0].NotPrincipal["AWS"].([]string)[1])
		t.Fail()
	}

	expectedEntries := 2
	if len(policy.Statements[0].NotPrincipal["AWS"].([]string)) != expectedEntries {
		t.Logf("Unexpected Statement NotPrincipal has: `%d` entries but: `%d` expected", len(policy.Statements[1].NotPrincipal["AWS"].([]string)), expectedEntries)
		t.Fail()
	}

	if policy.Statements[1].NotPrincipal["CanonicalUser"].([]string)[0] != expectedValueIndex0 {
		t.Logf("Unexpected Statement NotPrincipal: `%s` Statement NotPrincipal expected: `*`", policy.Statements[1].NotPrincipal["CanonicalUser"].([]string)[0])
		t.Fail()
	}

	if policy.Statements[1].NotPrincipal["CanonicalUser"].([]string)[1] != expectedValueIndex1 {
		t.Logf("Unexpected Statement NotPrincipal: `%s` Statement NotPrincipal expected: `*`", policy.Statements[1].NotPrincipal["CanonicalUser"].([]string)[1])
		t.Fail()
	}

	if len(policy.Statements[1].NotPrincipal["CanonicalUser"].([]string)) != expectedEntries {
		t.Logf("Unexpected Statement NotPrincipal has: `%d` entries but: `%d` expected", len(policy.Statements[1].NotPrincipal["CanonicalUser"].([]string)), expectedEntries)
		t.Fail()
	}

	if policy.Statements[2].NotPrincipal["Service"].([]string)[0] != expectedValueIndex0 {
		t.Logf("Unexpected Statement NotPrincipal: `%s` Statement NotPrincipal expected: `*`", policy.Statements[2].NotPrincipal["Service"].([]string)[0])
		t.Fail()
	}

	if policy.Statements[2].NotPrincipal["Service"].([]string)[1] != expectedValueIndex1 {
		t.Logf("Unexpected Statement NotPrincipal: `%s` Statement NotPrincipal expected: `*`", policy.Statements[2].NotPrincipal["Service"].([]string)[1])
		t.Fail()
	}

	if len(policy.Statements[2].NotPrincipal["Service"].([]string)) != expectedEntries {
		t.Logf("Unexpected Statement NotPrincipal has: `%d` entries but: `%d` expected", len(policy.Statements[2].NotPrincipal["Service"].([]string)), expectedEntries)
		t.Fail()
	}

	if policy.Statements[3].NotPrincipal["Any Value"].([]string)[0] != expectedValueIndex0 {
		t.Logf("Unexpected Statement NotPrincipal: `%s` Statement NotPrincipal expected: `*`", policy.Statements[3].NotPrincipal["Any Value"].([]string)[0])
		t.Fail()
	}

	if policy.Statements[3].NotPrincipal["Any Value"].([]string)[1] != expectedValueIndex1 {
		t.Logf("Unexpected Statement NotPrincipal: `%s` Statement NotPrincipal expected: `*`", policy.Statements[3].NotPrincipal["Any Value"].([]string)[1])
		t.Fail()
	}

	if len(policy.Statements[3].NotPrincipal["Any Value"].([]string)) != expectedEntries {
		t.Logf("Unexpected Statement NotPrincipal has: `%d` entries but: `%d` expected", len(policy.Statements[3].NotPrincipal["Any Value"].([]string)), expectedEntries)
		t.Fail()
	}
}

func testStatementNotPrincipalElementMultipleValueAtNotPrincipalFails(t *testing.T) {
	// Set up
	policyContent := `{ "Statement": [{ "NotPrincipal": ["111199991111", "999911119999"] }] }`

	// Test
	_, err := canonicalPolicy(policyContent)

	// Expected
	if err == nil {
		t.Fatal("An error is expected when parsing no content")
	}

	expectedErrorMsg := "canonicalPolicy failed unmarshalling source data: UnmarshalJSON failed for Statements (Array of Statement): %5B%7B+%22NotPrincipal%22%3A+%5B%22111199991111%22%2C+%22999911119999%22%5D+%7D%5D.  src: %7B+%22Statement%22%3A+%5B%7B+%22NotPrincipal%22%3A+%5B%22111199991111%22%2C+%22999911119999%22%5D+%7D%5D+%7D"

	if errorMsg := err.Error(); errorMsg != expectedErrorMsg {
		t.Fatalf("The error message returned is expected to be: %s", expectedErrorMsg)
	}
}

func TestStatementActionElement(t *testing.T) {
	t.Run("TestActionValuesAreLowercased", testActionValuesAreLowercased)
	t.Run("TestDuplicateValuesInActionAreRemoved", testDuplicateValuesInActionAreRemoved)
	t.Run("TestActionSortsValues", testActionSortsValues)
	t.Run("TestActionWhenValueIsNotAString", testActionWhenValueIsNotAString)
}

func testActionWhenValueIsNotAString(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "Action": 3
        }
     ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	expectedAction := "3"
	if policy.Statements[0].Action[0] != expectedAction {
		t.Logf("Unexpected Statement Action: `%s` Statement Action expected: `%s`", policy.Statements[0].Action[0], expectedAction)
		t.Fail()
	}
}

func testActionSortsValues(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "Action": ["S3:List*", "S3:Get*"]
        },
        {
          "Action": ["S3:List*", "S3:Get*", "EC2:Launch*"]
        }
     ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	expectedAction := "s3:get*"
	if policy.Statements[0].Action[0] != expectedAction {
		t.Logf("Unexpected Statement Action: `%s` Statement Action expected: `%s`", policy.Statements[0].Action[0], expectedAction)
		t.Fail()
	}

	expectedAction = "s3:list*"
	if policy.Statements[0].Action[1] != expectedAction {
		t.Logf("Unexpected Statement Action: `%s` Statement Action expected: `%s`", policy.Statements[0].Action[1], expectedAction)
		t.Fail()
	}

	expectedAction = "ec2:launch*"
	if policy.Statements[1].Action[0] != expectedAction {
		t.Logf("Unexpected Statement Action: `%s` Statement Action expected: `%s`", policy.Statements[1].Action[0], expectedAction)
		t.Fail()
	}

	expectedAction = "s3:get*"
	if policy.Statements[1].Action[1] != expectedAction {
		t.Logf("Unexpected Statement Action: `%s` Statement Action expected: `%s`", policy.Statements[1].Action[1], expectedAction)
		t.Fail()
	}

	expectedAction = "s3:list*"
	if policy.Statements[1].Action[2] != expectedAction {
		t.Logf("Unexpected Statement Action: `%s` Statement Action expected: `%s`", policy.Statements[1].Action[2], expectedAction)
		t.Fail()
	}
}

func testDuplicateValuesInActionAreRemoved(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "Action": ["S3:Get*", "S3:Get*"]
        },
        {
          "Action": ["S3:Get*", "S3:List*", "S3:Get*", "S3:List*"]
        }
     ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	expectedAction := "s3:get*"
	if policy.Statements[0].Action[0] != expectedAction {
		t.Logf("Unexpected Statement Action: `%s` Statement Action expected: `%s`", policy.Statements[0].Action[0], expectedAction)
		t.Fail()
	}

	if len(policy.Statements[0].Action) != 1 {
		t.Logf("Unexpected Statement Action - too many actions returned")
		t.Fail()
	}

	if policy.Statements[1].Action[0] != expectedAction {
		t.Logf("Unexpected Statement Action: `%s` Statement Action expected: `%s`", policy.Statements[1].Action[0], expectedAction)
		t.Fail()
	}

	expectedAction = "s3:list*"
	if policy.Statements[1].Action[1] != expectedAction {
		t.Logf("Unexpected Statement Action: `%s` Statement Action expected: `%s`", policy.Statements[1].Action[1], expectedAction)
		t.Fail()
	}

	if len(policy.Statements[1].Action) != 2 {
		t.Logf("Unexpected Statement Action - too many actions returned")
		t.Fail()
	}
}

func testActionValuesAreLowercased(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "Action": "S3:Get*"
        },
        {
          "Action": ["S3:Get*"]
        },
        {
          "Action": ["S3:Get*", "S3:List*"]
        }
     ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	expectedAction := "s3:get*"
	if policy.Statements[0].Action[0] != expectedAction {
		t.Logf("Unexpected Statement Action: `%s` Statement Action expected: `%s`", policy.Statements[0].Action[0], expectedAction)
		t.Fail()
	}

	if policy.Statements[1].Action[0] != expectedAction {
		t.Logf("Unexpected Statement Action: `%s` Statement Action expected: `%s`", policy.Statements[1].Action[0], expectedAction)
		t.Fail()
	}

	if policy.Statements[2].Action[0] != expectedAction {
		t.Logf("Unexpected Statement Action: `%s` Statement Action expected: `%s`", policy.Statements[2].Action[0], expectedAction)
		t.Fail()
	}

	expectedAction = "s3:list*"
	if policy.Statements[2].Action[1] != expectedAction {
		t.Logf("Unexpected Statement Action: `%s` Statement Action expected: `%s`", policy.Statements[2].Action[1], expectedAction)
		t.Fail()
	}
}

func TestStatementNotActionElement(t *testing.T) {
	t.Run("TestNotActionValuesAreLowercased", testNotActionValuesAreLowercased)
	t.Run("TestDuplicateValuesInNotActionAreRemoved", testDuplicateValuesInNotActionAreRemoved)
	t.Run("TestNotActionSortsValues", testNotActionSortsValues)
	t.Run("TestNotActionWhenValueIsNotAString", testNotActionWhenValueIsNotAString)
}

func testNotActionValuesAreLowercased(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "NotAction": "S3:Get*"
        },
        {
          "NotAction": ["S3:Get*"]
        },
        {
          "NotAction": ["S3:Get*", "S3:List*"]
        }
     ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	expectedNotAction := "s3:get*"
	if policy.Statements[0].NotAction[0] != expectedNotAction {
		t.Logf("Unexpected Statement NotAction: `%s` Statement NotAction expected: `%s`", policy.Statements[0].NotAction[0], expectedNotAction)
		t.Fail()
	}

	if policy.Statements[1].NotAction[0] != expectedNotAction {
		t.Logf("Unexpected Statement NotAction: `%s` Statement NotAction expected: `%s`", policy.Statements[1].NotAction[0], expectedNotAction)
		t.Fail()
	}

	if policy.Statements[2].NotAction[0] != expectedNotAction {
		t.Logf("Unexpected Statement NotAction: `%s` Statement NotAction expected: `%s`", policy.Statements[2].NotAction[0], expectedNotAction)
		t.Fail()
	}

	expectedNotAction = "s3:list*"
	if policy.Statements[2].NotAction[1] != expectedNotAction {
		t.Logf("Unexpected Statement NotAction: `%s` Statement NotAction expected: `%s`", policy.Statements[2].NotAction[1], expectedNotAction)
		t.Fail()
	}
}

func testDuplicateValuesInNotActionAreRemoved(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "NotAction": ["S3:Get*", "S3:Get*"]
        },
        {
          "NotAction": ["S3:Get*", "S3:List*", "S3:Get*", "S3:List*"]
        }
     ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	expectedNotAction := "s3:get*"
	if policy.Statements[0].NotAction[0] != expectedNotAction {
		t.Logf("Unexpected Statement NotAction: `%s` Statement NotAction expected: `%s`", policy.Statements[0].NotAction[0], expectedNotAction)
		t.Fail()
	}

	if len(policy.Statements[0].NotAction) != 1 {
		t.Logf("Unexpected Statement NotAction - too many actions returned")
		t.Fail()
	}

	if policy.Statements[1].NotAction[0] != expectedNotAction {
		t.Logf("Unexpected Statement NotAction: `%s` Statement NotAction expected: `%s`", policy.Statements[1].NotAction[0], expectedNotAction)
		t.Fail()
	}

	expectedNotAction = "s3:list*"
	if policy.Statements[1].NotAction[1] != expectedNotAction {
		t.Logf("Unexpected Statement NotAction: `%s` Statement NotAction expected: `%s`", policy.Statements[1].NotAction[1], expectedNotAction)
		t.Fail()
	}

	if len(policy.Statements[1].NotAction) != 2 {
		t.Logf("Unexpected Statement NotAction - too many actions returned")
		t.Fail()
	}
}

func testNotActionSortsValues(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "NotAction": ["S3:List*", "S3:Get*"]
        },
        {
          "NotAction": ["S3:List*", "S3:Get*", "EC2:Launch*"]
        }
     ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	expectedNotAction := "s3:get*"
	if policy.Statements[0].NotAction[0] != expectedNotAction {
		t.Logf("Unexpected Statement NotAction: `%s` Statement NotAction expected: `%s`", policy.Statements[0].NotAction[0], expectedNotAction)
		t.Fail()
	}

	expectedNotAction = "s3:list*"
	if policy.Statements[0].NotAction[1] != expectedNotAction {
		t.Logf("Unexpected Statement NotAction: `%s` Statement NotAction expected: `%s`", policy.Statements[0].NotAction[1], expectedNotAction)
		t.Fail()
	}

	expectedNotAction = "ec2:launch*"
	if policy.Statements[1].NotAction[0] != expectedNotAction {
		t.Logf("Unexpected Statement NotAction: `%s` Statement NotAction expected: `%s`", policy.Statements[1].NotAction[0], expectedNotAction)
		t.Fail()
	}

	expectedNotAction = "s3:get*"
	if policy.Statements[1].NotAction[1] != expectedNotAction {
		t.Logf("Unexpected Statement NotAction: `%s` Statement NotAction expected: `%s`", policy.Statements[1].NotAction[1], expectedNotAction)
		t.Fail()
	}

	expectedNotAction = "s3:list*"
	if policy.Statements[1].NotAction[2] != expectedNotAction {
		t.Logf("Unexpected Statement NotAction: `%s` Statement NotAction expected: `%s`", policy.Statements[1].NotAction[2], expectedNotAction)
		t.Fail()
	}
}

func testNotActionWhenValueIsNotAString(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "NotAction": 3
        }
     ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	expectedNotAction := "3"
	if policy.Statements[0].NotAction[0] != expectedNotAction {
		t.Logf("Unexpected Statement NotAction: `%s` Statement NotAction expected: `%s`", policy.Statements[0].NotAction[0], expectedNotAction)
		t.Fail()
	}
}

func TestStatementEffectElement(t *testing.T) {
	t.Run("TestStatementEffectElementAllStringsValid", testStatementEffectElementAllStringsValid)
	t.Run("TestStatementEffectElementForMultipleValuesFails", testStatementEffectElementForMultipleValuesFails)
	t.Run("TestStatementEffectElementValuesOtherThanStringsFails", testStatementEffectElementValuesOtherThanStringsFails)
}

func testStatementEffectElementValuesOtherThanStringsFails(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "Sid": 3
        }
     ]
    }
	`

	// Test
	_, err := canonicalPolicy(policyContent)

	// Expected
	if err == nil {
		t.Fatal("An error is expected when parsing no content")
	}

	expectedErrorMsg := "canonicalPolicy failed unmarshalling source data: UnmarshalJSON failed for Statements (Array of Statement): %5B%0A++++++++%7B%0A++++++++++%22Sid%22%3A+3%0A++++++++%7D%0A+++++%5D.  src: %0A++++%7B%0A++++++%22Statement%22%3A+%5B%0A++++++++%7B%0A++++++++++%22Sid%22%3A+3%0A++++++++%7D%0A+++++%5D%0A++++%7D%0A%09"

	if errorMsg := err.Error(); errorMsg != expectedErrorMsg {
		t.Fatalf("The error message returned is expected to be: %s", expectedErrorMsg)
	}
}

func testStatementEffectElementForMultipleValuesFails(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "Sid": ["Sid1", "Sid2"]
        }
      ]
    }
	`

	// Test
	_, err := canonicalPolicy(policyContent)

	// Expected
	if err == nil {
		t.Fatal("An error is expected when parsing no content")
	}

	expectedErrorMsg := "canonicalPolicy failed unmarshalling source data: UnmarshalJSON failed for Statements (Array of Statement): %5B%0A++++++++%7B%0A++++++++++%22Sid%22%3A+%5B%22Sid1%22%2C+%22Sid2%22%5D%0A++++++++%7D%0A++++++%5D.  src: %0A++++%7B%0A++++++%22Statement%22%3A+%5B%0A++++++++%7B%0A++++++++++%22Sid%22%3A+%5B%22Sid1%22%2C+%22Sid2%22%5D%0A++++++++%7D%0A++++++%5D%0A++++%7D%0A%09"

	if errorMsg := err.Error(); errorMsg != expectedErrorMsg {
		t.Fatalf("The error message returned is expected to be: %s", expectedErrorMsg)
	}
}

func testStatementEffectElementAllStringsValid(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "Effect": "Allow"
        },
        {
          "Effect": "Deny"
        },
        {
          "Effect": "CanBeAnyValue"
        }
      ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	var expectedStatement1Effect = "Allow"
	if policy.Statements[0].Effect != expectedStatement1Effect {
		t.Logf("Unexpected Statement Effect: `%s` Statement Effect expected: `%s`", policy.Statements[0].Effect, expectedStatement1Effect)
		t.Fail()
	}

	var expectedStatement2Effect = "Deny"
	if policy.Statements[1].Effect != expectedStatement2Effect {
		t.Logf("Unexpected Statement Effect: `%s` Statement Effect expected: `%s`", policy.Statements[1].Effect, expectedStatement2Effect)
		t.Fail()
	}

	var expectedStatement3Effect = "CanBeAnyValue"
	if policy.Statements[2].Effect != expectedStatement3Effect {
		t.Logf("Unexpected Statement Effect: `%s` Statement Effect expected: `%s`", policy.Statements[2].Effect, expectedStatement3Effect)
		t.Fail()
	}
}

func TestStatementSidElement(t *testing.T) {
	t.Run("TestStatementSidElementAllStringsValid", testStatementSidElementAllStringsValid)
	t.Run("TestStatementSidElementForMultipleValuesFails", testStatementSidElementForMultipleValuesFails)
	t.Run("TestStatementSidElementValuesOtherThanStringsFails", testStatementSidElementValuesOtherThanStringsFails)
}

func testStatementSidElementValuesOtherThanStringsFails(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "Sid": 3
        }
     ]
    }
	`

	// Test
	_, err := canonicalPolicy(policyContent)

	// Expected
	if err == nil {
		t.Fatal("An error is expected when parsing no content")
	}

	expectedErrorMsg := "canonicalPolicy failed unmarshalling source data: UnmarshalJSON failed for Statements (Array of Statement): %5B%0A++++++++%7B%0A++++++++++%22Sid%22%3A+3%0A++++++++%7D%0A+++++%5D.  src: %0A++++%7B%0A++++++%22Statement%22%3A+%5B%0A++++++++%7B%0A++++++++++%22Sid%22%3A+3%0A++++++++%7D%0A+++++%5D%0A++++%7D%0A%09"

	if errorMsg := err.Error(); errorMsg != expectedErrorMsg {
		t.Fatalf("The error message returned is expected to be: %s", expectedErrorMsg)
	}
}

func testStatementSidElementAllStringsValid(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "Sid": "Sid1"
        },
        {
          "Sid": "Sid2"
        }
      ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	var expectedStatement1Sid = "Sid1"
	if policy.Statements[0].Sid != expectedStatement1Sid {
		t.Logf("Unexpected Statement Sid: `%s` Statement Sid expected: `%s`", policy.Statements[0].Sid, expectedStatement1Sid)
		t.Fail()
	}

	var expectedStatement2Sid = "Sid2"
	if policy.Statements[1].Sid != expectedStatement2Sid {
		t.Logf("Unexpected Statement Sid: `%s` Statement Sid expected: `%s`", policy.Statements[1].Sid, expectedStatement2Sid)
		t.Fail()
	}
}

func testStatementSidElementForMultipleValuesFails(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "Sid": ["Sid1", "Sid2"]
        }
      ]
    }
	`

	// Test
	_, err := canonicalPolicy(policyContent)

	// Expected
	if err == nil {
		t.Fatal("An error is expected when parsing no content")
	}

	expectedErrorMsg := "canonicalPolicy failed unmarshalling source data: UnmarshalJSON failed for Statements (Array of Statement): %5B%0A++++++++%7B%0A++++++++++%22Sid%22%3A+%5B%22Sid1%22%2C+%22Sid2%22%5D%0A++++++++%7D%0A++++++%5D.  src: %0A++++%7B%0A++++++%22Statement%22%3A+%5B%0A++++++++%7B%0A++++++++++%22Sid%22%3A+%5B%22Sid1%22%2C+%22Sid2%22%5D%0A++++++++%7D%0A++++++%5D%0A++++%7D%0A%09"

	if errorMsg := err.Error(); errorMsg != expectedErrorMsg {
		t.Fatalf("The error message returned is expected to be: %s", expectedErrorMsg)
	}
}

func TestStatementResourceElement(t *testing.T) {
	t.Run("TestStatementResourceElementLeavesCasingUnchanged", testStatementResourceElementLeavesCasingUnchanged)
	t.Run("TestStatementResourceElementSorts", testStatementResourceElementSorts)
	t.Run("TestActionWhenValueIsNotAString", testResourceWhenValueIsNotAString)
}

func testResourceWhenValueIsNotAString(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "Resource": 3
        }
     ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	expectedValue := "3"
	if policy.Statements[0].Resource[0] != expectedValue {
		t.Logf("Unexpected Statement Resource: `%s` Statement Resource expected: `%s`", policy.Statements[0].Resource[0], expectedValue)
		t.Fail()
	}
}

func testStatementResourceElementSorts(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "Resource": ["ccc/*", "bbb/*", "aaa/*"]
        },
        {
          "Resource": ["arn:aws:dynamodb:us-east-2:account-ID-without-hyphens:table/*", "arn:aws:dynamodb:us-east-1:account-ID-without-hyphens:table/*"]
        }
      ]
    }
    `

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	expectedValue := "aaa/*"
	if policy.Statements[0].Resource[0] != expectedValue {
		t.Logf("Unexpected Statement Effect: `%s` Statement Resource expected: `%s`", policy.Statements[0].Resource[0], expectedValue)
		t.Fail()
	}

	expectedValue = "bbb/*"
	if policy.Statements[0].Resource[1] != expectedValue {
		t.Logf("Unexpected Statement Resource: `%s` Statement Resource expected: `%s`", policy.Statements[0].Resource[1], expectedValue)
		t.Fail()
	}

	expectedValue = "ccc/*"
	if policy.Statements[0].Resource[2] != expectedValue {
		t.Logf("Unexpected Statement Resource: `%s` Statement Resource expected: `%s`", policy.Statements[0].Resource[2], expectedValue)
		t.Fail()
	}

	expectedValue = "arn:aws:dynamodb:us-east-1:account-ID-without-hyphens:table/*"
	if policy.Statements[1].Resource[0] != expectedValue {
		t.Logf("Unexpected Statement Resource: `%s` Statement Resource expected: `%s`", policy.Statements[1].Resource[0], expectedValue)
		t.Fail()
	}

	expectedValue = "arn:aws:dynamodb:us-east-2:account-ID-without-hyphens:table/*"
	if policy.Statements[1].Resource[1] != expectedValue {
		t.Logf("Unexpected Statement Resource: `%s` Statement Resource expected: `%s`", policy.Statements[1].Resource[1], expectedValue)
		t.Fail()
	}
}

func testStatementResourceElementLeavesCasingUnchanged(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "Resource": "arn:aws:dynamodb:us-east-1:account-ID-without-hyphens:table/*"
        },
        {
          "Resource": ["arn:aws:dynamodb:us-east-1:account-ID-without-hyphens:table/*"]
        },
        {
          "Resource": ["arn:aws:dynamodb:us-east-1:account-ID-without-hyphens:table/*", "arn:aws:dynamodb:us-east-2:account-ID-without-hyphens:table/*"]
        }
      ]
    }
    `

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	expectedValue := "arn:aws:dynamodb:us-east-1:account-ID-without-hyphens:table/*"
	if policy.Statements[0].Resource[0] != expectedValue {
		t.Logf("Unexpected Statement Resource: `%s` Statement Resource expected: `%s`", policy.Statements[0].Resource[0], expectedValue)
		t.Fail()
	}

	if policy.Statements[1].Resource[0] != expectedValue {
		t.Logf("Unexpected Statement Resource: `%s` Statement Resource expected: `%s`", policy.Statements[1].Resource[0], expectedValue)
		t.Fail()
	}

	if policy.Statements[2].Resource[0] != expectedValue {
		t.Logf("Unexpected Statement Resource: `%s` Statement Resource expected: `%s`", policy.Statements[2].Resource[0], expectedValue)
		t.Fail()
	}

	expectedValue = "arn:aws:dynamodb:us-east-2:account-ID-without-hyphens:table/*"
	if policy.Statements[2].Resource[1] != expectedValue {
		t.Logf("Unexpected Statement Resource: `%s` Statement Resource expected: `%s`", policy.Statements[2].Resource[1], expectedValue)
		t.Fail()
	}
}

func TestStatementNotResourceElement(t *testing.T) {
	t.Run("TestStatementNotResourceElementLeavesCasingUnchanged", testStatementNotResourceElementLeavesCasingUnchanged)
	t.Run("TestStatementNotResourceElementSorts", testStatementNotResourceElementSorts)
	t.Run("TestNotResourceWhenValueIsNotAString", testNotResourceWhenValueIsNotAString)
}

func testNotResourceWhenValueIsNotAString(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "NotResource": 3
        }
     ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	expectedValue := "3"
	if policy.Statements[0].NotResource[0] != expectedValue {
		t.Logf("Unexpected Statement NotResource: `%s` Statement NotResource expected: `%s`", policy.Statements[0].NotResource[0], expectedValue)
		t.Fail()
	}
}

func testStatementNotResourceElementSorts(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "NotResource": ["ccc/*", "bbb/*", "aaa/*"]
        },
        {
          "NotResource": ["arn:aws:dynamodb:us-east-2:account-ID-without-hyphens:table/*", "arn:aws:dynamodb:us-east-1:account-ID-without-hyphens:table/*"]
        }
      ]
    }
    `

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	expectedValue := "aaa/*"
	if policy.Statements[0].NotResource[0] != expectedValue {
		t.Logf("Unexpected Statement NotResource: `%s` Statement NotResource expected: `%s`", policy.Statements[0].NotResource[0], expectedValue)
		t.Fail()
	}

	expectedValue = "bbb/*"
	if policy.Statements[0].NotResource[1] != expectedValue {
		t.Logf("Unexpected Statement NotResource: `%s` Statement NotResource expected: `%s`", policy.Statements[0].NotResource[1], expectedValue)
		t.Fail()
	}

	expectedValue = "ccc/*"
	if policy.Statements[0].NotResource[2] != expectedValue {
		t.Logf("Unexpected Statement NotResource: `%s` Statement NotResource expected: `%s`", policy.Statements[0].NotResource[2], expectedValue)
		t.Fail()
	}

	expectedValue = "arn:aws:dynamodb:us-east-1:account-ID-without-hyphens:table/*"
	if policy.Statements[1].NotResource[0] != expectedValue {
		t.Logf("Unexpected Statement NotResource: `%s` Statement NotResource expected: `%s`", policy.Statements[1].NotResource[0], expectedValue)
		t.Fail()
	}

	expectedValue = "arn:aws:dynamodb:us-east-2:account-ID-without-hyphens:table/*"
	if policy.Statements[1].NotResource[1] != expectedValue {
		t.Logf("Unexpected Statement NotResource: `%s` Statement NotResource expected: `%s`", policy.Statements[1].NotResource[1], expectedValue)
		t.Fail()
	}
}

func testStatementNotResourceElementLeavesCasingUnchanged(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "NotResource": "arn:aws:dynamodb:us-east-1:account-ID-without-hyphens:table/*"
        },
        {
          "NotResource": ["arn:aws:dynamodb:us-east-1:account-ID-without-hyphens:table/*"]
        },
        {
          "NotResource": ["arn:aws:dynamodb:us-east-1:account-ID-without-hyphens:table/*", "arn:aws:dynamodb:us-east-2:account-ID-without-hyphens:table/*"]
        }
      ]
    }
    `

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	expectedValue := "arn:aws:dynamodb:us-east-1:account-ID-without-hyphens:table/*"
	if policy.Statements[0].NotResource[0] != expectedValue {
		t.Logf("Unexpected Statement NotResource: `%s` Statement NotResource expected: `%s`", policy.Statements[0].NotResource[0], expectedValue)
		t.Fail()
	}

	if policy.Statements[1].NotResource[0] != expectedValue {
		t.Logf("Unexpected Statement NotResource: `%s` Statement NotResource expected: `%s`", policy.Statements[1].NotResource[0], expectedValue)
		t.Fail()
	}

	if policy.Statements[2].NotResource[0] != expectedValue {
		t.Logf("Unexpected Statement NotResource: `%s` Statement NotResource expected: `%s`", policy.Statements[2].NotResource[0], expectedValue)
		t.Fail()
	}

	expectedValue = "arn:aws:dynamodb:us-east-2:account-ID-without-hyphens:table/*"
	if policy.Statements[2].NotResource[1] != expectedValue {
		t.Logf("Unexpected Statement NotResource: `%s` Statement NotResource expected: `%s`", policy.Statements[2].NotResource[1], expectedValue)
		t.Fail()
	}
}

func TestStatementConditionElement(t *testing.T) {
	t.Run("TestSingleConditions", testSingleConditions)
	t.Run("TestMultipleConditions", testMultipleConditions)
	t.Run("TestSingleSubCondition", testSingleSubCondition)
	t.Run("TestMultipleSubConditions", testMultipleSubConditions)
	t.Run("TestMultipleValuesInSubConditionsInAscendingOrder", testMultipleValuesInSubConditionsInAscendingOrder)
	t.Run("TestMultipleValuesInSubConditionsInDecendingOrder", testMultipleValuesInSubConditionsInDecendingOrder)
	t.Run("TestRemoveDuplicatesFromMultipleValuesInSubConditions", testRemoveDuplicatesFromMultipleValuesInSubConditions)
	t.Run("TestValueThatAreNotStringInSubConditionValues", testValueThatAreNotStringInSubConditionValues)
	t.Run("testMultipleValuesThatAreNotStringInSubConditionValue", testMultipleValuesThatAreNotStringInSubConditionValue)
	t.Run("testRemoveDuplicatesFromMultipleValuesThatAreNotStringInSubConditionValue", testRemoveDuplicatesFromMultipleValuesThatAreNotStringInSubConditionValue)
	t.Run("TestAllSubConditionCategories", testAllSubConditionCategories)
}

func testSingleConditions(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "Condition": {
            "StringEquals": { "StringCondition": "Value" }
          }
        }
      ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	currentCondition := policy.Statements[0].Condition
	currentConditionCount := len(currentCondition)
	expectedConditionCount := 1
	if currentConditionCount != expectedConditionCount {
		t.Logf("Unexpected number of Conditions: `%d` Conditions expected: `%d`", currentConditionCount, expectedConditionCount)
		t.Fail()
	}

	currentSubCondition := currentCondition["StringEquals"].(map[string]interface{})
	currentSubConditionCount := len(currentSubCondition)
	expectedSubConditionCount := 1
	if currentSubConditionCount != expectedSubConditionCount {
		t.Logf("Unexpected number of Sub Conditions: `%d` Sub Conditions expected: `%d`", currentSubConditionCount, expectedSubConditionCount)
		t.Fail()
	}

	currentSubConditionValues := currentSubCondition["stringcondition"].([]string)
	currentSubConditionValuesCount := len(currentSubConditionValues)
	expectedSubConditionValuesCount := 1
	if currentSubConditionValuesCount != expectedSubConditionValuesCount {
		t.Logf("Unexpected number of Sub Conditions values: `%d` Sub Conditions expected values: `%d`", currentSubConditionValuesCount, expectedSubConditionValuesCount)
		t.Fail()
	}

	currentSubConditionValue := currentSubConditionValues[0]
	expectedSubConditionValue := "Value"
	if currentSubConditionValue != expectedSubConditionValue {
		t.Logf("Unexpected Statement Sub Condition value: `%s` Statement Sub Condition value expected: `%s`", currentSubConditionValue, expectedSubConditionValue)
		t.Fail()
	}
}

func testMultipleConditions(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "Condition": {
            "StringEquals": { "StringCondition": "Value" },
			"NumericEquals": { "NumericCondition": "10" }
          }
        }
      ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	currentCondition := policy.Statements[0].Condition
	currentConditionCount := len(currentCondition)
	expectedConditionCount := 2
	if currentConditionCount != expectedConditionCount {
		t.Logf("Unexpected number of Conditions: `%d` Conditions expected: `%d`", currentConditionCount, expectedConditionCount)
		t.Fail()
	}

	currentSubCondition := currentCondition["StringEquals"].(map[string]interface{})
	currentSubConditionCount := len(currentSubCondition)
	expectedSubConditionCount := 1
	if currentSubConditionCount != expectedSubConditionCount {
		t.Logf("Unexpected number of Sub Conditions: `%d` Sub Conditions expected: `%d`", currentSubConditionCount, expectedSubConditionCount)
		t.Fail()
	}

	currentSubConditionValues := currentSubCondition["stringcondition"].([]string)
	currentSubConditionValuesCount := len(currentSubConditionValues)
	expectedSubConditionValuesCount := 1
	if currentSubConditionValuesCount != expectedSubConditionValuesCount {
		t.Logf("Unexpected number of Sub Conditions values: `%d` Sub Conditions expected values: `%d`", currentSubConditionValuesCount, expectedSubConditionValuesCount)
		t.Fail()
	}

	currentSubConditionValue := currentSubConditionValues[0]
	expectedSubConditionValue := "Value"
	if currentSubConditionValue != expectedSubConditionValue {
		t.Logf("Unexpected Statement Sub Condition value: `%s` Statement Sub Condition value expected: `%s`", currentSubConditionValue, expectedSubConditionValue)
		t.Fail()
	}

	currentSubCondition = currentCondition["NumericEquals"].(map[string]interface{})
	currentSubConditionCount = len(currentSubCondition)
	expectedSubConditionCount = 1
	if currentSubConditionCount != expectedSubConditionCount {
		t.Logf("Unexpected number of Sub Conditions: `%d` Sub Conditions expected: `%d`", currentSubConditionCount, expectedSubConditionCount)
		t.Fail()
	}

	currentSubConditionValues = currentSubCondition["numericcondition"].([]string)
	currentSubConditionValuesCount = len(currentSubConditionValues)
	expectedSubConditionValuesCount = 1
	if currentSubConditionValuesCount != expectedSubConditionValuesCount {
		t.Logf("Unexpected number of Sub Conditions values: `%d` Sub Conditions expected values: `%d`", currentSubConditionValuesCount, expectedSubConditionValuesCount)
		t.Fail()
	}

	currentSubConditionValue = currentSubConditionValues[0]
	expectedSubConditionValue = "10"
	if currentSubConditionValue != expectedSubConditionValue {
		t.Logf("Unexpected Statement Sub Condition value: `%s` Statement Sub Condition value expected: `%s`", currentSubConditionValue, expectedSubConditionValue)
		t.Fail()
	}
}

func testSingleSubCondition(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "Condition": {
            "StringEquals": {
              "StringCondition1": "Value1"
            }
          }
        }
      ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	currentCondition := policy.Statements[0].Condition
	currentConditionCount := len(currentCondition)
	expectedConditionCount := 1
	if currentConditionCount != expectedConditionCount {
		t.Logf("Unexpected number of Conditions: `%d` Conditions expected: `%d`", currentConditionCount, expectedConditionCount)
		t.Fail()
	}

	currentSubCondition := currentCondition["StringEquals"].(map[string]interface{})
	currentSubConditionCount := len(currentSubCondition)
	expectedSubConditionCount := 1
	if currentSubConditionCount != expectedSubConditionCount {
		t.Logf("Unexpected number of Sub Conditions: `%d` Sub Conditions expected: `%d`", currentSubConditionCount, expectedSubConditionCount)
		t.Fail()
	}

	currentSubConditionValues := currentSubCondition["stringcondition1"].([]string)
	currentSubConditionValuesCount := len(currentSubConditionValues)
	expectedSubConditionValuesCount := 1
	if currentSubConditionValuesCount != expectedSubConditionValuesCount {
		t.Logf("Unexpected number of Sub Conditions values: `%d` Sub Conditions expected values: `%d`", currentSubConditionValuesCount, expectedSubConditionValuesCount)
		t.Fail()
	}

	currentSubConditionValue := currentSubConditionValues[0]
	expectedSubConditionValue := "Value1"
	if currentSubConditionValue != expectedSubConditionValue {
		t.Logf("Unexpected Statement Sub Condition value: `%s` Statement Sub Condition value expected: `%s`", currentSubConditionValue, expectedSubConditionValue)
		t.Fail()
	}
}

func testMultipleSubConditions(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "Condition": {
            "StringEquals": {
              "StringCondition1": "Value1",
              "StringCondition2": "Value2"
            }
          }
        }
      ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	currentCondition := policy.Statements[0].Condition
	currentConditionCount := len(currentCondition)
	expectedConditionCount := 1
	if currentConditionCount != expectedConditionCount {
		t.Logf("Unexpected number of Conditions: `%d` Conditions expected: `%d`", currentConditionCount, expectedConditionCount)
		t.Fail()
	}

	currentSubCondition := currentCondition["StringEquals"].(map[string]interface{})
	currentSubConditionCount := len(currentSubCondition)
	expectedSubConditionCount := 2
	if currentSubConditionCount != expectedSubConditionCount {
		t.Logf("Unexpected number of Sub Conditions: `%d` Sub Conditions expected: `%d`", currentSubConditionCount, expectedSubConditionCount)
		t.Fail()
	}

	currentSubConditionValues := currentSubCondition["stringcondition1"].([]string)
	currentSubConditionValuesCount := len(currentSubConditionValues)
	expectedSubConditionValuesCount := 1
	if currentSubConditionValuesCount != expectedSubConditionValuesCount {
		t.Logf("Unexpected number of Sub Conditions values: `%d` Sub Conditions expected values: `%d`", currentSubConditionValuesCount, expectedSubConditionValuesCount)
		t.Fail()
	}

	currentSubConditionValue := currentSubConditionValues[0]
	expectedSubConditionValue := "Value1"
	if currentSubConditionValue != expectedSubConditionValue {
		t.Logf("Unexpected Statement Sub Condition value: `%s` Statement Sub Condition value expected: `%s`", currentSubConditionValue, expectedSubConditionValue)
		t.Fail()
	}

	currentSubConditionValues = currentSubCondition["stringcondition2"].([]string)
	currentSubConditionValuesCount = len(currentSubConditionValues)
	expectedSubConditionValuesCount = 1
	if currentSubConditionValuesCount != expectedSubConditionValuesCount {
		t.Logf("Unexpected number of Sub Conditions values: `%d` Sub Conditions expected values: `%d`", currentSubConditionValuesCount, expectedSubConditionValuesCount)
		t.Fail()
	}

	currentSubConditionValue = currentSubConditionValues[0]
	expectedSubConditionValue = "Value2"
	if currentSubConditionValue != expectedSubConditionValue {
		t.Logf("Unexpected Statement Sub Condition value: `%s` Statement Sub Condition value expected: `%s`", currentSubConditionValue, expectedSubConditionValue)
		t.Fail()
	}
}

func testMultipleValuesInSubConditionsInAscendingOrder(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "Condition": { "StringEquals": { "StringCondition": ["AAAAA", "ZZZZZ"] } }
        }
      ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	currentCondition := policy.Statements[0].Condition
	currentConditionCount := len(currentCondition)
	expectedConditionCount := 1
	if currentConditionCount != expectedConditionCount {
		t.Logf("Unexpected number of Conditions: `%d` Conditions expected: `%d`", currentConditionCount, expectedConditionCount)
		t.Fail()
	}

	currentSubCondition := currentCondition["StringEquals"].(map[string]interface{})
	currentSubConditionCount := len(currentSubCondition)
	expectedSubConditionCount := 1
	if currentSubConditionCount != expectedSubConditionCount {
		t.Logf("Unexpected number of Sub Conditions: `%d` Sub Conditions expected: `%d`", currentSubConditionCount, expectedSubConditionCount)
		t.Fail()
	}

	currentSubConditionValues := currentSubCondition["stringcondition"].([]string)
	currentSubConditionValuesCount := len(currentSubConditionValues)
	expectedSubConditionValuesCount := 2
	if currentSubConditionValuesCount != expectedSubConditionValuesCount {
		t.Logf("Unexpected number of Sub Conditions values: `%d` Sub Conditions expected values: `%d`", currentSubConditionValuesCount, expectedSubConditionValuesCount)
		t.Fail()
	}

	currentSubConditionValue := currentSubConditionValues[0]
	expectedSubConditionValue := "AAAAA"
	if currentSubConditionValue != expectedSubConditionValue {
		t.Logf("Unexpected Statement Sub Condition value: `%s` Statement Sub Condition value expected: `%s`", currentSubConditionValue, expectedSubConditionValue)
		t.Fail()
	}

	currentSubConditionValue = currentSubConditionValues[1]
	expectedSubConditionValue = "ZZZZZ"
	if currentSubConditionValue != expectedSubConditionValue {
		t.Logf("Unexpected Statement Sub Condition value: `%s` Statement Sub Condition value expected: `%s`", currentSubConditionValue, expectedSubConditionValue)
		t.Fail()
	}
}

func testMultipleValuesInSubConditionsInDecendingOrder(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "Condition": { "StringEquals": { "StringCondition": ["ZZZZZ", "AAAAA"] } }
        }
      ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	currentCondition := policy.Statements[0].Condition
	currentConditionCount := len(currentCondition)
	expectedConditionCount := 1
	if currentConditionCount != expectedConditionCount {
		t.Logf("Unexpected number of Conditions: `%d` Conditions expected: `%d`", currentConditionCount, expectedConditionCount)
		t.Fail()
	}

	currentSubCondition := currentCondition["StringEquals"].(map[string]interface{})
	currentSubConditionCount := len(currentSubCondition)
	expectedSubConditionCount := 1
	if currentSubConditionCount != expectedSubConditionCount {
		t.Logf("Unexpected number of Sub Conditions: `%d` Sub Conditions expected: `%d`", currentSubConditionCount, expectedSubConditionCount)
		t.Fail()
	}

	currentSubConditionValues := currentSubCondition["stringcondition"].([]string)
	currentSubConditionValuesCount := len(currentSubConditionValues)
	expectedSubConditionValuesCount := 2
	if currentSubConditionValuesCount != expectedSubConditionValuesCount {
		t.Logf("Unexpected number of Sub Conditions values: `%d` Sub Conditions expected values: `%d`", currentSubConditionValuesCount, expectedSubConditionValuesCount)
		t.Fail()
	}

	currentSubConditionValue := currentSubConditionValues[0]
	expectedSubConditionValue := "AAAAA"
	if currentSubConditionValue != expectedSubConditionValue {
		t.Logf("Unexpected Statement Sub Condition value: `%s` Statement Sub Condition value expected: `%s`", currentSubConditionValue, expectedSubConditionValue)
		t.Fail()
	}

	currentSubConditionValue = currentSubConditionValues[1]
	expectedSubConditionValue = "ZZZZZ"
	if currentSubConditionValue != expectedSubConditionValue {
		t.Logf("Unexpected Statement Sub Condition value: `%s` Statement Sub Condition value expected: `%s`", currentSubConditionValue, expectedSubConditionValue)
		t.Fail()
	}
}

func testRemoveDuplicatesFromMultipleValuesInSubConditions(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "Condition": { "StringEquals": { "StringCondition": ["AAAAA", "AAAAA"] } }
        }
      ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	currentCondition := policy.Statements[0].Condition
	currentConditionCount := len(currentCondition)
	expectedConditionCount := 1
	if currentConditionCount != expectedConditionCount {
		t.Logf("Unexpected number of Conditions: `%d` Conditions expected: `%d`", currentConditionCount, expectedConditionCount)
		t.Fail()
	}

	currentSubCondition := currentCondition["StringEquals"].(map[string]interface{})
	currentSubConditionCount := len(currentSubCondition)
	expectedSubConditionCount := 1
	if currentSubConditionCount != expectedSubConditionCount {
		t.Logf("Unexpected number of Sub Conditions: `%d` Sub Conditions expected: `%d`", currentSubConditionCount, expectedSubConditionCount)
		t.Fail()
	}

	currentSubConditionValues := currentSubCondition["stringcondition"].([]string)
	currentSubConditionValuesCount := len(currentSubConditionValues)
	expectedSubConditionValuesCount := 1
	if currentSubConditionValuesCount != expectedSubConditionValuesCount {
		t.Logf("Unexpected number of Sub Conditions values: `%d` Sub Conditions expected values: `%d`", currentSubConditionValuesCount, expectedSubConditionValuesCount)
		t.Fail()
	}

	currentSubConditionValue := currentSubConditionValues[0]
	expectedSubConditionValue := "AAAAA"
	if currentSubConditionValue != expectedSubConditionValue {
		t.Logf("Unexpected Statement Sub Condition value: `%s` Statement Sub Condition value expected: `%s`", currentSubConditionValue, expectedSubConditionValue)
		t.Fail()
	}
}

func testValueThatAreNotStringInSubConditionValues(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "Condition": {
            "NumericEquals": { "NumericCondition": 10 },
            "Bool": { "BoolCondition": true }
          }
        }
      ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	currentCondition := policy.Statements[0].Condition
	currentConditionCount := len(currentCondition)
	expectedConditionCount := 2
	if currentConditionCount != expectedConditionCount {
		t.Logf("Unexpected number of Conditions: `%d` Conditions expected: `%d`", currentConditionCount, expectedConditionCount)
		t.Fail()
	}

	currentSubCondition := currentCondition["NumericEquals"].(map[string]interface{})
	currentSubConditionCount := len(currentSubCondition)
	expectedSubConditionCount := 1
	if currentSubConditionCount != expectedSubConditionCount {
		t.Logf("Unexpected number of Sub Conditions: `%d` Sub Conditions expected: `%d`", currentSubConditionCount, expectedSubConditionCount)
		t.Fail()
	}

	currentSubConditionValues := currentSubCondition["numericcondition"].([]string)
	currentSubConditionValuesCount := len(currentSubConditionValues)
	expectedSubConditionValuesCount := 1
	if currentSubConditionValuesCount != expectedSubConditionValuesCount {
		t.Logf("Unexpected number of Sub Conditions values: `%d` Sub Conditions expected values: `%d`", currentSubConditionValuesCount, expectedSubConditionValuesCount)
		t.Fail()
	}

	currentSubConditionValue := currentSubConditionValues[0]
	expectedSubConditionValue := "10"
	if currentSubConditionValue != expectedSubConditionValue {
		t.Logf("Unexpected Statement Sub Condition value: `%s` Statement Sub Condition value expected: `%s`", currentSubConditionValue, expectedSubConditionValue)
		t.Fail()
	}

	currentSubCondition = currentCondition["Bool"].(map[string]interface{})
	currentSubConditionCount = len(currentSubCondition)
	expectedSubConditionCount = 1
	if currentSubConditionCount != expectedSubConditionCount {
		t.Logf("Unexpected number of Sub Conditions: `%d` Sub Conditions expected: `%d`", currentSubConditionCount, expectedSubConditionCount)
		t.Fail()
	}

	currentSubConditionValues = currentSubCondition["boolcondition"].([]string)
	currentSubConditionValuesCount = len(currentSubConditionValues)
	expectedSubConditionValuesCount = 1
	if currentSubConditionValuesCount != expectedSubConditionValuesCount {
		t.Logf("Unexpected number of Sub Conditions values: `%d` Sub Conditions expected values: `%d`", currentSubConditionValuesCount, expectedSubConditionValuesCount)
		t.Fail()
	}

	currentSubConditionValue = currentSubConditionValues[0]
	expectedSubConditionValue = "true"
	if currentSubConditionValue != expectedSubConditionValue {
		t.Logf("Unexpected Statement Sub Condition value: `%s` Statement Sub Condition value expected: `%s`", currentSubConditionValue, expectedSubConditionValue)
		t.Fail()
	}
}

func testMultipleValuesThatAreNotStringInSubConditionValue(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "Condition": {
            "NumericEquals": { "NumericCondition": [10, 20] },
            "Bool": { "BoolCondition": [true, false] }
          }
        }
      ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	currentCondition := policy.Statements[0].Condition
	currentConditionCount := len(currentCondition)
	expectedConditionCount := 2
	if currentConditionCount != expectedConditionCount {
		t.Logf("Unexpected number of Conditions: `%d` Conditions expected: `%d`", currentConditionCount, expectedConditionCount)
		t.Fail()
	}

	currentSubCondition := currentCondition["NumericEquals"].(map[string]interface{})
	currentSubConditionCount := len(currentSubCondition)
	expectedSubConditionCount := 1
	if currentSubConditionCount != expectedSubConditionCount {
		t.Logf("Unexpected number of Sub Conditions: `%d` Sub Conditions expected: `%d`", currentSubConditionCount, expectedSubConditionCount)
		t.Fail()
	}

	currentSubConditionValues := currentSubCondition["numericcondition"].([]string)
	currentSubConditionValuesCount := len(currentSubConditionValues)
	expectedSubConditionValuesCount := 2
	if currentSubConditionValuesCount != expectedSubConditionValuesCount {
		t.Logf("Unexpected number of Sub Conditions values: `%d` Sub Conditions expected values: `%d`", currentSubConditionValuesCount, expectedSubConditionValuesCount)
		t.Fail()
	}

	currentSubConditionValue := currentSubConditionValues[0]
	expectedSubConditionValue := "10"
	if currentSubConditionValue != expectedSubConditionValue {
		t.Logf("Unexpected Statement Sub Condition value: `%s` Statement Sub Condition value expected: `%s`", currentSubConditionValue, expectedSubConditionValue)
		t.Fail()
	}

	currentSubConditionValue = currentSubConditionValues[1]
	expectedSubConditionValue = "20"
	if currentSubConditionValue != expectedSubConditionValue {
		t.Logf("Unexpected Statement Sub Condition value: `%s` Statement Sub Condition value expected: `%s`", currentSubConditionValue, expectedSubConditionValue)
		t.Fail()
	}

	currentSubCondition = currentCondition["Bool"].(map[string]interface{})
	currentSubConditionCount = len(currentSubCondition)
	expectedSubConditionCount = 1
	if currentSubConditionCount != expectedSubConditionCount {
		t.Logf("Unexpected number of Sub Conditions: `%d` Sub Conditions expected: `%d`", currentSubConditionCount, expectedSubConditionCount)
		t.Fail()
	}

	currentSubConditionValues = currentSubCondition["boolcondition"].([]string)
	currentSubConditionValuesCount = len(currentSubConditionValues)
	expectedSubConditionValuesCount = 2
	if currentSubConditionValuesCount != expectedSubConditionValuesCount {
		t.Logf("Unexpected number of Sub Conditions values: `%d` Sub Conditions expected values: `%d`", currentSubConditionValuesCount, expectedSubConditionValuesCount)
		t.Fail()
	}

	currentSubConditionValue = currentSubConditionValues[0]
	expectedSubConditionValue = "false"
	if currentSubConditionValue != expectedSubConditionValue {
		t.Logf("Unexpected Statement Sub Condition value: `%s` Statement Sub Condition value expected: `%s`", currentSubConditionValue, expectedSubConditionValue)
		t.Fail()
	}

	currentSubConditionValue = currentSubConditionValues[1]
	expectedSubConditionValue = "true"
	if currentSubConditionValue != expectedSubConditionValue {
		t.Logf("Unexpected Statement Sub Condition value: `%s` Statement Sub Condition value expected: `%s`", currentSubConditionValue, expectedSubConditionValue)
		t.Fail()
	}
}

func testRemoveDuplicatesFromMultipleValuesThatAreNotStringInSubConditionValue(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "Condition": {
            "NumericEquals": { "NumericCondition": [10, 10] },
            "Bool": { "BoolCondition": [true, true] }
          }
        }
      ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	currentCondition := policy.Statements[0].Condition
	currentConditionCount := len(currentCondition)
	expectedConditionCount := 2
	if currentConditionCount != expectedConditionCount {
		t.Logf("Unexpected number of Conditions: `%d` Conditions expected: `%d`", currentConditionCount, expectedConditionCount)
		t.Fail()
	}

	currentSubCondition := currentCondition["NumericEquals"].(map[string]interface{})
	currentSubConditionCount := len(currentSubCondition)
	expectedSubConditionCount := 1
	if currentSubConditionCount != expectedSubConditionCount {
		t.Logf("Unexpected number of Sub Conditions: `%d` Sub Conditions expected: `%d`", currentSubConditionCount, expectedSubConditionCount)
		t.Fail()
	}

	currentSubConditionValues := currentSubCondition["numericcondition"].([]string)
	currentSubConditionValuesCount := len(currentSubConditionValues)
	expectedSubConditionValuesCount := 1
	if currentSubConditionValuesCount != expectedSubConditionValuesCount {
		t.Logf("Unexpected number of Sub Conditions values: `%d` Sub Conditions expected values: `%d`", currentSubConditionValuesCount, expectedSubConditionValuesCount)
		t.Fail()
	}

	currentSubConditionValue := currentSubConditionValues[0]
	expectedSubConditionValue := "10"
	if currentSubConditionValue != expectedSubConditionValue {
		t.Logf("Unexpected Statement Sub Condition value: `%s` Statement Sub Condition value expected: `%s`", currentSubConditionValue, expectedSubConditionValue)
		t.Fail()
	}

	currentSubCondition = currentCondition["Bool"].(map[string]interface{})
	currentSubConditionCount = len(currentSubCondition)
	expectedSubConditionCount = 1
	if currentSubConditionCount != expectedSubConditionCount {
		t.Logf("Unexpected number of Sub Conditions: `%d` Sub Conditions expected: `%d`", currentSubConditionCount, expectedSubConditionCount)
		t.Fail()
	}

	currentSubConditionValues = currentSubCondition["boolcondition"].([]string)
	currentSubConditionValuesCount = len(currentSubConditionValues)
	expectedSubConditionValuesCount = 1
	if currentSubConditionValuesCount != expectedSubConditionValuesCount {
		t.Logf("Unexpected number of Sub Conditions values: `%d` Sub Conditions expected values: `%d`", currentSubConditionValuesCount, expectedSubConditionValuesCount)
		t.Fail()
	}

	currentSubConditionValue = currentSubConditionValues[0]
	expectedSubConditionValue = "true"
	if currentSubConditionValue != expectedSubConditionValue {
		t.Logf("Unexpected Statement Sub Condition value: `%s` Statement Sub Condition value expected: `%s`", currentSubConditionValue, expectedSubConditionValue)
		t.Fail()
	}
}

func testAllSubConditionCategories(t *testing.T) {
	// Set up
	policyContent := `
    {
      "Statement": [
        {
          "Condition": { "StringEquals": { "StringCondition": "Value" } }
        },
        {
          "Condition": { "NumericEquals": { "NumericCondition": "10" } }
        },
        {
          "Condition": { "DateGreaterThan": { "DateCondition": "2020-01-01T00:00:01Z" } }
        },
        {
          "Condition": { "Bool": { "BoolCondition": "false" } }
        },
        {
          "Condition": { "BinaryEquals": { "BinaryCondition": "QmluYXJ5VmFsdWVJbkJhc2U2NA==" } }
        },
        {
          "Condition": { "IpAddress": { "IpCondition": "203.0.113.0/24" } }
        },
        {
          "Condition": { "ArnEquals": { "ArnCondition": "arn:aws:sns:REGION:123456789012:TOPIC-ID" } }
        },
        {
          "Condition": { "Null": { "AnyCondition": "true" } }
        }
      ]
    }
	`

	// Test
	policyInterface, err := canonicalPolicy(policyContent)

	// Expected
	if err != nil {
		t.Fatalf("Unexpected error while parsing data: %s", err)
	}

	policy := policyInterface.(Policy)

	currentCondition := policy.Statements[0].Condition
	currentConditionCount := len(currentCondition)
	expectedConditionCount := 1
	if currentConditionCount != expectedConditionCount {
		t.Logf("Unexpected number of Conditions: `%d` Conditions expected: `%d`", currentConditionCount, expectedConditionCount)
		t.Fail()
	}

	currentSubCondition := currentCondition["StringEquals"].(map[string]interface{})
	currentSubConditionCount := len(currentSubCondition)
	expectedSubConditionCount := 1
	if currentSubConditionCount != expectedSubConditionCount {
		t.Logf("Unexpected number of Sub Conditions: `%d` Sub Conditions expected: `%d`", currentSubConditionCount, expectedSubConditionCount)
		t.Fail()
	}

	currentSubConditionValues := currentSubCondition["stringcondition"].([]string)
	currentSubConditionValuesCount := len(currentSubConditionValues)
	expectedSubConditionValuesCount := 1
	if currentSubConditionValuesCount != expectedSubConditionValuesCount {
		t.Logf("Unexpected number of Sub Conditions values: `%d` Sub Conditions expected values: `%d`", currentSubConditionValuesCount, expectedSubConditionValuesCount)
		t.Fail()
	}

	currentSubConditionValue := currentSubConditionValues[0]
	expectedSubConditionValue := "Value"
	if currentSubConditionValue != expectedSubConditionValue {
		t.Logf("Unexpected Statement Sub Condition value: `%s` Statement Sub Condition value expected: `%s`", currentSubConditionValue, expectedSubConditionValue)
		t.Fail()
	}

	currentCondition = policy.Statements[1].Condition
	currentConditionCount = len(currentCondition)
	expectedConditionCount = 1
	if currentConditionCount != expectedConditionCount {
		t.Logf("Unexpected number of Conditions: `%d` Conditions expected: `%d`", currentConditionCount, expectedConditionCount)
		t.Fail()
	}

	currentSubCondition = currentCondition["NumericEquals"].(map[string]interface{})
	currentSubConditionCount = len(currentSubCondition)
	expectedSubConditionCount = 1
	if currentSubConditionCount != expectedSubConditionCount {
		t.Logf("Unexpected number of Sub Conditions: `%d` Sub Conditions expected: `%d`", currentSubConditionCount, expectedSubConditionCount)
		t.Fail()
	}

	currentSubConditionValues = currentSubCondition["numericcondition"].([]string)
	currentSubConditionValuesCount = len(currentSubConditionValues)
	expectedSubConditionValuesCount = 1
	if currentSubConditionValuesCount != expectedSubConditionValuesCount {
		t.Logf("Unexpected number of Sub Conditions values: `%d` Sub Conditions expected values: `%d`", currentSubConditionValuesCount, expectedSubConditionValuesCount)
		t.Fail()
	}

	currentSubConditionValue = currentSubConditionValues[0]
	expectedSubConditionValue = "10"
	if currentSubConditionValue != expectedSubConditionValue {
		t.Logf("Unexpected Statement Sub Condition value: `%s` Statement Sub Condition value expected: `%s`", currentSubConditionValue, expectedSubConditionValue)
		t.Fail()
	}

	currentCondition = policy.Statements[2].Condition
	currentConditionCount = len(currentCondition)
	expectedConditionCount = 1
	if currentConditionCount != expectedConditionCount {
		t.Logf("Unexpected number of Conditions: `%d` Conditions expected: `%d`", currentConditionCount, expectedConditionCount)
		t.Fail()
	}

	currentSubCondition = currentCondition["DateGreaterThan"].(map[string]interface{})
	currentSubConditionCount = len(currentSubCondition)
	expectedSubConditionCount = 1
	if currentSubConditionCount != expectedSubConditionCount {
		t.Logf("Unexpected number of Sub Conditions: `%d` Sub Conditions expected: `%d`", currentSubConditionCount, expectedSubConditionCount)
		t.Fail()
	}

	currentSubConditionValues = currentSubCondition["datecondition"].([]string)
	currentSubConditionValuesCount = len(currentSubConditionValues)
	expectedSubConditionValuesCount = 1
	if currentSubConditionValuesCount != expectedSubConditionValuesCount {
		t.Logf("Unexpected number of Sub Conditions values: `%d` Sub Conditions expected values: `%d`", currentSubConditionValuesCount, expectedSubConditionValuesCount)
		t.Fail()
	}

	currentSubConditionValue = currentSubConditionValues[0]
	expectedSubConditionValue = "2020-01-01T00:00:01Z"
	if currentSubConditionValue != expectedSubConditionValue {
		t.Logf("Unexpected Statement Sub Condition value: `%s` Statement Sub Condition value expected: `%s`", currentSubConditionValue, expectedSubConditionValue)
		t.Fail()
	}

	currentCondition = policy.Statements[3].Condition
	currentConditionCount = len(currentCondition)
	expectedConditionCount = 1
	if currentConditionCount != expectedConditionCount {
		t.Logf("Unexpected number of Conditions: `%d` Conditions expected: `%d`", currentConditionCount, expectedConditionCount)
		t.Fail()
	}

	currentSubCondition = currentCondition["Bool"].(map[string]interface{})
	currentSubConditionCount = len(currentSubCondition)
	expectedSubConditionCount = 1
	if currentSubConditionCount != expectedSubConditionCount {
		t.Logf("Unexpected number of Sub Conditions: `%d` Sub Conditions expected: `%d`", currentSubConditionCount, expectedSubConditionCount)
		t.Fail()
	}

	currentSubConditionValues = currentSubCondition["boolcondition"].([]string)
	currentSubConditionValuesCount = len(currentSubConditionValues)
	expectedSubConditionValuesCount = 1
	if currentSubConditionValuesCount != expectedSubConditionValuesCount {
		t.Logf("Unexpected number of Sub Conditions values: `%d` Sub Conditions expected values: `%d`", currentSubConditionValuesCount, expectedSubConditionValuesCount)
		t.Fail()
	}

	currentSubConditionValue = currentSubConditionValues[0]
	expectedSubConditionValue = "false"
	if currentSubConditionValue != expectedSubConditionValue {
		t.Logf("Unexpected Statement Sub Condition value: `%s` Statement Sub Condition value expected: `%s`", currentSubConditionValue, expectedSubConditionValue)
		t.Fail()
	}

	currentCondition = policy.Statements[4].Condition
	currentConditionCount = len(currentCondition)
	expectedConditionCount = 1
	if currentConditionCount != expectedConditionCount {
		t.Logf("Unexpected number of Conditions: `%d` Conditions expected: `%d`", currentConditionCount, expectedConditionCount)
		t.Fail()
	}

	currentSubCondition = currentCondition["BinaryEquals"].(map[string]interface{})
	currentSubConditionCount = len(currentSubCondition)
	expectedSubConditionCount = 1
	if currentSubConditionCount != expectedSubConditionCount {
		t.Logf("Unexpected number of Sub Conditions: `%d` Sub Conditions expected: `%d`", currentSubConditionCount, expectedSubConditionCount)
		t.Fail()
	}

	currentSubConditionValues = currentSubCondition["binarycondition"].([]string)
	currentSubConditionValuesCount = len(currentSubConditionValues)
	expectedSubConditionValuesCount = 1
	if currentSubConditionValuesCount != expectedSubConditionValuesCount {
		t.Logf("Unexpected number of Sub Conditions values: `%d` Sub Conditions expected values: `%d`", currentSubConditionValuesCount, expectedSubConditionValuesCount)
		t.Fail()
	}

	currentSubConditionValue = currentSubConditionValues[0]
	expectedSubConditionValue = "QmluYXJ5VmFsdWVJbkJhc2U2NA=="
	if currentSubConditionValue != expectedSubConditionValue {
		t.Logf("Unexpected Statement Sub Condition value: `%s` Statement Sub Condition value expected: `%s`", currentSubConditionValue, expectedSubConditionValue)
		t.Fail()
	}

	currentCondition = policy.Statements[5].Condition
	currentConditionCount = len(currentCondition)
	expectedConditionCount = 1
	if currentConditionCount != expectedConditionCount {
		t.Logf("Unexpected number of Conditions: `%d` Conditions expected: `%d`", currentConditionCount, expectedConditionCount)
		t.Fail()
	}

	currentSubCondition = currentCondition["IpAddress"].(map[string]interface{})
	currentSubConditionCount = len(currentSubCondition)
	expectedSubConditionCount = 1
	if currentSubConditionCount != expectedSubConditionCount {
		t.Logf("Unexpected number of Sub Conditions: `%d` Sub Conditions expected: `%d`", currentSubConditionCount, expectedSubConditionCount)
		t.Fail()
	}

	currentSubConditionValues = currentSubCondition["ipcondition"].([]string)
	currentSubConditionValuesCount = len(currentSubConditionValues)
	expectedSubConditionValuesCount = 1
	if currentSubConditionValuesCount != expectedSubConditionValuesCount {
		t.Logf("Unexpected number of Sub Conditions values: `%d` Sub Conditions expected values: `%d`", currentSubConditionValuesCount, expectedSubConditionValuesCount)
		t.Fail()
	}

	currentSubConditionValue = currentSubConditionValues[0]
	expectedSubConditionValue = "203.0.113.0/24"
	if currentSubConditionValue != expectedSubConditionValue {
		t.Logf("Unexpected Statement Sub Condition value: `%s` Statement Sub Condition value expected: `%s`", currentSubConditionValue, expectedSubConditionValue)
		t.Fail()
	}

	currentCondition = policy.Statements[6].Condition
	currentConditionCount = len(currentCondition)
	expectedConditionCount = 1
	if currentConditionCount != expectedConditionCount {
		t.Logf("Unexpected number of Conditions: `%d` Conditions expected: `%d`", currentConditionCount, expectedConditionCount)
		t.Fail()
	}

	currentSubCondition = currentCondition["ArnEquals"].(map[string]interface{})
	currentSubConditionCount = len(currentSubCondition)
	expectedSubConditionCount = 1
	if currentSubConditionCount != expectedSubConditionCount {
		t.Logf("Unexpected number of Sub Conditions: `%d` Sub Conditions expected: `%d`", currentSubConditionCount, expectedSubConditionCount)
		t.Fail()
	}

	currentSubConditionValues = currentSubCondition["arncondition"].([]string)
	currentSubConditionValuesCount = len(currentSubConditionValues)
	expectedSubConditionValuesCount = 1
	if currentSubConditionValuesCount != expectedSubConditionValuesCount {
		t.Logf("Unexpected number of Sub Conditions values: `%d` Sub Conditions expected values: `%d`", currentSubConditionValuesCount, expectedSubConditionValuesCount)
		t.Fail()
	}

	currentSubConditionValue = currentSubConditionValues[0]
	expectedSubConditionValue = "arn:aws:sns:REGION:123456789012:TOPIC-ID"
	if currentSubConditionValue != expectedSubConditionValue {
		t.Logf("Unexpected Statement Sub Condition value: `%s` Statement Sub Condition value expected: `%s`", currentSubConditionValue, expectedSubConditionValue)
		t.Fail()
	}

	currentCondition = policy.Statements[7].Condition
	currentConditionCount = len(currentCondition)
	expectedConditionCount = 1
	if currentConditionCount != expectedConditionCount {
		t.Logf("Unexpected number of Conditions: `%d` Conditions expected: `%d`", currentConditionCount, expectedConditionCount)
		t.Fail()
	}

	currentSubCondition = currentCondition["Null"].(map[string]interface{})
	currentSubConditionCount = len(currentSubCondition)
	expectedSubConditionCount = 1
	if currentSubConditionCount != expectedSubConditionCount {
		t.Logf("Unexpected number of Sub Conditions: `%d` Sub Conditions expected: `%d`", currentSubConditionCount, expectedSubConditionCount)
		t.Fail()
	}

	currentSubConditionValues = currentSubCondition["anycondition"].([]string)
	currentSubConditionValuesCount = len(currentSubConditionValues)
	expectedSubConditionValuesCount = 1
	if currentSubConditionValuesCount != expectedSubConditionValuesCount {
		t.Logf("Unexpected number of Sub Conditions values: `%d` Sub Conditions expected values: `%d`", currentSubConditionValuesCount, expectedSubConditionValuesCount)
		t.Fail()
	}

	currentSubConditionValue = currentSubConditionValues[0]
	expectedSubConditionValue = "true"
	if currentSubConditionValue != expectedSubConditionValue {
		t.Logf("Unexpected Statement Sub Condition value: `%s` Statement Sub Condition value expected: `%s`", currentSubConditionValue, expectedSubConditionValue)
		t.Fail()
	}
}

// TODO: Add some tests to see if we have already canonicalised a policy, what does this function do.
