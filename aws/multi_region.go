package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/iancoleman/strcase"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

// MultiRegionGet :: helper function to execute list call for all regions
func MultiRegionList(listFunc plugin.HydrateFunc) plugin.HydrateFunc {
	// get regions from the connection config
	regions := getRegions()

	return func(ctx context.Context, queryData *plugin.QueryData, hydrateData *plugin.HydrateData) (interface{}, error) {
		// build a list of parameter maps - one per region
		paramsList := make([]map[string]string, len(regions))
		for i, region := range regions {
			paramsList[i] = map[string]string{"region": region}
		}
		return plugin.ListForEach(ctx, queryData, hydrateData, listFunc, paramsList)
	}
}

// MultiRegionGet :: helper function to execute get for all regions
func MultiRegionGet(getFunc plugin.HydrateFunc) plugin.HydrateFunc {
	// get regions from the connection config
	regions := getRegions()

	return func(ctx context.Context, queryData *plugin.QueryData, hydrateData *plugin.HydrateData) (interface{}, error) {
		// build a list of parameter maps - one per region
		paramsList := make([]map[string]string, len(regions))
		for i, region := range regions {
			paramsList[i] = map[string]string{"region": region}
		}
		return plugin.GetForEach(ctx, queryData, hydrateData, getFunc, paramsList)
	}
}

// dummy function to get regions - this wil be replace with connection config read
func getRegions() []string {
	// TODO read the actual conneciton config
	return []string{"us-east-1", "eu-west-1", "eu-west-2"}
}

// WrappedItem should be used to wrap the hydrated items so that they always include the region.
// Use the 'unwrap' transforms to unwrap from this struct into the columns
type WrappedItem struct {
	Region string
	Item   interface{}
}

func wrapItem(item interface{}, region string) *WrappedItem {
	return &WrappedItem{
		Region: region,
		Item:   item,
	}
}

// TRANSFORMS
func unwrapItem(ctx context.Context, d *transform.TransformData) (interface{}, error) {

	item, ok := d.HydrateItem.(*WrappedItem)
	if ok {
		return item.Item, nil
	}
	return item, nil
}

func unwrapFromCamel(_ context.Context, d *transform.TransformData) (interface{}, error) {
	unwrapped := d.HydrateItem

	item, ok := d.HydrateItem.(*WrappedItem)
	if ok {
		unwrapped = item.Item
	}

	propertyPath := strcase.ToCamel(d.ColumnName)

	fieldValue, ok := helpers.GetNestedFieldValueFromInterface(unwrapped, propertyPath)
	if !ok {
		return nil, fmt.Errorf("Failed to retrieve property path %s\n", propertyPath)
	}

	return fieldValue, nil
}

func unwrapFromField(_ context.Context, d *transform.TransformData) (interface{}, error) {
	unwrapped := d.HydrateItem

	item, ok := d.HydrateItem.(*WrappedItem)
	if ok {
		unwrapped = item.Item
	}

	propertyPath := d.Param.(string)

	fieldValue, ok := helpers.GetNestedFieldValueFromInterface(unwrapped, propertyPath)
	if !ok {
		return nil, fmt.Errorf("Failed to retrieve property path %s\n", propertyPath)
	}

	return fieldValue, nil
}

// utility
func errorIsNotFound(ctx context.Context, err error) bool {
	plugin.Logger(ctx).Trace("errorIsNotFound", "err", err)

	if awsErr, ok := err.(awserr.Error); ok {
		plugin.Logger(ctx).Trace("errorIsNotFound", "Code", awsErr.Code())
		return strings.HasSuffix(awsErr.Code(), ".NotFound")
	}
	return false
}
