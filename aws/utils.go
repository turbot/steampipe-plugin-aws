package aws

import (
	"context"
	"encoding/base64"
	"fmt"
	"math"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/sagemaker"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
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
	value := types.SafeString(fmt.Sprintf("%v", d.Value))
	if value == "" {
		return "false", nil
	}
	return value, nil
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

func getQualsValueByColumn(equalQuals plugin.KeyColumnQualMap, columnName string, dataType string) interface{} {
	var value interface{}
	for _, q := range equalQuals[columnName].Quals {
		if dataType == "string" {
			if q.Value.GetStringValue() != "" {
				value = q.Value.GetStringValue()
			} else {
				valList := getListValues(q.Value.GetListValue())
				if len(valList) > 0 {
					value = valList
				}
			}
		}
		if dataType == "boolean" {
			switch q.Operator {
			case "<>":
				value = "false"
			case "=":
				value = "true"
			}
		}
		if dataType == "int64" {
			value = q.Value.GetInt64Value()
			if q.Value.GetInt64Value() == 0 {
				valueSlice := make([]*string, 0)
				for _, value := range q.Value.GetListValue().Values {
					val := strconv.FormatInt(value.GetInt64Value(), 10)
					valueSlice = append(valueSlice, &val)
				}
				value = valueSlice
			}
		}
		if dataType == "double" {
			value = q.Value.GetDoubleValue()
			if q.Value.GetDoubleValue() == 0 {
				valueSlice := make([]*string, 0)
				for _, value := range q.Value.GetListValue().Values {
					val := strconv.FormatFloat(value.GetDoubleValue(), 'f', 4, 64)
					valueSlice = append(valueSlice, &val)
				}
				value = valueSlice
			}

		}
		if dataType == "ipaddr" {
			value = q.Value.GetInetValue().Addr
			if q.Value.GetInetValue().Addr == "" {
				valueSlice := make([]*string, 0)
				for _, value := range q.Value.GetListValue().Values {
					val := value.GetInetValue().Addr
					valueSlice = append(valueSlice, &val)
				}
				value = valueSlice
			}
		}
		if dataType == "cidr" {
			value = q.Value.GetInetValue().Cidr
			if q.Value.GetInetValue().Addr == "" {
				valueSlice := make([]*string, 0)
				for _, value := range q.Value.GetListValue().Values {
					val := value.GetInetValue().Cidr
					valueSlice = append(valueSlice, &val)
				}
				value = valueSlice
			}
		}
		if dataType == "time" {
			value = getListValues(q.Value.GetListValue())
			if len(getListValues(q.Value.GetListValue())) == 0 {
				value = q.Value.GetTimestampValue().AsTime()
			}
		}
	}
	return value
}

// handleNullIfZero :: handles empty slices and map convert them to null instead of the zero type
func handleEmptySliceAndMap(ctx context.Context, d *transform.TransformData) (any, error) {
	if d.Value == nil {
		return nil, nil
	}

	reflectVal := reflect.ValueOf(d.Value)
	switch reflectVal.Kind() {
	case reflect.Slice, reflect.Map:
		if reflectVal.Len() == 0 {
			return nil, nil
		}
	case reflect.Struct:
		if reflectVal == reflect.Zero(reflectVal.Type()) {
			return nil, nil
		}
	}

	return d.Value, nil
}
