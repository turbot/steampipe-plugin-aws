package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3control"
	"github.com/aws/aws-sdk-go-v2/service/s3control/types"
	"github.com/aws/smithy-go"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsS3AccessPoint(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_s3_access_point",
		Description: "AWS S3 Access Point",
		List: &plugin.ListConfig{
			Hydrate: listS3AccessPoints,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NoSuchAccessPoint", "InvalidParameter", "InvalidRequest"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{Name: "bucket_name", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "region"}),
			Hydrate:    getS3AccessPoint,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"NoSuchAccessPoint", "InvalidParameter", "InvalidRequest"}),
			},
		},
		GetMatrixItemFunc: BuildRegionList,
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
	// Get account details

	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_access_point.listS3AccessPoints", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	region := d.KeyColumnQualString(matrixKeyRegion)
	// Create Session
	svc, err := S3ControlClient(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_access_point.listS3AccessPoints", "client_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	maxItems := int32(100)

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			maxItems = limit
		}
	}

	// Minimum limit is 0 as per the doc https://docs.aws.amazon.com/AmazonS3/latest/API/API_control_ListAccessPoints.html#API_control_ListAccessPoints_RequestSyntax
	input := &s3control.ListAccessPointsInput{
		AccountId:  aws.String(commonColumnData.AccountId),
		MaxResults: maxItems,
	}

	equalQuals := d.KeyColumnQuals
	if equalQuals["bucket_name"] != nil {
		if equalQuals["bucket_name"].GetStringValue() != "" {
			input.Bucket = aws.String(equalQuals["bucket_name"].GetStringValue())
		}
	}

	paginator := s3control.NewListAccessPointsPaginator(svc, input, func(o *s3control.ListAccessPointsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_s3_access_point.listS3AccessPoints", "api_error", err)
			return nil, err
		}

		for _, accessPoint := range output.AccessPointList {
			d.StreamListItem(ctx, accessPoint)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getS3AccessPoint(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	matrixRegion := d.KeyColumnQualString(matrixKeyRegion)

	// Get account details

	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_access_point.getS3AccessPoint", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Create Session
	svc, err := S3ControlClient(ctx, d, matrixRegion)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_access_point.getS3AccessPoint", "client_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	var name, region string
	if h.Item != nil {
		name = *h.Item.(types.AccessPoint).Name
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
	item, err := svc.GetAccessPoint(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_access_point.getS3AccessPoint", "api_error", err)
		return nil, err
	}

	return item, nil
}

func getS3AccessPointPolicyStatus(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)

	// Get account details

	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_access_point.getS3AccessPointPolicyStatus", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Create Session
	svc, err := S3ControlClient(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_access_point.getS3AccessPointPolicyStatus", "client_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	accessPointName := accessPointName(h.Item)

	// Build params
	params := &s3control.GetAccessPointPolicyStatusInput{
		Name:      aws.String(accessPointName),
		AccountId: aws.String(commonColumnData.AccountId),
	}

	// execute list call
	op, err := svc.GetAccessPointPolicyStatus(ctx, params)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "NoSuchAccessPointPolicy" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_s3_access_point.getS3AccessPointPolicyStatus", "api_error", err)
		return nil, err
	}

	return op, nil
}

func getS3AccessPointPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)

	// Get account details

	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_access_point.getS3AccessPointPolicy", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Create Session
	svc, err := S3ControlClient(ctx, d, region)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_access_point.getS3AccessPointPolicy", "client_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region check
		return nil, nil
	}

	accessPointName := accessPointName(h.Item)

	// Build params
	params := &s3control.GetAccessPointPolicyInput{
		Name:      aws.String(accessPointName),
		AccountId: aws.String(commonColumnData.AccountId),
	}

	// execute list call
	op, err := svc.GetAccessPointPolicy(ctx, params)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "NoSuchAccessPointPolicy" {
				return nil, nil
			}
		}
		plugin.Logger(ctx).Error("aws_s3_access_point.getS3AccessPointPolicy", "api_error", err)
		return nil, err
	}

	return op, nil
}

func getAccessPointArn(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	accessPointName := accessPointName(h.Item)
	region := d.KeyColumnQualString(matrixKeyRegion)

	// Get account details

	commonData, err := getCommonColumnsCached(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_s3_access_point.getAccessPointArn", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)
	arn := "arn:" + commonColumnData.Partition + ":s3:" + region + ":" + commonColumnData.AccountId + ":accesspoint/" + accessPointName

	return arn, nil
}

func accessPointName(item interface{}) string {
	switch item := item.(type) {
	case types.AccessPoint:
		return *item.Name
	case *s3control.GetAccessPointOutput:
		return *item.Name
	}
	return ""
}
