package aws

type numberConversionTest struct {
	name     string
	input    string
	expected string
}

// func TestGetConditionalKeymapping(t *testing.T) {
// 	testCases := []numberConversionTest{
// 		{
// 			`any_value`,
// 			`{"Version":"2012-10-17","Statement":[{"Sid":"","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::111122223333:root"},"Action":"sts:AssumeRole","Condition":{"ForAnyValue:StringEquals":{"sts:ExternalId":["abc:cde","efg:hij"]}}}]}`,
// 			`{"sts:externalid":{"ForAnyValue:StringEquals":{"Key":"sts:externalid","Value":["efg:hij","abc:cde"]}}}`,
// 		},
// 		{
// 			`single_value`,
// 			`{"Version":"2012-10-17","Statement":[{"Sid":"","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::111122223333:root"},"Action":"sts:AssumeRole","Condition":{"StringEquals":{"aws:PrincipalAccount":"123456789012","aws:PrincipalOrgID":"o-123456789"}}}]}`,
// 			`{"aws:principalaccount":{"StringEquals":{"Key":"aws:principalaccount","Value":["123456789012"]}},"aws:principalorgid":{"StringEquals":{"Key":"aws:principalorgid","Value":["o-123456789"]}}}`,
// 		},
// 	}

// 	for _, test := range testCases {

// 		pol, err := canonicalPolicy(test.input)
// 		if err != nil {
// 			t.Errorf("Convert failed for case '%s': %v", pol, err)
// 		}

// 		condition := (pol.(Policy)).Statements[0].Condition

// 		newCondition, err := getConditionKeyMapping(&condition)
// 		if err != nil {
// 			t.Errorf("Failed to get condition key mapping '%s': %v", condition, err)
// 		}

// 		var input Conditions
// 		_ = json.Unmarshal([]byte(test.expected), &input)
// 		// strdata, _ := json.MarshalIndent(input, "", "\t")
// 		// output, _ := json.MarshalIndent(newCondition, "", "\t")
// 		// fmt.Printf("%s\n", string(strdata))

// 		if !reflect.DeepEqual(input, *newCondition) {
// 			t.Errorf("\nTest: '%s.%s' FAILED\nexpected:\n %v \ngot:\n %v \n", "TestConvertStatementCondition", test.name, input, *newCondition)
// 		}
// 		fmt.Printf("\nTest: '%s.%s' PASSED\noutput:\n %v\n", "TestConvertStatementCondition", test.name, input)
// 	}
// }
