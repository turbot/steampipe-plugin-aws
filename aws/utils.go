package aws

import (
	"context"
	"encoding/base64"
	"fmt"
	"math"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/sagemaker"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func ec2TagsToMap(tags []*ec2.Tag) (*map[string]string, error) {
	var turbotTagsMap map[string]string
	if tags == nil {
		return nil, nil
	}

	turbotTagsMap = map[string]string{}
	for _, i := range tags {
		turbotTagsMap[*i.Key] = *i.Value
	}

	return &turbotTagsMap, nil
}

func arnToAkas(_ context.Context, d *transform.TransformData) (interface{}, error) {
	arn := types.SafeString(d.Value)
	return []string{arn}, nil
}

func arnToTitle(_ context.Context, d *transform.TransformData) (interface{}, error) {
	arn := types.SafeString(d.Value)

	// get the resource title
	title := arn[strings.LastIndex(arn, ":")+1:]

	return title, nil
}

func convertTimestamp(_ context.Context, d *transform.TransformData) (interface{}, error) {
	epochTime := d.Value.(*int64)

	if epochTime != nil {
		timeInSec := math.Floor(float64(*epochTime) / 1000)
		unixTimestamp := time.Unix(int64(timeInSec), 0)
		timestampRFC3339Format := unixTimestamp.Format(time.RFC3339)
		return timestampRFC3339Format, nil
	}
	return nil, nil
}

func extractNameFromSqsQueueURL(queue string) (string, error) {
	//http://sqs.us-west-2.amazonaws.com/123456789012/queueName
	u, err := url.Parse(queue)
	if err != nil {
		return "", err
	}
	segments := strings.Split(u.Path, "/")
	if len(segments) != 3 {
		return "", fmt.Errorf("SQS Url not parsed correctly")
	}

	return segments[2], nil
}

func handleNilString(_ context.Context, d *transform.TransformData) (interface{}, error) {
	value := types.SafeString(d.Value)
	if value == "" {
		return "false", nil
	}
	return value, nil
}

func rowSourceFromIPPermission(group *ec2.SecurityGroup, permission *ec2.IpPermission, groupType string) []interface{} {
	var rowSource []interface{}

	// create 1 row per ip-range
	if permission.IpRanges != nil {
		for _, r := range permission.IpRanges {
			rowSource = append(rowSource, &vpcSecurityGroupRulesRowData{
				Group:           group,
				Permission:      permission,
				IPRange:         r,
				Ipv6Range:       nil,
				UserIDGroupPair: nil,
				Type:            groupType,
			})
		}
	}

	// create 1 row per ipv6-range
	if permission.Ipv6Ranges != nil {
		for _, r := range permission.Ipv6Ranges {
			rowSource = append(rowSource, &vpcSecurityGroupRulesRowData{
				Group:           group,
				Permission:      permission,
				IPRange:         nil,
				Ipv6Range:       r,
				UserIDGroupPair: nil,
				Type:            groupType,
			})
		}
	}

	// create 1 row per user id group pair
	if permission.UserIdGroupPairs != nil {
		for _, r := range permission.UserIdGroupPairs {
			rowSource = append(rowSource, &vpcSecurityGroupRulesRowData{
				Group:           group,
				Permission:      permission,
				IPRange:         nil,
				Ipv6Range:       nil,
				UserIDGroupPair: r,
				Type:            groupType,
			})
		}
	}

	return rowSource
}

func resourceInterfaceDescription(key string) string {
	switch key {
	case "akas":
		return "Array of globally unique identifier strings (also known as) for the resource."
	case "tags":
		return "A map of tags for the resource."
	case "title":
		return "Title of the resource."
	}
	return ""
}

func getLastPathElement(path string) string {
	if path == "" {
		return ""
	}
	pathItems := strings.Split(path, "/")
	return pathItems[len(pathItems)-1]
}

func lastPathElement(_ context.Context, d *transform.TransformData) (interface{}, error) {
	return getLastPathElement(types.SafeString(d.Value)), nil
}

func base64DecodedData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data, err := base64.StdEncoding.DecodeString(types.SafeString(d.Value))
	// check if CorruptInputError or invalid UTF-8
	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instancedata-add-user-data.html
	if err != nil {
		return nil, nil
	} else if !utf8.Valid(data) {
		return types.SafeString(d.Value), nil
	}
	return data, nil
}

// Transform function for sagemaker resources tags
func sageMakerTurbotTags(_ context.Context, d *transform.TransformData) (interface{},
	error) {
	tags := d.HydrateItem.([]*sagemaker.Tag)

	if tags != nil {
		turbotTagsMap := map[string]string{}
		for _, i := range tags {
			turbotTagsMap[*i.Key] = *i.Value
		}
		return turbotTagsMap, nil
	}

	return nil, nil
}
