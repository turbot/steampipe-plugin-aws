package aws

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/aws/aws-sdk-go-v2/service/accessanalyzer/types"
)

func main() {
	originalFields := reflect.TypeOf(types.AnalyzerSummary{})

	customColumns := tableAwsAccessAnalyzer(context.Background()).Columns

	customFieldsMap := make(map[string]bool)
	for _, column := range customColumns {
		customFieldsMap[column.Name] = true
	}

	for i := 0; i < originalFields.NumField(); i++ {
		field := originalFields.Field(i)
		if !customFieldsMap[toSnakeCase(field.Name)] && field.Name != "noSmithyDocumentSerde" {
			fmt.Printf("Missing field in custom format: %s\n", field.Name)
		}
	}
}

// toSnakeCase converts a string from CamelCase to snake_case.
func toSnakeCase(str string) string {
	// Add a space before all caps
	var spaceBuffer strings.Builder
	for i, rune := range str {
		if unicode.IsUpper(rune) && i != 0 {
			spaceBuffer.WriteRune(' ')
		}
		spaceBuffer.WriteRune(rune)
	}

	// Replace spaces with underscores and lower-case all letters
	snake := strings.ToLower(strings.ReplaceAll(spaceBuffer.String(), " ", "_"))

	// Handle special case where consecutive upper case letters are present
	re := regexp.MustCompile(`_([a-z])_([a-z])`)
	snake = re.ReplaceAllString(snake, "_$1$2")

	return snake
}
