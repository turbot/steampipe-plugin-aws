package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsSSMAssociation(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ssm_association",
		Description: "AWS SSM Association",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("association_id"),
			ShouldIgnoreError: isNotFoundError([]string{"AssociationDoesNotExist", "ValidationException"}),
			Hydrate:           getAwsSSMAssociation,
		},
		List: &plugin.ListConfig{
			Hydrate: listAwsSSMAssociations,
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "association_id",
				Description: "The ID created by the system when you create an association.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "association_name",
				Description: "The Name of association.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "document_name",
				Description: "The name of the Systems Manager document.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "date",
				Description: "The date when the association was made.",
				Hydrate:     getAwsSSMAssociation,
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "association_version",
				Description: "The association version.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "document_version",
				Description: "The version of the document used in the association.",
				Hydrate:     getAwsSSMAssociation,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "instance_id",
				Description: "The ID of the instance.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getAwsSSMAssociation,
			},
			{
				Name:        "last_execution_date",
				Description: "The date on which the association was last run.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "overview",
				Description: "Information about the association.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "schedule_expression",
				Description: "A cron expression that specifies a schedule when the association runs.",
				Hydrate:     getAwsSSMAssociation,
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "targets",
				Description: "A cron expression that specifies a schedule when the association runs.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "compliance_severity",
				Description: "A cron expression that specifies a schedule when the association runs.",
				Hydrate:     getAwsSSMAssociation,
				Type:        proto.ColumnType_STRING,
			},

			// Standard columns for all tables
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
				Hydrate:     getAwsSSMAssociationAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listAwsSSMAssociations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	plugin.Logger(ctx).Trace("listAwsSSMAssociations", "AWS_REGION", region)

	// Create session
	svc, err := SsmService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// List call
	err = svc.ListAssociationsPages(
		&ssm.ListAssociationsInput{},
		func(page *ssm.ListAssociationsOutput, isLast bool) bool {
			for _, association := range page.Associations {
				d.StreamListItem(ctx, association)
			}
			return !isLast
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getAwsSSMAssociation(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsSSMAssociation")

	// TODO put me in helper function
	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}
	var id string
	if h.Item != nil {
		id = associationID(h.Item)
	} else {
		id = d.KeyColumnQuals["association_id"].GetStringValue()
	}

	// Create Session
	svc, err := SsmService(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Build the params
	params := &ssm.DescribeAssociationInput{
		AssociationId: aws.String(id),
	}

	// Get call
	data, err := svc.DescribeAssociation(params)
	if err != nil {
		plugin.Logger(ctx).Debug("getAwsSSMAssociation", "ERROR", err)
		return nil, err
	}
	return data.AssociationDescription, nil
}

func getAwsSSMAssociationAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsSSMAssociationAkas")
	associationData := associationID(h.Item)
	c, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := c.(*awsCommonColumnData)
	aka := []string{"arn:" + commonColumnData.Partition + ":ssm:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":association/" + associationData}

	return aka, nil
}

func associationID(item interface{}) string {
	switch item.(type) {
	case *ssm.Association:
		return *item.(*ssm.Association).AssociationId
	case *ssm.AssociationDescription:
		return *item.(*ssm.AssociationDescription).AssociationId
	}
	return ""
}
