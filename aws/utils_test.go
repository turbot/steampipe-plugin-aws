package aws

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"testing"
)

func Test_calculateMaxLimit(t *testing.T) {
	type args[T interface{ ~int | ~int64 | ~int32 }] struct {
		maxItems T
		d        *plugin.QueryData
	}
	type testCase[T interface{ ~int | ~int64 | ~int32 }] struct {
		name string
		args args[T]
		want T
	}
	var limit int64 = 100
	tests := []testCase[int64]{
		{
			name: "larger than 100",
			args: args[int64]{
				maxItems: 200,
				d: &plugin.QueryData{
					QueryContext: &plugin.QueryContext{
						Limit: &limit,
					},
				},
			},
			want: 100,
		},
		{
			name: "less than 100",
			args: args[int64]{
				maxItems: 50,
				d: &plugin.QueryData{
					QueryContext: &plugin.QueryContext{
						Limit: &limit,
					},
				},
			},
			want: 50,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateMaxLimit(tt.args.maxItems, tt.args.d); got != tt.want {
				t.Errorf("calculateMaxLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}
