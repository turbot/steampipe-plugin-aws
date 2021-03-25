package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/wafv2"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsWafv2IpSet(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_wafv2_ip_set",
		Description: "AWS WAFv2 IP Set",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"id", "name", "scope"}),
			ShouldIgnoreError: isNotFoundError([]string{}),
			Hydrate:           getAwsWafv2IpSet,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsWafv2IpSets,
		},
		GetMatrixItem: BuildRegionList,
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the IP set.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the entity.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ARN"),
			},
			{
				Name:        "id",
				Description: "A unique identifier for the IP set.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "scope",
				Description: "Specifies the scope of the IP Set. Possible values are: 'REGIONAL' and 'CLOUDFRONT'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(ipSetLocation),
			},
			{
				Name:        "description",
				Description: "A description of the IP set that helps with identification.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ip_address_version",
				Description: "Specifies the IP address type. Possible values are: 'IPV4' and 'IPV6'.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsWafv2IpSet,
				Transform:   transform.FromField("IPAddressVersion"),
			},
			{
				Name:        "lock_token",
				Description: "A token used for optimistic locking.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "addresses",
				Description: "An array of strings that specify one or more IP addresses or blocks of IP addresses in Classless Inter-Domain Routing (CIDR) notation.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsWafv2IpSet,
			},
			// {
			// 	Name:        "tags_src",
			// 	Description: "A list of tags associated with the resource.",
			// 	Type:        proto.ColumnType_JSON,
			// 	Hydrate:     listTagsForAwsWafv2WebAcl,
			// 	Transform:   transform.FromField("TagInfoForResource.TagList"),
			// },

			// steampipe standard columns
			// {
			// 	Name:        "title",
			// 	Description: resourceInterfaceDescription("title"),
			// 	Type:        proto.ColumnType_STRING,
			// 	Transform:   transform.FromField("Name"),
			// },
			// {
			// 	Name:        "tags",
			// 	Description: resourceInterfaceDescription("tags"),
			// 	Type:        proto.ColumnType_JSON,
			// 	Hydrate:     listTagsForAwsWafv2WebAcl,
			// 	Transform:   transform.FromField("TagInfoForResource.TagList").Transform(webAclTagListToTurbotTags),
			// },
			// {
			// 	Name:        "akas",
			// 	Description: resourceInterfaceDescription("akas"),
			// 	Type:        proto.ColumnType_JSON,
			// 	Transform:   transform.FromField("ARN").Transform(arnToAkas),
			// },

			// aws standard columns
			{
				Name:        "partition",
				Description: "The AWS partition in which the resource is located (aws, aws-cn, or aws-us-gov).",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCommonColumns,
			},
			{
				Name:        "region",
				Description: "The AWS Region in which the resource is located.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(ipSetRegion),
			},
			{
				Name:        "account_id",
				Description: "The AWS Account ID in which the resource is located.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getCommonColumns,
			},
		},
	}
}

//// LIST FUNCTION

func listAwsWafv2IpSets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listAwsWafv2IpSets", "AWS_REGION", region)

	// Create session
	svc, err := WAFv2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List all regional web acls
	pagesLeft := true
	regionalIPSetParams := &wafv2.ListIPSetsInput{
		Scope: aws.String("REGIONAL"),
	}
	for pagesLeft {
		response, err := svc.ListIPSets(regionalIPSetParams)
		if err != nil {
			return nil, err
		}

		for _, regionalIPSets := range response.IPSets {
			d.StreamListItem(ctx, regionalIPSets)
		}

		if response.NextMarker != nil {
			pagesLeft = true
			regionalIPSetParams.NextMarker = response.NextMarker
		} else {
			pagesLeft = false
		}
	}

	// List all global web acls
	// To work with CloudFront, you must specify the Region US East (N. Virginia)
	if region == "us-east-1" {
		pagesLeft = true
		globalIPSetParams := &wafv2.ListIPSetsInput{
			Scope: aws.String("CLOUDFRONT"),
		}
		for pagesLeft {
			response, err := svc.ListIPSets(globalIPSetParams)
			if err != nil {
				return nil, err
			}

			for _, globalIPSets := range response.IPSets {
				d.StreamListItem(ctx, globalIPSets)
			}

			if response.NextMarker != nil {
				pagesLeft = true
				globalIPSetParams.NextMarker = response.NextMarker
			} else {
				pagesLeft = false
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAwsWafv2IpSet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsWafv2IpSet")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	var id, name, scope string
	if h.Item != nil {
		data := ipSetData(h.Item)
		id = data["ID"]
		name = data["Name"]
		locationType := strings.Split(strings.Split(string(data["Arn"]), ":")[5], "/")[0]

		if locationType == "regional" {
			scope = "REGIONAL"
		} else {
			scope = "CLOUDFRONT"
		}
	} else {
		id = d.KeyColumnQuals["id"].GetStringValue()
		name = d.KeyColumnQuals["name"].GetStringValue()
		scope = d.KeyColumnQuals["scope"].GetStringValue()
	}

	// Create Session
	svc, err := WAFv2Service(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// To work with CloudFront, you must specify the Region US East (N. Virginia)
	if strings.ToLower(scope) == "cloudfront" && region != "us-east-1" {
		return nil, nil
	}

	params := &wafv2.GetIPSetInput{
		Id:    aws.String(id),
		Name:  aws.String(name),
		Scope: aws.String(scope),
	}

	op, err := svc.GetIPSet(params)
	if err != nil {
		plugin.Logger(ctx).Debug("GetIPSet", "ERROR", err)
		return nil, err
	}

	return op.IPSet, nil
}

// func listTagsForAwsWafv2WebAcl(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 	plugin.Logger(ctx).Trace("listTagsForAwsWafv2WebAcl")

// 	// TODO put me in helper function
// 	var region string
// 	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
// 	if matrixRegion != nil {
// 		region = matrixRegion.(string)
// 	}
// 	data := webAclData(h.Item)
// 	locationType := strings.Split(strings.Split(string(data["Arn"]), ":")[5], "/")[0]

// 	// To work with CloudFront, you must specify the Region US East (N. Virginia)
// 	if locationType == "global" && region != "us-east-1" {
// 		return nil, nil
// 	}

// 	// Create session
// 	svc, err := WAFv2Service(ctx, d, region)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Build param
// 	param := &wafv2.ListTagsForResourceInput{
// 		ResourceARN: aws.String(data["Arn"]),
// 	}

// 	webAclTags, err := svc.ListTagsForResource(param)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return webAclTags, nil
// }

//// TRANSFORM FUNCTIONS

func ipSetLocation(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := ipSetData(d.HydrateItem)
	loc := strings.Split(strings.Split(string(data["Arn"]), ":")[5], "/")[0]
	if loc == "regional" {
		return "REGIONAL", nil
	}
	return "CLOUDFRONT", nil
}

// func webAclTagListToTurbotTags(ctx context.Context, d *transform.TransformData) (interface{}, error) {
// 	plugin.Logger(ctx).Trace("webAclTagListToTurbotTags")
// 	data := d.HydrateItem.(*wafv2.ListTagsForResourceOutput)

// 	if data.TagInfoForResource.TagList == nil || len(data.TagInfoForResource.TagList) < 1 {
// 		return nil, nil
// 	}

// 	// Mapping the resource tags inside turbotTags
// 	var turbotTagsMap map[string]string
// 	if data.TagInfoForResource.TagList != nil {
// 		turbotTagsMap = map[string]string{}
// 		for _, i := range data.TagInfoForResource.TagList {
// 			turbotTagsMap[*i.Key] = *i.Value
// 		}
// 	}

// 	return turbotTagsMap, nil
// }

func ipSetRegion(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := ipSetData(d.HydrateItem)
	loc := strings.Split(strings.Split(string(data["Arn"]), ":")[5], "/")[0]

	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	if loc == "global" {
		return "global", nil
	}
	return region, nil
}

func ipSetData(item interface{}) map[string]string {
	data := map[string]string{}
	switch item := item.(type) {
	case *wafv2.IPSet:
		data["ID"] = *item.Id
		data["Arn"] = *item.ARN
		data["Name"] = *item.Name
		data["Description"] = *item.Description
	case *wafv2.IPSetSummary:
		data["ID"] = *item.Id
		data["Arn"] = *item.ARN
		data["Name"] = *item.Name
		data["Description"] = *item.Description
	}
	return data
}
