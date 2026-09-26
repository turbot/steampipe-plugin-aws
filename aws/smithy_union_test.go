package aws

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentcorecontrol/types"
)

func toJSON(t *testing.T, v interface{}) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return string(b)
}

func TestFlattenSmithyUnionsNil(t *testing.T) {
	if got := flattenSmithyUnions(nil); got != nil {
		t.Fatalf("nil input: expected nil, got %#v", got)
	}
	var nilUnion types.AgentRuntimeArtifact
	if got := flattenSmithyUnions(nilUnion); got != nil {
		t.Fatalf("nil interface: expected nil, got %#v", got)
	}
	var nilPtr *types.NetworkConfiguration
	if got := flattenSmithyUnions(nilPtr); got != nil {
		t.Fatalf("nil pointer: expected nil, got %#v", got)
	}
}

func TestFlattenSmithyUnionsNestedUnion(t *testing.T) {
	var artifact types.AgentRuntimeArtifact = &types.AgentRuntimeArtifactMemberCodeConfiguration{
		Value: types.CodeConfiguration{
			Code: &types.CodeMemberS3{
				Value: types.S3Location{
					Bucket: aws.String("my-bucket"),
					Prefix: aws.String("agent.zip"),
				},
			},
			EntryPoint: []string{"main.py"},
			Runtime:    types.AgentManagedRuntimeTypePython312,
		},
	}

	got := toJSON(t, flattenSmithyUnions(artifact))
	want := `{"CodeConfiguration":{"Code":{"S3":{"Bucket":"my-bucket","Prefix":"agent.zip","VersionId":null}},"EntryPoint":["main.py"],"Runtime":"PYTHON_3_12"}}`
	if got != want {
		t.Fatalf("\n got: %s\nwant: %s", got, want)
	}
}

func TestFlattenSmithyUnionsContainerVariant(t *testing.T) {
	var artifact types.AgentRuntimeArtifact = &types.AgentRuntimeArtifactMemberContainerConfiguration{
		Value: types.ContainerConfiguration{ContainerUri: aws.String("123.dkr.ecr.us-east-1.amazonaws.com/agent:latest")},
	}

	got := toJSON(t, flattenSmithyUnions(artifact))
	want := `{"ContainerConfiguration":{"ContainerUri":"123.dkr.ecr.us-east-1.amazonaws.com/agent:latest"}}`
	if got != want {
		t.Fatalf("\n got: %s\nwant: %s", got, want)
	}
}

func TestFlattenSmithyUnionsAuthorizer(t *testing.T) {
	var auth types.AuthorizerConfiguration = &types.AuthorizerConfigurationMemberCustomJWTAuthorizer{
		Value: types.CustomJWTAuthorizerConfiguration{
			DiscoveryUrl:   aws.String("https://issuer.example/.well-known/openid-configuration"),
			AllowedClients: []string{"client-a"},
		},
	}

	flat, ok := flattenSmithyUnions(auth).(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", flattenSmithyUnions(auth))
	}
	inner, ok := flat["CustomJWTAuthorizer"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected CustomJWTAuthorizer key, got %v", flat)
	}
	if inner["DiscoveryUrl"] != "https://issuer.example/.well-known/openid-configuration" {
		t.Fatalf("DiscoveryUrl not preserved: %v", inner["DiscoveryUrl"])
	}
	if _, present := inner["Value"]; present {
		t.Fatalf("Value wrapper should have been removed: %v", inner)
	}
}

func TestFlattenSmithyUnionsPlainStructUnchanged(t *testing.T) {
	// Structs that are not union members keep their field names, and nested
	// pointers, maps, enums and times survive.
	ts := time.Date(2026, 9, 25, 12, 0, 0, 0, time.UTC)
	in := struct {
		Network *types.NetworkConfiguration
		Env     map[string]string
		When    time.Time
		Raw     []byte
	}{
		Network: &types.NetworkConfiguration{
			NetworkMode:       types.NetworkModeVpc,
			NetworkModeConfig: &types.VpcConfig{Subnets: []string{"subnet-1"}, SecurityGroups: []string{"sg-1"}},
		},
		Env:  map[string]string{"A": "1"},
		When: ts,
		Raw:  []byte("hi"),
	}

	got := toJSON(t, flattenSmithyUnions(in))
	want := `{"Env":{"A":"1"},"Network":{"NetworkMode":"VPC","NetworkModeConfig":{"RequireServiceS3Endpoint":null,"SecurityGroups":["sg-1"],"Subnets":["subnet-1"]}},"Raw":"aGk=","When":"2026-09-25T12:00:00Z"}`
	if got != want {
		t.Fatalf("\n got: %s\nwant: %s", got, want)
	}
}

func TestSmithyUnionVariantRejectsLookalikes(t *testing.T) {
	// Name matches the pattern but has two exported fields: not a union member.
	type FooMemberBar struct {
		Value string
		Other string
	}
	// Name matches but the single field is not called Value.
	type FooMemberBaz struct {
		Thing string
	}
	for _, v := range []interface{}{FooMemberBar{Value: "x"}, FooMemberBaz{Thing: "y"}} {
		flat, ok := flattenSmithyUnions(v).(map[string]interface{})
		if !ok {
			t.Fatalf("expected map for %T", v)
		}
		if _, isUnion := flat["Bar"]; isUnion {
			t.Fatalf("%T should not be treated as a union: %v", v, flat)
		}
		if _, isUnion := flat["Baz"]; isUnion {
			t.Fatalf("%T should not be treated as a union: %v", v, flat)
		}
	}
}

func TestFlattenSmithyUnionsListOfUnions(t *testing.T) {
	fs := []types.FilesystemConfiguration{
		&types.FilesystemConfigurationMemberSessionStorage{Value: types.SessionStorageConfiguration{MountPath: aws.String("/mnt/session")}},
		&types.FilesystemConfigurationMemberEfsAccessPoint{Value: types.EfsAccessPointConfiguration{AccessPointArn: aws.String("arn:aws:elasticfilesystem:us-east-1:123456789012:access-point/fsap-0abc")}},
		&types.FilesystemConfigurationMemberS3FilesAccessPoint{Value: types.S3FilesAccessPointConfiguration{AccessPointArn: aws.String("arn:aws:s3:us-east-1:123456789012:accesspoint/files")}},
	}

	got := toJSON(t, flattenSmithyUnions(fs))
	want := `[{"SessionStorage":{"MountPath":"/mnt/session"}},{"EfsAccessPoint":{"AccessPointArn":"arn:aws:elasticfilesystem:us-east-1:123456789012:access-point/fsap-0abc","MountPath":null}},{"S3FilesAccessPoint":{"AccessPointArn":"arn:aws:s3:us-east-1:123456789012:accesspoint/files","MountPath":null}}]`
	if got != want {
		t.Fatalf("\n got: %s\nwant: %s", got, want)
	}

	var empty []types.FilesystemConfiguration
	if got := flattenSmithyUnions(empty); got != nil {
		t.Fatalf("nil slice: expected nil, got %#v", got)
	}
}

func TestFlattenSmithyUnionsRequestHeaderAllowlist(t *testing.T) {
	var cfg types.RequestHeaderConfiguration = &types.RequestHeaderConfigurationMemberRequestHeaderAllowlist{Value: []string{"Authorization", "X-Tenant"}}
	got := toJSON(t, flattenSmithyUnions(cfg))
	want := `{"RequestHeaderAllowlist":["Authorization","X-Tenant"]}`
	if got != want {
		t.Fatalf("\n got: %s\nwant: %s", got, want)
	}
}
