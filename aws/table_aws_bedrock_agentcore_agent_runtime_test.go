package aws

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentcorecontrol"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentcorecontrol/types"
)

func TestIsValidBedrockAgentCoreAgentRuntimeId(t *testing.T) {
	testCases := []struct {
		name  string
		id    string
		valid bool
	}{
		{"real id", "steampipe_test_public-H3VjKs7uYz", true},
		{"single char name", "a-0123456789", true},
		{"max length name", strings.Repeat("a", 100) + "-0123456789", true},
		{"name too long", strings.Repeat("a", 101) + "-0123456789", false},
		{"empty", "", false},
		{"no suffix", "steampipe_test_public", false},
		{"suffix too short", "agent-abc123XYZ", false},
		{"suffix too long", "agent-abc123XYZ00", false},
		{"suffix with underscore", "agent-abc123XYZ_", false},
		{"name starts with digit", "1agent-abc123XYZ0", false},
		{"name starts with underscore", "_agent-abc123XYZ0", false},
		{"hyphen in name", "my-agent-abc123XYZ0", false},
		{"spaces", "bad id with spaces!", false},
		{"arn instead of id", "arn:aws:bedrock-agentcore:us-east-1:123456789012:runtime/agent-abc123XYZ0", false},
		{"unicode", "agént-abc123XYZ0", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isValidBedrockAgentCoreAgentRuntimeId(tc.id); got != tc.valid {
				t.Fatalf("isValidBedrockAgentCoreAgentRuntimeId(%q) = %v, want %v", tc.id, got, tc.valid)
			}
		})
	}
}

func TestBedrockAgentCoreAgentRuntimeArn(t *testing.T) {
	const arn = "arn:aws:bedrock-agentcore:us-east-1:123456789012:runtime/agent-abc123XYZ0"

	if got := bedrockAgentCoreAgentRuntimeArn(types.AgentRuntime{AgentRuntimeArn: aws.String(arn)}); got != arn {
		t.Fatalf("list item: got %q", got)
	}
	if got := bedrockAgentCoreAgentRuntimeArn(&bedrockagentcorecontrol.GetAgentRuntimeOutput{AgentRuntimeArn: aws.String(arn)}); got != arn {
		t.Fatalf("get output: got %q", got)
	}
	if got := bedrockAgentCoreAgentRuntimeArn(types.AgentRuntime{}); got != "" {
		t.Fatalf("nil arn on list item: got %q", got)
	}
	if got := bedrockAgentCoreAgentRuntimeArn(nil); got != "" {
		t.Fatalf("nil item: got %q", got)
	}
	if got := bedrockAgentCoreAgentRuntimeArn("unexpected"); got != "" {
		t.Fatalf("unexpected type: got %q", got)
	}
}
