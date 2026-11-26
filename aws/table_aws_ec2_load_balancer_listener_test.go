package aws

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
)

func TestListenerAttributesMap(t *testing.T) {
	got := listenerAttributesMap([]types.ListenerAttribute{
		{Key: aws.String("tcp.idle_timeout.seconds"), Value: aws.String("350")},
		{Key: aws.String("routing.http.response.server.enabled"), Value: aws.String("true")},
		{Key: aws.String("empty.value")},
		{Value: aws.String("no key")},
	})
	want := map[string]string{
		"tcp.idle_timeout.seconds":             "350",
		"routing.http.response.server.enabled": "true",
		"empty.value":                          "",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("listenerAttributesMap() = %v, want %v", got, want)
	}
	if got := listenerAttributesMap(nil); len(got) != 0 {
		t.Errorf("listenerAttributesMap(nil) = %v, want empty", got)
	}
}
