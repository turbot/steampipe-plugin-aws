package aws

import (
	"errors"
	"fmt"
	"testing"

	"github.com/aws/smithy-go"
)

func TestIsLocalEKSClusterUnsupportedError(t *testing.T) {
	localClusterErr := &smithy.GenericAPIError{
		Code:    "ClientException",
		Message: "This operation is not supported for local EKS clusters",
	}

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"local cluster", localClusterErr, true},
		{"wrapped local cluster", fmt.Errorf("operation error EKS: ListNodegroups: %w", localClusterErr), true},
		{"other client exception", &smithy.GenericAPIError{Code: "ClientException", Message: "Cluster status is not ACTIVE"}, false},
		{"access denied", &smithy.GenericAPIError{Code: "AccessDeniedException", Message: "This operation is not supported for local EKS clusters"}, false},
		{"not an API error", errors.New("This operation is not supported for local EKS clusters"), false},
		{"nil", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isLocalEKSClusterUnsupportedError(tt.err); got != tt.want {
				t.Fatalf("isLocalEKSClusterUnsupportedError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}
