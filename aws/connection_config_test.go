package aws

import (
	"reflect"
	"testing"

	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

func connectionWithConfig(config interface{}) *plugin.Connection {
	c := &plugin.Connection{Name: "test"}
	c.SetConfig(config)
	return c
}

func TestGetConfigNilConnection(t *testing.T) {
	got := GetConfig(nil)
	if !reflect.DeepEqual(got, awsConfig{}) {
		t.Errorf("GetConfig(nil) = %+v, want empty config", got)
	}
}

func TestGetConfigNilRawConfig(t *testing.T) {
	got := GetConfig(&plugin.Connection{Name: "test"})
	if !reflect.DeepEqual(got, awsConfig{}) {
		t.Errorf("GetConfig with nil raw config = %+v, want empty config", got)
	}
}

func TestGetConfigWrongType(t *testing.T) {
	got := GetConfig(connectionWithConfig("not an awsConfig"))
	if !reflect.DeepEqual(got, awsConfig{}) {
		t.Errorf("GetConfig with wrong config type = %+v, want empty config", got)
	}
}

func TestGetConfigNormalizesRegions(t *testing.T) {
	testCases := []struct {
		name     string
		regions  []string
		expected []string
	}{
		{"lower case unchanged", []string{"us-east-1", "eu-*"}, []string{"us-east-1", "eu-*"}},
		{"mixed case lowered", []string{"US-EAST-1", "Eu-West-2"}, []string{"us-east-1", "eu-west-2"}},
		{"wildcard unchanged", []string{"*"}, []string{"*"}},
		{"exclusion pattern lowered", []string{"*", "!ME-SOUTH-1"}, []string{"*", "!me-south-1"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := GetConfig(connectionWithConfig(awsConfig{Regions: tc.regions}))
			if !reflect.DeepEqual(got.Regions, tc.expected) {
				t.Errorf("GetConfig regions = %v, want %v", got.Regions, tc.expected)
			}
		})
	}
}

func TestGetConfigEmptyRegionsPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("GetConfig with regions = [] did not panic")
		}
	}()
	GetConfig(connectionWithConfig(awsConfig{Regions: []string{}}))
}

// The remaining config args should pass through GetConfig untouched.
func TestGetConfigPassesThroughOtherArgs(t *testing.T) {
	defaultRegion := "eu-west-2"
	profile := "myprofile"
	accessKey := "AKIA0000000000000000"
	secretKey := "secret"
	sessionToken := "token"
	maxRetries := 5
	minDelay := 50
	endpointUrl := "http://localhost:4566"
	forcePathStyle := true

	in := awsConfig{
		Regions:               []string{"us-east-1"},
		DefaultRegion:         &defaultRegion,
		Profile:               &profile,
		AccessKey:             &accessKey,
		SecretKey:             &secretKey,
		SessionToken:          &sessionToken,
		MaxErrorRetryAttempts: &maxRetries,
		MinErrorRetryDelay:    &minDelay,
		IgnoreErrorMessages:   []string{".*timeout.*"},
		IgnoreErrorCodes:      []string{"AccessDenied", "UnauthorizedOperation"},
		EndpointUrl:           &endpointUrl,
		S3ForcePathStyle:      &forcePathStyle,
	}

	got := GetConfig(connectionWithConfig(in))
	if !reflect.DeepEqual(got, in) {
		t.Errorf("GetConfig = %+v, want %+v", got, in)
	}
}

func TestNormalizeRegion(t *testing.T) {
	testCases := []struct {
		in       string
		expected string
	}{
		{"us-east-1", "us-east-1"},
		{"US-EAST-1", "us-east-1"},
		{"Us-Gov-West-1", "us-gov-west-1"},
		{"!ME-SOUTH-1", "!me-south-1"},
		{"*", "*"},
		{"", ""},
	}

	for _, tc := range testCases {
		if got := NormalizeRegion(tc.in); got != tc.expected {
			t.Errorf("NormalizeRegion(%q) = %q, want %q", tc.in, got, tc.expected)
		}
	}
}
