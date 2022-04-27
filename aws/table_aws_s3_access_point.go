package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3control"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsS3AccessPoint(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_s3_access_point",
		Description: "AWS S3 Access Point",
		List: &plugin.ListConfig{
			Hydrate: listS3AccessPoints,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "bucket_name", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "region"}),
			Hydrate:    getS3AccessPoint,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundErrorWithContext([]string{"NoSuchAccessPoint", "AccessDenied", "InvalidParameter"}),
			},
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "Specifies the name of the access point.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "access_point_arn",
				Description: "Amazon Resource Name (ARN) of the access point.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAccessPointArn,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "bucket_name",
				Description: "The name of the bucket associated with this access point.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Bucket"),
			},
			{
				Name:        "access_point_policy_is_public",
				Description: "Indicates whether this access point policy is public, or not.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getS3AccessPointPolicyStatus,
				Transform:   transform.FromField("PolicyStatus.IsPublic"),
				Default:     false,
			},
			{
				Name:        "block_public_acls",
				Description: "Specifies whether Amazon S3 should block public access control lists (ACLs) for buckets in this account.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getS3AccessPoint,
				Transform:   transform.FromField("PublicAccessBlockConfiguration.BlockPublicAcls"),
			},
			{
				Name:        "block_public_policy",
				Description: "Specifies whether Amazon S3 should block public bucket policies for buckets in this account.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getS3AccessPoint,
				Transform:   transform.FromField("PublicAccessBlockConfiguration.BlockPublicPolicy"),
			},
			{
				Name:        "ignore_public_acls",
				Description: "Specifies whether Amazon S3 should ignore public ACLs for buckets in this account.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getS3AccessPoint,
				Transform:   transform.FromField("PublicAccessBlockConfiguration.IgnorePublicAcls"),
			},
			{
				Name:        "restrict_public_buckets",
				Description: "Specifies whether Amazon S3 should restrict public bucket policies for buckets in this account.",
				Type:        proto.ColumnType_BOOL,
				Hydrate:     getS3AccessPoint,
				Transform:   transform.FromField("PublicAccessBlockConfiguration.RestrictPublicBuckets"),
			},
			{
				Name:        "creation_date",
				Description: "The date and time when the specified access point was created.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getS3AccessPoint,
			},
			{
				Name:        "network_origin",
				Description: "Indicates whether this access point allows access from the public internet. If VpcConfiguration is specified for this access point, then NetworkOrigin is VPC, and the access point doesn't allow access from the public internet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "vpc_id",
				Description: "Specifies the VPC ID from which the access point will only allow connections.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("VpcConfiguration.VpcId"),
			},
			{
				Name:        "policy",
				Description: "The access point policy associated with the specified access point.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getS3AccessPointPolicy,
				Transform:   transform.FromField("Policy").Transform(transform.UnmarshalYAML),
			},
			{
				Name:        "policy_std",
				Description: "Contains the policy in a canonical form for easier searching.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getS3AccessPointPolicy,
				Transform:   transform.FromField("Policy").Transform(policyToCanonical),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAccessPointArn,
				Transform:   transform.FromValue().Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listS3AccessPoints(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listS3AccessPoints")
	region := d.KeyColumnQualString(matrixKeyRegion)

	// Get account details
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Create Session
	svc, err := S3ControlService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	input := &s3control.ListAccessPointsInput{
		AccountId:  aws.String(commonColumnData.AccountId),
		MaxResults: aws.Int64(1000),
	}

	equalQuals := d.KeyColumnQuals
	if equalQuals["bucket_name"] != nil {
		if equalQuals["bucket_name"].GetStringValue() != "" {
			input.Bucket = aws.String(equalQuals["bucket_name"].GetStringValue())
		}
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.MaxResults {
			// Minimum limit is 0 as per the doc https://docs.aws.amazon.com/AmazonS3/latest/API/API_control_ListAccessPoints.html#API_control_ListAccessPoints_RequestSyntax
			input.MaxResults = limit
		}
	}

	err = svc.ListAccessPointsPages(
		input,
		func(page *s3control.ListAccessPointsOutput, isLast bool) bool {
			for _, accessPoint := range page.AccessPointList {
				d.StreamListItem(ctx, accessPoint)

				// Context may get cancelled due to manual cancellation or if the limit has been reached
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return !isLast
		},
	)
	return nil, err
}

//// HYDRATE FUNCTIONS

func getS3AccessPoint(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getS3AccessPoint")
	matrixRegion := d.KeyColumnQualString(matrixKeyRegion)

	// Get account details
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Create Session
	svc, err := S3ControlService(ctx, d, matrixRegion)
	if err != nil {
		return nil, err
	}

	var name, region string
	if h.Item != nil {
		name = *h.Item.(*s3control.AccessPoint).Name
		region = matrixRegion
	} else {
		name = d.KeyColumnQuals["name"].GetStringValue()
		region = d.KeyColumnQuals["region"].GetStringValue()
	}

	// Return nil, if given region doesn't match config region
	if region != matrixRegion {
		return nil, nil
	}

	// Build params
	params := &s3control.GetAccessPointInput{
		Name:      aws.String(name),
		AccountId: aws.String(commonColumnData.AccountId),
	}

	// execute list call
	item, err := svc.GetAccessPoint(params)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func getS3AccessPointPolicyStatus(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getS3AccessPointPolicyStatus")
	region := d.KeyColumnQualString(matrixKeyRegion)

	// Get account details
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Create Session
	svc, err := S3ControlService(ctx, d, region)
	if err != nil {
		return nil, err
	}
	accessPointName := accessPointName(h.Item)

	// Build params
	params := &s3control.GetAccessPointPolicyStatusInput{
		Name:      aws.String(accessPointName),
		AccountId: aws.String(commonColumnData.AccountId),
	}

	// execute list call
	op, err := svc.GetAccessPointPolicyStatus(params)
	if err != nil {
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "NoSuchAccessPointPolicy" {
				return nil, nil
			}
		}
		return nil, err
	}

	return op, nil
}

func getS3AccessPointPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getS3AccessPointPolicy")
	region := d.KeyColumnQualString(matrixKeyRegion)

	// Get account details
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Create Session
	svc, err := S3ControlService(ctx, d, region)
	if err != nil {
		return nil, err
	}
	accessPointName := accessPointName(h.Item)

	// Build params
	params := &s3control.GetAccessPointPolicyInput{
		Name:      aws.String(accessPointName),
		AccountId: aws.String(commonColumnData.AccountId),
	}

	// execute list call
	op, err := svc.GetAccessPointPolicy(params)
	if err != nil {
		if a, ok := err.(awserr.Error); ok {
			if a.Code() == "NoSuchAccessPointPolicy" {
				return nil, nil
			}
		}
		return nil, err
	}

	return op, nil
}

func getAccessPointArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	accessPointName := accessPointName(h.Item)
	region := d.KeyColumnQualString(matrixKeyRegion)

	// Get account details
	getCommonColumnsCached := plugin.HydrateFunc(getCommonColumns).WithCache()
	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)
	arn := "arn:" + commonColumnData.Partition + ":s3:" + region + ":" + commonColumnData.AccountId + ":accesspoint/" + accessPointName

	return arn, nil
}

func accessPointName(item interface{}) string {
	switch item := item.(type) {
	case *s3control.AccessPoint:
		return *item.Name
	case *s3control.GetAccessPointOutput:
		return *item.Name
	}
	return ""
}
