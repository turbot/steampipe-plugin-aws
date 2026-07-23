package aws

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	smithy "github.com/aws/smithy-go"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Build an error chain matching what the AWS SDK produces when a request
// fails at the transport level: OperationError -> RequestSendError ->
// url.Error -> net.OpError.
func transportError(op string, inner error) error {
	return &smithy.OperationError{
		ServiceID:     "EC2",
		OperationName: "DescribeSecurityGroups",
		Err: &smithyhttp.RequestSendError{
			Err: &url.Error{
				Op:  "Post",
				URL: "https://ec2.me-south-1.amazonaws.com/",
				Err: &net.OpError{
					Op:  op,
					Net: "tcp",
					Err: inner,
				},
			},
		},
	}
}

type timeoutError struct{}

func (timeoutError) Error() string   { return "i/o timeout" }
func (timeoutError) Timeout() bool   { return true }
func (timeoutError) Temporary() bool { return true }

func TestIsDialErrorRetryable(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		expected aws.Ternary
	}{
		{"dial i/o timeout", transportError("dial", timeoutError{}), aws.FalseTernary},
		{"dial connection refused", transportError("dial", os.NewSyscallError("connect", syscall.ECONNREFUSED)), aws.FalseTernary},
		{"dial no route to host", transportError("dial", os.NewSyscallError("connect", syscall.EHOSTUNREACH)), aws.FalseTernary},
		{"dial network unreachable", transportError("dial", os.NewSyscallError("connect", syscall.ENETUNREACH)), aws.FalseTernary},
		{"ephemeral port exhaustion stays retryable", transportError("dial", os.NewSyscallError("connect", syscall.EADDRNOTAVAIL)), aws.UnknownTernary},
		{"fd exhaustion stays retryable", transportError("dial", os.NewSyscallError("socket", syscall.EMFILE)), aws.UnknownTernary},
		{"system fd exhaustion stays retryable", transportError("dial", os.NewSyscallError("socket", syscall.ENFILE)), aws.UnknownTernary},
		{"read on established connection", transportError("read", syscall.ECONNRESET), aws.UnknownTernary},
		{"context cancellation is not a send error", &smithy.CanceledError{Err: context.Canceled}, aws.UnknownTernary},
		{"api error without transport failure", fmt.Errorf("api error Throttling: rate exceeded"), aws.UnknownTernary},
		{"nil-adjacent plain error", fmt.Errorf("something else"), aws.UnknownTernary},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := isDialErrorRetryable(tc.err)
			if got != tc.expected {
				t.Errorf("isDialErrorRetryable(%v) = %v, want %v", tc.err, got, tc.expected)
			}
		})
	}
}

func TestSharedHTTPClientDialTimeout(t *testing.T) {
	client, ok := sharedHTTPClient.(interface{ GetDialer() *net.Dialer })
	if !ok {
		t.Fatalf("sharedHTTPClient does not expose GetDialer")
	}
	if got := client.GetDialer().Timeout; got != 10*time.Second {
		t.Errorf("dialer timeout = %v, want %v", got, 10*time.Second)
	}
}
