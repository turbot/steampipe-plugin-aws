package aws

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codeconnections"
	"github.com/aws/aws-sdk-go-v2/service/codeconnections/types"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

func TestCodeConnectionsTurbotTags(t *testing.T) {
	tests := []struct {
		name string
		tags []types.Tag
		want interface{}
	}{
		{name: "omitted tags return SQL null"},
		{name: "empty tags return SQL null", tags: []types.Tag{}},
		{
			name: "populated tags preserve keys and empty values",
			tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
				{Key: aws.String("Owner"), Value: aws.String("")},
			},
			want: map[string]string{"Environment": "test", "Owner": ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := codeConnectionsTurbotTags(context.Background(), &transform.TransformData{
				HydrateItem: &codeconnections.ListTagsForResourceOutput{Tags: tt.tags},
			})
			if err != nil {
				t.Fatal(err)
			}
			// An actual nil interface is required for the SDK's SQL NULL path.
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}
