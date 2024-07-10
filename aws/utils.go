package aws

import (
	"context"
	"encoding/base64"
	"fmt"
	"math"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	sagemakerTypes "github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

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

// isControlCharacter checks if a character is a control character
func isControlCharacter(r rune) bool {
	return r < 32 || r == 127
}

// isPrintable checks if a character is printable.
func isPrintable(r rune) bool {
	return unicode.IsPrint(r) && !isControlCharacter(r)
}


// cleanString removes control characters and ensures valid UTF-8 encoding
func cleanString(data string) string {
	var builder strings.Builder
	for _, r := range data {
		// Only add valid and printable runes to the builder
		if utf8.ValidRune(r) && isPrintable(r) {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func base64DecodedData(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data, err := base64.StdEncoding.DecodeString(types.SafeString(d.Value))
	// check if CorruptInputError or invalid UTF-8
	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instancedata-add-user-data.html
	if err != nil {
		return nil, nil
	} 
	// Check for valid UTF-8
	if !utf8.Valid(data) {
		return types.SafeString(d.Value), nil
	}
	// Convert byte slice to string
	decodedString := string(data)
	// Clean the string by removing invalid UTF-8 and control characters
	processedData := cleanString(decodedString)
	return processedData, nil
}

// Transform function for sagemaker resources tags
func sageMakerTurbotTags(_ context.Context, d *transform.TransformData) (interface{},
	error) {
	tags := d.HydrateItem.([]sagemakerTypes.Tag)

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

func isAWSARN(s string) bool {

	// Define the AWS ARN pattern
	arnPattern := `^arn:(aws|aws-cn|aws-us-gov):[a-z0-9_-]+:[a-z0-9_-]+:\d+:([a-zA-Z0-9/-]+)?$`
	regex := regexp.MustCompile(arnPattern)

	// Test if the string matches the ARN pattern
	return regex.MatchString(s)
}

// The function ensures that the given engine is is not of type `docdb` or `neptune`
// You can get the available engines list from - https://docs.aws.amazon.com/cli/latest/reference/rds/describe-db-engine-versions.html
// Current supported RDS engine values as of 2023/08/07
func isSuppportedRDSEngine(engine string) bool {
	supportedEngines := []string{
		"aurora",
		"aurora-mysql",
		"aurora-postgresql",
		"custom-sqlserver-ee",
		"custom-sqlserver-se",
		"custom-sqlserver-web",
		"mariadb",
		"mysql",
		"oracle-ee",
		"oracle-ee-cdb",
		"oracle-se2",
		"oracle-se2-cdb",
		"postgres",
		"sqlserver-ee",
		"sqlserver-ex",
		"sqlserver-se",
		"sqlserver-web",
	}

	return helpers.StringSliceContains(supportedEngines, engine)
}

// Helper function for integer based environment variables.
func readEnvVarToInt(name string, defaultVal int) int {
	val := defaultVal
	envValue := os.Getenv(name)
	if envValue != "" {
		i, err := strconv.Atoi(envValue)
		if err == nil {
			val = i
		}
	}
	return val
}
