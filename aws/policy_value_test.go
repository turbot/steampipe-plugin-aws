package aws

import "testing"

func TestPolicyValueContainsNoWildcards(t *testing.T) {
	t.Run("TestContainsWithoutWildcardAndMatch", testContainsWithoutWildcardMatch)
	t.Run("TestContainsWithoutWildcardAndNoMatch1", testContainsWithoutWildcardNoMatch1)
	t.Run("TestContainsWithoutWildcardAndNoMatch2", testContainsWithoutWildcardNoMatch2)
	t.Run("TestContainsWithoutWildcardAndExtraCharacterMatch", testContainsWithoutWildcardAndExtraCharacterMatch)
	t.Run("TestContainsWithoutWildcardAndExtraCharacterNoMatch1", testContainsWithoutWildcardAndExtraCharacterNoMatch1)
	t.Run("TestContainsWithoutWildcardAndExtraCharacterNoMatch2", testContainsWithoutWildcardAndExtraCharacterNoMatch2)
}

func testContainsWithoutWildcardMatch(t *testing.T) {
	// Set up
	firstValue := "z"
	compareValue := "z"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithoutWildcardNoMatch1(t *testing.T) {
	// Set up
	firstValue := "z"
	compareValue := "a"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithoutWildcardNoMatch2(t *testing.T) {
	// Set up
	firstValue := "z"
	compareValue := "za"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithoutWildcardAndExtraCharacterMatch(t *testing.T) {
	// Set up
	firstValue := "za"
	compareValue := "za"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithoutWildcardAndExtraCharacterNoMatch1(t *testing.T) {
	// Set up
	firstValue := "za"
	compareValue := "z"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithoutWildcardAndExtraCharacterNoMatch2(t *testing.T) {
	// Set up
	firstValue := "za"
	compareValue := "zaa"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func TestPolicyValueContainsPartialWildcards(t *testing.T) {
	t.Run("TestContainsWithPartialWildcardMatch", testContainsWithPartialWildcardMatch1)
	t.Run("TestContainsWithPartialWildcardMatch", testContainsWithPartialWildcardMatch2)
	t.Run("TestContainsWithPartialWildcardNoMatch1", testContainsWithPartialWildcardNoMatch1)
	t.Run("TestContainsWithPartialWildcardNoMatch2", testContainsWithPartialWildcardNoMatch2)
	t.Run("TestContainsWithPartialWildcardNoMatch2", testContainsWithPartialWildcardNoMatch3)
	t.Run("TestContainsWithPartialWildcardNoMatch4", testContainsWithPartialWildcardNoMatch4)
	t.Run("TestContainsWithPartialWildcardNoMatch5", testContainsWithPartialWildcardNoMatch5)

	t.Run("TestContainsWithPartialWildcardAndCharacterAtStartMatch", testContainsWithPartialWildcardAndCharacterAtStartMatch)
	t.Run("TestContainsWithPartialWildcardAndCharacterAtStartNoMatch1", testContainsWithPartialWildcardAndCharacterAtStartNoMatch1)
	t.Run("TestContainsWithPartialWildcardAndCharacterAtStartNoMatch2", testContainsWithPartialWildcardAndCharacterAtStartNoMatch2)
	t.Run("TestContainsWithPartialWildcardAndCharacterAtEndMatch", testContainsWithPartialWildcardAndCharacterAtEndMatch)
	t.Run("TestContainsWithPartialWildcardAndCharacterAtEndNoMatch1", testContainsWithPartialWildcardAndCharacterAtEndNoMatch1)
	t.Run("TestContainsWithPartialWildcardAndCharacterAtEndNoMatch2", testContainsWithPartialWildcardAndCharacterAtEndNoMatch2)
	t.Run("TestContainsWithPartialWildcardAndCharacterAtStartAndEndMatch", testContainsWithPartialWildcardAndCharacterAtStartAndEndMatch)
	t.Run("TestContainsWithPartialWildcardAndCharacterAtStartAndEndNoMatch1", testContainsWithPartialWildcardAndCharacterAtStartAndEndNoMatch1)
	t.Run("TestContainsWithPartialWildcardAndCharacterAtStartAndEndNoMatch2", testContainsWithPartialWildcardAndCharacterAtStartAndEndNoMatch2)

	t.Run("TestContainsWithDoublePartialWildcardMatch", testContainsWithDoublePartialWildcardMatch)
	t.Run("TestContainsWithDoublePartialWildcardNoMatch1", testContainsWithDoublePartialWildcardNoMatch1)
	t.Run("TestContainsWithDoublePartialWildcardNoMatch2", testContainsWithDoublePartialWildcardNoMatch2)
	t.Run("TestContainsWithDoublePartialWildcardAndCharacterAtStartMatch", testContainsWithDoublePartialWildcardAndCharacterAtStartMatch)
	t.Run("TestContainsWithDoublePartialWildcardAndCharacterAtStartNoMatch1", testContainsWithDoublePartialWildcardAndCharacterAtStartNoMatch1)
	t.Run("TestContainsWithDoublePartialWildcardAndCharacterAtStartNoMatch2", testContainsWithDoublePartialWildcardAndCharacterAtStartNoMatch2)
	t.Run("TestContainsWithDoublePartialWildcardAndCharacterAtEndMatch", testContainsWithDoublePartialWildcardAndCharacterAtEndMatch)
	t.Run("TestContainsWithDoublePartialWildcardAndCharacterAtEndNoMatch1", testContainsWithDoublePartialWildcardAndCharacterAtEndNoMatch1)
	t.Run("TestContainsWithDoublePartialWildcardAndCharacterAtEndNoMatch2", testContainsWithDoublePartialWildcardAndCharacterAtEndNoMatch2)
	t.Run("TestContainsWithDoublePartialWildcardAndCharacterAtStartAndEndMatch", testContainsWithDoublePartialWildcardAndCharacterAtStartAndEndMatch)
	t.Run("TestContainsWithDoublePartialWildcardAndCharacterAtStartAndEndNoMatch1", testContainsWithDoublePartialWildcardAndCharacterAtStartAndEndNoMatch1)
	t.Run("TestContainsWithDoublePartialWildcardAndCharacterAtStartAndEndNoMatch2", testContainsWithDoublePartialWildcardAndCharacterAtStartAndEndNoMatch2)
}

func testContainsWithPartialWildcardMatch1(t *testing.T) {
	// Set up
	firstValue := "?"
	compareValue := "z"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithPartialWildcardMatch2(t *testing.T) {
	// Set up
	firstValue := "?"
	compareValue := "?"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithPartialWildcardNoMatch1(t *testing.T) {
	// Set up
	firstValue := "?"
	compareValue := ""

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithPartialWildcardNoMatch2(t *testing.T) {
	// Set up
	firstValue := "?"
	compareValue := "zz"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithPartialWildcardNoMatch3(t *testing.T) {
	// Set up
	firstValue := "?"
	compareValue := "??"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithPartialWildcardNoMatch4(t *testing.T) {
	// Set up
	firstValue := "?"
	compareValue := "?*"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithPartialWildcardNoMatch5(t *testing.T) {
	// Set up
	firstValue := "?"
	compareValue := "*"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithPartialWildcardAndCharacterAtStartMatch(t *testing.T) {
	// Set up
	firstValue := "a?"
	compareValue := "az"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithPartialWildcardAndCharacterAtStartNoMatch1(t *testing.T) {
	// Set up
	firstValue := "a?"
	compareValue := "a"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithPartialWildcardAndCharacterAtStartNoMatch2(t *testing.T) {
	// Set up
	firstValue := "a?"
	compareValue := "a"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithPartialWildcardAndCharacterAtEndMatch(t *testing.T) {
	// Set up
	firstValue := "?z"
	compareValue := "az"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithPartialWildcardAndCharacterAtEndNoMatch1(t *testing.T) {
	// Set up
	firstValue := "?z"
	compareValue := "azz"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithPartialWildcardAndCharacterAtEndNoMatch2(t *testing.T) {
	// Set up
	firstValue := "?z"
	compareValue := "a"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithPartialWildcardAndCharacterAtStartAndEndMatch(t *testing.T) {
	// Set up
	firstValue := "a?z"
	compareValue := "agz"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithPartialWildcardAndCharacterAtStartAndEndNoMatch1(t *testing.T) {
	// Set up
	firstValue := "a?z"
	compareValue := "ag"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithPartialWildcardAndCharacterAtStartAndEndNoMatch2(t *testing.T) {
	// Set up
	firstValue := "a?z"
	compareValue := "agzg"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithDoublePartialWildcardMatch(t *testing.T) {
	// Set up
	firstValue := "??"
	compareValue := "zz"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithDoublePartialWildcardNoMatch1(t *testing.T) {
	// Set up
	firstValue := "??"
	compareValue := "z"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithDoublePartialWildcardNoMatch2(t *testing.T) {
	// Set up
	firstValue := "??"
	compareValue := "zzz"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithDoublePartialWildcardAndCharacterAtStartMatch(t *testing.T) {
	// Set up
	firstValue := "a??"
	compareValue := "azz"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithDoublePartialWildcardAndCharacterAtStartNoMatch1(t *testing.T) {
	// Set up
	firstValue := "a??"
	compareValue := "az"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithDoublePartialWildcardAndCharacterAtStartNoMatch2(t *testing.T) {
	// Set up
	firstValue := "a??"
	compareValue := "azzz"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithDoublePartialWildcardAndCharacterAtEndMatch(t *testing.T) {
	// Set up
	firstValue := "??z"
	compareValue := "ggz"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithDoublePartialWildcardAndCharacterAtEndNoMatch1(t *testing.T) {
	// Set up
	firstValue := "??z"
	compareValue := "ggza"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithDoublePartialWildcardAndCharacterAtEndNoMatch2(t *testing.T) {
	// Set up
	firstValue := "??z"
	compareValue := "gg"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithDoublePartialWildcardAndCharacterAtStartAndEndMatch(t *testing.T) {
	// Set up
	firstValue := "a??z"
	compareValue := "aggz"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithDoublePartialWildcardAndCharacterAtStartAndEndNoMatch1(t *testing.T) {
	// Set up
	firstValue := "a??z"
	compareValue := "agg"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithDoublePartialWildcardAndCharacterAtStartAndEndNoMatch2(t *testing.T) {
	// Set up
	firstValue := "a??z"
	compareValue := "aggzg"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func TestPolicyValueContainsFullWildcards(t *testing.T) {
	t.Run("TestContainsWithFullWildcardMatch1", testContainsWithFullWildcardMatch1)
	t.Run("TestContainsWithFullWildcardMatch2", testContainsWithFullWildcardMatch2)
	t.Run("TestContainsWithFullWildcardMatch3", testContainsWithFullWildcardMatch3)
	t.Run("TestContainsWithFullWildcardMatch4", testContainsWithFullWildcardMatch4)

	t.Run("TestContainsWithFullWildcardAndCharacterAtStartMatch1", testContainsWithFullWildcardAndCharacterAtStartMatch1)
	t.Run("TestContainsWithFullWildcardAndCharacterAtStartMatch2", testContainsWithFullWildcardAndCharacterAtStartMatch2)
	t.Run("TestContainsWithFullWildcardAndCharacterAtStartMatch3", testContainsWithFullWildcardAndCharacterAtStartMatch3)
	t.Run("TestContainsWithFullWildcardAndCharacterAtStartNoMatch", testContainsWithFullWildcardAndCharacterAtStartNoMatch)
	t.Run("TestContainsWithFullWildcardAndCharacterAtStartNoMatch", testContainsWithFullWildcardAndCharacterAtStartNoMatch)
	t.Run("TestContainsWithFullWildcardAndCharacterAtEndMatch1", testContainsWithFullWildcardAndCharacterAtEndMatch1)
	t.Run("TestContainsWithFullWildcardAndCharacterAtEndMatch2", testContainsWithFullWildcardAndCharacterAtEndMatch2)
	t.Run("TestContainsWithFullWildcardAndCharacterAtEndMatch3", testContainsWithFullWildcardAndCharacterAtEndMatch3)
	t.Run("TestContainsWithFullWildcardAndCharacterAtEndNoMatch", testContainsWithFullWildcardAndCharacterAtEndNoMatch)
	t.Run("TestContainsWithFullWildcardAndCharacterAtStartAndEndMatch1", testContainsWithFullWildcardAndCharacterAtStartAndEndMatch1)
	t.Run("TestContainsWithFullWildcardAndCharacterAtStartAndEndMatch2", testContainsWithFullWildcardAndCharacterAtStartAndEndMatch2)
	t.Run("TestContainsWithFullWildcardAndCharacterAtStartAndEndMatch3", testContainsWithFullWildcardAndCharacterAtStartAndEndMatch3)
	t.Run("TestContainsWithFullWildcardAndCharacterAtStartAndEndNoMatch", testContainsWithFullWildcardAndCharacterAtStartAndEndNoMatch)

	t.Run("TestContainsWithDoubleFullWildcardMatch1", testContainsWithDoubleFullWildcardMatch1)
	t.Run("TestContainsWithDoubleFullWildcardMatch2", testContainsWithDoubleFullWildcardMatch2)
	t.Run("TestContainsWithDoubleFullWildcardAndCharacterAtStartMatch1", testContainsWithDoubleFullWildcardAndCharacterAtStartMatch1)
	t.Run("TestContainsWithDoubleFullWildcardAndCharacterAtStartMatch2", testContainsWithDoubleFullWildcardAndCharacterAtStartMatch2)
	t.Run("TestContainsWithDoubleFullWildcardAndCharacterAtStartMatch3", testContainsWithDoubleFullWildcardAndCharacterAtStartMatch3)
	t.Run("TestContainsWithDoubleFullWildcardAndCharacterAtStartNoMatch", testContainsWithDoubleFullWildcardAndCharacterAtStartNoMatch)
	t.Run("TestContainsWithDoubleFullWildcardAndCharacterAtEndMatch1", testContainsWithDoubleFullWildcardAndCharacterAtEndMatch1)
	t.Run("TestContainsWithDoubleFullWildcardAndCharacterAtEndMatch2", testContainsWithDoubleFullWildcardAndCharacterAtEndMatch2)
	t.Run("TestContainsWithDoubleFullWildcardAndCharacterAtEndMatch3", testContainsWithDoubleFullWildcardAndCharacterAtEndMatch3)
	t.Run("TestContainsWithDoubleFullWildcardAndCharacterAtEndNoMatch", testContainsWithDoubleFullWildcardAndCharacterAtEndNoMatch)
	t.Run("TestContainsWithDoubleFullWildcardAndCharacterAtStartAndEndMatch1", testContainsWithDoubleFullWildcardAndCharacterAtStartAndEndMatch1)
	t.Run("TestContainsWithDoubleFullWildcardAndCharacterAtStartAndEndMatch2", testContainsWithDoubleFullWildcardAndCharacterAtStartAndEndMatch2)
	t.Run("TestContainsWithDoubleFullWildcardAndCharacterAtStartAndEndMatch3", testContainsWithDoubleFullWildcardAndCharacterAtStartAndEndMatch3)
	t.Run("TestContainsWithDoubleFullWildcardAndCharacterAtStartAndEndNoMatch", testContainsWithDoubleFullWildcardAndCharacterAtStartAndEndNoMatch)

	t.Run("TestContainsWithTwoFullWildcardMatch1", testContainsWithTwoFullWildcardMatch1)
	t.Run("TestContainsWithTwoFullWildcardMatch2", testContainsWithTwoFullWildcardMatch2)
	t.Run("TestContainsWithTwoFullWildcardMatch3", testContainsWithTwoFullWildcardMatch3)
	t.Run("TestContainsWithTwoFullWildcardMatch4", testContainsWithTwoFullWildcardMatch4)
	t.Run("TestContainsWithTwoFullWildcardNoMatch1", testContainsWithTwoFullWildcardNoMatch1)
	t.Run("TestContainsWithTwoFullWildcardNoMatch2", testContainsWithTwoFullWildcardNoMatch2)
	t.Run("TestContainsWithTwoFullWildcardNoMatch3", testContainsWithTwoFullWildcardNoMatch3)
	t.Run("TestContainsWithTwoFullWildcardNoMatch4", testContainsWithTwoFullWildcardNoMatch4)
	t.Run("TestContainsWithTwoFullWildcardNoMatch5", testContainsWithTwoFullWildcardNoMatch5)

	t.Run("TestContainsWithRepeatingPatternMatch1", testContainsWithRepeatingPatternMatch1)
	t.Run("TestContainsWithRepeatingPatternMatch2", testContainsWithRepeatingPatternMatch2)
	t.Run("TestContainsWithRepeatingPatternMatch3", testContainsWithRepeatingPatternMatch3)
	t.Run("TestContainsWithRepeatingPatternMatch4", testContainsWithRepeatingPatternMatch4)
	t.Run("TestContainsWithRepeatingPatternMatch5", testContainsWithRepeatingPatternMatch5)

	t.Run("TestContainsWithRepeatingPatternNoMatch", testContainsWithRepeatingPatternNoMatch)
}

func testContainsWithFullWildcardMatch1(t *testing.T) {
	// Set up
	firstValue := "*"
	compareValue := ""

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithFullWildcardMatch2(t *testing.T) {
	// Set up
	firstValue := "*"
	compareValue := "z"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithFullWildcardMatch3(t *testing.T) {
	// Set up
	firstValue := "*"
	compareValue := "?"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithFullWildcardMatch4(t *testing.T) {
	// Set up
	firstValue := "*"
	compareValue := "*"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithFullWildcardAndCharacterAtStartMatch1(t *testing.T) {
	// Set up
	firstValue := "a*"
	compareValue := "a"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithFullWildcardAndCharacterAtStartMatch2(t *testing.T) {
	// Set up
	firstValue := "a*"
	compareValue := "aa"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithFullWildcardAndCharacterAtStartMatch3(t *testing.T) {
	// Set up
	firstValue := "a*"
	compareValue := "aaa"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithFullWildcardAndCharacterAtStartNoMatch(t *testing.T) {
	// Set up
	firstValue := "a*"
	compareValue := "b"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithFullWildcardAndCharacterAtEndMatch1(t *testing.T) {
	// Set up
	firstValue := "*a"
	compareValue := "ga"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithFullWildcardAndCharacterAtEndMatch2(t *testing.T) {
	// Set up
	firstValue := "*a"
	compareValue := "gaa"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithFullWildcardAndCharacterAtEndMatch3(t *testing.T) {
	// Set up
	firstValue := "*a"
	compareValue := "gaaa"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithFullWildcardAndCharacterAtEndNoMatch(t *testing.T) {
	// Set up
	firstValue := "*a"
	compareValue := "gb"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithFullWildcardAndCharacterAtStartAndEndMatch1(t *testing.T) {
	// Set up
	firstValue := "a*a"
	compareValue := "aga"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithFullWildcardAndCharacterAtStartAndEndMatch2(t *testing.T) {
	// Set up
	firstValue := "a*a"
	compareValue := "agga"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithFullWildcardAndCharacterAtStartAndEndMatch3(t *testing.T) {
	// Set up
	firstValue := "a*a"
	compareValue := "aggga"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithFullWildcardAndCharacterAtStartAndEndNoMatch(t *testing.T) {
	// Set up
	firstValue := "a*a"
	compareValue := "ag"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithDoubleFullWildcardMatch1(t *testing.T) {
	// Set up
	firstValue := "**"
	compareValue := ""

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithDoubleFullWildcardMatch2(t *testing.T) {
	// Set up
	firstValue := "**"
	compareValue := ""

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithDoubleFullWildcardAndCharacterAtStartMatch1(t *testing.T) {
	// Set up
	firstValue := "a**"
	compareValue := "aa"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithDoubleFullWildcardAndCharacterAtStartMatch2(t *testing.T) {
	// Set up
	firstValue := "a**"
	compareValue := "aaa"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithDoubleFullWildcardAndCharacterAtStartMatch3(t *testing.T) {
	// Set up
	firstValue := "a**"
	compareValue := "aaaa"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithDoubleFullWildcardAndCharacterAtStartNoMatch(t *testing.T) {
	// Set up
	firstValue := "a**"
	compareValue := "b"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithDoubleFullWildcardAndCharacterAtEndMatch1(t *testing.T) {
	// Set up
	firstValue := "**a"
	compareValue := "ga"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithDoubleFullWildcardAndCharacterAtEndMatch2(t *testing.T) {
	// Set up
	firstValue := "**a"
	compareValue := "gga"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithDoubleFullWildcardAndCharacterAtEndMatch3(t *testing.T) {
	// Set up
	firstValue := "**a"
	compareValue := "ggga"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithDoubleFullWildcardAndCharacterAtEndNoMatch(t *testing.T) {
	// Set up
	firstValue := "**a"
	compareValue := "gb"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithDoubleFullWildcardAndCharacterAtStartAndEndMatch1(t *testing.T) {
	// Set up
	firstValue := "a**a"
	compareValue := "aga"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithDoubleFullWildcardAndCharacterAtStartAndEndMatch2(t *testing.T) {
	// Set up
	firstValue := "a**a"
	compareValue := "agga"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithDoubleFullWildcardAndCharacterAtStartAndEndMatch3(t *testing.T) {
	// Set up
	firstValue := "a**a"
	compareValue := "aggga"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithDoubleFullWildcardAndCharacterAtStartAndEndNoMatch(t *testing.T) {
	// Set up
	firstValue := "a**a"
	compareValue := "agb"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithTwoFullWildcardMatch1(t *testing.T) {
	// Set up
	firstValue := "a*b*c"
	compareValue := "abc"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithTwoFullWildcardMatch2(t *testing.T) {
	// Set up
	firstValue := "a*b*c"
	compareValue := "agbc"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithTwoFullWildcardMatch3(t *testing.T) {
	// Set up
	firstValue := "a*b*c"
	compareValue := "agbgc"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithTwoFullWildcardMatch4(t *testing.T) {
	// Set up
	firstValue := "a*b*c"
	compareValue := "aggbggc"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithTwoFullWildcardNoMatch1(t *testing.T) {
	// Set up
	firstValue := "a*b*c"
	compareValue := "Abc"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithTwoFullWildcardNoMatch2(t *testing.T) {
	// Set up
	firstValue := "a*b*c"
	compareValue := "aBc"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithTwoFullWildcardNoMatch3(t *testing.T) {
	// Set up
	firstValue := "a*b*c"
	compareValue := "abC"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithTwoFullWildcardNoMatch4(t *testing.T) {
	// Set up
	firstValue := "a*b*c"
	compareValue := "abcd"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithTwoFullWildcardNoMatch5(t *testing.T) {
	// Set up
	firstValue := "a*b*c"
	compareValue := "Aabc"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithRepeatingPatternMatch1(t *testing.T) {
	// Set up
	firstValue := "a*bc*bc"
	compareValue := "aggggggbcggggggbc"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithRepeatingPatternMatch2(t *testing.T) {
	// Set up
	firstValue := "a*bc*bc"
	compareValue := "aggggggbcggggggbcggggggbc"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithRepeatingPatternMatch3(t *testing.T) {
	// Set up
	firstValue := "a*bc*bc"
	compareValue := "aggggggbhcggggggbcggggggbc"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithRepeatingPatternMatch4(t *testing.T) {
	// Set up
	firstValue := "a*bc*bc"
	compareValue := "aggggggbcggggggbhcggggggbc"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithRepeatingPatternMatch5(t *testing.T) {
	// Set up
	firstValue := "a*?bc*?bc"
	compareValue := "a?bcgbchbc"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func testContainsWithRepeatingPatternNoMatch(t *testing.T) {
	// Set up
	firstValue := "a*bc*bc"
	compareValue := "aggggggbcggggggbcggggggbhc"

	sourceValue := MakePolicyValue(firstValue)

	// Run
	result := sourceValue.Contains(compareValue)

	// Evaluate
	expectedResult := false
	if result != expectedResult {
		t.Fatalf("Expected result: %t but received result: %t", expectedResult, result)
	}
}

func TestPolicyValueIntersectionNoWildcards(t *testing.T) {
	t.Run("TestIntersectionWithoutWildcardAndMatch", testIntersectionWithoutWildcardMatch)
	t.Run("TestIntersectionWithoutWildcardAndNoMatch1", testIntersectionWithoutWildcardNoMatch1)
	t.Run("TestIntersectionWithoutWildcardAndNoMatch2", testIntersectionWithoutWildcardNoMatch2)
	t.Run("TestIntersectionWithoutWildcardAndExtraCharacterMatch", testIntersectionWithoutWildcardAndExtraCharacterMatch)
	t.Run("TestIntersectionWithoutWildcardAndExtraCharacterNoMatch1", testIntersectionWithoutWildcardAndExtraCharacterNoMatch1)
	t.Run("TestIntersectionWithoutWildcardAndExtraCharacterNoMatch2", testIntersectionWithoutWildcardAndExtraCharacterNoMatch2)
}

func testIntersectionWithoutWildcardMatch(t *testing.T) {
	// Set up
	firstValue := "z"
	secondValue := "z"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "z"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithoutWildcardNoMatch1(t *testing.T) {
	// Set up
	firstValue := "z"
	secondValue := "a"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithoutWildcardNoMatch2(t *testing.T) {
	// Set up
	firstValue := "z"
	secondValue := "za"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithoutWildcardAndExtraCharacterMatch(t *testing.T) {
	// Set up
	firstValue := "za"
	secondValue := "za"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "za"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithoutWildcardAndExtraCharacterNoMatch1(t *testing.T) {
	// Set up
	firstValue := "za"
	secondValue := "z"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithoutWildcardAndExtraCharacterNoMatch2(t *testing.T) {
	// Set up
	firstValue := "za"
	secondValue := "zaa"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func TestPolicyValueIntersectionPartialWildcards(t *testing.T) {
	t.Run("TestIntersectionWithPartialWildcardMatch", testIntersectionWithPartialWildcardMatch1)
	t.Run("TestIntersectionWithPartialWildcardMatch", testIntersectionWithPartialWildcardMatch2)
	t.Run("testIntersectionWithPartialWildcardMatch3", testIntersectionWithPartialWildcardMatch3)
	t.Run("testIntersectionWithPartialWildcardMatch4", testIntersectionWithPartialWildcardMatch4)
	t.Run("TestIntersectionWithPartialWildcardNoMatch1", testIntersectionWithPartialWildcardNoMatch1)
	t.Run("TestIntersectionWithPartialWildcardNoMatch2", testIntersectionWithPartialWildcardNoMatch2)
	t.Run("TestIntersectionWithPartialWildcardNoMatch2", testIntersectionWithPartialWildcardNoMatch3)

	t.Run("TestIntersectionWithPartialWildcardAndCharacterAtStartMatch", testIntersectionWithPartialWildcardAndCharacterAtStartMatch)
	t.Run("TestIntersectionWithPartialWildcardAndCharacterAtStartNoMatch1", testIntersectionWithPartialWildcardAndCharacterAtStartNoMatch1)
	t.Run("TestIntersectionWithPartialWildcardAndCharacterAtStartNoMatch2", testIntersectionWithPartialWildcardAndCharacterAtStartNoMatch2)
	t.Run("TestIntersectionWithPartialWildcardAndCharacterAtEndMatch", testIntersectionWithPartialWildcardAndCharacterAtEndMatch)
	t.Run("TestIntersectionWithPartialWildcardAndCharacterAtEndNoMatch1", testIntersectionWithPartialWildcardAndCharacterAtEndNoMatch1)
	t.Run("TestIntersectionWithPartialWildcardAndCharacterAtEndNoMatch2", testIntersectionWithPartialWildcardAndCharacterAtEndNoMatch2)
	t.Run("TestIntersectionWithPartialWildcardAndCharacterAtStartAndEndMatch", testIntersectionWithPartialWildcardAndCharacterAtStartAndEndMatch)
	t.Run("TestIntersectionWithPartialWildcardAndCharacterAtStartAndEndNoMatch1", testIntersectionWithPartialWildcardAndCharacterAtStartAndEndNoMatch1)
	t.Run("TestIntersectionWithPartialWildcardAndCharacterAtStartAndEndNoMatch2", testIntersectionWithPartialWildcardAndCharacterAtStartAndEndNoMatch2)

	t.Run("TestIntersectionWithDoublePartialWildcardMatch", testIntersectionWithDoublePartialWildcardMatch)
	t.Run("TestIntersectionWithDoublePartialWildcardNoMatch1", testIntersectionWithDoublePartialWildcardNoMatch1)
	t.Run("TestIntersectionWithDoublePartialWildcardNoMatch2", testIntersectionWithDoublePartialWildcardNoMatch2)
	t.Run("TestIntersectionWithDoublePartialWildcardAndCharacterAtStartMatch", testIntersectionWithDoublePartialWildcardAndCharacterAtStartMatch)
	t.Run("TestIntersectionWithDoublePartialWildcardAndCharacterAtStartNoMatch1", testIntersectionWithDoublePartialWildcardAndCharacterAtStartNoMatch1)
	t.Run("TestIntersectionWithDoublePartialWildcardAndCharacterAtStartNoMatch2", testIntersectionWithDoublePartialWildcardAndCharacterAtStartNoMatch2)
	t.Run("TestIntersectionWithDoublePartialWildcardAndCharacterAtEndMatch", testIntersectionWithDoublePartialWildcardAndCharacterAtEndMatch)
	t.Run("TestIntersectionWithDoublePartialWildcardAndCharacterAtEndNoMatch1", testIntersectionWithDoublePartialWildcardAndCharacterAtEndNoMatch1)
	t.Run("TestIntersectionWithDoublePartialWildcardAndCharacterAtEndNoMatch2", testIntersectionWithDoublePartialWildcardAndCharacterAtEndNoMatch2)
	t.Run("TestIntersectionWithDoublePartialWildcardAndCharacterAtStartAndEndMatch", testIntersectionWithDoublePartialWildcardAndCharacterAtStartAndEndMatch)
	t.Run("TestIntersectionWithDoublePartialWildcardAndCharacterAtStartAndEndNoMatch1", testIntersectionWithDoublePartialWildcardAndCharacterAtStartAndEndNoMatch1)
	t.Run("TestIntersectionWithDoublePartialWildcardAndCharacterAtStartAndEndNoMatch2", testIntersectionWithDoublePartialWildcardAndCharacterAtStartAndEndNoMatch2)
}

func testIntersectionWithPartialWildcardMatch1(t *testing.T) {
	// Set up
	firstValue := "?"
	secondValue := "z"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "z"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithPartialWildcardMatch2(t *testing.T) {
	// Set up
	firstValue := "?"
	secondValue := "?"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "?"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithPartialWildcardMatch3(t *testing.T) {
	// Set up
	firstValue := "?"
	secondValue := "?*"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "?"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithPartialWildcardMatch4(t *testing.T) {
	// Set up
	firstValue := "?"
	secondValue := "*"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "?"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithPartialWildcardNoMatch1(t *testing.T) {
	// Set up
	firstValue := "?"
	secondValue := ""

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithPartialWildcardNoMatch2(t *testing.T) {
	// Set up
	firstValue := "?"
	secondValue := "zz"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithPartialWildcardNoMatch3(t *testing.T) {
	// Set up
	firstValue := "?"
	secondValue := "??"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithPartialWildcardAndCharacterAtStartMatch(t *testing.T) {
	// Set up
	firstValue := "a?"
	secondValue := "az"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "az"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithPartialWildcardAndCharacterAtStartNoMatch1(t *testing.T) {
	// Set up
	firstValue := "a?"
	secondValue := "a"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithPartialWildcardAndCharacterAtStartNoMatch2(t *testing.T) {
	// Set up
	firstValue := "a?"
	secondValue := "a"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithPartialWildcardAndCharacterAtEndMatch(t *testing.T) {
	// Set up
	firstValue := "?z"
	secondValue := "az"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "az"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithPartialWildcardAndCharacterAtEndNoMatch1(t *testing.T) {
	// Set up
	firstValue := "?z"
	secondValue := "azz"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithPartialWildcardAndCharacterAtEndNoMatch2(t *testing.T) {
	// Set up
	firstValue := "?z"
	secondValue := "a"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithPartialWildcardAndCharacterAtStartAndEndMatch(t *testing.T) {
	// Set up
	firstValue := "a?z"
	secondValue := "agz"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "agz"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithPartialWildcardAndCharacterAtStartAndEndNoMatch1(t *testing.T) {
	// Set up
	firstValue := "a?z"
	secondValue := "ag"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithPartialWildcardAndCharacterAtStartAndEndNoMatch2(t *testing.T) {
	// Set up
	firstValue := "a?z"
	secondValue := "agzg"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithDoublePartialWildcardMatch(t *testing.T) {
	// Set up
	firstValue := "??"
	secondValue := "zz"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "zz"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithDoublePartialWildcardNoMatch1(t *testing.T) {
	// Set up
	firstValue := "??"
	secondValue := "z"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithDoublePartialWildcardNoMatch2(t *testing.T) {
	// Set up
	firstValue := "??"
	secondValue := "zzz"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithDoublePartialWildcardAndCharacterAtStartMatch(t *testing.T) {
	// Set up
	firstValue := "a??"
	secondValue := "azz"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "azz"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithDoublePartialWildcardAndCharacterAtStartNoMatch1(t *testing.T) {
	// Set up
	firstValue := "a??"
	secondValue := "az"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithDoublePartialWildcardAndCharacterAtStartNoMatch2(t *testing.T) {
	// Set up
	firstValue := "a??"
	secondValue := "azzz"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithDoublePartialWildcardAndCharacterAtEndMatch(t *testing.T) {
	// Set up
	firstValue := "??z"
	secondValue := "ggz"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "ggz"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithDoublePartialWildcardAndCharacterAtEndNoMatch1(t *testing.T) {
	// Set up
	firstValue := "??z"
	secondValue := "ggza"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithDoublePartialWildcardAndCharacterAtEndNoMatch2(t *testing.T) {
	// Set up
	firstValue := "??z"
	secondValue := "gg"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithDoublePartialWildcardAndCharacterAtStartAndEndMatch(t *testing.T) {
	// Set up
	firstValue := "a??z"
	secondValue := "aggz"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "aggz"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithDoublePartialWildcardAndCharacterAtStartAndEndNoMatch1(t *testing.T) {
	// Set up
	firstValue := "a??z"
	secondValue := "agg"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithDoublePartialWildcardAndCharacterAtStartAndEndNoMatch2(t *testing.T) {
	// Set up
	firstValue := "a??z"
	secondValue := "aggzg"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func TestPolicyValueIntersectionFullWildcards(t *testing.T) {
	t.Run("TestIntersectionWithFullWildcardMatch1", testIntersectionWithFullWildcardMatch1)
	t.Run("TestIntersectionWithFullWildcardMatch2", testIntersectionWithFullWildcardMatch2)
	t.Run("TestIntersectionWithFullWildcardMatch3", testIntersectionWithFullWildcardMatch3)
	t.Run("TestIntersectionWithFullWildcardMatch4", testIntersectionWithFullWildcardMatch4)

	t.Run("TestIntersectionWithFullWildcardAndCharacterAtStartMatch1", testIntersectionWithFullWildcardAndCharacterAtStartMatch1)
	t.Run("TestIntersectionWithFullWildcardAndCharacterAtStartMatch2", testIntersectionWithFullWildcardAndCharacterAtStartMatch2)
	t.Run("TestIntersectionWithFullWildcardAndCharacterAtStartMatch3", testIntersectionWithFullWildcardAndCharacterAtStartMatch3)
	t.Run("TestIntersectionWithFullWildcardAndCharacterAtStartNoMatch", testIntersectionWithFullWildcardAndCharacterAtStartNoMatch)
	t.Run("TestIntersectionWithFullWildcardAndCharacterAtStartNoMatch", testIntersectionWithFullWildcardAndCharacterAtStartNoMatch)
	t.Run("TestIntersectionWithFullWildcardAndCharacterAtEndMatch1", testIntersectionWithFullWildcardAndCharacterAtEndMatch1)
	t.Run("TestIntersectionWithFullWildcardAndCharacterAtEndMatch2", testIntersectionWithFullWildcardAndCharacterAtEndMatch2)
	t.Run("TestIntersectionWithFullWildcardAndCharacterAtEndMatch3", testIntersectionWithFullWildcardAndCharacterAtEndMatch3)
	t.Run("TestIntersectionWithFullWildcardAndCharacterAtEndNoMatch", testIntersectionWithFullWildcardAndCharacterAtEndNoMatch)
	t.Run("TestIntersectionWithFullWildcardAndCharacterAtStartAndEndMatch1", testIntersectionWithFullWildcardAndCharacterAtStartAndEndMatch1)
	t.Run("TestIntersectionWithFullWildcardAndCharacterAtStartAndEndMatch2", testIntersectionWithFullWildcardAndCharacterAtStartAndEndMatch2)
	t.Run("TestIntersectionWithFullWildcardAndCharacterAtStartAndEndMatch3", testIntersectionWithFullWildcardAndCharacterAtStartAndEndMatch3)
	t.Run("TestIntersectionWithFullWildcardAndCharacterAtStartAndEndNoMatch", testIntersectionWithFullWildcardAndCharacterAtStartAndEndNoMatch)

	t.Run("TestIntersectionWithDoubleFullWildcardMatch1", testIntersectionWithDoubleFullWildcardMatch1)
	t.Run("TestIntersectionWithDoubleFullWildcardMatch2", testIntersectionWithDoubleFullWildcardMatch2)
	t.Run("TestIntersectionWithDoubleFullWildcardAndCharacterAtStartMatch1", testIntersectionWithDoubleFullWildcardAndCharacterAtStartMatch1)
	t.Run("TestIntersectionWithDoubleFullWildcardAndCharacterAtStartMatch2", testIntersectionWithDoubleFullWildcardAndCharacterAtStartMatch2)
	t.Run("TestIntersectionWithDoubleFullWildcardAndCharacterAtStartMatch3", testIntersectionWithDoubleFullWildcardAndCharacterAtStartMatch3)
	t.Run("TestIntersectionWithDoubleFullWildcardAndCharacterAtStartNoMatch", testIntersectionWithDoubleFullWildcardAndCharacterAtStartNoMatch)
	t.Run("TestIntersectionWithDoubleFullWildcardAndCharacterAtEndMatch1", testIntersectionWithDoubleFullWildcardAndCharacterAtEndMatch1)
	t.Run("TestIntersectionWithDoubleFullWildcardAndCharacterAtEndMatch2", testIntersectionWithDoubleFullWildcardAndCharacterAtEndMatch2)
	t.Run("TestIntersectionWithDoubleFullWildcardAndCharacterAtEndMatch3", testIntersectionWithDoubleFullWildcardAndCharacterAtEndMatch3)
	t.Run("TestIntersectionWithDoubleFullWildcardAndCharacterAtEndNoMatch", testIntersectionWithDoubleFullWildcardAndCharacterAtEndNoMatch)
	t.Run("TestIntersectionWithDoubleFullWildcardAndCharacterAtStartAndEndMatch1", testIntersectionWithDoubleFullWildcardAndCharacterAtStartAndEndMatch1)
	t.Run("TestIntersectionWithDoubleFullWildcardAndCharacterAtStartAndEndMatch2", testIntersectionWithDoubleFullWildcardAndCharacterAtStartAndEndMatch2)
	t.Run("TestIntersectionWithDoubleFullWildcardAndCharacterAtStartAndEndMatch3", testIntersectionWithDoubleFullWildcardAndCharacterAtStartAndEndMatch3)
	t.Run("TestIntersectionWithDoubleFullWildcardAndCharacterAtStartAndEndNoMatch", testIntersectionWithDoubleFullWildcardAndCharacterAtStartAndEndNoMatch)

	t.Run("TestIntersectionWithTwoFullWildcardMatch1", testIntersectionWithTwoFullWildcardMatch1)
	t.Run("TestIntersectionWithTwoFullWildcardMatch2", testIntersectionWithTwoFullWildcardMatch2)
	t.Run("TestIntersectionWithTwoFullWildcardMatch3", testIntersectionWithTwoFullWildcardMatch3)
	t.Run("TestIntersectionWithTwoFullWildcardMatch4", testIntersectionWithTwoFullWildcardMatch4)
	t.Run("TestIntersectionWithTwoFullWildcardNoMatch1", testIntersectionWithTwoFullWildcardNoMatch1)
	t.Run("TestIntersectionWithTwoFullWildcardNoMatch2", testIntersectionWithTwoFullWildcardNoMatch2)
	t.Run("TestIntersectionWithTwoFullWildcardNoMatch3", testIntersectionWithTwoFullWildcardNoMatch3)
	t.Run("TestIntersectionWithTwoFullWildcardNoMatch4", testIntersectionWithTwoFullWildcardNoMatch4)
	t.Run("TestIntersectionWithTwoFullWildcardNoMatch5", testIntersectionWithTwoFullWildcardNoMatch5)

	t.Run("TestIntersectionWithRepeatingPatternMatch1", testIntersectionWithRepeatingPatternMatch1)
	t.Run("TestIntersectionWithRepeatingPatternMatch2", testIntersectionWithRepeatingPatternMatch2)
	t.Run("TestIntersectionWithRepeatingPatternMatch3", testIntersectionWithRepeatingPatternMatch3)
	t.Run("TestIntersectionWithRepeatingPatternMatch4", testIntersectionWithRepeatingPatternMatch4)
	t.Run("TestIntersectionWithRepeatingPatternMatch5", testIntersectionWithRepeatingPatternMatch5)

	t.Run("TestIntersectionWithRepeatingPatternNoMatch", testIntersectionWithRepeatingPatternNoMatch)
}

func testIntersectionWithFullWildcardMatch1(t *testing.T) {
	// Set up
	firstValue := "*"
	secondValue := ""

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithFullWildcardMatch2(t *testing.T) {
	// Set up
	firstValue := "*"
	secondValue := "z"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "z"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithFullWildcardMatch3(t *testing.T) {
	// Set up
	firstValue := "*"
	secondValue := "?"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "?"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithFullWildcardMatch4(t *testing.T) {
	// Set up
	firstValue := "*"
	secondValue := "*"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "*"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithFullWildcardAndCharacterAtStartMatch1(t *testing.T) {
	// Set up
	firstValue := "a*"
	secondValue := "a"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "a"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithFullWildcardAndCharacterAtStartMatch2(t *testing.T) {
	// Set up
	firstValue := "a*"
	secondValue := "aa"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "aa"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithFullWildcardAndCharacterAtStartMatch3(t *testing.T) {
	// Set up
	firstValue := "a*"
	secondValue := "aaa"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "aaa"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithFullWildcardAndCharacterAtStartNoMatch(t *testing.T) {
	// Set up
	firstValue := "a*"
	secondValue := "b"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithFullWildcardAndCharacterAtEndMatch1(t *testing.T) {
	// Set up
	firstValue := "*a"
	secondValue := "ga"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "ga"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithFullWildcardAndCharacterAtEndMatch2(t *testing.T) {
	// Set up
	firstValue := "*a"
	secondValue := "gaa"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "gaa"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithFullWildcardAndCharacterAtEndMatch3(t *testing.T) {
	// Set up
	firstValue := "*a"
	secondValue := "gaaa"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "gaaa"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithFullWildcardAndCharacterAtEndNoMatch(t *testing.T) {
	// Set up
	firstValue := "*a"
	secondValue := "gb"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithFullWildcardAndCharacterAtStartAndEndMatch1(t *testing.T) {
	// Set up
	firstValue := "a*a"
	secondValue := "aga"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "aga"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithFullWildcardAndCharacterAtStartAndEndMatch2(t *testing.T) {
	// Set up
	firstValue := "a*a"
	secondValue := "agga"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "agga"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithFullWildcardAndCharacterAtStartAndEndMatch3(t *testing.T) {
	// Set up
	firstValue := "a*a"
	secondValue := "aggga"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "aggga"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithFullWildcardAndCharacterAtStartAndEndNoMatch(t *testing.T) {
	// Set up
	firstValue := "a*a"
	secondValue := "ag"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithDoubleFullWildcardMatch1(t *testing.T) {
	// Set up
	firstValue := "**"
	secondValue := ""

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithDoubleFullWildcardMatch2(t *testing.T) {
	// Set up
	firstValue := "**"
	secondValue := ""

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithDoubleFullWildcardAndCharacterAtStartMatch1(t *testing.T) {
	// Set up
	firstValue := "a**"
	secondValue := "aa"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "aa"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithDoubleFullWildcardAndCharacterAtStartMatch2(t *testing.T) {
	// Set up
	firstValue := "a**"
	secondValue := "aaa"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "aaa"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithDoubleFullWildcardAndCharacterAtStartMatch3(t *testing.T) {
	// Set up
	firstValue := "a**"
	secondValue := "aaaa"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "aaaa"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithDoubleFullWildcardAndCharacterAtStartNoMatch(t *testing.T) {
	// Set up
	firstValue := "a**"
	secondValue := "b"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithDoubleFullWildcardAndCharacterAtEndMatch1(t *testing.T) {
	// Set up
	firstValue := "**a"
	secondValue := "ga"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "ga"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithDoubleFullWildcardAndCharacterAtEndMatch2(t *testing.T) {
	// Set up
	firstValue := "**a"
	secondValue := "gga"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "gga"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithDoubleFullWildcardAndCharacterAtEndMatch3(t *testing.T) {
	// Set up
	firstValue := "**a"
	secondValue := "ggga"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "ggga"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithDoubleFullWildcardAndCharacterAtEndNoMatch(t *testing.T) {
	// Set up
	firstValue := "**a"
	secondValue := "gb"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithDoubleFullWildcardAndCharacterAtStartAndEndMatch1(t *testing.T) {
	// Set up
	firstValue := "a**a"
	secondValue := "aga"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "aga"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithDoubleFullWildcardAndCharacterAtStartAndEndMatch2(t *testing.T) {
	// Set up
	firstValue := "a**a"
	secondValue := "agga"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "agga"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithDoubleFullWildcardAndCharacterAtStartAndEndMatch3(t *testing.T) {
	// Set up
	firstValue := "a**a"
	secondValue := "aggga"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "aggga"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithDoubleFullWildcardAndCharacterAtStartAndEndNoMatch(t *testing.T) {
	// Set up
	firstValue := "a**a"
	secondValue := "agb"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithTwoFullWildcardMatch1(t *testing.T) {
	// Set up
	firstValue := "a*b*c"
	secondValue := "abc"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "abc"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithTwoFullWildcardMatch2(t *testing.T) {
	// Set up
	firstValue := "a*b*c"
	secondValue := "agbc"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "agbc"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithTwoFullWildcardMatch3(t *testing.T) {
	// Set up
	firstValue := "a*b*c"
	secondValue := "agbgc"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "agbgc"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithTwoFullWildcardMatch4(t *testing.T) {
	// Set up
	firstValue := "a*b*c"
	secondValue := "aggbggc"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "aggbggc"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithTwoFullWildcardNoMatch1(t *testing.T) {
	// Set up
	firstValue := "a*b*c"
	secondValue := "Abc"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithTwoFullWildcardNoMatch2(t *testing.T) {
	// Set up
	firstValue := "a*b*c"
	secondValue := "aBc"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithTwoFullWildcardNoMatch3(t *testing.T) {
	// Set up
	firstValue := "a*b*c"
	secondValue := "abC"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithTwoFullWildcardNoMatch4(t *testing.T) {
	// Set up
	firstValue := "a*b*c"
	secondValue := "abcd"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithTwoFullWildcardNoMatch5(t *testing.T) {
	// Set up
	firstValue := "a*b*c"
	secondValue := "Aabc"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithRepeatingPatternMatch1(t *testing.T) {
	// Set up
	firstValue := "a*bc*bc"
	secondValue := "aggggggbcggggggbc"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "aggggggbcggggggbc"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithRepeatingPatternMatch2(t *testing.T) {
	// Set up
	firstValue := "a*bc*bc"
	secondValue := "aggggggbcggggggbcggggggbc"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "aggggggbcggggggbcggggggbc"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithRepeatingPatternMatch3(t *testing.T) {
	// Set up
	firstValue := "a*bc*bc"
	secondValue := "aggggggbhcggggggbcggggggbc"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "aggggggbhcggggggbcggggggbc"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithRepeatingPatternMatch4(t *testing.T) {
	// Set up
	firstValue := "a*bc*bc"
	secondValue := "aggggggbcggggggbhcggggggbc"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "aggggggbcggggggbhcggggggbc"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithRepeatingPatternMatch5(t *testing.T) {
	// Set up
	firstValue := "a*?bc*?bc"
	secondValue := "a?bcgbchbc"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := "a?bcgbchbc"
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}

func testIntersectionWithRepeatingPatternNoMatch(t *testing.T) {
	// Set up
	firstValue := "a*bc*bc"
	secondValue := "aggggggbcggggggbcggggggbhc"

	firstPolicyValue := MakePolicyValue(firstValue)
	secondPolicyValue := MakePolicyValue(secondValue)

	// Run
	result1 := firstPolicyValue.Intersection(secondPolicyValue)
	result2 := secondPolicyValue.Intersection(firstPolicyValue)

	// Evaluate
	expectedResult := ""
	if result1 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result1)
	}

	if result2 != expectedResult {
		t.Fatalf("Expected result: %s but received result: %s", expectedResult, result2)
	}
}
