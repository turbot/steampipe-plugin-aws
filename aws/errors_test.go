package aws

import (
	"errors"
	"fmt"
	"testing"

	"github.com/aws/smithy-go"
)

func TestHasAPIErrorCode(t *testing.T) {
	notFound := &smithy.GenericAPIError{Code: "DBSnapshotNotFound", Message: "DBSnapshot snap-1 not found."}
	cases := []struct {
		name  string
		err   error
		codes []string
		want  bool
	}{
		{"matching code", notFound, []string{"DBSnapshotNotFound"}, true},
		{"one of several codes", notFound, []string{"DBClusterSnapshotNotFoundFault", "DBSnapshotNotFound"}, true},
		{"wrapped error", fmt.Errorf("operation error RDS: %w", notFound), []string{"DBSnapshotNotFound"}, true},
		{"other API error", &smithy.GenericAPIError{Code: "AccessDenied"}, []string{"DBSnapshotNotFound"}, false},
		{"EBS snapshot not found", &smithy.GenericAPIError{Code: "InvalidSnapshot.NotFound"}, []string{"InvalidSnapshot.NotFound"}, true},
		{"prefix only does not match", &smithy.GenericAPIError{Code: "InvalidSnapshot.NotFoundSomething"}, []string{"InvalidSnapshot.NotFound"}, false},
		{"not an API error", errors.New("connection reset"), []string{"DBSnapshotNotFound"}, false},
		{"nil error", nil, []string{"DBSnapshotNotFound"}, false},
	}
	for _, c := range cases {
		if got := hasAPIErrorCode(c.err, c.codes...); got != c.want {
			t.Errorf("%s: hasAPIErrorCode() = %v, want %v", c.name, got, c.want)
		}
	}
}
